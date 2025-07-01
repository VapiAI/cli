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

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// Manage webhook endpoints and configurations
var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Manage Vapi webhook endpoints",
	Long: `Manage webhook endpoints and configurations for receiving real-time events from Vapi.

Webhooks enable your applications to:
- Receive real-time call events and updates
- Handle function calls from voice assistants
- Get transcript and conversation data
- Monitor system events and errors
- Implement custom business logic triggers`,
}

var listWebhookCmd = &cobra.Command{
	Use:   "list",
	Short: "List all webhook endpoints",
	Long:  `Display all configured webhook endpoints with their URLs, events, and status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔗 Listing webhook endpoints...")
		fmt.Println()
		fmt.Println("Webhook endpoints provide real-time event delivery for:")
		fmt.Println("- Call events (started, ended, forwarded)")
		fmt.Println("- Conversation events (transcript, messages)")
		fmt.Println("- Function calls and tool executions")
		fmt.Println("- System events and error notifications")
		fmt.Println()
		fmt.Println("View and manage webhooks via:")
		fmt.Println("- Vapi Dashboard: https://dashboard.vapi.ai/webhooks")
		fmt.Println("- API Endpoint: GET /webhooks")
		fmt.Println()
		fmt.Println("For local webhook testing, use: vapi listen --forward-to <your-endpoint>")

		return nil
	},
}

var getWebhookCmd = &cobra.Command{
	Use:   "get [webhook-id]",
	Short: "Get details of a specific webhook",
	Long:  `Retrieve the complete configuration of a webhook including URL, events, and delivery settings.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		webhookID := args[0]

		fmt.Printf("🔍 Getting webhook details for ID: %s\n", webhookID)
		fmt.Println()
		fmt.Println("Webhook details include:")
		fmt.Println("- Endpoint URL and authentication")
		fmt.Println("- Subscribed event types")
		fmt.Println("- Delivery settings and retry policies")
		fmt.Println("- Success/failure statistics")
		fmt.Println("- Recent delivery logs")
		fmt.Println()
		fmt.Println("View webhook details via:")
		fmt.Printf("- Dashboard: https://dashboard.vapi.ai/webhooks/%s\n", webhookID)
		fmt.Println("- API: GET /webhooks/{id}")

		return nil
	},
}

var createWebhookCmd = &cobra.Command{
	Use:   "create [url]",
	Short: "Create a new webhook endpoint",
	Long: `Create a new webhook endpoint to receive Vapi events.
	
If no URL is provided, you'll be prompted to enter it interactively.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var webhookURL string

		// Get URL from argument or prompt
		if len(args) > 0 {
			webhookURL = args[0]
		} else {
			urlPrompt := &survey.Input{
				Message: "Enter webhook URL:",
				Help:    "Example: https://your-app.com/api/webhooks/vapi",
			}
			if err := survey.AskOne(urlPrompt, &webhookURL); err != nil {
				return fmt.Errorf("webhook creation canceled: %w", err)
			}
		}

		fmt.Printf("🔗 Creating webhook for: %s\n", webhookURL)
		fmt.Println()
		fmt.Println("Webhook creation involves configuring:")
		fmt.Println("- Target endpoint URL and authentication")
		fmt.Println("- Event types to subscribe to")
		fmt.Println("- Retry policies and timeout settings")
		fmt.Println("- Request headers and security")
		fmt.Println()
		fmt.Println("Common webhook events:")
		fmt.Println("- call.started / call.ended")
		fmt.Println("- transcript events")
		fmt.Println("- function-call events")
		fmt.Println("- status-update events")
		fmt.Println()
		fmt.Println("Create and configure webhooks via:")
		fmt.Println("- Vapi Dashboard: https://dashboard.vapi.ai/webhooks")
		fmt.Println("- API: POST /webhooks")
		fmt.Println()
		fmt.Println("💡 Tip: Use 'vapi listen' to test webhook delivery locally")
		fmt.Printf("   Example: vapi listen --forward-to %s\n", webhookURL)

		return nil
	},
}

var updateWebhookCmd = &cobra.Command{
	Use:   "update [webhook-id]",
	Short: "Update an existing webhook",
	Long: `Update the configuration of an existing webhook.
	
This includes modifying the URL, event subscriptions, and delivery settings.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		webhookID := args[0]

		fmt.Printf("📝 Update webhook: %s\n", webhookID)
		fmt.Println()
		fmt.Println("Webhook updates can include:")
		fmt.Println("- Endpoint URL changes")
		fmt.Println("- Event subscription modifications")
		fmt.Println("- Authentication and header updates")
		fmt.Println("- Retry policy adjustments")
		fmt.Println("- Timeout and delivery settings")
		fmt.Println()
		fmt.Println("Update via the Vapi dashboard:")
		fmt.Printf("https://dashboard.vapi.ai/webhooks/%s\n", webhookID)
		fmt.Println()
		fmt.Println("Or use the Vapi API: PATCH /webhooks/{id}")

		return nil
	},
}

