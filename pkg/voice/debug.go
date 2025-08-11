package voice

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"time"
)

// AudioDebugger handles audio debugging and recording
type AudioDebugger struct {
	enabled       bool
	inputFile     *os.File
	outputFile    *os.File
	inputMutex    sync.Mutex
	outputMutex   sync.Mutex
	sampleRate    int
	channels      int
	bitsPerSample int

	// Timing and flow tracking
	lastInputTime     time.Time
	lastOutputTime    time.Time
	outputSampleCount int64
	silentChunks      int
	totalChunks       int
}

// NewAudioDebugger creates a new audio debugger
func NewAudioDebugger(enabled bool) *AudioDebugger {
	return &AudioDebugger{
		enabled:       enabled,
		sampleRate:    48000, // Match your audio pipeline
		channels:      1,     // Mono
		bitsPerSample: 16,
	}
}

// Start initializes debug recording files
func (d *AudioDebugger) Start() error {
	if !d.enabled {
		return nil
	}

	timestamp := time.Now().Format("20060102-150405")

	// Create input debug file
	inputPath := fmt.Sprintf("audio_debug_input_%s.wav", timestamp)
	// #nosec G304 -- This is intentional file creation for debugging
	inputFile, err := os.Create(inputPath)
	if err != nil {
		return fmt.Errorf("failed to create input debug file: %w", err)
	}
	d.inputFile = inputFile

	// Create output debug file
	outputPath := fmt.Sprintf("audio_debug_output_%s.wav", timestamp)
	// #nosec G304 -- This is intentional file creation for debugging
	outputFile, err := os.Create(outputPath)
	if err != nil {
		if err := inputFile.Close(); err != nil {
			fmt.Printf("Failed to close input file: %v\n", err)
		}
		return fmt.Errorf("failed to create output debug file: %w", err)
	}
	d.outputFile = outputFile

	// Write WAV headers (we'll update the size later)
	if err := d.writeWAVHeader(d.inputFile); err != nil {
		return fmt.Errorf("failed to write input WAV header: %w", err)
	}
	if err := d.writeWAVHeader(d.outputFile); err != nil {
		return fmt.Errorf("failed to write output WAV header: %w", err)
	}

	fmt.Printf("📝 Audio debugging enabled:\n")
	fmt.Printf("   Input:  %s\n", inputPath)
	fmt.Printf("   Output: %s\n", outputPath)

	return nil
}

// WriteInput writes input audio samples to debug file
func (d *AudioDebugger) WriteInput(samples []float32) {
	if !d.enabled || d.inputFile == nil {
		return
	}

	d.inputMutex.Lock()
	defer d.inputMutex.Unlock()

	// Convert float32 to int16 and write
	for _, sample := range samples {
		// Check for clipping in float domain
		if sample > 1.0 || sample < -1.0 {
			fmt.Printf("⚠️  Input clipping detected: %.3f\n", sample)
		}

		// Clamp to prevent overflow
		if sample > 1.0 {
			sample = 1.0
		} else if sample < -1.0 {
			sample = -1.0
		}

		// Convert to int16
		int16Sample := int16(sample * 32767.0)
		if err := binary.Write(d.inputFile, binary.LittleEndian, int16Sample); err != nil {
			fmt.Printf("Failed to write input sample: %v\n", err)
		}
	}
}

