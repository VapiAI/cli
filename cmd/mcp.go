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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/VapiAI/cli/pkg/analytics"
	"github.com/VapiAI/cli/pkg/config"
)

const vapiMCPURL = "https://mcp-v2.vapi.ai/api/mcp"

type MCPServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
	Type    string            `json:"type,omitempty"`
}

var mcpCmd = newMCPCommand()

func newMCPCommand() *cobra.Command {
	description := "Set up the Vapi v2 MCP server"
	command := &cobra.Command{
		Use:   "mcp",
		Short: description,
		Long: description + ".\n\nUses Streamable HTTP through mcp-remote; Node.js and npx are required.\n" +
			"Use 'vapi skills setup' to install public agent skills.",
		RunE: func(cmd *cobra.Command, args []string) error { return cmd.Help() },
	}
	command.AddCommand(&cobra.Command{
		Use:   "setup [ide]",
		Short: description,
		Long: description + ".\n\nSupported IDEs: cursor, windsurf, vscode. Omit the IDE to select detected IDEs.\n" +
			"The integration uses your current Vapi API key.",
		Args: cobra.MaximumNArgs(1),
		RunE: analytics.TrackCommandWrapper("mcp", "setup", func(cmd *cobra.Command, args []string) error {
			key := mcpAPIKey(cmd)
			if key == "" {
				return fmt.Errorf("a Vapi API key is required: run 'vapi login' or set VAPI_API_KEY")
			}
			return setupMCPIntegration(cmd, args, key)
		}),
	})
	command.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Check configuration without connecting to the server",
		Args:  cobra.NoArgs,
		RunE: analytics.TrackCommandWrapper("mcp", "status", func(cmd *cobra.Command, args []string) error {
			for _, ide := range []string{"cursor", "windsurf", "vscode"} {
				if err := printMCPStatus(cmd, ide); err != nil {
					return err
				}
			}
			return nil
		}),
	})
	return command
}

func mcpAPIKey(cmd *cobra.Command) string {
	if cmd.Flags().Changed("api-key") {
		key, _ := cmd.Flags().GetString("api-key")
		return key
	}
	if cfg := config.GetConfig(); cfg != nil {
		return cfg.GetActiveAPIKey()
	}
	return viper.GetString("api_key")
}

func setupMCPIntegration(cmd *cobra.Command, args []string, key string) error {
	ides := []string{}
	if len(args) == 1 && strings.ToLower(args[0]) != "auto" {
		ide := strings.ToLower(args[0])
		if _, err := mcpConfigPath(ide); err != nil {
			return err
		}
		ides = append(ides, ide)
	} else {
		for _, ide := range []string{"cursor", "windsurf", "vscode"} {
			configPath, err := mcpConfigPath(ide)
			if err != nil {
				return err
			}
			if info, err := os.Stat(filepath.Dir(configPath)); err == nil && info.IsDir() {
				ides = append(ides, ide)
			}
		}
		if len(ides) == 0 {
			return fmt.Errorf("no supported IDEs detected; run 'vapi mcp setup cursor', 'windsurf', or 'vscode' explicitly")
		}
		var selected []string
		if err := survey.AskOne(&survey.MultiSelect{
			Message: "Which IDEs should be configured?", Options: ides, Default: ides,
		}, &selected); err != nil {
			return fmt.Errorf("setup canceled: %w", err)
		}
		ides = selected
	}
	if len(ides) == 0 {
		return nil
	}
	for _, ide := range ides {
		configPath, err := mcpConfigPath(ide)
		if err != nil {
			return err
		}
		server := remoteMCPConfig(key, runtime.GOOS)
		section := "mcpServers"
		if ide == "vscode" {
			section = "servers"
			server.Type = "stdio"
		}
		if err := writeMCPConfig(configPath, section, server); err != nil {
			return fmt.Errorf("configure %s: %w", ide, err)
		}
		cmd.Printf("Configured %s (vapi): %s\n", ide, configPath)
	}
	cmd.Printf("Endpoint: %s (Streamable HTTP)\nRestart your IDE and enable vapi in MCP settings.\n", vapiMCPURL)
	cmd.Println("Your current API key is stored in the user MCP configuration. Keep the file private and re-run setup after changing accounts or rotating the key.")
	return nil
}

