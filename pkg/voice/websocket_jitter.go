package voice

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// WebSocketJitterBuffer provides adaptive jitter buffering for WebSocket audio
type WebSocketJitterBuffer struct {
	// Configuration
	targetDelay time.Duration
	maxDelay    time.Duration
	minDelay    time.Duration
	sampleRate  int

	// Buffer management
	audioBuffer [][]float32
	bufferMutex sync.RWMutex

	// Timing control
	lastWriteTime time.Time
	lastReadTime  time.Time
	readInterval  time.Duration

	// Adaptive delay
	currentDelay time.Duration
	delayMutex   sync.RWMutex

	// Control
	running  bool
	runMutex sync.RWMutex

	// Statistics
	packetsReceived int64
	packetsDropped  int64
	underruns       int64
	overruns        int64

	// Read ticker for consistent output
	ticker     *time.Ticker
	outputChan chan []float32
}

// WebSocketJitterConfig holds configuration for WebSocket jitter buffer
type WebSocketJitterConfig struct {
	SampleRate     int           // Audio sample rate (16000 for Vapi)
	MinDelay       time.Duration // Minimum buffer delay
	MaxDelay       time.Duration // Maximum buffer delay
	TargetDelay    time.Duration // Initial target delay
	PacketInterval time.Duration // Expected packet interval (20ms for Vapi)
}

// DefaultWebSocketJitterConfig returns optimized config for Vapi WebSocket
func DefaultWebSocketJitterConfig() *WebSocketJitterConfig {
	return &WebSocketJitterConfig{
		SampleRate:     16000,
		MinDelay:       40 * time.Millisecond,  // Minimum 40ms buffering
		MaxDelay:       200 * time.Millisecond, // Maximum 200ms buffering
		TargetDelay:    80 * time.Millisecond,  // Target 80ms - good for voice
		PacketInterval: 20 * time.Millisecond,  // Vapi sends 20ms packets
	}
}

// NewWebSocketJitterBuffer creates a new WebSocket-compatible jitter buffer
func NewWebSocketJitterBuffer(config *WebSocketJitterConfig) (*WebSocketJitterBuffer, error) {
	if config == nil {
		config = DefaultWebSocketJitterConfig()
	}

	jb := &WebSocketJitterBuffer{
		targetDelay:  config.TargetDelay,
		maxDelay:     config.MaxDelay,
		minDelay:     config.MinDelay,
		sampleRate:   config.SampleRate,
		currentDelay: config.TargetDelay,
		readInterval: config.PacketInterval,
		audioBuffer:  make([][]float32, 0, 50), // Pre-allocate for ~1 second
		outputChan:   make(chan []float32, 10),
	}

	return jb, nil
}

// WriteAudio adds audio samples to the jitter buffer
func (jb *WebSocketJitterBuffer) WriteAudio(samples []float32) error {
	if !jb.IsRunning() {
		return fmt.Errorf("jitter buffer not running")
	}

	now := time.Now()

	jb.bufferMutex.Lock()
	defer jb.bufferMutex.Unlock()

	// Copy samples to avoid any reference issues
	sampleCopy := make([]float32, len(samples))
	copy(sampleCopy, samples)

	// Add to buffer
	jb.audioBuffer = append(jb.audioBuffer, sampleCopy)
	jb.packetsReceived++
	jb.lastWriteTime = now

	// Check for buffer overflow
	maxBufferSize := int(jb.maxDelay / jb.readInterval)
	if len(jb.audioBuffer) > maxBufferSize {
		// Drop oldest packet
		jb.audioBuffer = jb.audioBuffer[1:]
		jb.overruns++
		if jb.overruns%25 == 0 {
			log.Printf("⚠️ Jitter buffer overrun: dropped oldest packet (total: %d)", jb.overruns)
		}
	}

	// Adaptive delay adjustment based on buffer fill
	jb.adjustDelay()

	return nil
}

// adjustDelay adapts the buffer delay based on current conditions
func (jb *WebSocketJitterBuffer) adjustDelay() {
	bufferSize := len(jb.audioBuffer)
	targetBufferSize := int(jb.targetDelay / jb.readInterval)

	jb.delayMutex.Lock()
	defer jb.delayMutex.Unlock()

	// Adjust target delay based on buffer fill
	if bufferSize < targetBufferSize/2 {
		// Buffer running low - increase delay slightly
		jb.currentDelay += 5 * time.Millisecond
		if jb.currentDelay > jb.maxDelay {
			jb.currentDelay = jb.maxDelay
		}
	} else if bufferSize > targetBufferSize*2 {
		// Buffer getting too full - decrease delay slightly
		jb.currentDelay -= 5 * time.Millisecond
		if jb.currentDelay < jb.minDelay {
			jb.currentDelay = jb.minDelay
		}
	}
}

