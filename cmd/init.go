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
	"os/exec"
	"path/filepath"
	"runtime"
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
	case integrations.ProjectTypeUnknown:
		// Unknown project type
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
	sdkInstalled := false
	if integrationOptions.InstallSDK {
		installCmd := project.GetSDKInstallCommand()
		if installCmd != "" {
			fmt.Printf("📦 Installing %s...\n", project.GetSDKPackage())

			// Check if it's a build file modification (Java)
			if strings.Contains(installCmd, "//") {
				fmt.Printf("   %s\n", installCmd)
			} else {
				// Split the command into parts
				parts := strings.Fields(installCmd)
				if len(parts) > 0 {
					// Validate command parts for security
					if !isValidCommand(parts[0]) {
						fmt.Printf("   ⚠️  Skipping potentially unsafe command: %s\n", installCmd)
					} else {
						// #nosec G204 - command is validated above
						cmd := exec.Command(parts[0], parts[1:]...)
						cmd.Dir = absPath
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr

						if err := cmd.Run(); err != nil {
							fmt.Printf("   ⚠️  Installation failed: %v\n", err)
							fmt.Printf("   You may need to run this manually: %s\n", installCmd)
						} else {
							fmt.Printf("   ✅ Successfully installed %s\n", project.GetSDKPackage())
							sdkInstalled = true
						}
					}
				}
			}
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
		case integrations.FrameworkVue:
			fmt.Printf("   ℹ️  Code generation for %s coming soon!\n", project.GetFrameworkName())
		case integrations.FrameworkAngular:
			fmt.Printf("   ℹ️  Code generation for %s coming soon!\n", project.GetFrameworkName())
		case integrations.FrameworkSvelte:
			fmt.Printf("   ℹ️  Code generation for %s coming soon!\n", project.GetFrameworkName())
		case integrations.FrameworkNuxt:
			fmt.Printf("   ℹ️  Code generation for %s coming soon!\n", project.GetFrameworkName())
		case integrations.FrameworkVanilla:
			fmt.Printf("   ℹ️  Code generation for %s coming soon!\n", project.GetFrameworkName())
		case integrations.FrameworkReactNative:
			fmt.Printf("   ℹ️  Code generation for %s coming soon!\n", project.GetFrameworkName())
		case integrations.FrameworkFlutter:
			fmt.Printf("   ℹ️  Code generation for %s coming soon!\n", project.GetFrameworkName())
		case integrations.FrameworkRuby:
			fmt.Printf("   ℹ️  Code generation for %s coming soon!\n", project.GetFrameworkName())
		case integrations.FrameworkJava:
			fmt.Printf("   ℹ️  Code generation for %s coming soon!\n", project.GetFrameworkName())
		case integrations.FrameworkCSharp:
			fmt.Printf("   ℹ️  Code generation for %s coming soon!\n", project.GetFrameworkName())
		case integrations.FrameworkUnknown:
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

	// SDK installation step (only show if not already installed)
	if integrationOptions.InstallSDK && !sdkInstalled {
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
	if project.ProjectType != integrations.ProjectTypeUnknown {
		fmt.Printf("%d. Framework-specific setup:\n", stepNum)

		switch project.Framework {
		case integrations.FrameworkReact:
			fmt.Println("   - Add the Vapi component to your app")
			fmt.Println("   - Configure the voice interface in your main component")
		case integrations.FrameworkNext:
			fmt.Println("   - Add Vapi to your Next.js pages or components")
			fmt.Println("   - Consider using the useVapi hook for state management")
		case integrations.FrameworkNode:
			fmt.Println("   - Set up webhook endpoints for call events")
			fmt.Println("   - Create assistant configuration and call management routes")
		case integrations.FrameworkPython:
			fmt.Println("   - Configure webhook handlers for call events")
			fmt.Println("   - Set up assistant and call management endpoints")
		case integrations.FrameworkVue:
			fmt.Println("   - Follow the Vue.js integration guide in documentation")
		case integrations.FrameworkAngular:
			fmt.Println("   - Follow the Angular integration guide in documentation")
		case integrations.FrameworkSvelte:
			fmt.Println("   - Follow the Svelte integration guide in documentation")
		case integrations.FrameworkNuxt:
			fmt.Println("   - Follow the Nuxt.js integration guide in documentation")
		case integrations.FrameworkRemix:
			fmt.Println("   - Follow the Remix integration guide in documentation")
		case integrations.FrameworkVanilla:
			fmt.Println("   - Follow the vanilla JavaScript integration guide in documentation")
		case integrations.FrameworkReactNative:
			fmt.Println("   - Follow the React Native integration guide in documentation")
		case integrations.FrameworkFlutter:
			fmt.Println("   - Follow the Flutter integration guide in documentation")
		case integrations.FrameworkGolang:
			fmt.Println("   - Set up webhook endpoints for call events")
			fmt.Println("   - Create assistant configuration and call management routes")
		case integrations.FrameworkRuby:
			fmt.Println("   - Follow the Ruby integration guide in documentation")
		case integrations.FrameworkJava:
			fmt.Println("   - Follow the Java integration guide in documentation")
		case integrations.FrameworkCSharp:
			fmt.Println("   - Follow the C#/.NET integration guide in documentation")
		case integrations.FrameworkUnknown:
			fmt.Println("   - Follow the documentation for your specific framework")
		}
		fmt.Println()
		stepNum++
	}

	// MCP Integration recommendation (v0.1.9 feature)
	fmt.Printf("%d. 🧠 Set up IDE Integration (Recommended):\n", stepNum)
	fmt.Println("   vapi mcp setup")
	fmt.Println("   This turns your IDE into a Vapi expert with:")
	fmt.Println("   • Real-time access to 138+ documentation pages")
	fmt.Println("   • Smart code suggestions and examples")
	fmt.Println("   • No more hallucinated API information")
	fmt.Println("   • Works with Cursor, Windsurf, and VSCode")
	fmt.Println()
	stepNum++

	// Testing and deployment
	fmt.Printf("%d. Test and deploy:\n", stepNum)
	fmt.Println("   - Test your voice integration locally")
	fmt.Println("   - Set up webhooks for production")
	fmt.Println("   - Deploy with your API keys configured")
	fmt.Println()

	// Offer to open documentation
	var openDocs bool
	docsPrompt := &survey.Confirm{
		Message: "Would you like to open the Vapi documentation?",
		Default: true,
	}

	if survey.AskOne(docsPrompt, &openDocs) == nil && openDocs {
		fmt.Println("📖 Opening documentation...")

		// Use the correct documentation URL
		docURL := "https://docs.vapi.ai/quickstart/introduction"

		// Try to open the browser
		var openCmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			openCmd = exec.Command("open", docURL)
		case "linux":
			openCmd = exec.Command("xdg-open", docURL)
		case "windows":
			openCmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", docURL)
		default:
			fmt.Printf("   Visit: %s\n", docURL)
			return nil
		}

		if openCmd != nil {
			if err := openCmd.Start(); err != nil {
				fmt.Printf("   Visit: %s\n", docURL)
			}
		}
	}

	return nil
}

// isValidCommand checks if a command is safe to execute
func isValidCommand(cmd string) bool {
	// Allow common package managers and build tools
	allowedCommands := []string{
		"npm", "yarn", "pnpm", "bun", // Node.js package managers
		"pip", "pip3", "poetry", "conda", // Python package managers
		"mvn", "gradle", // Java build tools
		"go", "cargo", // Go and Rust
		"composer", // PHP
		"bundle",   // Ruby
		"dotnet",   // .NET
	}

	for _, allowed := range allowedCommands {
		if cmd == allowed {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(initCmd)
}
