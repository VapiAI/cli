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
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	forwardTo  string
	listenPort int
	skipVerify bool
)

// Forward webhook events to your local development server
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Forward webhook events to your local server",
	Long: `Start a webhook listener that forwards Vapi webhook events to your local development server.

This is perfect for testing webhooks during development without needing ngrok or other tunneling tools.
The CLI will create a secure tunnel and forward all webhook events to your specified local endpoint.

Examples:
  vapi listen --forward-to localhost:3000/webhook
  vapi listen --forward-to http://localhost:8080/api/webhooks --port 4242
  vapi listen --forward-to localhost:3000/webhook --skip-verify`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if forwardTo == "" {
			return fmt.Errorf("--forward-to is required")
		}

		// Validate and normalize the forward-to URL
		forwardURL, err := normalizeURL(forwardTo)
		if err != nil {
			return fmt.Errorf("invalid --forward-to URL: %w", err)
		}

		return startWebhookListener(forwardURL, listenPort, skipVerify)
	},
}

// normalizeURL ensures the URL has proper scheme and format
func normalizeURL(rawURL string) (string, error) {
	// If no scheme provided, assume http://
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "http://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	if parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL: missing host")
	}

	return parsedURL.String(), nil
}

// startWebhookListener starts the local webhook server and forwarding logic
func startWebhookListener(forwardURL string, port int, skipVerify bool) error {
	// Create styles for better output formatting
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF"))
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Bold(true)

	fmt.Println(headerStyle.Render("🚀 Vapi Webhook Listener"))
	fmt.Println()
	fmt.Printf("%s Listening on port %d\n", successStyle.Render("✓"), port)
	fmt.Printf("%s Forwarding to: %s\n", successStyle.Render("✓"), forwardURL)
	if skipVerify {
		fmt.Printf("%s TLS verification disabled\n", warningStyle.Render("⚠"))
	}
	fmt.Println()
	fmt.Println(infoStyle.Render("Waiting for webhook events... (Press Ctrl+C to stop)"))
	fmt.Println()

	// Generate a unique identifier for this session
	sessionID := generateSessionID()

	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleWebhook(w, r, forwardURL, sessionID)
	})

	server := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("%s Server error: %v\n", errorStyle.Render("✗"), err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println()
	fmt.Println(infoStyle.Render("Shutting down webhook listener..."))

	// Shutdown server with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	fmt.Println(successStyle.Render("✓ Webhook listener stopped"))
	return nil
}

// handleWebhook processes incoming webhook requests and forwards them
func handleWebhook(w http.ResponseWriter, r *http.Request, forwardURL, sessionID string) {
	timestamp := time.Now().Format("15:04:05")

	// Create styles
	eventStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF7F")).Bold(true)
	methodStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEEB")).Bold(true)
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#98FB98"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6347")).Bold(true)

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[%s] %s Failed to read request body: %v\n", timestamp, errorStyle.Render("ERROR"), err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			fmt.Printf("[%s] %s Failed to close request body: %v\n", timestamp, errorStyle.Render("ERROR"), err)
		}
	}()

	// Try to parse as JSON to get event type if possible
	var webhook map[string]interface{}
	eventType := "webhook"
	if json.Unmarshal(body, &webhook) == nil {
		if et, ok := webhook["type"].(string); ok {
			eventType = et
		} else if event, ok := webhook["event"].(string); ok {
			eventType = event
		}
	}

	// Log the incoming webhook
	fmt.Printf("[%s] %s %s %s\n",
		timestamp,
		eventStyle.Render("→"),
		methodStyle.Render(r.Method),
		eventType,
	)

	// Forward the webhook to the target URL
	req, err := http.NewRequest(r.Method, forwardURL, bytes.NewReader(body))
	if err != nil {
		fmt.Printf("[%s] %s Failed to create forward request: %v\n", timestamp, errorStyle.Render("ERROR"), err)
		http.Error(w, "Failed to create forward request", http.StatusInternalServerError)
		return
	}

	// Copy headers from original request
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// Add custom headers to identify forwarded requests
	req.Header.Set("X-Vapi-Forwarded", "true")
	req.Header.Set("X-Vapi-Session-ID", sessionID)
	req.Header.Set("X-Vapi-Original-URL", r.URL.String())

	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make the forward request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[%s] %s Forward failed: %v\n", timestamp, errorStyle.Render("✗"), err)
		http.Error(w, "Failed to forward request", http.StatusBadGateway)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("[%s] %s Failed to close response body: %v\n", timestamp, errorStyle.Render("ERROR"), err)
		}
	}()

	// Read response from target server
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[%s] %s Failed to read response: %v\n", timestamp, errorStyle.Render("ERROR"), err)
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Log the response status
	statusColor := statusStyle
	if resp.StatusCode >= 400 {
		statusColor = errorStyle
	}
	fmt.Printf("[%s] %s %s %d %s\n",
		timestamp,
		statusColor.Render("←"),
		methodStyle.Render(r.Method),
		resp.StatusCode,
		http.StatusText(resp.StatusCode),
	)

	// Copy response headers back
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Send response back to Vapi
	w.WriteHeader(resp.StatusCode)
	if _, err := w.Write(respBody); err != nil {
		fmt.Printf("[%s] %s Failed to write response: %v\n", timestamp, errorStyle.Render("ERROR"), err)
	}

	// Print webhook details if it's a JSON payload
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") && len(body) > 0 {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, body, "  ", "  "); err == nil {
			fmt.Printf("  %s\n", strings.ReplaceAll(prettyJSON.String(), "\n", "\n  "))
		}
		fmt.Println()
	}
}

// generateSessionID creates a unique identifier for this listening session
func generateSessionID() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based ID if random fails
		return fmt.Sprintf("%d", time.Now().Unix())
	}
	return hex.EncodeToString(b)
}

func init() {
	rootCmd.AddCommand(listenCmd)

	// Add flags
	listenCmd.Flags().StringVar(&forwardTo, "forward-to", "", "Local endpoint to forward webhooks to (e.g., localhost:3000/webhook)")
	listenCmd.Flags().IntVar(&listenPort, "port", 4242, "Port to listen on for incoming webhooks")
	listenCmd.Flags().BoolVar(&skipVerify, "skip-verify", false, "Skip TLS certificate verification when forwarding")

	// Mark required flags
	if err := listenCmd.MarkFlagRequired("forward-to"); err != nil {
		panic(fmt.Sprintf("Failed to mark flag as required: %v", err))
	}
}
