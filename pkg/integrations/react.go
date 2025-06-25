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
)

type ReactProject struct {
	Path         string
	PackageJSON  *PackageJSON
	IsTypeScript bool
	IsNextJS     bool
	IsVite       bool
	HasTailwind  bool
}

type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// DetectReactProject analyzes the current directory to detect a React project
func DetectReactProject(projectPath string) (*ReactProject, error) {
	packageJSONPath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no package.json found in %s", projectPath)
	}

	// Read and parse package.json
	packageJSONPath = filepath.Clean(packageJSONPath)
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json: %w", err)
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	// Check for React
	isReact := false
	if pkg.Dependencies != nil {
		if _, hasReact := pkg.Dependencies["react"]; hasReact {
			isReact = true
		}
	}
	if pkg.DevDependencies != nil {
		if _, hasReact := pkg.DevDependencies["react"]; hasReact {
			isReact = true
		}
	}

	if !isReact {
		return nil, fmt.Errorf("not a React project - no React dependency found")
	}

	project := &ReactProject{
		Path:        projectPath,
		PackageJSON: &pkg,
	}

	project.detectProjectType()

	return project, nil
}

func (rp *ReactProject) detectProjectType() {
	// Check for TypeScript
	tsConfigPath := filepath.Join(rp.Path, "tsconfig.json")
	if _, err := os.Stat(tsConfigPath); err == nil {
		rp.IsTypeScript = true
	}

	// Check for Next.js
	if rp.PackageJSON.Dependencies != nil {
		if _, hasNext := rp.PackageJSON.Dependencies["next"]; hasNext {
			rp.IsNextJS = true
		}
	}

	// Check for Vite
	viteConfigPaths := []string{
		filepath.Join(rp.Path, "vite.config.js"),
		filepath.Join(rp.Path, "vite.config.ts"),
	}
	for _, path := range viteConfigPaths {
		if _, err := os.Stat(path); err == nil {
			rp.IsVite = true
			break
		}
	}

	// Check for Tailwind CSS
	tailwindConfigPaths := []string{
		filepath.Join(rp.Path, "tailwind.config.js"),
		filepath.Join(rp.Path, "tailwind.config.ts"),
	}
	for _, path := range tailwindConfigPaths {
		if _, err := os.Stat(path); err == nil {
			rp.HasTailwind = true
			break
		}
	}

	// Also check in dependencies
	if rp.PackageJSON.Dependencies != nil {
		if _, hasTailwind := rp.PackageJSON.Dependencies["tailwindcss"]; hasTailwind {
			rp.HasTailwind = true
		}
	}
	if rp.PackageJSON.DevDependencies != nil {
		if _, hasTailwind := rp.PackageJSON.DevDependencies["tailwindcss"]; hasTailwind {
			rp.HasTailwind = true
		}
	}
}

// InstallVapi adds Vapi to the React project
func (rp *ReactProject) InstallVapi() error {
	// Install Vapi Web SDK
	fmt.Println("Installing Vapi Web SDK...")

	// Add Vapi dependency to package.json
	if rp.PackageJSON.Dependencies == nil {
		rp.PackageJSON.Dependencies = make(map[string]string)
	}
	rp.PackageJSON.Dependencies["@vapi-ai/web"] = "^2.0.0"

	// Save updated package.json
	if err := rp.savePackageJSON(); err != nil {
		return fmt.Errorf("failed to update package.json: %w", err)
	}

	fmt.Println("✓ Added @vapi-ai/web to package.json")
	return nil
}

