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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/VapiAI/cli/pkg/analytics"
)

// IDE configuration structures
type MCPServerConfig struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type CursorMCPConfig struct {
	MCPServers map[string]MCPServerConfig `json:"mcpServers"`
}

type VSCodeMCPConfig struct {
	Servers map[string]VSCodeServerConfig `json:"servers"`
}

type VSCodeServerConfig struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Type    string   `json:"type"`
}

// MCP command sets up Model Context Protocol integration to turn IDEs into Vapi experts
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Set up MCP integration to turn your IDE into a Vapi expert",
	Long: `Set up Model Context Protocol (MCP) integration with your IDE.

This command configures your IDE (Cursor, Windsurf, VSCode, etc.) to have access to Vapi's complete knowledge base, 
including documentation, examples, and best practices. Once configured, your IDE's AI assistant will understand 
everything about Vapi and can help you build voice AI applications without hallucinating incorrect information.

The MCP server provides access to:
- Complete Vapi documentation
- Code examples and templates  
- Best practices and guides
- API reference and troubleshooting
- Feature announcements and updates`,
	RunE: analytics.TrackCommandWrapper("mcp", "", runMCPCommand),
}

var mcpSetupCmd = &cobra.Command{
	Use:   "setup [ide]",
	Short: "Set up MCP integration for your IDE",
	Long: `Set up MCP integration for your IDE to turn it into a Vapi expert.

Supported IDEs:
- cursor: Cursor IDE
- windsurf: Windsurf IDE  
- vscode: Visual Studio Code
- auto: Auto-detect and configure all found IDEs

The setup will:
1. Install or configure the Vapi MCP docs server
2. Create the appropriate configuration files for your IDE
3. Provide instructions for completing the setup`,
	Args: cobra.MaximumNArgs(1),
	RunE: analytics.TrackCommandWrapper("mcp", "setup", runMCPSetupCommand),
}

var mcpStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check MCP integration status",
	Long:  `Check the status of MCP integration across your IDEs and show configuration details.`,
	RunE:  analytics.TrackCommandWrapper("mcp", "status", runMCPStatusCommand),
}

func runMCPCommand(cmd *cobra.Command, args []string) error {
	// Display MCP information and guide user to setup
	fmt.Println("🧠 Vapi MCP Integration")
	fmt.Println()
	fmt.Println("Turn your IDE into a Vapi expert with Model Context Protocol (MCP)!")
	fmt.Println()
	fmt.Println("What MCP does:")
	fmt.Println("• Provides your IDE's AI with complete Vapi knowledge")
	fmt.Println("• Eliminates hallucination about Vapi features and APIs")
	fmt.Println("• Enables intelligent code suggestions and examples")
	fmt.Println("• Keeps your IDE updated with latest Vapi features")
	fmt.Println()
	fmt.Println("Supported IDEs:")
	fmt.Println("• Cursor")
	fmt.Println("• Windsurf")
	fmt.Println("• Visual Studio Code")
	fmt.Println()
	fmt.Println("Get started:")
	fmt.Println("  vapi mcp setup          # Auto-detect and configure IDEs")
	fmt.Println("  vapi mcp setup cursor   # Configure Cursor specifically")
	fmt.Println("  vapi mcp status         # Check current status")
	fmt.Println()

	return nil
}

func runMCPSetupCommand(cmd *cobra.Command, args []string) error {
	var targetIDE string
	if len(args) > 0 {
		targetIDE = strings.ToLower(args[0])
	} else {
		targetIDE = "auto"
	}

	fmt.Println("🔧 Setting up Vapi MCP Integration...")
	fmt.Println()

	if targetIDE == "auto" {
		return setupAutoDetectedIDEs()
	}

	return setupSpecificIDE(targetIDE)
}

func runMCPStatusCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("📊 Vapi MCP Integration Status")
	fmt.Println()

	// Check each IDE's configuration
	checkCursorConfig()
	checkWindsurfConfig()
	checkVSCodeConfig()

	return nil
}