// WriteOutput writes output audio samples to debug file
func (d *AudioDebugger) WriteOutput(samples []float32) {
	if !d.enabled || d.outputFile == nil {
		return
	}

	d.outputMutex.Lock()
	defer d.outputMutex.Unlock()

	// Track timing and detect gaps
	now := time.Now()
	if !d.lastOutputTime.IsZero() {
		timeSinceLastOutput := now.Sub(d.lastOutputTime)
		expectedInterval := time.Duration(float64(len(samples)) / float64(d.sampleRate) * float64(time.Second))

		// Detect significant gaps (more than 2x expected interval)
		if timeSinceLastOutput > expectedInterval*2 {
			fmt.Printf("🔇 OUTPUT GAP DETECTED: Expected %.2fms, got %.2fms (gap: %.2fms)\n",
				float64(expectedInterval.Nanoseconds())/1e6,
				float64(timeSinceLastOutput.Nanoseconds())/1e6,
				float64((timeSinceLastOutput-expectedInterval).Nanoseconds())/1e6)
		}
	}
	d.lastOutputTime = now
	d.outputSampleCount += int64(len(samples))

	// Check if this chunk is mostly silent
	var silentSamples int
	for _, sample := range samples {
		if sample > -0.001 && sample < 0.001 { // Very quiet threshold
			silentSamples++
		}
	}

	d.totalChunks++
	if float64(silentSamples)/float64(len(samples)) > 0.95 {
		d.silentChunks++
		if d.totalChunks%50 == 0 { // Log every 50 chunks
			fmt.Printf("🔇 Output silence rate: %d/%d chunks (%.1f%%) - Current chunk: %d/%d silent\n",
				d.silentChunks, d.totalChunks,
				float64(d.silentChunks)/float64(d.totalChunks)*100,
				silentSamples, len(samples))
		}
	}

	// Convert float32 to int16 and write
	for _, sample := range samples {
		// Check for clipping in float domain
		if sample > 1.0 || sample < -1.0 {
			fmt.Printf("⚠️  Output clipping detected: %.3f\n", sample)
		}

		// Clamp to prevent overflow
		if sample > 1.0 {
			sample = 1.0
		} else if sample < -1.0 {
			sample = -1.0
		}

		// Convert to int16
		int16Sample := int16(sample * 32767.0)
		if err := binary.Write(d.outputFile, binary.LittleEndian, int16Sample); err != nil {
			fmt.Printf("Failed to write output sample: %v\n", err)
		}
	}
}

// LogAudioStats logs statistics about audio samples
func (d *AudioDebugger) LogAudioStats(samples []float32, source string) {
	if !d.enabled || len(samples) == 0 {
		return
	}

	// Calculate RMS
	var sum float64
	var peak float32
	var clippedCount int

	for _, sample := range samples {
		sum += float64(sample * sample)

		absSample := sample
		if absSample < 0 {
			absSample = -absSample
		}

		if absSample > peak {
			peak = absSample
		}

		if sample > 1.0 || sample < -1.0 {
			clippedCount++
		}
	}

	rms := float32(sum / float64(len(samples)))

	if clippedCount > 0 || peak > 0.95 {
		fmt.Printf("🔊 %s Audio Stats: RMS=%.3f, Peak=%.3f, Clipped=%d/%d\n",
			source, rms, peak, clippedCount, len(samples))
	}
}

// Stop closes debug files and updates WAV headers
func (d *AudioDebugger) Stop() error {
	if !d.enabled {
		return nil
	}

	var errs []error

	if d.inputFile != nil {
		d.inputMutex.Lock()
		if err := d.updateWAVHeader(d.inputFile); err != nil {
			errs = append(errs, fmt.Errorf("failed to update input WAV header: %w", err))
		}
		if err := d.inputFile.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close input file: %w", err))
		}
		d.inputMutex.Unlock()
	}

	if d.outputFile != nil {
		d.outputMutex.Lock()
		if err := d.updateWAVHeader(d.outputFile); err != nil {
			errs = append(errs, fmt.Errorf("failed to update output WAV header: %w", err))
		}
		if err := d.outputFile.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close output file: %w", err))
		}
		d.outputMutex.Unlock()
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during stop: %v", errs)
	}

	fmt.Println("📝 Audio debug files saved")
	return nil
}

// writeWAVHeader writes a WAV file header
func (d *AudioDebugger) writeWAVHeader(file *os.File) error {
	// WAV header structure
	header := []byte{
		'R', 'I', 'F', 'F', // ChunkID
		0, 0, 0, 0, // ChunkSize (to be filled later)
		'W', 'A', 'V', 'E', // Format
		'f', 'm', 't', ' ', // Subchunk1ID
		16, 0, 0, 0, // Subchunk1Size (16 for PCM)
		1, 0, // AudioFormat (1 = PCM)
		byte(d.channels), byte(d.channels >> 8), // NumChannels
		byte(d.sampleRate), byte(d.sampleRate >> 8), byte(d.sampleRate >> 16), byte(d.sampleRate >> 24), // SampleRate
		0, 0, 0, 0, // ByteRate (to be calculated)
		0, 0, // BlockAlign (to be calculated)
		byte(d.bitsPerSample), byte(d.bitsPerSample >> 8), // BitsPerSample
		'd', 'a', 't', 'a', // Subchunk2ID
		0, 0, 0, 0, // Subchunk2Size (to be filled later)
	}

	// Calculate ByteRate and BlockAlign
	blockAlign := d.channels * d.bitsPerSample / 8
	byteRate := d.sampleRate * blockAlign

	// Update ByteRate
	// #nosec G115 -- byteRate is calculated from safe constants
	binary.LittleEndian.PutUint32(header[28:32], uint32(byteRate))
	// Update BlockAlign
	// #nosec G115 -- blockAlign is calculated from safe constants
	binary.LittleEndian.PutUint16(header[32:34], uint16(blockAlign))

	_, err := file.Write(header)
	return err
}

