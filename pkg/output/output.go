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
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

type OutputFormat string

const (
	FormatJSON  OutputFormat = "json"
	FormatTable OutputFormat = "table"
	FormatYAML  OutputFormat = "yaml"
)

func PrintJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func PrintTable(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print headers
	for i, h := range headers {
		if i > 0 {
			if _, err := fmt.Fprint(w, "\t"); err != nil {
				return
			}
		}
		if _, err := fmt.Fprint(w, h); err != nil {
			return
		}
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return
	}

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				if _, err := fmt.Fprint(w, "\t"); err != nil {
					return
				}
			}
			if _, err := fmt.Fprint(w, cell); err != nil {
				return
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return
		}
	}

	// Flush buffer
	if err := w.Flush(); err != nil {
		fmt.Printf("Warning: failed to flush table writer: %v\n", err)
	}
}
