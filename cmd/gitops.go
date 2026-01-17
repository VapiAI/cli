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
*/

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/VapiAI/cli/pkg/config"
	"github.com/VapiAI/cli/pkg/gitops"
)

var gitopsCmd = &cobra.Command{
	Use:   "gitops",
	Short: "Manage Vapi resources via GitOps",
	Long: `Manage Vapi resources (Assistants, Tools, Structured Outputs) declaratively via YAML files.

GitOps enables version control, code review, and team collaboration for your Vapi configuration.

Commands:
  init   - Initialize a new GitOps project structure
  pull   - Pull resources from Vapi to local YAML files
  apply  - Push local YAML changes to Vapi

Example workflow:
  vapi gitops init              # Create project structure
  vapi gitops pull              # Fetch existing resources
  # Edit YAML files...
  vapi gitops apply             # Push changes to Vapi`,
}

var gitopsInitCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a GitOps project structure",
	Long: `Initialize a new GitOps project with the required directory structure.

Creates:
  resources/
    assistants/     - Assistant YAML files
    tools/          - Tool YAML files
    structuredOutputs/ - Structured Output YAML files
  .vapi-state.json  - State file mapping resource IDs to UUIDs
  .gitignore        - Updated with gitops entries
  .env.example      - Example environment file`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGitopsInit,
}

var gitopsPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull resources from Vapi to local YAML files",
	Long: `Pull all resources from your Vapi account and save them as local YAML files.

This will:
- Fetch all Assistants, Tools, and Structured Outputs from Vapi
- Save them as YAML files in the resources/ directory
- Update .vapi-state.json with UUID mappings
- Resolve cross-references to use resource IDs instead of UUIDs

Existing files will be overwritten with the latest data from Vapi.`,
	RunE: runGitopsPull,
}

var gitopsApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply local YAML changes to Vapi",
	Long: `Apply local YAML resource files to your Vapi account.

This will:
- Load all YAML files from the resources/ directory
- Create new resources that don't exist in Vapi
- Update existing resources with local changes
- Delete resources that were removed locally (with confirmation)
- Resolve cross-references (e.g., tool IDs in assistants)

Resources are applied in dependency order:
  1. Tools
  2. Structured Outputs
  3. Assistants`,
	RunE: runGitopsApply,
}

func init() {
	rootCmd.AddCommand(gitopsCmd)
	gitopsCmd.AddCommand(gitopsInitCmd)
	gitopsCmd.AddCommand(gitopsPullCmd)
	gitopsCmd.AddCommand(gitopsApplyCmd)
}

func runGitopsInit(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	// Make path absolute
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if already initialized
	if gitops.IsGitOpsProject(absPath) {
		return fmt.Errorf("GitOps project already initialized at %s", absPath)
	}

	fmt.Println("🚀 Initializing Vapi GitOps project...")
	fmt.Println()

	// Create project structure
	if err := gitops.InitProject(absPath); err != nil {
		return err
	}

	// Create .env.example
	if err := gitops.CreateEnvExample(absPath); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("✅ GitOps project initialized!")
	fmt.Println()
	fmt.Println("📝 Next steps:")
	fmt.Println()
	fmt.Println("1. Configure your API key:")
	fmt.Println("   cp .env.example .env")
	fmt.Println("   # Edit .env and add your VAPI_TOKEN")
	fmt.Println()
	fmt.Println("2. Pull existing resources from Vapi:")
	fmt.Println("   vapi gitops pull")
	fmt.Println()
	fmt.Println("3. Or start creating resources in the YAML files")
	fmt.Println("   and apply them:")
	fmt.Println("   vapi gitops apply")

	return nil
}

func runGitopsPull(cmd *cobra.Command, args []string) error {
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if GitOps project
	if !gitops.IsGitOpsProject(projectPath) {
		return fmt.Errorf("not a GitOps project (run 'vapi gitops init' first)")
	}

	// Get API key
	apiKey, err := getAPIKey()
	if err != nil {
		return err
	}

	// Create config
	cfg := gitops.NewConfig(projectPath)
	cfg.APIKey = apiKey

	// Load CLI config to get base URL
	cliConfig, err := config.LoadConfig()
	if err == nil && cliConfig.GetAPIBaseURL() != "" {
		cfg.APIBaseURL = cliConfig.GetAPIBaseURL()
	}

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("🔄 Vapi GitOps Pull")
	fmt.Printf("   API: %s\n", cfg.APIBaseURL)
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// Create and run pull engine
	engine, err := gitops.NewPullEngine(cfg)
	if err != nil {
		return err
	}

	ctx := context.Background()
	stats, err := engine.Pull(ctx)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("✅ Pull complete!")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("📋 Summary:")
	for rt, s := range stats {
		fmt.Printf("   %s: %d new, %d existing\n", rt, s.Created, s.Updated)
	}

	return nil
}

func runGitopsApply(cmd *cobra.Command, args []string) error {
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if GitOps project
	if !gitops.IsGitOpsProject(projectPath) {
		return fmt.Errorf("not a GitOps project (run 'vapi gitops init' first)")
	}

	// Get API key
	apiKey, err := getAPIKey()
	if err != nil {
		return err
	}

	// Create config
	cfg := gitops.NewConfig(projectPath)
	cfg.APIKey = apiKey

	// Load CLI config to get base URL
	cliConfig, err := config.LoadConfig()
	if err == nil && cliConfig.GetAPIBaseURL() != "" {
		cfg.APIBaseURL = cliConfig.GetAPIBaseURL()
	}

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("🚀 Vapi GitOps Apply")
	fmt.Printf("   API: %s\n", cfg.APIBaseURL)
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// Create and run apply engine
	engine, err := gitops.NewApplyEngine(cfg)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := engine.Apply(ctx); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("✅ Apply complete!")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Printf("📋 Summary: %s\n", engine.Summary())

	return nil
}

// getAPIKey retrieves the API key from environment or config.
func getAPIKey() (string, error) {
	// Try environment variable first
	if key := os.Getenv("VAPI_TOKEN"); key != "" {
		return key, nil
	}

	// Try loading from .env file
	envFile := ".env"
	if data, err := os.ReadFile(envFile); err == nil {
		if key := parseEnvVar(string(data), "VAPI_TOKEN"); key != "" {
			return key, nil
		}
	}

	// Try CLI config
	cfg, err := config.LoadConfig()
	if err == nil && cfg.APIKey != "" {
		return cfg.APIKey, nil
	}

	return "", fmt.Errorf("API key not found. Set VAPI_TOKEN in .env or run 'vapi login'")
}

// parseEnvVar extracts a variable value from .env content.
func parseEnvVar(content, varName string) string {
	lines := splitLines(content)
	prefix := varName + "="

	for _, line := range lines {
		line = trimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if hasPrefix(line, prefix) {
			value := line[len(prefix):]
			// Remove quotes if present
			if len(value) >= 2 {
				if (value[0] == '"' && value[len(value)-1] == '"') ||
					(value[0] == '\'' && value[len(value)-1] == '\'') {
					value = value[1 : len(value)-1]
				}
			}
			return value
		}
	}
	return ""
}

// splitLines splits content into lines.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// trimSpace removes leading and trailing whitespace.
func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// hasPrefix checks if s starts with prefix.
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
