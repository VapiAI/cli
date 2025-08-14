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

Authors:

	Dan Goosewin <dan@vapi.ai>
*/
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	vapiclient "github.com/VapiAI/server-sdk-go/client"
	"github.com/VapiAI/server-sdk-go/option"

	"github.com/VapiAI/cli/pkg/config"
)

type VapiClient struct {
	client *vapiclient.Client
	config *config.Config
}

func NewVapiClient(apiKey string) (*VapiClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Load configuration to get environment settings
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Set API key from parameter
	cfg.APIKey = apiKey

	// Create client with environment-specific base URL
	options := []option.RequestOption{
		option.WithToken(apiKey),
	}

	// Add base URL if not production
	if baseURL := cfg.GetAPIBaseURL(); baseURL != "https://api.vapi.ai" {
		options = append(options, option.WithBaseURL(baseURL))
	}

	client := vapiclient.NewClient(options...)

	return &VapiClient{
		client: client,
		config: cfg,
	}, nil
}

func (v *VapiClient) GetClient() *vapiclient.Client {
	return v.client
}

func (v *VapiClient) GetConfig() *config.Config {
	return v.config
}

// DoRawJSON sends a raw JSON request to the Vapi API using the underlying client.
// path should be like "/assistants/<id>". method is e.g. "PATCH".
func (v *VapiClient) DoRawJSON(ctx context.Context, method string, path string, body []byte) (map[string]interface{}, error) {
	baseURL := strings.TrimRight(v.config.GetAPIBaseURL(), "/")
	rel := "/" + strings.TrimLeft(path, "/")
	url := baseURL + rel

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	} else {
		bodyReader = http.NoBody
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+v.config.GetActiveAPIKey())

	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = httpResp.Body.Close() }()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error %d: %s", httpResp.StatusCode, string(respBytes))
	}

	var resp map[string]interface{}
	if len(respBytes) > 0 {
		if err := json.Unmarshal(respBytes, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %w", err)
		}
	}
	return resp, nil
}
