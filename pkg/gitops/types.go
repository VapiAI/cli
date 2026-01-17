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

// ResourceType represents the type of Vapi resource.
type ResourceType string

const (
	ResourceTypeAssistants       ResourceType = "assistants"
	ResourceTypeTools            ResourceType = "tools"
	ResourceTypeStructuredOutputs ResourceType = "structuredOutputs"
)

// AllResourceTypes returns all supported resource types in dependency order.
func AllResourceTypes() []ResourceType {
	return []ResourceType{
		ResourceTypeTools,
		ResourceTypeStructuredOutputs,
		ResourceTypeAssistants,
	}
}

// DeleteOrder returns resource types in reverse dependency order for safe deletion.
func DeleteOrder() []ResourceType {
	return []ResourceType{
		ResourceTypeAssistants,
		ResourceTypeStructuredOutputs,
		ResourceTypeTools,
	}
}

// StateFile represents the mapping between resource IDs and Vapi UUIDs.
type StateFile struct {
	Assistants       map[string]string `json:"assistants"`
	Tools            map[string]string `json:"tools"`
	StructuredOutputs map[string]string `json:"structuredOutputs"`
}

// NewStateFile creates an empty state file.
func NewStateFile() *StateFile {
	return &StateFile{
		Assistants:       make(map[string]string),
		Tools:            make(map[string]string),
		StructuredOutputs: make(map[string]string),
	}
}

// GetSection returns the state section for a given resource type.
func (s *StateFile) GetSection(rt ResourceType) map[string]string {
	switch rt {
	case ResourceTypeAssistants:
		return s.Assistants
	case ResourceTypeTools:
		return s.Tools
	case ResourceTypeStructuredOutputs:
		return s.StructuredOutputs
	default:
		return nil
	}
}

// SetUUID sets the UUID for a resource ID in the appropriate section.
func (s *StateFile) SetUUID(rt ResourceType, resourceID, uuid string) {
	switch rt {
	case ResourceTypeAssistants:
		s.Assistants[resourceID] = uuid
	case ResourceTypeTools:
		s.Tools[resourceID] = uuid
	case ResourceTypeStructuredOutputs:
		s.StructuredOutputs[resourceID] = uuid
	}
}

// DeleteResource removes a resource from the state.
func (s *StateFile) DeleteResource(rt ResourceType, resourceID string) {
	switch rt {
	case ResourceTypeAssistants:
		delete(s.Assistants, resourceID)
	case ResourceTypeTools:
		delete(s.Tools, resourceID)
	case ResourceTypeStructuredOutputs:
		delete(s.StructuredOutputs, resourceID)
	}
}

// ResourceFile represents a loaded resource from a YAML file.
type ResourceFile struct {
	ResourceID   string                 // Filename without extension (e.g., "transfer-call" or "company-1/transfer-call")
	ResourceType ResourceType           // The type of resource
	FilePath     string                 // Full path to the YAML file
	Data         map[string]interface{} // Parsed YAML content
}

// LoadedResources holds all loaded resources by type.
type LoadedResources struct {
	Tools            []*ResourceFile
	StructuredOutputs []*ResourceFile
	Assistants       []*ResourceFile
}

// GetByType returns the resources for a given type.
func (lr *LoadedResources) GetByType(rt ResourceType) []*ResourceFile {
	switch rt {
	case ResourceTypeTools:
		return lr.Tools
	case ResourceTypeStructuredOutputs:
		return lr.StructuredOutputs
	case ResourceTypeAssistants:
		return lr.Assistants
	default:
		return nil
	}
}

// OrphanedResource represents a resource in state but not in filesystem.
type OrphanedResource struct {
	ResourceID string
	UUID       string
}

// VapiResource represents a resource fetched from the Vapi API.
type VapiResource struct {
	ID   string                 `json:"id"`
	Name string                 `json:"name,omitempty"`
	Data map[string]interface{} `json:"-"` // Raw JSON data
}

// ApplyResult holds the result of applying a resource.
type ApplyResult struct {
	ResourceID string
	UUID       string
	Created    bool
	Updated    bool
	Error      error
}

// PullStats tracks pull operation statistics.
type PullStats struct {
	Created int
	Updated int
}
