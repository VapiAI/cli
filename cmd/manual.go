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
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	versionpkg "github.com/VapiAI/cli/pkg/version"
)

// Generate manual pages for the Vapi CLI - hidden command for build process
var manualCmd = &cobra.Command{
	Use:    "manual",
	Short:  "Generate manual pages for Vapi CLI",
	Long:   `Generate Unix manual pages for the Vapi CLI and all subcommands.`,
	Hidden: true, // Hidden from normal help output
	RunE:   generateManualPages,
}

var manualOutputDir string

func generateManualPages(cmd *cobra.Command, args []string) error {
	fmt.Printf("📖 Generating manual pages for Vapi CLI v%s...\n", versionpkg.Get())

	// Ensure output directory exists
	if err := os.MkdirAll(manualOutputDir, 0o750); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Configure man page header
	header := &doc.GenManHeader{
		Title:   "VAPI",
		Section: "1", // User commands section
		Source:  fmt.Sprintf("Vapi CLI %s", versionpkg.Get()),
		Manual:  "Vapi CLI Manual",
		Date:    &[]time.Time{time.Now()}[0],
	}

	// Generate man pages for all commands
	fmt.Printf("   Output directory: %s\n", manualOutputDir)

	if err := doc.GenManTree(rootCmd, header, manualOutputDir); err != nil {
		return fmt.Errorf("failed to generate man pages: %w", err)
	}

	// List generated files
	files, err := filepath.Glob(filepath.Join(manualOutputDir, "*.1"))
	if err != nil {
		return fmt.Errorf("failed to list generated files: %w", err)
	}

	fmt.Printf("   Generated %d manual page(s):\n", len(files))
	for _, file := range files {
		fmt.Printf("     • %s\n", filepath.Base(file))
	}

	fmt.Println()
	fmt.Println("📝 Installation:")
	fmt.Printf("     sudo cp %s/*.1 /usr/local/share/man/man1/\n", manualOutputDir)
	fmt.Println("     sudo mandb  # Update man database")
	fmt.Println()
	fmt.Println("🔍 Usage:")
	fmt.Println("     man vapi")
	fmt.Println("     man vapi-assistant")
	fmt.Println("     man vapi-call")
	fmt.Println()

	return nil
}

func init() {
	manualCmd.Flags().StringVarP(&manualOutputDir, "output", "o", "./man",
		"Output directory for generated manual pages")

	rootCmd.AddCommand(manualCmd)
}
