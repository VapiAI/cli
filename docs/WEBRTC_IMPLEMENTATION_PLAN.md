# WebRTC Call Implementation Plan for Vapi CLI

## Overview
This document outlines the implementation plan for adding WebRTC calling functionality to the Vapi CLI using Pion WebRTC library and Daily.co as an intermediary service.

## Architecture

### High-Level Components
1. **CLI Command Interface** - New `vapi webrtc` command group
2. **WebRTC Client** - Pion-based WebRTC implementation 
3. **Daily.co Integration** - Room management and signaling via Daily.co API
4. **Vapi Integration** - Connect with existing Vapi assistant/call infrastructure
5. **Audio/Video Pipeline** - Handle media streams for voice/video calls

### Technology Stack
- **WebRTC Library**: Pion WebRTC v3 (github.com/pion/webrtc/v3)
- **Signaling Service**: Daily.co API
- **HTTP Client**: Standard Go net/http or existing client in codebase
- **Audio Processing**: Pion's built-in audio codecs (Opus, PCM)
- **CLI Framework**: Cobra (already in use)

## Debug Webhook System (Using Existing Infrastructure)

### Integration with Existing `vapi listen` Command
The WebRTC implementation will leverage the existing robust webhook infrastructure:

```go
type WebRTCDebugger struct {
    webhookURL   string        // URL for vapi listen forwarding
    events       chan WebhookEvent
    ui          *TerminalUI
    callID      string        // Track specific WebRTC call events
}

type WebhookEvent struct {
    Timestamp   time.Time     `json:"timestamp"`
    Type        string        `json:"type"`        // From existing webhook types
    CallID      string        `json:"call_id"`
    Data        interface{}   `json:"data"`
    SessionID   string        `json:"session_id"`
}
```

### Debug Integration Modes
1. **Auto-Start Listen Server** (`--debug`)
   - Automatically launches `vapi listen --forward-to localhost:3000/webhook`
   - Integrates webhook events into WebRTC terminal UI
   - Filters events by call ID for relevant debugging

2. **External Webhook Integration** (`--debug-webhook <url>`)
   - Uses existing webhook forwarding to external URL
   - Leverages existing authentication and retry logic
   - Maintains compatibility with current webhook tooling

3. **Existing File Logging**
   - Uses existing structured logging from `vapi listen`
   - Filters WebRTC-specific events for analysis

### Command Integration Examples
```bash
# WebRTC call with auto-debug (leverages existing listen command)
vapi call webrtc asst_12345 --debug
# Internally runs: vapi listen --forward-to localhost:3000/debug & 

# WebRTC call with external webhook (uses existing infrastructure)  
vapi call webrtc asst_12345 --debug-webhook http://localhost:8080/webhook

# WebRTC call with JSON config and debug
vapi call webrtc --config ./assistant.json --debug

# Manual setup using existing commands
vapi listen --forward-to localhost:3000/webhook &
vapi call webrtc asst_12345
```

### Terminal Flow Integration
```go
type TerminalUI struct {
    callStatus    *CallStatusView
    debugPanel    *DebugPanelView
    audioLevels   *AudioLevelsView
    controls      *ControlsView
}

// Real-time terminal layout
┌─ Call Status ────────────────────────────────────────────────────┐
│ 🟢 Connected to: Daily Room "test-call-1234"             │
│ 👤 Participants: You, Vapi Assistant                     │
│ ⏱️  Duration: 00:02:34                                    │
└──────────────────────────────────────────────────────────────────┘
┌─ Audio Levels ───────────────────────────────────────────────────┐
│ 🎤 Input:  ████████▒▒ 80%                               │
│ 🔊 Output: ██████▒▒▒▒ 60%                               │
└──────────────────────────────────────────────────────────────────┘
┌─ Debug Events ───────────────────────────────────────────────────┐
│ [14:23:45] POST /v1/calls → 201 Created                 │
│ [14:23:46] GET /v1/assistants/asst_123 → 200 OK         │
│ [14:23:47] WebSocket: connection established             │
│ [14:23:48] WebRTC: ICE candidate received               │
└──────────────────────────────────────────────────────────────────┘
┌─ Controls ───────────────────────────────────────────────────────┐
│ [m] Mute  [h] Hang up  [d] Toggle debug  [q] Quit       │
└──────────────────────────────────────────────────────────────────┘
```

