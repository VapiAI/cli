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
	"context"
	"fmt"
)

// PullEngine handles the pull workflow.
type PullEngine struct {
	config *Config
	client *APIClient
	state  *StateFile
}

// NewPullEngine creates a new pull engine.
func NewPullEngine(config *Config) (*PullEngine, error) {
	if err := config.ValidateForApply(); err != nil {
		return nil, err
	}

	state, err := LoadState(config.StateFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	client := NewAPIClient(config.APIBaseURL, config.APIKey)

	return &PullEngine{
		config: config,
		client: client,
		state:  state,
	}, nil
}

// Pull runs the full pull workflow.
func (e *PullEngine) Pull(ctx context.Context) (map[ResourceType]*PullStats, error) {
	stats := make(map[ResourceType]*PullStats)

	// Pull in dependency order (tools first, then structured outputs, then assistants)
	// This ensures references can be resolved correctly

	fmt.Println("\n📥 Pulling tools...")
	toolStats, err := e.pullResourceType(ctx, ResourceTypeTools)
	if err != nil {
		return nil, fmt.Errorf("failed to pull tools: %w", err)
	}
	stats[ResourceTypeTools] = toolStats

	fmt.Println("\n📥 Pulling structured outputs...")
	outputStats, err := e.pullResourceType(ctx, ResourceTypeStructuredOutputs)
	if err != nil {
		return nil, fmt.Errorf("failed to pull structured outputs: %w", err)
	}
	stats[ResourceTypeStructuredOutputs] = outputStats

	fmt.Println("\n📥 Pulling assistants...")
	assistantStats, err := e.pullResourceType(ctx, ResourceTypeAssistants)
	if err != nil {
		return nil, fmt.Errorf("failed to pull assistants: %w", err)
	}
	stats[ResourceTypeAssistants] = assistantStats

	// Save state
	if err := SaveState(e.config.StateFilePath, e.state); err != nil {
		return nil, fmt.Errorf("failed to save state: %w", err)
	}

	return stats, nil
}

// pullResourceType pulls all resources of a given type.
func (e *PullEngine) pullResourceType(ctx context.Context, rt ResourceType) (*PullStats, error) {
	stats := &PullStats{}

	// Fetch all resources from API
	resources, err := e.client.ListResources(ctx, rt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resources: %w", err)
	}

	fmt.Printf("   Found %d %s in Vapi\n", len(resources), rt)

	// Build reverse map for existing state
	reverseMap := BuildReverseMap(e.state, rt)
	existingIDs := make(map[string]bool)
	for id := range e.state.GetSection(rt) {
		existingIDs[id] = true
	}

	// Track new state section
	newSection := make(map[string]string)

	for _, resource := range resources {
		uuid, ok := resource["id"].(string)
		if !ok {
			continue
		}

		// Check if we already have this resource in state (by UUID)
		resourceID := reverseMap[uuid]
		isNew := resourceID == ""

		if isNew {
			// Generate new resource ID from name or UUID
			name, _ := resource["name"].(string)
			if name == "" {
				name = uuid[:8]
			}
			resourceID = GenerateUniqueResourceID(name, existingIDs)
			existingIDs[resourceID] = true
			stats.Created++
		} else {
			stats.Updated++
		}

		// Clean resource (remove server-managed fields)
		cleaned := CleanResourceForPull(resource)

		// Reverse resolve references (UUIDs → resource IDs)
		resolved := ReverseResolveReferences(cleaned, e.state)

		// Write to file
		filePath, err := WriteResourceFile(e.config, rt, resourceID, resolved)
		if err != nil {
			return nil, fmt.Errorf("failed to write %s: %w", resourceID, err)
		}

		icon := "📝"
		if isNew {
			icon = "✨"
		}
		fmt.Printf("   %s %s -> %s\n", icon, resourceID, filePath)

		// Update state
		newSection[resourceID] = uuid
	}

	// Update state with new mappings
	switch rt {
	case ResourceTypeTools:
		e.state.Tools = newSection
	case ResourceTypeStructuredOutputs:
		e.state.StructuredOutputs = newSection
	case ResourceTypeAssistants:
		e.state.Assistants = newSection
	}

	return stats, nil
}

// GetState returns the current state.
func (e *PullEngine) GetState() *StateFile {
	return e.state
}
