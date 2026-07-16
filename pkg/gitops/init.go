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

package gitops

import (
	"fmt"
	"os"
	"path/filepath"
)

// InitProject initializes a new GitOps project structure.
func InitProject(projectPath string) error {
	config := NewConfig(projectPath)

	// Create resources directories
	for _, rt := range AllResourceTypes() {
		dir := config.GetResourceDir(rt)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", rt, err)
		}
		fmt.Printf("  ✓ Created %s/\n", filepath.Join(DefaultResourcesDir, ResourceDirs[rt]))

		// Create .gitkeep to ensure empty directories are tracked
		gitkeep := filepath.Join(dir, ".gitkeep")
		if err := os.WriteFile(gitkeep, []byte{}, 0644); err != nil {
			return fmt.Errorf("failed to create .gitkeep: %w", err)
		}
	}

	// Create empty state file
	state := NewStateFile()
	if err := SaveState(config.StateFilePath, state); err != nil {
		return fmt.Errorf("failed to create state file: %w", err)
	}
	fmt.Printf("  ✓ Created %s\n", DefaultStateFile)

	// Create or update .gitignore
	if err := updateGitignore(projectPath); err != nil {
		return fmt.Errorf("failed to update .gitignore: %w", err)
	}
	fmt.Println("  ✓ Updated .gitignore")

	// Create example assistant
	if err := createExampleFiles(config); err != nil {
		return fmt.Errorf("failed to create example files: %w", err)
	}

	return nil
}

// updateGitignore adds gitops-related entries to .gitignore.
func updateGitignore(projectPath string) error {
	gitignorePath := filepath.Join(projectPath, ".gitignore")

	entries := []string{
		"# Vapi GitOps",
		".env",
		".env.*",
		"!.env.example",
	}

	// Read existing content
	existing := ""
	if data, err := os.ReadFile(gitignorePath); err == nil {
		existing = string(data)
	}

	// Check if already contains gitops entries
	if len(existing) > 0 && containsGitopsEntries(existing) {
		return nil
	}

	// Append entries
	content := existing
	if len(content) > 0 && content[len(content)-1] != '\n' {
		content += "\n"
	}
	content += "\n"
	for _, entry := range entries {
		content += entry + "\n"
	}

	return os.WriteFile(gitignorePath, []byte(content), 0644)
}

// containsGitopsEntries checks if .gitignore already has gitops entries.
func containsGitopsEntries(content string) bool {
	return len(content) > 0 && (contains(content, "# Vapi GitOps") || contains(content, ".env.*"))
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// createExampleFiles creates example resource files.
func createExampleFiles(config *Config) error {
	// Example tool
	toolData := map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "get_weather",
			"description": "Get the current weather for a location",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "The city name",
					},
				},
				"required": []string{"location"},
			},
		},
		"server": map[string]interface{}{
			"url": "https://your-api.com/weather",
		},
	}

	toolPath, err := WriteResourceFile(config, ResourceTypeTools, "example-get-weather", toolData)
	if err != nil {
		return err
	}
	fmt.Printf("  ✓ Created example tool: %s\n", toolPath)

	// Example assistant
	assistantData := map[string]interface{}{
		"name": "Example Assistant",
		"model": map[string]interface{}{
			"provider": "openai",
			"model":    "gpt-4o",
			"messages": []map[string]interface{}{
				{
					"role":    "system",
					"content": "You are a helpful assistant. Be concise and friendly.",
				},
			},
			"toolIds": []string{
				"example-get-weather",
			},
		},
		"firstMessage": "Hello! How can I help you today?",
		"voice": map[string]interface{}{
			"provider": "11labs",
			"voiceId":  "21m00Tcm4TlvDq8ikWAM",
		},
	}

	assistantPath, err := WriteResourceFile(config, ResourceTypeAssistants, "example-assistant", assistantData)
	if err != nil {
		return err
	}
	fmt.Printf("  ✓ Created example assistant: %s\n", assistantPath)

	return nil
}

// CreateEnvExample creates an example .env file.
func CreateEnvExample(projectPath string) error {
	envExamplePath := filepath.Join(projectPath, ".env.example")

	content := `# Vapi GitOps Configuration
# Copy this file to .env and fill in your values

# Your Vapi API token (required)
# Get it from https://dashboard.vapi.ai/account
VAPI_TOKEN=your-token-here

# API Base URL (optional, defaults to https://api.vapi.ai)
# VAPI_BASE_URL=https://api.vapi.ai
`

	if err := os.WriteFile(envExamplePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create .env.example: %w", err)
	}

	fmt.Println("  ✓ Created .env.example")
	return nil
}

// IsGitOpsProject checks if a directory is a GitOps project.
func IsGitOpsProject(projectPath string) bool {
	config := NewConfig(projectPath)

	// Check for resources directory
	if _, err := os.Stat(config.ResourcesDir); os.IsNotExist(err) {
		return false
	}

	// Check for at least one resource type directory
	for _, rt := range AllResourceTypes() {
		if _, err := os.Stat(config.GetResourceDir(rt)); err == nil {
			return true
		}
	}

	return false
}
