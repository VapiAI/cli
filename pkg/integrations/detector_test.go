package integrations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectProject(t *testing.T) {
	tests := []struct {
		name        string
		setupFiles  map[string]string
		expected    Framework
		projectType ProjectType
	}{
		{
			name: "detects Go project",
			setupFiles: map[string]string{
				"go.mod": `module example.com/test
go 1.21`,
			},
			expected:    FrameworkGolang,
			projectType: ProjectTypeBackend,
		},
		{
			name: "detects Python project",
			setupFiles: map[string]string{
				"requirements.txt": "flask==2.0.0\nrequests==2.31.0",
			},
			expected:    FrameworkPython,
			projectType: ProjectTypeBackend,
		},
		{
			name: "detects React project",
			setupFiles: map[string]string{
				"package.json": `{
					"name": "test-app",
					"dependencies": {
						"react": "^18.0.0",
						"react-dom": "^18.0.0"
					}
				}`,
			},
			expected:    FrameworkReact,
			projectType: ProjectTypeWeb,
		},
		{
			name: "detects Next.js project",
			setupFiles: map[string]string{
				"package.json": `{
					"name": "test-app",
					"dependencies": {
						"next": "^13.0.0",
						"react": "^18.0.0"
					}
				}`,
			},
			expected:    FrameworkNext,
			projectType: ProjectTypeWeb,
		},
		{
			name: "detects Node.js backend project",
			setupFiles: map[string]string{
				"package.json": `{
					"name": "api-server",
					"dependencies": {
						"express": "^4.0.0"
					}
				}`,
			},
			expected:    FrameworkNode,
			projectType: ProjectTypeBackend,
		},
		{
			name: "prioritizes Go over Python",
			setupFiles: map[string]string{
				"go.mod":           "module example.com/test\ngo 1.21",
				"requirements.txt": "flask==2.0.0",
			},
			expected:    FrameworkGolang,
			projectType: ProjectTypeBackend,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			defer func() {
				if err := os.RemoveAll(tmpDir); err != nil {
					t.Logf("Warning: failed to remove temp dir: %v", err)
				}
			}()

			// Create test files
			for filename, content := range tt.setupFiles {
				filePath := filepath.Join(tmpDir, filename)
				err := os.WriteFile(filePath, []byte(content), 0o644)
				require.NoError(t, err)
			}

			// Run detection
			project, err := DetectProject(tmpDir)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, project.Framework)
			assert.Equal(t, tt.projectType, project.ProjectType)
		})
	}
}

func TestDetectProjectFeatures(t *testing.T) {
	tmpDir := t.TempDir()
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	// Create a React project with TypeScript and Tailwind
	packageJSON := PackageJSON{
		Name: "test-app",
		Dependencies: map[string]string{
			"react": "^18.0.0",
		},
		DevDependencies: map[string]string{
			"tailwindcss": "^3.0.0",
			"typescript":  "^5.0.0",
		},
	}

	data, err := json.Marshal(packageJSON)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmpDir, "package.json"), data, 0o644)
	require.NoError(t, err)

	// Create tsconfig.json
	err = os.WriteFile(filepath.Join(tmpDir, "tsconfig.json"), []byte("{}"), 0o644)
	require.NoError(t, err)

	// Create tailwind.config.js
	err = os.WriteFile(filepath.Join(tmpDir, "tailwind.config.js"), []byte("module.exports = {}"), 0o644)
	require.NoError(t, err)

	// Run detection
	project, err := DetectProject(tmpDir)
	require.NoError(t, err)

	assert.True(t, project.IsTypeScript)
	assert.True(t, project.HasTailwind)
	assert.Equal(t, FrameworkReact, project.Framework)
}

func TestGetSDKPackage(t *testing.T) {
	tests := []struct {
		framework Framework
		expected  string
	}{
		{FrameworkReact, "@vapi-ai/web"},
		{FrameworkNext, "@vapi-ai/web"},
		{FrameworkReactNative, "@vapi-ai/react-native"},
		{FrameworkNode, "@vapi-ai/server-sdk"},
		{FrameworkPython, "vapi-python"},
		{FrameworkGolang, "github.com/VapiAI/vapi-go"},
		{FrameworkRuby, "vapi"},
		{FrameworkJava, "ai.vapi:vapi-java"},
		{FrameworkCSharp, "Vapi"},
		{FrameworkFlutter, "vapi_flutter"},
	}

	for _, tt := range tests {
		t.Run(string(tt.framework), func(t *testing.T) {
			project := &ProjectInfo{Framework: tt.framework}
			assert.Equal(t, tt.expected, project.GetSDKPackage())
		})
	}
}
