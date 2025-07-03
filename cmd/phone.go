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

// Manage phone numbers for making and receiving calls
var phoneCmd = &cobra.Command{
	Use:   "phone",
	Short: "Manage Vapi phone numbers",
	Long: `Manage your Vapi phone numbers for making and receiving calls.

Phone numbers are required for:
- Making outbound calls to customers
- Receiving inbound calls from customers  
- Setting up phone call campaigns
- Configuring call routing and forwarding`,
}

var listPhoneCmd = &cobra.Command{
	Use:   "list",
	Short: "List all phone numbers",
	Long:  `Display all phone numbers in your account with their status and configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("📞 Listing phone numbers...")

		ctx := context.Background()

		// Fetch phone numbers from the API
		phoneNumbers, err := vapiClient.GetClient().PhoneNumbers.List(ctx, nil)
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
			return fmt.Errorf("failed to list phone numbers: %w", err)
		}

		if len(phoneNumbers) == 0 {
			fmt.Println("No phone numbers found. Create one with 'vapi phone create'")
			return nil
		}

		// Display in a readable table format
		fmt.Printf("\nFound %d phone number(s):\n\n", len(phoneNumbers))
		fmt.Printf("%-36s %-16s %-25s %-8s %-20s\n", "ID", "Number", "Name", "Status", "Created")
		fmt.Printf("%-36s %-16s %-25s %-8s %-20s\n", "----", "------", "----", "------", "-------")

		for _, phoneNumber := range phoneNumbers {
			// Extract common fields from the union type
			fields := extractPhoneNumberFields(*phoneNumber)

			// Show assistant ID in a separate line if it exists
			fmt.Printf("%-36s %-16s %-25s %-8s %-20s\n",
				fields.ID, fields.Number, fields.Name, fields.Status, fields.CreatedAt)

			if fields.AssistantID != "None" {
				fmt.Printf("  └─ Assistant: %s\n", fields.AssistantID)
			}
		}

		return nil
	},
}

var getPhoneCmd = &cobra.Command{
	Use:   "get [phone-number-id]",
	Short: "Get details of a specific phone number",
	Long:  `Retrieve the complete configuration of a phone number including routing and settings.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		phoneNumberID := args[0]

		fmt.Printf("🔍 Getting phone number details for ID: %s\n", phoneNumberID)

		// Fetch the phone number configuration
		phoneNumber, err := vapiClient.GetClient().PhoneNumbers.Get(ctx, phoneNumberID)
		if err != nil {
			return fmt.Errorf("failed to get phone number: %w", err)
		}

		// Display as formatted JSON for easy reading
		if err := output.PrintJSON(phoneNumber); err != nil {
			return fmt.Errorf("failed to display phone number: %w", err)
		}

		return nil
	},
}

var createPhoneCmd = &cobra.Command{
	Use:   "create",
	Short: "Purchase a new phone number",
	Long: `Purchase a new phone number for your Vapi account.
	
Phone number purchase involves carrier integration and billing setup,
which is best done through the Vapi dashboard interface.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("📞 Phone number purchase...")
		fmt.Println()
		fmt.Println("Phone number purchase requires:")
		fmt.Println("- Country/region selection")
		fmt.Println("- Area code preferences")
		fmt.Println("- Billing configuration")
		fmt.Println("- Carrier integration setup")
		fmt.Println()
		fmt.Println("Purchase phone numbers through the Vapi dashboard:")
		fmt.Println("https://dashboard.vapi.ai/phone-numbers")
		fmt.Println()
		fmt.Println("Or use the Vapi API for programmatic phone number management.")

		return nil
	},
}

var updatePhoneCmd = &cobra.Command{
	Use:   "update [phone-number-id]",
	Short: "Update phone number configuration",
	Long: `Update the configuration of an existing phone number.
	
This includes routing settings, webhooks, and other phone number parameters.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		phoneNumberID := args[0]

		fmt.Printf("📝 Update phone number: %s\n", phoneNumberID)
		fmt.Println()
		fmt.Println("Phone number updates can include:")
		fmt.Println("- Inbound call routing")
		fmt.Println("- Webhook configurations")
		fmt.Println("- Assistant assignments")
		fmt.Println("- Call forwarding settings")
		fmt.Println()
		fmt.Println("Update via the Vapi dashboard:")
		fmt.Printf("https://dashboard.vapi.ai/phone-numbers/%s\n", phoneNumberID)
		fmt.Println()
		fmt.Println("Or use the Vapi API: PATCH /phone-number/{id}")

		return nil
	},
}

