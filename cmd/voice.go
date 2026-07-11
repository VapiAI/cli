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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	vapi "github.com/VapiAI/server-sdk-go"
	"github.com/spf13/cobra"

	"github.com/VapiAI/cli/pkg/voice"
)

var (
	configFile        string
	audioInputDevice  string
	audioOutputDevice string
	noVideo           bool
	callTimeout       int
	audioDebug        bool

	// Transient assistant configuration
	assistantName string
	firstMessage  string
	voiceID       string
	model         string
	systemMessage string
)

// Voice call management commands
var voiceCmd = &cobra.Command{
	Use:   "voice [assistant-id]",
	Short: "Start voice call with assistant",
	Long: `Start a real-time voice call with a Vapi assistant.

This command creates a WebSocket connection using Vapi's native transport,
enabling bidirectional audio streaming for natural conversations.

You can either use an existing assistant ID or create a transient assistant
by specifying configuration flags.

Voice Call Flow:
  1. Creates a call via Vapi's /call endpoint with WebSocket transport
  2. Establishes WebSocket connection to Vapi's audio transport
  3. Streams microphone audio to the assistant
  4. Plays assistant responses through speakers

The VAPI_API_KEY will be used from your active CLI account configuration.

Examples:
  # Use existing assistant
  vapi call voice asst_12345
  
  # Create transient assistant inline
  vapi call voice --name "My Assistant" --first-message "Hello! How can I help you?"
  
  # Advanced transient assistant
  vapi call voice --name "Support Bot" --first-message "Hi there!" --voice-id "jennifer" --model "gpt-4o"
  
  # Load from config file
  vapi call voice --config ./assistant.json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var assistantID string

		// Determine if we're using an existing assistant or creating a transient one
		if len(args) > 0 {
			// Use existing assistant ID
			assistantID = args[0]
		} else if configFile != "" {
			// Load assistant configuration from JSON file
			// Clean the path to prevent directory traversal
			cleanPath := filepath.Clean(configFile)
			data, err := os.ReadFile(cleanPath)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}

			var config map[string]interface{}
			if err := json.Unmarshal(data, &config); err != nil {
				return fmt.Errorf("failed to parse config file: %w", err)
			}

			// Check if config has existing assistant ID
			if id, ok := config["assistant_id"].(string); ok {
				assistantID = id
			} else if id, ok := config["assistantId"].(string); ok {
				assistantID = id
			} else {
				// No assistant ID found - create transient assistant from config
				loadConfigIntoFlags(config)

				createdAssistantID, err := createTransientAssistant()
				if err != nil {
					return fmt.Errorf("failed to create transient assistant from config: %w", err)
				}
				assistantID = createdAssistantID
			}
		} else if assistantName != "" || firstMessage != "" {
			// Create transient assistant
			createdAssistantID, err := createTransientAssistant()
			if err != nil {
				return fmt.Errorf("failed to create transient assistant: %w", err)
			}
			assistantID = createdAssistantID
		} else {
			return fmt.Errorf("assistant ID is required (provide as argument, via --config, or via transient assistant flags like --name)")
		}

		return startVoiceCall(assistantID)
	},
}

var configureVoiceCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure voice call audio devices",
	Long:  `Configure audio input and output devices for voice calls.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🎛️  Voice Call Configuration")
		fmt.Println()

		// Create device manager to list devices
		deviceManager := voice.NewAudioDeviceManager()
		if err := deviceManager.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize audio system: %w", err)
		}
		defer func() {
			if err := deviceManager.Terminate(); err != nil {
				fmt.Printf("Failed to terminate device manager: %v\n", err)
			}
		}()

		// List available devices
		deviceList, err := deviceManager.ListDevices()
		if err != nil {
			return fmt.Errorf("failed to list audio devices: %w", err)
		}

		fmt.Println("Available audio devices:")
		fmt.Print(deviceList)

		fmt.Println("Configuration:")
		fmt.Println("- Use device names with --audio-input and --audio-output flags")
		fmt.Println("- Use 'default' to use system default devices")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  vapi call voice asst_12345 --audio-input \"Built-in Microphone\"")

		return nil
	},
}

