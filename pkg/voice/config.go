package voice

import (
	"time"
)

// WebRTCConfig holds configuration for WebRTC functionality
type WebRTCConfig struct {
	// Vapi API Configuration
	VapiAPIKey       string `mapstructure:"vapi_api_key"`       // Private API key
	VapiPublicAPIKey string `mapstructure:"vapi_public_api_key"` // Public API key for /call/web
	VapiBaseURL      string `mapstructure:"vapi_base_url"`

	// Daily.co Configuration (legacy - now handled by Vapi)
	DailyAPIKey string `mapstructure:"daily_api_key"`
	DailyDomain string `mapstructure:"daily_domain"`

	// WebRTC Configuration
	STUNServers []string `mapstructure:"stun_servers"`
	TURNServers []string `mapstructure:"turn_servers"`

	// Audio Configuration
	AudioInputDevice  string `mapstructure:"audio_input_device"`
	AudioOutputDevice string `mapstructure:"audio_output_device"`
	SampleRate        int    `mapstructure:"sample_rate"`
	BufferSize        int    `mapstructure:"buffer_size"`

	// Call Configuration
	CallTimeout time.Duration `mapstructure:"call_timeout"`
	VideoEnabled bool         `mapstructure:"video_enabled"`
}

// DefaultWebRTCConfig returns default WebRTC configuration
func DefaultWebRTCConfig() *WebRTCConfig {
	return &WebRTCConfig{
		// Default Vapi API configuration
		VapiBaseURL: "https://api.vapi.ai",
		
		// Default to Vapi's Daily.co subdomain for WebRTC calls (legacy)
		DailyDomain: "vapi",
		STUNServers: []string{
			"stun:stun.l.google.com:19302",
			"stun:stun1.l.google.com:19302",
		},
		AudioInputDevice:  "default",
		AudioOutputDevice: "default",
		SampleRate:        48000,
		BufferSize:        480,
		CallTimeout:       30 * time.Minute,
		VideoEnabled:      false, // Audio-only by default
	}
}

// getAPIKey returns the Vapi API key for authentication
func (c *WebRTCConfig) getAPIKey() string {
	return c.VapiAPIKey
}

// getPublicAPIKey returns the Vapi public API key for /call/web endpoint
func (c *WebRTCConfig) getPublicAPIKey() string {
	if c.VapiPublicAPIKey != "" {
		return c.VapiPublicAPIKey
	}
	// Fallback to private key if no public key is set
	return c.VapiAPIKey
}

// getPrivateAPIKey returns the Vapi private API key for /call endpoint
func (c *WebRTCConfig) getPrivateAPIKey() string {
	return c.VapiAPIKey
}

// getAPIBaseURL returns the Vapi API base URL
func (c *WebRTCConfig) getAPIBaseURL() string {
	if c.VapiBaseURL == "" {
		return "https://api.vapi.ai"
	}
	return c.VapiBaseURL
}