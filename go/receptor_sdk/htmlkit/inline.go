// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package htmlkit

import (
	stdhtml "html"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// PartRef returns the part name a src refers to, if it is a part reference.
func PartRef(src string) (name string, ok bool) {
	src = strings.TrimSpace(src)
	if len(src) <= len(PartScheme) || !strings.EqualFold(src[:len(PartScheme)], PartScheme) {
		return "", false
	}
	return src[len(PartScheme):], true
}

// PartLookup finds a part by name: its mime and bytes, or ok false.
type PartLookup func(name string) (mime string, body []byte, ok bool)

// InlineParts turns every src="tr-part:<name>" of a page back into a data
// URI with the part's bytes, so the page stands on its own. A reference
// whose part is not found is left as it is. It returns the page and the
// names resolved, in page order without repeats.
func InlineParts(htmlBody string, lookup PartLookup) (string, []string) {
	var out strings.Builder
	out.Grow(len(htmlBody))
	var resolved []string
	seen := map[string]bool{}
	z := html.NewTokenizer(strings.NewReader(htmlBody))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err() != io.EOF {
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
		attrs := make([]html.Attribute, 0, 4)
		for hasAttr {
			var key, val []byte
			key, val, hasAttr = z.TagAttr()
			attrs = append(attrs, html.Attribute{Key: string(key), Val: string(val)})
		}
		rewritten := false
		for i, a := range attrs {
			if !strings.EqualFold(a.Key, "src") {
				continue
			}
			partName, isRef := PartRef(a.Val)
			if !isRef {
				break
			}
			mime, body, ok := lookup(partName)
			if !ok {
				break
			}
			attrs[i].Val = DataURI(mime, body)
			rewritten = true
			if !seen[partName] {
				seen[partName] = true
				resolved = append(resolved, partName)
			}
			break
		}
		if !rewritten {
			out.Write(raw)
			continue
		}
		out.WriteString("<img")
		for _, a := range attrs {
			out.WriteByte(' ')
			out.WriteString(a.Key)
			out.WriteString(`="`)
			out.WriteString(stdhtml.EscapeString(a.Val))
			out.WriteByte('"')
		}
		if tt == html.SelfClosingTagToken {
			out.WriteString("/>")
		} else {
			out.WriteByte('>')
		}
	}
	return out.String(), resolved
}
