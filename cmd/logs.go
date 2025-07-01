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
)

// View system and call logs for debugging and monitoring
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View Vapi system and call logs",
	Long: `View system logs, call logs, and debugging information for your Vapi account.

Logs are essential for:
- Debugging call issues and errors
- Monitoring system performance
- Analyzing call patterns and behavior
- Troubleshooting webhook delivery`,
}

var listLogsCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent system logs",
	Long:  `Display recent system logs including errors, warnings, and operational events.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("📋 Listing system logs...")
		fmt.Println()
		fmt.Println("System logs provide insights into:")
		fmt.Println("- API request/response patterns")
		fmt.Println("- Error conditions and debugging info")
		fmt.Println("- Performance metrics and timing")
		fmt.Println("- Webhook delivery status")
		fmt.Println()
		fmt.Println("Access detailed logs via:")
		fmt.Println("- Vapi Dashboard: https://dashboard.vapi.ai/logs")
		fmt.Println("- API Endpoint: GET /logs")
		fmt.Println("- Real-time monitoring tools")
		fmt.Println()
		fmt.Println("For call-specific logs, use: vapi logs calls")

		return nil
	},
}

var callLogsCmd = &cobra.Command{
	Use:   "calls [call-id]",
	Short: "View logs for specific calls",
	Long: `Display detailed logs for a specific call, including transcripts, events, and debugging information.
	
If no call-id is provided, shows logs for recent calls.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var callID string
		if len(args) > 0 {
			callID = args[0]
		}

		if callID != "" {
			fmt.Printf("📞 Call logs for: %s\n", callID)
		} else {
			fmt.Println("📞 Recent call logs...")
		}
		fmt.Println()
		fmt.Println("Call logs include:")
		fmt.Println("- Real-time transcripts and conversations")
		fmt.Println("- Assistant responses and reasoning")
		fmt.Println("- Function calls and tool executions")
		fmt.Println("- Audio processing and speech recognition")
		fmt.Println("- Error conditions and troubleshooting info")
		fmt.Println()
		fmt.Println("Access detailed call logs via:")
		fmt.Println("- Vapi Dashboard call details page")
		fmt.Println("- API: GET /call/{id}/logs")
		fmt.Println("- Real-time webhook events during calls")
		fmt.Println()
		if callID != "" {
			fmt.Printf("Direct link: https://dashboard.vapi.ai/calls/%s\n", callID)
		}

		return nil
	},
}

var errorLogsCmd = &cobra.Command{
	Use:   "errors",
	Short: "View recent error logs",
	Long:  `Display recent error logs to help with debugging and troubleshooting.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔥 Recent error logs...")
		fmt.Println()
		fmt.Println("Error logs help diagnose:")
		fmt.Println("- API authentication issues")
		fmt.Println("- Call setup and connection problems")
		fmt.Println("- Assistant configuration errors")
		fmt.Println("- Webhook delivery failures")
		fmt.Println("- Integration and tool call issues")
		fmt.Println()
		fmt.Println("Access error logs via:")
		fmt.Println("- Vapi Dashboard: https://dashboard.vapi.ai/logs?level=error")
		fmt.Println("- API Endpoint: GET /logs?level=error")
		fmt.Println("- Email notifications (if configured)")
		fmt.Println()
		fmt.Println("For immediate error notifications, configure webhooks")
		fmt.Println("to receive real-time error events in your application.")

		return nil
	},
}

var webhookLogsCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "View webhook delivery logs",
	Long:  `Display webhook delivery logs including successful deliveries, failures, and retry attempts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔗 Webhook delivery logs...")
		fmt.Println()
		fmt.Println("Webhook logs track:")
		fmt.Println("- Successful webhook deliveries")
		fmt.Println("- Failed delivery attempts and reasons")
		fmt.Println("- Retry attempts and backoff strategies")
		fmt.Println("- Response codes from your endpoints")
		fmt.Println("- Payload content and timing")
		fmt.Println()
		fmt.Println("Access webhook logs via:")
		fmt.Println("- Vapi Dashboard: https://dashboard.vapi.ai/webhooks")
		fmt.Println("- API Endpoint: GET /webhook-logs")
		fmt.Println()
		fmt.Println("Tip: Use 'vapi listen' to test webhook delivery locally")
		fmt.Println("This helps debug webhook handling without deployment.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.AddCommand(listLogsCmd)
	logsCmd.AddCommand(callLogsCmd)
	logsCmd.AddCommand(errorLogsCmd)
	logsCmd.AddCommand(webhookLogsCmd)
}
