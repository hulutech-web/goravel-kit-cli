package commands

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestUpdateEnvFile_ReplacesValues(t *testing.T) {
    tempDir, err := os.MkdirTemp("", "goravel-kit-cli-commands-env-*")
    if err != nil {
        t.Fatalf("failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    envPath := filepath.Join(tempDir, ".env")
    content := "APP_NAME=Goravel\nAPP_URL=http://localhost\nOTHER=1\n"
    if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
        t.Fatalf("failed to write .env: %v", err)
    }

    if err := updateEnvFile(tempDir, "my-app"); err != nil {
        t.Fatalf("updateEnvFile failed: %v", err)
    }

    updated, err := os.ReadFile(envPath)
    if err != nil {
        t.Fatalf("failed to read updated .env: %v", err)
    }
    got := string(updated)
    if !strings.Contains(got, "APP_NAME=my-app") {
        t.Fatalf("expected APP_NAME replaced, got: %s", got)
    }
    if !strings.Contains(got, "APP_URL=http://localhost:3000") {
        t.Fatalf("expected APP_URL replaced, got: %s", got)
    }
    if !strings.Contains(got, "OTHER=1") {
        t.Fatalf("unexpected modification of unrelated keys, got: %s", got)
    }
}

func TestUpdateEnvFile_NoEnvFile_NoError(t *testing.T) {
    tempDir, err := os.MkdirTemp("", "goravel-kit-cli-commands-env-missing-*")
    if err != nil {
        t.Fatalf("failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    if err := updateEnvFile(tempDir, "my-app"); err != nil {
        t.Fatalf("expected no error when .env missing, got: %v", err)
    }
}

func TestMoveDirectoryCrossPlatform(t *testing.T) {
    baseDir, err := os.MkdirTemp("", "goravel-kit-cli-commands-move-*")
    if err != nil {
        t.Fatalf("failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(baseDir)

    srcDir := filepath.Join(baseDir, "src")
    dstDir := filepath.Join(baseDir, "dst")

    if err := os.MkdirAll(filepath.Join(srcDir, "a", "b"), 0755); err != nil {
        t.Fatalf("failed to create nested dirs: %v", err)
    }
    if err := os.WriteFile(filepath.Join(srcDir, "a", "b", "file.txt"), []byte("data"), 0644); err != nil {
        t.Fatalf("failed to write file: %v", err)
    }

    if err := moveDirectoryCrossPlatform(srcDir, dstDir); err != nil {
        t.Fatalf("moveDirectoryCrossPlatform failed: %v", err)
    }

    if _, err := os.Stat(filepath.Join(dstDir, "a", "b", "file.txt")); err != nil {
        t.Fatalf("expected moved file to exist: %v", err)
    }
    if _, err := os.Stat(srcDir); !os.IsNotExist(err) {
        t.Fatalf("expected src dir removed; stat err=%v", err)
    }
}

func TestCopyFile(t *testing.T) {
    dir, err := os.MkdirTemp("", "goravel-kit-cli-commands-copy-*")
    if err != nil {
        t.Fatalf("failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(dir)

    src := filepath.Join(dir, "src.txt")
    dst := filepath.Join(dir, "dst.txt")
    data := []byte("hello world")
    if err := os.WriteFile(src, data, 0600); err != nil {
        t.Fatalf("failed to write src: %v", err)
    }

    if err := copyFile(src, dst); err != nil {
        t.Fatalf("copyFile failed: %v", err)
    }

    got, err := os.ReadFile(dst)
    if err != nil {
        t.Fatalf("failed to read dst: %v", err)
    }
    if string(got) != string(data) {
        t.Fatalf("copied content mismatch: got=%q want=%q", string(got), string(data))
    }

    srcInfo, _ := os.Stat(src)
    dstInfo, _ := os.Stat(dst)
    if dstInfo.Mode() != srcInfo.Mode() {
        t.Fatalf("expected file mode copied, got %v want %v", dstInfo.Mode(), srcInfo.Mode())
    }
}

func TestHasNextMirror(t *testing.T) {
    mirrors := []struct {
        name    string
        url     string
        sshURL  string
        enabled bool
    }{
        {name: "GitHub", enabled: true},
        {name: "Gitee", enabled: true},
    }

    if !hasNextMirror(mirrors, "GitHub") {
        t.Fatalf("expected next mirror after GitHub")
    }
    if hasNextMirror(mirrors, "Gitee") {
        t.Fatalf("expected no next mirror after last enabled mirror")
    }

    mirrorsDisabled := []struct {
        name    string
        url     string
        sshURL  string
        enabled bool
    }{
        {name: "GitHub", enabled: false},
        {name: "Gitee", enabled: true},
    }
    if hasNextMirror(mirrorsDisabled, "Gitee") {
        t.Fatalf("expected no next mirror when at last enabled entry")
    }
}


