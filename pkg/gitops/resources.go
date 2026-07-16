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
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadResources loads all YAML files for a resource type.
func LoadResources(config *Config, rt ResourceType) ([]*ResourceFile, error) {
	dir := config.GetResourceDir(rt)

	// If directory doesn't exist, return empty slice
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return []*ResourceFile{}, nil
	}

	var resources []*ResourceFile

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .yml and .yaml files
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yml" && ext != ".yaml" {
			return nil
		}

		resource, err := loadResourceFile(path, dir, rt)
		if err != nil {
			return fmt.Errorf("failed to load %s: %w", path, err)
		}

		resources = append(resources, resource)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load resources from %s: %w", dir, err)
	}

	return resources, nil
}

// loadResourceFile loads a single YAML resource file.
func loadResourceFile(path, baseDir string, rt ResourceType) (*ResourceFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var content map[string]interface{}
	if err := yaml.Unmarshal(data, &content); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Generate resource ID from relative path
	relPath, err := filepath.Rel(baseDir, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path: %w", err)
	}

	// Remove extension and convert to forward slashes for consistency
	resourceID := strings.TrimSuffix(relPath, filepath.Ext(relPath))
	resourceID = filepath.ToSlash(resourceID)

	return &ResourceFile{
		ResourceID:   resourceID,
		ResourceType: rt,
		FilePath:     path,
		Data:         content,
	}, nil
}

// LoadAllResources loads all resources of all types.
func LoadAllResources(config *Config) (*LoadedResources, error) {
	loaded := &LoadedResources{}

	tools, err := LoadResources(config, ResourceTypeTools)
	if err != nil {
		return nil, fmt.Errorf("failed to load tools: %w", err)
	}
	loaded.Tools = tools

	outputs, err := LoadResources(config, ResourceTypeStructuredOutputs)
	if err != nil {
		return nil, fmt.Errorf("failed to load structured outputs: %w", err)
	}
	loaded.StructuredOutputs = outputs

	assistants, err := LoadResources(config, ResourceTypeAssistants)
	if err != nil {
		return nil, fmt.Errorf("failed to load assistants: %w", err)
	}
	loaded.Assistants = assistants

	return loaded, nil
}

// WriteResourceFile writes a resource to a YAML file.
func WriteResourceFile(config *Config, rt ResourceType, resourceID string, data map[string]interface{}) (string, error) {
	dir := config.GetResourceDir(rt)
	filePath := filepath.Join(dir, resourceID+".yml")

	// Ensure directory exists (including nested directories)
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to YAML
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, yamlData, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

// GetResourceIDsFromFiles returns all resource IDs from loaded resources.
func GetResourceIDsFromFiles(resources []*ResourceFile) []string {
	ids := make([]string, len(resources))
	for i, r := range resources {
		ids[i] = r.ResourceID
	}
	return ids
}

// Slugify converts a name to a URL-friendly slug.
func Slugify(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace non-alphanumeric characters with hyphens
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		} else {
			result.WriteRune('-')
		}
	}

	// Clean up multiple hyphens and trim
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")

	return slug
}

// GenerateUniqueResourceID generates a unique resource ID from a name.
func GenerateUniqueResourceID(name string, existingIDs map[string]bool) string {
	base := Slugify(name)
	if base == "" {
		base = "resource"
	}

	resourceID := base
	counter := 1

	for existingIDs[resourceID] {
		resourceID = fmt.Sprintf("%s-%d", base, counter)
		counter++
	}

	return resourceID
}
