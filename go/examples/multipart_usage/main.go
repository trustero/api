// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/trustero/api/go/receptor_sdk/multipartkit"
)

func main() {
	// 1. Write Multipart Data Using MultipartBuilder
	multipartFilePath := "multipart_data.txt"
	dstFile, err := os.Create(multipartFilePath)
	if err != nil {
		log.Fatalf("Failed to create multipart file: %v", err)
	}
	defer dstFile.Close()

	bufferSize := multipartkit.DefaultBufferSize

	// Initialize the multipart builder
	builder, err := multipartkit.NewMultipartBuilder(dstFile, bufferSize)
	if err != nil {
		log.Fatalf("Failed to create multipart builder: %v", err)
	}
	// Get the boundary string
	boundary := builder.GetBoundary()

	// Add a file part
	err = builder.AddFile("test1.csv", "test1.csv", "application/csv")
	if err != nil {
		log.Fatalf("Failed to add file part: %v", err)
	}

	// Add another file part
	err = builder.AddFile("test2.docx", "test2.docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	if err != nil {
		log.Fatalf("Failed to add file part: %v", err)
	}

	// Finalize the multipart data
	err = builder.Finalize()
	if err != nil {
		log.Fatalf("Failed to finalize multipart builder: %v", err)
	}

	fmt.Println("Multipart data successfully written to", multipartFilePath)

	// 2. Read Multipart Data and Print Metadata Using MultipartReader
	srcFile, err := os.Open(multipartFilePath)
	if err != nil {
		log.Fatalf("Failed to open multipart file for reading: %v", err)
	}
	defer srcFile.Close()

	reader, err := multipartkit.NewMultipartReader(srcFile, boundary, bufferSize)
	if err != nil {
		log.Fatalf("Failed to initialize multipart reader: %v", err)
	}

	// Stream metadata for all parts and print it as JSON
	fmt.Println("\nMetadata for the Multipart File (JSON):")
	err = reader.MetadataJSON(os.Stdout)
	if err != nil {
		log.Fatalf("Failed to print metadata as JSON: %v", err)
	}
	fmt.Println()

}
