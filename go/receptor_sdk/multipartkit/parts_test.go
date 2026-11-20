// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package multipartkit

import (
	"bytes"
	"io"
	"mime"
	"os"
	"path/filepath"
	"testing"
)

func readParts(t *testing.T, envelope []byte, boundary string) (infos []PartInfo, oks []bool, bodies [][]byte) {
	t.Helper()
	r, err := NewMultipartReader(bytes.NewReader(envelope), boundary, DefaultBufferSize)
	if err != nil {
		t.Fatal(err)
	}
	for {
		p, err := r.NextPart()
		if err == io.EOF {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(p)
		info, ok := ParsePartInfo(p.Header)
		infos, oks, bodies = append(infos, info), append(oks, ok), append(bodies, body)
	}
}

// A document part written from bytes or a file reads back with its name,
// parent, mime, kind and locator; file and protobuf parts are not parts.
func TestDocumentPartRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	b, _ := NewMultipartBuilder(&buf, 0)
	if err := b.AddBytes("page.html", "page.html", "text/html", []byte("<p>hi</p>"), nil); err != nil {
		t.Fatal(err)
	}
	locator := map[string]string{"page_id": "123", "attachment_id": "att1"}
	if err := b.AddPartBytes(PartInfo{Name: "diagram.png", Parent: "page.html", ContentType: "image/png", Kind: "image", Locator: locator}, []byte("PNG")); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "photo.jpg")
	if err := os.WriteFile(path, []byte("JPG"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := b.AddPartFile(PartInfo{Name: "Übersicht ä.jpg", Parent: "page.html", ContentType: "image/jpeg", Kind: "image"}, path); err != nil {
		t.Fatal(err)
	}
	if err := b.Finalize(); err != nil {
		t.Fatal(err)
	}

	infos, oks, bodies := readParts(t, buf.Bytes(), b.GetBoundary())
	if len(infos) != 3 {
		t.Fatalf("want 3 parts, got %d", len(infos))
	}
	if oks[0] {
		t.Errorf("the file part parsed as a document part: %+v", infos[0])
	}
	if !oks[1] || !oks[2] {
		t.Fatalf("document parts not recognised: %v", oks)
	}
	got := infos[1]
	if got.Name != "diagram.png" || got.Parent != "page.html" || got.ContentType != "image/png" || got.Kind != "image" {
		t.Errorf("part 1 read back as %+v", got)
	}
	if got.Locator["page_id"] != "123" || got.Locator["attachment_id"] != "att1" {
		t.Errorf("locator read back as %v", got.Locator)
	}
	if string(bodies[1]) != "PNG" || string(bodies[2]) != "JPG" {
		t.Errorf("bodies read back as %q %q", bodies[1], bodies[2])
	}
	if infos[2].Name != "Übersicht ä.jpg" {
		t.Errorf("non-ASCII name read back as %q", infos[2].Name)
	}
	if infos[2].Locator != nil {
		t.Errorf("a part without a locator read back with %v", infos[2].Locator)
	}
}

func TestDocumentPartNeedsAName(t *testing.T) {
	var buf bytes.Buffer
	b, _ := NewMultipartBuilder(&buf, 0)
	if err := b.AddPartBytes(PartInfo{Parent: "page.html"}, []byte("x")); err == nil {
		t.Error("a nameless part was accepted")
	}
}

// A file name with a quote still reads back, and a header value from page
// content cannot add header lines.
func TestFilePartNamesAndHeadersStaySafe(t *testing.T) {
	var buf bytes.Buffer
	b, _ := NewMultipartBuilder(&buf, 0)
	if err := b.AddBytes(`say "hi".html`, `say "hi".html`, "text/html", []byte("<p>"), nil); err != nil {
		t.Fatal(err)
	}
	if err := b.AddPartBytes(PartInfo{Name: "x.png", Parent: `say "hi".html`, ContentType: "image/png\r\nX-Injected: 1", Kind: "image\n"}, []byte("PNG")); err != nil {
		t.Fatal(err)
	}
	if err := b.Finalize(); err != nil {
		t.Fatal(err)
	}
	r, _ := NewMultipartReader(bytes.NewReader(buf.Bytes()), b.GetBoundary(), DefaultBufferSize)
	file, err := r.NextPart()
	if err != nil {
		t.Fatal(err)
	}
	if _, params, err := parseDisposition(file.Header.Get("Content-Disposition")); err != nil || params["name"] != `say "hi".html` {
		t.Errorf("file disposition %q read back as %v (%v)", file.Header.Get("Content-Disposition"), params, err)
	}
	part, err := r.NextPart()
	if err != nil {
		t.Fatal(err)
	}
	info, ok := ParsePartInfo(part.Header)
	if !ok || info.Parent != `say "hi".html` || info.ContentType != "image/pngX-Injected:1" || info.Kind != "image" {
		t.Errorf("part read back as %+v", info)
	}
	if part.Header.Get("X-Injected") != "" {
		t.Error("a header was injected through the content type")
	}
}

func parseDisposition(v string) (string, map[string]string, error) { return mime.ParseMediaType(v) }
