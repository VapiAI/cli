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
	Long:  `Print detailed version information about the Vapi CLI.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("vapi version %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built at: %s\n", date)
		fmt.Printf("  built by: %s\n", builtBy)
		fmt.Printf("  go version: %s\n", runtime.Version())
		fmt.Printf("  platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

		// Show environment information for developers
		if cfg, err := config.LoadConfig(); err == nil {
			if !cfg.IsProduction() {
				fmt.Printf("  environment: %s\n", cfg.GetEnvironment())
				fmt.Printf("  api url: %s\n", cfg.GetAPIBaseURL())
			}
		}

		// Check for updates (non-blocking)
		go func() {
			if shouldCheckForUpdates() {
				if release, hasUpdate, err := checkForUpdates(); err == nil && hasUpdate {
					fmt.Printf("\n🚀 Update available: %s → %s\n", version, release.TagName)
					fmt.Printf("   Run 'vapi update' to upgrade\n")

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

func init() {
	rootCmd.AddCommand(versionCmd)
}