// GenerateVapiComponents creates Vapi integration components
func (rp *ReactProject) GenerateVapiComponents() error {
	// Determine the source directory
	srcDir := "src"
	if rp.IsNextJS {
		// For Next.js, components usually go in components/ or src/components/
		if _, err := os.Stat(filepath.Join(rp.Path, "src")); err == nil {
			srcDir = "src"
		} else {
			srcDir = "."
		}
	}

	componentsDir := filepath.Join(rp.Path, srcDir, "components", "vapi")
	if err := os.MkdirAll(componentsDir, 0o750); err != nil {
		return fmt.Errorf("failed to create components directory: %w", err)
	}

	// Generate Vapi hook
	hookFile := "useVapi"
	if rp.IsTypeScript {
		hookFile += ".ts"
	} else {
		hookFile += ".js"
	}

	hookContent := rp.generateVapiHook()
	hookPath := filepath.Join(componentsDir, hookFile)
	if err := os.WriteFile(hookPath, []byte(hookContent), 0o600); err != nil {
		return fmt.Errorf("failed to write Vapi hook: %w", err)
	}

	// Generate Vapi component
	componentFile := "VapiButton"
	if rp.IsTypeScript {
		componentFile += ".tsx"
	} else {
		componentFile += ".jsx"
	}

	componentContent := rp.generateVapiComponent()
	componentPath := filepath.Join(componentsDir, componentFile)
	if err := os.WriteFile(componentPath, []byte(componentContent), 0o600); err != nil {
		return fmt.Errorf("failed to write Vapi component: %w", err)
	}

	// Generate example usage
	exampleFile := "VapiExample"
	if rp.IsTypeScript {
		exampleFile += ".tsx"
	} else {
		exampleFile += ".jsx"
	}

	exampleContent := rp.generateVapiExample()
	examplePath := filepath.Join(componentsDir, exampleFile)
	if err := os.WriteFile(examplePath, []byte(exampleContent), 0o600); err != nil {
		return fmt.Errorf("failed to write Vapi example: %w", err)
	}

	fmt.Printf("✓ Generated Vapi components in %s\n", componentsDir)
	return nil
}

// GenerateEnvTemplate creates a .env template with Vapi configuration
func (rp *ReactProject) GenerateEnvTemplate() error {
	envPath := filepath.Join(rp.Path, ".env.example")

	envContent := `# Vapi Configuration
REACT_APP_VAPI_PUBLIC_KEY=your_public_key_here
REACT_APP_VAPI_ASSISTANT_ID=your_assistant_id_here

# Optional: Vapi Server URL (defaults to production)
# REACT_APP_VAPI_BASE_URL=https://api.vapi.ai
`

	if rp.IsNextJS {
		envContent = `# Vapi Configuration
NEXT_PUBLIC_VAPI_PUBLIC_KEY=your_public_key_here
NEXT_PUBLIC_VAPI_ASSISTANT_ID=your_assistant_id_here

# Optional: Vapi Server URL (defaults to production)
# NEXT_PUBLIC_VAPI_BASE_URL=https://api.vapi.ai
`
	}

	if err := os.WriteFile(envPath, []byte(envContent), 0o600); err != nil {
		return fmt.Errorf("failed to create .env.example: %w", err)
	}

	fmt.Println("✓ Created .env.example with Vapi configuration")
	return nil
}

func (rp *ReactProject) savePackageJSON() error {
	data, err := json.MarshalIndent(rp.PackageJSON, "", "  ")
	if err != nil {
		return err
	}

	packageJSONPath := filepath.Join(rp.Path, "package.json")
	return os.WriteFile(packageJSONPath, data, 0o600)
}

