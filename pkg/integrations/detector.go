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
package integrations

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Framework string

const (
	// Frontend Frameworks
	FrameworkReact   Framework = "react"
	FrameworkVue     Framework = "vue"
	FrameworkAngular Framework = "angular"
	FrameworkSvelte  Framework = "svelte"
	FrameworkNext    Framework = "nextjs"
	FrameworkNuxt    Framework = "nuxtjs"
	FrameworkRemix   Framework = "remix"
	FrameworkVanilla Framework = "vanilla"

	// Mobile Frameworks
	FrameworkReactNative Framework = "react-native"
	FrameworkFlutter     Framework = "flutter"

	// Backend Languages/Frameworks
	FrameworkPython Framework = "python"
	FrameworkGolang Framework = "go"
	FrameworkRuby   Framework = "ruby"
	FrameworkJava   Framework = "java"
	FrameworkCSharp Framework = "csharp"
	FrameworkNode   Framework = "node"

	FrameworkUnknown Framework = "unknown"
)

type ProjectType string

const (
	ProjectTypeWeb     ProjectType = "web"
	ProjectTypeMobile  ProjectType = "mobile"
	ProjectTypeBackend ProjectType = "backend"
	ProjectTypeUnknown ProjectType = "unknown"
)

type ProjectInfo struct {
	Path         string
	Framework    Framework
	ProjectType  ProjectType
	PackageJSON  *PackageJSON
	IsTypeScript bool
	HasTailwind  bool
	UsesPnpm     bool
	UsesYarn     bool
	UsesBun      bool
	BuildTool    string // vite, webpack, etc.

	// Backend specific fields
	PythonVersion string
	GoVersion     string
	JavaBuildTool string // maven, gradle
	DotNetVersion string
}

// DetectProject analyzes the current directory to detect project type
func DetectProject(projectPath string) (*ProjectInfo, error) {
	project := &ProjectInfo{
		Path:        projectPath,
		Framework:   FrameworkUnknown,
		ProjectType: ProjectTypeUnknown,
	}

	// Check for backend/mobile project markers in order of specificity
	// More specific markers first to avoid false positives

	// Go - go.mod is very specific to Go projects
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); err == nil {
		project.Framework = FrameworkGolang
		project.ProjectType = ProjectTypeBackend
		// Try to read Go version from go.mod
		if data, err := os.ReadFile(filepath.Join(projectPath, "go.mod")); err == nil {
			// Simple version extraction (could be improved)
			content := string(data)
			if content != "" {
				lines := strings.Split(content, "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "go ") {
						project.GoVersion = strings.TrimSpace(strings.TrimPrefix(line, "go"))
						break
					}
				}
			}
		}
		return project, nil
	}

	// Java - pom.xml and build.gradle are specific to Java
	if _, err := os.Stat(filepath.Join(projectPath, "pom.xml")); err == nil {
		project.Framework = FrameworkJava
		project.ProjectType = ProjectTypeBackend
		project.JavaBuildTool = "maven"
		return project, nil
	}
	if _, err := os.Stat(filepath.Join(projectPath, "build.gradle")); err == nil ||
		(func() bool {
			_, err := os.Stat(filepath.Join(projectPath, "build.gradle.kts"))
			return err == nil
		})() {
		project.Framework = FrameworkJava
		project.ProjectType = ProjectTypeBackend
		project.JavaBuildTool = "gradle"
		return project, nil
	}

	// Ruby - Gemfile is specific to Ruby
	if _, err := os.Stat(filepath.Join(projectPath, "Gemfile")); err == nil {
		project.Framework = FrameworkRuby
		project.ProjectType = ProjectTypeBackend
		return project, nil
	}

	// C#/.NET - .csproj and .sln are specific to .NET
	files, _ := os.ReadDir(projectPath)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".csproj") || strings.HasSuffix(file.Name(), ".sln") {
			project.Framework = FrameworkCSharp
			project.ProjectType = ProjectTypeBackend
			return project, nil
		}
	}

	// Flutter - pubspec.yaml is specific to Flutter/Dart
	if _, err := os.Stat(filepath.Join(projectPath, "pubspec.yaml")); err == nil {
		project.Framework = FrameworkFlutter
		project.ProjectType = ProjectTypeMobile
		return project, nil
	}

	// Python - Check after more specific markers since requirements.txt is common
	if _, err := os.Stat(filepath.Join(projectPath, "requirements.txt")); err == nil ||
		(func() bool {
			_, err1 := os.Stat(filepath.Join(projectPath, "setup.py"))
			_, err2 := os.Stat(filepath.Join(projectPath, "pyproject.toml"))
			_, err3 := os.Stat(filepath.Join(projectPath, "Pipfile"))
			return err1 == nil || err2 == nil || err3 == nil
		})() {
		project.Framework = FrameworkPython
		project.ProjectType = ProjectTypeBackend
		return project, nil
	}

	// If no backend/mobile markers found, check for Node.js/frontend projects
	packageJSONPath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		// Read and parse package.json
		data, err := os.ReadFile(packageJSONPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read package.json: %w", err)
		}

		var pkg PackageJSON
		if err := json.Unmarshal(data, &pkg); err != nil {
			return nil, fmt.Errorf("failed to parse package.json: %w", err)
		}

		project.PackageJSON = &pkg

		// Detect framework from package.json
		project.detectFramework()

		// Set project type based on detected framework
		switch project.Framework {
		case FrameworkReactNative:
			project.ProjectType = ProjectTypeMobile
		case FrameworkReact, FrameworkVue, FrameworkAngular, FrameworkSvelte,
			FrameworkNext, FrameworkNuxt, FrameworkRemix, FrameworkVanilla:
			project.ProjectType = ProjectTypeWeb
		default:
			// If no specific frontend framework detected, it's likely a Node.js backend
			project.Framework = FrameworkNode
			project.ProjectType = ProjectTypeBackend
		}

		// Detect other characteristics
		project.detectProjectFeatures()

		// Detect package manager
		project.detectPackageManager()

		return project, nil
	}

	// Check for vanilla HTML project
	indexPath := filepath.Join(projectPath, "index.html")
	if _, err := os.Stat(indexPath); err == nil {
		project.Framework = FrameworkVanilla
		project.ProjectType = ProjectTypeWeb
		return project, nil
	}

	return nil, fmt.Errorf("could not detect project type in %s", projectPath)
}

