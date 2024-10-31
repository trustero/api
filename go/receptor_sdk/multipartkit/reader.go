// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package multipartkit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"strings"
)

// MultipartReader reads multipart data from a stream and supports configurable buffer size.
type MultipartReader struct {
	reader     *multipart.Reader
	bufferSize int
}

// PartMetadata describes metadata for each part of the multipart message.
type PartMetadata struct {
	PartName string                 `json:"part_name"`
	PartType string                 `json:"part_type"`
	Headers  map[string]interface{} `json:"headers"`
}

// NewMultipartReader creates a new instance of MultipartReader.
func NewMultipartReader(r io.Reader, boundary string, bufferSize int) (*MultipartReader, error) {
	if boundary == "" {
		return nil, fmt.Errorf("boundary cannot be empty")
	}
	if bufferSize <= 0 {
		bufferSize = DefaultBufferSize
	}

	bufReader := bufio.NewReader(r)

	mr := &MultipartReader{
		bufferSize: bufferSize,
		reader:     multipart.NewReader(bufReader, boundary),
	}

	return mr, nil
}

// NextPart returns the next part of the multipart stream using the current reader.
func (mr *MultipartReader) NextPart() (*multipart.Part, error) {
	return mr.reader.NextPart()
}

// MetadataJSON streams the JSON representation of the metadata for all parts in the multipart message.
func (mr *MultipartReader) MetadataJSON(w io.Writer) error {
	var metadataList []PartMetadata

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to get next part: %v", err)
		}

		contentDisposition := part.Header.Get("Content-Disposition")
		_, dispositionParams, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			continue
		}

		partName, exists := dispositionParams["name"]
		if !exists {
			partName = "Unknown"
		}
		partType := "Unknown"
		contentType := part.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/protobuf") {
			partType = "Protobuf"
		} else if contentDisposition != "" {
			partType = "File"
		}

		headersMap := headerToMap(part.Header)
		metadata := PartMetadata{
			PartName: partName,
			PartType: partType,
			Headers:  headersMap,
		}
		metadataList = append(metadataList, metadata)
	}

	jsonData, err := json.MarshalIndent(metadataList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON metadata: %v", err)
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSON metadata: %v", err)
	}
	return nil
}

// headerToMap converts MIMEHeader into a flat map[string]interface{} for JSON encoding.
func headerToMap(header textproto.MIMEHeader) map[string]interface{} {
	flatHeaders := make(map[string]interface{})
	for key, values := range header {
		if len(values) == 1 {
			flatHeaders[key] = values[0]
		} else {
			flatHeaders[key] = values
		}
	}
	return flatHeaders
}