func setupAutoDetectedIDEs() error {
	var detectedIDEs []string
	var setupActions []func() error

	// Check for Cursor
	if cursorConfigExists() {
		detectedIDEs = append(detectedIDEs, "Cursor")
		setupActions = append(setupActions, setupCursorMCP)
	}

	// Check for Windsurf
	if windsurfConfigExists() {
		detectedIDEs = append(detectedIDEs, "Windsurf")
		setupActions = append(setupActions, setupWindsurfMCP)
	}

	// Check for VSCode
	if vscodeConfigExists() {
		detectedIDEs = append(detectedIDEs, "VSCode")
		setupActions = append(setupActions, setupVSCodeMCP)
	}

	if len(detectedIDEs) == 0 {
		fmt.Println("❌ No supported IDEs detected.")
		fmt.Println()
		fmt.Println("Please install one of the following IDEs:")
		fmt.Println("• Cursor: https://cursor.sh")
		fmt.Println("• Windsurf: https://codeium.com/windsurf")
		fmt.Println("• VSCode: https://code.visualstudio.com")
		fmt.Println()
		fmt.Println("Then run 'vapi mcp setup' again.")
		return nil
	}

	fmt.Printf("🎯 Detected IDEs: %s\n", strings.Join(detectedIDEs, ", "))
	fmt.Println()

	// Ask user which IDEs to configure
	var selectedIDEs []string
	prompt := &survey.MultiSelect{
		Message: "Which IDEs would you like to configure with Vapi MCP?",
		Options: detectedIDEs,
		Default: detectedIDEs, // Select all by default
	}

	if err := survey.AskOne(prompt, &selectedIDEs); err != nil {
		return fmt.Errorf("setup canceled: %w", err)
	}

	if len(selectedIDEs) == 0 {
		fmt.Println("No IDEs selected. Setup canceled.")
		return nil
	}

	// Execute setup for selected IDEs
	for i, ide := range detectedIDEs {
		for _, selected := range selectedIDEs {
			if ide == selected {
				fmt.Printf("⚙️  Setting up %s...\n", ide)
				if err := setupActions[i](); err != nil {
					fmt.Printf("❌ Failed to setup %s: %v\n", ide, err)
				} else {
					fmt.Printf("✅ Successfully configured %s\n", ide)
				}
				fmt.Println()
				break
			}
		}
	}

	displayPostSetupInstructions()
	return nil
}

func setupSpecificIDE(ide string) error {
	switch ide {
	case "cursor":
		fmt.Println("⚙️  Setting up Cursor...")
		if err := setupCursorMCP(); err != nil {
			return fmt.Errorf("failed to setup Cursor: %w", err)
		}
		fmt.Println("✅ Successfully configured Cursor")

	case "windsurf":
		fmt.Println("⚙️  Setting up Windsurf...")
		if err := setupWindsurfMCP(); err != nil {
			return fmt.Errorf("failed to setup Windsurf: %w", err)
		}
		fmt.Println("✅ Successfully configured Windsurf")

	case "vscode":
		fmt.Println("⚙️  Setting up VSCode...")
		if err := setupVSCodeMCP(); err != nil {
			return fmt.Errorf("failed to setup VSCode: %w", err)
		}
		fmt.Println("✅ Successfully configured VSCode")

	default:
		return fmt.Errorf("unsupported IDE: %s. Supported: cursor, windsurf, vscode", ide)
	}

	fmt.Println()
	displayPostSetupInstructions()
	return nil
}

// setupCursorLikeMCP sets up MCP for IDEs that use the Cursor-style configuration format
func setupCursorLikeMCP(configPath, ideName string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath), 0o750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read existing config or create new one
	var config CursorMCPConfig
	// #nosec G304 - configPath is constructed safely by helper functions
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			// If file exists but is invalid JSON, backup and create new
			backupPath := configPath + ".backup." + fmt.Sprintf("%d", time.Now().Unix())
			if renameErr := os.Rename(configPath, backupPath); renameErr != nil {
				fmt.Printf("⚠️  Warning: Failed to backup invalid config: %v\n", renameErr)
			} else {
				fmt.Printf("⚠️  Backed up invalid config to %s\n", backupPath)
			}
		}
	}

	// Initialize config if needed
	if config.MCPServers == nil {
		config.MCPServers = make(map[string]MCPServerConfig)
	}

	// Add Vapi MCP server
	config.MCPServers["vapi"] = MCPServerConfig{
		Command: "npx",
		Args:    []string{"-y", "@vapi-ai/mcp-docs-server"},
	}

	// Write config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("📝 Created/updated %s config: %s\n", ideName, configPath)
	return nil
}

