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
	"strings"
)

// ResolveReferences resolves resource ID references to UUIDs in a resource's data.
// This handles toolIds, structuredOutputIds, assistant_ids, and assistantId in destinations.
func ResolveReferences(data map[string]interface{}, state *StateFile) map[string]interface{} {
	result := deepCopyMap(data)

	// Resolve toolIds in model
	if model, ok := result["model"].(map[string]interface{}); ok {
		if toolIds, ok := model["toolIds"].([]interface{}); ok {
			model["toolIds"] = resolveIDList(toolIds, state.Tools)
		}
	}

	// Resolve structuredOutputIds in artifactPlan
	if artifactPlan, ok := result["artifactPlan"].(map[string]interface{}); ok {
		if outputIds, ok := artifactPlan["structuredOutputIds"].([]interface{}); ok {
			artifactPlan["structuredOutputIds"] = resolveIDList(outputIds, state.StructuredOutputs)
		}
	}

	// Resolve assistant_ids in structured outputs
	if assistantIds, ok := result["assistant_ids"].([]interface{}); ok {
		result["assistant_ids"] = resolveIDList(assistantIds, state.Assistants)
	}

	// Resolve assistantId in destinations (handoff tools)
	if destinations, ok := result["destinations"].([]interface{}); ok {
		result["destinations"] = resolveDestinations(destinations, state.Assistants)
	}

	return result
}

// resolveIDList resolves a list of resource IDs to UUIDs.
func resolveIDList(ids []interface{}, stateMap map[string]string) []interface{} {
	resolved := make([]interface{}, len(ids))
	for i, id := range ids {
		if strID, ok := id.(string); ok {
			resolved[i] = resolveID(strID, stateMap)
		} else {
			resolved[i] = id
		}
	}
	return resolved
}

// resolveID resolves a single resource ID to UUID.
// Strips inline comments (## ...) before resolving.
func resolveID(id string, stateMap map[string]string) string {
	// Strip inline comments
	id = stripComment(id)

	// Look up in state map
	if uuid, ok := stateMap[id]; ok {
		return uuid
	}

	// Return original if not found (might be an actual UUID)
	return id
}

// resolveDestinations resolves assistantId in destination objects.
func resolveDestinations(destinations []interface{}, assistantMap map[string]string) []interface{} {
	resolved := make([]interface{}, len(destinations))
	for i, dest := range destinations {
		if destMap, ok := dest.(map[string]interface{}); ok {
			resolvedDest := deepCopyMap(destMap)
			if assistantID, ok := resolvedDest["assistantId"].(string); ok {
				resolvedDest["assistantId"] = resolveID(assistantID, assistantMap)
			}
			resolved[i] = resolvedDest
		} else {
			resolved[i] = dest
		}
	}
	return resolved
}

// stripComment removes inline comments (## ...) from a reference.
func stripComment(ref string) string {
	if idx := strings.Index(ref, "##"); idx != -1 {
		return strings.TrimSpace(ref[:idx])
	}
	return strings.TrimSpace(ref)
}

// ResolveAssistantIDs resolves a list of assistant IDs to UUIDs.
func ResolveAssistantIDs(ids []string, state *StateFile) []string {
	resolved := make([]string, 0, len(ids))
	for _, id := range ids {
		cleanID := stripComment(id)
		if uuid, ok := state.Assistants[cleanID]; ok {
			resolved = append(resolved, uuid)
		}
	}
	return resolved
}

// ExtractReferencedIDs extracts all referenced resource IDs from a resource's data.
// Used for orphan detection.
func ExtractReferencedIDs(data map[string]interface{}) *ReferencedIDs {
	refs := &ReferencedIDs{
		Tools:            []string{},
		StructuredOutputs: []string{},
		Assistants:       []string{},
	}

	// Extract toolIds from model
	if model, ok := data["model"].(map[string]interface{}); ok {
		if toolIds, ok := model["toolIds"].([]interface{}); ok {
			for _, id := range toolIds {
				if strID, ok := id.(string); ok {
					refs.Tools = append(refs.Tools, stripComment(strID))
				}
			}
		}
	}

	// Extract structuredOutputIds from artifactPlan
	if artifactPlan, ok := data["artifactPlan"].(map[string]interface{}); ok {
		if outputIds, ok := artifactPlan["structuredOutputIds"].([]interface{}); ok {
			for _, id := range outputIds {
				if strID, ok := id.(string); ok {
					refs.StructuredOutputs = append(refs.StructuredOutputs, stripComment(strID))
				}
			}
		}
	}

	// Extract assistant_ids
	if assistantIds, ok := data["assistant_ids"].([]interface{}); ok {
		for _, id := range assistantIds {
			if strID, ok := id.(string); ok {
				refs.Assistants = append(refs.Assistants, stripComment(strID))
			}
		}
	}

	// Extract assistantId from destinations
	if destinations, ok := data["destinations"].([]interface{}); ok {
		for _, dest := range destinations {
			if destMap, ok := dest.(map[string]interface{}); ok {
				if assistantID, ok := destMap["assistantId"].(string); ok {
					refs.Assistants = append(refs.Assistants, stripComment(assistantID))
				}
			}
		}
	}

	return refs
}

// ReferencedIDs holds lists of referenced resource IDs.
type ReferencedIDs struct {
	Tools            []string
	StructuredOutputs []string
	Assistants       []string
}

