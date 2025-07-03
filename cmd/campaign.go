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
)

// Main campaign command
var campaignCmd = &cobra.Command{
	Use:   "campaign",
	Short: "Manage AI phone campaigns",
	Long:  `Create and manage campaigns for automated AI phone calls at scale.`,
}

// Campaign list command
var campaignListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all campaigns",
	Long:  `Display all campaigns in your Vapi account with their status and details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Fetch all campaigns from the API
		campaignsResponse, err := vapiClient.GetClient().Campaigns.CampaignControllerFindAll(ctx, &vapi.CampaignControllerFindAllRequest{})
		if err != nil {
			return fmt.Errorf("failed to list campaigns: %w", err)
		}

		campaigns := campaignsResponse.Results
		if len(campaigns) == 0 {
			fmt.Println("No campaigns found. Create one with 'vapi campaign create'")
			return nil
		}

		// Display campaigns in a formatted table
		fmt.Printf("%-36s %-20s %-15s %-15s %-20s\n", "ID", "Name", "Status", "Calls Ended", "Created")
		fmt.Println(strings.Repeat("-", 110))

		for _, campaign := range campaigns {
			status := string(campaign.Status)
			callsEnded := fmt.Sprintf("%.0f", campaign.CallsCounterEnded)
			created := campaign.CreatedAt.Format("2006-01-02 15:04")

			fmt.Printf("%-36s %-20s %-15s %-15s %-20s\n",
				campaign.Id,
				truncateString(campaign.Name, 20),
				status,
				callsEnded,
				created)
		}

		return nil
	},
}

// Campaign create command
var campaignCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new campaign",
	Long:  `Create a new campaign for automated AI phone calls.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Interactive campaign creation
		var name string
		namePrompt := &survey.Input{
			Message: "Campaign name:",
			Help:    "A name to identify your campaign",
		}
		if err := survey.AskOne(namePrompt, &name, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

		// Choose between assistant or workflow
		var useWorkflow bool
		workflowPrompt := &survey.Confirm{
			Message: "Use workflow instead of assistant?",
			Default: false,
			Help:    "Workflows allow visual conversation flow design",
		}
		if err := survey.AskOne(workflowPrompt, &useWorkflow); err != nil {
			return err
		}

		var assistantId, workflowId *string

		if useWorkflow {
			// Fetch workflows
			workflows, err := vapiClient.GetClient().Workflow.WorkflowControllerFindAll(ctx)
			if err != nil {
				return fmt.Errorf("failed to fetch workflows: %w", err)
			}

			if len(workflows) == 0 {
				fmt.Println("No workflows found. Create one first with 'vapi workflow create'")
				return nil
			}

			// Let user select a workflow
			var workflowOptions []string
			workflowMap := make(map[string]string)
			for _, workflow := range workflows {
				label := fmt.Sprintf("%s (ID: %s)", workflow.Name, workflow.Id)
				workflowOptions = append(workflowOptions, label)
				workflowMap[label] = workflow.Id
			}

			var selectedWorkflow string
			workflowSelectPrompt := &survey.Select{
				Message: "Select workflow:",
				Options: workflowOptions,
			}
			if err := survey.AskOne(workflowSelectPrompt, &selectedWorkflow); err != nil {
				return err
			}
			id := workflowMap[selectedWorkflow]
			workflowId = &id
		} else {
			// Fetch assistants
			assistants, err := vapiClient.GetClient().Assistants.List(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to fetch assistants: %w", err)
			}

			if len(assistants) == 0 {
				fmt.Println("No assistants found. Create one first with 'vapi assistant create'")
				return nil
			}

			// Let user select an assistant
			var assistantOptions []string
			assistantMap := make(map[string]string)
			for _, assistant := range assistants {
				name := "Unnamed"
				if assistant.Name != nil {
					name = *assistant.Name
				}
				label := fmt.Sprintf("%s (ID: %s)", name, assistant.Id)
				assistantOptions = append(assistantOptions, label)
				assistantMap[label] = assistant.Id
			}

			var selectedAssistant string
			assistantSelectPrompt := &survey.Select{
				Message: "Select assistant:",
				Options: assistantOptions,
			}
			if err := survey.AskOne(assistantSelectPrompt, &selectedAssistant); err != nil {
				return err
			}
			id := assistantMap[selectedAssistant]
			assistantId = &id
		}

		// Fetch phone numbers
		phoneNumbers, err := vapiClient.GetClient().PhoneNumbers.List(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to fetch phone numbers: %w", err)
		}

		if len(phoneNumbers) == 0 {
			fmt.Println("\nNo phone numbers found. You need to purchase a phone number first.")
			fmt.Println("Visit https://dashboard.vapi.ai to purchase phone numbers.")
			return nil
		}

		// Let user select a phone number
		var phoneOptions []string
		phoneMap := make(map[string]string)
		for _, phone := range phoneNumbers {
			// Handle the union type - extract common fields from each phone number type
			var phoneId, phoneNumber, phoneName string

			if phone.ByoPhoneNumber != nil {
				phoneId = phone.ByoPhoneNumber.Id
				if phone.ByoPhoneNumber.Number != nil {
					phoneNumber = *phone.ByoPhoneNumber.Number
				}
				if phone.ByoPhoneNumber.Name != nil {
					phoneName = *phone.ByoPhoneNumber.Name
				}
			} else if phone.TwilioPhoneNumber != nil {
				phoneId = phone.TwilioPhoneNumber.Id
				phoneNumber = phone.TwilioPhoneNumber.Number
				if phone.TwilioPhoneNumber.Name != nil {
					phoneName = *phone.TwilioPhoneNumber.Name
				}
			} else if phone.VonagePhoneNumber != nil {
				phoneId = phone.VonagePhoneNumber.Id
				phoneNumber = phone.VonagePhoneNumber.Number
				if phone.VonagePhoneNumber.Name != nil {
					phoneName = *phone.VonagePhoneNumber.Name
				}
			} else if phone.VapiPhoneNumber != nil {
				phoneId = phone.VapiPhoneNumber.Id
				if phone.VapiPhoneNumber.Number != nil {
					phoneNumber = *phone.VapiPhoneNumber.Number
				}
				if phone.VapiPhoneNumber.Name != nil {
					phoneName = *phone.VapiPhoneNumber.Name
				}
			} else if phone.TelnyxPhoneNumber != nil {
				phoneId = phone.TelnyxPhoneNumber.Id
				phoneNumber = phone.TelnyxPhoneNumber.Number
				if phone.TelnyxPhoneNumber.Name != nil {
					phoneName = *phone.TelnyxPhoneNumber.Name
				}
			}

			// Create display label
			label := phoneNumber
			if phoneName != "" {
				label = fmt.Sprintf("%s (%s)", phoneNumber, phoneName)
			}
			if label == "" {
				label = phoneId // Fallback to ID if no number
			}

			phoneOptions = append(phoneOptions, label)
			phoneMap[label] = phoneId
		}

		var selectedPhone string
		phoneSelectPrompt := &survey.Select{
			Message: "Select phone number to use for calls:",
			Options: phoneOptions,
		}
		if err := survey.AskOne(phoneSelectPrompt, &selectedPhone); err != nil {
			return err
		}
		phoneNumberId := phoneMap[selectedPhone]

		// For now, we'll create a basic campaign without customers
		// Advanced features like customer lists and scheduling should be done via dashboard
		createRequest := &vapi.CreateCampaignDto{
			Name:          name,
			AssistantId:   assistantId,
			WorkflowId:    workflowId,
			PhoneNumberId: phoneNumberId,
			Customers:     []*vapi.CreateCustomerDto{}, // Empty for now
		}

		// Create the campaign
		campaign, err := vapiClient.GetClient().Campaigns.CampaignControllerCreate(ctx, createRequest)
		if err != nil {
			return fmt.Errorf("failed to create campaign: %w", err)
		}

		fmt.Println("✅ Campaign created successfully!")
		fmt.Printf("ID: %s\n", campaign.Id)
		fmt.Printf("Name: %s\n", campaign.Name)
		fmt.Printf("Status: %s\n", campaign.Status)
		fmt.Println("\nNote: To add customers and configure scheduling, visit the Vapi Dashboard:")
		fmt.Printf("https://dashboard.vapi.ai/campaigns/%s\n", campaign.Id)

		return nil
	},
}