## Terminal Flow Design

### Non-blocking Terminal UI
```go
type CallInterface struct {
    done        chan bool
    keyEvents   chan rune
    uiUpdates   chan UIUpdate
    callEvents  chan CallEvent
}

// Goroutine structure
func (c *CallInterface) Run() {
    go c.handleKeyInput()     // Non-blocking keyboard input
    go c.handleCallEvents()   // WebRTC/Vapi event processing
    go c.handleUIUpdates()    // Terminal display updates
    go c.handleWebhookEvents() // Debug webhook processing
    
    // Main event loop
    for {
        select {
        case key := <-c.keyEvents:
            c.handleKeyPress(key)
        case event := <-c.callEvents:
            c.updateCallStatus(event)
        case update := <-c.uiUpdates:
            c.refreshDisplay(update)
        case <-c.done:
            return
        }
    }
}
```

### Key Controls During Call
- `m`: Toggle mute/unmute
- `h`: Hang up call
- `d`: Toggle debug panel visibility
- `v`: Adjust volume levels
- `r`: Start/stop recording
- `t`: Show call transcript
- `q`: Quit (with confirmation)
- `↑/↓`: Scroll through debug events

### Terminal State Management
```go
type TerminalState struct {
    mode         DisplayMode  // Normal, Debug, Transcript
    callActive   bool
    muted        bool
    recording    bool
    debugVisible bool
    scrollPos    int
}
```

## Implementation Phases

### Phase 1: Core Infrastructure
1. **Add Dependencies**
   ```go
   // Add to go.mod
   github.com/pion/webrtc/v3 v3.x.x
   github.com/pion/interceptor v0.x.x
   ```

2. **Create WebRTC Package Structure**
   ```
   pkg/webrtc/
   ├── client.go          // Main WebRTC client
   ├── daily.go           // Daily.co API integration
   ├── signaling.go       // WebRTC signaling handling
   ├── media.go           // Audio/video stream management
   ├── config.go          // WebRTC configuration
   ├── audio.go           // PortAudio integration
   ├── devices.go         // Audio device management
   ├── api.go             // Vapi API request/response handling
   ├── terminal.go        // Terminal UI management
   └── diagnostics.go     // Connection diagnostics
   ```

3. **Daily.co Integration**
   ```go
   type DailyClient struct {
       apiKey    string
       domain    string
       httpClient *http.Client
   }
   
   type Room struct {
       Name        string                 `json:"name"`
       URL         string                 `json:"url"`
       Config      *RoomConfig           `json:"config,omitempty"`
       CreatedAt   time.Time             `json:"created_at"`
       Privacy     string                `json:"privacy"` // "public" | "private"
   }
   
   type RoomConfig struct {
       MaxParticipants int                `json:"max_participants"`
       EnableChat      bool               `json:"enable_chat"`
       EnableRecording bool               `json:"enable_recording"`
       AudioOnly       bool               `json:"audio_only"`
   }
   
   // Room management methods
   func (d *DailyClient) CreateRoom(name string, config *RoomConfig) (*Room, error)
   func (d *DailyClient) GetRoom(name string) (*Room, error)  
   func (d *DailyClient) DeleteRoom(name string) error
   func (d *DailyClient) GenerateToken(roomName string, props *TokenProperties) (string, error)
   ```
   
   **Authentication Flow:**
   1. Create room via Daily.co REST API with API key
   2. Generate meeting token for secure room access
   3. Connect to Daily.co WebSocket with token
   4. Handle room events and participant management

### Phase 2: CLI Commands
Add new command group under existing `call` command:

