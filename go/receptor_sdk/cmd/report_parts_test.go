// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package cmd

import (
	"bytes"
	"io"
	"mime"
	"os"
	"strings"
	"testing"

	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_sdk/multipartkit"
	"github.com/trustero/api/go/receptor_v1"
)

type envelopePart struct {
	disposition string
	params      map[string]string
	contentType string
	body        string
}

func readEnvelope(t *testing.T, path, contentType string) []envelopePart {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	boundary := strings.TrimSpace(strings.Split(contentType, "boundary=")[1])
	r, err := multipartkit.NewMultipartReader(f, boundary, multipartkit.DefaultBufferSize)
	if err != nil {
		t.Fatal(err)
	}
	var parts []envelopePart
	for {
		p, err := r.NextPart()
		if err == io.EOF {
			return parts
		}
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(p)
		disposition, params, _ := mime.ParseMediaType(p.Header.Get("Content-Disposition"))
		parts = append(parts, envelopePart{disposition: disposition, params: params, contentType: p.Header.Get("Content-Type"), body: string(body)})
	}
}

func findingWithDoc(doc *receptor_v1.Document) *receptor_v1.Finding {
	return &receptor_v1.Finding{
		ReceptorType: "trr-test",
		Evidences: []*receptor_v1.Evidence{{
			Caption:      "Page",
			EvidenceType: &receptor_v1.Evidence_Doc{Doc: doc},
		}},
	}
}

// A document without parts produces the envelope receptors have always
// sent: the finding, the file, the sources, and nothing else. Receptors
// built on this SDK against an older server depend on that.
func TestMultipartEvidenceWithoutPartsIsUnchanged(t *testing.T) {
	doc := &receptor_v1.Document{Mime: "text/html", Body: []byte("<p>hi</p>"), FileName: "page.html"}
	contentType, path, err := multipartEvidence(findingWithDoc(doc), nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)
	if !strings.HasPrefix(contentType, "multipart/tr-mixed; text/html; boundary=") {
		t.Errorf("content type %q", contentType)
	}
	parts := readEnvelope(t, path, contentType)
	if len(parts) != 3 {
		t.Fatalf("want finding, file, sources; got %d parts: %+v", len(parts), parts)
	}
	if parts[0].disposition != "protobuf" || parts[0].params["name"] != "receptor_v1.Finding" {
		t.Errorf("part 0: %+v", parts[0])
	}
	if parts[1].disposition != "file" || parts[1].params["name"] != "page.html" || parts[1].params["filename"] != "page.html" || parts[1].body != "<p>hi</p>" {
		t.Errorf("part 1: %+v", parts[1])
	}
	if parts[2].disposition != "protobuf" || parts[2].params["name"] != "receptor_v1.Sources" {
		t.Errorf("part 2: %+v", parts[2])
	}
}

// Parts follow the file part and precede the sources, named after the file
// part they belong to, whether the document is in memory or streamed.
func TestMultipartEvidenceWritesDocumentParts(t *testing.T) {
	dir := t.TempDir()
	imgPath := dir + "/photo.jpg"
	if err := os.WriteFile(imgPath, []byte("JPG"), 0o600); err != nil {
		t.Fatal(err)
	}
	sdkDoc := receptor_sdk.Document{
		Mime: "text/html", Body: []byte(`<img src="tr-part:diagram.png">`),
		Parts: []receptor_sdk.DocumentPart{
			{Name: "diagram.png", Mime: "image/png", Kind: receptor_sdk.PartKindImage, Body: []byte("PNG"), Locator: map[string]string{"page_id": "1"}},
			{Name: "photo.jpg", Mime: "image/jpeg", Kind: receptor_sdk.PartKindImage, StreamFilePath: imgPath},
		},
	}
	// An in-memory document without a file name is named after the caption.
	parts := documentParts("Page", sdkDoc)
	if len(parts) != 2 || parts[0].Info.Parent != "Page" || parts[1].Path != imgPath {
		t.Fatalf("documentParts: %+v", parts)
	}

	doc := &receptor_v1.Document{Mime: sdkDoc.Mime, Body: sdkDoc.Body}
	contentType, path, err := multipartEvidence(findingWithDoc(doc), nil, parts, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)
	got := readEnvelope(t, path, contentType)
	if len(got) != 5 {
		t.Fatalf("want finding, file, 2 parts, sources; got %d: %+v", len(got), got)
	}
	if got[1].disposition != "file" || got[1].params["name"] != "Page" {
		t.Errorf("file part: %+v", got[1])
	}
	for i, want := range []struct{ name, ct, body string }{{"diagram.png", "image/png", "PNG"}, {"photo.jpg", "image/jpeg", "JPG"}} {
		p := got[2+i]
		if p.disposition != multipartkit.DispositionPart || p.params["name"] != want.name || p.params["parent"] != "Page" || p.contentType != want.ct || p.body != want.body {
			t.Errorf("part %d: %+v", i, p)
		}
		if _, hasFilename := p.params["filename"]; hasFilename {
			t.Errorf("part %d carries filename=, which the app reads as a file part", i)
		}
	}
	if got[4].params["name"] != "receptor_v1.Sources" {
		t.Errorf("last part: %+v", got[4])
	}
}

// The file part's hash covers its parts, so a changed picture changes the
// hash Trustero compares with the previous scan, while a document without
// parts keeps the plain hash of its bytes.
func TestMultipartEvidenceHashCoversParts(t *testing.T) {
	page := []byte(`<img src="tr-part:diagram.png">`)
	plainHash, _ := multipartkit.ComputeHash(bytes.NewReader(page), 1024)
	hashFor := func(picture []byte) string {
		doc := receptor_sdk.Document{Mime: "text/html", Body: page, FileName: "page.html",
			Parts: []receptor_sdk.DocumentPart{{Name: "diagram.png", Mime: "image/png", Kind: receptor_sdk.PartKindImage, Body: picture}}}
		contentType, path, err := multipartEvidence(findingWithDoc(&receptor_v1.Document{Mime: doc.Mime, Body: doc.Body, FileName: doc.FileName}), nil, documentParts("Page", doc), nil)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(path)
		return readEnvelopeHashes(t, path, contentType)[1]
	}
	v1, v2 := hashFor([]byte("PNG v1")), hashFor([]byte("PNG v2"))
	if v1 == v2 || v1 == plainHash {
		t.Errorf("hashes: v1 %s, v2 %s, plain %s", v1, v2, plainHash)
	}
	if again := hashFor([]byte("PNG v1")); again != v1 {
		t.Errorf("the hash is not stable: %s then %s", v1, again)
	}

	contentType, path, err := multipartEvidence(findingWithDoc(&receptor_v1.Document{Mime: "text/html", Body: page, FileName: "page.html"}), nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)
	if got := readEnvelopeHashes(t, path, contentType)[1]; got != plainHash {
		t.Errorf("a document without parts should keep the plain hash: %s vs %s", got, plainHash)
	}
}

func readEnvelopeHashes(t *testing.T, path, contentType string) []string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	boundary := strings.TrimSpace(strings.Split(contentType, "boundary=")[1])
	r, _ := multipartkit.NewMultipartReader(f, boundary, multipartkit.DefaultBufferSize)
	var hashes []string
	for {
		p, err := r.NextPart()
		if err == io.EOF {
			return hashes
		}
		if err != nil {
			t.Fatal(err)
		}
		_, _ = io.Copy(io.Discard, p)
		hashes = append(hashes, p.Header.Get("Content-Hash"))
	}
}
