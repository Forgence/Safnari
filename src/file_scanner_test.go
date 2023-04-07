// file_scanner_test.go
package main

import (
	"context"
	"os"
	"testing"
)

func TestTraverseFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create a temporary file in the tempDir
	file, err := os.CreateTemp(tempDir, "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	fileProcessor := make(chan FileData, 100)
	defer close(fileProcessor)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go traverseFiles(ctx, tempDir, fileProcessor, true, true, "", 0, nil, nil, false, nil)
	for fileData := range fileProcessor {
		if fileData.FilePath == "" {
			t.Error("File path is empty")
		}
	}
}
