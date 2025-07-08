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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/VapiAI/cli/pkg/config"
)

var (
	// Version information set by main
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// SetVersion sets the version information for the CLI
func SetVersion(v, c, d, b string) {
	version = v
	commit = c
	date = d
	builtBy = b
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print detailed version information about the Vapi CLI including MCP server compatibility.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("🔧 Vapi CLI v%s\n", version)
		fmt.Printf("\n📊 Build Information:\n")
		fmt.Printf("  Commit: %s\n", commit)
		fmt.Printf("  Built at: %s\n", date)
		fmt.Printf("  Built by: %s\n", builtBy)
		fmt.Printf("  Go version: %s\n", runtime.Version())
		fmt.Printf("  Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

		// Show MCP server compatibility
		fmt.Printf("\n🧠 MCP Integration:\n")
		fmt.Printf("  Compatible MCP server: @vapi-ai/mcp-docs-server@%s\n", version)
		fmt.Printf("  MCP setup: vapi mcp setup\n")
		fmt.Printf("  MCP status: vapi mcp status\n")

		// Show environment information for developers
		if cfg, err := config.LoadConfig(); err == nil {
			if !cfg.IsProduction() {
				fmt.Printf("\n🛠️  Development Environment:\n")
				fmt.Printf("  Environment: %s\n", cfg.GetEnvironment())
				fmt.Printf("  API URL: %s\n", cfg.GetAPIBaseURL())
			}
		}

		// Show installation info
		fmt.Printf("\n📦 Installation:\n")
		execPath, err := os.Executable()
		if err == nil {
			fmt.Printf("  Executable: %s\n", execPath)
		}

		// Check if man pages are available
		if runtime.GOOS != "windows" {
			fmt.Printf("  Manual pages: ")
			if isManPageInstalled() {
				fmt.Printf("✅ Available (man vapi)\n")
			} else {
				fmt.Printf("❌ Not installed (run 'vapi install-man-pages')\n")
			}
		}

		// Check for updates (non-blocking)
		go func() {
			if shouldCheckForUpdates() {
				if release, hasUpdate, err := checkForUpdates(); err == nil && hasUpdate {
					fmt.Printf("\n🚀 Update Available:\n")
					fmt.Printf("  Current: v%s\n", version)
					fmt.Printf("  Latest: %s\n", release.TagName)
					fmt.Printf("  Run: vapi update\n")

					// Update last check time
					updateLastCheckTime()
				}
			}
		}()

		return nil
	},
}

// shouldCheckForUpdates returns true if we should check for updates (once per day)
func shouldCheckForUpdates() bool {
	configDir, err := getConfigDir()
	if err != nil {
		return true // Check if we can't determine config dir
	}

	lastCheckFile := filepath.Join(configDir, ".vapi-last-update-check")
	info, err := os.Stat(lastCheckFile)
	if err != nil {
		return true // Check if file doesn't exist
	}

	// Check if more than 24 hours have passed
	return time.Since(info.ModTime()) > 24*time.Hour
}

// updateLastCheckTime updates the timestamp of the last update check
func updateLastCheckTime() {
	configDir, err := getConfigDir()
	if err != nil {
		return
	}

	// Ensure config directory exists (0750 permissions for security)
	// #nosec G301 - config directory needs to be accessible by user
	if err := os.MkdirAll(configDir, 0o750); err != nil {
		return
	}

	lastCheckFile := filepath.Join(configDir, ".vapi-last-update-check")
	// #nosec G304 - lastCheckFile path is controlled by this function
	file, err := os.Create(lastCheckFile)
	if err != nil {
		return
	}
	_ = file.Close() // Ignore close errors for background operations
}

// getConfigDir returns the configuration directory for the CLI
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".vapi-cli"), nil
}

// isManPageInstalled checks if the vapi man page is installed
func isManPageInstalled() bool {
	// Check common man page locations
	manDirs := []string{
		"/usr/local/share/man/man1",
		"/usr/share/man/man1",
		"/opt/homebrew/share/man/man1",
	}

	for _, dir := range manDirs {
		manFile := filepath.Join(dir, "vapi.1")
		if _, err := os.Stat(manFile); err == nil {
			return true
		}
	}

	return false
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
