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
	"strings"
	"time"

	"github.com/VapiAI/cli/pkg/config"
)

type AuthManager struct {
	callbackPort int
	authCode     chan string
	authError    chan error
	orgName      string
	orgID        string
	email        string
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

		// Try to capture optional organization name from query params if provided by dashboard
		// These parameters are optional and will be ignored if not present
		if org := r.URL.Query().Get("org_name"); org != "" {
			a.orgName = org
		} else if org := r.URL.Query().Get("organization"); org != "" {
			a.orgName = org
		} else if org := r.URL.Query().Get("org"); org != "" {
			a.orgName = org
		}

		// Organization ID for reliable deduplication
		if orgID := r.URL.Query().Get("org_id"); orgID != "" {
			a.orgID = orgID
		}

		// Optional email
		if email := r.URL.Query().Get("email"); email != "" {
			a.email = email
		} else if email := r.URL.Query().Get("user_email"); email != "" {
			a.email = email
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
	return LoginWithAccountName("")
}

// LoginWithAccountName performs authentication and optionally saves with a specific account name
func LoginWithAccountName(accountName string) error {
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

	// De-duplicate: if an existing account matches org ID (preferred), organization name, or API key,
	// update it instead of creating a new one
	if cfg.Accounts != nil {
		// First prefer matching by org ID when provided (most reliable)
		if authManager.orgID != "" {
			for existingName, existing := range cfg.Accounts {
				if existing.OrgID != authManager.orgID {
					continue
				}
				cfg.Accounts[existingName] = config.Account{
					APIKey:       apiKey,
					Organization: authManager.orgName,
					OrgID:        authManager.orgID,
					Environment:  cfg.GetEnvironment(),
					LoginTime:    time.Now().Format(time.RFC3339),
					Email:        firstNonEmpty(authManager.email, existing.Email),
				}
				cfg.ActiveAccount = existingName

				if err := config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("failed to save API key: %w", err)
				}

				fmt.Printf("\n✅ Re-authenticated existing account '%s'!\n", existingName)
				fmt.Printf("Organization: %s\n", authManager.orgName)
				if authManager.email != "" {
					fmt.Printf("Email: %s\n", authManager.email)
				}
				if len(cfg.Accounts) > 1 {
					fmt.Println("💡 Use 'vapi auth switch' to switch between accounts")
				}
				return nil
			}
		}

		// Then try matching by organization name when provided
		if authManager.orgName != "" {
			for existingName, existing := range cfg.Accounts {
				if !strings.EqualFold(existing.Organization, authManager.orgName) {
					continue
				}
				cfg.Accounts[existingName] = config.Account{
					APIKey:       apiKey,
					Organization: authManager.orgName,
					OrgID:        authManager.orgID,
					Environment:  cfg.GetEnvironment(),
					LoginTime:    time.Now().Format(time.RFC3339),
					Email:        firstNonEmpty(authManager.email, existing.Email),
				}
				cfg.ActiveAccount = existingName

				if err := config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("failed to save API key: %w", err)
				}

				fmt.Printf("\n✅ Re-authenticated existing account '%s'!\n", existingName)
				fmt.Printf("Organization: %s\n", authManager.orgName)
				if authManager.email != "" {
					fmt.Printf("Email: %s\n", authManager.email)
				}
				if len(cfg.Accounts) > 1 {
					fmt.Println("💡 Use 'vapi auth switch' to switch between accounts")
				}
				return nil
			}
		}

		// Fallback: match by API key if the same key already exists
		for existingName, existing := range cfg.Accounts {
			if existing.APIKey != apiKey {
				continue
			}
			// Update existing account's login time and organization (if newly available)
			orgName := authManager.orgName
			if orgName == "" {
				orgName = existing.Organization
			}

			cfg.Accounts[existingName] = config.Account{
				APIKey:       apiKey,
				Organization: orgName,
				OrgID:        authManager.orgID,
				Environment:  cfg.GetEnvironment(),
				LoginTime:    time.Now().Format(time.RFC3339),
				Email:        firstNonEmpty(authManager.email, existing.Email),
			}
			// Set as active account
			cfg.ActiveAccount = existingName

			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to save API key: %w", err)
			}

			fmt.Printf("\n✅ Re-authenticated existing account '%s'!\n", existingName)
			if orgName != "" {
				fmt.Printf("Organization: %s\n", orgName)
			}
			if authManager.email != "" {
				fmt.Printf("Email: %s\n", authManager.email)
			}
			if len(cfg.Accounts) > 1 {
				fmt.Println("💡 Use 'vapi auth switch' to switch between accounts")
			}
			return nil
		}
	}

	// Generate account name if not provided
	if accountName == "" {
		accountName = fmt.Sprintf("account-%d", time.Now().Unix())
	}

	// Add as new account (supports multiple accounts). Use org info if we captured it.
	cfg.AddAccount(accountName, apiKey, authManager.orgName, authManager.orgID, authManager.email)

	// For backward compatibility, also set legacy APIKey field if it's the first account
	if len(cfg.Accounts) == 1 {
		cfg.APIKey = apiKey
	}

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save API key: %w", err)
	}

	fmt.Printf("\n✅ Successfully authenticated as '%s'! Your API key has been saved.\n", accountName)
	if authManager.orgName != "" {
		fmt.Printf("Organization: %s\n", authManager.orgName)
	}
	if authManager.email != "" {
		fmt.Printf("Email: %s\n", authManager.email)
	}
	if len(cfg.Accounts) > 1 {
		fmt.Println("💡 Use 'vapi auth switch' to switch between accounts")
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// Logout clears the stored authentication credentials
func Logout() error {
	return LogoutAccount("")
}

// LogoutAccount logs out from a specific account, or all accounts if accountName is empty
func LogoutAccount(accountName string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = &config.Config{}
	}

	if accountName == "" {
		// Logout from all accounts
		cfg.Accounts = make(map[string]config.Account)
		cfg.ActiveAccount = ""
		cfg.APIKey = "" // Also clear legacy API key

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to clear API keys: %w", err)
		}

		fmt.Println("🔓 Successfully logged out from all accounts!")
		return nil
	}

	// Logout from specific account
	if err := cfg.RemoveAccount(accountName); err != nil {
		return fmt.Errorf("failed to remove account: %w", err)
	}

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("🔓 Successfully logged out from account '%s'!\n", accountName)

	// Show remaining accounts
	remaining := cfg.ListAccounts()
	if len(remaining) > 0 {
		fmt.Printf("Active account: %s\n", cfg.ActiveAccount)
		fmt.Printf("Remaining accounts: %d\n", len(remaining))
	} else {
		fmt.Println("No accounts remaining. Use 'vapi auth login' to authenticate.")
	}

	return nil
}

