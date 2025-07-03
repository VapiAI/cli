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

// Manage visual conversation flows with branching logic and variable extraction
var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage Vapi workflows",
	Long: `Manage your Vapi workflows.

Workflows are visual conversation flows that:
- Create deterministic conversation paths
- Branch based on conditions
- Extract and use variables
- Integrate with external APIs
- Handle complex multi-step interactions`,
}

var listWorkflowCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workflows",
	Long:  `Display all workflows in your account with their IDs, names, and basic configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		fmt.Println("📋 Listing workflows...")

		// Fetch all workflows from the API
		workflows, err := vapiClient.GetClient().Workflow.WorkflowControllerFindAll(ctx)
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
			return fmt.Errorf("failed to list workflows: %w", err)
		}

		if len(workflows) == 0 {
			fmt.Println("No workflows found. Create one with 'vapi workflow create'")
			return nil
		}

		// Display in a readable table format
		fmt.Printf("\nFound %d workflow(s):\n\n", len(workflows))
		fmt.Printf("%-36s %-30s %-20s\n", "ID", "Name", "Created")
		fmt.Printf("%-36s %-30s %-20s\n", "----", "----", "-------")

		for _, workflow := range workflows {
			name := workflow.Name
			if name == "" {
				name = "Unnamed"
			}

			created := workflow.CreatedAt.Format("2006-01-02 15:04")

			fmt.Printf("%-36s %-30s %-20s\n", workflow.Id, name, created)
		}

		return nil
	},
}

var createWorkflowCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new workflow",
	Long: `Create a new Vapi workflow interactively.
	
For visual workflow building with nodes, edges, and advanced configuration, 
use the Vapi dashboard at https://dashboard.vapi.ai`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔄 Create a new Vapi workflow")
		fmt.Println()

		// Basic workflow configuration via interactive prompts
		var config struct {
			Name        string
			Description string
		}

		questions := []*survey.Question{
			{
				Name:     "Name",
				Prompt:   &survey.Input{Message: "Workflow name:"},
				Validate: survey.Required,
			},
			{
				Name: "Description",
				Prompt: &survey.Input{
					Message: "Workflow description (optional):",
					Default: "A voice conversation workflow",
				},
			},
		}

		if err := survey.Ask(questions, &config); err != nil {
			return fmt.Errorf("workflow creation canceled: %w", err)
		}

		fmt.Println()
		fmt.Println("ℹ️  Note: This creates a basic workflow structure.")
		fmt.Println("   For visual workflow building with nodes and edges, use the Vapi dashboard.")
		fmt.Println()

		// Confirm before creation
		var confirmCreate bool
		confirmPrompt := &survey.Confirm{
			Message: "Create workflow with these settings?",
			Default: true,
		}

		if err := survey.AskOne(confirmPrompt, &confirmCreate); err != nil || !confirmCreate {
			fmt.Println("Creation canceled.")
			return nil
		}

		fmt.Println("\n🔄 Creating workflow...")

		ctx := context.Background()

		// Create a basic workflow with a simple conversation node
		isStart := true
		startNodeName := "start"

		conversationNode := &vapi.ConversationNode{
			Name:    startNodeName,
			IsStart: &isStart,
			Prompt:  &config.Description,
		}

		// Wrap the conversation node in the union type
		startNode := &vapi.CreateWorkflowDtoNodesItem{
			ConversationNode: conversationNode,
		}

		// Create the workflow via API
		createRequest := &vapi.CreateWorkflowDto{
			Name:  config.Name,
			Nodes: []*vapi.CreateWorkflowDtoNodesItem{startNode},
			Edges: []*vapi.Edge{}, // Empty edges for a single-node workflow
		}

		workflow, err := vapiClient.GetClient().Workflow.WorkflowControllerCreate(ctx, createRequest)
		if err != nil {
			return fmt.Errorf("failed to create workflow: %w", err)
		}

		fmt.Println("✅ Workflow created successfully!")
		fmt.Printf("ID: %s\n", workflow.Id)
		fmt.Printf("Name: %s\n", config.Name)
		fmt.Printf("Description: %s\n", config.Description)
		fmt.Println()
		fmt.Println("Your workflow is now available in the dashboard for visual editing:")
		fmt.Printf("https://dashboard.vapi.ai/workflows/%s\n", workflow.Id)
		fmt.Println()
		fmt.Println("💡 Tip: Use the visual builder to add more nodes, conditions, and connections!")

		return nil
	},
}

var getWorkflowCmd = &cobra.Command{
	Use:   "get [workflow-id]",
	Short: "Get details of a specific workflow",
	Long:  `Retrieve the full configuration of a workflow including nodes, edges, and settings.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		workflowID := args[0]

		fmt.Printf("🔍 Getting workflow details for ID: %s\n", workflowID)

		// Fetch the workflow configuration
		workflow, err := vapiClient.GetClient().Workflow.WorkflowControllerFindOne(ctx, workflowID)
		if err != nil {
			return fmt.Errorf("failed to get workflow: %w", err)
		}

		// Display as formatted JSON for easy reading
		if err := output.PrintJSON(workflow); err != nil {
			return fmt.Errorf("failed to display workflow: %w", err)
		}

		return nil
	},
}

var updateWorkflowCmd = &cobra.Command{
	Use:   "update [workflow-id]",
	Short: "Update an existing workflow",
	Long: `Update a workflow's configuration.

Complex updates involving nodes, edges, conditions, or advanced settings 
are best done through the Vapi dashboard at https://dashboard.vapi.ai`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workflowID := args[0]

		fmt.Printf("📝 Update workflow: %s\n", workflowID)
		fmt.Println()
		fmt.Println("Workflow updates are best done through the Vapi dashboard where you can:")
		fmt.Println("- Visually design conversation flows")
		fmt.Println("- Add and connect nodes")
		fmt.Println("- Configure conditions and variables")
		fmt.Println("- Test workflows interactively")
		fmt.Println()
		fmt.Println("Visit: https://dashboard.vapi.ai/workflows")

		return nil
	},
}

// nolint:dupl // Delete commands follow a similar pattern across resources
var deleteWorkflowCmd = &cobra.Command{
	Use:   "delete [workflow-id]",
	Short: "Delete a workflow",
	Long:  `Permanently delete a workflow. This cannot be undone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		workflowID := args[0]

		// Require explicit confirmation for destructive actions
		var confirmDelete bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to delete workflow %s?", workflowID),
			Default: false,
		}

		if err := survey.AskOne(prompt, &confirmDelete); err != nil {
			return fmt.Errorf("deletion canceled: %w", err)
		}

		if !confirmDelete {
			fmt.Println("Deletion canceled.")
			return nil
		}

		fmt.Printf("🗑️  Deleting workflow with ID: %s\n", workflowID)

		// Execute deletion via API
		_, err := vapiClient.GetClient().Workflow.WorkflowControllerDelete(ctx, workflowID)
		if err != nil {
			return fmt.Errorf("failed to delete workflow: %w", err)
		}

		fmt.Println("✅ Workflow deleted successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(workflowCmd)
	workflowCmd.AddCommand(listWorkflowCmd)
	workflowCmd.AddCommand(createWorkflowCmd)
	workflowCmd.AddCommand(getWorkflowCmd)
	workflowCmd.AddCommand(updateWorkflowCmd)
	workflowCmd.AddCommand(deleteWorkflowCmd)
}
