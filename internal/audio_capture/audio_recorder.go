package audiocapture

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"sync"
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
	defaultCtx      *malgo.AllocatedContext
	blackHoleCtx    *malgo.AllocatedContext
	defaultConfig   malgo.DeviceConfig
	blackHoleConfig malgo.DeviceConfig

	micDevice       *malgo.Device // NEW: mic capture device
	blackHoleDevice *malgo.Device // NEW: output (BlackHole) capture device

	micCh chan []byte // NEW: mic PCM chunks
	bhCh  chan []byte // NEW: BlackHole PCM chunks

	mx sync.Mutex

	capturedData  []byte
	isRecording   bool
	volume        float32 // Volume multiplier (0.0 to 1.0)
	sampleRate    uint32
	channels      uint32
	bitsPerSample uint16
}

// NewAudioRecorder creates a new audio recorder instance
func NewAudioRecorder(sampleRate uint32, channels uint32, bitsPerSample uint16) (*AudioRecorder, error) {
	defaultCtx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", strings.TrimSuffix(message, "\n"))
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize context: %w", err)
	}

	defaultConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	defaultConfig.Capture.Format = malgo.FormatS16 // 16-bit PCM
	defaultConfig.Capture.Channels = channels
	defaultConfig.SampleRate = sampleRate
	defaultConfig.Alsa.NoMMap = 1

	ar := &AudioRecorder{
		defaultCtx:    defaultCtx,
		defaultConfig: defaultConfig,
		capturedData:  make([]byte, 0),
		isRecording:   false,
		sampleRate:    sampleRate,
		channels:      channels,
		volume:        1.0, // Default to 100% volume
		bitsPerSample: bitsPerSample,
		micCh:         make(chan []byte, 64), // NEW
		bhCh:          make(chan []byte, 64), // NEW
	}

	blackHoleCtx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", strings.TrimSuffix(message, "\n"))
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize context: %w", err)
	}

	ar.blackHoleCtx = blackHoleCtx
	blackHoleConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	blackHoleConfig.Capture.Format = malgo.FormatS16 // 16-bit PCM
	blackHoleConfig.Capture.Channels = channels
	blackHoleConfig.SampleRate = sampleRate
	device := ar.getBlackHoleDevice()
	blackHoleConfig.Capture.DeviceID = device.ID.Pointer()
	blackHoleConfig.Alsa.NoMMap = 1

	ar.blackHoleConfig = blackHoleConfig

	return ar, nil
}

func (ar *AudioRecorder) getBlackHoleDevice() (device malgo.DeviceInfo) {
	captureInfos, err := ar.blackHoleCtx.Devices(malgo.Capture)
	if err != nil {
		return
	}

	for _, captureInfo := range captureInfos {
		if strings.Contains(strings.ToLower(captureInfo.Name()), "blackhole") {
			return captureInfo
		}
	}

	return
}

