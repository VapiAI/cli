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

	"github.com/AlecAivazis/survey/v2"
	vapi "github.com/VapiAI/server-sdk-go"
	"github.com/spf13/cobra"

	"github.com/VapiAI/cli/pkg/output"
)

// Manage AI voice assistants that handle phone calls and conversations
var assistantCmd = &cobra.Command{
	Use:   "assistant",
	Short: "Manage Vapi assistants",
	Long: `Manage your Vapi AI assistants.

Assistants are AI-powered agents that can:
- Make and receive phone calls
- Handle natural conversations
- Integrate with your tools and APIs
- Process voice input in real-time`,
}

var listAssistantCmd = &cobra.Command{
	Use:   "list",
	Short: "List all assistants",
	Long:  `Display all assistants in your account with their IDs, names, and metadata.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		fmt.Println("📋 Listing assistants...")

		// Fetch up to 50 assistants from the API
		listRequest := &vapi.AssistantsListRequest{
			Limit: vapi.Float64(50),
		}

		assistants, err := vapiClient.GetClient().Assistants.List(ctx, listRequest)
		if err != nil {
			return fmt.Errorf("failed to list assistants: %w", err)
		}

		if len(assistants) == 0 {
			fmt.Println("No assistants found. Create one with 'vapi assistant create'")
			return nil
		}

		// Display in a readable table format
		fmt.Printf("\nFound %d assistant(s):\n\n", len(assistants))
		fmt.Printf("%-36s %-30s %-20s\n", "ID", "Name", "Created")
		fmt.Printf("%-36s %-30s %-20s\n", "----", "----", "-------")

		for _, assistant := range assistants {
			name := "Unnamed"
			if assistant.Name != nil {
				name = *assistant.Name
			}

			created := assistant.CreatedAt.Format("2006-01-02 15:04")

			fmt.Printf("%-36s %-30s %-20s\n", assistant.Id, name, created)
		}

		return nil
	},
}

var createAssistantCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new assistant",
	Long: `Create a new Vapi assistant interactively.
	
For advanced configuration (voice selection, model parameters, tools), 
use the Vapi dashboard at https://dashboard.vapi.ai`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🤖 Create a new Vapi assistant")
		fmt.Println()

		// Basic assistant configuration via interactive prompts
		var config struct {
			Name         string
			FirstMessage string
		}

		questions := []*survey.Question{
			{
				Name:     "Name",
				Prompt:   &survey.Input{Message: "Assistant name:"},
				Validate: survey.Required,
			},
			{
				Name: "FirstMessage",
				Prompt: &survey.Input{
					Message: "First message (greeting):",
					Default: "Hello! How can I help you today?",
				},
			},
		}

		if err := survey.Ask(questions, &config); err != nil {
			return fmt.Errorf("assistant creation canceled: %w", err)
		}

		fmt.Println()
		fmt.Println("ℹ️  Note: This creates a basic assistant with default settings.")
		fmt.Println("   For advanced configuration (model, voice, prompts), use the Vapi dashboard.")
		fmt.Println()

		// Confirm before creation
		var confirmCreate bool
		confirmPrompt := &survey.Confirm{
			Message: "Create assistant with these settings?",
			Default: true,
		}

		if err := survey.AskOne(confirmPrompt, &confirmCreate); err != nil || !confirmCreate {
			fmt.Println("Creation canceled.")
			return nil
		}

		fmt.Println("\n🔄 Creating assistant...")
		fmt.Println("Note: Full assistant creation with model/voice configuration is available in the Vapi dashboard.")
		fmt.Printf("\nCreated assistant would have:\n")
		fmt.Printf("- Name: %s\n", config.Name)
		fmt.Printf("- First Message: %s\n", config.FirstMessage)
		fmt.Println("\nVisit https://dashboard.vapi.ai to create and configure assistants.")

		return nil
	},
}

var getAssistantCmd = &cobra.Command{
	Use:   "get [assistant-id]",
	Short: "Get details of a specific assistant",
	Long:  `Retrieve the full configuration of an assistant including voice, model, and tool settings.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		assistantID := args[0]

		fmt.Printf("🔍 Getting assistant details for ID: %s\n", assistantID)

		// Fetch the assistant configuration
		assistant, err := vapiClient.GetClient().Assistants.Get(ctx, assistantID)
		if err != nil {
			return fmt.Errorf("failed to get assistant: %w", err)
		}

		// Display as formatted JSON for easy reading
		if err := output.PrintJSON(assistant); err != nil {
			return fmt.Errorf("failed to display assistant: %w", err)
		}

		return nil
	},
}

var updateAssistantCmd = &cobra.Command{
	Use:   "update [assistant-id]",
	Short: "Update an existing assistant",
	Long: `Update an assistant's configuration.

Complex updates involving voice models, tools, or advanced settings 
are best done through the Vapi dashboard at https://dashboard.vapi.ai`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		assistantID := args[0]

		fmt.Printf("📝 Update assistant: %s\n", assistantID)
		fmt.Println()
		fmt.Println("Assistant updates are best done through the Vapi dashboard where you can:")
		fmt.Println("- Configure model settings (GPT-4, Claude, etc.)")
		fmt.Println("- Select and customize voices")
		fmt.Println("- Set up tools and functions")
		fmt.Println("- Configure advanced behaviors")
		fmt.Println()
		fmt.Println("Visit: https://dashboard.vapi.ai/assistants")

		return nil
	},
}

var deleteAssistantCmd = &cobra.Command{
	Use:   "delete [assistant-id]",
	Short: "Delete an assistant",
	Long:  `Permanently delete an assistant. This cannot be undone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		assistantID := args[0]

		// Require explicit confirmation for destructive actions
		var confirmDelete bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to delete assistant %s?", assistantID),
			Default: false,
		}

		if err := survey.AskOne(prompt, &confirmDelete); err != nil {
			return fmt.Errorf("deletion canceled: %w", err)
		}

		if !confirmDelete {
			fmt.Println("Deletion canceled.")
			return nil
		}

		fmt.Printf("🗑️  Deleting assistant with ID: %s\n", assistantID)

		// Execute deletion via API
		_, err := vapiClient.GetClient().Assistants.Delete(ctx, assistantID)
		if err != nil {
			return fmt.Errorf("failed to delete assistant: %w", err)
		}

		fmt.Println("✅ Assistant deleted successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(assistantCmd)
	assistantCmd.AddCommand(listAssistantCmd)
	assistantCmd.AddCommand(createAssistantCmd)
	assistantCmd.AddCommand(getAssistantCmd)
	assistantCmd.AddCommand(updateAssistantCmd)
	assistantCmd.AddCommand(deleteAssistantCmd)
}
