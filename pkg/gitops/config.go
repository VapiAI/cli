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
	"os"
	"path/filepath"
)

const (
	// DefaultResourcesDir is the default directory for resource YAML files.
	DefaultResourcesDir = "resources"

	// DefaultStateFile is the default state file name.
	DefaultStateFile = ".vapi-state.json"

	// DefaultAPIBaseURL is the default Vapi API base URL.
	DefaultAPIBaseURL = "https://api.vapi.ai"
)

// ResourceDirs maps resource types to their directory names.
var ResourceDirs = map[ResourceType]string{
	ResourceTypeAssistants:       "assistants",
	ResourceTypeTools:            "tools",
	ResourceTypeStructuredOutputs: "structuredOutputs",
}

// APIEndpoints maps resource types to their API endpoints.
var APIEndpoints = map[ResourceType]string{
	ResourceTypeAssistants:       "/assistant",
	ResourceTypeTools:            "/tool",
	ResourceTypeStructuredOutputs: "/structured-output",
}

// UpdateExcludedKeys lists fields that cannot be updated after creation.
var UpdateExcludedKeys = map[ResourceType][]string{
	ResourceTypeTools:            {"type"},
	ResourceTypeAssistants:       {},
	ResourceTypeStructuredOutputs: {"type"},
}

// ExcludedFieldsOnPull lists fields to remove when pulling resources.
var ExcludedFieldsOnPull = []string{
	"id",
	"orgId",
	"createdAt",
	"updatedAt",
	"analyticsMetadata",
	"isDeleted",
}

// Config holds the GitOps configuration.
type Config struct {
	// ProjectPath is the root path of the GitOps project.
	ProjectPath string

	// ResourcesDir is the directory containing resource YAML files.
	ResourcesDir string

	// StateFilePath is the path to the state file.
	StateFilePath string

	// APIBaseURL is the Vapi API base URL.
	APIBaseURL string

	// APIKey is the Vapi API key.
	APIKey string
}

// NewConfig creates a new Config with default values.
func NewConfig(projectPath string) *Config {
	return &Config{
		ProjectPath:   projectPath,
		ResourcesDir:  filepath.Join(projectPath, DefaultResourcesDir),
		StateFilePath: filepath.Join(projectPath, DefaultStateFile),
		APIBaseURL:    DefaultAPIBaseURL,
	}
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	// Check if project path exists
	if _, err := os.Stat(c.ProjectPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", c.ProjectPath)
	}

	return nil
}

// ValidateForApply checks if the configuration is valid for apply operations.
func (c *Config) ValidateForApply() error {
	if err := c.Validate(); err != nil {
		return err
	}

	if c.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	// Check if resources directory exists
	if _, err := os.Stat(c.ResourcesDir); os.IsNotExist(err) {
		return fmt.Errorf("resources directory does not exist: %s (run 'vapi gitops init' first)", c.ResourcesDir)
	}

	return nil
}

// GetResourceDir returns the full path to a resource type's directory.
func (c *Config) GetResourceDir(rt ResourceType) string {
	return filepath.Join(c.ResourcesDir, ResourceDirs[rt])
}

// GetAPIEndpoint returns the API endpoint for a resource type.
func GetAPIEndpoint(rt ResourceType) string {
	return APIEndpoints[rt]
}

// IsExcludedOnUpdate checks if a field should be excluded on update.
func IsExcludedOnUpdate(rt ResourceType, field string) bool {
	excluded := UpdateExcludedKeys[rt]
	for _, f := range excluded {
		if f == field {
			return true
		}
	}
	return false
}

// IsExcludedOnPull checks if a field should be excluded when pulling.
func IsExcludedOnPull(field string) bool {
	for _, f := range ExcludedFieldsOnPull {
		if f == field {
			return true
		}
	}
	return false
}

// RemoveExcludedKeys removes fields that cannot be updated.
func RemoveExcludedKeys(data map[string]interface{}, rt ResourceType) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		if !IsExcludedOnUpdate(rt, k) {
			result[k] = v
		}
	}
	return result
}