func (rp *ReactProject) generateVapiHook() string {
	if rp.IsTypeScript {
		return `import { useState, useCallback, useEffect } from 'react';
import Vapi from '@vapi-ai/web';

interface VapiConfig {
  publicKey: string;
  assistantId: string;
  baseUrl?: string;
}

interface VapiState {
  isSessionActive: boolean;
  isLoading: boolean;
  error: string | null;
}

export const useVapi = (config: VapiConfig) => {
  const [vapi, setVapi] = useState<Vapi | null>(null);
  const [state, setState] = useState<VapiState>({
    isSessionActive: false,
    isLoading: false,
    error: null,
  });

  useEffect(() => {
    const vapiInstance = new Vapi(config.publicKey, config.baseUrl);
    setVapi(vapiInstance);

    const handleCallStart = () => {
      setState(prev => ({ ...prev, isSessionActive: true, isLoading: false }));
    };

    const handleCallEnd = () => {
      setState(prev => ({ ...prev, isSessionActive: false, isLoading: false }));
    };

    const handleError = (error: any) => {
      setState(prev => ({ ...prev, error: error.message, isLoading: false }));
    };

    vapiInstance.on('call-start', handleCallStart);
    vapiInstance.on('call-end', handleCallEnd);
    vapiInstance.on('error', handleError);

    return () => {
      vapiInstance.off('call-start', handleCallStart);
      vapiInstance.off('call-end', handleCallEnd);
      vapiInstance.off('error', handleError);
    };
  }, [config.publicKey, config.baseUrl]);

  const startCall = useCallback(async () => {
    if (!vapi) return;

    setState(prev => ({ ...prev, isLoading: true, error: null }));
    
    try {
      await vapi.start(config.assistantId);
    } catch (error: any) {
      setState(prev => ({ ...prev, error: error.message, isLoading: false }));
    }
  }, [vapi, config.assistantId]);

  const endCall = useCallback(() => {
    if (!vapi) return;
    vapi.stop();
  }, [vapi]);

  return {
    startCall,
    endCall,
    ...state,
  };
};
`
	}

	return `import { useState, useCallback, useEffect } from 'react';
import Vapi from '@vapi-ai/web';

export const useVapi = (config) => {
  const [vapi, setVapi] = useState(null);
  const [state, setState] = useState({
    isSessionActive: false,
    isLoading: false,
    error: null,
  });

  useEffect(() => {
    const vapiInstance = new Vapi(config.publicKey, config.baseUrl);
    setVapi(vapiInstance);

    const handleCallStart = () => {
      setState(prev => ({ ...prev, isSessionActive: true, isLoading: false }));
    };

    const handleCallEnd = () => {
      setState(prev => ({ ...prev, isSessionActive: false, isLoading: false }));
    };

    const handleError = (error) => {
      setState(prev => ({ ...prev, error: error.message, isLoading: false }));
    };

    vapiInstance.on('call-start', handleCallStart);
    vapiInstance.on('call-end', handleCallEnd);
    vapiInstance.on('error', handleError);

    return () => {
      vapiInstance.off('call-start', handleCallStart);
      vapiInstance.off('call-end', handleCallEnd);
      vapiInstance.off('error', handleError);
    };
  }, [config.publicKey, config.baseUrl]);

  const startCall = useCallback(async () => {
    if (!vapi) return;

    setState(prev => ({ ...prev, isLoading: true, error: null }));
    
    try {
      await vapi.start(config.assistantId);
    } catch (error) {
      setState(prev => ({ ...prev, error: error.message, isLoading: false }));
    }
  }, [vapi, config.assistantId]);

  const endCall = useCallback(() => {
    if (!vapi) return;
    vapi.stop();
  }, [vapi]);

  return {
    startCall,
    endCall,
    ...state,
  };
};
`
}

