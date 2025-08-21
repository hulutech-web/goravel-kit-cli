package utils

import (
    "os"
    "path/filepath"
    "testing"
)

func TestDirectoryAndFileExists(t *testing.T) {
    tempDir, err := os.MkdirTemp("", "goravel-kit-cli-utils-file-*")
    if err != nil {
        t.Fatalf("failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    if !DirectoryExists(tempDir) {
        t.Fatalf("expected DirectoryExists(%q) to be true", tempDir)
    }

    filePath := filepath.Join(tempDir, "test.txt")
    if err := os.WriteFile(filePath, []byte("hello"), 0644); err != nil {
        t.Fatalf("failed to write file: %v", err)
    }

    if !FileExists(filePath) {
        t.Fatalf("expected FileExists(%q) to be true", filePath)
    }

    if DirectoryExists(filePath) {
        t.Fatalf("expected DirectoryExists(%q) to be false for a file", filePath)
    }

    if FileExists(tempDir) {
        t.Fatalf("expected FileExists(%q) to be false for a directory", tempDir)
    }

    nonExist := filepath.Join(tempDir, "does-not-exist")
    if DirectoryExists(nonExist) {
        t.Fatalf("expected DirectoryExists(%q) to be false for non-existing path", nonExist)
    }
    if FileExists(nonExist) {
        t.Fatalf("expected FileExists(%q) to be false for non-existing path", nonExist)
    }
}

func TestMoveAndRemoveDirectory(t *testing.T) {
    baseDir, err := os.MkdirTemp("", "goravel-kit-cli-utils-move-*")
    if err != nil {
        t.Fatalf("failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(baseDir)

    srcDir := filepath.Join(baseDir, "src")
    dstDir := filepath.Join(baseDir, "dst")

    if err := os.MkdirAll(srcDir, 0755); err != nil {
        t.Fatalf("failed to create src dir: %v", err)
    }

    nested := filepath.Join(srcDir, "nested")
    if err := os.MkdirAll(nested, 0755); err != nil {
        t.Fatalf("failed to create nested dir: %v", err)
    }

    contentPath := filepath.Join(nested, "file.txt")
    if err := os.WriteFile(contentPath, []byte("content"), 0644); err != nil {
        t.Fatalf("failed to create nested file: %v", err)
    }

    if err := MoveDirectory(srcDir, dstDir); err != nil {
        t.Fatalf("MoveDirectory failed: %v", err)
    }

    if !DirectoryExists(dstDir) {
        t.Fatalf("expected destination dir to exist after move")
    }

    movedFile := filepath.Join(dstDir, "nested", "file.txt")
    if !FileExists(movedFile) {
        t.Fatalf("expected moved file to exist: %s", movedFile)
    }

    if DirectoryExists(srcDir) {
        t.Fatalf("expected source dir to not exist after move")
    }

    if err := RemoveDirectory(dstDir); err != nil {
        t.Fatalf("RemoveDirectory failed: %v", err)
    }
    if DirectoryExists(dstDir) {
        t.Fatalf("expected destination dir to be removed")
    }
}


