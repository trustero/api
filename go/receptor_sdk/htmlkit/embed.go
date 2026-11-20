// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

// Package htmlkit embeds the images of an HTML page a receptor collects.
//
// Pages exported from a provider (a Confluence page, a Notion page) refer to
// their images by URLs that only work with the provider's credentials, or
// that expire. A receptor downloads each image at scan time so the stored
// page is self-contained. EmbedImages does the walking, the budgets and the
// rewriting; the receptor supplies the download as a Fetcher, since only it
// knows the provider's URLs and authentication.
//
// Two output modes:
//
//   - ModeDataURI (the default) inlines every image as a data URI. This is
//     what receptors have always produced.
//   - ModeParts rewrites each image to src="tr-part:<name>" and returns the
//     images as receptor_sdk.DocumentPart values for the Document's Parts.
//     Trustero stores the page and its images together, lists and indexes
//     each image on its own, and resolves the references when it stores the
//     page, so the page still renders. Pictures too small to be evidence
//     (icons, emoji) are inlined as data URIs instead of becoming parts.
package htmlkit

import (
	"bytes"
	"encoding/base64"
	"fmt"
	stdhtml "html"
	"image"
	_ "image/gif"  // image header decoders for the part size floor
	_ "image/jpeg" //
	_ "image/png"  //
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/trustero/api/go/receptor_sdk"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Mode selects how an embedded image is written into the page.
type Mode int

const (
	// ModeDataURI inlines every image as a data URI.
	ModeDataURI Mode = iota
	// ModeParts refers to each image as a document part (src="tr-part:<name>")
	// and returns the parts; images under the part size floor are inlined.
	ModeParts
)

// PartScheme is the URI scheme a page uses to refer to one of its parts by
// name: src="tr-part:diagram.png".
const PartScheme = "tr-part:"

// Defaults for Options. They match what the Confluence and Notion receptors
// have used for pages, and Trustero's own floor for pictures worth keeping
// as parts.
const (
	DefaultMaxImages        = 40
	DefaultMaxImageBytes    = 12 * 1024 * 1024
	DefaultMaxPageBytes     = 40 * 1024 * 1024
	DefaultMinPartBytes     = 8 * 1024
	DefaultMinPartDimension = 150
)

// Tag is one <img> of the page, as the Fetcher sees it.
type Tag struct {
	// Index is the tag's position among the page's <img> tags, from 0.
	Index int
	// Src is the tag's src attribute, unescaped and trimmed.
	Src string
	// Attrs holds every attribute of the tag, keys lower-cased, values
	// unescaped and trimmed.
	Attrs map[string]string
	// MaxBytes is the most the image may be to fit the budgets left: the
	// per-image cap, or what is left of the page's. A Fetcher can stop a
	// download that grows past it.
	MaxBytes int
}

// Image is a downloaded image.
type Image struct {
	Body []byte
	// Mime is the image's type; sniffed from Body when empty.
	Mime string
	// Name is the part's name and file name in ModeParts; "image-<n>.<ext>"
	// when empty. Unused in ModeDataURI.
	Name string
	// Locator says where the image came from at the provider, for the part.
	Locator map[string]string
}

// Result is what a Fetcher makes of a tag.
type Result struct {
	// Image is the downloaded image, or nil when the tag is not embedded.
	Image *Image
	// FallbackText, when Image is nil, replaces the whole tag with text
	// (escaped). Empty leaves the tag as it is.
	FallbackText string
}

// Fetcher downloads the image of one tag. It is called once per distinct
// src; a tag whose src was already fetched reuses the result.
type Fetcher func(tag Tag) Result

// Options bound and shape the embedding. Zero values take the defaults.
type Options struct {
	Mode Mode
	// MaxImages caps how many images are embedded.
	MaxImages int
	// MaxImageBytes caps one image; a larger download is not embedded.
	MaxImageBytes int
	// MaxPageBytes caps the images of the page together.
	MaxPageBytes int
	// MinPartBytes and MinPartDimension are the floor under which an image
	// is inlined rather than made a part in ModeParts: smaller than
	// MinPartBytes, or shorter than MinPartDimension on both sides.
	// Negative values disable the floor.
	MinPartBytes     int
	MinPartDimension int
}

func (o Options) withDefaults() Options {
	if o.MaxImages == 0 {
		o.MaxImages = DefaultMaxImages
	}
	if o.MaxImageBytes == 0 {
		o.MaxImageBytes = DefaultMaxImageBytes
	}
	if o.MaxPageBytes == 0 {
		o.MaxPageBytes = DefaultMaxPageBytes
	}
	if o.MinPartBytes == 0 {
		o.MinPartBytes = DefaultMinPartBytes
	}
	if o.MinPartDimension == 0 {
		o.MinPartDimension = DefaultMinPartDimension
	}
	return o
}

// Output is the rewritten page and what was embedded.
type Output struct {
	HTML string
	// Parts are the images written as parts, in page order (ModeParts only).
	Parts []receptor_sdk.DocumentPart
	// Embedded counts the images embedded, inline or as parts.
	Embedded int
	// Bytes is the size of the embedded images together.
	Bytes int
}

// keptAttrs are the <img> attributes copied to the rewritten tag. Provider
// attributes (data-*, class) are dropped on purpose: they carry URLs that
// only work with the provider's credentials, and the page should not keep
// references to them.
var keptAttrs = []string{"alt", "title", "width", "height"}

// EmbedImages rewrites the <img> tags of htmlBody. Every tag is offered to
// fetch except those whose src is a cid:, javascript: or tr-part: reference
// and those already inlined (data:); in ModeParts an inlined picture above
// the part size floor is made a part without a fetch. The rest of the page
// is copied byte for byte.
func EmbedImages(htmlBody string, fetch Fetcher, opts Options) Output {
	opts = opts.withDefaults()
	e := &embedder{opts: opts, fetch: fetch, bySrc: map[string]*fetched{}, names: map[string]bool{}}
	var out strings.Builder
	out.Grow(len(htmlBody))
	z := html.NewTokenizer(strings.NewReader(htmlBody))
	index := 0
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err() != io.EOF {
				// The tokenizer gives up only on a reader error, which a
				// string never has; keep the page whole regardless.
				out.Write(z.Raw())
			}
			break
		}
		raw := z.Raw()
		if tt != html.StartTagToken && tt != html.SelfClosingTagToken {
			out.Write(raw)
			continue
		}
		name, hasAttr := z.TagName()
		if atom.Lookup(name) != atom.Img {
			out.Write(raw)
			continue
		}
		tag := Tag{Index: index, Attrs: map[string]string{}}
		index++
		for hasAttr {
			// TagAttr hands back the key lower-cased and the value unescaped.
			var key, val []byte
			key, val, hasAttr = z.TagAttr()
			tag.Attrs[string(key)] = strings.TrimSpace(string(val))
		}
		tag.Src = tag.Attrs["src"]
		if replacement, ok := e.rewrite(tag); ok {
			out.WriteString(replacement)
		} else {
			out.Write(raw)
		}
	}
	return Output{HTML: out.String(), Parts: e.parts, Embedded: e.count, Bytes: e.bytes}
}

