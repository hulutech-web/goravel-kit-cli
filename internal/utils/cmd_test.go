package utils

import (
    "os/exec"
    "path/filepath"
    "testing"
)

func TestNewCommandWithDir(t *testing.T) {
    dir := t.TempDir()
    cmd := NewCommandWithDir("pwd", []string{}, dir)
    if cmd.Dir != dir {
        t.Fatalf("expected cmd.Dir=%q, got %q", dir, cmd.Dir)
    }

    // Ensure the command is constructed and can run
    out, err := cmd.Output()
    if err != nil {
        // Some systems may not have `pwd` as an external command; fallback to `go env PWD` check
        // We only assert command creation contract, not execution environment.
        if _, ok := err.(*exec.Error); ok {
            t.Skip("pwd not available on this system; skipping execution check")
        } else {
            t.Fatalf("unexpected error running command: %v", err)
        }
    } else {
        got := string(out)
        if !filepath.IsAbs(dir) {
            t.Fatalf("temp dir should be absolute: %q", dir)
        }
        // On macOS/Linux, output ends with a newline
        if !containsPathLine(got, dir) {
            t.Fatalf("expected output to contain working dir; got=%q want contains %q", got, dir)
        }
    }
}

func containsPathLine(output, dir string) bool {
    // Normalize newlines; simple contains is sufficient for our purpose.
    return len(output) > 0 && (output == dir || output == dir+"\n" || output == dir+"\r\n" || (len(output) >= len(dir) && output[:len(dir)] == dir))
}