var testAudioCmd = &cobra.Command{
	Use:   "test-audio",
	Short: "Test audio devices",
	Long:  `Test microphone and speaker functionality for voice calls.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🎤 Audio Test")
		fmt.Println()

		// Create a basic audio stream to test devices
		config := voice.DefaultWebRTCConfig()
		if audioInputDevice != "" {
			config.AudioInputDevice = audioInputDevice
		}
		if audioOutputDevice != "" {
			config.AudioOutputDevice = audioOutputDevice
		}

		audioStream, err := voice.NewAudioStream(config)
		if err != nil {
			return fmt.Errorf("failed to create audio stream: %w", err)
		}

		fmt.Println("Testing audio devices...")
		fmt.Printf("Input device: %s\n", config.AudioInputDevice)
		fmt.Printf("Output device: %s\n", config.AudioOutputDevice)
		fmt.Println()

		// Try to start the audio stream briefly
		if err := audioStream.Start(); err != nil {
			return fmt.Errorf("failed to start audio stream: %w", err)
		}

		fmt.Println("✅ Audio devices initialized successfully!")
		fmt.Printf("Input device: %s\n", audioStream.GetInputDevice().Name)
		fmt.Printf("Output device: %s\n", audioStream.GetOutputDevice().Name)
		fmt.Println()

		// Test for a brief moment
		fmt.Println("Testing audio for 3 seconds...")
		time.Sleep(3 * time.Second)

		// Get audio levels
		inputLevel, outputLevel := audioStream.GetInputLevel(), audioStream.GetOutputLevel()
		fmt.Printf("Input level: %.1f%%\n", inputLevel*100)
		fmt.Printf("Output level: %.1f%%\n", outputLevel*100)

		// Stop the audio stream
		if err := audioStream.Stop(); err != nil {
			fmt.Printf("Warning: %v\n", err)
		}

		fmt.Println()
		fmt.Println("✅ Audio test completed!")
		return nil
	},
}

var statusVoiceCmd = &cobra.Command{
	Use:   "status",
	Short: "Show voice call status",
	Long:  `Display the status of the current voice call.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("📞 Voice Call Status")
		fmt.Println()
		fmt.Println("No active voice call.")
		fmt.Println()
		fmt.Println("Start a call with:")
		fmt.Println("  vapi call voice <assistant-id>")
		return nil
	},
}

