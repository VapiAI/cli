package voice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	vapiclient "github.com/VapiAI/server-sdk-go/client"
)

// CallStatus represents the current state of a voice call
type CallStatus string

const (
	CallStatusIdle         CallStatus = "idle"
	CallStatusConnecting   CallStatus = "connecting"
	CallStatusConnected    CallStatus = "connected"
	CallStatusDisconnected CallStatus = "disconnected"
	CallStatusFailed       CallStatus = "failed"
)

// CallState holds the current state of a voice call
type CallState struct {
	CallID       string
	AssistantID  string
	Status       CallStatus
	StartTime    time.Time
	WebSocketURL string
}

// APIRequest represents a request to the Vapi API
type APIRequest struct {
	Method    string
	URL       string
	Headers   map[string]string
	Body      interface{}
	Timestamp time.Time
}

// APIResponse represents a response from the Vapi API
type APIResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       interface{}
	Duration   time.Duration
	Timestamp  time.Time
}

// VoiceClient manages voice calls with Vapi WebSocket transport
type VoiceClient struct {
	config     *WebRTCConfig
	vapiClient *vapiclient.Client
	callState  *CallState

	// Audio pipeline
	audioStream *AudioStream

	// WebSocket signaling
	signaling *VapiWebSocket

	// Audio processing
	audioProcessor *WebSocketAudioProcessor

	// Echo cancellation state
	lastSpeakerSamples []float32

	// Event channels
	requestLog  chan APIRequest
	responseLog chan APIResponse
	callEvents  chan CallEvent
}

// CallEvent represents events during a voice call
type CallEvent struct {
	Type      string
	Data      interface{}
	Timestamp time.Time
}

// NewVoiceClient creates a new voice client
func NewVoiceClient(config *WebRTCConfig, vapiClient *vapiclient.Client) (*VoiceClient, error) {
	if config == nil {
		config = DefaultWebRTCConfig()
	}

	// Create audio stream
	audioStream, err := NewAudioStream(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio stream: %w", err)
	}

	// Create WebSocket signaling client
	signaling := NewVapiWebSocket()

	// Create audio processor
	audioProcessor, err := NewWebSocketAudioProcessor()
	if err != nil {
		return nil, fmt.Errorf("failed to create audio processor: %w", err)
	}

	return &VoiceClient{
		config:             config,
		vapiClient:         vapiClient,
		audioStream:        audioStream,
		signaling:          signaling,
		audioProcessor:     audioProcessor,
		lastSpeakerSamples: make([]float32, 0),
		callState: &CallState{
			Status: CallStatusIdle,
		},
		requestLog:  make(chan APIRequest, 100),
		responseLog: make(chan APIResponse, 100),
		callEvents:  make(chan CallEvent, 100),
	}, nil
}

// StartCall initiates a voice call with the specified assistant
func (c *VoiceClient) StartCall(assistantID string) error {
	c.callState.Status = CallStatusConnecting
	c.callState.AssistantID = assistantID
	c.callState.StartTime = time.Now()

	// 1. Create WebSocket call via Vapi's /call endpoint with WebSocket transport
	call, err := c.createVapiWebSocketCall(assistantID)
	if err != nil {
		c.callState.Status = CallStatusFailed
		return fmt.Errorf("failed to create Vapi WebSocket call: %w", err)
	}

	// Update call state from Vapi response
	c.callState.CallID = call.Id
	c.callState.WebSocketURL = call.RoomURL

	// 2. Connect to Vapi WebSocket transport
	if err := c.signaling.Connect(call.RoomURL); err != nil {
		c.callState.Status = CallStatusFailed
		return fmt.Errorf("failed to connect to WebSocket transport: %w", err)
	}

	// Start monitoring signaling events
	go c.handleSignalingEvents()

	// 3. Start audio stream
	if err := c.audioStream.Start(); err != nil {
		c.callState.Status = CallStatusFailed
		return fmt.Errorf("failed to start audio stream: %w", err)
	}

	// 4. Reset and start audio processing
	c.audioProcessor.Reset()

	// 5. Start streaming microphone audio to WebSocket
	go c.streamMicrophoneAudio()

	c.callState.Status = CallStatusConnected

	// Emit call started event
	c.callEvents <- CallEvent{
		Type:      "call_started",
		Data:      c.callState,
		Timestamp: time.Now(),
	}

	return nil
}

// WebSocketCallRequest represents the request structure for /call endpoint with WebSocket transport
type WebSocketCallRequest struct {
	AssistantID string `json:"assistantId"`
	Transport   struct {
		Provider    string `json:"provider"`
		AudioFormat struct {
			Format     string `json:"format"`
			Container  string `json:"container"`
			SampleRate int    `json:"sampleRate"`
		} `json:"audioFormat"`
	} `json:"transport"`
}

