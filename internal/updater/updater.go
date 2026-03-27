package updater

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/MobileOps-Team/mobileops-cli/internal/config"
)

const (
	githubRepo     = "MobileOps-Team/mobileops-cli"
	releasesURL    = "https://api.github.com/repos/" + githubRepo + "/releases/latest"
	cacheDuration  = 24 * time.Hour
	cacheFileName  = "version-check.json"
	httpTimeout    = 5 * time.Second
	updateTimeout  = 60 * time.Second
)

// versionCache holds the cached version check result.
type versionCache struct {
	LatestVersion string `json:"latest_version"`
	CheckedAt     string `json:"checked_at"`
}

// githubRelease is a minimal representation of a GitHub release.
type githubRelease struct {
	TagName string `json:"tag_name"`
}

// cacheFilePath returns the path to the version check cache file.
func cacheFilePath() string {
	return filepath.Join(config.ConfigDir(), cacheFileName)
}

// CheckVersionBackground performs a non-blocking version check.
// It returns the latest version string via the channel, or empty string if
// the check fails or the cache is fresh and current version is up to date.
func CheckVersionBackground(currentVersion string) <-chan string {
	ch := make(chan string, 1)
	go func() {
		defer func() {
			select {
			case ch <- "":
			default:
			}
		}()

		latest, err := cachedLatestVersion()
		if err == nil && latest != "" {
			if isNewer(latest, currentVersion) {
				ch <- latest
				return
			}
			return
		}

		latest, err = fetchLatestVersion(httpTimeout)
		if err != nil {
			return
		}

		_ = writeCache(latest)

		if isNewer(latest, currentVersion) {
			ch <- latest
			return
		}
	}()
	return ch
}

func FetchLatestVersion() (string, error) {
	return fetchLatestVersion(updateTimeout)
}

func fetchLatestVersion(timeout time.Duration) (string, error) {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(releasesURL)
	if err != nil {
		return "", fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse release info: %w", err)
	}

	version := strings.TrimPrefix(release.TagName, "v")
	if version == "" {
		return "", fmt.Errorf("no version found in release")
	}

	return version, nil
}

func cachedLatestVersion() (string, error) {
	data, err := os.ReadFile(cacheFilePath())
	if err != nil {
		return "", err
	}

	var cache versionCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return "", err
	}

	checkedAt, err := time.Parse(time.RFC3339, cache.CheckedAt)
	if err != nil {
		return "", err
	}

	if time.Since(checkedAt) > cacheDuration {
		return "", fmt.Errorf("cache expired")
	}

	return cache.LatestVersion, nil
}

func writeCache(version string) error {
	dir := config.ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	cache := versionCache{
		LatestVersion: version,
		CheckedAt:     time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFilePath(), data, 0644)
}

func isNewer(latest, current string) bool {
	latestParts := parseVersion(latest)
	currentParts := parseVersion(current)

	for i := 0; i < 3; i++ {
		if latestParts[i] > currentParts[i] {
			return true
		}
		if latestParts[i] < currentParts[i] {
			return false
		}
	}
	return false
}

func parseVersion(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	var result [3]int
	for i, p := range parts {
		if i >= 3 {
			break
		}
		fmt.Sscanf(p, "%d", &result[i])
	}
	return result
}

func SelfUpdate(currentVersion string) error {
	latest, err := FetchLatestVersion()
	if err != nil {
		return err
	}

	if !isNewer(latest, currentVersion) {
		fmt.Printf("Already up to date (v%s)\n", currentVersion)
		return nil
	}

	fmt.Printf("Updating v%s -> v%s ...\n", currentVersion, latest)

	osName := runtime.GOOS
	archName := runtime.GOARCH

	downloadURL := fmt.Sprintf(
		"https://github.com/%s/releases/download/v%s/mobileops-cli_%s_%s_%s.tar.gz",
		githubRepo, latest, latest, osName, archName,
	)

	client := &http.Client{Timeout: updateTimeout}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d (URL: %s)", resp.StatusCode, downloadURL)
	}

	binaryData, err := extractBinaryFromTarGz(resp.Body, "mobileops")
	if err != nil {
		return fmt.Errorf("failed to extract binary: %w", err)
	}

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to determine executable path: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	dir := filepath.Dir(execPath)
	tmpFile, err := os.CreateTemp(dir, "mobileops-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.Write(binaryData); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write new binary: %w", err)
	}
	tmpFile.Close()

	info, err := os.Stat(execPath)
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to stat current binary: %w", err)
	}
	if err := os.Chmod(tmpPath, info.Mode()); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	if err := os.Rename(tmpPath, execPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	_ = writeCache(latest)

	fmt.Printf("Successfully updated to v%s\n", latest)
	return nil
}

func extractBinaryFromTarGz(r io.Reader, binaryPrefix string) ([]byte, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to open gzip: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar: %w", err)
		}

		name := filepath.Base(header.Name)
		if header.Typeflag == tar.TypeReg && strings.HasPrefix(name, binaryPrefix) && !strings.Contains(name, ".") {
			data, err := io.ReadAll(tr)
			if err != nil {
				return nil, fmt.Errorf("failed to read binary from archive: %w", err)
			}
			return data, nil
		}
	}

	return nil, fmt.Errorf("binary not found in archive")
}
