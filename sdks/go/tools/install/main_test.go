package main

import (
	"os"
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

func TestNativeReleaseTag(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "version.go"), []byte("package mosscore\n\nconst NativeLibReleaseTag = \"c-sdk-v1.2.3\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	tag, err := nativeReleaseTag(root)
	if err != nil {
		t.Fatalf("nativeReleaseTag returned an error: %v", err)
	}
	if tag != "c-sdk-v1.2.3" {
		t.Fatalf("nativeReleaseTag() = %q, want c-sdk-v1.2.3", tag)
	}
}

func TestInstallReceipt(t *testing.T) {
	path := filepath.Join(t.TempDir(), installReceiptName)
	if err := writeReceipt(path, "libmoss-v0.9.0-x86_64-unknown-linux-gnu.tar.gz", "abc123"); err != nil {
		t.Fatalf("writeReceipt returned an error: %v", err)
	}
	if !receiptMatches(path, "libmoss-v0.9.0-x86_64-unknown-linux-gnu.tar.gz", "abc123") {
		t.Fatal("receiptMatches() = false, want true")
	}
	if receiptMatches(path, "libmoss-v0.9.0-x86_64-unknown-linux-gnu.tar.gz", "different") {
		t.Fatal("receiptMatches() = true for a different checksum, want false")
	}
}
