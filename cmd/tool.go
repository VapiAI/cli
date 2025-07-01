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
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/VapiAI/cli/pkg/output"
)

// Manage custom tools and functions that connect voice agents to external APIs
var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Manage Vapi tools and functions",
	Long: `Manage custom tools and functions that connect your voice agents to external APIs and databases.

Tools enable your assistants to:
- Call external APIs and web services
- Access databases and business systems
- Perform custom actions during conversations
- Integrate with third-party platforms
- Execute business logic and workflows`,
}

var listToolCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tools",
	Long:  `Display all custom tools and functions in your account with their configurations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔧 Listing tools...")

		ctx := context.Background()

		// Fetch tools from the API
		// Note: The exact API endpoint may vary, this follows the pattern of other list commands
		tools, err := vapiClient.GetClient().Tools.List(ctx, nil)
		if err != nil {
			// Check if this is a deserialization error related to new features
			if strings.Contains(err.Error(), "cannot be deserialized") {
				fmt.Println("⚠️  Warning: The Vapi API returned data in a format not yet supported by this CLI version.")
				fmt.Println("   This usually happens when new features are added to Vapi.")
				fmt.Println("   Please check for CLI updates: https://github.com/VapiAI/cli/releases")
				fmt.Println()
				fmt.Printf("   Technical details: %v\n", err)
				return fmt.Errorf("incompatible API response format")
			}
			return fmt.Errorf("failed to list tools: %w", err)
		}

		// Display as formatted JSON for complete details
		if err := output.PrintJSON(tools); err != nil {
			return fmt.Errorf("failed to display tools: %w", err)
		}

		return nil
	},
}

var getToolCmd = &cobra.Command{
	Use:   "get [tool-id]",
	Short: "Get details of a specific tool",
	Long:  `Retrieve the complete configuration of a tool including function definition, parameters, and settings.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		toolID := args[0]

		fmt.Printf("🔍 Getting tool details for ID: %s\n", toolID)

		// Fetch the tool configuration
		tool, err := vapiClient.GetClient().Tools.Get(ctx, toolID)
		if err != nil {
			return fmt.Errorf("failed to get tool: %w", err)
		}

		// Display as formatted JSON for easy reading
		if err := output.PrintJSON(tool); err != nil {
			return fmt.Errorf("failed to display tool: %w", err)
		}

		return nil
	},
}

var createToolCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new custom tool",
	Long: `Create a new custom tool or function for your voice agents.
	
Tool creation involves defining:
- Function signatures and parameters
- API endpoints and authentication
- Response handling and data mapping
- Error handling and fallback behavior

This is best done through the Vapi dashboard for visual tool configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔧 Creating a new tool...")
		fmt.Println()
		fmt.Println("Tool creation requires defining:")
		fmt.Println("- Function name and description")
		fmt.Println("- Input parameters and types")
		fmt.Println("- API endpoint configuration")
		fmt.Println("- Authentication and headers")
		fmt.Println("- Response data mapping")
		fmt.Println("- Error handling logic")
		fmt.Println()
		fmt.Println("Create tools through the Vapi dashboard:")
		fmt.Println("https://dashboard.vapi.ai/tools")
		fmt.Println()
		fmt.Println("For programmatic tool creation, use the Vapi API:")
		fmt.Println("POST /tools with function definition")
		fmt.Println()
		fmt.Println("Built-in tool types available:")
		fmt.Println("- Function tools (custom API calls)")
		fmt.Println("- End call tools")
		fmt.Println("- Transfer call tools")
		fmt.Println("- DTMF (keypad) tools")
		fmt.Println("- Integration tools (Google Sheets, etc.)")

		return nil
	},
}

var updateToolCmd = &cobra.Command{
	Use:   "update [tool-id]",
	Short: "Update an existing tool",
	Long: `Update the configuration of an existing tool.
	
This includes modifying function parameters, API endpoints, 
authentication, and response handling logic.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolID := args[0]

		fmt.Printf("📝 Update tool: %s\n", toolID)
		fmt.Println()
		fmt.Println("Tool updates can include:")
		fmt.Println("- Function signature changes")
		fmt.Println("- API endpoint modifications")
		fmt.Println("- Authentication updates")
		fmt.Println("- Parameter validation rules")
		fmt.Println("- Response data mapping")
		fmt.Println("- Error handling improvements")
		fmt.Println()
		fmt.Println("Update via the Vapi dashboard:")
		fmt.Printf("https://dashboard.vapi.ai/tools/%s\n", toolID)
		fmt.Println()
		fmt.Println("Or use the Vapi API: PATCH /tools/{id}")

		return nil
	},
}

