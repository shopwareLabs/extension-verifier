package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "copyfiles-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("Failed to remove temporary directory: %v", err)
		}
	}()

	// Create source directory structure
	srcDir := filepath.Join(tempDir, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create a normal file
	normalFile := filepath.Join(srcDir, "normal.txt")
	if err := os.WriteFile(normalFile, []byte("normal content"), 0644); err != nil {
		t.Fatalf("Failed to create normal file: %v", err)
	}

	// Create a .devenv directory with a file
	devenvDir := filepath.Join(srcDir, ".devenv")
	if err := os.MkdirAll(devenvDir, 0755); err != nil {
		t.Fatalf("Failed to create .devenv directory: %v", err)
	}
	devenvFile := filepath.Join(devenvDir, "devenv.txt")
	if err := os.WriteFile(devenvFile, []byte("devenv content"), 0644); err != nil {
		t.Fatalf("Failed to create file in .devenv: %v", err)
	}

	// Create a .direnv directory with a file
	direnvDir := filepath.Join(srcDir, ".direnv")
	if err := os.MkdirAll(direnvDir, 0755); err != nil {
		t.Fatalf("Failed to create .direnv directory: %v", err)
	}
	direnvFile := filepath.Join(direnvDir, "direnv.txt")
	if err := os.WriteFile(direnvFile, []byte("direnv content"), 0644); err != nil {
		t.Fatalf("Failed to create file in .direnv: %v", err)
	}

	// Create a regular subdirectory with a file
	subDir := filepath.Join(srcDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	subFile := filepath.Join(subDir, "sub.txt")
	if err := os.WriteFile(subFile, []byte("sub content"), 0644); err != nil {
		t.Fatalf("Failed to create file in subdirectory: %v", err)
	}

	// Create destination directory
	dstDir := filepath.Join(tempDir, "dst")

	// Copy files from src to dst
	if err := copyFiles(srcDir, dstDir); err != nil {
		t.Fatalf("copyFiles failed: %v", err)
	}

	// Check if normal file was copied
	dstNormalFile := filepath.Join(dstDir, "normal.txt")
	if _, err := os.Stat(dstNormalFile); os.IsNotExist(err) {
		t.Errorf("Normal file was not copied")
	}

	// Check if file in subdirectory was copied
	dstSubFile := filepath.Join(dstDir, "subdir", "sub.txt")
	if _, err := os.Stat(dstSubFile); os.IsNotExist(err) {
		t.Errorf("File in subdirectory was not copied")
	}

	// Check if .devenv directory was excluded
	dstDevenvDir := filepath.Join(dstDir, ".devenv")
	if _, err := os.Stat(dstDevenvDir); !os.IsNotExist(err) {
		t.Errorf(".devenv directory was not excluded")
	}

	// Check if .direnv directory was excluded
	dstDirenvDir := filepath.Join(dstDir, ".direnv")
	if _, err := os.Stat(dstDirenvDir); !os.IsNotExist(err) {
		t.Errorf(".direnv directory was not excluded")
	}
}