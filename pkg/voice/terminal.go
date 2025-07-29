package voice

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// TerminalUI manages the terminal interface for voice calls
type TerminalUI struct {
	client     *VoiceClient
	done       chan bool
	keyEvents  chan rune
	uiUpdates  chan UIUpdate
	callEvents chan CallEvent

	// Styles
	successStyle lipgloss.Style
	errorStyle   lipgloss.Style
	infoStyle    lipgloss.Style
	headerStyle  lipgloss.Style
}

// UIUpdate represents a terminal UI update
type UIUpdate struct {
	Type string
	Data interface{}
}

// NewTerminalUI creates a new terminal UI manager
func NewTerminalUI(client *VoiceClient) *TerminalUI {
	return &TerminalUI{
		client:     client,
		done:       make(chan bool),
		keyEvents:  make(chan rune),
		uiUpdates:  make(chan UIUpdate),
		callEvents: make(chan CallEvent),

		// Initialize styles
		successStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true),
		errorStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true),
		infoStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF")),
		headerStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Bold(true),
	}
}

// Run starts the terminal UI
func (ui *TerminalUI) Run() error {
	// Display initial header
	ui.displayHeader()

	// Set up signal handling for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start event monitoring goroutines
	go ui.monitorCallEvents()
	go ui.handleKeyboardInput()

	// Main event loop
	for {
		select {
		case <-c:
			// Interrupt signal received
			fmt.Println(ui.infoStyle.Render("\nShutting down..."))
			return ui.shutdown()

		case event := <-ui.callEvents:
			ui.handleCallEvent(event)

		case update := <-ui.uiUpdates:
			ui.handleUIUpdate(update)

		case <-ui.done:
			return nil
		}
	}
}

// displayHeader shows the initial UI header
func (ui *TerminalUI) displayHeader() {
	fmt.Println(ui.headerStyle.Render("🚀 Vapi Voice Call"))
	fmt.Println()
	fmt.Println(ui.infoStyle.Render("Starting voice call..."))
	fmt.Println(ui.infoStyle.Render("Press Ctrl+C to end the call"))
	fmt.Println()
}

// monitorCallEvents monitors call events from the voice client
func (ui *TerminalUI) monitorCallEvents() {
	for event := range ui.client.GetCallEvents() {
		ui.callEvents <- event
	}
}

// handleCallEvent processes call events
func (ui *TerminalUI) handleCallEvent(event CallEvent) {
	timestamp := event.Timestamp.Format("15:04:05")

	switch event.Type {
	case "call_started":
		fmt.Printf("[%s] %s Call started successfully\n",
			timestamp, ui.successStyle.Render("✓"))
		ui.displayCallStatus()

	case "call_ended":
		fmt.Printf("[%s] %s Call ended\n",
			timestamp, ui.infoStyle.Render("•"))
		ui.done <- true

	case "ice_connection_state_change":
		state := event.Data.(string)
		fmt.Printf("[%s] %s Connection state: %s\n",
			timestamp, ui.infoStyle.Render("•"), state)

	case "ice_candidate":
		fmt.Printf("[%s] %s Connection negotiation\n",
			timestamp, ui.infoStyle.Render("•"))

	case "offer_sent":
		fmt.Printf("[%s] %s Audio connection established\n",
			timestamp, ui.infoStyle.Render("•"))

	case "room_connected":
		fmt.Printf("[%s] %s Connected to Vapi WebSocket transport\n",
			timestamp, ui.successStyle.Render("✓"))

	case "participant_joined":
		fmt.Printf("[%s] %s Participant joined call\n",
			timestamp, ui.successStyle.Render("✓"))

	case "connection_error":
		fmt.Printf("[%s] %s Connection error: %v\n",
			timestamp, ui.errorStyle.Render("✗"), event.Data)

	case "signaling_room_joined":
		fmt.Printf("[%s] %s Vapi WebSocket connected\n",
			timestamp, ui.successStyle.Render("✓"))

	default:
		// Show all events for debugging
		if event.Type != "" {
			fmt.Printf("[%s] %s %s\n",
				timestamp, ui.infoStyle.Render("•"), event.Type)
		}
	}
}

// handleUIUpdate processes UI updates
func (ui *TerminalUI) handleUIUpdate(update UIUpdate) {
	switch update.Type {
	case "status_update":
		ui.displayCallStatus()
	case "error":
		fmt.Printf("%s %v\n", ui.errorStyle.Render("✗"), update.Data)
	}
}

// displayCallStatus shows current call status
func (ui *TerminalUI) displayCallStatus() {
	state := ui.client.GetCallState()

	fmt.Println(ui.headerStyle.Render("📞 Call Status"))
	fmt.Printf("  Call ID: %s\n", state.CallID)
	fmt.Printf("  Assistant: %s\n", state.AssistantID)
	fmt.Printf("  Status: %s\n", ui.formatStatus(state.Status))
	fmt.Printf("  Duration: %s\n", ui.formatDuration(state.StartTime))

	if state.WebSocketURL != "" {
		fmt.Printf("  Room: %s\n", state.CallID)
		fmt.Printf("  WebSocket URL: %s\n", state.WebSocketURL)
	}

	// Display audio status
	if ui.client.IsAudioRunning() {
		fmt.Printf("  Audio: %s\n", ui.successStyle.Render("Active"))
	} else {
		fmt.Printf("  Audio: %s\n", ui.errorStyle.Render("Inactive"))
	}

	fmt.Println()
}

// formatStatus formats call status with appropriate colors
func (ui *TerminalUI) formatStatus(status CallStatus) string {
	switch status {
	case CallStatusConnected:
		return ui.successStyle.Render(string(status))
	case CallStatusFailed, CallStatusDisconnected:
		return ui.errorStyle.Render(string(status))
	case CallStatusIdle, CallStatusConnecting:
		return ui.infoStyle.Render(string(status))
	default:
		return ui.infoStyle.Render(string(status))
	}
}

// formatDuration formats call duration
func (ui *TerminalUI) formatDuration(startTime time.Time) string {
	if startTime.IsZero() {
		return "00:00:00"
	}

	duration := time.Since(startTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// handleKeyboardInput handles keyboard input (placeholder for future interactive features)
func (ui *TerminalUI) handleKeyboardInput() {
	// This is a placeholder for future keyboard input handling
	// For now, we rely on signal handling for termination
}

// shutdown gracefully shuts down the terminal UI
func (ui *TerminalUI) shutdown() error {
	fmt.Println(ui.infoStyle.Render("Ending voice call..."))

	// End the call if still active
	if ui.client.GetCallState().Status == CallStatusConnected {
		if err := ui.client.EndCall(); err != nil {
			fmt.Printf("%s Failed to end call: %v\n", ui.errorStyle.Render("✗"), err)
			// Don't return error, continue with shutdown
		}
	}

	// Give a brief moment for cleanup to complete
	time.Sleep(200 * time.Millisecond)

	fmt.Println(ui.successStyle.Render("✓ Voice call ended successfully"))

	// Force exit the process
	os.Exit(0)
	return nil // This line will never be reached, but Go requires it
}