// StartRecording begins capturing audio from the mic + BlackHole and merges them
func (ar *AudioRecorder) StartRecording() error {
	if ar.isRecording {
		return fmt.Errorf("recording is already in progress")
	}

	// Clear previous recording data
	ar.capturedData = make([]byte, 0)
	ar.isRecording = true

	// Recreate channels in case of multiple recordings
	ar.micCh = make(chan []byte, 64)
	ar.bhCh = make(chan []byte, 64)

	// --- MIC CALLBACK ---
	micCallbacks := malgo.DeviceCallbacks{
		Data: func(_, in []byte, _ uint32) {
			if !ar.isRecording || len(in) == 0 {
				return
			}
			buf := make([]byte, len(in))
			copy(buf, in)
			select {
			case ar.micCh <- buf:
			default:
				// channel full - drop
			}
		},
	}

	// --- BLACKHOLE CALLBACK ---
	bhCallbacks := malgo.DeviceCallbacks{
		Data: func(_, in []byte, _ uint32) {
			if !ar.isRecording || len(in) == 0 {
				return
			}
			buf := make([]byte, len(in))
			copy(buf, in)
			select {
			case ar.bhCh <- buf:
			default:
				// channel full - drop
			}
		},
	}

	// Init mic device
	micDevice, err := malgo.InitDevice(ar.defaultCtx.Context, ar.defaultConfig, micCallbacks)
	if err != nil {
		ar.isRecording = false
		return fmt.Errorf("failed to initialize mic device: %w", err)
	}
	ar.micDevice = micDevice

	// Init BlackHole device
	bhDevice, err := malgo.InitDevice(ar.blackHoleCtx.Context, ar.blackHoleConfig, bhCallbacks)
	if err != nil {
		ar.isRecording = false
		ar.micDevice.Uninit()
		ar.micDevice = nil
		return fmt.Errorf("failed to initialize BlackHole device: %w", err)
	}
	ar.blackHoleDevice = bhDevice

	// Start both devices
	if err := ar.micDevice.Start(); err != nil {
		ar.isRecording = false
		ar.micDevice.Uninit()
		ar.blackHoleDevice.Uninit()
		ar.micDevice = nil
		ar.blackHoleDevice = nil
		return fmt.Errorf("failed to start mic device: %w", err)
	}

	if err := ar.blackHoleDevice.Start(); err != nil {
		ar.isRecording = false
		ar.micDevice.Stop()
		ar.micDevice.Uninit()
		ar.blackHoleDevice.Uninit()
		ar.micDevice = nil
		ar.blackHoleDevice = nil
		return fmt.Errorf("failed to start BlackHole device: %w", err)
	}

	fmt.Printf("Recording started... Sample Rate: %d Hz, Channels: %d\n", ar.sampleRate, ar.channels)

	// Mixer goroutine: merge mic + BlackHole into capturedData
	go func() {
		for {
			micBuf, ok1 := <-ar.micCh
			bhBuf, ok2 := <-ar.bhCh

			if !ok1 || !ok2 {
				return
			}

			ar.mx.Lock()
			ar.capturedData = append(ar.capturedData, mixInt16Stereo(micBuf, bhBuf)...)
			ar.mx.Unlock()
		}
	}()

	// Watcher to uninit devices once recording stops (extra safety)
	go func() {
		for ar.isRecording {
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return nil
}

// StopRecording stops the audio capture
func (ar *AudioRecorder) StopRecording() error {
	if !ar.isRecording {
		return fmt.Errorf("no recording in progress")
	}

	ar.isRecording = false

	// Stop and uninit devices
	if ar.micDevice != nil {
		_ = ar.micDevice.Stop()
		ar.micDevice.Uninit()
		ar.micDevice = nil
	}
	if ar.blackHoleDevice != nil {
		_ = ar.blackHoleDevice.Stop()
		ar.blackHoleDevice.Uninit()
		ar.blackHoleDevice = nil
	}

	// Close channels so mixer goroutine can exit
	if ar.micCh != nil {
		close(ar.micCh)
	}
	if ar.bhCh != nil {
		close(ar.bhCh)
	}

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

	captureInfos, err := ar.defaultCtx.Devices(malgo.Capture)
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

	infos, err := ar.defaultCtx.Devices(malgo.Capture)
	if err != nil {
		return fmt.Errorf("failed to get capture devices: %w", err)
	}

	for _, info := range infos {
		if info.ID.String() == deviceID {
			// Update device config with the specific device
			ar.defaultConfig.Capture.DeviceID = info.ID.Pointer()
			fmt.Printf("Input device set to: %s (%s)\n", info.Name(), info.ID.String())
			return nil
		}
	}

	return fmt.Errorf("input device with ID %s not found", deviceID)
}

func mixInt16Stereo(a, b []byte) []byte {
	// Per-stream gains – tune these if needed
	const gainA = 0.5
	const gainB = 0.5

	// Both empty
	if len(a) == 0 && len(b) == 0 {
		return nil
	}
	// Only one present – just copy it
	if len(a) == 0 {
		out := make([]byte, len(b))
		copy(out, b)
		return out
	}
	if len(b) == 0 {
		out := make([]byte, len(a))
		copy(out, a)
		return out
	}

	samplesA := len(a) / 2
	samplesB := len(b) / 2
	maxSamples := samplesA
	if samplesB > maxSamples {
		maxSamples = samplesB
	}

	out := make([]byte, maxSamples*2)

	for i := 0; i < maxSamples; i++ {
		var s1, s2 int16

		if i < samplesA {
			s1 = int16(binary.LittleEndian.Uint16(a[i*2:]))
		}
		if i < samplesB {
			s2 = int16(binary.LittleEndian.Uint16(b[i*2:]))
		}

		// apply gain and mix
		v := float32(s1)*gainA + float32(s2)*gainB
		iv := int32(v)

		// clamp to int16
		if iv > 32767 {
			iv = 32767
		} else if iv < -32768 {
			iv = -32768
		}

		binary.LittleEndian.PutUint16(out[i*2:], uint16(int16(iv)))
	}

	return out
}