func setupCursorMCP() error {
	return setupCursorLikeMCP(getCursorConfigPath(), "Cursor")
}

func setupWindsurfMCP() error {
	return setupCursorLikeMCP(getWindsurfConfigPath(), "Windsurf")
}

func setupVSCodeMCP() error {
	configPath := getVSCodeConfigPath()

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath), 0o750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read existing config or create new one
	var config VSCodeMCPConfig
	// #nosec G304 - configPath is constructed safely by helper functions
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			// If file exists but is invalid JSON, backup and create new
			backupPath := configPath + ".backup." + fmt.Sprintf("%d", time.Now().Unix())
			if renameErr := os.Rename(configPath, backupPath); renameErr != nil {
				fmt.Printf("⚠️  Warning: Failed to backup invalid config: %v\n", renameErr)
			} else {
				fmt.Printf("⚠️  Backed up invalid config to %s\n", backupPath)
			}
		}
	}

	// Initialize config if needed
	if config.Servers == nil {
		config.Servers = make(map[string]VSCodeServerConfig)
	}

	// Add Vapi MCP server (VSCode uses different format)
	var command string
	var args []string

	if runtime.GOOS == "windows" {
		command = "cmd"
		args = []string{"/c", "npx", "-y", "@vapi-ai/mcp-docs-server"}
	} else {
		command = "npx"
		args = []string{"-y", "@vapi-ai/mcp-docs-server"}
	}

	config.Servers["vapi"] = VSCodeServerConfig{
		Command: command,
		Args:    args,
		Type:    "stdio",
	}

	// Write config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("📝 Created/updated VSCode config: %s\n", configPath)
	return nil
}

// Helper functions to check IDE existence and get config paths
func cursorConfigExists() bool {
	return fileExists(getCursorConfigPath()) || dirExists(getCursorDir())
}

func windsurfConfigExists() bool {
	return fileExists(getWindsurfConfigPath()) || dirExists(getWindsurfDir())
}

func vscodeConfigExists() bool {
	return fileExists(getVSCodeConfigPath()) || dirExists(getVSCodeDir())
}

func getCursorDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".cursor")
}

func getWindsurfDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".codeium", "windsurf")
}

func getVSCodeDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".vscode")
}

func getCursorConfigPath() string {
	return filepath.Join(getCursorDir(), "mcp.json")
}

func getWindsurfConfigPath() string {
	return filepath.Join(getWindsurfDir(), "mcp_config.json")
}

