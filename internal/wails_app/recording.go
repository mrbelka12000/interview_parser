package wails_app

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	audiocapture "github.com/mrbelka12000/interview_parser/internal/audio_capture"
)

// RecordingResult represents the result of audio recording operations
type RecordingResult struct {
	Success  bool    `json:"success"`
	Message  string  `json:"message"`
	FilePath string  `json:"filePath,omitempty"`
	Duration float64 `json:"duration,omitempty"`
	DataSize int     `json:"dataSize,omitempty"`
}

// StartAudioRecording starts recording audio from microphone
func (a *App) StartAudioRecording() (*RecordingResult, error) {
	if a.audioRecorder == nil {
		return &RecordingResult{
			Success: false,
			Message: "Audio recorder not initialized",
		}, nil
	}

	err := a.audioRecorder.Start()
	if err != nil {
		return &RecordingResult{
			Success: false,
			Message: fmt.Sprintf("Failed to start recording: %s", err),
		}, nil
	}

	return &RecordingResult{
		Success: true,
		Message: "Recording started successfully",
	}, nil
}

// StopAudioRecording stops the audio recording and returns the recording info
func (a *App) StopAudioRecording() (*RecordingResult, error) {
	if a.audioRecorder == nil {
		return &RecordingResult{
			Success: false,
			Message: "Audio recorder not initialized",
		}, nil
	}

	err := a.audioRecorder.Stop()
	if err != nil {
		return &RecordingResult{
			Success: false,
			Message: fmt.Sprintf("Failed to stop recording: %s", err),
		}, nil
	}

	// Get audio info
	info := a.audioRecorder.GetAudioInfo()

	return &RecordingResult{
		Success:  true,
		Message:  "Recording stopped successfully",
		DataSize: info["data_size"].(int),
		Duration: info["duration_seconds"].(float64),
	}, nil
}

// SaveRecording saves the recorded audio to a file
func (a *App) SaveRecording(filename string) (*RecordingResult, error) {
	if a.audioRecorder == nil {
		return &RecordingResult{
			Success: false,
			Message: "Audio recorder not initialized",
		}, nil
	}

	if filename == "" {
		// Generate default filename with timestamp
		timestamp := time.Now().Format("20060102_150405")
		filename = fmt.Sprintf("interview_recording_%s.wav", timestamp)
	}

	// Ensure filename has .wav extension
	if !strings.HasSuffix(strings.ToLower(filename), ".wav") {
		filename += ".wav"
	}

	// Save to default directory
	filePath := filepath.Join(a.cfg.DefaultAnalyzeCallDir, filename)

	err := a.audioRecorder.SaveAsWAV(filePath)
	if err != nil {
		return &RecordingResult{
			Success: false,
			Message: fmt.Sprintf("Failed to save recording: %s", err),
		}, nil
	}

	// Get audio info
	info := a.audioRecorder.GetAudioInfo()

	return &RecordingResult{
		Success:  true,
		Message:  "Recording saved successfully",
		FilePath: filePath,
		Duration: info["duration_seconds"].(float64),
		DataSize: info["data_size"].(int),
	}, nil
}

// DeviceResult represents the result of device operations
type DeviceResult struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message,omitempty"`
	Devices []audiocapture.AudioDevice `json:"devices,omitempty"`
	Device  *audiocapture.AudioDevice  `json:"device,omitempty"`
}

// GetInputDevices returns a list of available input devices
func (a *App) GetInputDevices() (*DeviceResult, error) {
	if a.audioRecorder == nil {
		return &DeviceResult{
			Success: false,
			Message: "Audio recorder not initialized",
		}, nil
	}

	devices, err := a.audioRecorder.GetInputDevices()
	if err != nil {
		return &DeviceResult{
			Success: false,
			Message: fmt.Sprintf("Failed to get input devices: %s", err),
		}, nil
	}

	return &DeviceResult{
		Success: true,
		Message: "Input devices retrieved successfully",
		Devices: devices,
	}, nil
}

// SetAudioInputDevice sets the audio input device for recording
func (a *App) SetAudioInputDevice(deviceID string) (*DeviceResult, error) {
	if a.audioRecorder == nil {
		return &DeviceResult{
			Success: false,
			Message: "Audio recorder not initialized",
		}, nil
	}

	err := a.audioRecorder.SetInputDeviceByID(deviceID)
	if err != nil {
		return &DeviceResult{
			Success: false,
			Message: fmt.Sprintf("Failed to set input device: %s", err),
		}, nil
	}

	return &DeviceResult{
		Success: true,
		Message: "Input device set successfully",
	}, nil
}

// GetRecordingStatus returns the current recording status
func (a *App) GetRecordingStatus() (*RecordingResult, error) {
	if a.audioRecorder == nil {
		return &RecordingResult{
			Success: false,
			Message: "Audio recorder not initialized",
		}, nil
	}

	info := a.audioRecorder.GetAudioInfo()

	return &RecordingResult{
		Success:  true,
		Message:  "Status retrieved successfully",
		DataSize: info["data_size"].(int),
		Duration: info["duration_seconds"].(float64),
	}, nil
}
