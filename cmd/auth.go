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
	"strings"

	"github.com/spf13/cobra"

	"github.com/VapiAI/cli/pkg/analytics"
	"github.com/VapiAI/cli/pkg/auth"
	"github.com/VapiAI/cli/pkg/config"
)

// Manage authentication and account switching
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication and account switching",
	Long: `Manage your Vapi authentication credentials and account switching.

This is especially useful for users who work with multiple organizations
or need to switch between different Vapi accounts.`,
}

// Login command under auth
var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Vapi using browser-based login",
	Long: `Opens your browser to authenticate with Vapi.

This secure authentication flow:
1. Opens dashboard.vapi.ai in your browser
2. Runs a local server to receive the auth token
3. Saves your API key for future CLI commands`,
	RunE: analytics.TrackCommandWrapper("auth", "login", func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔐 Authenticating with Vapi...")
		fmt.Println()

		// Start the browser-based authentication flow
		// The Login() function handles saving the API key
		if err := auth.Login(); err != nil {
			return err
		}

		fmt.Println("\nYou can now use all Vapi CLI commands.")
		fmt.Println("• List assistants: vapi assistant list")
		fmt.Println("• View call history: vapi call list")
		fmt.Println("• Integrate with projects: vapi init")

		return nil
	}),
}

// Logout command to clear stored credentials
var logoutCmd = &cobra.Command{
	Use:   "logout [account-name]",
	Short: "Clear stored authentication credentials",
	Long: `Clear your stored API key and logout from a Vapi account.

Usage:
  vapi auth logout                  # Logout from active account
  vapi auth logout [account-name]   # Logout from specific account
  vapi auth logout --all            # Logout from all accounts

This is useful when:
- Switching to a different Vapi organization
- Changing accounts
- Clearing credentials for security reasons

After logout, you'll need to run 'vapi auth login' to authenticate again.`,
	RunE: analytics.TrackCommandWrapper("auth", "logout", func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")

		if all {
			fmt.Println("🔓 Logging out from all accounts...")
			fmt.Println()

			if err := auth.LogoutAccount(""); err != nil {
				return fmt.Errorf("logout failed: %w", err)
			}
		} else {
			accountName := ""
			if len(args) > 0 {
				accountName = args[0]
			}

			if accountName == "" {
				fmt.Println("🔓 Logging out from active account...")
			} else {
				fmt.Printf("🔓 Logging out from account '%s'...\n", accountName)
			}
			fmt.Println()

			if err := auth.LogoutAccount(accountName); err != nil {
				return fmt.Errorf("logout failed: %w", err)
			}
		}

		return nil
	}),
}

// Status command to show current authentication state
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	Long: `Display information about your current authentication status.

Shows:
- Whether you're currently authenticated
- Which account/organization you're logged into
- API key source (config file vs environment variable)
- Current environment (production, staging, development)
- API endpoints being used

This is helpful for debugging authentication issues or confirming
which account you're currently using.`,
	RunE: analytics.TrackCommandWrapper("auth", "status", func(cmd *cobra.Command, args []string) error {
		status, err := auth.GetStatus()
		if err != nil {
			return fmt.Errorf("failed to get auth status: %w", err)
		}

		fmt.Println("🔐 Authentication Status")
		fmt.Println()

		// Authentication status
		if status.IsAuthenticated {
			fmt.Println("✅ Authenticated")
		} else {
			fmt.Println("❌ Not authenticated")
		}

		// API key information
		fmt.Printf("API Key: ")
		if status.APIKeySet {
			fmt.Printf("✅ Set (via %s)\n", status.APIKeySource)
		} else {
			fmt.Printf("❌ Not set\n")
		}

		// Environment information
		fmt.Printf("Environment: %s\n", status.Environment)
		fmt.Printf("API URL: %s\n", status.BaseURL)
		fmt.Printf("Dashboard: %s\n", status.DashboardURL)

		// Multiple accounts information
		if status.TotalAccounts > 0 {
			fmt.Println()
			fmt.Printf("Accounts (%d total):\n", status.TotalAccounts)
			for accountName, account := range status.Accounts {
				active := ""
				if accountName == status.ActiveAccount {
					active = " ✓ (active)"
				}
				fmt.Printf("  • %s%s", accountName, active)
				if account.Organization != "" {
					fmt.Printf(" - %s", account.Organization)
				}
				if account.LoginTime != "" {
					fmt.Printf(" (logged in: %s)", account.LoginTime[:10]) // Show just the date
				}
				fmt.Println()
			}
		}

		fmt.Println()

		// Action suggestions
		if !status.IsAuthenticated {
			fmt.Println("💡 To authenticate, run: vapi auth login")
		} else {
			if status.TotalAccounts > 1 {
				fmt.Println("💡 To switch accounts, run: vapi auth switch")
			}
			fmt.Println("💡 To add another account, run: vapi auth login")
			fmt.Println("💡 To logout, run: vapi auth logout")
		}

		return nil
	}),
}

