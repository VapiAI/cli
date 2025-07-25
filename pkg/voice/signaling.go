package voice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// VapiWebSocket handles WebSocket communication with Vapi transport
type VapiWebSocket struct {
	conn       *websocket.Conn
	wsURL      string
	events     chan SignalingEvent
	
	// Control
	connected  bool
	mutex      sync.RWMutex
	done       chan struct{}
}

// SignalingEvent represents a signaling event
type SignalingEvent struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	From      string      `json:"from,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// WebSocket message types for Vapi transport
const (
	MSG_ROOM_JOINED         = "room-joined"
	MSG_ERROR               = "error"
)

// NewVapiWebSocket creates a new Vapi WebSocket client
func NewVapiWebSocket() *VapiWebSocket {
	return &VapiWebSocket{
		events: make(chan SignalingEvent, 100),
		done:   make(chan struct{}),
	}
}

// Connect connects to Vapi WebSocket transport
func (s *VapiWebSocket) Connect(wsURL string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.connected {
		return fmt.Errorf("already connected to WebSocket transport")
	}

	if wsURL == "" {
		return fmt.Errorf("WebSocket URL is required")
	}

	s.wsURL = wsURL
	
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second
	
	// Add authentication headers for Vapi WebSocket
	headers := http.Header{}
	
	conn, resp, err := dialer.Dial(wsURL, headers)
	if err != nil {
		if resp != nil {
			if body, readErr := io.ReadAll(resp.Body); readErr == nil {
				return fmt.Errorf("WebSocket handshake failed (status %d): %s", resp.StatusCode, string(body))
			}
		}
		return fmt.Errorf("failed to connect to Vapi WebSocket: %w", err)
	}

	s.conn = conn
	s.connected = true

	// Start message handling for Vapi transport events
	go s.handleMessages()

	return nil
}


// handleMessages processes incoming WebSocket messages
func (s *VapiWebSocket) handleMessages() {
	defer func() {
		if r := recover(); r != nil {
			// Panic recovery - websocket connection was closed
		}
		s.mutex.Lock()
		s.connected = false
		if s.conn != nil {
			s.conn.Close()
		}
		s.mutex.Unlock()
	}()

	for {
		select {
		case <-s.done:
			return
		default:
			// Read message from WebSocket (blocking)
			messageType, data, err := s.conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					// Normal closure
					return
				}
				
				// Send error event
				s.events <- SignalingEvent{
					Type:      "websocket_error",
					Data:      err.Error(),
					Timestamp: time.Now(),
				}
				return
			}

			if messageType == websocket.TextMessage {
				s.handleTextMessage(data)
			} else if messageType == websocket.BinaryMessage {
				s.handleBinaryMessage(data)
			}
		}
	}
}

// handleTextMessage processes JSON control messages from Vapi WebSocket transport
func (s *VapiWebSocket) handleTextMessage(data []byte) {
	
	var message map[string]interface{}
	if err := json.Unmarshal(data, &message); err != nil {
		s.events <- SignalingEvent{
			Type:      "parse_error",
			Data:      string(data),
			Timestamp: time.Now(),
		}
		return
	}

	// Vapi WebSocket transport messages
	// Common types: speech-update, transcript, function-call, hang, etc.
	msgType := "vapi_transport_event"
	if eventType, ok := message["type"].(string); ok {
		msgType = eventType
	}

	// Create signaling event for Vapi transport
	event := SignalingEvent{
		Type:      msgType,
		Data:      message,
		Timestamp: time.Now(),
	}


	// Send event to listeners
	select {
	case s.events <- event:
	default:
		// Channel full, drop event
	}
}

// handleBinaryMessage processes binary audio data from Vapi WebSocket transport
func (s *VapiWebSocket) handleBinaryMessage(data []byte) {
	
	// Binary data is PCM audio from the assistant
	// Convert to float32 samples for audio playback
	if len(data)%2 != 0 {
		return
	}
	
	// Convert PCM 16-bit little-endian to float32 samples
	samples := make([]float32, len(data)/2)
	for i := 0; i < len(samples); i++ {
		// Read 16-bit little-endian sample correctly
		low := uint16(data[i*2])
		high := uint16(data[i*2+1])
		sample := int16(low | (high << 8))
		// Convert to float32 (-1.0 to 1.0) with proper scaling
		samples[i] = float32(sample) / 32767.0
	}
	
	// Send audio samples to output stream via event
	s.events <- SignalingEvent{
		Type:      "audio_data",
		Data:      samples,
		Timestamp: time.Now(),
	}
}

// SendAudioData sends binary audio data to Vapi WebSocket transport
func (s *VapiWebSocket) SendAudioData(samples []float32) error {
	s.mutex.RLock()
	conn := s.conn
	connected := s.connected
	s.mutex.RUnlock()

	if !connected || conn == nil {
		return fmt.Errorf("not connected to WebSocket transport")
	}

	// Convert float32 samples to PCM 16-bit little-endian
	data := make([]byte, len(samples)*2)
	for i, sample := range samples {
		// Clamp to [-1.0, 1.0] and convert to int16
		if sample > 1.0 {
			sample = 1.0
		} else if sample < -1.0 {
			sample = -1.0
		}
		
		pcmSample := int16(sample * 32767.0)
		
		// Write as little-endian
		data[i*2] = byte(pcmSample & 0xFF)
		data[i*2+1] = byte((pcmSample >> 8) & 0xFF)
	}

	return conn.WriteMessage(websocket.BinaryMessage, data)
}


// sendMessage sends a message to the WebSocket connection
func (s *VapiWebSocket) sendMessage(message map[string]interface{}) error {
	s.mutex.RLock()
	conn := s.conn
	connected := s.connected
	s.mutex.RUnlock()

	if !connected || conn == nil {
		return fmt.Errorf("not connected to signaling server")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

// GetEvents returns the events channel
func (s *VapiWebSocket) GetEvents() <-chan SignalingEvent {
	return s.events
}

// IsConnected returns true if connected to the signaling server
func (s *VapiWebSocket) IsConnected() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.connected
}

// Close closes the signaling connection
func (s *VapiWebSocket) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.connected {
		return nil
	}

	// Set connected to false first to stop message reading
	s.connected = false

	// Close connection immediately to interrupt any blocking reads
	var err error
	if s.conn != nil {
		err = s.conn.Close()
		s.conn = nil
	}

	// Signal shutdown to handleMessage goroutine
	select {
	case <-s.done:
		// already closed
	default:
		close(s.done)
	}

	return err
}