func (p *ProjectInfo) detectFramework() {
	deps := make(map[string]bool)

	// Combine dependencies and devDependencies
	if p.PackageJSON.Dependencies != nil {
		for dep := range p.PackageJSON.Dependencies {
			deps[dep] = true
		}
	}
	if p.PackageJSON.DevDependencies != nil {
		for dep := range p.PackageJSON.DevDependencies {
			deps[dep] = true
		}
	}

	// Check for frameworks in order of specificity
	switch {
	case deps["react-native"]:
		p.Framework = FrameworkReactNative
	case deps["next"]:
		p.Framework = FrameworkNext
	case deps["nuxt"] || deps["nuxt3"]:
		p.Framework = FrameworkNuxt
	case deps["@remix-run/react"]:
		p.Framework = FrameworkRemix
	case deps["react"]:
		p.Framework = FrameworkReact
	case deps["vue"]:
		p.Framework = FrameworkVue
	case deps["@angular/core"]:
		p.Framework = FrameworkAngular
	case deps["svelte"]:
		p.Framework = FrameworkSvelte
	default:
		// Check for vanilla JS project indicators
		if deps["vite"] || deps["webpack"] || deps["@rspack/cli"] || deps["@rspack/core"] || deps["parcel"] {
			p.Framework = FrameworkVanilla
		}
	}
}

func (p *ProjectInfo) detectProjectFeatures() {
	// Check for TypeScript
	tsConfigPath := filepath.Join(p.Path, "tsconfig.json")
	if _, err := os.Stat(tsConfigPath); err == nil {
		p.IsTypeScript = true
	}

	// Check for Tailwind CSS
	tailwindPaths := []string{
		filepath.Join(p.Path, "tailwind.config.js"),
		filepath.Join(p.Path, "tailwind.config.ts"),
	}
	for _, path := range tailwindPaths {
		if _, err := os.Stat(path); err == nil {
			p.HasTailwind = true
			break
		}
	}

	// Also check in dependencies
	if p.PackageJSON.Dependencies != nil {
		if _, hasTailwind := p.PackageJSON.Dependencies["tailwindcss"]; hasTailwind {
			p.HasTailwind = true
		}
	}
	if p.PackageJSON.DevDependencies != nil {
		if _, hasTailwind := p.PackageJSON.DevDependencies["tailwindcss"]; hasTailwind {
			p.HasTailwind = true
		}
	}

	// Detect build tool
	if p.PackageJSON.DevDependencies != nil {
		switch {
		case p.PackageJSON.DevDependencies["vite"] != "":
			p.BuildTool = "vite"
		case p.PackageJSON.DevDependencies["webpack"] != "":
			p.BuildTool = "webpack"
		case p.PackageJSON.DevDependencies["@rspack/cli"] != "" || p.PackageJSON.DevDependencies["@rspack/core"] != "":
			p.BuildTool = "rspack"
		case p.PackageJSON.DevDependencies["parcel"] != "":
			p.BuildTool = "parcel"
		case p.PackageJSON.DevDependencies["esbuild"] != "":
			p.BuildTool = "esbuild"
		}
	}

	// Also check for config files if build tool not detected from dependencies
	if p.BuildTool == "" {
		configFiles := map[string]string{
			"vite.config.js":    "vite",
			"vite.config.ts":    "vite",
			"webpack.config.js": "webpack",
			"webpack.config.ts": "webpack",
			"rspack.config.js":  "rspack",
			"rspack.config.ts":  "rspack",
		}

		for configFile, tool := range configFiles {
			if _, err := os.Stat(filepath.Join(p.Path, configFile)); err == nil {
				p.BuildTool = tool
				break
			}
		}
	}
}

