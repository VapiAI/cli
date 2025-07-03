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
	vapi "github.com/VapiAI/server-sdk-go"
	"github.com/spf13/cobra"

	"github.com/VapiAI/cli/pkg/output"
)

// Manage text-based chat conversations with Vapi assistants
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Manage Vapi chat conversations",
	Long: `Manage your text-based chat conversations with Vapi assistants.

Chat conversations allow you to:
- Have text-based conversations with AI assistants
- Test assistant behavior without voice calls
- Create streaming and non-streaming chat sessions
- Manage chat history and sessions`,
}

var listChatCmd = &cobra.Command{
	Use:   "list",
	Short: "List all chat conversations",
	Long:  `Display all chat conversations in your account with their IDs, status, and metadata.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("💬 Listing chat conversations...")

		ctx := context.Background()

		// Fetch up to 50 chats from the API
		listRequest := &vapi.ChatsListRequest{
			Limit: vapi.Float64(50),
		}

		chats, err := vapiClient.GetClient().Chats.List(ctx, listRequest)
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
			return fmt.Errorf("failed to list chats: %w", err)
		}

		chatResults := chats.Results
		if len(chatResults) == 0 {
			fmt.Println("No chat conversations found. Create one with 'vapi chat create'")
			return nil
		}

		// Display in a readable table format
		fmt.Printf("\nFound %d chat conversation(s):\n\n", len(chatResults))
		fmt.Printf("%-36s %-30s %-36s %-20s\n", "ID", "Name", "Assistant ID", "Created")
		fmt.Printf("%-36s %-30s %-36s %-20s\n", "----", "----", "------------", "-------")

		for _, chat := range chatResults {
			name := "Unnamed"
			if chat.Name != nil && *chat.Name != "" {
				name = *chat.Name
			}

			assistantId := "None"
			if chat.AssistantId != nil && *chat.AssistantId != "" {
				assistantId = *chat.AssistantId
			}

			created := chat.CreatedAt.Format("2006-01-02 15:04")

			fmt.Printf("%-36s %-30s %-36s %-20s\n", chat.Id, name, assistantId, created)
		}

		return nil
	},
}

var createChatCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new chat conversation",
	Long: `Create a new text-based chat conversation with a Vapi assistant.
	
Chat creation requires specific parameters and is best done through the Vapi dashboard
where you can configure assistants, messages, and streaming options.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("💬 Creating a new chat conversation...")

		// For now, provide guidance to use the dashboard for complex chat creation
		fmt.Println("Chat creation requires specific parameters:")
		fmt.Println("- Assistant ID or Workflow ID")
		fmt.Println("- Initial message content")
		fmt.Println("- Streaming configuration")
		fmt.Println()
		fmt.Println("For full chat functionality, use the Vapi dashboard:")
		fmt.Println("https://dashboard.vapi.ai/chats")
		fmt.Println()
		fmt.Println("You can also use the Vapi SDKs for programmatic chat creation.")

		return nil
	},
}

var getChatCmd = &cobra.Command{
	Use:   "get [chat-id]",
	Short: "Get details of a specific chat conversation",
	Long:  `Retrieve the complete history and details of a chat conversation including all messages.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		chatID := args[0]

		fmt.Printf("🔍 Getting chat conversation details for ID: %s\n", chatID)

		// Fetch the chat conversation
		chat, err := vapiClient.GetClient().Chats.Get(ctx, chatID)
		if err != nil {
			return fmt.Errorf("failed to get chat: %w", err)
		}

		// Display as formatted JSON for complete details
		if err := output.PrintJSON(chat); err != nil {
			return fmt.Errorf("failed to display chat: %w", err)
		}

		return nil
	},
}

var deleteChatCmd = &cobra.Command{
	Use:   "delete [chat-id]",
	Short: "Delete a chat conversation",
	Long:  `Permanently delete a chat conversation and all its messages. This cannot be undone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		chatID := args[0]

		// Require explicit confirmation for destructive actions
		confirmed, err := confirmDeletion("chat conversation", chatID)
		if err != nil {
			return err
		}
		if !confirmed {
			return nil
		}

		fmt.Printf("🗑️  Deleting chat conversation with ID: %s\n", chatID)

		// Execute deletion via API
		_, err = vapiClient.GetClient().Chats.Delete(ctx, chatID)
		if err != nil {
			return fmt.Errorf("failed to delete chat: %w", err)
		}

		fmt.Println("✅ Chat conversation deleted successfully")
		return nil
	},
}

var continueChatCmd = &cobra.Command{
	Use:   "continue [chat-id] [message]",
	Short: "Continue an existing chat conversation",
	Long: `Send a new message to continue an existing chat conversation.
	
This operation requires specific API parameters and is best done programmatically 
using the Vapi SDKs or through the dashboard interface.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		chatID := args[0]
		message := args[1]

		fmt.Printf("💬 Continue chat conversation: %s\n", chatID)
		fmt.Printf("📝 Message: %s\n", message)
		fmt.Println()
		fmt.Println("Chat continuation requires using the Vapi API directly.")
		fmt.Println("Use the Vapi SDKs for programmatic chat interactions:")
		fmt.Println("- Node.js: @vapi-ai/server-sdk")
		fmt.Println("- Python: vapi-python")
		fmt.Println("- Go: github.com/VapiAI/server-sdk-go")
		fmt.Println()
		fmt.Println("Or use the Vapi dashboard: https://dashboard.vapi.ai/chats")

		return nil
	},
}

// confirmDeletion prompts the user for confirmation before destructive actions
func confirmDeletion(itemType, itemID string) (bool, error) {
	var confirmDelete bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to delete %s %s?", itemType, itemID),
		Default: false,
	}

	if err := survey.AskOne(prompt, &confirmDelete); err != nil {
		return false, fmt.Errorf("deletion canceled: %w", err)
	}

	if !confirmDelete {
		fmt.Println("Deletion canceled.")
		return false, nil
	}

	return true, nil
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.AddCommand(listChatCmd)
	chatCmd.AddCommand(createChatCmd)
	chatCmd.AddCommand(getChatCmd)
	chatCmd.AddCommand(deleteChatCmd)
	chatCmd.AddCommand(continueChatCmd)
}
