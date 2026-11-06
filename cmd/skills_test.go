package cmd

import (
	"context"
	"reflect"
	"runtime"
	"testing"

	"github.com/spf13/cobra"
)

func TestSkillsSetupForwardsInstallerOptions(t *testing.T) {
	command := &cobra.Command{}
	command.SetContext(context.Background())
	args := []string{"--global", "--agent", "cursor", "--skill", "create-assistant", "--yes"}
	installer := skillsInstallCommand(command, args)
	want := append([]string{"npx", "-y", "skills", "add", "VapiAI/skills"}, args...)
	if runtime.GOOS == "windows" {
		want = append([]string{"cmd", "/c"}, want...)
	}
	if !reflect.DeepEqual(installer.Args, want) {
		t.Fatalf("got %v; want %v", installer.Args, want)
	}
	if len(args) != 6 || args[0] != "--global" {
		t.Fatal("input options mutated")
	}
}
