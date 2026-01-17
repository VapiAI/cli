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

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Test_Resource", "test-resource"},
		{"My  Assistant!", "my-assistant"},
		{"UPPERCASE", "uppercase"},
		{"special@#$chars", "special-chars"},
		{"leading---dashes", "leading-dashes"},
		{"trailing---", "trailing"},
		{"123-numbers", "123-numbers"},
		{"multiple___underscores", "multiple-underscores"},
		{"", ""},
	}

	for _, tt := range tests {
		result := Slugify(tt.input)
		if result != tt.expected {
			t.Errorf("Slugify(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestGenerateUniqueResourceID(t *testing.T) {
	existing := map[string]bool{
		"my-resource":   true,
		"my-resource-1": true,
	}

	// Test generating ID for new resource
	result := GenerateUniqueResourceID("New Resource", existing)
	if result != "new-resource" {
		t.Errorf("GenerateUniqueResourceID for new resource = %q, expected %q", result, "new-resource")
	}

	// Test generating ID for existing resource (should get suffix)
	result = GenerateUniqueResourceID("My Resource", existing)
	if result != "my-resource-2" {
		t.Errorf("GenerateUniqueResourceID for existing = %q, expected %q", result, "my-resource-2")
	}
}

func TestWriteAndLoadResourceFile(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitops-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := NewConfig(tempDir)

	// Create tools directory
	toolsDir := cfg.GetResourceDir(ResourceTypeTools)
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		t.Fatalf("Failed to create tools dir: %v", err)
	}

	// Write a test resource
	testData := map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "test_function",
			"description": "A test function",
		},
	}

	filePath, err := WriteResourceFile(cfg, ResourceTypeTools, "test-tool", testData)
	if err != nil {
		t.Fatalf("WriteResourceFile failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("File was not created at %s", filePath)
	}

	// Read back the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Verify content contains expected values
	contentStr := string(content)
	if !containsStr(contentStr, "test_function") {
		t.Error("File content should contain 'test_function'")
	}
	if !containsStr(contentStr, "A test function") {
		t.Error("File content should contain 'A test function'")
	}
}

func TestLoadResources(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitops-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := NewConfig(tempDir)

	// Create tools directory with a test file
	toolsDir := cfg.GetResourceDir(ResourceTypeTools)
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		t.Fatalf("Failed to create tools dir: %v", err)
	}

	// Create a nested directory
	nestedDir := filepath.Join(toolsDir, "production")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	// Write test files
	tool1 := `type: function
function:
  name: tool1
  description: First tool
`
	tool2 := `type: function
function:
  name: tool2
  description: Second tool
`

	if err := os.WriteFile(filepath.Join(toolsDir, "tool1.yaml"), []byte(tool1), 0644); err != nil {
		t.Fatalf("Failed to write tool1: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nestedDir, "tool2.yaml"), []byte(tool2), 0644); err != nil {
		t.Fatalf("Failed to write tool2: %v", err)
	}

	// Load resources
	resources, err := LoadResources(cfg, ResourceTypeTools)
	if err != nil {
		t.Fatalf("LoadResources failed: %v", err)
	}

	if len(resources) != 2 {
		t.Errorf("LoadResources returned %d resources, expected 2", len(resources))
	}

	// Verify resource IDs
	ids := make(map[string]bool)
	for _, r := range resources {
		ids[r.ResourceID] = true
	}

	if !ids["tool1"] {
		t.Error("Missing resource ID 'tool1'")
	}
	if !ids["production-tool2"] {
		t.Error("Missing resource ID 'production-tool2' (nested should include folder prefix)")
	}
}

func TestLoadAllResources(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitops-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := NewConfig(tempDir)

	// Create all resource directories
	for _, rt := range AllResourceTypes() {
		dir := cfg.GetResourceDir(rt)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir for %s: %v", rt, err)
		}
	}

	// Write test files
	assistantYAML := `name: Test Assistant
model:
  provider: openai
  model: gpt-4o
`
	toolYAML := `type: function
function:
  name: test_tool
`
	outputYAML := `name: Test Output
schema:
  type: object
`

	if err := os.WriteFile(
		filepath.Join(cfg.GetResourceDir(ResourceTypeAssistants), "test-assistant.yaml"),
		[]byte(assistantYAML), 0644,
	); err != nil {
		t.Fatalf("Failed to write assistant: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(cfg.GetResourceDir(ResourceTypeTools), "test-tool.yaml"),
		[]byte(toolYAML), 0644,
	); err != nil {
		t.Fatalf("Failed to write tool: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(cfg.GetResourceDir(ResourceTypeStructuredOutputs), "test-output.yaml"),
		[]byte(outputYAML), 0644,
	); err != nil {
		t.Fatalf("Failed to write output: %v", err)
	}

	// Load all resources
	loaded, err := LoadAllResources(cfg)
	if err != nil {
		t.Fatalf("LoadAllResources failed: %v", err)
	}

	if len(loaded.Assistants) != 1 {
		t.Errorf("Expected 1 assistant, got %d", len(loaded.Assistants))
	}
	if len(loaded.Tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(loaded.Tools))
	}
	if len(loaded.StructuredOutputs) != 1 {
		t.Errorf("Expected 1 structured output, got %d", len(loaded.StructuredOutputs))
	}
}

func TestGetResourceIDsFromFiles(t *testing.T) {
	files := []*ResourceFile{
		{ResourceID: "resource-1"},
		{ResourceID: "resource-2"},
		{ResourceID: "resource-3"},
	}

	ids := GetResourceIDsFromFiles(files)

	if len(ids) != 3 {
		t.Errorf("GetResourceIDsFromFiles returned %d IDs, expected 3", len(ids))
	}

	expected := map[string]bool{
		"resource-1": true,
		"resource-2": true,
		"resource-3": true,
	}

	for _, id := range ids {
		if !expected[id] {
			t.Errorf("Unexpected ID: %s", id)
		}
	}
}

// Helper function
func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