// Campaign get command
var campaignGetCmd = &cobra.Command{
	Use:   "get [campaign-id]",
	Short: "Get campaign details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		campaignID := args[0]

		// Fetch the campaign
		campaign, err := vapiClient.GetClient().Campaigns.CampaignControllerFindOne(ctx, campaignID)
		if err != nil {
			return fmt.Errorf("failed to get campaign: %w", err)
		}

		// Display campaign details
		fmt.Printf("Campaign Details:\n")
		fmt.Printf("================\n")
		fmt.Printf("ID: %s\n", campaign.Id)
		fmt.Printf("Name: %s\n", campaign.Name)
		fmt.Printf("Status: %s\n", campaign.Status)
		if campaign.EndedReason != nil {
			fmt.Printf("Ended Reason: %s\n", *campaign.EndedReason)
		}
		fmt.Printf("Phone Number ID: %s\n", campaign.PhoneNumberId)

		if campaign.AssistantId != nil {
			fmt.Printf("Assistant ID: %s\n", *campaign.AssistantId)
		}
		if campaign.WorkflowId != nil {
			fmt.Printf("Workflow ID: %s\n", *campaign.WorkflowId)
		}

		fmt.Printf("Calls Ended: %.0f\n", campaign.CallsCounterEnded)
		fmt.Printf("Created: %s\n", campaign.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", campaign.UpdatedAt.Format("2006-01-02 15:04:05"))

		if len(campaign.Customers) > 0 {
			fmt.Printf("\nCustomers: %d\n", len(campaign.Customers))
		}

		if campaign.SchedulePlan != nil {
			fmt.Printf("\nSchedule Plan configured\n")
		}

		return nil
	},
}

