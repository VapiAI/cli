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
package main

import (
	"github.com/VapiAI/cli/cmd"
	"github.com/VapiAI/cli/pkg/version"
)

// Build variables set by goreleaser
var (
	buildVersion = "dev" // Will be overridden by -ldflags
	commit       = "none"
	date         = "unknown"
	builtBy      = "unknown"
)

func main() {
	// Set build information in the version package
	version.SetBuildInfo(buildVersion, commit, date, builtBy)

	// Set version information for the CLI commands
	cmd.SetVersion(version.Get(), version.GetCommit(), version.GetDate(), version.GetBuiltBy())

	// Execute the CLI
	cmd.Execute()
}
