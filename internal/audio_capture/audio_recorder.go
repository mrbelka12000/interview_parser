package audiocapture

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/gen2brain/malgo"
)

// WAV header structure for PCM audio
type WAVHeader struct {
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

// AudioDevice represents an audio device information
type AudioDevice struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsInput   bool   `json:"isInput"`
	IsOutput  bool   `json:"isOutput"`
	IsDefault bool   `json:"isDefault"`
}

// AudioRecorder handles audio recording and saving
type AudioRecorder struct {
	ctx           *malgo.AllocatedContext
	deviceConfig  malgo.DeviceConfig
	capturedData  []byte
	isRecording   bool
	volume        float32 // Volume multiplier (0.0 to 1.0)
	sampleRate    uint32
	channels      uint32
	bitsPerSample uint16
}

// NewAudioRecorder creates a new audio recorder instance
func NewAudioRecorder(sampleRate uint32, channels uint32, bitsPerSample uint16) (*AudioRecorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize context: %w", err)
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16 // 16-bit PCM
	deviceConfig.Capture.Channels = channels
	deviceConfig.SampleRate = sampleRate
	deviceConfig.Alsa.NoMMap = 1

	return &AudioRecorder{
		ctx:           ctx,
		deviceConfig:  deviceConfig,
		capturedData:  make([]byte, 0),
		isRecording:   false,
		sampleRate:    sampleRate,
		channels:      channels,
		volume:        1.0, // Default to 100% volume
		bitsPerSample: bitsPerSample,
	}, nil
}

// StartRecording begins capturing audio from the default microphone
func (ar *AudioRecorder) StartRecording() error {
	if ar.isRecording {
		return fmt.Errorf("recording is already in progress")
	}

	// Clear previous recording data
	ar.capturedData = make([]byte, 0)

	callbacks := malgo.DeviceCallbacks{
		Data: func(_, pSample []byte, _ uint32) {
			ar.capturedData = append(ar.capturedData, pSample...)
		},
	}

	device, err := malgo.InitDevice(ar.ctx.Context, ar.deviceConfig, callbacks)
	if err != nil {
		return fmt.Errorf("failed to initialize device: %w", err)
	}

	err = device.Start()
	if err != nil {
		return fmt.Errorf("failed to start device: %w", err)
	}

	ar.isRecording = true
	fmt.Printf("Recording started... Sample Rate: %d Hz, Channels: %d\n", ar.sampleRate, ar.channels)

	// Store device reference for stopping
	go func() {
		for ar.isRecording {
			time.Sleep(100 * time.Millisecond)
		}
		device.Uninit()
	}()

	return nil
}

// StopRecording stops the audio capture
func (ar *AudioRecorder) StopRecording() error {
	if !ar.isRecording {
		return fmt.Errorf("no recording in progress")
	}

	ar.isRecording = false
	fmt.Printf("Recording stopped. Captured %d bytes\n", len(ar.capturedData))
	return nil
}

// SaveAsWAV saves the captured audio as a WAV file
func (ar *AudioRecorder) SaveAsWAV(filename string) error {
	if len(ar.capturedData) == 0 {
		return fmt.Errorf("no audio data to save")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Calculate sizes
	dataSize := uint32(len(ar.capturedData))
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
		NumChannels:   uint16(ar.channels),
		SampleRate:    ar.sampleRate,
		ByteRate:      ar.sampleRate * uint32(ar.channels) * uint32(ar.bitsPerSample) / 8,
		BlockAlign:    uint16(ar.channels) * ar.bitsPerSample / 8,
		BitsPerSample: ar.bitsPerSample,
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: dataSize,
	}

	// Write header
	err = binary.Write(file, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write audio data
	_, err = file.Write(ar.applyVolume(ar.capturedData))
	if err != nil {
		return fmt.Errorf("failed to write audio data: %w", err)
	}

	fmt.Printf("Audio saved as WAV file: %s\n", filename)
	return nil
}

// GetCapturedData returns the captured audio data
func (ar *AudioRecorder) GetCapturedData() []byte {
	return ar.capturedData
}

// GetAudioInfo returns information about the captured audio
func (ar *AudioRecorder) GetAudioInfo() map[string]interface{} {
	duration := float64(len(ar.capturedData)) / float64(ar.sampleRate*ar.channels*2) // 2 bytes per sample for 16-bit
	return map[string]interface{}{
		"sample_rate":      ar.sampleRate,
		"channels":         ar.channels,
		"bits_per_sample":  ar.bitsPerSample,
		"data_size":        len(ar.capturedData),
		"duration_seconds": duration,
	}
}

// applyVolume applies volume to audio samples
func (ar *AudioRecorder) applyVolume(audioData []byte) []byte {
	if ar.volume == 1.0 {
		return audioData // No change needed
	}

	// Convert to int16 samples, apply volume, convert back
	samples := make([]int16, len(audioData)/2)
	for i := 0; i < len(samples); i++ {
		// Little-endian 16-bit samples
		sample := int16(audioData[i*2]) | int16(audioData[i*2+1])<<8
		sample = int16(float32(sample) * ar.volume)
		samples[i] = sample
	}

	// Convert back to bytes
	result := make([]byte, len(audioData))
	for i, sample := range samples {
		result[i*2] = byte(sample)
		result[i*2+1] = byte(sample >> 8)
	}

	return result
}

// GetInputDevices returns a list of available input devices
func (ar *AudioRecorder) GetInputDevices() ([]AudioDevice, error) {
	var devices []AudioDevice

	captureInfos, err := ar.ctx.Devices(malgo.Capture)
	if err != nil {
		return nil, fmt.Errorf("failed to get capture devices: %w", err)
	}

	for _, info := range captureInfos {
		device := AudioDevice{
			ID:        info.ID.String(),
			Name:      info.Name(),
			IsInput:   true,
			IsOutput:  false,
			IsDefault: info.IsDefault > 0,
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// SetInputDeviceByID sets the audio input device by device ID for recording
func (ar *AudioRecorder) SetInputDeviceByID(deviceID string) error {
	if ar.isRecording {
		return fmt.Errorf("cannot change input device while recording")
	}

	infos, err := ar.ctx.Devices(malgo.Capture)
	if err != nil {
		return fmt.Errorf("failed to get capture devices: %w", err)
	}

	for _, info := range infos {
		if info.ID.String() == deviceID {
			// Update device config with the specific device
			ar.deviceConfig.Capture.DeviceID = info.ID.Pointer()
			fmt.Printf("Input device set to: %s (%s)\n", info.Name(), info.ID.String())
			return nil
		}
	}

	return fmt.Errorf("input device with ID %s not found", deviceID)
}