func (rp *ReactProject) generateVapiComponent() string {
	envPrefix := "REACT_APP_"
	if rp.IsNextJS {
		envPrefix = "NEXT_PUBLIC_"
	}

	if rp.IsTypeScript {
		if rp.HasTailwind {
			return fmt.Sprintf(`import React from 'react';
import { useVapi } from './useVapi';

interface VapiButtonProps {
  publicKey?: string;
  assistantId?: string;
  baseUrl?: string;
  className?: string;
  children?: React.ReactNode;
}

export const VapiButton: React.FC<VapiButtonProps> = ({
  publicKey = process.env.%sVAPI_PUBLIC_KEY,
  assistantId = process.env.%sVAPI_ASSISTANT_ID,
  baseUrl = process.env.%sVAPI_BASE_URL,
  className,
  children,
}) => {
  const { startCall, endCall, isSessionActive, isLoading, error } = useVapi({
    publicKey: publicKey || '',
    assistantId: assistantId || '',
    baseUrl,
  });

  const handleClick = () => {
    if (isSessionActive) {
      endCall();
    } else {
      startCall();
    }
  };

  if (!publicKey || !assistantId) {
    return (
      <div style={{ color: 'red', padding: '8px' }}>
        Missing Vapi configuration. Please set %sVAPI_PUBLIC_KEY and %sVAPI_ASSISTANT_ID environment variables.
      </div>
    );
  }

  return (
    <>
      <button
        onClick={handleClick}
        disabled={isLoading}
        className={className || "bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded disabled:opacity-50"}
      >
        {children || (isLoading ? 'Connecting...' : isSessionActive ? 'End Call' : 'Start Call')}
      </button>
      {error && (
        <div style={{ color: 'red', marginTop: '8px', fontSize: '14px' }}>
          Error: {error}
        </div>
      )}
    </>
  );
};
`, envPrefix, envPrefix, envPrefix, envPrefix, envPrefix)
		} else {
			// Non-Tailwind version with inline styles
			return fmt.Sprintf(`import React from 'react';
import { useVapi } from './useVapi';

interface VapiButtonProps {
  publicKey?: string;
  assistantId?: string;
  baseUrl?: string;
  className?: string;
  children?: React.ReactNode;
}

export const VapiButton: React.FC<VapiButtonProps> = ({
  publicKey = process.env.%sVAPI_PUBLIC_KEY,
  assistantId = process.env.%sVAPI_ASSISTANT_ID,
  baseUrl = process.env.%sVAPI_BASE_URL,
  className,
  children,
}) => {
  const { startCall, endCall, isSessionActive, isLoading, error } = useVapi({
    publicKey: publicKey || '',
    assistantId: assistantId || '',
    baseUrl,
  });

  const handleClick = () => {
    if (isSessionActive) {
      endCall();
    } else {
      startCall();
    }
  };

  if (!publicKey || !assistantId) {
    return (
      <div style={{ color: 'red', padding: '8px' }}>
        Missing Vapi configuration. Please set %sVAPI_PUBLIC_KEY and %sVAPI_ASSISTANT_ID environment variables.
      </div>
    );
  }

  const buttonStyle = {
    backgroundColor: isSessionActive ? '#dc2626' : '#3b82f6',
    color: 'white',
    border: 'none',
    padding: '8px 16px',
    borderRadius: '4px',
    fontWeight: 'bold' as const,
    cursor: isLoading ? 'not-allowed' : 'pointer',
    opacity: isLoading ? 0.5 : 1,
  };

  return (
    <>
      <button
        onClick={handleClick}
        disabled={isLoading}
        className={className}
        style={className ? undefined : buttonStyle}
      >
        {children || (isLoading ? 'Connecting...' : isSessionActive ? 'End Call' : 'Start Call')}
      </button>
      {error && (
        <div style={{ color: 'red', marginTop: '8px', fontSize: '14px' }}>
          Error: {error}
        </div>
      )}
    </>
  );
};
`, envPrefix, envPrefix, envPrefix, envPrefix, envPrefix)
		}
	}

	// JavaScript version
	if rp.HasTailwind {
		return fmt.Sprintf(`import React from 'react';
import { useVapi } from './useVapi';

export const VapiButton = ({
  publicKey = process.env.%sVAPI_PUBLIC_KEY,
  assistantId = process.env.%sVAPI_ASSISTANT_ID,
  baseUrl = process.env.%sVAPI_BASE_URL,
  className,
  children,
}) => {
  const { startCall, endCall, isSessionActive, isLoading, error } = useVapi({
    publicKey: publicKey || '',
    assistantId: assistantId || '',
    baseUrl,
  });

  const handleClick = () => {
    if (isSessionActive) {
      endCall();
    } else {
      startCall();
    }
  };

  if (!publicKey || !assistantId) {
    return (
      <div style={{ color: 'red', padding: '8px' }}>
        Missing Vapi configuration. Please set %sVAPI_PUBLIC_KEY and %sVAPI_ASSISTANT_ID environment variables.
      </div>
    );
  }

  return (
    <>
      <button
        onClick={handleClick}
        disabled={isLoading}
        className={className || "bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded disabled:opacity-50"}
      >
        {children || (isLoading ? 'Connecting...' : isSessionActive ? 'End Call' : 'Start Call')}
      </button>
      {error && (
        <div style={{ color: 'red', marginTop: '8px', fontSize: '14px' }}>
          Error: {error}
        </div>
      )}
    </>
  );
};
`, envPrefix, envPrefix, envPrefix, envPrefix, envPrefix)
	} else {
		// Non-Tailwind JavaScript version
		return fmt.Sprintf(`import React from 'react';
import { useVapi } from './useVapi';

export const VapiButton = ({
  publicKey = process.env.%sVAPI_PUBLIC_KEY,
  assistantId = process.env.%sVAPI_ASSISTANT_ID,
  baseUrl = process.env.%sVAPI_BASE_URL,
  className,
  children,
}) => {
  const { startCall, endCall, isSessionActive, isLoading, error } = useVapi({
    publicKey: publicKey || '',
    assistantId: assistantId || '',
    baseUrl,
  });

  const handleClick = () => {
    if (isSessionActive) {
      endCall();
    } else {
      startCall();
    }
  };

  if (!publicKey || !assistantId) {
    return (
      <div style={{ color: 'red', padding: '8px' }}>
        Missing Vapi configuration. Please set %sVAPI_PUBLIC_KEY and %sVAPI_ASSISTANT_ID environment variables.
      </div>
    );
  }

  const buttonStyle = {
    backgroundColor: isSessionActive ? '#dc2626' : '#3b82f6',
    color: 'white',
    border: 'none',
    padding: '8px 16px',
    borderRadius: '4px',
    fontWeight: 'bold',
    cursor: isLoading ? 'not-allowed' : 'pointer',
    opacity: isLoading ? 0.5 : 1,
  };

  return (
    <>
      <button
        onClick={handleClick}
        disabled={isLoading}
        className={className}
        style={className ? undefined : buttonStyle}
      >
        {children || (isLoading ? 'Connecting...' : isSessionActive ? 'End Call' : 'Start Call')}
      </button>
      {error && (
        <div style={{ color: 'red', marginTop: '8px', fontSize: '14px' }}>
          Error: {error}
        </div>
      )}
    </>
  );
};
`, envPrefix, envPrefix, envPrefix, envPrefix, envPrefix)
	}
}

