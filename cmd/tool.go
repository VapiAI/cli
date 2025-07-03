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

	vapi "github.com/VapiAI/server-sdk-go"
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

		if len(tools) == 0 {
			fmt.Println("No tools found. Create one with 'vapi tool create'")
			return nil
		}

		// Display in a readable table format
		fmt.Printf("\nFound %d tool(s):\n\n", len(tools))
		fmt.Printf("%-36s %-30s %-20s %-20s\n", "ID", "Name", "Type", "Created")
		fmt.Printf("%-36s %-30s %-20s %-20s\n", "----", "----", "----", "-------")

		for _, tool := range tools {
			// Extract common fields from the union type
			id, name, toolType, createdAt := extractToolFields(tool)

			fmt.Printf("%-36s %-30s %-20s %-20s\n",
				id, name, toolType, createdAt)
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
		confirmed, err := confirmDeletion("tool", fmt.Sprintf("%s (this will remove it from all assistants)", toolID))
		if err != nil {
			return err
		}
		if !confirmed {
			return nil
		}

		fmt.Printf("🗑️  Deleting tool with ID: %s\n", toolID)

		// Execute deletion via API
		_, err = vapiClient.GetClient().Tools.Delete(ctx, toolID)
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

// extractToolFields extracts common fields from the union type tool
func extractToolFields(tool *vapi.ToolsListResponseItem) (id, name, toolType, createdAt string) {
	// Handle FunctionTool
	if funcTool := tool.GetFunctionTool(); funcTool != nil {
		id = funcTool.GetId()
		name = "Unknown"
		if function := funcTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "Function"
		createdAt = funcTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle EndCallTool
	if endCallTool := tool.GetEndCallTool(); endCallTool != nil {
		id = endCallTool.GetId()
		name = "Unknown"
		if function := endCallTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "End Call"
		createdAt = endCallTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle TransferCallTool
	if transferTool := tool.GetTransferCallTool(); transferTool != nil {
		id = transferTool.GetId()
		name = "Unknown"
		if function := transferTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "Transfer Call"
		createdAt = transferTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle DtmfTool
	if dtmfTool := tool.GetDtmfTool(); dtmfTool != nil {
		id = dtmfTool.GetId()
		name = "Unknown"
		if function := dtmfTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "DTMF"
		createdAt = dtmfTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle ApiRequestTool
	if apiTool := tool.GetApiRequestTool(); apiTool != nil {
		id = apiTool.GetId()
		name = "Unknown"
		if function := apiTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "API Request"
		createdAt = apiTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle OutputTool
	if outputTool := tool.GetOutputTool(); outputTool != nil {
		id = outputTool.GetId()
		name = "Unknown"
		if function := outputTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "Output"
		createdAt = outputTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle McpTool
	if mcpTool := tool.GetMcpTool(); mcpTool != nil {
		id = mcpTool.GetId()
		name = "Unknown"
		if function := mcpTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "MCP"
		createdAt = mcpTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle BashTool
	if bashTool := tool.GetBashTool(); bashTool != nil {
		id = bashTool.GetId()
		name = "Unknown"
		if function := bashTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "Bash"
		createdAt = bashTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle ComputerTool
	if computerTool := tool.GetComputerTool(); computerTool != nil {
		id = computerTool.GetId()
		name = "Unknown"
		if function := computerTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "Computer"
		createdAt = computerTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle TextEditorTool
	if textTool := tool.GetTextEditorTool(); textTool != nil {
		id = textTool.GetId()
		name = "Unknown"
		if function := textTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "Text Editor"
		createdAt = textTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle QueryTool
	if queryTool := tool.GetQueryTool(); queryTool != nil {
		id = queryTool.GetId()
		name = "Unknown"
		if function := queryTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "Query"
		createdAt = queryTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle SlackSendMessageTool
	if slackTool := tool.GetSlackSendMessageTool(); slackTool != nil {
		id = slackTool.GetId()
		name = "Unknown"
		if function := slackTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "Slack"
		createdAt = slackTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle SmsTool
	if smsTool := tool.GetSmsTool(); smsTool != nil {
		id = smsTool.GetId()
		name = "Unknown"
		if function := smsTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "SMS"
		createdAt = smsTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle GhlTool
	if ghlTool := tool.GetGhlTool(); ghlTool != nil {
		id = ghlTool.GetId()
		name = "Unknown"
		if function := ghlTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "GoHighLevel"
		createdAt = ghlTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle MakeTool
	if makeTool := tool.GetMakeTool(); makeTool != nil {
		id = makeTool.GetId()
		name = "Unknown"
		if function := makeTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "Make"
		createdAt = makeTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle Google Calendar tools
	if calTool := tool.GetGoogleCalendarCreateEventTool(); calTool != nil {
		id = calTool.GetId()
		name = "Unknown"
		if function := calTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "Google Calendar"
		createdAt = calTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle Google Calendar Check Availability
	if calAvailTool := tool.GetGoogleCalendarCheckAvailabilityTool(); calAvailTool != nil {
		id = calAvailTool.GetId()
		name = "Unknown"
		if function := calAvailTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "Google Calendar"
		createdAt = calAvailTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle Google Sheets tools
	if sheetsTool := tool.GetGoogleSheetsRowAppendTool(); sheetsTool != nil {
		id = sheetsTool.GetId()
		name = "Unknown"
		if function := sheetsTool.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "Google Sheets"
		createdAt = sheetsTool.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle GoHighLevel Calendar Availability
	if ghlCalAvail := tool.GetGoHighLevelCalendarAvailabilityTool(); ghlCalAvail != nil {
		id = ghlCalAvail.GetId()
		name = "Unknown"
		if function := ghlCalAvail.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "GHL Calendar"
		createdAt = ghlCalAvail.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle GoHighLevel Calendar Event Create
	if ghlCalEvent := tool.GetGoHighLevelCalendarEventCreateTool(); ghlCalEvent != nil {
		id = ghlCalEvent.GetId()
		name = "Unknown"
		if function := ghlCalEvent.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "GHL Calendar"
		createdAt = ghlCalEvent.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle GoHighLevel Contact Create
	if ghlContactCreate := tool.GetGoHighLevelContactCreateTool(); ghlContactCreate != nil {
		id = ghlContactCreate.GetId()
		name = "Unknown"
		if function := ghlContactCreate.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "GHL Contact"
		createdAt = ghlContactCreate.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Handle GoHighLevel Contact Get
	if ghlContactGet := tool.GetGoHighLevelContactGetTool(); ghlContactGet != nil {
		id = ghlContactGet.GetId()
		name = "Unknown"
		if function := ghlContactGet.GetFunction(); function != nil {
			name = function.GetName()
		}
		toolType = "GHL Contact"
		createdAt = ghlContactGet.GetCreatedAt().Format("2006-01-02 15:04")
		return
	}

	// Fallback if no tool type is set
	return "Unknown", "Unknown", "Unknown", "Unknown"
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
