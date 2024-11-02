// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package multipartkit

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"mime/multipart"
	"os"
)

// DefaultBufferSize defines a default buffer size (5 MB) for writing parts.
const DefaultBufferSize = 5 * 1024 * 1024 // 5 MB

// MultipartBuilder is responsible for building multipart data with writers and buffer size.
type MultipartBuilder struct {
	writer     *multipart.Writer
	bufferSize int
}

// NewMultipartBuilder initializes a new MultipartBuilder using an io.Writer and buffer size.
func NewMultipartBuilder(w io.Writer, bufferSize int) (*MultipartBuilder, error) {
	if bufferSize <= 0 {
		bufferSize = DefaultBufferSize
	}
	writer := multipart.NewWriter(w)
	return &MultipartBuilder{
		writer:     writer,
		bufferSize: bufferSize,
	}, nil
}

// Boundary retrieves the boundary used by the multipart writer.
// This Boundary string is used to separate parts in the multipart stream
// and should be saved separately for reading the multipart data to follow
// the RFC 822 standard for multipart messages.
func (mb *MultipartBuilder) GetBoundary() string {
	return mb.writer.Boundary()
}

// AddProtobuf writes a single or slice of Protobuf messages as a part of the multipart stream,
// including Content-Size header.
func (mb *MultipartBuilder) AddProtobuf(partName string, pb interface{}) error {
	switch v := pb.(type) {
	case proto.Message:
		return mb.addSingleOrMultipleProtobuf(partName, []proto.Message{v})

	case []proto.Message:
		return mb.addSingleOrMultipleProtobuf(partName, v)

	default:
		return fmt.Errorf("unsupported type: expected proto.Message or []proto.Message, got %T", pb)
	}
}

func (mb *MultipartBuilder) addSingleOrMultipleProtobuf(partName string, pbs []proto.Message) error {
	var concatenatedData []byte
	for _, pb := range pbs {
		marshaledData, err := proto.Marshal(pb)
		if err != nil {
			return fmt.Errorf("failed to marshal protobuf message: %v", err)
		}
		concatenatedData = append(concatenatedData, marshaledData...)
	}

	contentSize := len(concatenatedData)
	partWriter, err := mb.writer.CreatePart(map[string][]string{
		"Content-Disposition": {fmt.Sprintf(`protobuf; name="%s"`, partName)},
		"Content-Type":        {"application/protobuf"},
		"Content-Size":        {fmt.Sprintf("%d", contentSize)},
	})
	if err != nil {
		return fmt.Errorf("failed to create multipart part: %v", err)
	}
	_, err = partWriter.Write(concatenatedData)
	if err != nil {
		return fmt.Errorf("failed to write protobuf part: %v", err)
	}

	return nil
}

// AddFile writes the content of a file as a part of the multipart stream with Content-Size header.
func (mb *MultipartBuilder) AddFile(partName, filePath, contentType string) error {
	// Open the file to add to the multipart stream.
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Get file size by checking FileInfo of opened file
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	contentSize := fileInfo.Size()

	// Create a part writer with Content-Size header
	partWriter, err := mb.writer.CreatePart(map[string][]string{
		"Content-Disposition": {fmt.Sprintf(`file; name="%s"; filename="%s"`, partName, filePath)},
		"Content-Type":        {contentType},
		"Content-Size":        {fmt.Sprintf("%d", contentSize)},
	})
	if err != nil {
		return fmt.Errorf("failed to create multipart part: %v", err)
	}

	// Write file content into the part.
	buf := make([]byte, mb.bufferSize)
	_, err = io.CopyBuffer(partWriter, file, buf)
	if err != nil {
		return fmt.Errorf("failed to write file part: %v", err)
	}

	return nil
}

// AddBytes writes a file (provided as bytes) into the multipart builder with the given filename and MIME type.
func (mb *MultipartBuilder) AddBytes(partName, fileName, contentType string, data []byte) error {
	// Create a part writer with Content-Disposition, Content-Type, and Content-Length headers
	partWriter, err := mb.writer.CreatePart(map[string][]string{
		"Content-Disposition": {fmt.Sprintf(`file; name="%s"; filename="%s"`, partName, fileName)},
		"Content-Type":        {contentType},
		"Content-Length":      {fmt.Sprintf("%d", len(data))},
	})
	if err != nil {
		return fmt.Errorf("failed to create multipart part: %v", err)
	}
	_, err = partWriter.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write bytes to part: %v", err)
	}

	return nil
}

// Finalize writes the ending boundary and closes the multipart writer.
func (mb *MultipartBuilder) Finalize() error {
	return mb.writer.Close()
}