// ReadAudio reads processed audio samples from the jitter buffer
func (jb *WebSocketJitterBuffer) ReadAudio(numSamples int) []float32 {
	if !jb.IsRunning() {
		return make([]float32, numSamples) // Return silence
	}

	// Try to get samples from output channel with timeout
	select {
	case samples := <-jb.outputChan:
		// Resize to requested length if needed
		if len(samples) == numSamples {
			return samples
		}

		result := make([]float32, numSamples)
		if len(samples) > 0 {
			copy(result, samples)
		}
		return result

	case <-time.After(10 * time.Millisecond):
		// Timeout - return silence to prevent blocking
		jb.underruns++
		if jb.underruns%50 == 0 {
			log.Printf("⚠️ Jitter buffer underrun: no data available (total: %d)", jb.underruns)
		}
		return make([]float32, numSamples)
	}
}

// Start begins jitter buffer operation
func (jb *WebSocketJitterBuffer) Start() error {
	jb.runMutex.Lock()
	defer jb.runMutex.Unlock()

	if jb.running {
		return fmt.Errorf("jitter buffer already running")
	}

	jb.running = true

	// Start read ticker for consistent output timing
	jb.ticker = time.NewTicker(jb.readInterval)
	go jb.readLoop()

	// Start stats monitoring
	go jb.monitorStats()

	log.Printf("🎵 WebSocket Jitter Buffer started (target delay: %v, interval: %v)",
		jb.targetDelay, jb.readInterval)
	return nil
}

// readLoop continuously reads from buffer and outputs at regular intervals
func (jb *WebSocketJitterBuffer) readLoop() {
	defer jb.ticker.Stop()

	initialDelay := jb.currentDelay
	log.Printf("🎵 Jitter buffer starting with %v initial delay", initialDelay)

	// Initial delay before starting to read
	time.Sleep(initialDelay)

	for range jb.ticker.C {
		if !jb.IsRunning() {
			return
		}

		jb.bufferMutex.RLock()
		bufferLen := len(jb.audioBuffer)

		if bufferLen > 0 {
			// Get the oldest packet
			samples := jb.audioBuffer[0]

			// Remove from buffer
			jb.bufferMutex.RUnlock()
			jb.bufferMutex.Lock()
			if len(jb.audioBuffer) > 0 {
				jb.audioBuffer = jb.audioBuffer[1:]
			}
			jb.bufferMutex.Unlock()

			// Send to output channel (non-blocking)
			select {
			case jb.outputChan <- samples:
				jb.lastReadTime = time.Now()
			default:
				// Output channel full - drop this packet
				jb.packetsDropped++
			}
		} else {
			jb.bufferMutex.RUnlock()
			// Buffer empty - output silence
			silence := make([]float32, 320) // 20ms at 16kHz
			select {
			case jb.outputChan <- silence:
			default:
				// Output channel full - just skip
			}
		}
	}
}

// monitorStats logs periodic statistics
func (jb *WebSocketJitterBuffer) monitorStats() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !jb.IsRunning() {
			return
		}
		jb.logStats()
	}
}

// logStats logs current buffer statistics
func (jb *WebSocketJitterBuffer) logStats() {
	jb.bufferMutex.RLock()
	bufferSize := len(jb.audioBuffer)
	jb.bufferMutex.RUnlock()

	jb.delayMutex.RLock()
	currentDelay := jb.currentDelay
	jb.delayMutex.RUnlock()

	outputQueueSize := len(jb.outputChan)

	log.Printf("📊 WebSocket Jitter Buffer Stats: Buffer: %d packets, Delay: %v, Output queue: %d/10, Received: %d, Dropped: %d, Underruns: %d, Overruns: %d",
		bufferSize, currentDelay, outputQueueSize, jb.packetsReceived, jb.packetsDropped, jb.underruns, jb.overruns)
}

// Stop stops the jitter buffer
func (jb *WebSocketJitterBuffer) Stop() error {
	jb.runMutex.Lock()
	defer jb.runMutex.Unlock()

	if !jb.running {
		return nil
	}

	jb.running = false

	if jb.ticker != nil {
		jb.ticker.Stop()
	}

	log.Printf("🎵 WebSocket Jitter Buffer stopped")
	return nil
}

// IsRunning returns true if the jitter buffer is running
func (jb *WebSocketJitterBuffer) IsRunning() bool {
	jb.runMutex.RLock()
	defer jb.runMutex.RUnlock()
	return jb.running
}

// GetStats returns current jitter buffer statistics
func (jb *WebSocketJitterBuffer) GetStats() map[string]interface{} {
	jb.bufferMutex.RLock()
	bufferSize := len(jb.audioBuffer)
	jb.bufferMutex.RUnlock()

	jb.delayMutex.RLock()
	currentDelay := jb.currentDelay
	jb.delayMutex.RUnlock()

	return map[string]interface{}{
		"buffer_size":       bufferSize,
		"current_delay_ms":  currentDelay.Milliseconds(),
		"target_delay_ms":   jb.targetDelay.Milliseconds(),
		"packets_received":  jb.packetsReceived,
		"packets_dropped":   jb.packetsDropped,
		"underruns":         jb.underruns,
		"overruns":          jb.overruns,
		"output_queue_size": len(jb.outputChan),
		"running":           jb.IsRunning(),
	}
}
