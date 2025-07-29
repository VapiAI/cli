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
	enabled      bool
	inputFile    *os.File
	outputFile   *os.File
	inputMutex   sync.Mutex
	outputMutex  sync.Mutex
	sampleRate   int
	channels     int
	bitsPerSample int
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
	inputFile, err := os.Create(inputPath)
	if err != nil {
		return fmt.Errorf("failed to create input debug file: %w", err)
	}
	d.inputFile = inputFile
	
	// Create output debug file
	outputPath := fmt.Sprintf("audio_debug_output_%s.wav", timestamp)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		inputFile.Close()
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
		binary.Write(d.inputFile, binary.LittleEndian, int16Sample)
	}
}

// WriteOutput writes output audio samples to debug file
func (d *AudioDebugger) WriteOutput(samples []float32) {
	if !d.enabled || d.outputFile == nil {
		return
	}
	
	d.outputMutex.Lock()
	defer d.outputMutex.Unlock()
	
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
		binary.Write(d.outputFile, binary.LittleEndian, int16Sample)
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
	binary.LittleEndian.PutUint32(header[28:32], uint32(byteRate))
	// Update BlockAlign
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
	file.Seek(4, 0)
	binary.Write(file, binary.LittleEndian, uint32(fileSize-8))
	
	// Update Subchunk2Size (file size - 44)
	file.Seek(40, 0)
	binary.Write(file, binary.LittleEndian, uint32(fileSize-44))
	
	return nil
}