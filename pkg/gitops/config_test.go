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

func TestNewConfig(t *testing.T) {
	cfg := NewConfig("/test/project")

	if cfg.ProjectPath != "/test/project" {
		t.Errorf("ProjectPath = %s, expected /test/project", cfg.ProjectPath)
	}

	expectedResources := filepath.Join("/test/project", DefaultResourcesDir)
	if cfg.ResourcesDir != expectedResources {
		t.Errorf("ResourcesDir = %s, expected %s", cfg.ResourcesDir, expectedResources)
	}

	expectedState := filepath.Join("/test/project", DefaultStateFile)
	if cfg.StateFilePath != expectedState {
		t.Errorf("StateFilePath = %s, expected %s", cfg.StateFilePath, expectedState)
	}

	if cfg.APIBaseURL != DefaultAPIBaseURL {
		t.Errorf("APIBaseURL = %s, expected %s", cfg.APIBaseURL, DefaultAPIBaseURL)
	}
}

func TestConfig_GetResourceDir(t *testing.T) {
	cfg := NewConfig("/test/project")

	tests := []struct {
		rt       ResourceType
		expected string
	}{
		{ResourceTypeAssistants, filepath.Join("/test/project", DefaultResourcesDir, "assistants")},
		{ResourceTypeTools, filepath.Join("/test/project", DefaultResourcesDir, "tools")},
		{ResourceTypeStructuredOutputs, filepath.Join("/test/project", DefaultResourcesDir, "structuredOutputs")},
	}

	for _, tt := range tests {
		result := cfg.GetResourceDir(tt.rt)
		if result != tt.expected {
			t.Errorf("GetResourceDir(%s) = %s, expected %s", tt.rt, result, tt.expected)
		}
	}
}

func TestConfig_ValidateForApply(t *testing.T) {
	// Create temp directory for testing
	tempDir, err := os.MkdirTemp("", "gitops-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := NewConfig(tempDir)
	cfg.APIKey = "test-key"

	// Create resources directory
	if err := os.MkdirAll(cfg.ResourcesDir, 0755); err != nil {
		t.Fatalf("Failed to create resources dir: %v", err)
	}

	// Should pass now
	if err := cfg.ValidateForApply(); err != nil {
		t.Errorf("ValidateForApply should pass with valid config, got: %v", err)
	}

	// Test missing API key
	cfg.APIKey = ""
	if err := cfg.ValidateForApply(); err == nil {
		t.Error("ValidateForApply should fail with empty API key")
	}
}

func TestGetAPIEndpoint(t *testing.T) {
	tests := []struct {
		rt       ResourceType
		expected string
	}{
		{ResourceTypeAssistants, "/assistant"},
		{ResourceTypeTools, "/tool"},
		{ResourceTypeStructuredOutputs, "/structured-output"},
	}

	for _, tt := range tests {
		result := GetAPIEndpoint(tt.rt)
		if result != tt.expected {
			t.Errorf("GetAPIEndpoint(%s) = %s, expected %s", tt.rt, result, tt.expected)
		}
	}
}

func TestIsExcludedOnPull(t *testing.T) {
	excludedKeys := []string{"id", "createdAt", "updatedAt", "orgId"}
	for _, key := range excludedKeys {
		if !IsExcludedOnPull(key) {
			t.Errorf("IsExcludedOnPull(%s) should return true", key)
		}
	}

	allowedKeys := []string{"name", "model", "voice", "type"}
	for _, key := range allowedKeys {
		if IsExcludedOnPull(key) {
			t.Errorf("IsExcludedOnPull(%s) should return false", key)
		}
	}
}

func TestIsExcludedOnUpdate(t *testing.T) {
	excludedKeys := []string{"id", "createdAt", "updatedAt", "orgId"}
	for _, key := range excludedKeys {
		if !IsExcludedOnUpdate(key) {
			t.Errorf("IsExcludedOnUpdate(%s) should return true", key)
		}
	}

	allowedKeys := []string{"name", "model", "voice", "type"}
	for _, key := range allowedKeys {
		if IsExcludedOnUpdate(key) {
			t.Errorf("IsExcludedOnUpdate(%s) should return false", key)
		}
	}
}

func TestRemoveExcludedKeys(t *testing.T) {
	data := map[string]interface{}{
		"id":        "uuid-123",
		"name":      "Test Resource",
		"createdAt": "2025-01-01",
		"updatedAt": "2025-01-02",
		"model":     "gpt-4o",
	}

	result := RemoveExcludedKeys(data, ResourceTypeAssistants)

	if _, exists := result["id"]; exists {
		t.Error("id should be removed")
	}
	if _, exists := result["createdAt"]; exists {
		t.Error("createdAt should be removed")
	}
	if _, exists := result["updatedAt"]; exists {
		t.Error("updatedAt should be removed")
	}
	if result["name"] != "Test Resource" {
		t.Error("name should be preserved")
	}
	if result["model"] != "gpt-4o" {
		t.Error("model should be preserved")
	}
}
