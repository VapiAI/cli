package analytics

import (
	"crypto/sha256"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/VapiAI/cli/pkg/config"
	"github.com/posthog/posthog-go"
)

const (
	// PostHog project configuration
	PostHogAPIKey = "phc_oX8UUqsZVs6ifWJRTZ0yNtwTv852ccqRsk09SsZbHHb"
	PostHogHost   = "https://us.i.posthog.com"

	// Event names following PostHog conventions
	EventCommandExecuted = "cli_command_executed"
	EventCommandFailed   = "cli_command_failed"
	EventSessionStarted  = "cli_session_started"
	EventError           = "cli_error"
)

// Client wraps PostHog client with privacy controls
type Client struct {
	client     posthog.Client
	enabled    bool
	distinctID string
}

var globalClient *Client

// Initialize creates and configures the analytics client
func Initialize() {
	// Check opt-out mechanisms (multiple ways to disable)
	if isOptedOut() {
		globalClient = &Client{enabled: false}
		return
	}

	// Create PostHog client with privacy-focused configuration
	client, err := posthog.NewWithConfig(
		PostHogAPIKey,
		posthog.Config{
			Endpoint:                  PostHogHost,
			FeatureFlagRequestTimeout: time.Second * 2, // Short timeout
		},
	)
	if err != nil {
		// Fail silently - analytics should never break the CLI
		globalClient = &Client{enabled: false}
		return
	}

	globalClient = &Client{
		client:     client,
		enabled:    true,
		distinctID: generateAnonymousID(),
	}

	// Track session start
	TrackEvent(EventSessionStarted, map[string]interface{}{
		"version":    getVersion(),
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"go_version": runtime.Version(),
	})
}

// Close properly shuts down the analytics client
func Close() {
	if globalClient != nil && globalClient.enabled && globalClient.client != nil {
		globalClient.client.Close()
	}
}

// TrackEvent sends an analytics event (no-op if opted out)
func TrackEvent(event string, properties map[string]interface{}) {
	if globalClient == nil || !globalClient.enabled {
		return
	}

	// Always add anonymous tracking flag and timestamp
	if properties == nil {
		properties = make(map[string]interface{})
	}

	// Ensure anonymous tracking (no PII)
	properties["$process_person_profile"] = false
	properties["anonymous"] = true
	properties["timestamp"] = time.Now().UTC().Format(time.RFC3339)

	// Send event asynchronously (non-blocking)
	globalClient.client.Enqueue(posthog.Capture{
		DistinctId: globalClient.distinctID,
		Event:      event,
		Properties: properties,
	})
}

// TrackCommand tracks CLI command usage
func TrackCommand(command string, subcommand string, success bool, duration time.Duration, errorMsg string) {
	properties := map[string]interface{}{
		"command":     command,
		"subcommand":  subcommand,
		"success":     success,
		"duration_ms": duration.Milliseconds(),
	}

	eventName := EventCommandExecuted
	if !success {
		eventName = EventCommandFailed
		if errorMsg != "" {
			// Hash error message to avoid PII while preserving error patterns
			properties["error_hash"] = hashString(errorMsg)
			properties["error_type"] = categorizeError(errorMsg)
		}
	}

	TrackEvent(eventName, properties)
}

// TrackError tracks errors with anonymized information
func TrackError(errorMsg string, context map[string]interface{}) {
	properties := map[string]interface{}{
		"error_hash": hashString(errorMsg),
		"error_type": categorizeError(errorMsg),
	}

	// Add context without PII
	for k, v := range context {
		if isSafeProperty(k) {
			properties[k] = v
		}
	}

	TrackEvent(EventError, properties)
}

// isOptedOut checks various opt-out mechanisms
func isOptedOut() bool {
	// Check environment variables (multiple options for user preference)
	envVars := []string{
		"VAPI_DISABLE_ANALYTICS",
		"VAPI_NO_TELEMETRY",
		"DISABLE_TELEMETRY",
		"DO_NOT_TRACK",
	}

	for _, env := range envVars {
		if value := os.Getenv(env); value != "" {
			// Any non-empty value means opt-out (1, true, yes, etc.)
			return true
		}
	}

	// Check config file setting
	cfg := config.GetConfig()
	if cfg != nil && cfg.DisableAnalytics {
		return true
	}

	return false
}

// generateAnonymousID creates a stable anonymous identifier
func generateAnonymousID() string {
	// Use a combination of system info that's stable but not personally identifiable
	// This allows us to track usage patterns without identifying users
	identifier := fmt.Sprintf("%s-%s-%s",
		runtime.GOOS,
		runtime.GOARCH,
		getInstallationID(),
	)

	hash := sha256.Sum256([]byte(identifier))
	return fmt.Sprintf("cli_%x", hash[:8]) // Use first 8 bytes of hash
}

// getInstallationID generates a stable installation identifier
func getInstallationID() string {
	// Try to get a stable identifier from config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "unknown"
	}

	// Hash the config directory path for a stable but anonymous ID
	hash := sha256.Sum256([]byte(configDir))
	return fmt.Sprintf("%x", hash[:4])
}

// hashString creates a hash for error categorization without exposing PII
func hashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash[:8])
}

// categorizeError provides error categories for analytics
func categorizeError(errorMsg string) string {
	lower := strings.ToLower(errorMsg)

	switch {
	case strings.Contains(lower, "network") || strings.Contains(lower, "connection"):
		return "network"
	case strings.Contains(lower, "auth") || strings.Contains(lower, "unauthorized"):
		return "auth"
	case strings.Contains(lower, "not found") || strings.Contains(lower, "404"):
		return "not_found"
	case strings.Contains(lower, "timeout"):
		return "timeout"
	case strings.Contains(lower, "rate limit"):
		return "rate_limit"
	case strings.Contains(lower, "validation") || strings.Contains(lower, "invalid"):
		return "validation"
	case strings.Contains(lower, "permission") || strings.Contains(lower, "forbidden"):
		return "permission"
	default:
		return "other"
	}
}

// isSafeProperty checks if a property is safe to include (no PII)
func isSafeProperty(key string) bool {
	safeKeys := map[string]bool{
		"command":     true,
		"subcommand":  true,
		"success":     true,
		"duration_ms": true,
		"error_type":  true,
		"version":     true,
		"os":          true,
		"arch":        true,
	}
	return safeKeys[key]
}

// getVersion returns the CLI version
func getVersion() string {
	// This will be set by the build process
	if version := os.Getenv("VAPI_VERSION"); version != "" {
		return version
	}
	return "dev"
}

// IsEnabled returns whether analytics is enabled
func IsEnabled() bool {
	return globalClient != nil && globalClient.enabled
}
