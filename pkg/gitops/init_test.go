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
	"os"
	"path/filepath"
	"testing"
)

func TestInitProject(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitops-init-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize project
	if err := InitProject(tempDir); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	// Verify resources directories exist
	cfg := NewConfig(tempDir)
	for _, rt := range AllResourceTypes() {
		dir := cfg.GetResourceDir(rt)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Resource directory not created: %s", dir)
		}

		// Verify .gitkeep exists
		gitkeep := filepath.Join(dir, ".gitkeep")
		if _, err := os.Stat(gitkeep); os.IsNotExist(err) {
			t.Errorf(".gitkeep not created in: %s", dir)
		}
	}

	// Verify state file exists
	if _, err := os.Stat(cfg.StateFilePath); os.IsNotExist(err) {
		t.Error("State file not created")
	}

	// Verify .gitignore exists and has content
	gitignore := filepath.Join(tempDir, ".gitignore")
	if _, err := os.Stat(gitignore); os.IsNotExist(err) {
		t.Error(".gitignore not created")
	}

	content, err := os.ReadFile(gitignore)
	if err != nil {
		t.Fatalf("Failed to read .gitignore: %v", err)
	}

	gitignoreContent := string(content)
	if !containsStr(gitignoreContent, ".env") {
		t.Error(".gitignore should contain .env")
	}
	if !containsStr(gitignoreContent, "!.env.example") {
		t.Error(".gitignore should contain !.env.example")
	}

	// Verify example files were created
	exampleTool := filepath.Join(cfg.GetResourceDir(ResourceTypeTools), "example-get-weather.yaml")
	if _, err := os.Stat(exampleTool); os.IsNotExist(err) {
		t.Error("Example tool file not created")
	}

	exampleAssistant := filepath.Join(cfg.GetResourceDir(ResourceTypeAssistants), "example-assistant.yaml")
	if _, err := os.Stat(exampleAssistant); os.IsNotExist(err) {
		t.Error("Example assistant file not created")
	}
}

func TestInitProject_AlreadyInitialized(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitops-init-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// First init should succeed
	if err := InitProject(tempDir); err != nil {
		t.Fatalf("First InitProject failed: %v", err)
	}

	// Second init should also succeed (idempotent)
	// It will overwrite existing files
	if err := InitProject(tempDir); err != nil {
		t.Fatalf("Second InitProject failed: %v", err)
	}
}

func TestIsGitOpsProject(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitops-check-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Not a gitops project initially
	if IsGitOpsProject(tempDir) {
		t.Error("Should not be a GitOps project initially")
	}

	// Initialize project
	if err := InitProject(tempDir); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	// Should be a gitops project now
	if !IsGitOpsProject(tempDir) {
		t.Error("Should be a GitOps project after init")
	}
}

func TestCreateEnvExample(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitops-env-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .env.example
	if err := CreateEnvExample(tempDir); err != nil {
		t.Fatalf("CreateEnvExample failed: %v", err)
	}

	// Verify file exists
	envExample := filepath.Join(tempDir, ".env.example")
	if _, err := os.Stat(envExample); os.IsNotExist(err) {
		t.Error(".env.example not created")
	}

	// Verify content
	content, err := os.ReadFile(envExample)
	if err != nil {
		t.Fatalf("Failed to read .env.example: %v", err)
	}

	contentStr := string(content)
	if !containsStr(contentStr, "VAPI_TOKEN") {
		t.Error(".env.example should contain VAPI_TOKEN")
	}
	if !containsStr(contentStr, "VAPI_BASE_URL") {
		t.Error(".env.example should contain VAPI_BASE_URL")
	}
}

func TestUpdateGitignore_Existing(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitops-gitignore-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	gitignorePath := filepath.Join(tempDir, ".gitignore")

	// Create existing .gitignore
	existingContent := "node_modules/\n*.log\n"
	if err := os.WriteFile(gitignorePath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("Failed to write initial .gitignore: %v", err)
	}

	// Initialize project (which updates .gitignore)
	if err := InitProject(tempDir); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	// Verify content
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("Failed to read .gitignore: %v", err)
	}

	contentStr := string(content)

	// Should contain original content
	if !containsStr(contentStr, "node_modules/") {
		t.Error(".gitignore should preserve original content")
	}

	// Should contain gitops entries
	if !containsStr(contentStr, ".env") {
		t.Error(".gitignore should contain .env")
	}
}

func TestUpdateGitignore_AlreadyHasEntries(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitops-gitignore-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	gitignorePath := filepath.Join(tempDir, ".gitignore")

	// Create .gitignore with gitops entries already
	existingContent := "# Vapi GitOps\n.env\n.env.*\n!.env.example\n"
	if err := os.WriteFile(gitignorePath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("Failed to write initial .gitignore: %v", err)
	}

	originalLen := len(existingContent)

	// Initialize project (which updates .gitignore)
	if err := InitProject(tempDir); err != nil {
		t.Fatalf("InitProject failed: %v", err)
	}

	// Verify content wasn't duplicated
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("Failed to read .gitignore: %v", err)
	}

	// File should not have grown significantly (entries should not be duplicated)
	if len(content) > originalLen*2 {
		t.Error(".gitignore entries should not be duplicated")
	}
}
