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
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	defaultReleaseTag  = "c-sdk-v0.9.0"
	defaultRepo        = "usemoss/moss"
	bindingsModulePath = "github.com/usemoss/moss/sdks/go/bindings"
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

	root, err := resolveWritableBindingsDir(*bindingsDir)
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
		libDir := filepath.Join(root, "lib", p.id)
		if err := installPlatform(root, libDir, true, baseURL, version, p, checksums, *force); err != nil {
			fatal(fmt.Errorf("%s: %w", p.id, err))
		}
	}

	fmt.Printf("Installed Moss C SDK %s\n", *releaseTag)
	fmt.Printf("Native libraries installed under %s\n", root)
}

// resolveWritableBindingsDir returns the bindings package CGO will compile.
// Downloaded Go modules are read-only, so an external consumer is vendored
// before installing the native library beside the bindings source.
func resolveWritableBindingsDir(explicit string) (string, error) {
	root, err := resolveBindingsDir(explicit)
	if err != nil {
		return "", err
	}
	if isWritableDir(root) {
		return root, nil
	}
	if explicit != "" || strings.TrimSpace(os.Getenv("MOSS_BINDINGS_DIR")) != "" {
		return "", fmt.Errorf("bindings directory %s is not writable", root)
	}

	goMod, err := goEnv("GOMOD")
	if err != nil || goMod == os.DevNull || goMod == "" {
		return "", fmt.Errorf("the downloaded bindings module is read-only; run this command from a Go module that imports the Moss SDK")
	}
	if out, err := exec.Command("go", "mod", "vendor").CombinedOutput(); err != nil {
		return "", fmt.Errorf("vendor Moss bindings for native installation: %w\n%s", err, strings.TrimSpace(string(out)))
	}

	root, err = bindingsDirFromGoList("-mod=vendor")
	if err != nil {
		return "", fmt.Errorf("locate vendored Moss bindings after `go mod vendor`: %w", err)
	}
	if !isWritableDir(root) {
		return "", fmt.Errorf("vendored bindings directory %s is not writable", root)
	}
	return root, nil
}

func resolveBindingsDir(explicit string) (string, error) {
	if explicit != "" {
		return filepath.Abs(explicit)
	}
	if env := strings.TrimSpace(os.Getenv("MOSS_BINDINGS_DIR")); env != "" {
		return filepath.Abs(env)
	}
	if dir, err := bindingsDirFromGoList(); err == nil {
		return dir, nil
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to locate bindings directory; set MOSS_BINDINGS_DIR or run from a module that requires github.com/usemoss/moss/sdks/go/bindings")
	}
	return filepath.Abs(filepath.Join(filepath.Dir(file), "..", "..", "bindings"))
}

func bindingsDirFromGoList(args ...string) (string, error) {
	args = append([]string{"list"}, args...)
	args = append(args, "-f", "{{.Dir}}", bindingsModulePath)
	cmd := exec.Command("go", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	dir := strings.TrimSpace(string(out))
	if dir == "" {
		return "", errors.New("bindings package directory not found")
	}
	return filepath.Abs(dir)
}

func isWritableDir(dir string) bool {
	if dir == "" {
		return false
	}
	test := filepath.Join(dir, ".moss-install-write-test")
	if err := os.WriteFile(test, []byte("ok"), 0o644); err != nil {
		return false
	}
	_ = os.Remove(test)
	return true
}

func goEnv(name string) (string, error) {
	out, err := exec.Command("go", "env", name).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
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

func installPlatform(bindingsRoot, libDir string, installHeader bool, baseURL, version string, p platform, checksums map[string]string, force bool) error {
	archive := fmt.Sprintf("libmoss-v%s-%s.tar.gz", version, p.triple)
	wantChecksum, ok := checksums[archive]
	if !ok {
		return fmt.Errorf("checksum not found for %s", archive)
	}

	destDir := libDir
	destLib := filepath.Join(destDir, p.libFile)
	headerPath := filepath.Join(bindingsRoot, "include", "libmoss.h")

	if !force {
		if err := verifyFileChecksum(destLib, wantChecksum); err == nil {
			if !installHeader {
				fmt.Printf("skip %s (already installed)\n", p.id)
				return nil
			}
			if _, err := os.Stat(headerPath); err == nil {
				fmt.Printf("skip %s (already installed)\n", p.id)
				return nil
			}
		}
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}
	if installHeader {
		if err := os.MkdirAll(filepath.Dir(headerPath), 0o755); err != nil {
			return err
		}
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

	if installHeader {
		if err := copyFile(filepath.Join(extractedRoot, "include", "libmoss.h"), headerPath); err != nil {
			return err
		}
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