// Switch active account command
var switchCmd = &cobra.Command{
	Use:   "switch [account-name]",
	Short: "Switch active Vapi account",
	Long: `Switch between multiple authenticated Vapi accounts.

If no account name is provided, you'll be prompted to choose from available accounts.
This is useful when working with multiple organizations or environments.`,
	RunE: analytics.TrackCommandWrapper("auth", "switch", func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		accounts := cfg.ListAccounts()
		if len(accounts) == 0 {
			fmt.Println("❌ No accounts found. Use 'vapi auth login' to authenticate.")
			return nil
		}

		if len(accounts) == 1 {
			// Only one account, just confirm it's active
			for accountName := range accounts {
				if err := cfg.SetActiveAccount(accountName); err != nil {
					return fmt.Errorf("failed to set active account: %w", err)
				}
				if err := config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("failed to save config: %w", err)
				}
				fmt.Printf("✅ Using account '%s' (only account available)\n", accountName)
				return nil
			}
		}

		var targetAccount string

		if len(args) > 0 {
			targetAccount = args[0]
		} else {
			// Interactive selection
			fmt.Println("Available accounts:")
			for accountName, account := range accounts {
				active := ""
				if accountName == cfg.ActiveAccount {
					active = " (active)"
				}
				fmt.Printf("  • %s%s", accountName, active)
				if account.Organization != "" {
					fmt.Printf(" - %s", account.Organization)
				}
				fmt.Println()
			}
			fmt.Print("\nEnter account name to switch to: ")
			if _, err := fmt.Scanln(&targetAccount); err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
		}

		if targetAccount == "" {
			return fmt.Errorf("no account specified")
		}

		if err := cfg.SetActiveAccount(targetAccount); err != nil {
			return err
		}

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✅ Switched to account '%s'\n", targetAccount)
		return nil
	}),
}

// Token command to show current API key
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Display the authentication token for the active account",
	Long: `Print the API key that the CLI uses for authenticating with Vapi.

This shows the actual API key being used (masked for security).
Useful for debugging authentication issues or confirming which token is active.`,
	RunE: analytics.TrackCommandWrapper("auth", "token", func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		apiKey := cfg.GetActiveAPIKey()
		if apiKey == "" {
			fmt.Println("❌ No API key found. Use 'vapi auth login' to authenticate.")
			return nil
		}

		source := cfg.GetAPIKeySource()
		fmt.Printf("API Key: %s\n", maskAPIKey(apiKey))
		fmt.Printf("Source: %s\n", source)

		if cfg.ActiveAccount != "" {
			fmt.Printf("Account: %s\n", cfg.ActiveAccount)
		}

		return nil
	}),
}

// Who am I command to show current user information (if we can get it from API)
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user and organization information",
	Long: `Display information about the currently authenticated user and organization.

This command makes an API call to fetch:
- Your user information
- Current organization details
- Account permissions and roles

Requires valid authentication to work.`,
	RunE: analytics.TrackCommandWrapper("auth", "whoami", func(cmd *cobra.Command, args []string) error {
		// Check authentication status first
		status, err := auth.GetStatus()
		if err != nil {
			return fmt.Errorf("failed to check authentication status: %w", err)
		}

		if !status.IsAuthenticated {
			return fmt.Errorf("not authenticated - run 'vapi login' first")
		}

		fmt.Println("👤 Current User Information")
		fmt.Println()

		// TODO: When the Vapi API has a /me or /user endpoint, we can fetch real user info
		// For now, we'll show what we can determine from the authentication status

		fmt.Printf("Environment: %s\n", status.Environment)
		fmt.Printf("API URL: %s\n", status.BaseURL)
		fmt.Printf("Dashboard: %s\n", status.DashboardURL)
		fmt.Printf("API Key Source: %s\n", status.APIKeySource)

		if status.ActiveAccount != "" {
			fmt.Printf("Active Account: %s\n", status.ActiveAccount)
		}

		// Extract some info from the API key if possible (safely)
		if status.APIKeySet {
			// Vapi API keys typically start with a prefix - we can show that safely
			fmt.Printf("API Key: %s\n", maskAPIKey(getAPIKeyForDisplay()))
		}

		fmt.Println()
		fmt.Printf("💡 For organization details, visit: %s\n", status.DashboardURL)

		return nil
	}),
}

// Helper function to safely mask API key for display
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return strings.Repeat("*", len(apiKey))
	}
	// Show first 4 characters and mask the rest
	return apiKey[:4] + strings.Repeat("*", len(apiKey)-4)
}

// Helper function to get API key for display purposes
func getAPIKeyForDisplay() string {
	cfg, err := config.LoadConfig()
	if err != nil {
		return ""
	}

	return cfg.GetActiveAPIKey()
}

func init() {
	// Add flags to logout command
	logoutCmd.Flags().Bool("all", false, "Logout from all accounts")

	// Add subcommands to auth
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(switchCmd)
	authCmd.AddCommand(tokenCmd)
	authCmd.AddCommand(whoamiCmd)

	// Add auth command to root
	rootCmd.AddCommand(authCmd)
}