// ReverseResolveReferences resolves UUIDs back to resource IDs (for pull).
func ReverseResolveReferences(data map[string]interface{}, state *StateFile) map[string]interface{} {
	result := deepCopyMap(data)

	// Build reverse maps
	toolsMap := BuildReverseMap(state, ResourceTypeTools)
	assistantsMap := BuildReverseMap(state, ResourceTypeAssistants)
	outputsMap := BuildReverseMap(state, ResourceTypeStructuredOutputs)

	// Reverse resolve toolIds in model
	if model, ok := result["model"].(map[string]interface{}); ok {
		if toolIds, ok := model["toolIds"].([]interface{}); ok {
			model["toolIds"] = reverseResolveList(toolIds, toolsMap)
		}
	}

	// Reverse resolve structuredOutputIds in artifactPlan
	if artifactPlan, ok := result["artifactPlan"].(map[string]interface{}); ok {
		if outputIds, ok := artifactPlan["structuredOutputIds"].([]interface{}); ok {
			artifactPlan["structuredOutputIds"] = reverseResolveList(outputIds, outputsMap)
		}
	}

	// Reverse resolve assistant_ids
	if assistantIds, ok := result["assistant_ids"].([]interface{}); ok {
		result["assistant_ids"] = reverseResolveList(assistantIds, assistantsMap)
	}

	// Reverse resolve assistantId in destinations
	if destinations, ok := result["destinations"].([]interface{}); ok {
		result["destinations"] = reverseResolveDestinations(destinations, assistantsMap)
	}

	return result
}

// reverseResolveList resolves UUIDs back to resource IDs.
func reverseResolveList(ids []interface{}, reverseMap map[string]string) []interface{} {
	resolved := make([]interface{}, len(ids))
	for i, id := range ids {
		if uuid, ok := id.(string); ok {
			if resourceID, ok := reverseMap[uuid]; ok {
				resolved[i] = resourceID
			} else {
				resolved[i] = uuid
			}
		} else {
			resolved[i] = id
		}
	}
	return resolved
}

// reverseResolveDestinations resolves UUIDs in destinations back to resource IDs.
func reverseResolveDestinations(destinations []interface{}, assistantMap map[string]string) []interface{} {
	resolved := make([]interface{}, len(destinations))
	for i, dest := range destinations {
		if destMap, ok := dest.(map[string]interface{}); ok {
			resolvedDest := deepCopyMap(destMap)
			if uuid, ok := resolvedDest["assistantId"].(string); ok {
				if resourceID, ok := assistantMap[uuid]; ok {
					resolvedDest["assistantId"] = resourceID
				}
			}
			resolved[i] = resolvedDest
		} else {
			resolved[i] = dest
		}
	}
	return resolved
}

// deepCopyMap creates a deep copy of a map.
func deepCopyMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		switch val := v.(type) {
		case map[string]interface{}:
			result[k] = deepCopyMap(val)
		case []interface{}:
			result[k] = deepCopySlice(val)
		default:
			result[k] = v
		}
	}
	return result
}

// deepCopySlice creates a deep copy of a slice.
func deepCopySlice(s []interface{}) []interface{} {
	result := make([]interface{}, len(s))
	for i, v := range s {
		switch val := v.(type) {
		case map[string]interface{}:
			result[i] = deepCopyMap(val)
		case []interface{}:
			result[i] = deepCopySlice(val)
		default:
			result[i] = v
		}
	}
	return result
}

// StripUnresolvedAssistantDestinations removes destinations with unresolved assistant IDs.
// Used during initial tool creation when assistants don't exist yet.
func StripUnresolvedAssistantDestinations(resolved, original map[string]interface{}) map[string]interface{} {
	destResolved, okResolved := resolved["destinations"].([]interface{})
	destOriginal, okOriginal := original["destinations"].([]interface{})

	if !okResolved || !okOriginal {
		return resolved
	}

	filtered := make([]interface{}, 0)
	for i, dest := range destResolved {
		if i >= len(destOriginal) {
			continue
		}

		destMap, ok := dest.(map[string]interface{})
		if !ok {
			filtered = append(filtered, dest)
			continue
		}

		origMap, ok := destOriginal[i].(map[string]interface{})
		if !ok {
			filtered = append(filtered, dest)
			continue
		}

		resolvedID, hasResolved := destMap["assistantId"].(string)
		originalID, hasOriginal := origMap["assistantId"].(string)

		// Keep if no assistantId or if it was resolved (different from original)
		if !hasResolved || !hasOriginal {
			filtered = append(filtered, dest)
			continue
		}

		cleanOriginal := stripComment(originalID)
		if resolvedID != cleanOriginal {
			// Was resolved, keep it
			filtered = append(filtered, dest)
		}
		// Otherwise, skip this destination (unresolved)
	}

	result := deepCopyMap(resolved)
	result["destinations"] = filtered
	return result
}

// ContainsReference checks if a reference ID is in a list.
func ContainsReference(refs []string, target string) bool {
	for _, ref := range refs {
		if ref == target {
			return true
		}
	}
	return false
}

// ValidationError represents a reference validation error.
type ValidationError struct {
	ResourceID   string
	ResourceType ResourceType
	Message      string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s/%s: %s", e.ResourceType, e.ResourceID, e.Message)
}
