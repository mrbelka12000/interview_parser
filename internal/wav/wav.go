package wav

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

type (
	// Writer header structure for PCM audio
	Writer struct {
		sampleRate    uint32
		channels      uint32
		bitsPerSample uint16
	}

	WAVHeader struct {
		ChunkID       [4]byte // "RIFF"
		ChunkSize     uint32
		Format        [4]byte // "WAVE"
		Subchunk1ID   [4]byte // "fmt "
		Subchunk1Size uint32  // 16 for PCM
		AudioFormat   uint16  // 1 for PCM
		NumChannels   uint16
		SampleRate    uint32
		ByteRate      uint32
		BlockAlign    uint16
		BitsPerSample uint16
		Subchunk2ID   [4]byte // "data"
		Subchunk2Size uint32
	}
)

func NewWriter(sampleRate, channels uint32, bitsPerSample uint16) *Writer {
	return &Writer{
		sampleRate:    sampleRate,
		channels:      channels,
		bitsPerSample: bitsPerSample,
	}
}

// SaveAsWAV saves the captured audio as a WAV file
func (w *Writer) SaveAsWAV(filename string, capturedData []byte) error {
	if len(capturedData) == 0 {
		return fmt.Errorf("no audio data to save")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Calculate sizes
	dataSize := uint32(len(capturedData))
	headerSize := uint32(44)
	totalSize := headerSize + dataSize

	// Create WAV header
	header := WAVHeader{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     totalSize - 8, // Total size - 8 bytes for ChunkID and ChunkSize
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16, // PCM format
		AudioFormat:   1,  // PCM
		NumChannels:   uint16(w.channels),
		SampleRate:    w.sampleRate,
		ByteRate:      w.sampleRate * uint32(w.channels) * uint32(w.bitsPerSample) / 8,
		BlockAlign:    uint16(w.channels) * w.bitsPerSample / 8,
		BitsPerSample: w.bitsPerSample,
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: dataSize,
	}

	// Write header
	err = binary.Write(file, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write audio data
	_, err = file.Write(capturedData)
	if err != nil {
		return fmt.Errorf("failed to write audio data: %w", err)
	}

	log.Printf("[I] Audio saved as WAV file: %s\n", filename)
	return nil
}