func (rp *ReactProject) generateVapiExample() string {
	if rp.IsTypeScript {
		return `import React from 'react';
import { VapiButton } from './VapiButton';

export const VapiExample: React.FC = () => {
  return (
    <div style={{ padding: '20px', maxWidth: '600px', margin: '0 auto' }}>
      <h1>Vapi Integration Example</h1>
      <p>
        This example demonstrates how to integrate Vapi into your React application.
        Make sure to set your environment variables in your .env file.
      </p>
      
      <div style={{ marginTop: '20px' }}>
        <h2>Basic Usage</h2>
        <VapiButton />
      </div>

      <div style={{ marginTop: '20px' }}>
        <h2>Custom Button</h2>
        <VapiButton className="custom-vapi-button">
          🎤 Talk to AI Assistant
        </VapiButton>
      </div>

      <div style={{ marginTop: '20px' }}>
        <h2>Setup Instructions</h2>
        <ol>
          <li>Copy .env.example to .env</li>
          <li>Add your Vapi public key and assistant ID</li>
          <li>Install dependencies: npm install</li>
          <li>Start your development server</li>
        </ol>
      </div>
    </div>
  );
};
`
	}

	return `import React from 'react';
import { VapiButton } from './VapiButton';

export const VapiExample = () => {
  return (
    <div style={{ padding: '20px', maxWidth: '600px', margin: '0 auto' }}>
      <h1>Vapi Integration Example</h1>
      <p>
        This example demonstrates how to integrate Vapi into your React application.
        Make sure to set your environment variables in your .env file.
      </p>
      
      <div style={{ marginTop: '20px' }}>
        <h2>Basic Usage</h2>
        <VapiButton />
      </div>

      <div style={{ marginTop: '20px' }}>
        <h2>Custom Button</h2>
        <VapiButton className="custom-vapi-button">
          🎤 Talk to AI Assistant
        </VapiButton>
      </div>

      <div style={{ marginTop: '20px' }}>
        <h2>Setup Instructions</h2>
        <ol>
          <li>Copy .env.example to .env</li>
          <li>Add your Vapi public key and assistant ID</li>
          <li>Install dependencies: npm install</li>
          <li>Start your development server</li>
        </ol>
      </div>
    </div>
  );
};
`
}

// GenerateReactIntegration creates React/Next.js components and hooks for Vapi integration
func GenerateReactIntegration(projectPath string, info *ProjectInfo) error {
	// Determine the components directory based on project structure
	componentsDir := filepath.Join(projectPath, "src", "components")
	hooksDir := filepath.Join(projectPath, "src", "hooks")

	// For Next.js app directory structure
	if info.Framework == FrameworkNext {
		// Check if using app directory
		if _, err := os.Stat(filepath.Join(projectPath, "app")); err == nil {
			componentsDir = filepath.Join(projectPath, "app", "components")
			hooksDir = filepath.Join(projectPath, "app", "hooks")
		}
	}

	// Create directories if they don't exist
	if err := os.MkdirAll(componentsDir, 0o750); err != nil {
		return fmt.Errorf("failed to create components directory: %w", err)
	}
	if err := os.MkdirAll(hooksDir, 0o750); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// Generate files based on TypeScript preference
	ext := "jsx"
	if info.IsTypeScript {
		ext = "tsx"
	}

	// Generate useVapi hook
	if err := generateUseVapiHook(hooksDir, ext, info.IsTypeScript); err != nil {
		return err
	}

	// Generate VapiButton component
	if err := generateVapiButton(componentsDir, ext, info.IsTypeScript, info.HasTailwind); err != nil {
		return err
	}

	// Generate example component
	if err := generateVapiExample(componentsDir, ext, info.IsTypeScript); err != nil {
		return err
	}

	// Generate environment template
	if err := generateEnvTemplate(projectPath); err != nil {
		return err
	}

	return nil
}