// Campaign update command
var campaignUpdateCmd = &cobra.Command{
	Use:   "update [campaign-id]",
	Short: "Update a campaign",
	Long:  `Update campaign details. Note: Some fields can only be updated when campaign is not in progress.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		campaignID := args[0]

		// Fetch current campaign
		campaign, err := vapiClient.GetClient().Campaigns.CampaignControllerFindOne(ctx, campaignID)
		if err != nil {
			return fmt.Errorf("failed to get campaign: %w", err)
		}

		// Check if campaign is in progress
		if campaign.Status == vapi.CampaignStatusInProgress {
			// Only allow ending the campaign
			var endCampaign bool
			endPrompt := &survey.Confirm{
				Message: "Campaign is in progress. Do you want to end it?",
				Default: false,
			}
			if err := survey.AskOne(endPrompt, &endCampaign); err != nil {
				return err
			}

			if endCampaign {
				status := "ended"
				updateRequest := &vapi.UpdateCampaignDto{
					Status: &status,
				}

				updated, err := vapiClient.GetClient().Campaigns.CampaignControllerUpdate(ctx, campaignID, updateRequest)
				if err != nil {
					return fmt.Errorf("failed to end campaign: %w", err)
				}

				fmt.Println("✅ Campaign ended successfully!")
				fmt.Printf("Status: %s\n", updated.Status)
				return nil
			}

			fmt.Println("Campaign update canceled.")
			return nil
		}

		fmt.Println("\nFor complex campaign updates (customers, scheduling, etc.), visit:")
		fmt.Printf("https://dashboard.vapi.ai/campaigns/%s\n", campaignID)

		return nil
	},
}

// Campaign delete command
var campaignDeleteCmd = &cobra.Command{
	Use:   "delete [campaign-id]",
	Short: "Delete a campaign",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		campaignID := args[0]

		// Confirm deletion
		var confirm bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to delete campaign %s?", campaignID),
			Default: false,
		}
		if err := survey.AskOne(prompt, &confirm); err != nil {
			return err
		}

		if !confirm {
			fmt.Println("Deletion canceled.")
			return nil
		}

		// Execute deletion via API
		_, err := vapiClient.GetClient().Campaigns.CampaignControllerRemove(ctx, campaignID)
		if err != nil {
			return fmt.Errorf("failed to delete campaign: %w", err)
		}

		fmt.Println("✅ Campaign deleted successfully!")
		return nil
	},
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func init() {
	rootCmd.AddCommand(campaignCmd)
	campaignCmd.AddCommand(campaignListCmd)
	campaignCmd.AddCommand(campaignCreateCmd)
	campaignCmd.AddCommand(campaignGetCmd)
	campaignCmd.AddCommand(campaignUpdateCmd)
	campaignCmd.AddCommand(campaignDeleteCmd)
}
