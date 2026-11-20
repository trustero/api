// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package multipartkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/textproto"
	"os"
	"strings"
)

// A document's parts (receptor_sdk.DocumentPart) travel in the evidence
// envelope as parts of their own, after the file part they belong to:
//
//	Content-Disposition: part; name="diagram.png"; parent="page.html"
//	Content-Type: image/png
//	Content-Kind: image
//	Content-Locator: {"attachment_id":"att1","page_id":"123"}
//	Content-Size: 48213
//	Content-Hash: <url-safe base64 sha256>
//
// The disposition type "part" is new; a reader that only knows "protobuf"
// and "file" skips these parts, so an envelope with parts is still read by
// older servers (they store the document without its parts).
const (
	// DispositionPart is the Content-Disposition type of a document part.
	DispositionPart = "part"
	// HeaderContentKind carries the part's kind (receptor_sdk.PartKind).
	HeaderContentKind = "Content-Kind"
	// HeaderContentLocator carries the part's locator as a JSON object of
	// string values.
	HeaderContentLocator = "Content-Locator"
)

// PartInfo describes one document part in the envelope.
type PartInfo struct {
	Name        string            // Part name, unique within its parent document
	Parent      string            // Name of the file part the part belongs to
	ContentType string            // Mime type of the part
	Kind        string            // receptor_sdk.PartKind as a string
	Locator     map[string]string // Where the part came from at the provider
}

// AddPartBytes writes a document part held in memory.
func (mb *MultipartBuilder) AddPartBytes(info PartInfo, data []byte) error {
	contentHash, err := ComputeHash(bytes.NewReader(data), min(mb.bufferSize, max(len(data), 1)))
	if err != nil {
		return fmt.Errorf("failed to compute hash for part %q: %v", info.Name, err)
	}
	headers, err := partHeaders(info, int64(len(data)), contentHash)
	if err != nil {
		return err
	}
	partWriter, err := mb.writer.CreatePart(headers)
	if err != nil {
		return fmt.Errorf("failed to create multipart part: %v", err)
	}
	if _, err := io.Copy(partWriter, bytes.NewReader(data)); err != nil {
		return fmt.Errorf("failed to write part %q: %v", info.Name, err)
	}
	return nil
}

// AddPartFile writes a document part from a file, streaming it.
func (mb *MultipartBuilder) AddPartFile(info PartInfo, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open part file: %v", err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get part file info: %v", err)
	}
	contentHash, err := ComputeHash(file, mb.bufferSize)
	if err != nil {
		return fmt.Errorf("failed to compute hash for part %q: %v", info.Name, err)
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek part file: %v", err)
	}
	headers, err := partHeaders(info, fileInfo.Size(), contentHash)
	if err != nil {
		return err
	}
	partWriter, err := mb.writer.CreatePart(headers)
	if err != nil {
		return fmt.Errorf("failed to create multipart part: %v", err)
	}
	if _, err := io.CopyBuffer(partWriter, file, make([]byte, mb.bufferSize)); err != nil {
		return fmt.Errorf("failed to write part %q: %v", info.Name, err)
	}
	return nil
}

func partHeaders(info PartInfo, size int64, contentHash string) (map[string][]string, error) {
	if strings.TrimSpace(info.Name) == "" {
		return nil, fmt.Errorf("a document part needs a name")
	}
	params := map[string]string{"name": info.Name}
	if info.Parent != "" {
		params["parent"] = info.Parent
	}
	// FormatMediaType quotes and, for a name outside ASCII, RFC 2231-encodes
	// the parameters, which ParseMediaType on the other side undoes. It gives
	// up ("") on a value it cannot represent at all.
	disposition := mime.FormatMediaType(DispositionPart, params)
	if disposition == "" {
		return nil, fmt.Errorf("part name %q or parent %q cannot be put in a Content-Disposition", info.Name, info.Parent)
	}
	contentType := headerToken(info.ContentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	kind := headerToken(info.Kind)
	headers := map[string][]string{
		"Content-Disposition": {disposition},
		"Content-Type":        {contentType},
		"Content-Size":        {fmt.Sprintf("%d", size)},
		"Content-Hash":        {contentHash},
	}
	if kind != "" {
		headers[HeaderContentKind] = []string{kind}
	}
	if len(info.Locator) > 0 {
		locator, err := json.Marshal(info.Locator)
		if err != nil {
			return nil, fmt.Errorf("failed to encode locator of part %q: %v", info.Name, err)
		}
		headers[HeaderContentLocator] = []string{string(locator)}
	}
	return headers, nil
}

// headerToken keeps a header value to one line of printable ASCII; the
// values come from receptor code, but a mime read off page content could
// carry a line break and inject a header.
func headerToken(v string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x21 || r > 0x7e {
			return -1
		}
		return r
	}, strings.TrimSpace(v))
}

// ParsePartInfo reads a document part's description from its headers. ok is
// false for any part that is not a document part (a protobuf or file part).
func ParsePartInfo(h textproto.MIMEHeader) (info PartInfo, ok bool) {
	dispositionType, params, err := mime.ParseMediaType(h.Get("Content-Disposition"))
	if err != nil || !strings.EqualFold(dispositionType, DispositionPart) {
		return PartInfo{}, false
	}
	info = PartInfo{
		Name:        params["name"],
		Parent:      params["parent"],
		ContentType: h.Get("Content-Type"),
		Kind:        strings.ToLower(strings.TrimSpace(h.Get(HeaderContentKind))),
	}
	if raw := h.Get(HeaderContentLocator); raw != "" {
		locator := map[string]string{}
		if err := json.Unmarshal([]byte(raw), &locator); err == nil {
			info.Locator = locator
		}
	}
	return info, info.Name != ""
}