type fetched struct {
	img      *Image
	fallback string
	// src the image was written with, once embedded: a data URI or a
	// tr-part reference, so the same picture twice on a page is one part.
	src string
}

type embedder struct {
	opts  Options
	fetch Fetcher
	bySrc map[string]*fetched
	names map[string]bool
	parts []receptor_sdk.DocumentPart
	count int
	bytes int
}

// rewrite returns the tag's replacement, or false to leave the tag as it is.
func (e *embedder) rewrite(tag Tag) (string, bool) {
	lower := strings.ToLower(tag.Src)
	if strings.HasPrefix(lower, "cid:") || strings.HasPrefix(lower, "javascript:") || strings.HasPrefix(lower, PartScheme) {
		return "", false
	}
	if strings.HasPrefix(lower, "data:") {
		// Already inlined. In parts mode a picture worth a part becomes one;
		// anything else stays as it is, byte for byte.
		if e.opts.Mode != ModeParts {
			return "", false
		}
		if f, ok := e.bySrc[tag.Src]; ok {
			return e.replacement(tag, f)
		}
		img := decodeDataURI(tag.Src)
		if img == nil || e.belowFloor(img) || e.count >= e.opts.MaxImages || len(img.Body) > e.opts.MaxImageBytes || len(img.Body) > e.opts.MaxPageBytes-e.bytes {
			return "", false
		}
		e.count++
		e.bytes += len(img.Body)
		f := &fetched{img: img, src: e.embedPart(img)}
		e.bySrc[tag.Src] = f
		return e.replacement(tag, f)
	}
	if tag.Src != "" {
		if f, ok := e.bySrc[tag.Src]; ok {
			if f.img != nil && e.opts.Mode == ModeDataURI {
				// A repeated picture is inlined again, so it costs the page
				// budget again; past the budget the tag is left as it is.
				if e.bytes >= e.opts.MaxPageBytes {
					return "", false
				}
				e.bytes += len(f.img.Body)
			}
			return e.replacement(tag, f)
		}
	}
	if e.count >= e.opts.MaxImages || e.bytes >= e.opts.MaxPageBytes {
		return "", false
	}
	f := &fetched{}
	if e.fetch != nil {
		tag.MaxBytes = e.opts.MaxImageBytes
		if remaining := e.opts.MaxPageBytes - e.bytes; remaining < tag.MaxBytes {
			tag.MaxBytes = remaining
		}
		r := e.fetch(tag)
		f.img, f.fallback = r.Image, r.FallbackText
	}
	if f.img != nil {
		if len(f.img.Body) == 0 || len(f.img.Body) > e.opts.MaxImageBytes || len(f.img.Body) > e.opts.MaxPageBytes-e.bytes {
			f.img = nil
		}
	}
	if f.img != nil {
		if f.img.Mime == "" {
			f.img.Mime = http.DetectContentType(f.img.Body)
		}
		if !strings.HasPrefix(strings.ToLower(f.img.Mime), "image/") {
			f.img = nil
		}
	}
	if f.img != nil {
		e.count++
		e.bytes += len(f.img.Body)
		f.src = e.embed(f.img)
	}
	if tag.Src != "" {
		e.bySrc[tag.Src] = f
	}
	return e.replacement(tag, f)
}