func getVSCodeConfigPath() string {
	return filepath.Join(getVSCodeDir(), "mcp.json")
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func dirExists(path string) bool {
	if stat, err := os.Stat(path); err != nil || !stat.IsDir() {
		return false
	}
	return true
}

func checkCursorConfig() {
	fmt.Println("🔍 Cursor IDE:")
	configPath := getCursorConfigPath()

	if !fileExists(configPath) {
		fmt.Println("  ❌ Not configured")
		fmt.Printf("  📍 Expected config at: %s\n", configPath)
		fmt.Println("  💡 Run 'vapi mcp setup cursor' to configure")
	} else {
		// Check if Vapi MCP is configured
		data, err := os.ReadFile(configPath) // #nosec G304 - configPath is constructed safely by helper functions
		if err != nil {
			fmt.Printf("  ⚠️  Config file exists but cannot be read: %v\n", err)
			return
		}

		var config CursorMCPConfig
		if err := json.Unmarshal(data, &config); err != nil {
			fmt.Printf("  ⚠️  Config file exists but contains invalid JSON: %v\n", err)
			return
		}

		if server, exists := config.MCPServers["vapi"]; exists {
			fmt.Println("  ✅ Configured with Vapi MCP server")
			fmt.Printf("  📍 Config: %s\n", configPath)
			fmt.Printf("  🔧 Command: %s %s\n", server.Command, strings.Join(server.Args, " "))
		} else {
			fmt.Println("  ⚠️  Config exists but Vapi MCP not found")
			fmt.Println("  💡 Run 'vapi mcp setup cursor' to add Vapi MCP")
		}
	}
	fmt.Println()
}

func checkWindsurfConfig() {
	fmt.Println("🔍 Windsurf IDE:")
	configPath := getWindsurfConfigPath()

	if !fileExists(configPath) {
		fmt.Println("  ❌ Not configured")
		fmt.Printf("  📍 Expected config at: %s\n", configPath)
		fmt.Println("  💡 Run 'vapi mcp setup windsurf' to configure")
	} else {
		// Check if Vapi MCP is configured
		data, err := os.ReadFile(configPath) // #nosec G304 - configPath is constructed safely by helper functions
		if err != nil {
			fmt.Printf("  ⚠️  Config file exists but cannot be read: %v\n", err)
			return
		}

		var config CursorMCPConfig
		if err := json.Unmarshal(data, &config); err != nil {
			fmt.Printf("  ⚠️  Config file exists but contains invalid JSON: %v\n", err)
			return
		}

		if server, exists := config.MCPServers["vapi"]; exists {
			fmt.Println("  ✅ Configured with Vapi MCP server")
			fmt.Printf("  📍 Config: %s\n", configPath)
			fmt.Printf("  🔧 Command: %s %s\n", server.Command, strings.Join(server.Args, " "))
		} else {
			fmt.Println("  ⚠️  Config exists but Vapi MCP not found")
			fmt.Println("  💡 Run 'vapi mcp setup windsurf' to add Vapi MCP")
		}
	}
	fmt.Println()
}

func checkVSCodeConfig() {
	fmt.Println("🔍 VSCode:")
	configPath := getVSCodeConfigPath()

	if !fileExists(configPath) {
		fmt.Println("  ❌ Not configured")
		fmt.Printf("  📍 Expected config at: %s\n", configPath)
		fmt.Println("  💡 Run 'vapi mcp setup vscode' to configure")
	} else {
		// Check if Vapi MCP is configured
		data, err := os.ReadFile(configPath) // #nosec G304 - configPath is constructed safely by helper functions
		if err != nil {
			fmt.Printf("  ⚠️  Config file exists but cannot be read: %v\n", err)
			return
		}

		var config VSCodeMCPConfig
		if err := json.Unmarshal(data, &config); err != nil {
			fmt.Printf("  ⚠️  Config file exists but contains invalid JSON: %v\n", err)
			return
		}

		if server, exists := config.Servers["vapi"]; exists {
			fmt.Println("  ✅ Configured with Vapi MCP server")
			fmt.Printf("  📍 Config: %s\n", configPath)
			fmt.Printf("  🔧 Command: %s %s\n", server.Command, strings.Join(server.Args, " "))
		} else {
			fmt.Println("  ⚠️  Config exists but Vapi MCP not found")
			fmt.Println("  💡 Run 'vapi mcp setup vscode' to add Vapi MCP")
		}
	}
	fmt.Println()
}

func displayPostSetupInstructions() {
	fmt.Println("🎉 Setup Complete!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Restart your IDE to load the new MCP configuration")
	fmt.Println("2. In your IDE, look for MCP server status in settings")
	fmt.Println("3. Enable the Vapi MCP server if it's not automatically enabled")
	fmt.Println("4. Start a new chat and ask your AI assistant about Vapi!")
	fmt.Println()
	fmt.Println("Try asking your AI:")
	fmt.Println("• \"How do I create a voice assistant with Vapi?\"")
	fmt.Println("• \"Show me Vapi examples for phone calls\"")
	fmt.Println("• \"What are the latest Vapi features?\"")
	fmt.Println()
	fmt.Println("Your IDE AI assistant now has access to:")
	fmt.Println("• Complete Vapi documentation")
	fmt.Println("• Code examples and templates")
	fmt.Println("• Best practices and guides")
	fmt.Println("• API reference and troubleshooting")
	fmt.Println()
	fmt.Println("💡 Pro tip: Your AI will no longer hallucinate Vapi information!")
	fmt.Println()
	fmt.Println("Need help? Check 'vapi mcp status' or visit https://docs.vapi.ai")
}

func init() {
	// Add subcommands
	mcpCmd.AddCommand(mcpSetupCmd)
	mcpCmd.AddCommand(mcpStatusCmd)

	// Add to root command
	rootCmd.AddCommand(mcpCmd)
}
