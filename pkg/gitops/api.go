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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// APIClient handles HTTP requests to the Vapi API.
type APIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewAPIClient creates a new API client.
func NewAPIClient(baseURL, apiKey string) *APIClient {
	return &APIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Request sends an HTTP request to the Vapi API.
func (c *APIClient) Request(ctx context.Context, method, path string, body interface{}) (map[string]interface{}, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	// Handle empty responses
	if len(respBody) == 0 {
		return nil, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

// Get sends a GET request.
func (c *APIClient) Get(ctx context.Context, path string) (map[string]interface{}, error) {
	return c.Request(ctx, http.MethodGet, path, nil)
}

// Post sends a POST request.
func (c *APIClient) Post(ctx context.Context, path string, body interface{}) (map[string]interface{}, error) {
	return c.Request(ctx, http.MethodPost, path, body)
}

// Patch sends a PATCH request.
func (c *APIClient) Patch(ctx context.Context, path string, body interface{}) (map[string]interface{}, error) {
	return c.Request(ctx, http.MethodPatch, path, body)
}

// Delete sends a DELETE request.
func (c *APIClient) Delete(ctx context.Context, path string) error {
	_, err := c.Request(ctx, http.MethodDelete, path, nil)
	return err
}

// ListResources fetches all resources of a given type.
func (c *APIClient) ListResources(ctx context.Context, rt ResourceType) ([]map[string]interface{}, error) {
	endpoint := GetAPIEndpoint(rt)
	url := c.baseURL + endpoint

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

// CreateResource creates a new resource.
func (c *APIClient) CreateResource(ctx context.Context, rt ResourceType, data map[string]interface{}) (map[string]interface{}, error) {
	endpoint := GetAPIEndpoint(rt)
	return c.Post(ctx, endpoint, data)
}

// UpdateResource updates an existing resource.
func (c *APIClient) UpdateResource(ctx context.Context, rt ResourceType, uuid string, data map[string]interface{}) (map[string]interface{}, error) {
	endpoint := GetAPIEndpoint(rt) + "/" + uuid
	// Remove excluded keys for update
	cleanData := RemoveExcludedKeys(data, rt)
	return c.Patch(ctx, endpoint, cleanData)
}

// DeleteResource deletes a resource.
func (c *APIClient) DeleteResource(ctx context.Context, rt ResourceType, uuid string) error {
	endpoint := GetAPIEndpoint(rt) + "/" + uuid
	return c.Delete(ctx, endpoint)
}

// CleanResourceForPull removes server-managed fields from a resource.
func CleanResourceForPull(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		if !IsExcludedOnPull(k) && v != nil {
			result[k] = v
		}
	}
	return result
}