// generateUseVapiHook creates the useVapi custom hook
func generateUseVapiHook(dir, ext string, isTypeScript bool) error {
	content := ""
	if isTypeScript {
		content = `import { useState, useCallback, useEffect } from 'react';
import Vapi from '@vapi-ai/web';

interface VapiConfig {
  publicKey: string;
  assistantId: string;
  baseUrl?: string;
}

interface VapiState {
  isSessionActive: boolean;
  isLoading: boolean;
  error: string | null;
}

export const useVapi = (config: VapiConfig) => {
  const [vapi, setVapi] = useState<Vapi | null>(null);
  const [state, setState] = useState<VapiState>({
    isSessionActive: false,
    isLoading: false,
    error: null,
  });

  useEffect(() => {
    const vapiInstance = new Vapi(config.publicKey, config.baseUrl);
    setVapi(vapiInstance);

    const handleCallStart = () => {
      setState(prev => ({ ...prev, isSessionActive: true, isLoading: false }));
    };

    const handleCallEnd = () => {
      setState(prev => ({ ...prev, isSessionActive: false, isLoading: false }));
    };

    const handleError = (error: any) => {
      setState(prev => ({ ...prev, error: error.message, isLoading: false }));
    };

    vapiInstance.on('call-start', handleCallStart);
    vapiInstance.on('call-end', handleCallEnd);
    vapiInstance.on('error', handleError);

    return () => {
      vapiInstance.off('call-start', handleCallStart);
      vapiInstance.off('call-end', handleCallEnd);
      vapiInstance.off('error', handleError);
    };
  }, [config.publicKey, config.baseUrl]);

  const startCall = useCallback(async () => {
    if (!vapi) return;

    setState(prev => ({ ...prev, isLoading: true, error: null }));
    
    try {
      await vapi.start(config.assistantId);
    } catch (error: any) {
      setState(prev => ({ ...prev, error: error.message, isLoading: false }));
    }
  }, [vapi, config.assistantId]);

  const endCall = useCallback(() => {
    if (!vapi) return;
    vapi.stop();
  }, [vapi]);

  return {
    startCall,
    endCall,
    ...state,
  };
};
`
	} else {
		content = `import { useState, useCallback, useEffect } from 'react';
import Vapi from '@vapi-ai/web';

export const useVapi = (config) => {
  const [vapi, setVapi] = useState(null);
  const [state, setState] = useState({
    isSessionActive: false,
    isLoading: false,
    error: null,
  });

  useEffect(() => {
    const vapiInstance = new Vapi(config.publicKey, config.baseUrl);
    setVapi(vapiInstance);

    const handleCallStart = () => {
      setState(prev => ({ ...prev, isSessionActive: true, isLoading: false }));
    };

    const handleCallEnd = () => {
      setState(prev => ({ ...prev, isSessionActive: false, isLoading: false }));
    };

    const handleError = (error) => {
      setState(prev => ({ ...prev, error: error.message, isLoading: false }));
    };

    vapiInstance.on('call-start', handleCallStart);
    vapiInstance.on('call-end', handleCallEnd);
    vapiInstance.on('error', handleError);

    return () => {
      vapiInstance.off('call-start', handleCallStart);
      vapiInstance.off('call-end', handleCallEnd);
      vapiInstance.off('error', handleError);
    };
  }, [config.publicKey, config.baseUrl]);

  const startCall = useCallback(async () => {
    if (!vapi) return;

    setState(prev => ({ ...prev, isLoading: true, error: null }));
    
    try {
      await vapi.start(config.assistantId);
    } catch (error) {
      setState(prev => ({ ...prev, error: error.message, isLoading: false }));
    }
  }, [vapi, config.assistantId]);

  const endCall = useCallback(() => {
    if (!vapi) return;
    vapi.stop();
  }, [vapi]);

  return {
    startCall,
    endCall,
    ...state,
  };
};
`
	}

	filename := fmt.Sprintf("useVapi.%s", ext)
	return os.WriteFile(filepath.Join(dir, filename), []byte(content), 0o600)
}