func (p *ProjectInfo) detectPackageManager() {
	// Check for lock files in order of preference
	lockFiles := map[string]string{
		"bun.lock":          "bun", // New text-based lockfile
		"bun.lockb":         "bun", // Old binary lockfile
		"pnpm-lock.yaml":    "pnpm",
		"yarn.lock":         "yarn",
		"package-lock.json": "npm",
	}

	for lockFile, manager := range lockFiles {
		if _, err := os.Stat(filepath.Join(p.Path, lockFile)); err == nil {
			p.UsesPnpm = manager == "pnpm"
			p.UsesYarn = manager == "yarn"
			p.UsesBun = manager == "bun"
			return
		}
	}

	// Default to npm if no lock file found
}

func (p *ProjectInfo) GetPackageManager() string {
	switch {
	case p.UsesPnpm:
		return "pnpm"
	case p.UsesYarn:
		return "yarn"
	case p.UsesBun:
		return "bun"
	default:
		return "npm"
	}
}

func (p *ProjectInfo) GetInstallCommand() string {
	pm := p.GetPackageManager()
	switch pm {
	case "yarn":
		return "yarn"
	case "pnpm":
		return "pnpm install"
	case "bun":
		return "bun install"
	default:
		return "npm install"
	}
}

func (p *ProjectInfo) GetAddCommand(pkg string) string {
	pm := p.GetPackageManager()
	switch pm {
	case "yarn":
		return fmt.Sprintf("yarn add %s", pkg)
	case "pnpm":
		return fmt.Sprintf("pnpm add %s", pkg)
	case "bun":
		return fmt.Sprintf("bun add %s", pkg)
	default:
		return fmt.Sprintf("npm install %s", pkg)
	}
}

func (p *ProjectInfo) GetFrameworkName() string {
	switch p.Framework {
	case FrameworkReact:
		return "React"
	case FrameworkVue:
		return "Vue"
	case FrameworkAngular:
		return "Angular"
	case FrameworkSvelte:
		return "Svelte"
	case FrameworkNext:
		return "Next.js"
	case FrameworkNuxt:
		return "Nuxt"
	case FrameworkRemix:
		return "Remix"
	case FrameworkVanilla:
		return "Vanilla JavaScript"
	case FrameworkReactNative:
		return "React Native"
	case FrameworkFlutter:
		return "Flutter"
	case FrameworkPython:
		return "Python"
	case FrameworkGolang:
		return "Go"
	case FrameworkRuby:
		return "Ruby"
	case FrameworkJava:
		return "Java"
	case FrameworkCSharp:
		return "C#/.NET"
	case FrameworkNode:
		return "Node.js"
	default:
		return "Unknown"
	}
}

// GetSDKPackage returns the appropriate Vapi SDK package for the framework
func (p *ProjectInfo) GetSDKPackage() string {
	switch p.Framework {
	case FrameworkReact, FrameworkVue, FrameworkAngular, FrameworkSvelte,
		FrameworkNext, FrameworkNuxt, FrameworkRemix, FrameworkVanilla:
		return "@vapi-ai/web"
	case FrameworkReactNative:
		return "@vapi-ai/react-native"
	case FrameworkNode:
		return "@vapi-ai/server-sdk"
	case FrameworkPython:
		return "vapi-python"
	case FrameworkGolang:
		return "github.com/VapiAI/vapi-go"
	case FrameworkRuby:
		return "vapi"
	case FrameworkJava:
		return "ai.vapi:vapi-java"
	case FrameworkCSharp:
		return "Vapi"
	case FrameworkFlutter:
		return "vapi_flutter"
	default:
		return ""
	}
}

// GetInstallCommand returns the install command for the SDK
func (p *ProjectInfo) GetSDKInstallCommand() string {
	sdk := p.GetSDKPackage()
	if sdk == "" {
		return ""
	}

	switch p.Framework {
	case FrameworkPython:
		return fmt.Sprintf("pip install %s", sdk)
	case FrameworkGolang:
		return fmt.Sprintf("go get %s", sdk)
	case FrameworkRuby:
		return fmt.Sprintf("gem install %s", sdk)
	case FrameworkJava:
		if p.JavaBuildTool == "gradle" {
			return fmt.Sprintf("// Add to build.gradle:\nimplementation '%s:latest.release'", sdk)
		}
		return "// Add to pom.xml:\n<dependency>\n  <groupId>ai.vapi</groupId>\n  <artifactId>vapi-java</artifactId>\n  <version>LATEST</version>\n</dependency>"
	case FrameworkCSharp:
		return fmt.Sprintf("dotnet add package %s", sdk)
	case FrameworkFlutter:
		return fmt.Sprintf("flutter pub add %s", sdk)
	default:
		// Node.js based frameworks
		return p.GetAddCommand(sdk)
	}
}
