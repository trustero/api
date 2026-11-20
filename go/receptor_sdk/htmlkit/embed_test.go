// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package htmlkit

import (
	"bytes"
	"image"
	"image/png"
	"strings"
	"testing"

	"github.com/trustero/api/go/receptor_sdk"
)

// pngBytes encodes a square png; random pixels so a 200px picture clears
// the part size floor (DefaultMinPartBytes) while a 16px icon does not.
func pngBytes(t *testing.T, side int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, side, side))
	x := uint32(1)
	for i := range img.Pix {
		x = x*1664525 + 1013904223
		img.Pix[i] = byte(x >> 24)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

type fakeFetch struct {
	images map[string]*Image
	calls  []Tag
}

func (f *fakeFetch) fetch(tag Tag) Result {
	f.calls = append(f.calls, tag)
	if img, ok := f.images[tag.Src]; ok {
		copy := *img
		return Result{Image: &copy}
	}
	return Result{}
}

const page = `<!DOCTYPE html>
<!-- a comment with <img src="x"> -->
<html><head><title>T</title><script>var s = '<img src="in-script">';</script></head>
<body><h1>Page</h1>
<p>Big: <img src="https://wiki/big.png" alt="Architecture &amp; flows" class="confluence-embedded-image" data-linked-resource-id="att1"></p>
<p>Small: <IMG SRC="https://wiki/icon.png" width="16" height="16"/></p>
<p>Missing: <img src="https://wiki/missing.png"></p>
<p>Again: <img src="https://wiki/big.png"></p>
</body></html>`

// A page nothing is embedded in comes back byte for byte, including the
// comment and the script that look like image tags.
func TestEmbedImagesLeavesAPageWithoutImagesUntouched(t *testing.T) {
	f := &fakeFetch{}
	out := EmbedImages(page, f.fetch, Options{})
	if out.HTML != page {
		t.Errorf("page changed:\n%s", out.HTML)
	}
	if len(f.calls) != 3 {
		t.Errorf("want one fetch per distinct src (3), got %d", len(f.calls))
	}
	if out.Embedded != 0 || len(out.Parts) != 0 {
		t.Errorf("embedded %d, parts %d", out.Embedded, len(out.Parts))
	}
}

func TestEmbedImagesDataURIMode(t *testing.T) {
	big, icon := pngBytes(t, 200), pngBytes(t, 16)
	f := &fakeFetch{images: map[string]*Image{"https://wiki/big.png": {Body: big}, "https://wiki/icon.png": {Body: icon}}}
	out := EmbedImages(page, f.fetch, Options{})
	if len(out.Parts) != 0 {
		t.Errorf("data URI mode made %d parts", len(out.Parts))
	}
	// The big picture is inlined twice, so it counts twice against the page.
	if out.Embedded != 2 || out.Bytes != 2*len(big)+len(icon) {
		t.Errorf("embedded %d images, %d bytes", out.Embedded, out.Bytes)
	}
	want := `<img src="` + DataURI("image/png", big) + `" alt="Architecture &amp; flows">`
	if !strings.Contains(out.HTML, want) {
		t.Errorf("big image not inlined with its alt and without provider attributes:\n%s", out.HTML)
	}
	if strings.Contains(out.HTML, "data-linked-resource-id") || strings.Contains(out.HTML, "class=") {
		t.Errorf("provider attributes kept:\n%s", out.HTML)
	}
	if !strings.Contains(out.HTML, `<img src="`+DataURI("image/png", icon)+`" width="16" height="16">`) {
		t.Errorf("icon not inlined with its size:\n%s", out.HTML)
	}
	if strings.Count(out.HTML, DataURI("image/png", big)) != 2 {
		t.Errorf("the second use of the big image was not inlined from the first fetch")
	}
	if !strings.Contains(out.HTML, `<img src="https://wiki/missing.png">`) {
		t.Errorf("a failed download did not leave the tag alone")
	}
	if strings.Count(out.HTML, "<script>var s = '<img src=\"in-script\">';</script>") != 1 {
		t.Errorf("script text was touched")
	}
}

func TestEmbedImagesPartsMode(t *testing.T) {
	big, icon := pngBytes(t, 200), pngBytes(t, 16)
	f := &fakeFetch{images: map[string]*Image{
		"https://wiki/big.png":  {Body: big, Name: "diagram.png", Locator: map[string]string{"attachment_id": "att1"}},
		"https://wiki/icon.png": {Body: icon},
	}}
	out := EmbedImages(page, f.fetch, Options{Mode: ModeParts})
	if len(out.Parts) != 1 {
		t.Fatalf("want the big image as the one part, got %d: %+v", len(out.Parts), out.Parts)
	}
	p := out.Parts[0]
	if p.Name != "diagram.png" || p.Mime != "image/png" || p.Kind != receptor_sdk.PartKindImage || !bytes.Equal(p.Body, big) || p.Locator["attachment_id"] != "att1" {
		t.Errorf("part: %+v", p)
	}
	if strings.Count(out.HTML, `src="tr-part:diagram.png"`) != 2 {
		t.Errorf("both uses of the big image should refer to the one part:\n%s", out.HTML)
	}
	if !strings.Contains(out.HTML, DataURI("image/png", icon)) {
		t.Errorf("the icon under the floor should be inlined, not a part")
	}
	if out.Embedded != 2 {
		t.Errorf("embedded %d", out.Embedded)
	}
}

// An image already inlined becomes a part without a fetch; names that
// clash are numbered; an image without a name is numbered by its part.
func TestEmbedImagesPartsFromDataURIsAndNames(t *testing.T) {
	big := pngBytes(t, 200)
	other := pngBytes(t, 210)
	html := `<img src="` + DataURI("image/png", big) + `" alt="a"><img src="https://x/1"><img src="https://x/2"><img src="https://x/3"><img src="data:image/gif;base64,R0lGOD">`
	f := &fakeFetch{images: map[string]*Image{
		"https://x/1": {Body: other, Name: "../pics/chart.png"},
		"https://x/2": {Body: pngBytes(t, 220), Name: "chart.png"},
		"https://x/3": {Body: pngBytes(t, 230)},
	}}
	out := EmbedImages(html, f.fetch, Options{Mode: ModeParts})
	if len(f.calls) != 3 {
		t.Errorf("the data URI image should not be fetched; %d fetches", len(f.calls))
	}
	var names []string
	for _, p := range out.Parts {
		names = append(names, p.Name)
	}
	if got := strings.Join(names, ","); got != "image-1.png,chart.png,chart-2.png,image-4.png" {
		t.Errorf("part names: %s", got)
	}
	if !strings.HasPrefix(out.HTML, `<img src="tr-part:image-1.png" alt="a">`) {
		t.Errorf("data URI image not rewritten to its part:\n%s", out.HTML)
	}
	if !strings.HasSuffix(out.HTML, `<img src="data:image/gif;base64,R0lGOD">`) {
		t.Errorf("an inlined icon under the floor should stay byte for byte:\n%s", out.HTML)
	}
}

func TestEmbedImagesFallbackTextAndSkippedSchemes(t *testing.T) {
	calls := 0
	fetch := func(tag Tag) Result {
		calls++
		return Result{FallbackText: "<:smile:>"}
	}
	html := `<p><img src="cid:abc"><img src="javascript:x"><img src="VBScript:y"><img src="tr-part:z.png"><img src="//wiki/rel.png"><img src="/wiki/path.png"><img src="https://wiki/emoji.png" data-emoji-fallback=":smile:"></p>`
	out := EmbedImages(html, fetch, Options{})
	if calls != 3 {
		t.Errorf("only web and relative sources are fetched; %d fetches", calls)
	}
	if out.HTML != `<p><img src="cid:abc"><img src="javascript:x"><img src="VBScript:y"><img src="tr-part:z.png">&lt;:smile:&gt;&lt;:smile:&gt;&lt;:smile:&gt;</p>` {
		t.Errorf("got %s", out.HTML)
	}
	for src, want := range map[string]string{"https://x/a": "https", "HTTP://x": "http", "/wiki/a.png": "", "a.png": "", "mailto:x": "mailto", ":x": "", "data:image/png;base64,AA": "data", "c:\\x": "c"} {
		if got := srcScheme(src); got != want {
			t.Errorf("srcScheme(%q) = %q, want %q", src, got, want)
		}
	}
}

func TestEmbedImagesBudgets(t *testing.T) {
	big := pngBytes(t, 200)
	html := `<img src="https://x/1"><img src="https://x/2"><img src="https://x/3">`
	f := &fakeFetch{images: map[string]*Image{"https://x/1": {Body: big}, "https://x/2": {Body: big}, "https://x/3": {Body: big}}}

	out := EmbedImages(html, f.fetch, Options{MaxImages: 2})
	if out.Embedded != 2 || !strings.HasSuffix(out.HTML, `<img src="https://x/3">`) {
		t.Errorf("MaxImages: embedded %d:\n%s", out.Embedded, out.HTML)
	}

	f.calls = nil
	out = EmbedImages(html, f.fetch, Options{MaxPageBytes: len(big) + 10})
	if out.Embedded != 1 {
		t.Errorf("MaxPageBytes: embedded %d", out.Embedded)
	}
	if len(f.calls) < 2 || f.calls[1].MaxBytes != 10 {
		t.Errorf("the second fetch should be told the 10 bytes left, got %+v", f.calls)
	}

	out = EmbedImages(html, f.fetch, Options{MaxImageBytes: len(big) - 1})
	if out.Embedded != 0 || out.HTML != html {
		t.Errorf("MaxImageBytes: embedded %d:\n%s", out.Embedded, out.HTML)
	}

	notImage := &fakeFetch{images: map[string]*Image{"https://x/1": {Body: []byte("%PDF-1.4 not a picture")}}}
	out = EmbedImages(`<img src="https://x/1">`, notImage.fetch, Options{})
	if out.Embedded != 0 {
		t.Errorf("a download that is not an image was embedded")
	}
}

func TestInlineParts(t *testing.T) {
	big := pngBytes(t, 200)
	html := `<p><img src="tr-part:diagram.png" alt="A &amp; B"/><img src="TR-PART:missing.png"><img src="https://x/1"></p>`
	lookup := func(name string) (string, []byte, bool) {
		if name == "diagram.png" {
			return "image/png", big, true
		}
		return "", nil, false
	}
	out, resolved := InlineParts(html, lookup)
	want := `<p><img src="` + DataURI("image/png", big) + `" alt="A &amp; B"/><img src="TR-PART:missing.png"><img src="https://x/1"></p>`
	if out != want {
		t.Errorf("got %s", out)
	}
	if len(resolved) != 1 || resolved[0] != "diagram.png" {
		t.Errorf("resolved %v", resolved)
	}
	if name, ok := PartRef(" tr-part:x.png "); !ok || name != "x.png" {
		t.Errorf("PartRef: %q %v", name, ok)
	}
	if _, ok := PartRef("tr-part:"); ok {
		t.Error("an empty reference parsed")
	}
}

// A part name from a provider file name may hold a quote; it is kept out of
// the attribute and the header, and the reference still resolves.
func TestEmbedImagesEscapesPartReferences(t *testing.T) {
	big := pngBytes(t, 200)
	f := &fakeFetch{images: map[string]*Image{"https://x/1": {Body: big, Name: `say "hi" <now>.png`}}}
	out := EmbedImages(`<img src="https://x/1">`, f.fetch, Options{Mode: ModeParts})
	if len(out.Parts) != 1 || out.Parts[0].Name != "say _hi_ _now_.png" {
		t.Fatalf("parts %+v", out.Parts)
	}
	if out.HTML != `<img src="tr-part:say _hi_ _now_.png">` {
		t.Errorf("got %s", out.HTML)
	}
	resolved, names := InlineParts(out.HTML, func(name string) (string, []byte, bool) {
		if name == "say _hi_ _now_.png" {
			return "image/png", big, true
		}
		return "", nil, false
	})
	if len(names) != 1 || !strings.Contains(resolved, DataURI("image/png", big)) {
		t.Errorf("reference did not resolve: %v\n%s", names, resolved)
	}
}

// Attribute values arrive unescaped exactly once, an upper-case data: scheme
// is read like a lower-case one, and a mime that is not a plain token is
// not made a part.
func TestEmbedImagesAttributeDecodingAndDataURIs(t *testing.T) {
	big := pngBytes(t, 200)
	var seen []string
	fetch := func(tag Tag) Result {
		seen = append(seen, tag.Src)
		return Result{Image: &Image{Body: big}}
	}
	html := `<img src="https://x/a.png?v=1&amp;notify=0" alt="A &amp;amp; B">` +
		`<img src="DATA:image/png;base64,` + strings.TrimPrefix(DataURI("image/png", big), "data:image/png;base64,") + `">` +
		`<img src="data:image/png&#10;X-Injected: 1;base64,AAAA">`
	out := EmbedImages(html, fetch, Options{Mode: ModeParts})
	if len(seen) != 1 || seen[0] != "https://x/a.png?v=1&notify=0" {
		t.Errorf("fetched src decoded as %q", seen)
	}
	if !strings.Contains(out.HTML, `alt="A &amp;amp; B"`) {
		t.Errorf("alt re-escaped wrongly:\n%s", out.HTML)
	}
	if len(out.Parts) != 2 {
		t.Fatalf("want the fetched picture and the DATA: one as parts, got %d", len(out.Parts))
	}
	if !strings.HasSuffix(out.HTML, `<img src="data:image/png&#10;X-Injected: 1;base64,AAAA">`) {
		t.Errorf("a data URI with a bad mime should be left alone:\n%s", out.HTML)
	}
}
