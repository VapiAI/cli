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
package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"time"

	"github.com/VapiAI/cli/pkg/config"
)

type AuthManager struct {
	callbackPort int
	authCode     chan string
	authError    chan error
}

func NewAuthManager() *AuthManager {
	return &AuthManager{
		callbackPort: 0, // Will be assigned dynamically
		authCode:     make(chan string, 1),
		authError:    make(chan error, 1),
	}
}

// Authenticate opens the browser for OAuth-style authentication
func (a *AuthManager) Authenticate() (string, error) {
	// Start local server to receive callback
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("failed to start callback server: %w", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			fmt.Printf("Warning: failed to close listener: %v\n", err)
		}
	}()

	// Get the assigned port
	a.callbackPort = listener.Addr().(*net.TCPAddr).Port

	// Generate state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	// Build auth URL
	authURL := a.buildAuthURL(state)

	// Start HTTP server in background
	server := &http.Server{
		Handler:           a.createCallbackHandler(state),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			a.authError <- fmt.Errorf("callback server error: %w", err)
		}
	}()

	// Open browser
	fmt.Println("🔐 Opening browser for authentication...")
	fmt.Printf("If browser doesn't open automatically, visit:\n%s\n\n", authURL)

	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Warning: Failed to open browser automatically: %v\n", err)
	}

	// Wait for auth response with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	select {
	case apiKey := <-a.authCode:
		// Shutdown server
		if err := server.Shutdown(context.Background()); err != nil {
			fmt.Printf("Warning: failed to shutdown server: %v\n", err)
		}
		return apiKey, nil

	case err := <-a.authError:
		if shutdownErr := server.Shutdown(context.Background()); shutdownErr != nil {
			fmt.Printf("Warning: failed to shutdown server: %v\n", shutdownErr)
		}
		return "", err

	case <-ctx.Done():
		if err := server.Shutdown(context.Background()); err != nil {
			fmt.Printf("Warning: failed to shutdown server: %v\n", err)
		}
		return "", fmt.Errorf("authentication timeout")
	}
}

func (a *AuthManager) buildAuthURL(state string) string {
	// Load config to get environment-specific dashboard URL
	cfg, err := config.LoadConfig()
	if err != nil {
		// Fallback to production if config loading fails
		cfg = &config.Config{}
		cfg.Environment = "production"
		cfg.DashboardURL = "https://dashboard.vapi.ai"
	}

	baseURL := fmt.Sprintf("%s/auth/cli", cfg.GetDashboardURL())
	params := url.Values{}
	params.Add("state", state)
	params.Add("redirect_uri", fmt.Sprintf("http://localhost:%d/callback", a.callbackPort))

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

func (a *AuthManager) createCallbackHandler(expectedState string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only handle callback path
		if r.URL.Path != "/callback" {
			http.NotFound(w, r)
			return
		}

		// Verify state
		state := r.URL.Query().Get("state")
		if state != expectedState {
			a.authError <- fmt.Errorf("invalid state parameter")
			a.writeErrorPage(w, "Invalid state parameter")
			return
		}

		// Get API key from query params
		apiKey := r.URL.Query().Get("api_key")
		if apiKey == "" {
			a.authError <- fmt.Errorf("no API key received")
			a.writeErrorPage(w, "No API key received")
			return
		}

		// Send success response
		a.writeSuccessPage(w)

		// Send API key through channel
		a.authCode <- apiKey
	})
}

func (a *AuthManager) writeSuccessPage(w http.ResponseWriter) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Vapi CLI - Authentication Successful</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            display: flex;
            align-items: center;
            justify-content: center;
            height: 100vh;
            margin: 0;
            background: #0a0a0a;
            color: #fff;
        }
        .container {
            text-align: center;
            padding: 40px;
            background: #1a1a1a;
            border-radius: 12px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
        }
        h1 { 
            color: #4ADE80; 
            margin-bottom: 16px;
        }
        p { 
            color: #a1a1aa; 
            margin-top: 8px;
        }
        .icon {
            font-size: 48px;
            margin-bottom: 16px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="icon">✅</div>
        <h1>Authentication Successful!</h1>
        <p>You can now close this window and return to your terminal.</p>
    </div>
    <script>
        setTimeout(() => window.close(), 3000);
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(html)); err != nil {
		fmt.Printf("Warning: failed to write response: %v\n", err)
	}
}

func (a *AuthManager) writeErrorPage(w http.ResponseWriter, message string) {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Vapi CLI - Authentication Failed</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            display: flex;
            align-items: center;
            justify-content: center;
            height: 100vh;
            margin: 0;
            background: #0a0a0a;
            color: #fff;
        }
        .container {
            text-align: center;
            padding: 40px;
            background: #1a1a1a;
            border-radius: 12px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
        }
        h1 { 
            color: #EF4444; 
            margin-bottom: 16px;
        }
        p { 
            color: #a1a1aa; 
            margin-top: 8px;
        }
        .icon {
            font-size: 48px;
            margin-bottom: 16px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="icon">❌</div>
        <h1>Authentication Failed</h1>
        <p>%s</p>
        <p>Please try again or check your terminal for more information.</p>
    </div>
</body>
</html>`, message)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusBadRequest)
	if _, err := w.Write([]byte(html)); err != nil {
		fmt.Printf("Warning: failed to write response: %v\n", err)
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func openBrowser(targetURL string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/c", "start", targetURL).Start()
	case "darwin":
		return exec.Command("open", targetURL).Start()
	default: // Linux and others
		return exec.Command("xdg-open", targetURL).Start()
	}
}

// Login performs the browser-based authentication flow
func Login() error {
	authManager := NewAuthManager()

	apiKey, err := authManager.Authenticate()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Save the API key to config
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = &config.Config{}
	}

	cfg.APIKey = apiKey

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save API key: %w", err)
	}

	fmt.Println("\n✅ Successfully authenticated! Your API key has been saved.")
	return nil
}
