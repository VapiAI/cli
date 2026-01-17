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
	"testing"
)

func TestNewStateFile(t *testing.T) {
	state := NewStateFile()

	if state.Assistants == nil {
		t.Error("Assistants map should be initialized")
	}
	if state.Tools == nil {
		t.Error("Tools map should be initialized")
	}
	if state.StructuredOutputs == nil {
		t.Error("StructuredOutputs map should be initialized")
	}
}

func TestStateFile_GetSection(t *testing.T) {
	state := NewStateFile()
	state.Assistants["test-assistant"] = "uuid-1"
	state.Tools["test-tool"] = "uuid-2"
	state.StructuredOutputs["test-output"] = "uuid-3"

	tests := []struct {
		rt       ResourceType
		expected string
	}{
		{ResourceTypeAssistants, "uuid-1"},
		{ResourceTypeTools, "uuid-2"},
		{ResourceTypeStructuredOutputs, "uuid-3"},
	}

	for _, tt := range tests {
		section := state.GetSection(tt.rt)
		if section == nil {
			t.Errorf("GetSection(%s) returned nil", tt.rt)
			continue
		}

		var key string
		switch tt.rt {
		case ResourceTypeAssistants:
			key = "test-assistant"
		case ResourceTypeTools:
			key = "test-tool"
		case ResourceTypeStructuredOutputs:
			key = "test-output"
		}

		if section[key] != tt.expected {
			t.Errorf("GetSection(%s)[%s] = %s, expected %s", tt.rt, key, section[key], tt.expected)
		}
	}
}

func TestStateFile_SetUUID(t *testing.T) {
	state := NewStateFile()

	state.SetUUID(ResourceTypeAssistants, "my-assistant", "uuid-a")
	state.SetUUID(ResourceTypeTools, "my-tool", "uuid-t")
	state.SetUUID(ResourceTypeStructuredOutputs, "my-output", "uuid-o")

	if state.Assistants["my-assistant"] != "uuid-a" {
		t.Errorf("SetUUID for assistant failed, got %s", state.Assistants["my-assistant"])
	}
	if state.Tools["my-tool"] != "uuid-t" {
		t.Errorf("SetUUID for tool failed, got %s", state.Tools["my-tool"])
	}
	if state.StructuredOutputs["my-output"] != "uuid-o" {
		t.Errorf("SetUUID for structured output failed, got %s", state.StructuredOutputs["my-output"])
	}
}

func TestStateFile_DeleteResource(t *testing.T) {
	state := NewStateFile()
	state.Assistants["to-delete"] = "uuid-1"
	state.Tools["to-delete"] = "uuid-2"
	state.StructuredOutputs["to-delete"] = "uuid-3"

	state.DeleteResource(ResourceTypeAssistants, "to-delete")
	state.DeleteResource(ResourceTypeTools, "to-delete")
	state.DeleteResource(ResourceTypeStructuredOutputs, "to-delete")

	if _, exists := state.Assistants["to-delete"]; exists {
		t.Error("DeleteResource for assistant failed")
	}
	if _, exists := state.Tools["to-delete"]; exists {
		t.Error("DeleteResource for tool failed")
	}
	if _, exists := state.StructuredOutputs["to-delete"]; exists {
		t.Error("DeleteResource for structured output failed")
	}
}

func TestAllResourceTypes(t *testing.T) {
	types := AllResourceTypes()

	if len(types) != 3 {
		t.Errorf("AllResourceTypes() returned %d types, expected 3", len(types))
	}

	expected := map[ResourceType]bool{
		ResourceTypeAssistants:        true,
		ResourceTypeTools:             true,
		ResourceTypeStructuredOutputs: true,
	}

	for _, rt := range types {
		if !expected[rt] {
			t.Errorf("Unexpected resource type: %s", rt)
		}
	}
}

func TestApplyOrder(t *testing.T) {
	order := ApplyOrder()

	// Tools should come before assistants (dependency order)
	if len(order) != 3 {
		t.Errorf("ApplyOrder() returned %d items, expected 3", len(order))
	}

	// First should be tools (no dependencies)
	if order[0] != ResourceTypeTools {
		t.Errorf("First in ApplyOrder should be tools, got %s", order[0])
	}

	// Last should be assistants (depends on tools and outputs)
	if order[2] != ResourceTypeAssistants {
		t.Errorf("Last in ApplyOrder should be assistants, got %s", order[2])
	}
}

func TestDeleteOrder(t *testing.T) {
	order := DeleteOrder()

	// Assistants should be deleted first (reverse of apply)
	if len(order) != 3 {
		t.Errorf("DeleteOrder() returned %d items, expected 3", len(order))
	}

	// First should be assistants (to remove references)
	if order[0] != ResourceTypeAssistants {
		t.Errorf("First in DeleteOrder should be assistants, got %s", order[0])
	}

	// Last should be tools
	if order[2] != ResourceTypeTools {
		t.Errorf("Last in DeleteOrder should be tools, got %s", order[2])
	}
}

func TestLoadedResources_GetByType(t *testing.T) {
	loaded := &LoadedResources{
		Assistants: []*ResourceFile{
			{ResourceID: "assistant-1"},
		},
		Tools: []*ResourceFile{
			{ResourceID: "tool-1"},
			{ResourceID: "tool-2"},
		},
		StructuredOutputs: []*ResourceFile{
			{ResourceID: "output-1"},
		},
	}

	tests := []struct {
		rt       ResourceType
		expected int
	}{
		{ResourceTypeAssistants, 1},
		{ResourceTypeTools, 2},
		{ResourceTypeStructuredOutputs, 1},
	}

	for _, tt := range tests {
		result := loaded.GetByType(tt.rt)
		if len(result) != tt.expected {
			t.Errorf("GetByType(%s) returned %d items, expected %d", tt.rt, len(result), tt.expected)
		}
	}
}

func TestPullStats(t *testing.T) {
	stats := &PullStats{
		Created: 5,
		Updated: 3,
	}

	if stats.Created != 5 {
		t.Errorf("Created = %d, expected 5", stats.Created)
	}
	if stats.Updated != 3 {
		t.Errorf("Updated = %d, expected 3", stats.Updated)
	}
}
