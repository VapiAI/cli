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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Manage CLI settings - API keys, default values, and environment configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Vapi CLI configuration",
	Long: `Configure the Vapi CLI settings.

The CLI looks for configuration in:
1. Environment variables (VAPI_API_KEY, etc.)
2. ./.vapi-cli.yaml (project-specific)
3. ~/.vapi-cli.yaml (global defaults)`,
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Long:  `Retrieve a specific configuration value or all settings if no key is provided.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Display all configuration values with their sources
			fmt.Println("Current configuration:")
			fmt.Println()

			settings := viper.AllSettings()
			for key, value := range settings {
				fmt.Printf("%s: %v\n", key, value)
			}

			fmt.Println()
			fmt.Printf("Config file: %s\n", viper.ConfigFileUsed())

			return nil
		}

		// Get specific configuration value
		key := args[0]
		value := viper.Get(key)

		if value == nil {
			fmt.Printf("Configuration key '%s' not found\n", key)
			return nil
		}

		fmt.Printf("%s: %v\n", key, value)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set a configuration value in the CLI config file.

Common settings:
- api_key: Your Vapi API key
- default_assistant_id: Default assistant for commands
- output_format: json, yaml, or table`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		// Update configuration
		viper.Set(key, value)

		// Save to config file
		if err := viper.WriteConfig(); err != nil {
			// If config file doesn't exist, create it
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}
		}

		fmt.Printf("✓ Set %s = %s\n", key, value)
		fmt.Printf("Configuration saved to: %s\n", viper.ConfigFileUsed())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
