// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package multipartkit

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"reflect"

	"google.golang.org/protobuf/proto"
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
func (mb *MultipartBuilder) GetBoundary() string {
	return mb.writer.Boundary()
}

// ComputeHash streams data from the provided `io.Reader` and calculates the SHA-256 hash.
// The hash is returned as UrlSafe base64 encoded string.
func ComputeHash(reader io.Reader, bufferSize int) (string, error) {
	hash := sha256.New()

	_, err := io.CopyBuffer(hash, reader, make([]byte, bufferSize))
	if err != nil {
		return "", fmt.Errorf("failed to read and calculate hash: %v", err)
	}

	hashValue := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return hashValue, nil
}

// AddProtobuf writes a single or slice of Protobuf messages as a part of the multipart stream,
// including Content-Size and Content-Hash headers.
func (mb *MultipartBuilder) AddProtobuf(partName string, pb interface{}) error {
	switch v := pb.(type) {
	case proto.Message:
		return mb.addSingleOrMultipleProtobuf(partName, []proto.Message{v})

	case []proto.Message:
		return mb.addSingleOrMultipleProtobuf(partName, v)

	default:
		// Check for the case where pb is a slice of pointers to protobuf messages.
		rv := reflect.ValueOf(pb)
		if rv.Kind() == reflect.Slice {
			protoMessages := make([]proto.Message, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				elem := rv.Index(i).Interface()

				// Check if this element implements proto.Message
				if msg, ok := elem.(proto.Message); ok {
					protoMessages[i] = msg
				} else {
					return fmt.Errorf("unsupported element in slice: expected proto.Message at index %d, got %T", i, elem)
				}
			}
			return mb.addSingleOrMultipleProtobuf(partName, protoMessages)
		}

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

	reader := bytes.NewReader(concatenatedData)

	// Compute the hash before creating the part
	contentHash, err := ComputeHash(bytes.NewReader(concatenatedData), mb.bufferSize)
	if err != nil {
		return fmt.Errorf("failed to compute hash for protobuf part: %v", err)
	}

	// Create the part with all headers, including Content-Hash
	partWriter, err := mb.writer.CreatePart(map[string][]string{
		"Content-Disposition": {fmt.Sprintf(`protobuf; name="%s"`, partName)},
		"Content-Type":        {"application/protobuf"},
		"Content-Size":        {fmt.Sprintf("%d", len(concatenatedData))},
		"Content-Hash":        {contentHash},
	})
	if err != nil {
		return fmt.Errorf("failed to create multipart part: %v", err)
	}

	_, err = io.Copy(partWriter, reader)
	if err != nil {
		return fmt.Errorf("failed to write protobuf data: %v", err)
	}

	return nil
}

// AddFile writes the content of a file as a part of the multipart stream with Content-Size and Content-Hash headers.
func (mb *MultipartBuilder) AddFile(partName, filePath, contentType string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// Compute the file hash before creating the part
	_, err = file.Seek(0, io.SeekStart) // Ensure we are at the start of the file
	if err != nil {
		return fmt.Errorf("failed to seek file: %v", err)
	}
	contentHash, err := ComputeHash(file, mb.bufferSize)
	if err != nil {
		return fmt.Errorf("failed to compute hash for file: %v", err)
	}

	_, err = file.Seek(0, io.SeekStart) // Rewind the file to the beginning before writing
	if err != nil {
		return fmt.Errorf("failed to seek file: %v", err)
	}

	// Create the part with all headers, including Content-Hash
	partWriter, err := mb.writer.CreatePart(map[string][]string{
		"Content-Disposition": {fmt.Sprintf(`file; name="%s"; filename="%s"`, partName, filePath)},
		"Content-Type":        {contentType},
		"Content-Size":        {fmt.Sprintf("%d", fileInfo.Size())},
		"Content-Hash":        {contentHash},
	})
	if err != nil {
		return fmt.Errorf("failed to create multipart part: %v", err)
	}

	_, err = io.CopyBuffer(partWriter, file, make([]byte, mb.bufferSize))
	if err != nil {
		return fmt.Errorf("failed to write file part: %v", err)
	}

	return nil
}

// AddBytes writes a file (provided as bytes) into the multipart builder with the given filename and MIME type.
func (mb *MultipartBuilder) AddBytes(partName, fileName, contentType string, data []byte) error {
	reader := bytes.NewReader(data)

	// Compute the hash for the provided byte array
	contentHash, err := ComputeHash(reader, mb.bufferSize)
	if err != nil {
		return fmt.Errorf("failed to compute hash for data: %v", err)
	}

	// Rewind the reader after calculating the hash
	reader.Seek(0, io.SeekStart)

	// Create the part with all headers, including Content-Hash
	partWriter, err := mb.writer.CreatePart(map[string][]string{
		"Content-Disposition": {fmt.Sprintf(`file; name="%s"; filename="%s"`, partName, fileName)},
		"Content-Type":        {contentType},
		"Content-Length":      {fmt.Sprintf("%d", len(data))},
		"Content-Hash":        {contentHash},
	})
	if err != nil {
		return fmt.Errorf("failed to create multipart part: %v", err)
	}

	_, err = io.Copy(partWriter, reader)
	if err != nil {
		return fmt.Errorf("failed to write data part: %v", err)
	}

	return nil
}

// Finalize writes the ending boundary and closes the multipart writer.
func (mb *MultipartBuilder) Finalize() error {
	return mb.writer.Close()
}
