package audiocapture

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gen2brain/malgo"

	"github.com/mrbelka12000/interview_parser/internal/wav"
)

// AudioDevice represents an audio device information
type AudioDevice struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsInput   bool   `json:"isInput"`
	IsOutput  bool   `json:"isOutput"`
	IsDefault bool   `json:"isDefault"`
}

// AudioCapturer handles audio capturing and saving
type AudioCapturer struct {
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
	isCapturing   bool
	sampleRate    uint32
	channels      uint32
	bitsPerSample uint16

	wav *wav.Writer
}

// NewAudioCapturer creates a new audio recorder instance
func NewAudioCapturer(sampleRate uint32, channels uint32, bitsPerSample uint16) (*AudioCapturer, error) {
	defaultCtx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		log.Printf("[I] AUDIO_CAPTURER <%v>\n", strings.TrimSuffix(message, "\n"))
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize context: %w", err)
	}

	defaultConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	defaultConfig.Capture.Format = malgo.FormatS16 // 16-bit PCM
	defaultConfig.Capture.Channels = channels
	defaultConfig.SampleRate = sampleRate
	defaultConfig.Alsa.NoMMap = 1

	ar := &AudioCapturer{
		defaultCtx:    defaultCtx,
		defaultConfig: defaultConfig,
		capturedData:  make([]byte, 0),
		isCapturing:   false,
		sampleRate:    sampleRate,
		channels:      channels,
		bitsPerSample: bitsPerSample,
		micCh:         make(chan []byte, 64), // NEW
		bhCh:          make(chan []byte, 64), // NEW
		wav:           wav.NewWriter(sampleRate, channels, bitsPerSample),
	}

	blackHoleCtx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		log.Printf("[I] LOG <%v>\n", strings.TrimSuffix(message, "\n"))
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

func (ar *AudioCapturer) getBlackHoleDevice() (device malgo.DeviceInfo) {
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

// Start begins capturing audio from the mic + BlackHole and merges them
func (ar *AudioCapturer) Start() error {
	if ar.isCapturing {
		return fmt.Errorf("capturing is already in progress")
	}

	// Clear previous capturing data
	ar.capturedData = make([]byte, 0)
	ar.isCapturing = true

	// Recreate channels in case of multiple capturing
	ar.micCh = make(chan []byte, 64)
	ar.bhCh = make(chan []byte, 64)

	// --- MIC CALLBACK ---
	micCallbacks := malgo.DeviceCallbacks{
		Data: func(_, in []byte, _ uint32) {
			if !ar.isCapturing || len(in) == 0 {
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
			if !ar.isCapturing || len(in) == 0 {
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
		ar.isCapturing = false
		return fmt.Errorf("failed to initialize mic device: %w", err)
	}
	ar.micDevice = micDevice

	// Init BlackHole device
	bhDevice, err := malgo.InitDevice(ar.blackHoleCtx.Context, ar.blackHoleConfig, bhCallbacks)
	if err != nil {
		ar.isCapturing = false
		ar.micDevice.Uninit()
		ar.micDevice = nil
		return fmt.Errorf("failed to initialize BlackHole device: %w", err)
	}
	ar.blackHoleDevice = bhDevice

	// Start both devices
	if err := ar.micDevice.Start(); err != nil {
		ar.isCapturing = false
		ar.micDevice.Uninit()
		ar.blackHoleDevice.Uninit()
		ar.micDevice = nil
		ar.blackHoleDevice = nil
		return fmt.Errorf("failed to start mic device: %w", err)
	}

	if err := ar.blackHoleDevice.Start(); err != nil {
		ar.isCapturing = false
		ar.micDevice.Stop()
		ar.micDevice.Uninit()
		ar.blackHoleDevice.Uninit()
		ar.micDevice = nil
		ar.blackHoleDevice = nil
		return fmt.Errorf("failed to start BlackHole device: %w", err)
	}

	log.Printf("[I] Capturing started... Sample Rate: %d Hz, Channels: %d\n", ar.sampleRate, ar.channels)

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

	return nil
}

// Stop stops the audio capture
func (ar *AudioCapturer) Stop() error {
	if !ar.isCapturing {
		return fmt.Errorf("no capturing in progress")
	}

	ar.isCapturing = false

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

	log.Printf("[I] Capturing stopped. Captured %d bytes\n", len(ar.capturedData))
	return nil
}

// GetAudioInfo returns information about the captured audio
func (ar *AudioCapturer) GetAudioInfo() map[string]interface{} {
	duration := float64(len(ar.capturedData)) / float64(ar.sampleRate*ar.channels*2) // 2 bytes per sample for 16-bit
	return map[string]interface{}{
		"sample_rate":      ar.sampleRate,
		"channels":         ar.channels,
		"bits_per_sample":  ar.bitsPerSample,
		"data_size":        len(ar.capturedData),
		"duration_seconds": duration,
	}
}

// GetInputDevices returns a list of available input devices
func (ar *AudioCapturer) GetInputDevices() ([]AudioDevice, error) {
	var devices []AudioDevice

	captureInfos, err := ar.defaultCtx.Devices(malgo.Capture)
	if err != nil {
		return nil, fmt.Errorf("failed to get capture devices: %w", err)
	}

	for _, info := range captureInfos {
		if strings.Contains(strings.ToLower(info.Name()), "blackhole") || strings.Contains(strings.ToLower(info.Name()), "pods") {
			continue
		}

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

// SetInputDeviceByID sets the audio input device by device ID for capturing
func (ar *AudioCapturer) SetInputDeviceByID(deviceID string) error {
	if ar.isCapturing {
		return fmt.Errorf("cannot change input device while capturing")
	}

	infos, err := ar.defaultCtx.Devices(malgo.Capture)
	if err != nil {
		return fmt.Errorf("failed to get capture devices: %w", err)
	}

	for _, info := range infos {
		if info.ID.String() == deviceID {
			// Update device config with the specific device
			ar.defaultConfig.Capture.DeviceID = info.ID.Pointer()
			log.Printf("[I] Input device set to: %s (%s)\n", info.Name(), info.ID.String())
			return nil
		}
	}

	return fmt.Errorf("input device with ID %s not found", deviceID)
}

func (ar *AudioCapturer) SaveAsWAV(filePath string) error {
	return ar.wav.SaveAsWAV(filePath, ar.capturedData)
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