// generateVapiButton creates the VapiButton component
func generateVapiButton(dir, ext string, isTypeScript, hasTailwind bool) error {
	envPrefix := "REACT_APP_"
	// TODO: Check for Next.js to use NEXT_PUBLIC_ prefix

	content := ""
	if isTypeScript {
		if hasTailwind {
			content = fmt.Sprintf(`import React from 'react';
import { useVapi } from '../hooks/useVapi';

interface VapiButtonProps {
  publicKey?: string;
  assistantId?: string;
  baseUrl?: string;
  className?: string;
  children?: React.ReactNode;
}

export const VapiButton: React.FC<VapiButtonProps> = ({
  publicKey = process.env.%sVAPI_PUBLIC_KEY,
  assistantId = process.env.%sVAPI_ASSISTANT_ID,
  baseUrl = process.env.%sVAPI_BASE_URL,
  className,
  children,
}) => {
  const { startCall, endCall, isSessionActive, isLoading, error } = useVapi({
    publicKey: publicKey || '',
    assistantId: assistantId || '',
    baseUrl,
  });

  const handleClick = () => {
    if (isSessionActive) {
      endCall();
    } else {
      startCall();
    }
  };

  if (!publicKey || !assistantId) {
    return (
      <div className="text-red-500 p-2">
        Missing Vapi configuration. Please set environment variables.
      </div>
    );
  }

  return (
    <>
      <button
        onClick={handleClick}
        disabled={isLoading}
        className={className || "bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded disabled:opacity-50"}
      >
        {children || (isLoading ? 'Connecting...' : isSessionActive ? 'End Call' : 'Start Call')}
      </button>
      {error && (
        <div className="text-red-500 mt-2 text-sm">
          Error: {error}
        </div>
      )}
    </>
  );
};
`, envPrefix, envPrefix, envPrefix)
		} else {
			// TypeScript without Tailwind
			content = fmt.Sprintf(`import React from 'react';
import { useVapi } from '../hooks/useVapi';

interface VapiButtonProps {
  publicKey?: string;
  assistantId?: string;
  baseUrl?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
}

export const VapiButton: React.FC<VapiButtonProps> = ({
  publicKey = process.env.%sVAPI_PUBLIC_KEY,
  assistantId = process.env.%sVAPI_ASSISTANT_ID,
  baseUrl = process.env.%sVAPI_BASE_URL,
  style,
  children,
}) => {
  const { startCall, endCall, isSessionActive, isLoading, error } = useVapi({
    publicKey: publicKey || '',
    assistantId: assistantId || '',
    baseUrl,
  });

  const handleClick = () => {
    if (isSessionActive) {
      endCall();
    } else {
      startCall();
    }
  };

  if (!publicKey || !assistantId) {
    return (
      <div style={{ color: 'red', padding: '8px' }}>
        Missing Vapi configuration. Please set environment variables.
      </div>
    );
  }

  const buttonStyle: React.CSSProperties = {
    backgroundColor: '#3B82F6',
    color: 'white',
    fontWeight: 'bold',
    padding: '8px 16px',
    borderRadius: '4px',
    border: 'none',
    cursor: isLoading ? 'not-allowed' : 'pointer',
    opacity: isLoading ? 0.5 : 1,
    ...style
  };

  return (
    <>
      <button
        onClick={handleClick}
        disabled={isLoading}
        style={buttonStyle}
      >
        {children || (isLoading ? 'Connecting...' : isSessionActive ? 'End Call' : 'Start Call')}
      </button>
      {error && (
        <div style={{ color: 'red', marginTop: '8px', fontSize: '14px' }}>
          Error: {error}
        </div>
      )}
    </>
  );
};
`, envPrefix, envPrefix, envPrefix)
		}
	} else {
		// JavaScript version
		content = fmt.Sprintf(`import React from 'react';
import { useVapi } from '../hooks/useVapi';

export const VapiButton = ({
  publicKey = process.env.%sVAPI_PUBLIC_KEY,
  assistantId = process.env.%sVAPI_ASSISTANT_ID,
  baseUrl = process.env.%sVAPI_BASE_URL,
  style,
  children,
}) => {
  const { startCall, endCall, isSessionActive, isLoading, error } = useVapi({
    publicKey: publicKey || '',
    assistantId: assistantId || '',
    baseUrl,
  });

  const handleClick = () => {
    if (isSessionActive) {
      endCall();
    } else {
      startCall();
    }
  };

  if (!publicKey || !assistantId) {
    return (
      <div style={{ color: 'red', padding: '8px' }}>
        Missing Vapi configuration. Please set environment variables.
      </div>
    );
  }

  const buttonStyle = {
    backgroundColor: '#3B82F6',
    color: 'white',
    fontWeight: 'bold',
    padding: '8px 16px',
    borderRadius: '4px',
    border: 'none',
    cursor: isLoading ? 'not-allowed' : 'pointer',
    opacity: isLoading ? 0.5 : 1,
    ...style
  };

  return (
    <>
      <button
        onClick={handleClick}
        disabled={isLoading}
        style={buttonStyle}
      >
        {children || (isLoading ? 'Connecting...' : isSessionActive ? 'End Call' : 'Start Call')}
      </button>
      {error && (
        <div style={{ color: 'red', marginTop: '8px', fontSize: '14px' }}>
          Error: {error}
        </div>
      )}
    </>
  );
};
`, envPrefix, envPrefix, envPrefix)
	}

	filename := fmt.Sprintf("VapiButton.%s", ext)
	return os.WriteFile(filepath.Join(dir, filename), []byte(content), 0o600)
}

