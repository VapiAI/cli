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

	"os/exec"

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

// installManPagesCmd installs manual pages for the current user
var installManPagesCmd = &cobra.Command{
	Use:   "install-man-pages",
	Short: "Install manual pages for the Vapi CLI",
	Long: `Install Unix manual pages for the Vapi CLI to enable 'man vapi' command.

This command will:
1. Generate the latest manual pages for all commands
2. Install them to the appropriate system location
3. Update the manual database

Requires administrative privileges for system-wide installation.`,
	RunE: installManualPages,
}

func installManualPages(cmd *cobra.Command, args []string) error {
	fmt.Printf("📖 Installing Vapi CLI manual pages...\n")

	// Create a temporary directory for generated pages
	tmpDir := filepath.Join(os.TempDir(), "vapi-man-pages")
	if err := os.MkdirAll(tmpDir, 0o750); err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up temp directory

	// Generate man pages
	fmt.Printf("   📝 Generating manual pages...\n")

	header := &doc.GenManHeader{
		Title:   "VAPI",
		Section: "1",
		Source:  fmt.Sprintf("Vapi CLI %s", versionpkg.Get()),
		Manual:  "Vapi CLI Manual",
		Date:    &[]time.Time{time.Now()}[0],
	}

	if err := doc.GenManTree(rootCmd, header, tmpDir); err != nil {
		return fmt.Errorf("failed to generate man pages: %w", err)
	}

	// Find appropriate installation directory
	manDir := getManPageInstallDir()
	if manDir == "" {
		return fmt.Errorf("could not find suitable man page directory")
	}

	// Check if directory exists and is writable
	if err := ensureManDirWritable(manDir); err != nil {
		fmt.Printf("   ⚠️  Need administrative privileges to install to %s\n", manDir)
		fmt.Printf("   💡 Suggestions:\n")
		fmt.Printf("      • Run with sudo: sudo vapi install-man-pages\n")
		fmt.Printf("      • Or manually copy files from %s to %s\n", tmpDir, manDir)
		return err
	}

	// Copy man pages to installation directory
	fmt.Printf("   📋 Installing to %s...\n", manDir)

	files, err := filepath.Glob(filepath.Join(tmpDir, "*.1"))
	if err != nil {
		return fmt.Errorf("failed to find generated man pages: %w", err)
	}

	for _, file := range files {
		fileName := filepath.Base(file)
		destPath := filepath.Join(manDir, fileName)

		if err := copyFile(file, destPath); err != nil {
			return fmt.Errorf("failed to copy %s: %w", fileName, err)
		}

		// Set appropriate permissions
		if err := os.Chmod(destPath, 0o644); err != nil {
			fmt.Printf("   ⚠️  Warning: failed to set permissions on %s\n", fileName)
		}
	}

	// Update man database
	fmt.Printf("   🔄 Updating manual database...\n")
	if err := updateManDatabase(); err != nil {
		fmt.Printf("   ⚠️  Warning: failed to update man database: %v\n", err)
		fmt.Printf("   💡 You may need to run 'sudo mandb' manually\n")
	}

	fmt.Printf("   ✅ Installed %d manual page(s)\n", len(files))
	fmt.Println()
	fmt.Println("📚 Manual pages installed successfully!")
	fmt.Println()
	fmt.Println("🔍 Usage:")
	fmt.Println("     man vapi          # Main CLI manual")
	fmt.Println("     man vapi-call     # Call management")
	fmt.Println("     man vapi-mcp      # MCP integration")
	fmt.Println("     man vapi-assistant# Assistant management")
	fmt.Println()
	fmt.Println("💡 Try 'man vapi' to test the installation!")

	return nil
}

func getManPageInstallDir() string {
	// Check common man page directories in order of preference
	candidates := []string{
		"/usr/local/share/man/man1",    // Most common for user-installed tools
		"/opt/homebrew/share/man/man1", // Homebrew on Apple Silicon
		"/usr/share/man/man1",          // System-wide
	}

	for _, dir := range candidates {
		if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
			return dir
		}
	}

	return ""
}

func ensureManDirWritable(dir string) error {
	// Check if directory exists
	stat, err := os.Stat(dir)
	if err != nil {
		// Try to create the directory
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("directory does not exist and cannot be created: %w", err)
		}
		return nil
	}

	if !stat.IsDir() {
		return fmt.Errorf("path exists but is not a directory")
	}

	// Test write permissions by creating a temporary file
	testFile := filepath.Join(dir, ".vapi-test-write")
	if file, err := os.Create(testFile); err != nil {
		return fmt.Errorf("directory is not writable: %w", err)
	} else {
		file.Close()
		os.Remove(testFile) // Clean up test file
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = srcFile.WriteTo(dstFile)
	return err
}

func updateManDatabase() error {
	// Try to update the man database
	commands := [][]string{
		{"mandb"},        // Most common
		{"makewhatis"},   // Some systems
		{"catman", "-w"}, // Alternative
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err == nil {
			return nil // Success
		}
	}

	return fmt.Errorf("no suitable man database update command found")
}

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
	rootCmd.AddCommand(installManPagesCmd)
}