// updateWAVHeader updates the WAV header with the correct file size
func (d *AudioDebugger) updateWAVHeader(file *os.File) error {
	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()

	// Update ChunkSize (file size - 8)
	if _, err := file.Seek(4, 0); err != nil {
		return fmt.Errorf("failed to seek to chunk size position: %w", err)
	}
	// #nosec G115 -- fileSize is from file stat, safe for WAV header
	if err := binary.Write(file, binary.LittleEndian, uint32(fileSize-8)); err != nil {
		return fmt.Errorf("failed to write chunk size: %w", err)
	}

	// Update Subchunk2Size (file size - 44)
	if _, err := file.Seek(40, 0); err != nil {
		return fmt.Errorf("failed to seek to subchunk size position: %w", err)
	}
	// #nosec G115 -- fileSize is from file stat, safe for WAV header
	if err := binary.Write(file, binary.LittleEndian, uint32(fileSize-44)); err != nil {
		return fmt.Errorf("failed to write subchunk size: %w", err)
	}

	return nil
}

// LogWebSocketAudio logs detailed information about incoming WebSocket audio
func (d *AudioDebugger) LogWebSocketAudio(samples []float32, timestamp time.Time) {
	if !d.enabled || len(samples) == 0 {
		return
	}

	// Check for timing gaps in WebSocket audio
	if !d.lastInputTime.IsZero() {
		timeSinceLastWS := timestamp.Sub(d.lastInputTime)
		expectedInterval := time.Duration(float64(len(samples)) / 16000.0 * float64(time.Second)) // 16kHz from Vapi

		if timeSinceLastWS > expectedInterval*3 {
			fmt.Printf("🌐 WEBSOCKET AUDIO GAP: Expected %.2fms, got %.2fms (gap: %.2fms)\n",
				float64(expectedInterval.Nanoseconds())/1e6,
				float64(timeSinceLastWS.Nanoseconds())/1e6,
				float64((timeSinceLastWS-expectedInterval).Nanoseconds())/1e6)
		}
	}
	d.lastInputTime = timestamp

	// Analyze audio content
	var silentSamples, clippedSamples int
	var peak, rms float32
	for _, sample := range samples {
		if sample > -0.001 && sample < 0.001 {
			silentSamples++
		}
		if sample > 1.0 || sample < -1.0 {
			clippedSamples++
		}

		absSample := sample
		if absSample < 0 {
			absSample = -absSample
		}
		if absSample > peak {
			peak = absSample
		}
		rms += sample * sample
	}
	rms /= float32(len(samples))

	silenceRate := float64(silentSamples) / float64(len(samples))

	// Log if significant silence or other issues
	if silenceRate > 0.9 || clippedSamples > 0 || peak > 0.95 {
		fmt.Printf("🌐 WebSocket Audio: %d samples, %.1f%% silent, peak=%.3f, rms=%.3f, clipped=%d\n",
			len(samples), silenceRate*100, peak, rms, clippedSamples)
	}
}

// LogBufferState logs the current state of audio buffers
func (d *AudioDebugger) LogBufferState(inputAvailable, outputAvailable, inputSize, outputSize int) {
	if !d.enabled {
		return
	}

	inputFill := float64(inputAvailable) / float64(inputSize) * 100
	outputFill := float64(outputAvailable) / float64(outputSize) * 100

	// Log if buffers are getting too full or too empty
	if inputFill < 10 || inputFill > 90 || outputFill < 10 || outputFill > 90 {
		fmt.Printf("📊 Buffer State: Input %.1f%% (%d/%d), Output %.1f%% (%d/%d)\n",
			inputFill, inputAvailable, inputSize,
			outputFill, outputAvailable, outputSize)
	}

	// Warn about potential underruns
	if outputFill < 5 {
		fmt.Printf("⚠️  OUTPUT BUFFER UNDERRUN RISK: Only %.1f%% filled (%d/%d samples)\n",
			outputFill, outputAvailable, outputSize)
	}
}

// LogAudioFlow provides a comprehensive view of the audio pipeline state
func (d *AudioDebugger) LogAudioFlow(stage string, sampleCount int, timestamp time.Time) {
	if !d.enabled {
		return
	}

	fmt.Printf("🎵 Audio Flow [%s]: %d samples at %s\n",
		stage, sampleCount, timestamp.Format("15:04:05.000"))
}