// WebSocketCallResponse represents the response from /call endpoint with WebSocket transport
type WebSocketCallResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	AssistantID string `json:"assistantId"`
	Transport   struct {
		Provider         string `json:"provider"`
		WebsocketCallURL string `json:"websocketCallUrl"` // The WebSocket URL for audio transport
	} `json:"transport"`
	CreatedAt time.Time `json:"createdAt"`
}

// Call represents a Vapi call for WebSocket transport
type Call struct {
	Id          string
	AssistantID string
	Status      string
	RoomURL     string
	RoomName    string
	JoinToken   string
	ListenURL   string // Vapi WebSocket for monitoring
	ControlURL  string // Vapi control URL
}

// createVapiWebSocketCall creates a WebSocket call via Vapi's /call endpoint with WebSocket transport
func (c *VoiceClient) createVapiWebSocketCall(assistantID string) (*Call, error) {
	// Prepare the request payload for WebSocket transport
	payload := WebSocketCallRequest{
		AssistantID: assistantID,
		Transport: struct {
			Provider    string `json:"provider"`
			AudioFormat struct {
				Format     string `json:"format"`
				Container  string `json:"container"`
				SampleRate int    `json:"sampleRate"`
			} `json:"audioFormat"`
		}{
			Provider: "vapi.websocket",
			AudioFormat: struct {
				Format     string `json:"format"`
				Container  string `json:"container"`
				SampleRate int    `json:"sampleRate"`
			}{
				Format:     "pcm_s16le",
				Container:  "raw",
				SampleRate: 16000, // Request 16kHz from Vapi (their default)
			},
		},
	}

	// Marshal the request payload
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal WebSocket call request: %w", err)
	}

	// Get the API base URL from config
	baseURL := c.config.getAPIBaseURL()
	url := baseURL + "/call"

	// Use private API key for call creation
	privateKey := c.config.getPrivateAPIKey()

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create WebSocket call request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+privateKey)

	// Log the API request
	requestLog := APIRequest{
		Method:    "POST",
		URL:       url,
		Headers:   map[string]string{"Authorization": "Bearer " + privateKey[:10] + "...", "Content-Type": "application/json"},
		Body:      payload,
		Timestamp: time.Now(),
	}
	select {
	case c.requestLog <- requestLog:
	default:
		// Channel full, drop log
	}

	// Make the HTTP request
	startTime := time.Now()
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create WebSocket call: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // Error handling would complicate deferred cleanup

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// Try to read error response body for more details
		var errorBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorBody); err == nil {
			// Log error response
			responseLog := APIResponse{
				StatusCode: resp.StatusCode,
				Headers:    make(map[string]string),
				Body:       errorBody,
				Duration:   time.Since(startTime),
				Timestamp:  time.Now(),
			}
			select {
			case c.responseLog <- responseLog:
			default:
				// Channel full, drop log
			}
			return nil, fmt.Errorf("WebSocket call creation failed with status %d: %v", resp.StatusCode, errorBody)
		}
		return nil, fmt.Errorf("WebSocket call creation failed with status: %d", resp.StatusCode)
	}

	// Read raw response to see the actual structure
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if err := resp.Body.Close(); err != nil {
		fmt.Printf("Failed to close response body: %v\n", err)
	}

	// Log successful response
	var responseBody map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		fmt.Printf("Failed to unmarshal response body: %v\n", err)
	}
	responseLog := APIResponse{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string]string),
		Body:       responseBody,
		Duration:   time.Since(startTime),
		Timestamp:  time.Now(),
	}
	select {
	case c.responseLog <- responseLog:
	default:
		// Channel full, drop log
	}

	// Parse the response
	var wsCallResp WebSocketCallResponse
	if err := json.Unmarshal(bodyBytes, &wsCallResp); err != nil {
		return nil, fmt.Errorf("failed to decode WebSocket call response: %w", err)
	}

	// Convert to our internal Call structure
	call := &Call{
		Id:          wsCallResp.ID,
		AssistantID: wsCallResp.AssistantID,
		Status:      wsCallResp.Status,
		RoomURL:     wsCallResp.Transport.WebsocketCallURL, // Use WebSocket URL as room URL
		RoomName:    wsCallResp.ID,                         // Use call ID as room name
		JoinToken:   "",                                    // No token needed for WebSocket transport
		ListenURL:   wsCallResp.Transport.WebsocketCallURL, // WebSocket URL for transport
		ControlURL:  "",                                    // No separate control URL for WebSocket transport
	}

	return call, nil
}

