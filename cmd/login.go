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

	"github.com/VapiAI/cli/pkg/analytics"
	"github.com/VapiAI/cli/pkg/auth"
)

// Authenticate with Vapi using browser-based login flow
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Vapi using secure browser login",
	Long: `🔐 Authenticate with Vapi using secure browser-based login

This command will:
1. Open your default browser to the Vapi authentication page
2. Allow you to sign in with your existing Vapi account
3. Automatically save your authentication credentials locally
4. Enable access to all Vapi CLI features

🆕 New to Vapi? 
  Sign up for free at: https://dashboard.vapi.ai

🔑 Alternative Authentication:
  You can also set your API key directly:
  export VAPI_API_KEY=your_api_key_here

🛡️  Security:
  Your credentials are stored securely in your system's credential store
  and are only used to authenticate with Vapi's APIs.`,
	RunE: analytics.TrackCommandWrapper("auth", "login", func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔐 Authenticating with Vapi...")
		fmt.Println()

		// Start the browser-based authentication flow
		// The Login() function handles saving the API key
		if err := auth.Login(); err != nil {
			return err
		}

		fmt.Println("\n✅ Authentication successful!")
		fmt.Println()
		fmt.Println("🎯 Next steps:")
		fmt.Println("• List assistants: vapi assistant list")
		fmt.Println("• View call history: vapi call list")
		fmt.Println("• Initialize project: vapi init")
		fmt.Println("• Set up IDE integration: vapi mcp setup")
		fmt.Println()
		fmt.Println("💡 For more commands, try: vapi --help")

		return nil
	}),
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
