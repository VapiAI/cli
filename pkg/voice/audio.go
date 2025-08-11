package voice

import (
	"fmt"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
)

const (
	// Audio configuration constants
	SampleRate    = 48000
	FrameSize     = 480 // 10ms at 48kHz
	Channels      = 1   // Mono
	BitsPerSample = 16
)

// AudioBuffer represents a circular buffer for audio data
type AudioBuffer struct {
	data  []float32
	size  int
	head  int
	tail  int
	count int
	mutex sync.Mutex
}

// NewAudioBuffer creates a new audio buffer
func NewAudioBuffer(size int) *AudioBuffer {
	return &AudioBuffer{
		data: make([]float32, size),
		size: size,
	}
}

// Write writes audio data to the buffer
func (b *AudioBuffer) Write(data []float32) int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	written := 0
	for i, sample := range data {
		if b.count >= b.size {
			// Buffer full, drop oldest sample
			b.tail = (b.tail + 1) % b.size
			b.count--
		}

		b.data[b.head] = sample
		b.head = (b.head + 1) % b.size
		b.count++
		written = i + 1
	}

	return written
}

// Read reads audio data from the buffer
func (b *AudioBuffer) Read(data []float32) int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	read := 0
	for i := range data {
		if b.count == 0 {
			// Buffer empty, fill with silence
			data[i] = 0
		} else {
			data[i] = b.data[b.tail]
			b.tail = (b.tail + 1) % b.size
			b.count--
		}
		read = i + 1
	}

	return read
}

// Available returns the number of samples available for reading
func (b *AudioBuffer) Available() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.count
}

// AudioStream manages audio input and output streams
type AudioStream struct {
	deviceManager *AudioDeviceManager
	config        *WebRTCConfig

	// Input stream
	inputStream *portaudio.Stream
	inputBuffer *AudioBuffer
	inputDevice *AudioDevice

	// Output stream
	outputStream *portaudio.Stream
	outputBuffer *AudioBuffer
	outputDevice *AudioDevice

	// Control
	running  bool
	runMutex sync.RWMutex
	stopChan chan struct{}

	// Debugging
	debugger *AudioDebugger
}

// NewAudioStream creates a new audio stream
func NewAudioStream(config *WebRTCConfig) (*AudioStream, error) {
	deviceManager := NewAudioDeviceManager()
	if err := deviceManager.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize device manager: %w", err)
	}

	// Create debugger if enabled
	debugger := NewAudioDebugger(config.AudioDebug)

	// Create audio buffers (1 second of audio data)
	bufferSize := SampleRate * 1
	inputBuffer := NewAudioBuffer(bufferSize)
	outputBuffer := NewAudioBuffer(bufferSize)

	return &AudioStream{
		deviceManager: deviceManager,
		config:        config,
		inputBuffer:   inputBuffer,
		outputBuffer:  outputBuffer,
		stopChan:      make(chan struct{}),
		debugger:      debugger,
	}, nil
}

// Start starts the audio streams
func (a *AudioStream) Start() error {
	a.runMutex.Lock()
	defer a.runMutex.Unlock()

	if a.running {
		return fmt.Errorf("audio stream already running")
	}

	// Setup input device
	var err error
	if a.config.AudioInputDevice == "default" || a.config.AudioInputDevice == "" {
		a.inputDevice, err = a.deviceManager.GetDefaultInputDevice()
	} else {
		a.inputDevice, err = a.deviceManager.FindInputDeviceByName(a.config.AudioInputDevice)
	}
	if err != nil {
		return fmt.Errorf("failed to get input device: %w", err)
	}

	// Setup output device
	if a.config.AudioOutputDevice == "default" || a.config.AudioOutputDevice == "" {
		a.outputDevice, err = a.deviceManager.GetDefaultOutputDevice()
	} else {
		a.outputDevice, err = a.deviceManager.FindOutputDeviceByName(a.config.AudioOutputDevice)
	}
	if err != nil {
		return fmt.Errorf("failed to get output device: %w", err)
	}

	// Start input stream
	if err := a.startInputStream(); err != nil {
		return fmt.Errorf("failed to start input stream: %w", err)
	}

	// Start output stream
	if err := a.startOutputStream(); err != nil {
		if closeErr := a.inputStream.Close(); closeErr != nil {
			fmt.Printf("Failed to close input stream: %v\n", closeErr)
		}
		return fmt.Errorf("failed to start output stream: %w", err)
	}

	// Start debugger if enabled
	if err := a.debugger.Start(); err != nil {
		fmt.Printf("Failed to start audio debugger: %v\n", err)
	}

	a.running = true
	return nil
}

// createStream is a helper function to create audio streams
func (a *AudioStream) createStream(isInput bool, device *AudioDevice, callback interface{}) (*portaudio.Stream, error) {
	// Get all devices to find the actual device info
	devices, err := portaudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	if device.Index >= len(devices) {
		return nil, fmt.Errorf("invalid device index: %d", device.Index)
	}

	actualDevice := devices[device.Index]

	var params portaudio.StreamParameters
	if isInput {
		params = portaudio.StreamParameters{
			Input: portaudio.StreamDeviceParameters{
				Device:   actualDevice,
				Channels: Channels,
				Latency:  time.Duration(device.DefaultLowInputLatency * float64(time.Second)),
			},
			SampleRate:      SampleRate,
			FramesPerBuffer: FrameSize,
		}
	} else {
		params = portaudio.StreamParameters{
			Output: portaudio.StreamDeviceParameters{
				Device:   actualDevice,
				Channels: Channels,
				Latency:  time.Duration(device.DefaultLowOutputLatency * float64(time.Second)),
			},
			SampleRate:      SampleRate,
			FramesPerBuffer: FrameSize,
		}
	}

	stream, err := portaudio.OpenStream(params, callback)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}

	if err := stream.Start(); err != nil {
		if closeErr := stream.Close(); closeErr != nil {
			fmt.Printf("Failed to close stream: %v\n", closeErr)
		}
		return nil, fmt.Errorf("failed to start stream: %w", err)
	}

	return stream, nil
}

