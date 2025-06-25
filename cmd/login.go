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

	"github.com/VapiAI/cli/pkg/auth"
)

// Authenticate with Vapi via browser-based OAuth flow
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Vapi using browser-based login",
	Long: `Opens your browser to authenticate with Vapi.

This secure authentication flow:
1. Opens dashboard.vapi.ai in your browser
2. Runs a local server to receive the auth token
3. Saves your API key for future CLI commands`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔐 Authenticating with Vapi...")
		fmt.Println()

		// Start the browser-based authentication flow
		// The Login() function handles saving the API key
		if err := auth.Login(); err != nil {
			return err
		}

		fmt.Println("\nYou can now use all Vapi CLI commands.")
		fmt.Println("• List assistants: vapi assistant list")
		fmt.Println("• View call history: vapi call list")
		fmt.Println("• Integrate with projects: vapi init")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
