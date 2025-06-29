/*
Copyright © 2025 Vapi, Inc.

Licensed under the MIT License (the "License");
you may not use this file except in compliance with the License.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

Authors:

	Dan Goosewin <dan@vapi.ai>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	Draft       bool   `json:"draft"`
	Prerelease  bool   `json:"prerelease"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the Vapi CLI to the latest version",
	Long: `Check for and install the latest version of the Vapi CLI.

This command will:
- Check GitHub releases for the latest version
- Compare with your current version
- Download and install the update if available
- Preserve your current configuration`,
	RunE: runUpdateCommand,
}

var checkUpdateCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for available updates",
	Long:  `Check if a newer version of the Vapi CLI is available without installing it.`,
	RunE:  runCheckUpdateCommand,
}

func runUpdateCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("🔄 Checking for Vapi CLI updates...")

	latestRelease, hasUpdate, err := checkForUpdates()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !hasUpdate {
		fmt.Printf("✅ You're already running the latest version: %s\n", version)
		return nil
	}

	fmt.Printf("🆕 New version available: %s (current: %s)\n", latestRelease.TagName, version)
	fmt.Printf("📅 Released: %s\n", formatReleaseDate(latestRelease.PublishedAt))

	if latestRelease.Body != "" {
		fmt.Println("\n📝 Release notes:")
		// Show only the first few lines of release notes to avoid overwhelming output
		lines := strings.Split(strings.TrimSpace(latestRelease.Body), "\n")
		maxLines := 10
		for i, line := range lines {
			if i >= maxLines {
				fmt.Println("   ... (see full release notes at: https://github.com/VapiAI/cli/releases)")
				break
			}
			fmt.Printf("   %s\n", line)
		}
	}

	fmt.Println("\n🔄 Installing update...")

	if err := installUpdate(latestRelease); err != nil {
		return fmt.Errorf("failed to install update: %w", err)
	}

	fmt.Printf("✅ Successfully updated to version %s!\n", latestRelease.TagName)
	fmt.Println("🔄 Please restart your terminal or run 'hash -r' to use the new version.")

	return nil
}

func runCheckUpdateCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Checking for updates...")

	latestRelease, hasUpdate, err := checkForUpdates()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !hasUpdate {
		fmt.Printf("✅ You're running the latest version: %s\n", version)
		return nil
	}

	fmt.Printf("🆕 New version available: %s (current: %s)\n", latestRelease.TagName, version)
	fmt.Printf("📅 Released: %s\n", formatReleaseDate(latestRelease.PublishedAt))

	if latestRelease.Body != "" {
		fmt.Println("\n📝 Release notes:")
		// Show only the first few lines of release notes
		lines := strings.Split(strings.TrimSpace(latestRelease.Body), "\n")
		maxLines := 5
		for i, line := range lines {
			if i >= maxLines {
				fmt.Println("   ... (see full release notes at: https://github.com/VapiAI/cli/releases)")
				break
			}
			fmt.Printf("   %s\n", line)
		}
	}

	fmt.Println("\n💡 Run 'vapi update' to install the latest version.")

	return nil
}

func checkForUpdates() (*GitHubRelease, bool, error) {
	// Get latest release from GitHub API
	resp, err := http.Get("https://api.github.com/repos/VapiAI/cli/releases/latest")
	if err != nil {
		return nil, false, err
	}
	defer func() {
		_ = resp.Body.Close() // Ignore close errors for API calls
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, false, err
	}

	// Skip draft and prerelease versions
	if release.Draft || release.Prerelease {
		return nil, false, fmt.Errorf("latest release is a draft or prerelease")
	}

	// Compare versions using semantic version comparison
	currentVersion := strings.TrimPrefix(version, "v")
	latestVersion := strings.TrimPrefix(release.TagName, "v")

	// Skip if current version is "dev" (development build)
	if currentVersion == "dev" {
		hasUpdate := latestVersion != "dev"
		return &release, hasUpdate, nil
	}

	// Skip if latest version is "dev" (shouldn't happen with GitHub releases)
	if latestVersion == "dev" {
		return &release, false, nil
	}

	// Compare semantic versions
	hasUpdate, err := isNewerVersion(latestVersion, currentVersion)
	if err != nil {
		// Fallback to string comparison if parsing fails
		hasUpdate = latestVersion != currentVersion
	}

	return &release, hasUpdate, nil
}

// isNewerVersion compares two semantic version strings and returns true if newer > current
func isNewerVersion(newer, current string) (bool, error) {
	newerParts, err := parseVersion(newer)
	if err != nil {
		return false, err
	}

	currentParts, err := parseVersion(current)
	if err != nil {
		return false, err
	}

	// Compare major.minor.patch
	for i := 0; i < 3; i++ {
		if newerParts[i] > currentParts[i] {
			return true, nil
		}
		if newerParts[i] < currentParts[i] {
			return false, nil
		}
	}

	// Versions are equal
	return false, nil
}

// parseVersion parses a semantic version string (e.g., "1.2.3") into [major, minor, patch]
func parseVersion(version string) ([3]int, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return [3]int{}, fmt.Errorf("invalid version format: %s", version)
	}

	var result [3]int
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return [3]int{}, fmt.Errorf("invalid version number: %s", part)
		}
		result[i] = num
	}

	return result, nil
}

func installUpdate(release *GitHubRelease) error {
	// Find the appropriate asset for the current platform
	// Note: GoReleaser uses "cli" as the project name, not "vapi"
	assetName := fmt.Sprintf("cli_%s_%s", getOSName(), getArchName())
	if runtime.GOOS == "windows" {
		assetName += ".zip"
	} else {
		assetName += ".tar.gz"
	}

	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		// Debug: show available assets
		fmt.Printf("Available assets:\n")
		for _, asset := range release.Assets {
			fmt.Printf("  - %s\n", asset.Name)
		}
		return fmt.Errorf("no compatible binary found for %s/%s (looking for: %s)", runtime.GOOS, runtime.GOARCH, assetName)
	}

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "vapi-update-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
	}()

	// Download the archive
	fmt.Println("📥 Downloading update...")
	archivePath := tmpDir + "/" + getArchiveFileName()
	if err := downloadFile(downloadURL, archivePath); err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}

	// Extract the archive
	fmt.Println("📦 Extracting update...")
	if err := extractArchive(archivePath, tmpDir); err != nil {
		return fmt.Errorf("failed to extract update: %w", err)
	}

	// Find the binary in the extracted files
	binaryPath := tmpDir + "/vapi"
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found in downloaded archive")
	}

	// Get current executable path
	currentPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Replace the current binary
	fmt.Println("🔄 Installing update...")
	if err := replaceExecutable(binaryPath, currentPath); err != nil {
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	return nil
}

func getOSName() string {
	switch runtime.GOOS {
	case "darwin":
		return "Darwin"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return runtime.GOOS
	}
}

func getArchName() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64"
	case "386":
		return "i386"
	case "arm64":
		return "arm64"
	case "arm":
		return "arm"
	default:
		return runtime.GOARCH
	}
}

func getArchiveFileName() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("cli_%s_%s.zip", getOSName(), getArchName())
	}
	return fmt.Sprintf("cli_%s_%s.tar.gz", getOSName(), getArchName())
}

func downloadFile(url, filepath string) error {
	// #nosec G107 - URL is from GitHub releases API, considered safe
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close() // Ignore close errors for downloads
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// #nosec G304 - filepath is controlled by this function
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close() // Ignore close errors for temporary files
	}()

	_, err = io.Copy(file, resp.Body)
	return err
}

func extractArchive(archivePath, destDir string) error {
	if runtime.GOOS == "windows" {
		// For Windows, we'd need to implement ZIP extraction
		// For now, provide instructions
		return fmt.Errorf("automatic update not yet supported on Windows. Please download manually from: https://github.com/VapiAI/cli/releases")
	}

	// Use tar for Unix-like systems
	cmd := exec.Command("tar", "-xzf", archivePath, "-C", destDir)
	return cmd.Run()
}

func replaceExecutable(newPath, currentPath string) error {
	// Make the new binary executable (executable permissions required)
	// #nosec G302 - executable files need 0755 permissions
	if err := os.Chmod(newPath, 0o755); err != nil {
		return err
	}

	// On Unix-like systems, we can replace the file directly
	if runtime.GOOS != "windows" {
		return os.Rename(newPath, currentPath)
	}

	// On Windows, we need to handle the case where the current executable might be locked
	backupPath := currentPath + ".backup"
	if err := os.Rename(currentPath, backupPath); err != nil {
		return err
	}

	if err := os.Rename(newPath, currentPath); err != nil {
		// Try to restore backup
		_ = os.Rename(backupPath, currentPath) // Ignore restore errors
		return err
	}

	// Remove backup
	_ = os.Remove(backupPath) // Ignore cleanup errors
	return nil
}

func formatReleaseDate(dateStr string) string {
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t.Format("January 2, 2006")
	}
	return dateStr
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.AddCommand(checkUpdateCmd)
}
