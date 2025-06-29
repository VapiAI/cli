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
package version

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Version information that can be overridden at build time
var (
	// These will be set by goreleaser during release builds
	version = "" // Will be overridden by -ldflags
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// Get returns the current version, preferring build-time version over VERSION file
func Get() string {
	// If version was set at build time (e.g., by goreleaser), use that
	if version != "" && version != "dev" {
		return version
	}

	// Try to read from VERSION file
	if v := readVersionFile(); v != "" {
		return v
	}

	// Fallback to dev
	return "dev"
}

// readVersionFile attempts to read the VERSION file from the project root
func readVersionFile() string {
	// Get the directory of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	// Navigate to project root (pkg/version -> project root)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	versionFile := filepath.Join(projectRoot, "VERSION")

	// Security: Validate that the file path is what we expect
	// Only allow reading VERSION file from project root
	if filepath.Base(versionFile) != "VERSION" {
		return ""
	}

	// Additional security check: ensure the path doesn't contain traversal attempts
	cleanPath := filepath.Clean(versionFile)
	if cleanPath != versionFile {
		return ""
	}

	// Read the file
	data, err := os.ReadFile(cleanPath) // #nosec G304 - path is validated above
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

// GetCommit returns the git commit hash
func GetCommit() string {
	return commit
}

// GetDate returns the build date
func GetDate() string {
	return date
}

// GetBuiltBy returns who/what built this binary
func GetBuiltBy() string {
	return builtBy
}

// SetBuildInfo sets the build information (called from main)
func SetBuildInfo(v, c, d, b string) {
	if v != "" {
		version = v
	}
	if c != "" {
		commit = c
	}
	if d != "" {
		date = d
	}
	if b != "" {
		builtBy = b
	}
}
