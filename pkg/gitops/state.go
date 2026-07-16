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
	"encoding/json"
	"fmt"
	"os"
)

// LoadState loads the state file from disk.
// If the file doesn't exist, returns an empty state.
func LoadState(path string) (*StateFile, error) {
	state := NewStateFile()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return state, nil
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	if err := json.Unmarshal(data, state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	// Ensure all maps are initialized
	if state.Assistants == nil {
		state.Assistants = make(map[string]string)
	}
	if state.Tools == nil {
		state.Tools = make(map[string]string)
	}
	if state.StructuredOutputs == nil {
		state.StructuredOutputs = make(map[string]string)
	}

	return state, nil
}

// SaveState saves the state file to disk.
func SaveState(path string, state *StateFile) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// BuildReverseMap creates a UUID -> resourceID map for a resource type.
func BuildReverseMap(state *StateFile, rt ResourceType) map[string]string {
	section := state.GetSection(rt)
	result := make(map[string]string)
	for resourceID, uuid := range section {
		result[uuid] = resourceID
	}
	return result
}

// GetResourceIDs returns all resource IDs for a given type.
func GetResourceIDs(state *StateFile, rt ResourceType) []string {
	section := state.GetSection(rt)
	ids := make([]string, 0, len(section))
	for id := range section {
		ids = append(ids, id)
	}
	return ids
}
