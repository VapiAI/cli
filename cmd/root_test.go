package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name: "shows help when no args",
			args: []string{},
			contains: []string{
				"Available Commands:",
				"assistant",
				"call",
				"config",
				"init",
				"login",
			},
		},
		{
			name: "shows help with --help flag",
			args: []string{"--help"},
			contains: []string{
				"Available Commands:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := rootCmd
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			assert.NoError(t, err)

			output := buf.String()
			for _, expected := range tt.contains {
				assert.Contains(t, output, expected)
			}
		})
	}
}

func TestAuthValidation(t *testing.T) {
	// Test that commands requiring auth fail without API key
	authRequiredCommands := [][]string{
		{"assistant", "list"},
		{"call", "list"},
	}

	for _, args := range authRequiredCommands {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			// Clear viper config to ensure no API key
			viper.Reset()
			viper.Set("api_key", "")

			cmd := rootCmd
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(args)

			err := cmd.Execute()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not authenticated")
		})
	}
}

func TestNoAuthValidation(t *testing.T) {
	// Test that these commands work without auth
	noAuthCommands := [][]string{
		{"login"},
		{"config", "get"},
		{"init"},
		{"--help"},
	}

	for _, args := range noAuthCommands {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			// Skip interactive commands in tests
			if args[0] == "init" || args[0] == "login" {
				t.Skipf("Skipping %s command test (interactive)", args[0])
			}

			// Clear viper config to ensure no API key
			viper.Reset()
			viper.Set("api_key", "")

			cmd := rootCmd
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(args)

			// We don't check for NoError here because the command might fail
			// for other reasons, we just want to ensure it doesn't fail with
			// "not authenticated"
			err := cmd.Execute()
			if err != nil {
				assert.NotContains(t, err.Error(), "not authenticated")
			}
		})
	}
}