// endVapiCall sends a DELETE request to Vapi to properly end the call
func (c *VoiceClient) endVapiCall(callID string) error {
	// Get the API base URL from config
	baseURL := c.config.getAPIBaseURL()
	url := baseURL + "/call/" + callID

	// Use private API key for call termination
	privateKey := c.config.getPrivateAPIKey()

	// Create DELETE request
	req, err := http.NewRequest("DELETE", url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create end call request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+privateKey)

	// Log the API request
	requestLog := APIRequest{
		Method:    "DELETE",
		URL:       url,
		Headers:   map[string]string{"Authorization": "Bearer " + privateKey[:10] + "..."},
		Body:      nil,
		Timestamp: time.Now(),
	}
	select {
	case c.requestLog <- requestLog:
	default:
		// Channel full, drop log
	}

	// Make the HTTP request
	startTime := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send end call request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // Error handling would complicate deferred cleanup

	// Log response
	responseLog := APIResponse{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string]string),
		Body:       nil,
		Duration:   time.Since(startTime),
		Timestamp:  time.Now(),
	}
	select {
	case c.responseLog <- responseLog:
	default:
		// Channel full, drop log
	}

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("end call request failed with status: %d", resp.StatusCode)
	}

	return nil
}

// EndCall terminates the current voice call
func (c *VoiceClient) EndCall() error {
	if c.callState.Status == CallStatusIdle {
		return fmt.Errorf("no active call to end")
	}

	// Send DELETE request to Vapi to properly end the call
	if c.callState.CallID != "" {
		if err := c.endVapiCall(c.callState.CallID); err != nil {
			fmt.Printf("Warning: failed to end Vapi call: %v\n", err)
			// Continue with local cleanup even if API call fails
		}
	}

	// Stop audio stream
	if c.audioStream != nil {
		if err := c.audioStream.Stop(); err != nil {
			fmt.Printf("Warning: failed to stop audio stream: %v\n", err)
		}
	}

	// Close signaling connection
	if c.signaling != nil {
		if err := c.signaling.Close(); err != nil {
			fmt.Printf("Warning: failed to close signaling: %v\n", err)
		}
	}

	// Reset call state
	c.callState.Status = CallStatusIdle
	c.callState.WebSocketURL = ""

	// Emit call ended event
	c.callEvents <- CallEvent{
		Type:      "call_ended",
		Data:      c.callState,
		Timestamp: time.Now(),
	}

	return nil
}

// GetCallState returns the current call state
func (c *VoiceClient) GetCallState() *CallState {
	return c.callState
}

// GetRequestLog returns the API request log channel
func (c *VoiceClient) GetRequestLog() <-chan APIRequest {
	return c.requestLog
}

// GetResponseLog returns the API response log channel
func (c *VoiceClient) GetResponseLog() <-chan APIResponse {
	return c.responseLog
}

// GetCallEvents returns the call events channel
func (c *VoiceClient) GetCallEvents() <-chan CallEvent {
	return c.callEvents
}

// GetAudioLevels returns current input and output audio levels
func (c *VoiceClient) GetAudioLevels() (input, output float32) {
	if c.audioStream == nil {
		return 0.0, 0.0
	}

	return c.audioStream.GetInputLevel(), c.audioStream.GetOutputLevel()
}

// IsAudioRunning returns true if audio stream is active
func (c *VoiceClient) IsAudioRunning() bool {
	if c.audioStream == nil {
		return false
	}

	return c.audioStream.IsRunning()
}

// ResetAudioProcessor resets the audio processor's internal state
func (c *VoiceClient) ResetAudioProcessor() {
	if c.audioProcessor != nil {
		c.audioProcessor.Reset()
	}
}

// SetNoiseGateThreshold adjusts the noise gate sensitivity
func (c *VoiceClient) SetNoiseGateThreshold(threshold float32) {
	if c.audioProcessor != nil {
		c.audioProcessor.SetNoiseGateThreshold(threshold)
	}
}

// SetEchoLearningRate adjusts the echo cancellation learning rate
func (c *VoiceClient) SetEchoLearningRate(rate float32) {
	if c.audioProcessor != nil {
		c.audioProcessor.SetLearningRate(rate)
	}
}

