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

func TestStripInlineComment(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"my-tool ## This is a comment", "my-tool"},
		{"my-tool##comment", "my-tool"},
		{"my-tool", "my-tool"},
		{"  my-tool  ## comment  ", "my-tool"},
		{"", ""},
	}

	for _, tt := range tests {
		result := StripInlineComment(tt.input)
		if result != tt.expected {
			t.Errorf("StripInlineComment(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestBuildReverseMap(t *testing.T) {
	state := NewStateFile()
	state.Tools["my-tool"] = "uuid-tool-1"
	state.Tools["other-tool"] = "uuid-tool-2"
	state.Assistants["my-assistant"] = "uuid-assistant-1"

	// Test tools reverse map
	toolsMap := BuildReverseMap(state, ResourceTypeTools)
	if toolsMap["uuid-tool-1"] != "my-tool" {
		t.Errorf("Expected uuid-tool-1 -> my-tool, got %s", toolsMap["uuid-tool-1"])
	}
	if toolsMap["uuid-tool-2"] != "other-tool" {
		t.Errorf("Expected uuid-tool-2 -> other-tool, got %s", toolsMap["uuid-tool-2"])
	}

	// Test assistants reverse map
	assistantsMap := BuildReverseMap(state, ResourceTypeAssistants)
	if assistantsMap["uuid-assistant-1"] != "my-assistant" {
		t.Errorf("Expected uuid-assistant-1 -> my-assistant, got %s", assistantsMap["uuid-assistant-1"])
	}
}

func TestResolveToolIDs(t *testing.T) {
	state := NewStateFile()
	state.Tools["weather-tool"] = "uuid-weather"
	state.Tools["calendar-tool"] = "uuid-calendar"

	toolIDs := []string{
		"weather-tool ## Get weather",
		"calendar-tool",
		"unknown-tool",
	}

	resolved := ResolveToolIDs(toolIDs, state)

	if len(resolved) != 3 {
		t.Errorf("Expected 3 resolved IDs, got %d", len(resolved))
	}

	if resolved[0] != "uuid-weather" {
		t.Errorf("Expected uuid-weather, got %s", resolved[0])
	}
	if resolved[1] != "uuid-calendar" {
		t.Errorf("Expected uuid-calendar, got %s", resolved[1])
	}
	// Unknown tools should be kept as-is (might be a real UUID)
	if resolved[2] != "unknown-tool" {
		t.Errorf("Expected unknown-tool (unchanged), got %s", resolved[2])
	}
}

func TestResolveAssistantIDs(t *testing.T) {
	state := NewStateFile()
	state.Assistants["support-bot"] = "uuid-support"
	state.Assistants["sales-bot"] = "uuid-sales"

	assistantIDs := []string{
		"support-bot",
		"sales-bot ## Sales assistant",
	}

	resolved := ResolveAssistantIDs(assistantIDs, state)

	if len(resolved) != 2 {
		t.Errorf("Expected 2 resolved IDs, got %d", len(resolved))
	}

	if resolved[0] != "uuid-support" {
		t.Errorf("Expected uuid-support, got %s", resolved[0])
	}
	if resolved[1] != "uuid-sales" {
		t.Errorf("Expected uuid-sales, got %s", resolved[1])
	}
}

func TestResolveReferences(t *testing.T) {
	state := NewStateFile()
	state.Tools["my-tool"] = "uuid-tool-123"
	state.StructuredOutputs["my-output"] = "uuid-output-456"

	data := map[string]interface{}{
		"name": "Test Assistant",
		"model": map[string]interface{}{
			"provider": "openai",
			"model":    "gpt-4o",
			"toolIds": []interface{}{
				"my-tool ## Weather tool",
			},
		},
		"structuredOutputIds": []interface{}{
			"my-output",
		},
	}

	resolved := ResolveReferences(data, state)

	// Check that original data is not modified
	originalToolIds := data["model"].(map[string]interface{})["toolIds"].([]interface{})
	if originalToolIds[0] != "my-tool ## Weather tool" {
		t.Error("Original data should not be modified")
	}

	// Check resolved values
	model := resolved["model"].(map[string]interface{})
	toolIds := model["toolIds"].([]interface{})
	if toolIds[0] != "uuid-tool-123" {
		t.Errorf("Expected uuid-tool-123, got %v", toolIds[0])
	}

	outputIds := resolved["structuredOutputIds"].([]interface{})
	if outputIds[0] != "uuid-output-456" {
		t.Errorf("Expected uuid-output-456, got %v", outputIds[0])
	}
}

func TestReverseResolveReferences(t *testing.T) {
	state := NewStateFile()
	state.Tools["my-tool"] = "uuid-tool-123"
	state.Assistants["my-assistant"] = "uuid-assistant-789"

	data := map[string]interface{}{
		"name": "Test Tool",
		"destinations": []interface{}{
			map[string]interface{}{
				"type":        "assistant",
				"assistantId": "uuid-assistant-789",
			},
		},
	}

	resolved := ReverseResolveReferences(data, state)

	destinations := resolved["destinations"].([]interface{})
	dest := destinations[0].(map[string]interface{})

	if dest["assistantId"] != "my-assistant" {
		t.Errorf("Expected my-assistant, got %v", dest["assistantId"])
	}
}

func TestExtractReferencedIDs(t *testing.T) {
	data := map[string]interface{}{
		"model": map[string]interface{}{
			"toolIds": []interface{}{
				"tool-1",
				"tool-2 ## comment",
			},
		},
		"structuredOutputIds": []interface{}{
			"output-1",
		},
		"destinations": []interface{}{
			map[string]interface{}{
				"assistantId": "assistant-1",
			},
		},
	}

	refs := ExtractReferencedIDs(data)

	if len(refs.Tools) != 2 {
		t.Errorf("Expected 2 tool refs, got %d", len(refs.Tools))
	}
	if refs.Tools[0] != "tool-1" || refs.Tools[1] != "tool-2" {
		t.Errorf("Unexpected tool refs: %v", refs.Tools)
	}

	if len(refs.StructuredOutputs) != 1 {
		t.Errorf("Expected 1 output ref, got %d", len(refs.StructuredOutputs))
	}
	if refs.StructuredOutputs[0] != "output-1" {
		t.Errorf("Unexpected output ref: %s", refs.StructuredOutputs[0])
	}

	if len(refs.Assistants) != 1 {
		t.Errorf("Expected 1 assistant ref, got %d", len(refs.Assistants))
	}
	if refs.Assistants[0] != "assistant-1" {
		t.Errorf("Unexpected assistant ref: %s", refs.Assistants[0])
	}
}

func TestContainsReference(t *testing.T) {
	refs := []string{"tool-1", "tool-2", "tool-3"}

	if !ContainsReference(refs, "tool-2") {
		t.Error("ContainsReference should return true for existing ref")
	}

	if ContainsReference(refs, "tool-4") {
		t.Error("ContainsReference should return false for non-existing ref")
	}
}

func TestStripUnresolvedAssistantDestinations(t *testing.T) {
	state := NewStateFile()
	state.Assistants["existing-assistant"] = "uuid-exists"

	original := map[string]interface{}{
		"destinations": []interface{}{
			map[string]interface{}{
				"type":        "assistant",
				"assistantId": "new-assistant ## Comment",
			},
			map[string]interface{}{
				"type":        "number",
				"phoneNumber": "+1234567890",
			},
		},
	}

	resolved := map[string]interface{}{
		"destinations": []interface{}{
			map[string]interface{}{
				"type":        "assistant",
				"assistantId": "new-assistant", // Not resolved, assistant doesn't exist
			},
			map[string]interface{}{
				"type":        "number",
				"phoneNumber": "+1234567890",
			},
		},
	}

	result := StripUnresolvedAssistantDestinations(resolved, original)

	destinations := result["destinations"].([]interface{})
	if len(destinations) != 1 {
		t.Errorf("Expected 1 destination after stripping, got %d", len(destinations))
	}

	// Only the phone number destination should remain
	dest := destinations[0].(map[string]interface{})
	if dest["type"] != "number" {
		t.Errorf("Expected number destination to remain, got %v", dest["type"])
	}
}

func TestIsUUID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"550E8400-E29B-41D4-A716-446655440000", true},
		{"my-resource-id", false},
		{"not-a-uuid", false},
		{"", false},
		{"550e8400-e29b-41d4-a716", false}, // Too short
	}

	for _, tt := range tests {
		result := IsUUID(tt.input)
		if result != tt.expected {
			t.Errorf("IsUUID(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}
