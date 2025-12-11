package audio_recorder

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gen2brain/malgo"

	"github.com/mrbelka12000/interview_parser/internal/wav"
)

type AudioRecorder struct {
	ctx    *malgo.AllocatedContext
	config malgo.DeviceConfig

	micDevice *malgo.Device // NEW: mic capture device

	mx sync.Mutex

	recordedData  []byte
	isRecording   bool
	sampleRate    uint32
	channels      uint32
	bitsPerSample uint16

	ch chan []byte

	wav *wav.Writer
}

func NewAudioRecorder(sampleRate uint32, channels uint32, bitsPerSample uint16, outCh chan []byte) (*AudioRecorder, error) {
	defaultCtx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		log.Printf("[I] AUDIO_RECORDER <%v>\n", strings.TrimSuffix(message, "\n"))
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize context: %w", err)
	}

	defaultConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	defaultConfig.Capture.Format = malgo.FormatS16 // 16-bit PCM
	defaultConfig.Capture.Channels = channels
	defaultConfig.SampleRate = sampleRate
	defaultConfig.Alsa.NoMMap = 1

	return &AudioRecorder{
		ctx:           defaultCtx,
		config:        defaultConfig,
		sampleRate:    sampleRate,
		channels:      channels,
		bitsPerSample: bitsPerSample,
		ch:            outCh,
		wav:           wav.NewWriter(sampleRate, channels, bitsPerSample),
	}, nil
}

func (ar *AudioRecorder) Start() error {
	if ar.isRecording {
		return fmt.Errorf("recording is already in progress")
	}

	ar.recordedData = make([]byte, 0)
	ar.isRecording = true

	callback := malgo.DeviceCallbacks{
		Data: func(_, pInputSamples []byte, _ uint32) {
			ar.ch <- pInputSamples
			ar.recordedData = append(ar.recordedData, pInputSamples...)
		},
	}
	// Init mic device
	micDevice, err := malgo.InitDevice(ar.ctx.Context, ar.config, callback)
	if err != nil {
		ar.isRecording = false
		return fmt.Errorf("failed to initialize mic device: %w", err)
	}
	ar.micDevice = micDevice

	if err := ar.micDevice.Start(); err != nil {
		ar.isRecording = false
		ar.micDevice.Uninit()
		ar.micDevice = nil
		return fmt.Errorf("failed to start mic device: %w", err)
	}

	log.Printf("[I] Recording started... Sample Rate: %d Hz, Channels: %d\n", ar.sampleRate, ar.channels)

	go func() {
		for ar.isRecording {
			time.Sleep(100 * time.Millisecond)
		}

		ar.micDevice.Uninit()
	}()

	return nil
}

func (ar *AudioRecorder) Stop() error {

	if !ar.isRecording {
		return fmt.Errorf("no recording in progress")
	}

	ar.isRecording = false
	if ar.micDevice != nil {
		ar.micDevice.Uninit()
	}

	log.Printf("[I] Recording stopped. Recorder %d bytes\n", len(ar.recordedData))

	return nil
}

// GetAudioInfo returns information about the captured audio
func (ar *AudioRecorder) GetAudioInfo() map[string]interface{} {
	duration := float64(len(ar.recordedData)) / float64(ar.sampleRate*ar.channels*2) // 2 bytes per sample for 16-bit
	return map[string]interface{}{
		"sample_rate":      ar.sampleRate,
		"channels":         ar.channels,
		"bits_per_sample":  ar.bitsPerSample,
		"data_size":        len(ar.recordedData),
		"duration_seconds": duration,
	}
}
