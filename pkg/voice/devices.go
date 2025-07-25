package voice

import (
	"fmt"
	"strings"

	"github.com/gordonklaus/portaudio"
)

// AudioDevice represents an audio input or output device
type AudioDevice struct {
	Index                int
	Name                 string
	MaxInputChannels     int
	MaxOutputChannels    int
	DefaultSampleRate    float64
	DefaultLowInputLatency   float64
	DefaultLowOutputLatency  float64
	IsDefault            bool
}

// AudioDeviceManager manages audio device enumeration and selection
type AudioDeviceManager struct {
	inputDevices  []AudioDevice
	outputDevices []AudioDevice
	initialized   bool
}

// NewAudioDeviceManager creates a new audio device manager
func NewAudioDeviceManager() *AudioDeviceManager {
	return &AudioDeviceManager{
		inputDevices:  make([]AudioDevice, 0),
		outputDevices: make([]AudioDevice, 0),
		initialized:   false,
	}
}

// Initialize initializes PortAudio and enumerates devices
func (m *AudioDeviceManager) Initialize() error {
	if m.initialized {
		return nil
	}

	// Initialize PortAudio
	if err := portaudio.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize PortAudio: %w", err)
	}

	// Enumerate devices
	if err := m.enumerateDevices(); err != nil {
		portaudio.Terminate()
		return fmt.Errorf("failed to enumerate audio devices: %w", err)
	}

	m.initialized = true
	return nil
}

// Terminate terminates PortAudio
func (m *AudioDeviceManager) Terminate() error {
	if !m.initialized {
		return nil
	}

	if err := portaudio.Terminate(); err != nil {
		return fmt.Errorf("failed to terminate PortAudio: %w", err)
	}

	m.initialized = false
	return nil
}

// enumerateDevices discovers all available audio devices
func (m *AudioDeviceManager) enumerateDevices() error {
	// Get default devices
	defaultInput, err := portaudio.DefaultInputDevice()
	if err != nil {
		// Default input device might not be available
		defaultInput = nil
	}

	defaultOutput, err := portaudio.DefaultOutputDevice()
	if err != nil {
		return fmt.Errorf("failed to get default output device: %w", err)
	}

	// Get all devices
	devices, err := portaudio.Devices()
	if err != nil {
		return fmt.Errorf("failed to get audio devices: %w", err)
	}

	// Clear existing device lists
	m.inputDevices = make([]AudioDevice, 0)
	m.outputDevices = make([]AudioDevice, 0)

	// Process each device
	for i, device := range devices {
		audioDevice := AudioDevice{
			Index:                   i,
			Name:                    device.Name,
			MaxInputChannels:        device.MaxInputChannels,
			MaxOutputChannels:       device.MaxOutputChannels,
			DefaultSampleRate:       device.DefaultSampleRate,
			DefaultLowInputLatency:  device.DefaultLowInputLatency.Seconds(),
			DefaultLowOutputLatency: device.DefaultLowOutputLatency.Seconds(),
			IsDefault:              false,
		}

		// Check if this is the default input device
		if defaultInput != nil && device == defaultInput {
			audioDevice.IsDefault = true
		}

		// Check if this is the default output device
		if device == defaultOutput {
			audioDevice.IsDefault = true
		}

		// Add to appropriate device list
		if device.MaxInputChannels > 0 {
			m.inputDevices = append(m.inputDevices, audioDevice)
		}
		if device.MaxOutputChannels > 0 {
			m.outputDevices = append(m.outputDevices, audioDevice)
		}
	}

	return nil
}

// GetInputDevices returns all available input devices
func (m *AudioDeviceManager) GetInputDevices() ([]AudioDevice, error) {
	if !m.initialized {
		if err := m.Initialize(); err != nil {
			return nil, err
		}
	}
	return m.inputDevices, nil
}

// GetOutputDevices returns all available output devices
func (m *AudioDeviceManager) GetOutputDevices() ([]AudioDevice, error) {
	if !m.initialized {
		if err := m.Initialize(); err != nil {
			return nil, err
		}
	}
	return m.outputDevices, nil
}

// GetDefaultInputDevice returns the default input device
func (m *AudioDeviceManager) GetDefaultInputDevice() (*AudioDevice, error) {
	devices, err := m.GetInputDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.IsDefault {
			return &device, nil
		}
	}

	// If no default found, return the first available input device
	if len(devices) > 0 {
		return &devices[0], nil
	}

	return nil, fmt.Errorf("no input devices available")
}

// GetDefaultOutputDevice returns the default output device
func (m *AudioDeviceManager) GetDefaultOutputDevice() (*AudioDevice, error) {
	devices, err := m.GetOutputDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.IsDefault {
			return &device, nil
		}
	}

	// If no default found, return the first available output device
	if len(devices) > 0 {
		return &devices[0], nil
	}

	return nil, fmt.Errorf("no output devices available")
}

// FindInputDeviceByName finds an input device by name (case-insensitive partial match)
func (m *AudioDeviceManager) FindInputDeviceByName(name string) (*AudioDevice, error) {
	devices, err := m.GetInputDevices()
	if err != nil {
		return nil, err
	}

	name = strings.ToLower(name)

	// First try exact match
	for _, device := range devices {
		if strings.ToLower(device.Name) == name {
			return &device, nil
		}
	}

	// Then try partial match
	for _, device := range devices {
		if strings.Contains(strings.ToLower(device.Name), name) {
			return &device, nil
		}
	}

	return nil, fmt.Errorf("input device not found: %s", name)
}

// FindOutputDeviceByName finds an output device by name (case-insensitive partial match)
func (m *AudioDeviceManager) FindOutputDeviceByName(name string) (*AudioDevice, error) {
	devices, err := m.GetOutputDevices()
	if err != nil {
		return nil, err
	}

	name = strings.ToLower(name)

	// First try exact match
	for _, device := range devices {
		if strings.ToLower(device.Name) == name {
			return &device, nil
		}
	}

	// Then try partial match
	for _, device := range devices {
		if strings.Contains(strings.ToLower(device.Name), name) {
			return &device, nil
		}
	}

	return nil, fmt.Errorf("output device not found: %s", name)
}

// ListDevices returns a formatted string listing all audio devices
func (m *AudioDeviceManager) ListDevices() (string, error) {
	if !m.initialized {
		if err := m.Initialize(); err != nil {
			return "", err
		}
	}

	var result strings.Builder

	result.WriteString("🎤 Input Devices:\n")
	for _, device := range m.inputDevices {
		defaultStr := ""
		if device.IsDefault {
			defaultStr = " (default)"
		}
		result.WriteString(fmt.Sprintf("  [%d] %s%s - %d channels, %.0f Hz\n",
			device.Index, device.Name, defaultStr, device.MaxInputChannels, device.DefaultSampleRate))
	}

	result.WriteString("\n🔊 Output Devices:\n")
	for _, device := range m.outputDevices {
		defaultStr := ""
		if device.IsDefault {
			defaultStr = " (default)"
		}
		result.WriteString(fmt.Sprintf("  [%d] %s%s - %d channels, %.0f Hz\n",
			device.Index, device.Name, defaultStr, device.MaxOutputChannels, device.DefaultSampleRate))
	}

	return result.String(), nil
}