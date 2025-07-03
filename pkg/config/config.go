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
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey           string             `mapstructure:"api_key"` // Legacy single API key (for backward compatibility)
	BaseURL          string             `mapstructure:"base_url"`
	DashboardURL     string             `mapstructure:"dashboard_url"`
	Environment      string             `mapstructure:"environment"`
	Timeout          int                `mapstructure:"timeout"`
	DisableAnalytics bool               `mapstructure:"disable_analytics"`
	Accounts         map[string]Account `mapstructure:"accounts"`       // Multiple accounts support
	ActiveAccount    string             `mapstructure:"active_account"` // Which account is currently active
}

// Account represents a single authenticated account/organization
type Account struct {
	APIKey       string `mapstructure:"apikey"`
	Organization string `mapstructure:"organization,omitempty"` // Organization name if available
	Environment  string `mapstructure:"environment,omitempty"`  // Per-account environment override
	LoginTime    string `mapstructure:"logintime,omitempty"`    // When this account was last authenticated
}

// Environment configuration
type Environment struct {
	Name         string
	APIBaseURL   string
	DashboardURL string
}

var environments = map[string]Environment{
	"production": {
		Name:         "production",
		APIBaseURL:   "https://api.vapi.ai",
		DashboardURL: "https://dashboard.vapi.ai",
	},
	"staging": {
		Name:         "staging",
		APIBaseURL:   "https://api.staging.vapi.ai",
		DashboardURL: "https://dashboard.staging.vapi.ai",
	},
	"development": {
		Name:         "development",
		APIBaseURL:   "http://localhost:3000",
		DashboardURL: "http://localhost:3001",
	},
}

func LoadConfig() (*Config, error) {
	// Set config name and paths
	viper.SetConfigName(".vapi-cli")
	viper.SetConfigType("yaml")

	// Add config search paths
	viper.AddConfigPath(".")
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home)
	}

	// Set environment variable prefix
	viper.SetEnvPrefix("VAPI")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("timeout", 30)
	viper.SetDefault("environment", "production")
	viper.SetDefault("disable_analytics", false)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found is okay, we'll create one if needed
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Apply environment-specific URLs if not explicitly set
	if err := config.applyEnvironment(); err != nil {
		return nil, err
	}

	return &config, nil
}

// applyEnvironment sets URLs based on the environment configuration
func (c *Config) applyEnvironment() error {
	// Determine environment from multiple sources (priority order):
	// 1. Environment variable VAPI_ENV
	// 2. Environment field in config
	// 3. Default to production

	env := c.Environment
	if envVar := os.Getenv("VAPI_ENV"); envVar != "" {
		env = envVar
	}
	if env == "" {
		env = "production"
	}

	// Normalize environment name
	env = strings.ToLower(env)
	if env == "dev" || env == "local" {
		env = "development"
	}
	if env == "stage" {
		env = "staging"
	}
	if env == "prod" {
		env = "production"
	}

	// Get environment configuration
	envConfig, exists := environments[env]
	if !exists {
		return fmt.Errorf("unknown environment: %s (valid: production, staging, development)", env)
	}

	// Update environment field
	c.Environment = envConfig.Name

	// Set URLs from environment config (unless explicitly overridden)
	// Check if URLs are default/empty or match the previous environment
	if c.BaseURL == "" || c.isDefaultURL(c.BaseURL) {
		c.BaseURL = envConfig.APIBaseURL
	}
	if c.DashboardURL == "" || c.isDefaultURL(c.DashboardURL) {
		c.DashboardURL = envConfig.DashboardURL
	}

	return nil
}

// isDefaultURL checks if a URL is a default URL from any environment
func (c *Config) isDefaultURL(url string) bool {
	for _, env := range environments {
		if url == env.APIBaseURL || url == env.DashboardURL {
			return true
		}
	}
	return false
}

// GetAPIBaseURL returns the API base URL, respecting environment overrides
func (c *Config) GetAPIBaseURL() string {
	// Check for explicit override via environment variable
	if override := os.Getenv("VAPI_API_BASE_URL"); override != "" {
		return override
	}
	return c.BaseURL
}

// GetDashboardURL returns the dashboard URL, respecting environment overrides
func (c *Config) GetDashboardURL() string {
	// Check for explicit override via environment variable
	if override := os.Getenv("VAPI_DASHBOARD_URL"); override != "" {
		return override
	}
	return c.DashboardURL
}

