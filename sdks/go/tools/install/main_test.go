package main

import (
	"path/filepath"
	"testing"
)

func TestSafeArchiveTarget(t *testing.T) {
	root := t.TempDir()

	target, err := safeArchiveTarget(root, "include/libmoss.h")
	if err != nil {
		t.Fatalf("safeArchiveTarget returned an error: %v", err)
	}
	want := filepath.Join(root, "include", "libmoss.h")
	if target != want {
		t.Fatalf("safeArchiveTarget() = %q, want %q", target, want)
	}

	for _, input := range []string{"..", "../outside", "/outside"} {
		if _, err := safeArchiveTarget(root, input); err == nil {
			t.Errorf("safeArchiveTarget(%q) succeeded, want error", input)
		}
	}
}
