// Command install downloads Moss C SDK static libraries from GitHub Releases
// into sdks/go/bindings for local CGO linking.
package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	defaultReleaseTag = "c-sdk-v0.9.0"
	defaultRepo       = "usemoss/moss"
)

type platform struct {
	id      string
	triple  string
	libFile string
	srcLib  string
	ldflags string
}

var platforms = []platform{
	{
		id: "linux-amd64", triple: "x86_64-unknown-linux-gnu",
		libFile: "libmoss.a", srcLib: "lib/libmoss.a",
	},
	{
		id: "linux-arm64", triple: "aarch64-unknown-linux-gnu",
		libFile: "libmoss.a", srcLib: "lib/libmoss.a",
	},
	{
		id: "darwin-arm64", triple: "aarch64-apple-darwin",
		libFile: "libmoss.a", srcLib: "lib/libmoss.a",
	},
	{
		id: "windows-amd64", triple: "x86_64-pc-windows-msvc",
		libFile: "moss.lib", srcLib: "lib/moss.lib",
	},
}

func main() {
	all := flag.Bool("all", false, "install libraries for all supported platforms")
	releaseTag := flag.String("release", defaultReleaseTag, "C SDK GitHub release tag")
	repo := flag.String("repo", defaultRepo, "GitHub repository (owner/name)")
	bindingsDir := flag.String("bindings", "", "bindings directory (default: auto-detect)")
	force := flag.Bool("force", false, "re-download even if the library is already installed")
	flag.Parse()

	root, err := resolveBindingsDir(*bindingsDir)
	if err != nil {
		fatal(err)
	}

	version := strings.TrimPrefix(*releaseTag, "c-sdk-")
	version = strings.TrimPrefix(version, "v")
	baseURL := fmt.Sprintf("https://github.com/%s/releases/download/%s", *repo, *releaseTag)

	checksums, err := fetchChecksums(baseURL)
	if err != nil {
		fatal(err)
	}

	targets, err := selectPlatforms(*all)
	if err != nil {
		fatal(err)
	}

	for _, p := range targets {
		if err := installPlatform(root, baseURL, version, p, checksums, *force); err != nil {
			fatal(fmt.Errorf("%s: %w", p.id, err))
		}
	}

	fmt.Printf("Installed Moss C SDK %s into %s\n", *releaseTag, root)
}

func resolveBindingsDir(explicit string) (string, error) {
	if explicit != "" {
		return filepath.Abs(explicit)
	}
	if env := strings.TrimSpace(os.Getenv("MOSS_BINDINGS_DIR")); env != "" {
		return filepath.Abs(env)
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to locate install tool path")
	}
	return filepath.Abs(filepath.Join(filepath.Dir(file), "..", "..", "bindings"))
}

func selectPlatforms(all bool) ([]platform, error) {
	if all {
		return platforms, nil
	}
	id, err := currentPlatformID()
	if err != nil {
		return nil, err
	}
	for _, p := range platforms {
		if p.id == id {
			return []platform{p}, nil
		}
	}
	return nil, fmt.Errorf("unsupported platform %s/%s", runtime.GOOS, runtime.GOARCH)
}

func currentPlatformID() (string, error) {
	switch runtime.GOOS + "/" + runtime.GOARCH {
	case "linux/amd64":
		return "linux-amd64", nil
	case "linux/arm64":
		return "linux-arm64", nil
	case "darwin/arm64":
		return "darwin-arm64", nil
	case "windows/amd64":
		return "windows-amd64", nil
	default:
		return "", fmt.Errorf("unsupported platform %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}

func fetchChecksums(baseURL string) (map[string]string, error) {
	resp, err := http.Get(baseURL + "/checksums-sha256.txt")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("checksums: HTTP %d", resp.StatusCode)
	}

	checksums := map[string]string{}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		parts := strings.Fields(strings.TrimSpace(scanner.Text()))
		if len(parts) == 2 {
			checksums[parts[1]] = parts[0]
		}
	}
	return checksums, scanner.Err()
}

func installPlatform(bindingsRoot, baseURL, version string, p platform, checksums map[string]string, force bool) error {
	archive := fmt.Sprintf("libmoss-v%s-%s.tar.gz", version, p.triple)
	wantChecksum, ok := checksums[archive]
	if !ok {
		return fmt.Errorf("checksum not found for %s", archive)
	}

	destDir := filepath.Join(bindingsRoot, "lib", p.id)
	destLib := filepath.Join(destDir, p.libFile)
	headerPath := filepath.Join(bindingsRoot, "include", "libmoss.h")

	if !force {
		if err := verifyFileChecksum(destLib, wantChecksum); err == nil {
			if _, err := os.Stat(headerPath); err == nil {
				fmt.Printf("skip %s (already installed)\n", p.id)
				return nil
			}
		}
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(headerPath), 0o755); err != nil {
		return err
	}

	tmp, err := os.MkdirTemp("", "moss-install-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	archivePath := filepath.Join(tmp, archive)
	if err := downloadFile(baseURL+"/"+archive, archivePath, wantChecksum); err != nil {
		return err
	}

	extractedRoot := filepath.Join(tmp, fmt.Sprintf("libmoss-v%s-%s", version, p.triple))
	if err := extractTarGz(archivePath, extractedRoot, version, p.triple); err != nil {
		return err
	}

	if err := copyFile(filepath.Join(extractedRoot, "include", "libmoss.h"), headerPath); err != nil {
		return err
	}
	if err := copyFile(filepath.Join(extractedRoot, p.srcLib), destLib); err != nil {
		return err
	}

	fmt.Printf("installed %s -> %s\n", p.id, destLib)
	return nil
}

func extractTarGz(archivePath, destRoot, version, triple string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	prefix := fmt.Sprintf("libmoss-v%s-%s/", version, triple)
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if hdr.Typeflag != tar.TypeReg || !strings.HasPrefix(hdr.Name, prefix) {
			continue
		}
		rel := strings.TrimPrefix(hdr.Name, prefix)
		target := filepath.Join(destRoot, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := writeFile(target, tr, hdr.Mode); err != nil {
			return err
		}
	}
}

func writeFile(path string, r io.Reader, mode int64) error {
	out, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileMode(mode))
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, r)
	return err
}

func fileMode(mode int64) os.FileMode {
	if mode == 0 {
		return 0o644
	}
	return os.FileMode(mode)
}

func downloadFile(url, dest, wantChecksum string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: HTTP %d", url, resp.StatusCode)
	}

	tmp, err := os.CreateTemp(filepath.Dir(dest), ".download-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

	hasher := sha256.New()
	if _, err := io.Copy(io.MultiWriter(tmp, hasher), resp.Body); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	got := hex.EncodeToString(hasher.Sum(nil))
	if got != wantChecksum {
		os.Remove(tmpPath)
		return fmt.Errorf("checksum mismatch for %s: got %s want %s", filepath.Base(dest), got, wantChecksum)
	}
	return os.Rename(tmpPath, dest)
}

func verifyFileChecksum(path, want string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return err
	}
	if hex.EncodeToString(hasher.Sum(nil)) != want {
		return fmt.Errorf("checksum mismatch")
	}
	return nil
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "moss install: %v\n", err)
	os.Exit(1)
}