// handleSignalingEvents processes events from Vapi WebSocket signaling
func (c *VoiceClient) handleSignalingEvents() {
	for event := range c.signaling.GetEvents() {
		// Skip noisy audio_data events from being logged
		if event.Type == "audio_data" {
			// Handle audio data directly without forwarding as call event
			if samples, ok := event.Data.([]float32); ok {
				// Debug: Check for clipping in incoming audio
				var clippedCount int
				var maxSample float32
				for _, s := range samples {
					if s < 0 {
						if s < maxSample {
							maxSample = -s
						}
					} else if s > maxSample {
						maxSample = s
					}
					if s > 1.0 || s < -1.0 {
						clippedCount++
					}
				}
				if clippedCount > 0 || maxSample > 0.95 {
					fmt.Printf("⚠️  Incoming Vapi audio: %d samples, %d clipped, peak=%.3f\n",
						len(samples), clippedCount, maxSample)
				}

				// Store speaker samples for echo cancellation
				c.lastSpeakerSamples = samples

				// Vapi sends 16kHz audio, we need to upsample to 48kHz
				// TODO: This simple 3x upsampling by repeating samples causes poor audio quality
				// Should use proper interpolation or resampling library
				upsampled := make([]float32, len(samples)*3)
				for i := 0; i < len(samples); i++ {
					// Repeat each sample 3 times for simple upsampling
					// This causes aliasing and distortion!
					upsampled[i*3] = samples[i]
					upsampled[i*3+1] = samples[i]
					upsampled[i*3+2] = samples[i]
				}
				c.audioStream.WriteAudio(upsampled)
			}
			continue
		}

		// Skip excessive logging events
		if event.Type == "model-output" || event.Type == "voice-input" {
			continue
		}

		// Forward other signaling events as call events (for logging)
		callEvent := CallEvent{
			Type:      "signaling_" + event.Type,
			Data:      event.Data,
			Timestamp: event.Timestamp,
		}

		select {
		case c.callEvents <- callEvent:
		default:
			// Channel full, drop event
		}

		// Handle specific signaling events
		switch event.Type {
		case "room_joined":
			c.callEvents <- CallEvent{
				Type:      "room_connected",
				Data:      "Successfully connected to Vapi WebSocket transport",
				Timestamp: time.Now(),
			}

		case "participant_joined":
			c.callEvents <- CallEvent{
				Type:      "participant_joined",
				Data:      event.Data,
				Timestamp: time.Now(),
			}

		case "speech-update":
			// Handle speech status updates
			c.callEvents <- CallEvent{
				Type:      "speech_update",
				Data:      event.Data,
				Timestamp: time.Now(),
			}

		case "transcript":
			// Handle transcript events
			c.callEvents <- CallEvent{
				Type:      "transcript",
				Data:      event.Data,
				Timestamp: time.Now(),
			}

		case "webrtc_error", "daily_error", "websocket_error":
			c.callEvents <- CallEvent{
				Type:      "connection_error",
				Data:      event.Data,
				Timestamp: time.Now(),
			}
		}
	}
}

// streamMicrophoneAudio continuously streams audio from microphone to Vapi WebSocket
func (c *VoiceClient) streamMicrophoneAudio() {
	// Buffer for audio samples
	// AudioStream uses 48kHz, but Vapi expects 16kHz
	const audioStreamSampleRate = 48000
	const vapiSampleRate = 16000
	const chunkDurationMs = 20
	const audioStreamSamplesPerChunk = (audioStreamSampleRate * chunkDurationMs) / 1000 // 960 samples at 48kHz
	const vapiSamplesPerChunk = (vapiSampleRate * chunkDurationMs) / 1000               // 320 samples at 16kHz

	audioBuffer := make([]float32, vapiSamplesPerChunk)

	for c.callState.Status == CallStatusConnected || c.callState.Status == CallStatusConnecting {
		// Read audio from microphone
		if c.audioStream.IsRunning() {
			// Get audio samples from input stream at 48kHz
			inputSamples := c.audioStream.ReadAudio(audioStreamSamplesPerChunk)
			if len(inputSamples) > 0 {
				// Downsample from 48kHz to 16kHz (take every 3rd sample)
				for i := 0; i < vapiSamplesPerChunk && i*3 < len(inputSamples); i++ {
					audioBuffer[i] = inputSamples[i*3]
				}

				// Apply audio processing (echo cancellation and noise reduction)
				processedAudio := c.audioProcessor.ProcessAudio(audioBuffer, c.lastSpeakerSamples)

				// Send processed audio to Vapi WebSocket
				if c.signaling != nil && c.signaling.IsConnected() {
					if err := c.signaling.SendAudioData(processedAudio); err != nil {
						fmt.Printf("Failed to send audio data: %v\n", err)
					}
				}
			}
		}

		// Sleep for chunk duration (20ms)
		time.Sleep(time.Duration(chunkDurationMs) * time.Millisecond)
	}
}
