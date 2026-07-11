package voice

import (
	"fmt"
	"math"

	"github.com/gorilla/websocket"
)

// WebSocketAudioProcessor handles audio processing with basic echo cancellation
type WebSocketAudioProcessor struct {
	// Echo cancellation state
	echoBuffer     []float32
	echoBufferSize int
	adaptiveFilter []float32
	learningRate   float32
	// Noise gate parameters
	noiseGateThreshold float32
	gateRatio          float32
}

// AudioPacket represents the structure for WebSocket audio data
type AudioPacket struct {
	MicSamples     []float32 `json:"micSamples"`
	SpeakerSamples []float32 `json:"speakerSamples,omitempty"`
	Timestamp      int64     `json:"timestamp"`
}

// NewWebSocketAudioProcessor creates a new audio processor with basic echo cancellation
func NewWebSocketAudioProcessor() (*WebSocketAudioProcessor, error) {
	const echoBufferSizeMs = 200 // 200ms echo buffer
	const sampleRate = 16000     // 16kHz sample rate (Vapi's format)
	echoBufferSize := (sampleRate * echoBufferSizeMs) / 1000
	return &WebSocketAudioProcessor{
		echoBuffer:         make([]float32, echoBufferSize),
		echoBufferSize:     echoBufferSize,
		adaptiveFilter:     make([]float32, 128), // 128-tap adaptive filter
		learningRate:       0.01,
		noiseGateThreshold: 0.01, // -40dB noise gate
		gateRatio:          0.1,  // 10:1 ratio
	}, nil
}

// ProcessAudio applies basic echo cancellation and noise reduction
func (wap *WebSocketAudioProcessor) ProcessAudio(micInput, speakerOutput []float32) []float32 {
	if len(micInput) == 0 {
		return micInput
	}
	processed := make([]float32, len(micInput))
	copy(processed, micInput)
	// Apply basic echo cancellation if we have speaker output
	if len(speakerOutput) > 0 {
		processed = wap.applyEchoCancellation(processed, speakerOutput)
	}
	// Apply noise gate
	processed = wap.applyNoiseGate(processed)
	return processed
}

// applyEchoCancellation implements a basic adaptive echo cancellation algorithm
func (wap *WebSocketAudioProcessor) applyEchoCancellation(micInput, speakerOutput []float32) []float32 {
	result := make([]float32, len(micInput))
	for i, sample := range micInput {
		// Store speaker output in echo buffer (circular buffer)
		if len(speakerOutput) > i {
			bufferIdx := (i) % wap.echoBufferSize
			wap.echoBuffer[bufferIdx] = speakerOutput[i]
		}
		// Estimate echo using adaptive filter
		var echoEstimate float32
		filterLen := len(wap.adaptiveFilter)
		for j := 0; j < filterLen && j < wap.echoBufferSize; j++ {
			bufferIdx := (i - j + wap.echoBufferSize) % wap.echoBufferSize
			echoEstimate += wap.adaptiveFilter[j] * wap.echoBuffer[bufferIdx]
		}
		// Subtract estimated echo from microphone input
		result[i] = sample - echoEstimate
		// Update adaptive filter using LMS algorithm
		errorSignal := result[i]
		for j := 0; j < filterLen && j < wap.echoBufferSize; j++ {
			bufferIdx := (i - j + wap.echoBufferSize) % wap.echoBufferSize
			wap.adaptiveFilter[j] += wap.learningRate * errorSignal * wap.echoBuffer[bufferIdx]
		}
	}
	return result
}

// applyNoiseGate applies a simple noise gate to reduce background noise
func (wap *WebSocketAudioProcessor) applyNoiseGate(input []float32) []float32 {
	result := make([]float32, len(input))
	for i, sample := range input {
		amplitude := float32(math.Abs(float64(sample)))
		if amplitude > wap.noiseGateThreshold {
			// Above threshold - pass through
			result[i] = sample
		} else {
			// Below threshold - apply gate ratio
			result[i] = sample * wap.gateRatio
		}
	}
	return result
}

// HandleWebSocket processes WebSocket connections with audio processing
func (wap *WebSocketAudioProcessor) HandleWebSocket(ws *websocket.Conn) error {
	defer func() {
		if err := ws.Close(); err != nil {
			fmt.Printf("Failed to close websocket: %v\n", err)
		}
	}()
	for {
		// Read audio data from WebSocket
		var audioData AudioPacket
		err := ws.ReadJSON(&audioData)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return fmt.Errorf("websocket read error: %w", err)
			}
			break
		}
		// Process audio with echo cancellation and noise reduction
		processed := wap.ProcessAudio(audioData.MicSamples, audioData.SpeakerSamples)
		// Send processed audio back via WebSocket
		response := AudioPacket{
			MicSamples: processed,
			Timestamp:  audioData.Timestamp,
		}
		if err := ws.WriteJSON(response); err != nil {
			return fmt.Errorf("websocket write error: %w", err)
		}
	}
	return nil
}

// Reset clears the processor's internal state
func (wap *WebSocketAudioProcessor) Reset() {
	// Clear echo buffer
	for i := range wap.echoBuffer {
		wap.echoBuffer[i] = 0
	}
	// Reset adaptive filter
	for i := range wap.adaptiveFilter {
		wap.adaptiveFilter[i] = 0
	}
}

// SetNoiseGateThreshold adjusts the noise gate sensitivity
func (wap *WebSocketAudioProcessor) SetNoiseGateThreshold(threshold float32) {
	wap.noiseGateThreshold = threshold
}

// SetLearningRate adjusts the adaptive filter learning rate
func (wap *WebSocketAudioProcessor) SetLearningRate(rate float32) {
	wap.learningRate = rate
}
