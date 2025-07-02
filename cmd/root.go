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
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/VapiAI/cli/pkg/analytics"
	"github.com/VapiAI/cli/pkg/client"
	"github.com/VapiAI/cli/pkg/config"
)

var (
	cfgFile     string
	vapiClient  *client.VapiClient
	bannerShown bool
)

// ASCII art banner
var asciiArt = `
  0000000          0000000            0000000                  00000000000000000           0000000  
00000000000      00000000000        0000000000               0000000000000000000000      00000000000
00000000000     000000000000       000000000000             000000000000000000000000     00000000000
0000000000000  0000000000000      00000000000000            000000000000000000000000     00000000000
0000000000000  000000000000      0000000000000000           0000000000000000000000000    00000000000
 00000000000000000000000000     000000000000000000          000000000000000000000000     00000000000
  000000000000000000000000     00000000000000000000         000000000000000000000000     00000000000
   0000000000000000000000      000000000000000000000        00000000000000000000000      00000000000
    00000000000000000000      00000000000000000000000       000000000000000000000        00000000000
     000000000000000000      0000000000000000000000000      0000000000000000             00000000000
      0000000000000000      000000000000000000000000000     000000000000                 00000000000
       00000000000000      00000000000000000000000000000    000000000000                 00000000000
        000000000000       00000000000000000000000000000    000000000000                 00000000000
         0000000000        00000000000000000000000000000     0000000000                  00000000000
          00000000           0000000000000000000000000        00000000                     0000000  
`

// Display banner with styling
func displayBanner() {
	if bannerShown {
		return
	}
	bannerShown = true

	// Create styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#62F6B5")). // Brand green color
		Bold(true)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A8F5D7")). // Light green
		Italic(true)

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")). // Gray
		Faint(true)

	// Print banner
	fmt.Println(titleStyle.Render(asciiArt))
	fmt.Println(subtitleStyle.Render("Voice AI for developers"))
	fmt.Println(versionStyle.Render("v" + version))
	fmt.Println()
}

// The main CLI command that displays help when run without subcommands
var rootCmd = &cobra.Command{
	Use:   "vapi",
	Short: "Voice AI for developers - Vapi CLI",
	Long:  `The official CLI for Vapi - build voice AI agents that make phone calls`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Track command execution
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime)
			analytics.TrackCommand("vapi", "", true, duration, "")
		}()

		// Check if --version flag is set
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Printf("vapi version %s\n", version)
			analytics.TrackEvent("version_displayed", map[string]interface{}{
				"version": version,
			})
			return nil
		}

		// Always display banner when running root command without subcommands
		displayBanner()
		// Display help by default when no subcommand is provided
		return cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Set up PersistentPreRunE here to avoid initialization cycle
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Skip validation for root command with no subcommands (just showing help)
		if cmd.Parent() == nil && len(args) == 0 && len(cmd.Commands()) > 0 {
			return nil
		}

		// Skip API key validation for commands that don't need it
		skipAuthCommands := []string{"login", "config", "init", "completion", "help", "version", "update", "mcp", "auth"}
		for _, skipCmd := range skipAuthCommands {
			if cmd.Name() == skipCmd || (cmd.Parent() != nil && cmd.Parent().Name() == skipCmd) {
				return nil
			}
		}

		// Validate API key is configured
		apiKey := viper.GetString("api_key")
		if apiKey == "" {
			printAuthPrompt()
			return fmt.Errorf("not authenticated")
		}

		// Initialize the Vapi client for API commands
		var err error
		vapiClient, err = client.NewVapiClient(apiKey)
		if err != nil {
			return fmt.Errorf("failed to initialize Vapi client: %w", err)
		}

		return nil
	}

	// Set up PersistentPostRunE for cleanup
	rootCmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		// Ensure analytics client is properly closed
		analytics.Close()
		return nil
	}

	// Global flag for config file location
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.vapi-cli.yaml or $HOME/.vapi-cli.yaml)")

	// Global flag for API key override
	rootCmd.PersistentFlags().String("api-key", "", "Vapi API key")
	if err := viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api-key")); err != nil {
		fmt.Printf("Warning: failed to bind api-key flag: %v\n", err)
	}

	// Add version flag
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")
}

// Execute runs the root command - this is the main entry point
func Execute() {
	// Execute the CLI
	if err := rootCmd.Execute(); err != nil {
		analytics.TrackError(err.Error(), map[string]interface{}{
			"command": "root",
		})
		analytics.Close()
		os.Exit(1)
	}

	// Close analytics on successful completion
	analytics.Close()
}

// Initialize viper configuration from file and environment
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in current directory first, then home
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".vapi-cli")
	}

	// Environment variables take precedence
	viper.AutomaticEnv()
	viper.SetEnvPrefix("VAPI")

	// Read config file if it exists (ignore errors)
	if err := viper.ReadInConfig(); err == nil {
		// Config file was found and read successfully
		if viper.GetBool("debug") {
			fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
		}
	}

	// Load configuration and set global config
	cfg, err := config.LoadConfig()
	if err != nil {
		// Don't fail the CLI if config loading fails
		if viper.GetBool("debug") {
			fmt.Printf("Warning: failed to load config: %v\n", err)
		}
	} else {
		config.SetConfig(cfg)
	}

	// Initialize analytics after config is loaded
	analytics.Initialize()
}

// Display instructions for authentication
func printAuthPrompt() {
	fmt.Println("🔒 Authentication required")
	fmt.Println()
	fmt.Println("You need to authenticate with Vapi to use this command.")
	fmt.Println()
	fmt.Println("Run: vapi login")
	fmt.Println()
	fmt.Println("Or set your API key manually:")
	fmt.Println("  export VAPI_API_KEY=your_api_key_here")
	fmt.Println()
}