var deletePhoneCmd = &cobra.Command{
	Use:   "delete [phone-number-id]",
	Short: "Release a phone number",
	Long:  `Release a phone number from your account. This will stop billing and make the number unavailable.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		phoneNumberID := args[0]

		// Require explicit confirmation for destructive actions
		confirmed, err := confirmDeletion("phone number", phoneNumberID)
		if err != nil {
			return err
		}
		if !confirmed {
			return nil
		}

		fmt.Printf("🗑️  Releasing phone number with ID: %s\n", phoneNumberID)

		// Execute deletion via API
		_, err = vapiClient.GetClient().PhoneNumbers.Delete(ctx, phoneNumberID)
		if err != nil {
			return fmt.Errorf("failed to release phone number: %w", err)
		}

		fmt.Println("✅ Phone number released successfully")
		fmt.Println("Note: Billing for this number will stop within 24 hours")
		return nil
	},
}

// PhoneNumberFields represents the common fields of a phone number
type PhoneNumberFields struct {
	ID          string
	Number      string
	Name        string
	Status      string
	AssistantID string
	CreatedAt   string
}

// extractPhoneNumberFields extracts common fields from the union type phone number
func extractPhoneNumberFields(phoneNumber vapi.PhoneNumbersListResponseItem) PhoneNumberFields {
	// Handle VapiPhoneNumber
	if vapiPhone := phoneNumber.GetVapiPhoneNumber(); vapiPhone != nil {
		status := "Unknown"
		if vapiPhone.GetStatus() != nil {
			status = string(*vapiPhone.GetStatus())
		}
		return PhoneNumberFields{
			ID:          vapiPhone.GetId(),
			Number:      getStringValue(vapiPhone.GetNumber()),
			Name:        getStringValue(vapiPhone.GetName()),
			Status:      status,
			AssistantID: getStringValue(vapiPhone.GetAssistantId()),
			CreatedAt:   vapiPhone.GetCreatedAt().Format("2006-01-02 15:04"),
		}
	}

	// Handle TwilioPhoneNumber
	if twilioPhone := phoneNumber.GetTwilioPhoneNumber(); twilioPhone != nil {
		status := "Unknown"
		if twilioPhone.GetStatus() != nil {
			status = string(*twilioPhone.GetStatus())
		}
		return PhoneNumberFields{
			ID:          twilioPhone.GetId(),
			Number:      twilioPhone.GetNumber(),
			Name:        getStringValue(twilioPhone.GetName()),
			Status:      status,
			AssistantID: getStringValue(twilioPhone.GetAssistantId()),
			CreatedAt:   twilioPhone.GetCreatedAt().Format("2006-01-02 15:04"),
		}
	}

	// Handle VonagePhoneNumber
	if vonagePhone := phoneNumber.GetVonagePhoneNumber(); vonagePhone != nil {
		status := "Unknown"
		if vonagePhone.GetStatus() != nil {
			status = string(*vonagePhone.GetStatus())
		}
		return PhoneNumberFields{
			ID:          vonagePhone.GetId(),
			Number:      vonagePhone.GetNumber(),
			Name:        getStringValue(vonagePhone.GetName()),
			Status:      status,
			AssistantID: getStringValue(vonagePhone.GetAssistantId()),
			CreatedAt:   vonagePhone.GetCreatedAt().Format("2006-01-02 15:04"),
		}
	}

	// Handle TelnyxPhoneNumber
	if telnyxPhone := phoneNumber.GetTelnyxPhoneNumber(); telnyxPhone != nil {
		status := "Unknown"
		if telnyxPhone.GetStatus() != nil {
			status = string(*telnyxPhone.GetStatus())
		}
		return PhoneNumberFields{
			ID:          telnyxPhone.GetId(),
			Number:      telnyxPhone.GetNumber(),
			Name:        getStringValue(telnyxPhone.GetName()),
			Status:      status,
			AssistantID: getStringValue(telnyxPhone.GetAssistantId()),
			CreatedAt:   telnyxPhone.GetCreatedAt().Format("2006-01-02 15:04"),
		}
	}

	// Handle ByoPhoneNumber
	if byoPhone := phoneNumber.GetByoPhoneNumber(); byoPhone != nil {
		status := "Unknown"
		if byoPhone.GetStatus() != nil {
			status = string(*byoPhone.GetStatus())
		}
		return PhoneNumberFields{
			ID:          byoPhone.GetId(),
			Number:      getStringValue(byoPhone.GetNumber()),
			Name:        getStringValue(byoPhone.GetName()),
			Status:      status,
			AssistantID: getStringValue(byoPhone.GetAssistantId()),
			CreatedAt:   byoPhone.GetCreatedAt().Format("2006-01-02 15:04"),
		}
	}

	// Fallback if no phone number type is set
	return PhoneNumberFields{
		ID:          "Unknown",
		Number:      "Unknown",
		Name:        "Unknown",
		Status:      "Unknown",
		AssistantID: "Unknown",
		CreatedAt:   "Unknown",
	}
}

// getStringValue safely extracts string from pointer or returns fallback
func getStringValue(ptr *string) string {
	if ptr != nil && *ptr != "" {
		return *ptr
	}
	return "None"
}

func init() {
	rootCmd.AddCommand(phoneCmd)
	phoneCmd.AddCommand(listPhoneCmd)
	phoneCmd.AddCommand(getPhoneCmd)
	phoneCmd.AddCommand(createPhoneCmd)
	phoneCmd.AddCommand(updatePhoneCmd)
	phoneCmd.AddCommand(deletePhoneCmd)
}
