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

// Manage phone calls - list call history, get recordings, and create outbound calls
var callCmd = &cobra.Command{
	Use:   "call",
	Short: "Manage Vapi calls",
	Long: `Manage your Vapi phone calls.

View call history, access recordings and transcripts, 
and initiate new outbound calls programmatically.`,
}

var listCallsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all calls",
	Long:  `Display your call history including status, duration, and participants.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		fmt.Println("Listing calls...")

		// Fetch up to 50 calls from the API
		listRequest := &vapi.CallsListRequest{
			Limit: vapi.Float64(50),
		}

		calls, err := vapiClient.GetClient().Calls.List(ctx, listRequest)
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
			return fmt.Errorf("failed to list calls: %w", err)
		}

		// Display as formatted JSON for complete details
		if err := output.PrintJSON(calls); err != nil {
			return fmt.Errorf("failed to display calls: %w", err)
		}

		return nil
	},
}

var createCallCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new call",
	Long: `Initiate an outbound call.

Required parameters:
- Phone number ID (from your Vapi account)
- Customer phone number
- Assistant ID to handle the call`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Creating a new call...")

		// Full call creation requires specific parameters
		// This is a placeholder showing what's needed
		fmt.Println("Call creation requires specific parameters:")
		fmt.Println("- phoneNumberId: Your Vapi phone number ID")
		fmt.Println("- customer: { number: \"+1234567890\" }")
		fmt.Println("- assistantId: Your assistant ID")
		fmt.Println()
		fmt.Println("Use the Vapi dashboard or API documentation to create calls with proper parameters.")
		return nil
	},
}

var getCallCmd = &cobra.Command{
	Use:   "get [call-id]",
	Short: "Get details of a specific call",
	Long:  `Retrieve complete details for a call including transcript, recording URL, and metadata.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		callID := args[0]

		fmt.Printf("Getting call with ID: %s\n", callID)

		// Fetch detailed call information
		call, err := vapiClient.GetClient().Calls.Get(ctx, callID)
		if err != nil {
			return fmt.Errorf("failed to get call: %w", err)
		}

		// Display complete call details
		if err := output.PrintJSON(call); err != nil {
			return fmt.Errorf("failed to display call: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(callCmd)
	callCmd.AddCommand(listCallsCmd)
	callCmd.AddCommand(createCallCmd)
	callCmd.AddCommand(getCallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// callCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// callCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
