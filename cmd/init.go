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
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/VapiAI/cli/pkg/integrations"
)

var initCmd = &cobra.Command{
	Use:   "init [project-path]",
	Short: "Initialize Vapi integration in your project",
	Long: `Initialize Vapi integration in your existing web project.
    
This interactive command will:
- Automatically detect your project framework
- Install the appropriate Vapi SDK
- Generate framework-specific components
- Create configuration templates
- Guide you through the setup process`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInitCommand,
}

func runInitCommand(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	// Make path absolute
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Check if directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", absPath)
	}

	fmt.Println("🚀 Welcome to Vapi Setup!")
	fmt.Println()

	// Step 1: Detect project
	fmt.Println("📋 Analyzing your project...")
	project, err := integrations.DetectProject(absPath)
	if err != nil {
		return fmt.Errorf("failed to detect project: %w", err)
	}

	// Display detected information
	fmt.Printf("✓ Detected: %s project", project.GetFrameworkName())
	if project.PackageJSON != nil && project.PackageJSON.Name != "" {
		fmt.Printf(" (%s)", project.PackageJSON.Name)
	}
	fmt.Println()

	// Display additional project info
	switch project.ProjectType {
	case integrations.ProjectTypeWeb:
		if project.IsTypeScript {
			fmt.Println("✓ TypeScript support detected")
		}
		if project.HasTailwind {
			fmt.Println("✓ Tailwind CSS detected")
		}
		if project.BuildTool != "" {
			fmt.Printf("✓ Build tool: %s\n", project.BuildTool)
		}
		if project.PackageJSON != nil {
			fmt.Printf("✓ Package manager: %s\n", project.GetPackageManager())
		}
	case integrations.ProjectTypeBackend:
		if project.GoVersion != "" {
			fmt.Printf("✓ Go version: %s\n", project.GoVersion)
		}
		if project.JavaBuildTool != "" {
			fmt.Printf("✓ Build tool: %s\n", project.JavaBuildTool)
		}
		if project.IsTypeScript {
			fmt.Println("✓ TypeScript support detected")
		}
	case integrations.ProjectTypeMobile:
		// Mobile specific info
	}
	fmt.Println()

	// Step 2: Confirm or override framework detection
	frameworkChoices := []string{
		"Use detected: " + project.GetFrameworkName(),
		"--- Frontend Frameworks ---",
		"React",
		"Vue",
		"Angular",
		"Svelte",
		"Next.js",
		"Nuxt",
		"Remix",
		"Vanilla JavaScript",
		"--- Mobile Frameworks ---",
		"React Native",
		"Flutter",
		"--- Backend Languages ---",
		"Node.js/TypeScript",
		"Python",
		"Go",
		"Ruby",
		"Java",
		"C#/.NET",
	}

	var frameworkChoice string
	prompt := &survey.Select{
		Message: "Confirm your framework/language:",
		Options: frameworkChoices,
		Default: frameworkChoices[0],
		Filter: func(filter string, value string, index int) bool {
			// Don't filter section headers
			if strings.Contains(value, "---") {
				return false
			}
			// Always show the default option
			if index == 0 {
				return true
			}
			// Filter based on user input
			return strings.Contains(strings.ToLower(value), strings.ToLower(filter))
		},
	}

	if err := survey.AskOne(prompt, &frameworkChoice); err != nil {
		return fmt.Errorf("framework selection canceled: %w", err)
	}

	// Update framework if user selected different one
	if frameworkChoice != frameworkChoices[0] && !strings.Contains(frameworkChoice, "---") {
		switch frameworkChoice {
		case "React":
			project.Framework = integrations.FrameworkReact
			project.ProjectType = integrations.ProjectTypeWeb
		case "Vue":
			project.Framework = integrations.FrameworkVue
			project.ProjectType = integrations.ProjectTypeWeb
		case "Angular":
			project.Framework = integrations.FrameworkAngular
			project.ProjectType = integrations.ProjectTypeWeb
		case "Svelte":
			project.Framework = integrations.FrameworkSvelte
			project.ProjectType = integrations.ProjectTypeWeb
		case "Next.js":
			project.Framework = integrations.FrameworkNext
			project.ProjectType = integrations.ProjectTypeWeb
		case "Nuxt":
			project.Framework = integrations.FrameworkNuxt
			project.ProjectType = integrations.ProjectTypeWeb
		case "Remix":
			project.Framework = integrations.FrameworkRemix
			project.ProjectType = integrations.ProjectTypeWeb
		case "Vanilla JavaScript":
			project.Framework = integrations.FrameworkVanilla
			project.ProjectType = integrations.ProjectTypeWeb
		case "React Native":
			project.Framework = integrations.FrameworkReactNative
			project.ProjectType = integrations.ProjectTypeMobile
		case "Flutter":
			project.Framework = integrations.FrameworkFlutter
			project.ProjectType = integrations.ProjectTypeMobile
		case "Node.js/TypeScript":
			project.Framework = integrations.FrameworkNode
			project.ProjectType = integrations.ProjectTypeBackend
		case "Python":
			project.Framework = integrations.FrameworkPython
			project.ProjectType = integrations.ProjectTypeBackend
		case "Go":
			project.Framework = integrations.FrameworkGolang
			project.ProjectType = integrations.ProjectTypeBackend
		case "Ruby":
			project.Framework = integrations.FrameworkRuby
			project.ProjectType = integrations.ProjectTypeBackend
		case "Java":
			project.Framework = integrations.FrameworkJava
			project.ProjectType = integrations.ProjectTypeBackend
		case "C#/.NET":
			project.Framework = integrations.FrameworkCSharp
			project.ProjectType = integrations.ProjectTypeBackend
		}
	}

	// Step 3: Integration options (different based on project type)
	var integrationOptions struct {
		InstallSDK         bool
		GenerateExamples   bool
		SetupEnvironment   bool
		GenerateComponents bool // Only for frontend
	}

	questions := []*survey.Question{
		{
			Name: "InstallSDK",
			Prompt: &survey.Confirm{
				Message: fmt.Sprintf("Install %s SDK?", project.GetSDKPackage()),
				Default: true,
			},
		},
		{
			Name: "GenerateExamples",
			Prompt: &survey.Confirm{
				Message: "Generate example code?",
				Default: true,
			},
		},
		{
			Name: "SetupEnvironment",
			Prompt: &survey.Confirm{
				Message: "Create environment configuration template?",
				Default: true,
			},
		},
	}

	// Add components question only for frontend projects
	if project.ProjectType == integrations.ProjectTypeWeb {
		questions = append(questions[:2], append([]*survey.Question{{
			Name: "GenerateComponents",
			Prompt: &survey.Confirm{
				Message: "Generate Vapi components and hooks?",
				Default: true,
			},
		}}, questions[2:]...)...)
	}

	if err := survey.Ask(questions, &integrationOptions); err != nil {
		return fmt.Errorf("setup canceled: %w", err)
	}

	fmt.Println()
	fmt.Println("🔧 Setting up Vapi integration...")
	fmt.Println()

	// Execute selected options
	if integrationOptions.InstallSDK {
		installCmd := project.GetSDKInstallCommand()
		if installCmd != "" {
			fmt.Printf("📦 Installing %s...\n", project.GetSDKPackage())
			fmt.Printf("   Run: %s\n", installCmd)
		}
	}

	if integrationOptions.GenerateExamples || integrationOptions.GenerateComponents || integrationOptions.SetupEnvironment {
		fmt.Println("🎨 Generating Vapi integration files...")

		// Call appropriate generation function based on framework
		var genErr error
		switch project.Framework {
		case integrations.FrameworkReact, integrations.FrameworkNext, integrations.FrameworkRemix:
			if integrationOptions.GenerateComponents {
				genErr = integrations.GenerateReactIntegration(absPath, project)
			}
		case integrations.FrameworkPython:
			genErr = integrations.GeneratePythonIntegration(absPath, project)
		case integrations.FrameworkGolang:
			genErr = integrations.GenerateGoIntegration(absPath, project)
		case integrations.FrameworkNode:
			genErr = integrations.GenerateNodeIntegration(absPath, project)
		// TODO: Add other framework integrations
		default:
			fmt.Printf("   ℹ️  Code generation for %s coming soon!\n", project.GetFrameworkName())
		}

		if genErr != nil {
			fmt.Printf("   ⚠️  Error generating files: %v\n", genErr)
		} else if genErr == nil && project.Framework != integrations.FrameworkUnknown {
			fmt.Println("   ✓ Generated example files")
		}
	}

	// Step 4: Final instructions
	fmt.Println()
	fmt.Println("🎉 Vapi integration setup complete!")
	fmt.Println()
	fmt.Println("📝 Next steps:")
	fmt.Println()

	stepNum := 1

	// SDK installation step
	if integrationOptions.InstallSDK {
		installCmd := project.GetSDKInstallCommand()
		if strings.Contains(installCmd, "//") {
			// For Java/build file modifications
			fmt.Printf("%d. Add Vapi SDK to your build file:\n", stepNum)
			fmt.Printf("   %s\n\n", installCmd)
		} else {
			fmt.Printf("%d. Install dependencies:\n", stepNum)
			fmt.Printf("   %s\n\n", installCmd)
		}
		stepNum++
	}

	// Environment setup step
	if integrationOptions.SetupEnvironment {
		fmt.Printf("%d. Configure your environment:\n", stepNum)
		fmt.Println("   - Copy .env.example to .env")
		fmt.Println("   - Add your Vapi API key")
		if project.ProjectType == integrations.ProjectTypeWeb || project.ProjectType == integrations.ProjectTypeMobile {
			fmt.Println("   - Add your Vapi public key and assistant ID")
		}
		fmt.Println()
		stepNum++
	}

	// Framework-specific instructions
	switch project.ProjectType {
	case integrations.ProjectTypeWeb:
		fmt.Printf("%d. Import and use Vapi in your application\n", stepNum)
		fmt.Println("   - Check the generated components in your project")
		fmt.Println("   - Use the VapiButton component or useVapi hook")
	case integrations.ProjectTypeBackend:
		fmt.Printf("%d. Run the example code:\n", stepNum)
		switch project.Framework {
		case integrations.FrameworkPython:
			fmt.Println("   python vapi_examples/basic_example.py")
		case integrations.FrameworkGolang:
			fmt.Println("   go run examples/vapi/basic_example.go")
		case integrations.FrameworkNode:
			if project.IsTypeScript {
				fmt.Println("   npx tsx vapi-examples/basic-example.ts")
			} else {
				fmt.Println("   node vapi-examples/basic-example.js")
			}
		}
	case integrations.ProjectTypeMobile:
		fmt.Printf("%d. Check the generated example code\n", stepNum)
	}
	fmt.Println()

	// Offer to open documentation
	var openDocs bool
	docsPrompt := &survey.Confirm{
		Message: "Would you like to open the Vapi documentation?",
		Default: true,
	}

	if survey.AskOne(docsPrompt, &openDocs) == nil && openDocs {
		fmt.Println("📖 Opening documentation...")
		// TODO: Open browser to framework-specific Vapi docs
		fmt.Printf("   Visit: https://docs.vapi.ai/sdk/%s\n", strings.ToLower(project.GetFrameworkName()))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
