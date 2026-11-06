package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/VapiAI/cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func readMCPTestConfig(t *testing.T, file string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func TestMCPSetupMigratesLegacyDocsAndPreservesSettings(t *testing.T) {
	file := filepath.Join(t.TempDir(), "mcp.json")
	original := `{"inputs":[{"id":"secret","type":"promptString","password":true}],"mcpServers":{"other":{"url":"https://example.test/mcp","headers":{"Authorization":"keep"},"enabled":false},"vapi":{"command":"npx","args":["-y","@vapi-ai/mcp-docs-server"]}}}`
	if err := os.WriteFile(file, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	before := readMCPTestConfig(t, file)
	for i := 0; i < 2; i++ {
		if err := writeMCPConfig(file, "mcpServers", remoteMCPConfig("test-key", "darwin")); err != nil {
			t.Fatal(err)
		}
	}

	after := readMCPTestConfig(t, file)
	beforeOther, _ := json.Marshal(before["mcpServers"].(map[string]any)["other"])
	afterOther, _ := json.Marshal(after["mcpServers"].(map[string]any)["other"])
	if string(beforeOther) != string(afterOther) {
		t.Fatal("unrelated server changed")
	}
	inputs, _ := json.Marshal(after["inputs"])
	if string(inputs) != `[{"id":"secret","password":true,"type":"promptString"}]` {
		t.Fatalf("inputs changed: %s", inputs)
	}
	servers := after["mcpServers"].(map[string]any)
	if len(servers) != 2 {
		t.Fatalf("unexpected server entries: %v", servers)
	}
	apiJSON, _ := json.Marshal(servers["vapi"])
	if !strings.Contains(string(apiJSON), vapiMCPURL) || !strings.Contains(string(apiJSON), "test-key") {
		t.Fatalf("missing API endpoint or auth: %s", apiJSON)
	}

	info, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
		t.Fatalf("key file mode: %v", info.Mode())
	}
}

func TestMCPSetupRefusesInvalidConfigs(t *testing.T) {
	for _, original := range []string{`{broken`, `null`, `[]`, `{"servers":null}`, `{"servers":[]}`} {
		t.Run(original, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "mcp.json")
			if err := os.WriteFile(file, []byte(original), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := writeMCPConfig(file, "servers", remoteMCPConfig("key", "linux")); err == nil {
				t.Fatal("expected invalid config error")
			}
			data, _ := os.ReadFile(file)
			if string(data) != original {
				t.Fatal("invalid config was overwritten")
			}
		})
	}
}

func TestMCPGeneratedTransports(t *testing.T) {
	for _, goos := range []string{"darwin", "linux", "windows"} {
		server := remoteMCPConfig("key", goos)
		if !containsMCPArg(server.Args, "http-only") {
			t.Fatal("legacy fallback allowed")
		}
		if !containsMCPArg(server.Args, "Authorization: Bearer ${VAPI_API_KEY}") || server.Env["VAPI_API_KEY"] != "key" {
			t.Fatal("missing bearer environment")
		}
		if goos == "windows" && (server.Command != "cmd" || server.Args[0] != "/c" || server.Args[1] != "npx") {
			t.Fatal("missing Windows launcher")
		}
	}
}

func TestMCPStatusDoesNotReportLegacyDocsAsAPI(t *testing.T) {
	userDir := t.TempDir()
	t.Setenv("HOME", userDir)
	t.Setenv("USERPROFILE", userDir)
	file, err := mcpConfigPath("cursor")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeMCPConfig(file, "mcpServers", MCPServerConfig{Command: "npx", Args: []string{"@vapi-ai/mcp-docs-server"}}); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	command := &cobra.Command{}
	command.SetOut(&output)
	if err := printMCPStatus(command, "cursor"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Legacy local docs server") || strings.Contains(output.String(), "Configured (") {
		t.Fatalf("misleading status: %s", output.String())
	}
}

func TestPublicSetupDoesNotRequireVapiAuthentication(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	if err := rootCmd.PersistentPreRunE(skillsCmd, nil); err != nil {
		t.Fatalf("skills setup required auth: %v", err)
	}
}

func TestAPISetupRequiresKeyBeforeWriting(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv("VAPI_API_KEY", "")
	previous := config.GetConfig()
	config.SetConfig(nil)
	t.Cleanup(func() { config.SetConfig(previous) })
	setup, _, err := mcpCmd.Find([]string{"setup"})
	if err != nil {
		t.Fatal(err)
	}
	if err := setup.RunE(setup, []string{"cursor"}); err == nil || !strings.Contains(err.Error(), "API key is required") {
		t.Fatalf("missing auth error: %v", err)
	}
}

func TestMCPAPIKeyUsesActiveAccountAndExplicitOverride(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv("VAPI_API_KEY", "")
	previous := config.GetConfig()
	t.Cleanup(func() { config.SetConfig(previous) })
	config.SetConfig(&config.Config{APIKey: "legacy", ActiveAccount: "current", Accounts: map[string]config.Account{"current": {APIKey: "active"}}})
	viper.Set("api_key", "legacy")
	command := &cobra.Command{}
	command.Flags().String("api-key", "", "")
	if got := mcpAPIKey(command); got != "active" {
		t.Fatalf("used %q instead of active account", got)
	}
	if err := command.Flags().Set("api-key", "override"); err != nil {
		t.Fatal(err)
	}
	if got := mcpAPIKey(command); got != "override" {
		t.Fatalf("explicit override ignored: %q", got)
	}
}

func TestMCPSetupWritesV2IntegrationForEachIDE(t *testing.T) {
	userDir := t.TempDir()
	t.Setenv("HOME", userDir)
	t.Setenv("USERPROFILE", userDir)
	t.Setenv("APPDATA", userDir)
	t.Setenv("XDG_CONFIG_HOME", userDir)
	command := &cobra.Command{}
	var output bytes.Buffer
	command.SetOut(&output)
	for _, ide := range []string{"cursor", "windsurf", "vscode"} {
		if err := setupMCPIntegration(command, []string{ide}, "test-only-key"); err != nil {
			t.Fatal(err)
		}

		file, err := mcpConfigPath(ide)
		if err != nil {
			t.Fatal(err)
		}
		section := "mcpServers"
		if ide == "vscode" {
			section = "servers"
		}
		servers := readMCPTestConfig(t, file)[section].(map[string]any)
		if len(servers) != 1 {
			t.Fatalf("%s: expected one v2 entry", ide)
		}
		if ide == "vscode" && servers["vapi"].(map[string]any)["type"] != "stdio" {
			t.Fatal("missing VS Code transport type")
		}
	}
	if strings.Contains(output.String(), "test-only-key") {
		t.Fatal("API key leaked to setup output")
	}
}