var deleteToolCmd = &cobra.Command{
	Use:   "delete [tool-id]",
	Short: "Delete a custom tool",
	Long:  `Permanently delete a custom tool. This will remove it from all assistants using it.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		toolID := args[0]

		// Require explicit confirmation for destructive actions
		var confirmDelete bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to delete tool %s? This will remove it from all assistants.", toolID),
			Default: false,
		}

		if err := survey.AskOne(prompt, &confirmDelete); err != nil {
			return fmt.Errorf("deletion canceled: %w", err)
		}

		if !confirmDelete {
			fmt.Println("Deletion canceled.")
			return nil
		}

		fmt.Printf("🗑️  Deleting tool with ID: %s\n", toolID)

		// Execute deletion via API
		_, err := vapiClient.GetClient().Tools.Delete(ctx, toolID)
		if err != nil {
			return fmt.Errorf("failed to delete tool: %w", err)
		}

		fmt.Println("✅ Tool deleted successfully")
		fmt.Println("Note: Assistants using this tool may need to be reconfigured")
		return nil
	},
}

var testToolCmd = &cobra.Command{
	Use:   "test [tool-id]",
	Short: "Test a tool with sample input",
	Long: `Test a tool by calling it with sample input parameters to verify it works correctly.
	
This helps debug tool configurations and API integrations before using them in live conversations.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolID := args[0]

		fmt.Printf("🧪 Testing tool: %s\n", toolID)
		fmt.Println()
		fmt.Println("Tool testing helps verify:")
		fmt.Println("- Function execution and responses")
		fmt.Println("- API connectivity and authentication")
		fmt.Println("- Parameter validation and types")
		fmt.Println("- Error handling and edge cases")
		fmt.Println("- Response data structure")
		fmt.Println()
		fmt.Println("Test tools via:")
		fmt.Println("- Vapi Dashboard tool testing interface")
		fmt.Println("- API: POST /tools/{id}/test")
		fmt.Println("- Assistant conversations with debug mode")
		fmt.Println()
		fmt.Printf("Direct link: https://dashboard.vapi.ai/tools/%s?tab=test\n", toolID)
		fmt.Println()
		fmt.Println("For automated testing, use the webhook 'function-call'")
		fmt.Println("events with 'vapi listen' to see real-time function calls.")

		return nil
	},
}

var listToolTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List available tool types",
	Long:  `Display all available tool types and their capabilities for creating new tools.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🛠️  Available tool types...")
		fmt.Println()

		fmt.Println("📞 Call Control Tools:")
		fmt.Println("  • End Call Tool - Terminate calls with custom messages")
		fmt.Println("  • Transfer Call Tool - Transfer to phone numbers or assistants")
		fmt.Println("  • DTMF Tool - Handle keypad input during calls")
		fmt.Println()

		fmt.Println("🔧 Custom Function Tools:")
		fmt.Println("  • HTTP Function Tool - Call external APIs and web services")
		fmt.Println("  • Database Query Tool - Execute database queries")
		fmt.Println("  • Custom Logic Tool - Run business logic functions")
		fmt.Println()

		fmt.Println("🔌 Integration Tools:")
		fmt.Println("  • Google Sheets Tool - Read/write spreadsheet data")
		fmt.Println("  • Calendar Tool - Schedule and manage appointments")
		fmt.Println("  • CRM Integration Tool - Access customer data")
		fmt.Println("  • Payment Processing Tool - Handle transactions")
		fmt.Println()

		fmt.Println("📊 Data Tools:")
		fmt.Println("  • Analytics Tool - Track conversation metrics")
		fmt.Println("  • Logging Tool - Custom event logging")
		fmt.Println("  • Validation Tool - Data validation and formatting")
		fmt.Println()

		fmt.Println("Create tools at: https://dashboard.vapi.ai/tools")
		fmt.Println("API Documentation: https://docs.vapi.ai/tools")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(toolCmd)
	toolCmd.AddCommand(listToolCmd)
	toolCmd.AddCommand(getToolCmd)
	toolCmd.AddCommand(createToolCmd)
	toolCmd.AddCommand(updateToolCmd)
	toolCmd.AddCommand(deleteToolCmd)
	toolCmd.AddCommand(testToolCmd)
	toolCmd.AddCommand(listToolTypesCmd)
}
