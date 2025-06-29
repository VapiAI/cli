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
	"testing"
)

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		newer    string
		current  string
		expected bool
		hasError bool
	}{
		// Basic cases
		{"1.0.1", "1.0.0", true, false},
		{"1.1.0", "1.0.0", true, false},
		{"2.0.0", "1.0.0", true, false},
		{"1.0.0", "1.0.0", false, false},
		{"1.0.0", "1.0.1", false, false},
		{"1.0.0", "1.1.0", false, false},
		{"1.0.0", "2.0.0", false, false},

		// Real-world scenarios
		{"0.0.2", "0.0.1", true, false},
		{"0.1.0", "0.0.2", true, false},
		{"1.0.0", "0.9.9", true, false},
		{"0.0.1", "0.0.2", false, false}, // This is the key test case

		// Edge cases
		{"10.0.0", "9.0.0", true, false},
		{"1.10.0", "1.9.0", true, false},
		{"1.0.10", "1.0.9", true, false},

		// Error cases
		{"1.0", "1.0.0", false, true},
		{"1.0.0", "1.0", false, true},
		{"1.0.a", "1.0.0", false, true},
		{"", "1.0.0", false, true},
	}

	for _, test := range tests {
		result, err := isNewerVersion(test.newer, test.current)

		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for isNewerVersion(%q, %q), but got none", test.newer, test.current)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for isNewerVersion(%q, %q): %v", test.newer, test.current, err)
			continue
		}

		if result != test.expected {
			t.Errorf("isNewerVersion(%q, %q) = %v, expected %v", test.newer, test.current, result, test.expected)
		}
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		version  string
		expected [3]int
		hasError bool
	}{
		{"1.0.0", [3]int{1, 0, 0}, false},
		{"0.0.1", [3]int{0, 0, 1}, false},
		{"10.20.30", [3]int{10, 20, 30}, false},
		{"1.0", [3]int{}, true},
		{"1.0.0.0", [3]int{}, true},
		{"1.0.a", [3]int{}, true},
		{"", [3]int{}, true},
	}

	for _, test := range tests {
		result, err := parseVersion(test.version)

		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for parseVersion(%q), but got none", test.version)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for parseVersion(%q): %v", test.version, err)
			continue
		}

		if result != test.expected {
			t.Errorf("parseVersion(%q) = %v, expected %v", test.version, result, test.expected)
		}
	}
}