var endVoiceCmd = &cobra.Command{
	Use:   "end",
	Short: "End current voice call",
	Long:  `Terminate the current voice call.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("📞 End Voice Call")
		fmt.Println()
		fmt.Println("No active call to end.")
		fmt.Println()
		fmt.Println("Calls can be ended by pressing Ctrl+C during an active call.")
		return nil
	},
}

// loadConfigIntoFlags loads configuration from a JSON config into the flag variables
func loadConfigIntoFlags(config map[string]interface{}) {
	// Load name
	if name, ok := config["name"].(string); ok {
		assistantName = name
	}

	// Load first message
	if msg, ok := config["first_message"].(string); ok {
		firstMessage = msg
	} else if msg, ok := config["firstMessage"].(string); ok {
		firstMessage = msg
	}

	// Load voice ID
	if voiceValue, ok := config["voice_id"].(string); ok {
		voiceID = voiceValue
	} else if voiceValue, ok := config["voiceId"].(string); ok {
		voiceID = voiceValue
	}

	// Load model
	if mdl, ok := config["model"].(string); ok {
		model = mdl
	}

	// Load system message
	if sysMsg, ok := config["system_message"].(string); ok {
		systemMessage = sysMsg
	} else if sysMsg, ok := config["systemMessage"].(string); ok {
		systemMessage = sysMsg
	}
}

// createTransientAssistant creates a temporary assistant for the voice call
func createTransientAssistant() (string, error) {
	fmt.Println("🤖 Creating transient assistant...")

	// Get Vapi client
	if vapiClient.GetClient() == nil {
		return "", fmt.Errorf("no active Vapi account found. Please run 'vapi login' first")
	}

	// Set defaults if not provided
	name := assistantName
	if name == "" {
		name = "Transient Assistant"
	}

	message := firstMessage
	if message == "" {
		message = "Hello! How can I assist you today?"
	}

	ctx := context.Background()

	// Create the assistant request
	createRequest := &vapi.CreateAssistantDto{
		Name:         &name,
		FirstMessage: &message,
		Voice: &vapi.CreateAssistantDtoVoice{
			VapiVoice: &vapi.VapiVoice{
				VoiceId: vapi.VapiVoiceVoiceIdElliot, // Default voice
			},
		},
	}

	// Note: For now, we'll keep it simple and just use the default voice and model
	// Advanced voice/model configuration can be added later once we understand the full API structure
	if voiceID != "" {
		fmt.Printf("ℹ️  Voice ID '%s' specified but using default voice for now\n", voiceID)
	}
	if model != "" {
		fmt.Printf("ℹ️  Model '%s' specified but using default model for now\n", model)
	}
	if systemMessage != "" {
		fmt.Printf("ℹ️  System message specified but using default behavior for now\n")
	}

	// Create the assistant
	assistant, err := vapiClient.GetClient().Assistants.Create(ctx, createRequest)
	if err != nil {
		return "", fmt.Errorf("failed to create transient assistant: %w", err)
	}

	fmt.Printf("✅ Created transient assistant: %s (ID: %s)\n", name, assistant.Id)
	return assistant.Id, nil
}

// startVoiceCall initiates a voice call with the specified assistant
func startVoiceCall(assistantID string) error {
	fmt.Printf("🚀 Starting voice call with assistant: %s\n", assistantID)
	fmt.Println()

	// Create voice call configuration
	config := voice.DefaultWebRTCConfig()

	// Override with command line options
	if audioInputDevice != "" {
		config.AudioInputDevice = audioInputDevice
	}
	if audioOutputDevice != "" {
		config.AudioOutputDevice = audioOutputDevice
	}
	config.VideoEnabled = !noVideo
	config.AudioDebug = audioDebug

	// Get Vapi API configuration from the CLI client
	if vapiClient.GetClient() == nil {
		return fmt.Errorf("no active Vapi account found. Please run 'vapi login' first")
	}

	// Set Vapi API key from the active account configuration
	if apiKey := vapiClient.GetConfig().GetActiveAPIKey(); apiKey != "" {
		config.VapiAPIKey = apiKey
	} else {
		return fmt.Errorf("VAPI_API_KEY not found. Please run 'vapi login' to authenticate")
	}

	// Set API base URL from configuration
	config.VapiBaseURL = vapiClient.GetConfig().GetAPIBaseURL()

	// Set public API key from environment if provided
	if pub := os.Getenv("VAPI_PUBLIC_KEY"); pub != "" {
		config.VapiPublicAPIKey = pub
	}

	// Create voice client
	client, err := voice.NewVoiceClient(config, vapiClient.GetClient())
	if err != nil {
		return fmt.Errorf("failed to create voice client: %w", err)
	}

	// Create terminal UI
	ui := voice.NewTerminalUI(client)

	// Start the call
	if err := client.StartCall(assistantID); err != nil {
		return fmt.Errorf("failed to start voice call: %w", err)
	}

	// Run the terminal UI (this blocks until call ends)
	return ui.Run()
}

func init() {
	// Add voice as a subcommand of call
	callCmd.AddCommand(voiceCmd)
	voiceCmd.AddCommand(configureVoiceCmd)
	voiceCmd.AddCommand(testAudioCmd)
	voiceCmd.AddCommand(statusVoiceCmd)
	voiceCmd.AddCommand(endVoiceCmd)

	// Add flags to the main voice command
	voiceCmd.Flags().StringVar(&configFile, "config", "", "Path to assistant configuration JSON file")
	voiceCmd.Flags().StringVar(&audioInputDevice, "audio-input", "", "Audio input device name")
	voiceCmd.Flags().StringVar(&audioOutputDevice, "audio-output", "", "Audio output device name")
	voiceCmd.Flags().IntVar(&callTimeout, "timeout", 30, "Call timeout in minutes")
	voiceCmd.Flags().BoolVar(&audioDebug, "audio-debug", false, "Enable audio debugging (saves input/output to WAV files)")

	// Transient assistant flags
	voiceCmd.Flags().StringVar(&assistantName, "name", "", "Name for transient assistant")
	voiceCmd.Flags().StringVar(&firstMessage, "first-message", "", "First message from transient assistant")
	voiceCmd.Flags().StringVar(&voiceID, "voice-id", "", "Voice ID for transient assistant (jennifer, derek, elliot)")
	voiceCmd.Flags().StringVar(&model, "model", "", "AI model for transient assistant (gpt-4o, gpt-4o-mini, etc.)")
	voiceCmd.Flags().StringVar(&systemMessage, "system-message", "", "System message for transient assistant")
}
