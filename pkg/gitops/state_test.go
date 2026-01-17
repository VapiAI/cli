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

func TestSaveAndLoadState(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitops-state-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	statePath := filepath.Join(tempDir, ".vapi-state.json")

	// Create and populate state
	state := NewStateFile()
	state.Assistants["my-assistant"] = "uuid-assistant-123"
	state.Tools["my-tool"] = "uuid-tool-456"
	state.StructuredOutputs["my-output"] = "uuid-output-789"

	// Save state
	if err := SaveState(statePath, state); err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Error("State file was not created")
	}

	// Load state
	loaded, err := LoadState(statePath)
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}

	// Verify loaded content
	if loaded.Assistants["my-assistant"] != "uuid-assistant-123" {
		t.Errorf("Assistants mismatch: expected uuid-assistant-123, got %s", loaded.Assistants["my-assistant"])
	}
	if loaded.Tools["my-tool"] != "uuid-tool-456" {
		t.Errorf("Tools mismatch: expected uuid-tool-456, got %s", loaded.Tools["my-tool"])
	}
	if loaded.StructuredOutputs["my-output"] != "uuid-output-789" {
		t.Errorf("StructuredOutputs mismatch: expected uuid-output-789, got %s", loaded.StructuredOutputs["my-output"])
	}
}

func TestLoadState_NonExistent(t *testing.T) {
	// Loading non-existent state should return empty state
	loaded, err := LoadState("/nonexistent/path/.vapi-state.json")
	if err != nil {
		t.Fatalf("LoadState should not error for non-existent file: %v", err)
	}

	if loaded.Assistants == nil {
		t.Error("Assistants map should be initialized")
	}
	if loaded.Tools == nil {
		t.Error("Tools map should be initialized")
	}
	if loaded.StructuredOutputs == nil {
		t.Error("StructuredOutputs map should be initialized")
	}

	if len(loaded.Assistants) != 0 {
		t.Error("Assistants should be empty")
	}
}

func TestLoadState_InvalidJSON(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitops-state-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	statePath := filepath.Join(tempDir, ".vapi-state.json")

	// Write invalid JSON
	if err := os.WriteFile(statePath, []byte("not valid json"), 0644); err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}

	// Load should fail
	_, err = LoadState(statePath)
	if err == nil {
		t.Error("LoadState should fail for invalid JSON")
	}
}

func TestSaveState_CreatesDirectory(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitops-state-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Path with nested directory that doesn't exist
	statePath := filepath.Join(tempDir, "nested", "dir", ".vapi-state.json")

	state := NewStateFile()
	state.Tools["test"] = "uuid-test"

	// Save should create directory
	if err := SaveState(statePath, state); err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Error("State file was not created in nested directory")
	}
}

func TestStateFile_Concurrent(t *testing.T) {
	state := NewStateFile()

	// Simulate concurrent operations (basic test)
	done := make(chan bool, 3)

	go func() {
		for i := 0; i < 100; i++ {
			state.SetUUID(ResourceTypeAssistants, "assistant", "uuid")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			state.SetUUID(ResourceTypeTools, "tool", "uuid")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			state.GetSection(ResourceTypeAssistants)
		}
		done <- true
	}()

	<-done
	<-done
	<-done
}