func remoteMCPConfig(key, goos string) MCPServerConfig {
	server := MCPServerConfig{Command: "npx", Args: []string{"-y", "mcp-remote", vapiMCPURL, "--transport", "http-only"}}
	server.Args = append(server.Args, "--header", "Authorization: Bearer ${VAPI_API_KEY}")
	server.Env = map[string]string{"VAPI_API_KEY": key}
	if goos == "windows" {
		server.Command = "cmd"
		server.Args = append([]string{"/c", "npx"}, server.Args...)
	}
	return server
}

// Update one entry without decoding away other servers' fields or editor settings.
// Invalid JSON is left untouched so setup cannot silently replace a user's config.
func writeMCPConfig(configPath, section string, server MCPServerConfig) error {
	root := map[string]json.RawMessage{}
	data, err := os.ReadFile(configPath) // #nosec G304 -- path is constructed by mcpConfigPath
	if err == nil {
		if err := json.Unmarshal(data, &root); err != nil || root == nil {
			return fmt.Errorf("invalid JSON object in %s; fix the file before running setup", configPath)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	servers := map[string]json.RawMessage{}
	if existing, ok := root[section]; ok {
		if err := json.Unmarshal(existing, &servers); err != nil || servers == nil {
			return fmt.Errorf("invalid %s object in %s", section, configPath)
		}
	}
	encoded, err := json.Marshal(server)
	if err != nil {
		return err
	}
	servers["vapi"] = encoded
	root[section], err = json.Marshal(servers)
	if err != nil {
		return err
	}
	output, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(configPath), 0o750); err != nil {
		return err
	}
	// Replace only after a complete write, so a failure cannot truncate existing settings.
	file, err := os.CreateTemp(filepath.Dir(configPath), ".vapi-mcp-*.json")
	if err != nil {
		return err
	}
	temporaryPath := file.Name()
	defer os.Remove(temporaryPath) //nolint:errcheck // already renamed on success
	if _, err := file.Write(append(output, '\n')); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, configPath)
}

func printMCPStatus(cmd *cobra.Command, ide string) error {
	configPath, err := mcpConfigPath(ide)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(configPath) // #nosec G304 -- path is constructed by mcpConfigPath
	if errors.Is(err, os.ErrNotExist) {
		cmd.Printf("%s: Not configured (%s)\n", ide, configPath)
		return nil
	}
	if err != nil {
		return err
	}
	var root map[string]json.RawMessage
	if err := json.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("read %s: %w", configPath, err)
	}
	section := "mcpServers"
	if ide == "vscode" {
		section = "servers"
	}
	servers := map[string]MCPServerConfig{}
	if data, ok := root[section]; ok {
		if err := json.Unmarshal(data, &servers); err != nil {
			return fmt.Errorf("read %s: %w", configPath, err)
		}
	}
	if server, ok := servers["vapi"]; ok {
		if strings.Contains(strings.Join(server.Args, " "), "@vapi-ai/mcp-docs-server") {
			cmd.Printf("%s: Legacy local docs server; run 'vapi mcp setup %s' to replace it with the v2 connection (%s)\n", ide, ide, configPath)
		} else if !containsMCPArg(server.Args, vapiMCPURL) {
			cmd.Printf("%s: vapi entry uses a different server; run 'vapi mcp setup %s' (%s)\n", ide, ide, configPath)
		} else {
			cmd.Printf("%s: Configured (%s)\n", ide, configPath)
		}
	} else {
		cmd.Printf("%s: Not configured (%s)\n", ide, configPath)
	}
	return nil
}

func containsMCPArg(args []string, value string) bool {
	for _, arg := range args {
		if arg == value {
			return true
		}
	}
	return false
}

func mcpConfigPath(ide string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch ide {
	case "cursor":
		return filepath.Join(homeDir, ".cursor", "mcp.json"), nil
	case "windsurf":
		return filepath.Join(homeDir, ".codeium", "windsurf", "mcp_config.json"), nil
	case "vscode":
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(configDir, "Code", "User", "mcp.json"), nil
	default:
		return "", fmt.Errorf("unsupported IDE %q; supported: cursor, windsurf, vscode", ide)
	}
}

func init() { rootCmd.AddCommand(mcpCmd) }
