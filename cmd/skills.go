package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/VapiAI/cli/pkg/analytics"
)

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "Install public Vapi agent skills",
}

var skillsSetupCmd = &cobra.Command{
	Use:   "setup [skills installer options]",
	Short: "Install public skills from VapiAI/skills",
	Long: `Install public Vapi agent skills using 'npx -y skills add VapiAI/skills'.

Requires Node.js 22.20 or later with npm/npx. No Vapi API key is needed to install.
The installer lets you choose skills and agents; installation is project-local by default.
Options are forwarded to the skills installer, including --global, --agent, --skill, and --yes.
This installs skill files. It does not configure an MCP connection.`,
	Example: `  vapi skills setup
  vapi skills setup --agent cursor --skill create-assistant
  vapi skills setup --global --agent claude-code --yes`,
	DisableFlagParsing: true,
	RunE:               analytics.TrackCommandWrapper("skills", "setup", runSkillsSetup),
}

func runSkillsSetup(cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return cmd.Help()
		}
	}
	installer := skillsInstallCommand(cmd, args)
	installer.Stdin = cmd.InOrStdin()
	installer.Stdout = cmd.OutOrStdout()
	installer.Stderr = cmd.ErrOrStderr()
	if err := installer.Run(); err != nil {
		return fmt.Errorf("skills installation failed (requires Node.js 22.20+ and npx): %w", err)
	}
	return nil
}

func skillsInstallCommand(cmd *cobra.Command, args []string) *exec.Cmd {
	installerArgs := append([]string{"-y", "skills", "add", "VapiAI/skills"}, args...)
	// Use the Windows command shim for npx.cmd; launch npx directly on other platforms.
	if runtime.GOOS == "windows" {
		return exec.CommandContext(cmd.Context(), "cmd", append([]string{"/c", "npx"}, installerArgs...)...)
	}
	return exec.CommandContext(cmd.Context(), "npx", installerArgs...) // #nosec G204 -- executable/source are fixed; options are explicitly forwarded
}

func init() { skillsCmd.AddCommand(skillsSetupCmd); rootCmd.AddCommand(skillsCmd) }