// generateVapiExample creates an example component showing Vapi usage
func generateVapiExample(dir, ext string, isTypeScript bool) error {
	content := ""
	if isTypeScript {
		content = `import React from 'react';
import { VapiButton } from './VapiButton';

export const VapiExample: React.FC = () => {
  return (
    <div style={{ padding: '20px' }}>
      <h2>Vapi Voice Assistant Example</h2>
      <p>Click the button below to start a voice conversation:</p>
      
      <VapiButton />
      
      <div style={{ marginTop: '20px' }}>
        <h3>How it works:</h3>
        <ol>
          <li>Make sure you have set up your environment variables</li>
          <li>Click "Start Call" to begin the conversation</li>
          <li>Speak naturally with the assistant</li>
          <li>Click "End Call" when finished</li>
        </ol>
      </div>
    </div>
  );
};
`
	} else {
		content = `import React from 'react';
import { VapiButton } from './VapiButton';

export const VapiExample = () => {
  return (
    <div style={{ padding: '20px' }}>
      <h2>Vapi Voice Assistant Example</h2>
      <p>Click the button below to start a voice conversation:</p>
      
      <VapiButton />
      
      <div style={{ marginTop: '20px' }}>
        <h3>How it works:</h3>
        <ol>
          <li>Make sure you have set up your environment variables</li>
          <li>Click "Start Call" to begin the conversation</li>
          <li>Speak naturally with the assistant</li>
          <li>Click "End Call" when finished</li>
        </ol>
      </div>
    </div>
  );
};
`
	}

	filename := fmt.Sprintf("VapiExample.%s", ext)
	return os.WriteFile(filepath.Join(dir, filename), []byte(content), 0o600)
}

// generateEnvTemplate creates environment template file
func generateEnvTemplate(projectPath string) error {
	// Check if it's a Next.js project
	envPrefix := "REACT_APP_"
	if _, err := os.Stat(filepath.Join(projectPath, "next.config.js")); err == nil {
		envPrefix = "NEXT_PUBLIC_"
	} else if _, err := os.Stat(filepath.Join(projectPath, "next.config.mjs")); err == nil {
		envPrefix = "NEXT_PUBLIC_"
	}

	content := fmt.Sprintf(`# Vapi Configuration
%sVAPI_PUBLIC_KEY=your_public_key_here
%sVAPI_ASSISTANT_ID=your_assistant_id_here

# Optional: Vapi Server URL (defaults to production)
# %sVAPI_BASE_URL=https://api.vapi.ai
`, envPrefix, envPrefix, envPrefix)

	return os.WriteFile(filepath.Join(projectPath, ".env.example"), []byte(content), 0o600)
}