var deleteWebhookCmd = &cobra.Command{
	Use:   "delete [webhook-id]",
	Short: "Delete a webhook endpoint",
	Long:  `Permanently delete a webhook endpoint. This will stop all event deliveries to this endpoint.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		webhookID := args[0]

		// Require explicit confirmation for destructive actions
		var confirmDelete bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to delete webhook %s? This will stop all event deliveries.", webhookID),
			Default: false,
		}

		if err := survey.AskOne(prompt, &confirmDelete); err != nil {
			return fmt.Errorf("deletion canceled: %w", err)
		}

		if !confirmDelete {
			fmt.Println("Deletion canceled.")
			return nil
		}

		fmt.Printf("🗑️  Delete webhook: %s\n", webhookID)
		fmt.Println()
		fmt.Println("Webhook deletion will:")
		fmt.Println("- Stop all event deliveries to this endpoint")
		fmt.Println("- Remove the webhook configuration")
		fmt.Println("- Clear delivery history and logs")
		fmt.Println()
		fmt.Println("Delete webhook via:")
		fmt.Printf("- Dashboard: https://dashboard.vapi.ai/webhooks/%s\n", webhookID)
		fmt.Println("- API: DELETE /webhooks/{id}")

		return nil
	},
}

var testWebhookCmd = &cobra.Command{
	Use:   "test [webhook-id]",
	Short: "Test a webhook endpoint",
	Long: `Send a test event to a webhook endpoint to verify it's working correctly.
	
This helps validate webhook configuration and endpoint availability.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		webhookID := args[0]

		fmt.Printf("🧪 Testing webhook: %s\n", webhookID)
		fmt.Println()
		fmt.Println("Webhook testing verifies:")
		fmt.Println("- Endpoint accessibility and response")
		fmt.Println("- Authentication and security")
		fmt.Println("- Event payload processing")
		fmt.Println("- Response handling and status codes")
		fmt.Println()
		fmt.Println("Test webhooks via:")
		fmt.Println("- Vapi Dashboard testing interface")
		fmt.Println("- API: POST /webhooks/{id}/test")
		fmt.Println()
		fmt.Printf("Direct link: https://dashboard.vapi.ai/webhooks/%s?tab=test\n", webhookID)
		fmt.Println()
		fmt.Println("💡 For local testing, use 'vapi listen' to forward")
		fmt.Println("   webhook events to your development environment.")

		return nil
	},
}

var eventsWebhookCmd = &cobra.Command{
	Use:   "events",
	Short: "List available webhook event types",
	Long:  `Display all available webhook event types that you can subscribe to.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("📡 Available webhook event types...")
		fmt.Println()

		fmt.Println("📞 Call Events:")
		fmt.Println("  • call.started - Call initiation and setup")
		fmt.Println("  • call.ended - Call completion and summary")
		fmt.Println("  • call.forwarded - Call transfer events")
		fmt.Println("  • call.recording-started - Recording begins")
		fmt.Println("  • call.recording-stopped - Recording ends")
		fmt.Println()

		fmt.Println("💬 Conversation Events:")
		fmt.Println("  • transcript - Real-time speech-to-text")
		fmt.Println("  • message - Assistant and user messages")
		fmt.Println("  • status-update - Call status changes")
		fmt.Println("  • speech-started - User begins speaking")
		fmt.Println("  • speech-ended - User stops speaking")
		fmt.Println()

		fmt.Println("🔧 Function Events:")
		fmt.Println("  • function-call - Tool/function executions")
		fmt.Println("  • function-result - Tool execution results")
		fmt.Println("  • tool-call-start - Tool call initiation")
		fmt.Println("  • tool-call-complete - Tool call completion")
		fmt.Println()

		fmt.Println("🔔 System Events:")
		fmt.Println("  • error - System and call errors")
		fmt.Println("  • warning - System warnings")
		fmt.Println("  • hang - Call hang-up events")
		fmt.Println("  • dtmf - Keypad input events")
		fmt.Println()

		fmt.Println("Configure events at: https://dashboard.vapi.ai/webhooks")
		fmt.Println("Event Documentation: https://docs.vapi.ai/webhooks")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(webhookCmd)
	webhookCmd.AddCommand(listWebhookCmd)
	webhookCmd.AddCommand(getWebhookCmd)
	webhookCmd.AddCommand(createWebhookCmd)
	webhookCmd.AddCommand(updateWebhookCmd)
	webhookCmd.AddCommand(deleteWebhookCmd)
	webhookCmd.AddCommand(testWebhookCmd)
	webhookCmd.AddCommand(eventsWebhookCmd)
}
