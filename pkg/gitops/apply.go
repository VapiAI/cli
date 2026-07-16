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

// ApplyEngine handles the apply workflow.
type ApplyEngine struct {
	config *Config
	client *APIClient
	state  *StateFile
}

// NewApplyEngine creates a new apply engine.
func NewApplyEngine(config *Config) (*ApplyEngine, error) {
	if err := config.ValidateForApply(); err != nil {
		return nil, err
	}

	state, err := LoadState(config.StateFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	client := NewAPIClient(config.APIBaseURL, config.APIKey)

	return &ApplyEngine{
		config: config,
		client: client,
		state:  state,
	}, nil
}

// Apply runs the full apply workflow.
func (e *ApplyEngine) Apply(ctx context.Context) error {
	// Load all resources
	loaded, err := LoadAllResources(e.config)
	if err != nil {
		return fmt.Errorf("failed to load resources: %w", err)
	}

	// Delete orphaned resources first
	if err := e.deleteOrphaned(ctx, loaded); err != nil {
		return fmt.Errorf("failed to delete orphaned resources: %w", err)
	}

	// Apply in dependency order: tools → structured outputs → assistants
	fmt.Println("\n🔧 Applying tools...")
	for _, tool := range loaded.Tools {
		if err := e.applyResource(ctx, tool); err != nil {
			return fmt.Errorf("failed to apply tool %s: %w", tool.ResourceID, err)
		}
	}

	fmt.Println("\n📊 Applying structured outputs...")
	for _, output := range loaded.StructuredOutputs {
		if err := e.applyResource(ctx, output); err != nil {
			return fmt.Errorf("failed to apply structured output %s: %w", output.ResourceID, err)
		}
	}

	fmt.Println("\n🤖 Applying assistants...")
	for _, assistant := range loaded.Assistants {
		if err := e.applyResource(ctx, assistant); err != nil {
			return fmt.Errorf("failed to apply assistant %s: %w", assistant.ResourceID, err)
		}
	}

	// Second pass: Link resources to assistants (now that assistants exist)
	fmt.Println("\n🔗 Linking tools to assistant destinations...")
	if err := e.updateToolAssistantRefs(ctx, loaded.Tools); err != nil {
		return fmt.Errorf("failed to update tool assistant refs: %w", err)
	}

	fmt.Println("\n🔗 Linking structured outputs to assistants...")
	if err := e.updateStructuredOutputAssistantRefs(ctx, loaded.StructuredOutputs); err != nil {
		return fmt.Errorf("failed to update structured output assistant refs: %w", err)
	}

	// Save state
	if err := SaveState(e.config.StateFilePath, e.state); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}

// applyResource applies a single resource.
func (e *ApplyEngine) applyResource(ctx context.Context, resource *ResourceFile) error {
	section := e.state.GetSection(resource.ResourceType)
	existingUUID := section[resource.ResourceID]

	// Resolve references
	payload := ResolveReferences(resource.Data, e.state)

	if existingUUID != "" {
		// Update existing resource
		fmt.Printf("  🔄 Updating %s: %s (%s)\n", resource.ResourceType, resource.ResourceID, existingUUID)
		_, err := e.client.UpdateResource(ctx, resource.ResourceType, existingUUID, payload)
		if err != nil {
			return err
		}
	} else {
		// Create new resource
		// For tools with assistant destinations, strip unresolved assistantIds
		createPayload := payload
		if resource.ResourceType == ResourceTypeTools {
			createPayload = StripUnresolvedAssistantDestinations(payload, resource.Data)
		}

		// For structured outputs, remove assistant_ids for initial creation
		if resource.ResourceType == ResourceTypeStructuredOutputs {
			delete(createPayload, "assistant_ids")
		}

		fmt.Printf("  ✨ Creating %s: %s\n", resource.ResourceType, resource.ResourceID)
		result, err := e.client.CreateResource(ctx, resource.ResourceType, createPayload)
		if err != nil {
			return err
		}

		uuid, ok := result["id"].(string)
		if !ok {
			return fmt.Errorf("no ID returned from create")
		}

		e.state.SetUUID(resource.ResourceType, resource.ResourceID, uuid)
	}

	return nil
}

// deleteOrphaned deletes resources that are in state but not in filesystem.
func (e *ApplyEngine) deleteOrphaned(ctx context.Context, loaded *LoadedResources) error {
	fmt.Println("\n🗑️  Checking for deleted resources...")

	// Find orphaned resources
	orphaned := make(map[ResourceType][]OrphanedResource)

	for _, rt := range DeleteOrder() {
		fileIDs := GetResourceIDsFromFiles(loaded.GetByType(rt))
		stateSection := e.state.GetSection(rt)

		for resourceID, uuid := range stateSection {
			found := false
			for _, fileID := range fileIDs {
				if fileID == resourceID {
					found = true
					break
				}
			}
			if !found {
				orphaned[rt] = append(orphaned[rt], OrphanedResource{
					ResourceID: resourceID,
					UUID:       uuid,
				})
			}
		}
	}

	// Check for orphan references before deleting
	var errors []string
	for rt, orphans := range orphaned {
		for _, orphan := range orphans {
			refs := findReferencingResources(orphan.ResourceID, rt, loaded)
			if len(refs) > 0 {
				errors = append(errors, fmt.Sprintf(
					"Cannot delete %s \"%s\" - still referenced by: %v",
					rt, orphan.ResourceID, refs,
				))
			}
		}
	}

	if len(errors) > 0 {
		fmt.Println("\n❌ Orphan reference errors:")
		for _, err := range errors {
			fmt.Printf("   %s\n", err)
		}
		return fmt.Errorf("cannot delete resources that are still referenced")
	}

	// Delete orphaned resources in reverse dependency order
	for _, rt := range DeleteOrder() {
		for _, orphan := range orphaned[rt] {
			fmt.Printf("  🗑️  Deleting %s: %s (%s)\n", rt, orphan.ResourceID, orphan.UUID)
			if err := e.client.DeleteResource(ctx, rt, orphan.UUID); err != nil {
				return fmt.Errorf("failed to delete %s %s: %w", rt, orphan.ResourceID, err)
			}
			e.state.DeleteResource(rt, orphan.ResourceID)
		}
	}

	return nil
}

// findReferencingResources finds resources that reference a given resource ID.
func findReferencingResources(targetID string, targetType ResourceType, loaded *LoadedResources) []string {
	var refs []string

	checkResource := func(resource *ResourceFile) {
		extracted := ExtractReferencedIDs(resource.Data)

		switch targetType {
		case ResourceTypeTools:
			if ContainsReference(extracted.Tools, targetID) {
				refs = append(refs, string(resource.ResourceType)+"/"+resource.ResourceID)
			}
		case ResourceTypeStructuredOutputs:
			if ContainsReference(extracted.StructuredOutputs, targetID) {
				refs = append(refs, string(resource.ResourceType)+"/"+resource.ResourceID)
			}
		case ResourceTypeAssistants:
			if ContainsReference(extracted.Assistants, targetID) {
				refs = append(refs, string(resource.ResourceType)+"/"+resource.ResourceID)
			}
		}
	}

	for _, r := range loaded.Assistants {
		checkResource(r)
	}
	for _, r := range loaded.StructuredOutputs {
		checkResource(r)
	}
	for _, r := range loaded.Tools {
		checkResource(r)
	}

	return refs
}

// updateToolAssistantRefs updates tools with assistant destination references.
func (e *ApplyEngine) updateToolAssistantRefs(ctx context.Context, tools []*ResourceFile) error {
	for _, tool := range tools {
		destinations, ok := tool.Data["destinations"].([]interface{})
		if !ok {
			continue
		}

		hasAssistantRefs := false
		for _, dest := range destinations {
			if destMap, ok := dest.(map[string]interface{}); ok {
				if _, hasID := destMap["assistantId"].(string); hasID {
					hasAssistantRefs = true
					break
				}
			}
		}

		if !hasAssistantRefs {
			continue
		}

		uuid := e.state.Tools[tool.ResourceID]
		if uuid == "" {
			continue
		}

		// Resolve destinations now that all assistants exist
		resolved := ResolveReferences(tool.Data, e.state)

		fmt.Printf("  🔗 Linking tool %s to assistant destinations\n", tool.ResourceID)
		_, err := e.client.Patch(ctx, GetAPIEndpoint(ResourceTypeTools)+"/"+uuid, map[string]interface{}{
			"destinations": resolved["destinations"],
		})
		if err != nil {
			return fmt.Errorf("failed to update tool %s: %w", tool.ResourceID, err)
		}
	}

	return nil
}

// updateStructuredOutputAssistantRefs updates structured outputs with assistant references.
func (e *ApplyEngine) updateStructuredOutputAssistantRefs(ctx context.Context, outputs []*ResourceFile) error {
	for _, output := range outputs {
		assistantIDs, ok := output.Data["assistant_ids"].([]interface{})
		if !ok || len(assistantIDs) == 0 {
			continue
		}

		uuid := e.state.StructuredOutputs[output.ResourceID]
		if uuid == "" {
			continue
		}

		// Convert to string slice
		ids := make([]string, 0, len(assistantIDs))
		for _, id := range assistantIDs {
			if strID, ok := id.(string); ok {
				ids = append(ids, strID)
			}
		}

		// Resolve assistant IDs
		resolvedIDs := ResolveAssistantIDs(ids, e.state)
		if len(resolvedIDs) == 0 {
			continue
		}

		fmt.Printf("  🔗 Linking structured output %s to assistants\n", output.ResourceID)
		_, err := e.client.Patch(ctx, GetAPIEndpoint(ResourceTypeStructuredOutputs)+"/"+uuid, map[string]interface{}{
			"assistantIds": resolvedIDs,
		})
		if err != nil {
			return fmt.Errorf("failed to update structured output %s: %w", output.ResourceID, err)
		}
	}

	return nil
}

// GetState returns the current state.
func (e *ApplyEngine) GetState() *StateFile {
	return e.state
}

// Summary returns a summary of the current state.
func (e *ApplyEngine) Summary() string {
	return fmt.Sprintf(
		"Tools: %d, Structured Outputs: %d, Assistants: %d",
		len(e.state.Tools),
		len(e.state.StructuredOutputs),
		len(e.state.Assistants),
	)
}