```
vapi call webrtc <assistant-id> [options]        // Start WebRTC call with assistant
vapi call webrtc --config <assistant-config.json> [options]  // Start with JSON config
vapi call webrtc configure                       // Configure audio devices
vapi call webrtc test-audio                     // Test microphone/speakers
vapi call webrtc status                          // Show current call status
vapi call webrtc end                             // End current WebRTC call
vapi call webrtc diagnostics                     // Connection diagnostics
```

**Primary Command Usage:**
```bash
# Start call with assistant ID
vapi call webrtc asst_12345 --debug-webhook http://localhost:3000/webhook

# Start call with JSON config
vapi call webrtc --config ./my-assistant.json --debug
```

**Command Flags:**
- `--room-name`: Custom room name (default: auto-generated)
- `--debug-webhook`: URL to receive debug request/response data
- `--debug`: Enable debug mode with local webhook server
- `--audio-input`: Specific audio input device
- `--audio-output`: Specific audio output device
- `--config`: Assistant configuration JSON file
- `--no-video`: Audio-only mode
- `--record`: Enable call recording

### Phase 3: WebRTC Implementation
1. **Peer Connection Setup**
   - Initialize Pion WebRTC peer connection
   - Configure ICE servers and STUN/TURN
   - Handle offer/answer exchange via Daily.co

2. **Media Handling**
   - Audio input/output (microphone/speakers)
   - Optional video support
   - Integration with Vapi's voice processing

3. **Signaling Protocol**
   - WebSocket connection to Daily.co
   - Handle ICE candidates exchange
   - Room state management

### Phase 4: Vapi Integration
1. **Assistant Connection**
   - Route audio to/from Vapi assistant
   - Handle call events and state changes
   - Integrate with existing Vapi call infrastructure

2. **Call Management**
   - Link WebRTC calls with Vapi call records
   - Transcript and recording integration
   - Billing and analytics

## File Structure Changes

### New Files to Create
```
cmd/webrtc.go                    // WebRTC CLI commands
pkg/webrtc/client.go            // Main WebRTC client
pkg/webrtc/daily.go             // Daily.co API client
pkg/webrtc/signaling.go         // WebRTC signaling
pkg/webrtc/media.go             // Media stream handling
pkg/webrtc/config.go            // Configuration
pkg/webrtc/audio.go             // PortAudio integration
pkg/webrtc/devices.go           // Audio device management
pkg/webrtc/api.go               // Vapi API request/response handling
pkg/webrtc/terminal.go          // Terminal UI management
pkg/webrtc/diagnostics.go       // Connection diagnostics
```

### Modified Files
```
cmd/call.go                     // Add WebRTC subcommands
go.mod                         // Add Pion WebRTC dependencies
```

## Dependencies

### Required Go Modules
```go
// Core WebRTC
github.com/pion/webrtc/v3 v3.2.40           // Core WebRTC implementation
github.com/pion/interceptor v0.1.25         // WebRTC interceptors
github.com/pion/opus v0.4.0                 // Opus audio codec
github.com/pion/rtp v1.8.2                  // RTP packet handling

// Audio System
github.com/gordonklaus/portaudio latest     // Cross-platform audio I/O
github.com/yourusername/go-audio latest     // Audio format conversion

// Networking
github.com/gorilla/websocket v1.5.1         // WebSocket for signaling

// Utilities
github.com/google/uuid v1.6.0               // Room ID generation
github.com/fatih/color v1.15.0              // Terminal colors for status
```

### Daily.co API Requirements
- Daily.co API key for room management
- WebSocket endpoint for real-time signaling
- REST API for room creation/management

## Configuration

### Environment Variables
```bash
DAILY_API_KEY=your_daily_api_key
DAILY_DOMAIN=your_daily_domain.daily.co
WEBRTC_STUN_SERVERS=stun:stun.l.google.com:19302
WEBRTC_TURN_SERVERS=turn:your-turn-server.com
WEBRTC_AUDIO_INPUT_DEVICE=default
WEBRTC_AUDIO_OUTPUT_DEVICE=default
```