// embed writes the image as a part or a data URI and returns the src to use.
func (e *embedder) embed(img *Image) string {
	if e.opts.Mode == ModeParts && !e.belowFloor(img) {
		return e.embedPart(img)
	}
	return DataURI(img.Mime, img.Body)
}

// embedPart records the image as a part and returns its reference.
func (e *embedder) embedPart(img *Image) string {
	name := e.uniqueName(img)
	e.parts = append(e.parts, receptor_sdk.DocumentPart{
		Name:    name,
		Mime:    img.Mime,
		Kind:    receptor_sdk.PartKindImage,
		Body:    img.Body,
		Locator: img.Locator,
	})
	return PartScheme + name
}

func (e *embedder) replacement(tag Tag, f *fetched) (string, bool) {
	if f.img == nil {
		if f.fallback != "" {
			return stdhtml.EscapeString(f.fallback), true
		}
		return "", false
	}
	var b strings.Builder
	b.WriteString(`<img src="`)
	b.WriteString(stdhtml.EscapeString(f.src))
	b.WriteByte('"')
	for _, key := range keptAttrs {
		if val, ok := tag.Attrs[key]; ok && val != "" {
			fmt.Fprintf(&b, ` %s="%s"`, key, stdhtml.EscapeString(val))
		}
	}
	b.WriteByte('>')
	return b.String(), true
}

// belowFloor reports whether an image is too small to be a part: under
// MinPartBytes, or under MinPartDimension on both sides. Dimensions come
// from the header of the formats Go decodes; others are judged by size.
func (e *embedder) belowFloor(img *Image) bool {
	if e.opts.MinPartBytes < 0 && e.opts.MinPartDimension < 0 {
		return false
	}
	if e.opts.MinPartBytes > 0 && len(img.Body) < e.opts.MinPartBytes {
		return true
	}
	if e.opts.MinPartDimension > 0 {
		if cfg, _, err := image.DecodeConfig(bytes.NewReader(img.Body)); err == nil {
			return cfg.Width < e.opts.MinPartDimension && cfg.Height < e.opts.MinPartDimension
		}
	}
	return false
}

// uniqueName gives the part a name no other part of the page has: the
// image's own, else image-<n>.<ext>, with a counter before the extension on
// a clash.
func (e *embedder) uniqueName(img *Image) string {
	name := strings.ToValidUTF8(strings.TrimSpace(img.Name), "�")
	name = path.Base(strings.ReplaceAll(name, "\\", "/"))
	// The name goes into an attribute and a MIME header; quotes, angle
	// brackets and control characters have no place in either.
	name = strings.Map(func(r rune) rune {
		if r == '"' || r == '<' || r == '>' || r < 0x20 || r == 0x7f {
			return '_'
		}
		return r
	}, name)
	if name == "." || name == "/" || name == "" {
		name = fmt.Sprintf("image-%d%s", len(e.parts)+1, extensionFor(img.Mime))
	}
	candidate := name
	for n := 2; e.names[candidate]; n++ {
		ext := path.Ext(name)
		candidate = fmt.Sprintf("%s-%d%s", strings.TrimSuffix(name, ext), n, ext)
	}
	e.names[candidate] = true
	return candidate
}

// DataURI inlines bytes as a data URI.
func DataURI(mime string, body []byte) string {
	return fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(body))
}

// decodeDataURI reads an inlined image back; nil when the URI is not a
// base64 image.
func decodeDataURI(src string) *Image {
	header, data, ok := strings.Cut(src[len("data:"):], ",")
	if !ok {
		return nil
	}
	mime, params, _ := strings.Cut(header, ";")
	mime = strings.ToLower(strings.TrimSpace(mime))
	// The mime becomes a part's Content-Type header, so only a plain
	// type/subtype token is accepted; anything else came from page content.
	if !imageMimeToken.MatchString(mime) || !strings.Contains(strings.ToLower(params), "base64") {
		return nil
	}
	body, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		if body, err = base64.RawStdEncoding.DecodeString(data); err != nil {
			return nil
		}
	}
	return &Image{Body: body, Mime: mime}
}

// imageMimeToken is the shape of an image mime a page may name.
var imageMimeToken = regexp.MustCompile(`^image/[a-z0-9][a-z0-9.+-]*$`)

var extensions = map[string]string{
	"image/png":     ".png",
	"image/jpeg":    ".jpg",
	"image/gif":     ".gif",
	"image/webp":    ".webp",
	"image/svg+xml": ".svg",
	"image/bmp":     ".bmp",
	"image/tiff":    ".tiff",
}

func extensionFor(mime string) string {
	mime, _, _ = strings.Cut(strings.ToLower(mime), ";")
	if ext, ok := extensions[strings.TrimSpace(mime)]; ok {
		return ext
	}
	return ""
}