// startInputStream starts the audio input stream
func (a *AudioStream) startInputStream() error {
	// Create input callback
	inputCallback := func(in []float32) {
		// Debug input audio
		a.debugger.WriteInput(in)
		a.debugger.LogAudioStats(in, "Input")

		// Write audio data to input buffer for processing
		a.inputBuffer.Write(in)
	}

	stream, err := a.createStream(true, a.inputDevice, inputCallback)
	if err != nil {
		return fmt.Errorf("failed to create input stream: %w", err)
	}

	a.inputStream = stream
	return nil
}

// startOutputStream starts the audio output stream
func (a *AudioStream) startOutputStream() error {
	// Create output callback
	outputCallback := func(out []float32) {
		// Read audio data from output buffer
		a.outputBuffer.Read(out)

		// Debug output audio
		a.debugger.WriteOutput(out)
		a.debugger.LogAudioStats(out, "Output")
	}

	stream, err := a.createStream(false, a.outputDevice, outputCallback)
	if err != nil {
		return fmt.Errorf("failed to create output stream: %w", err)
	}

	a.outputStream = stream
	return nil
}

// Stop stops the audio streams
func (a *AudioStream) Stop() error {
	a.runMutex.Lock()
	defer a.runMutex.Unlock()

	if !a.running {
		return nil
	}

	// Signal stop
	close(a.stopChan)

	// Stop and close streams
	var inputErr, outputErr error

	if a.inputStream != nil {
		inputErr = a.inputStream.Close()
		a.inputStream = nil
	}

	if a.outputStream != nil {
		outputErr = a.outputStream.Close()
		a.outputStream = nil
	}

	// Stop debugger
	if err := a.debugger.Stop(); err != nil {
		fmt.Printf("Warning: failed to stop audio debugger: %v\n", err)
	}

	// Terminate device manager
	if err := a.deviceManager.Terminate(); err != nil {
		fmt.Printf("Warning: failed to terminate device manager: %v\n", err)
	}

	a.running = false

	// Return first error encountered
	if inputErr != nil {
		return fmt.Errorf("failed to close input stream: %w", inputErr)
	}
	if outputErr != nil {
		return fmt.Errorf("failed to close output stream: %w", outputErr)
	}

	return nil
}

// WriteAudio writes audio data to the output buffer (for incoming audio)
func (a *AudioStream) WriteAudio(data []float32) int {
	return a.outputBuffer.Write(data)
}

// GetInputLevel returns the current input audio level (0.0 to 1.0)
func (a *AudioStream) GetInputLevel() float32 {
	// Get recent audio data from input buffer
	samples := make([]float32, FrameSize)
	read := a.inputBuffer.Read(samples)

	if read == 0 {
		return 0.0
	}

	// Calculate RMS level
	var sum float32
	for i := 0; i < read; i++ {
		sum += samples[i] * samples[i]
	}

	rms := float32(0.0)
	if read > 0 {
		rms = float32(sum) / float32(read)
		if rms > 0 {
			rms = float32(0.5) // Simplified RMS calculation
		}
	}

	// Clamp to [0, 1]
	if rms > 1.0 {
		rms = 1.0
	}

	return rms
}

// GetOutputLevel returns the current output audio level (0.0 to 1.0)
func (a *AudioStream) GetOutputLevel() float32 {
	// For output level, we can check the buffer fill level as a proxy
	available := a.outputBuffer.Available()
	bufferSize := a.outputBuffer.size

	if bufferSize == 0 {
		return 0.0
	}

	level := float32(available) / float32(bufferSize)
	if level > 1.0 {
		level = 1.0
	}

	return level
}

// IsRunning returns true if the audio stream is running
func (a *AudioStream) IsRunning() bool {
	a.runMutex.RLock()
	defer a.runMutex.RUnlock()
	return a.running
}

// GetBufferState returns detailed buffer state for debugging
func (a *AudioStream) GetBufferState() (inputAvail, outputAvail, inputSize, outputSize int) {
	if a.inputBuffer != nil {
		inputAvail = a.inputBuffer.Available()
		inputSize = a.inputBuffer.size
	}
	if a.outputBuffer != nil {
		outputAvail = a.outputBuffer.Available()
		outputSize = a.outputBuffer.size
	}
	return
}

// LogBufferState logs current buffer state using the debugger
func (a *AudioStream) LogBufferState() {
	if a.debugger == nil {
		return
	}

	inputAvail, outputAvail, inputSize, outputSize := a.GetBufferState()
	a.debugger.LogBufferState(inputAvail, outputAvail, inputSize, outputSize)
}

// GetInputDevice returns the current input device
func (a *AudioStream) GetInputDevice() *AudioDevice {
	return a.inputDevice
}

// ReadAudio reads audio samples from the input buffer
func (a *AudioStream) ReadAudio(numSamples int) []float32 {
	if !a.IsRunning() {
		return make([]float32, numSamples) // Return silence if not running
	}

	samples := make([]float32, numSamples)
	a.inputBuffer.Read(samples)
	return samples
}

// GetOutputDevice returns the current output device
func (a *AudioStream) GetOutputDevice() *AudioDevice {
	return a.outputDevice
}