// IsProduction returns true if running against production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetEnvironment returns the current environment name
func (c *Config) GetEnvironment() string {
	if c.Environment == "" {
		return "production"
	}
	return c.Environment
}

func SaveConfig(config *Config) error {
	viper.Set("api_key", config.APIKey)
	viper.Set("base_url", config.BaseURL)
	viper.Set("dashboard_url", config.DashboardURL)
	viper.Set("environment", config.Environment)
	viper.Set("timeout", config.Timeout)
	viper.Set("disable_analytics", config.DisableAnalytics)
	viper.Set("accounts", config.Accounts)
	viper.Set("active_account", config.ActiveAccount)

	// Save to home directory for persistence
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, ".vapi-cli.yaml")
	return viper.WriteConfigAs(configPath)
}

// GetActiveAPIKey returns the API key for the currently active account
func (c *Config) GetActiveAPIKey() string {
	// Check environment variable first (always takes precedence)
	if envKey := os.Getenv("VAPI_API_KEY"); envKey != "" {
		return envKey
	}

	// If we have multiple accounts and an active account is set
	if c.Accounts != nil && c.ActiveAccount != "" {
		if account, exists := c.Accounts[c.ActiveAccount]; exists {
			return account.APIKey
		}
	}

	// Fall back to legacy single API key
	return c.APIKey
}

// GetActiveAccount returns the currently active account
func (c *Config) GetActiveAccount() *Account {
	if c.Accounts != nil && c.ActiveAccount != "" {
		if account, exists := c.Accounts[c.ActiveAccount]; exists {
			return &account
		}
	}
	return nil
}

// AddAccount adds a new account or updates an existing one
func (c *Config) AddAccount(accountKey, apiKey, organization string) {
	if c.Accounts == nil {
		c.Accounts = make(map[string]Account)
	}

	c.Accounts[accountKey] = Account{
		APIKey:       apiKey,
		Organization: organization,
		Environment:  c.Environment, // Use current environment as default
		LoginTime:    time.Now().Format(time.RFC3339),
	}

	// Set as active account if it's the first one or no active account is set
	if c.ActiveAccount == "" || len(c.Accounts) == 1 {
		c.ActiveAccount = accountKey
	}
}

// SetActiveAccount switches to the specified account
func (c *Config) SetActiveAccount(accountKey string) error {
	if c.Accounts == nil {
		return fmt.Errorf("no accounts configured")
	}

	// Check if account exists (regardless of whether it has an API key)
	if _, exists := c.Accounts[accountKey]; !exists {
		return fmt.Errorf("account '%s' not found", accountKey)
	}

	c.ActiveAccount = accountKey
	return nil
}

// ListAccounts returns all configured accounts
func (c *Config) ListAccounts() map[string]Account {
	if c.Accounts == nil {
		return make(map[string]Account)
	}
	return c.Accounts
}

// RemoveAccount removes an account
func (c *Config) RemoveAccount(accountKey string) error {
	if c.Accounts == nil {
		return fmt.Errorf("no accounts configured")
	}

	if _, exists := c.Accounts[accountKey]; !exists {
		return fmt.Errorf("account '%s' not found", accountKey)
	}

	delete(c.Accounts, accountKey)

	// If we removed the active account, pick a new one or clear it
	if c.ActiveAccount == accountKey {
		if len(c.Accounts) > 0 {
			// Pick the first remaining account
			for key := range c.Accounts {
				c.ActiveAccount = key
				break
			}
		} else {
			c.ActiveAccount = ""
		}
	}

	return nil
}

// GetAPIKeySource returns where the current API key comes from
func (c *Config) GetAPIKeySource() string {
	if os.Getenv("VAPI_API_KEY") != "" {
		return "environment variable"
	}

	if c.Accounts != nil && c.ActiveAccount != "" {
		if _, exists := c.Accounts[c.ActiveAccount]; exists {
			return fmt.Sprintf("account '%s'", c.ActiveAccount)
		}
	}

	if c.APIKey != "" {
		return "config file (legacy)"
	}

	return "not set"
}

var globalConfig *Config

// GetConfig returns the global configuration instance
func GetConfig() *Config {
	return globalConfig
}

// SetConfig sets the global configuration instance
func SetConfig(config *Config) {
	globalConfig = config
}