### CLI Configuration
Extend existing config.go to include:
```go
type WebRTCConfig struct {
    DailyAPIKey    string `mapstructure:"daily_api_key"`
    DailyDomain    string `mapstructure:"daily_domain"`
    STUNServers    []string `mapstructure:"stun_servers"`
    TURNServers    []string `mapstructure:"turn_servers"`
    AudioCodec     string `mapstructure:"audio_codec"` // opus, pcm
    VideoEnabled   bool   `mapstructure:"video_enabled"`
}
```

## Implementation Steps

### Step 1: Setup and Dependencies
1. Add Pion WebRTC to go.mod
2. Create basic pkg/webrtc package structure
3. Implement Daily.co API client for room management

### Step 2: CLI Commands
1. Create cmd/webrtc.go with basic command structure
2. Implement room creation and joining commands
3. Add configuration handling for Daily.co credentials

### Step 3: WebRTC Core
1. Implement basic peer connection setup
2. Add signaling via Daily.co WebSocket
3. Handle offer/answer exchange and ICE candidates

### Step 4: Media Pipeline
1. **Audio Device Setup**
   ```go
   // Initialize PortAudio
   portaudio.Initialize()
   defer portaudio.Terminate()
   
   // Enumerate audio devices
   devices, err := portaudio.Devices()
   ```

2. **Audio Input Pipeline**
   ```go
   // Microphone -> PCM Buffer -> Opus Encoder -> WebRTC Track
   inputStream := setupAudioInput(selectedDevice)
   opusEncoder := opus.NewEncoder(48000, 1, opus.AppVoIP)
   audioTrack := setupWebRTCAudioTrack()
   ```

3. **Audio Output Pipeline**
   ```go
   // WebRTC Track -> Opus Decoder -> PCM Buffer -> Speakers
   outputStream := setupAudioOutput(selectedDevice)
   opusDecoder := opus.NewDecoder(48000, 1)
   ```

4. **Route audio to/from Vapi assistant**
   - Bidirectional audio stream routing
   - Real-time audio processing and forwarding

### Step 5: Integration and Testing
1. Connect WebRTC calls with Vapi call management
2. Add call state tracking and events
3. Test end-to-end call scenarios

## Security Considerations

1. **API Key Management**: Secure storage of Daily.co API keys
2. **Media Encryption**: Ensure DTLS/SRTP encryption is enabled
3. **Authentication**: Validate room access and user permissions
4. **Network Security**: Proper STUN/TURN server configuration

## Testing Strategy

1. **Unit Tests**: Individual component testing
2. **Integration Tests**: Daily.co API integration
3. **End-to-End Tests**: Full call scenarios
4. **Performance Tests**: Media quality and latency

## Potential Challenges

1. **Audio Routing**: Complex audio pipeline between WebRTC and Vapi
2. **NAT Traversal**: STUN/TURN server configuration
3. **Cross-Platform**: Audio device handling across different OS
4. **Error Handling**: Robust connection failure recovery
5. **Synchronization**: Managing call state between WebRTC and Vapi

## Success Metrics

1. Successful peer-to-peer connection establishment
2. Clear audio quality with low latency
3. Reliable connection through NAT/firewalls
4. Seamless integration with existing Vapi workflows
5. Proper call state management and recording

## Future Enhancements

1. **Video Support**: Add video calling capabilities
2. **Screen Sharing**: Implement screen sharing via WebRTC
3. **Multi-party Calls**: Support for conference calls
4. **Recording**: Direct WebRTC call recording
5. **Mobile Support**: Extend to mobile platforms via Go Mobile

## Resources

- [Pion WebRTC Documentation](https://pkg.go.dev/github.com/pion/webrtc/v3)
- [Daily.co API Documentation](https://docs.daily.co/)
- [WebRTC Standards](https://webrtc.org/)
- [Pion Examples](https://github.com/pion/example-webrtc-applications)