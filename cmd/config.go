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
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/VapiAI/cli/pkg/analytics"
	"github.com/VapiAI/cli/pkg/config"
)

// Manage CLI settings - API keys, default values, and environment configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  `View and manage your Vapi CLI configuration settings.`,
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get configuration value(s)",
	Long:  `Display current configuration. Optionally specify a key to get a specific value.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if len(args) > 0 {
			// Get specific key
			key := args[0]
			switch key {
			case "api_key":
				if cfg.APIKey == "" {
					fmt.Println("api_key: (not set)")
				} else {
					// Show masked API key for security
					fmt.Printf("api_key: %s***\n", cfg.APIKey[:8])
				}
			case "base_url":
				fmt.Printf("base_url: %s\n", cfg.GetAPIBaseURL())
			case "dashboard_url":
				fmt.Printf("dashboard_url: %s\n", cfg.GetDashboardURL())
			case "environment":
				fmt.Printf("environment: %s\n", cfg.GetEnvironment())
			case "timeout":
				fmt.Printf("timeout: %d\n", cfg.Timeout)
			default:
				return fmt.Errorf("unknown config key: %s", key)
			}
		} else {
			// Show all configuration
			fmt.Println("Current configuration:")
			fmt.Println()

			if cfg.APIKey == "" {
				fmt.Println("api_key: (not set)")
			} else {
				fmt.Printf("api_key: %s***\n", cfg.APIKey[:8])
			}

			fmt.Printf("environment: %s\n", cfg.GetEnvironment())
			fmt.Printf("base_url: %s\n", cfg.GetAPIBaseURL())
			fmt.Printf("dashboard_url: %s\n", cfg.GetDashboardURL())
			fmt.Printf("timeout: %d\n", cfg.Timeout)

			// Show environment variables if set (for developers)
			if envVars := getRelevantEnvVars(); len(envVars) > 0 {
				fmt.Println("\nEnvironment variables:")
				for _, envVar := range envVars {
					fmt.Printf("%s: %s\n", envVar.Name, envVar.Value)
				}
			}

			fmt.Println()
			fmt.Printf("Config file: %s\n", viper.ConfigFileUsed())
		}

		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set configuration value",
	Long:  `Set a configuration value. Available keys: api_key, timeout, environment`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg, err := config.LoadConfig()
		if err != nil {
			cfg = &config.Config{} // Create new config if loading fails
		}

		switch key {
		case "api_key":
			cfg.APIKey = value
		case "timeout":
			// Parse timeout as integer
			var timeout int
			if _, err := fmt.Sscanf(value, "%d", &timeout); err != nil {
				return fmt.Errorf("timeout must be a number: %w", err)
			}
			cfg.Timeout = timeout
		case "environment":
			// Validate environment
			validEnvs := []string{"production", "staging", "development"}
			value = strings.ToLower(value)
			if value == "prod" {
				value = "production"
			}
			if value == "stage" {
				value = "staging"
			}
			if value == "dev" || value == "local" {
				value = "development"
			}

			valid := false
			for _, env := range validEnvs {
				if value == env {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid environment: %s (valid: %s)", value, strings.Join(validEnvs, ", "))
			}
			cfg.Environment = value
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✅ Set %s = %s\n", key, value)
		return nil
	},
}

// Hidden command for developers to manage environments
var configEnvCmd = &cobra.Command{
	Use:    "env [environment]",
	Short:  "Switch environment (development use)",
	Long:   `Switch between environments. For development use only.`,
	Hidden: true, // Hide from normal help output
	Args:   cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Show current environment and available options
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			fmt.Printf("Current environment: %s\n", cfg.GetEnvironment())
			fmt.Printf("API Base URL: %s\n", cfg.GetAPIBaseURL())
			fmt.Printf("Dashboard URL: %s\n", cfg.GetDashboardURL())
			fmt.Println()
			fmt.Println("Available environments:")
			fmt.Println("  production  - https://api.vapi.ai")
			fmt.Println("  staging     - https://api.staging.vapi.ai")
			fmt.Println("  development - http://localhost:3000")
			fmt.Println()
			fmt.Println("Usage: vapi config env <environment>")
			return nil
		}

		environment := args[0]
		return configSetCmd.RunE(cmd, []string{"environment", environment})
	},
}

type EnvVar struct {
	Name  string
	Value string
}

func getRelevantEnvVars() []EnvVar {
	var envVars []EnvVar
	relevantVars := []string{
		"VAPI_ENV",
		"VAPI_API_KEY",
		"VAPI_API_BASE_URL",
		"VAPI_DASHBOARD_URL",
	}

	for _, varName := range relevantVars {
		if value := os.Getenv(varName); value != "" {
			// Mask sensitive values
			displayValue := value
			if strings.Contains(varName, "API_KEY") && len(value) > 8 {
				displayValue = value[:8] + "***"
			}
			envVars = append(envVars, EnvVar{Name: varName, Value: displayValue})
		}
	}

	return envVars
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configEnvCmd)

	// Add analytics subcommand
	analyticsCmd := &cobra.Command{
		Use:   "analytics",
		Short: "Manage analytics preferences",
		Long:  `Configure whether the CLI sends anonymous usage analytics to help improve the product.`,
	}

	analyticsStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show current analytics status",
		RunE: analytics.TrackCommandWrapper("config", "analytics-status", func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()

			fmt.Println("📊 Analytics Status")
			fmt.Println()

			if analytics.IsEnabled() {
				fmt.Println("✅ Analytics: ENABLED")
				fmt.Println("   Anonymous usage data is being collected to help improve the CLI")
			} else {
				fmt.Println("🚫 Analytics: DISABLED")
				fmt.Println("   No usage data is being collected")
			}

			fmt.Println()
			fmt.Println("Configuration:")

			if cfg != nil && cfg.DisableAnalytics {
				fmt.Println("  • Config file: disabled")
			} else {
				fmt.Println("  • Config file: enabled (default)")
			}

			// Check environment variables
			envDisabled := false
			envVars := []string{"VAPI_DISABLE_ANALYTICS", "VAPI_NO_TELEMETRY", "DISABLE_TELEMETRY", "DO_NOT_TRACK"}
			for _, env := range envVars {
				if os.Getenv(env) != "" {
					fmt.Printf("  • Environment (%s): disabled\n", env)
					envDisabled = true
					break
				}
			}
			if !envDisabled {
				fmt.Println("  • Environment: enabled (default)")
			}

			fmt.Println()
			fmt.Println("Data collected (when enabled):")
			fmt.Println("  • Command usage patterns (anonymous)")
			fmt.Println("  • Error types and frequencies (hashed)")
			fmt.Println("  • Performance metrics")
			fmt.Println("  • Operating system and architecture")
			fmt.Println("  • CLI version information")
			fmt.Println()
			fmt.Println("Data NOT collected:")
			fmt.Println("  • API keys or sensitive credentials")
			fmt.Println("  • File contents or personal data")
			fmt.Println("  • User-identifiable information")
			fmt.Println("  • Specific error messages (only hashed patterns)")

			return nil
		}),
	}

	analyticsEnableCmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable analytics collection",
		RunE: analytics.TrackCommandWrapper("config", "analytics-enable", func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()
			if cfg == nil {
				cfg = &config.Config{}
			}

			cfg.DisableAnalytics = false

			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			// Update global config
			config.SetConfig(cfg)

			fmt.Println("✅ Analytics enabled")
			fmt.Println("   Anonymous usage data will be collected to help improve the CLI")
			fmt.Println("   You can disable this anytime with: vapi config analytics disable")

			return nil
		}),
	}

	analyticsDisableCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable analytics collection",
		RunE: analytics.TrackCommandWrapper("config", "analytics-disable", func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()
			if cfg == nil {
				cfg = &config.Config{}
			}

			cfg.DisableAnalytics = true

			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			// Update global config
			config.SetConfig(cfg)

			fmt.Println("🚫 Analytics disabled")
			fmt.Println("   No usage data will be collected")
			fmt.Println("   You can re-enable this anytime with: vapi config analytics enable")

			return nil
		}),
	}

	analyticsCmd.AddCommand(analyticsStatusCmd)
	analyticsCmd.AddCommand(analyticsEnableCmd)
	analyticsCmd.AddCommand(analyticsDisableCmd)

	configCmd.AddCommand(analyticsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