// AuthStatus represents the current authentication status
type AuthStatus struct {
	IsAuthenticated bool
	APIKeySet       bool
	APIKeySource    string
	Environment     string
	BaseURL         string
	DashboardURL    string
	ActiveAccount   string
	TotalAccounts   int
	Accounts        map[string]config.Account
}

// GetStatus returns the current authentication status
func GetStatus() (*AuthStatus, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	status := &AuthStatus{
		Environment:   cfg.GetEnvironment(),
		BaseURL:       cfg.GetAPIBaseURL(),
		DashboardURL:  cfg.GetDashboardURL(),
		ActiveAccount: cfg.ActiveAccount,
		TotalAccounts: len(cfg.ListAccounts()),
		Accounts:      cfg.ListAccounts(),
	}

	// Check if API key is set from various sources
	apiKey := cfg.GetActiveAPIKey()
	apiKeySource := cfg.GetAPIKeySource()

	status.APIKeySet = apiKey != ""
	status.APIKeySource = apiKeySource

	// Determine if authenticated:
	// 1. If API key from environment variable, always authenticated if not empty
	// 2. If from account, authenticated if active account exists and has valid API key
	// 3. If from legacy config, authenticated if API key is not empty
	if apiKey != "" {
		switch apiKeySource {
		case "environment variable":
			status.IsAuthenticated = true
		case "config file (legacy)": // #nosec G101 - This is not a hardcoded credential but a source description
			status.IsAuthenticated = true
		default:
			// From account - check if active account exists and has API key
			if status.TotalAccounts > 0 && status.ActiveAccount != "" {
				if activeAccount := cfg.GetActiveAccount(); activeAccount != nil && activeAccount.APIKey != "" {
					status.IsAuthenticated = true
				} else {
					status.IsAuthenticated = false
				}
			} else {
				status.IsAuthenticated = false
			}
		}
	} else {
		status.IsAuthenticated = false
	}

	return status, nil
}
