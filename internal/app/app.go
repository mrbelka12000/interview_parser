package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/mrbelka12000/interview_parser/internal"
	audiocapture "github.com/mrbelka12000/interview_parser/internal/audio_capture"
	"github.com/mrbelka12000/interview_parser/internal/client"
	"github.com/mrbelka12000/interview_parser/internal/config"
	"github.com/mrbelka12000/interview_parser/internal/parser"
)

// FileInfo represents information about a file
type FileInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	IsDir     bool   `json:"isDir"`
	Size      int64  `json:"size"`
	Extension string `json:"extension"`
}

// App struct
type App struct {
	ctx           context.Context
	cfg           *config.Config
	aiClient      *client.Client
	parser        *parser.Parser
	audioRecorder *audiocapture.AudioRecorder
}

// NewApp creates a new App application struct
func NewApp(cfg *config.Config) *App {
	audioRecorder, err := audiocapture.NewAudioRecorder(cfg.AudioSampleRate, cfg.AudioChannels, cfg.AudioBitrate)
	if err != nil {
		log.Println(fmt.Sprintf("Error creating audio recorder %v", err))
	}

	return &App{
		cfg:           cfg,
		parser:        parser.NewParser(cfg),
		aiClient:      client.New(cfg),
		audioRecorder: audioRecorder,
	}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// GetFiles returns a list of files in the current directory
func (a *App) GetFiles() ([]FileInfo, error) {
	var files []FileInfo

	// Read directory contents
	entries, err := os.ReadDir(a.cfg.DefaultDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue // Skip files that can't be accessed
		}

		if strings.HasSuffix(entry.Name(), ".db") {
			continue
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext == "" && !entry.IsDir() {
			ext = "file"
		}

		fileInfo := FileInfo{
			Name:      info.Name(),
			Path:      filepath.Join(a.cfg.DefaultDir, info.Name()),
			IsDir:     entry.IsDir(),
			Size:      info.Size(),
			Extension: ext,
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// TranscriptionResult represents the result of transcription and analysis
type TranscriptionResult struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	TranscriptPath string `json:"transcriptPath,omitempty"`
	AnalysisPath   string `json:"analysisPath,omitempty"`
}

// ProcessFileForTranscription handles file upload and processing using the parser logic
func (a *App) ProcessFileForTranscription(filePath string, loadChunks bool) (*TranscriptionResult, error) {
	fmt.Printf("Processing file %s\n", filePath)
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &TranscriptionResult{
			Success: false,
			Message: fmt.Sprintf("file does not exist: %s", filePath),
		}, nil
	}

	apiKey, err := internal.GetOpenAIAPIKeyFromDB(a.cfg)
	if err != nil || apiKey == "" {
		return &TranscriptionResult{
			Success: false,
			Message: "No API Key provided",
		}, nil
	}

	// Create a copy of config with updated paths for this specific file
	cfg := *a.cfg
	cfg.LoadChunks = loadChunks
	cfg.OpenAIAPIKey = apiKey

	dir, err := os.ReadDir(cfg.DefaultDir)
	if err != nil {
		return &TranscriptionResult{
			Success: false,
			Message: fmt.Sprintf("failed to read working dir: %s", err),
		}, nil
	}

	// Generate unique output paths for this transcription
	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	transcriptPath := filepath.Join(cfg.DefaultTranscriptDir, fmt.Sprintf("%s_transcript_%v.txt", baseName, len(dir)))
	analyzePath := filepath.Join(cfg.DefaultAnalyzeDir, fmt.Sprintf("%s_analysis_%v.md", baseName, len(dir)))

	// Send initial progress
	a.sendProgress(5, "Loading file...", "Reading audio/video file...")

	// Step 1: Split into chunks or load existing chunks
	var chunks []string
	if cfg.LoadChunks {
		a.sendProgress(15, "Loading chunks...", "Loading existing chunk files...")
		chunks, err = a.parser.LoadChunks(&cfg)
		if err != nil {
			return &TranscriptionResult{
				Success: false,
				Message: fmt.Sprintf("failed to load chunks: %s", err),
			}, nil
		}
	} else {
		a.sendProgress(15, "Splitting into chunks...", "Dividing file into manageable segments...")
		chunks, err = a.parser.SplitIntoChunks(&cfg, filePath)
		if err != nil {
			return &TranscriptionResult{
				Success: false,
				Message: fmt.Sprintf("failed to split into chunks: %s", err),
			}, nil
		}
	}

	if len(chunks) == 0 {
		return &TranscriptionResult{
			Success: false,
			Message: "no chunks to process",
		}, nil
	}

	// Step 2: Transcribe chunks using the provided parser logic
	a.sendProgress(25, "Transcribing audio...", "Converting speech to text using AI...")

	var (
		wg            sync.WaitGroup
		mx            sync.Mutex
		workers       = make(chan struct{}, cfg.ParallelWorkers)
		collectedText = make([]string, len(chunks))
		completed     = 0
	)

	for i, chunk := range chunks {
		wg.Add(1)
		workers <- struct{}{}
		go func(ind int, chunkVar string) {
			defer func() {
				<-workers
				wg.Done()
			}()

			textFromChunk, err := a.aiClient.Transcribe(a.ctx, chunkVar)
			if err != nil {
				log.Printf("Error transcribing chunk %d: %v", ind, err)
				return
			}

			mx.Lock()
			collectedText[ind] = textFromChunk
			completed++
			// Update progress based on completed chunks
			progress := 25 + int(float64(completed)/float64(len(chunks))*35) // 25% to 60%
			a.sendProgress(progress, "Processing chunks...", fmt.Sprintf("Processed %d/%d audio segments...", completed, len(chunks)))
			mx.Unlock()

		}(i, chunk)
	}

	wg.Wait()

	if errors.Is(a.ctx.Err(), context.Canceled) {
		log.Println("[i] cancelled by signal, skip analyze")
		return &TranscriptionResult{
			Success: false,
			Message: "processing was cancelled",
		}, nil
	}

	// Step 3: Join and format transcript
	a.sendProgress(65, "Formatting transcript...", "Organizing and formatting text...")
	transcript := strings.Join(collectedText, "")
	transcript = a.parser.FormatText(transcript)
	// Step 4: Save transcript
	a.sendProgress(75, "Saving transcript...", "Writing transcript file...")
	if err = a.parser.SaveTranscript(transcriptPath, transcript); err != nil {
		return &TranscriptionResult{
			Success: false,
			Message: fmt.Sprintf("failed to save transcript: %s", err),
		}, nil
	}

	transcriptBatches := a.parser.BatchTranscript(transcript)
	analyzeRespBatches := make([][]client.Question, len(transcriptBatches))
	completed = 0
	for i, batch := range transcriptBatches {
		wg.Add(1)
		workers <- struct{}{}
		go func(ind int, b string) {
			defer func() {
				<-workers
				wg.Done()
			}()

			analyzeRespTmp, err := a.aiClient.AnalyzeTranscript(a.ctx, b)
			if err != nil {
				log.Printf("Error analyzing transcript %d: %v", ind, err)
				return
			}

			mx.Lock()
			analyzeRespBatches[ind] = analyzeRespTmp.Questions
			completed++
			progress := 85 + int(float64(completed)/float64(len(transcriptBatches))*10) // 85% to 95%

			a.sendProgress(progress, "Analyzing transcript...", fmt.Sprintf("Analyzed %d/%d tranascript segments...", completed, len(transcriptBatches)))
			mx.Unlock()

		}(i, batch)

	}

	wg.Wait()
	close(workers)

	var analyzeResp client.AnalyzeResponse
	for _, batch := range analyzeRespBatches {
		analyzeResp.Questions = append(analyzeResp.Questions, batch...)
	}

	// Step 6: Save analysis response
	a.sendProgress(97, "Saving analysis...", "Writing analysis file...")
	if err = a.parser.SaveAnalyzeResponse(analyzePath, analyzeResp); err != nil {
		return &TranscriptionResult{
			Success: false,
			Message: fmt.Sprintf("failed to save analysis: %s", err),
		}, nil
	}

	// Final progress
	a.sendProgress(100, "Complete!", "Processing finished successfully!")

	return &TranscriptionResult{
		Success:        true,
		Message:        "File processed successfully",
		TranscriptPath: transcriptPath,
		AnalysisPath:   analyzePath,
	}, nil
}

// sendProgress sends progress update to frontend
func (a *App) sendProgress(percentage int, stage, details string) {
	report := ProgressReport{
		Percentage: percentage,
		Stage:      stage,
		Details:    details,
		IsComplete: percentage >= 100,
	}

	// Send event to frontend
	runtime.EventsEmit(a.ctx, "progress", report)
}

// ProcessFile handles file click actions
func (a *App) ProcessFile(filePath string) (*FileInfo, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	if info.IsDir() {
		return &FileInfo{
			Name:  info.Name(),
			Path:  filepath.Join(filePath, info.Name()),
			IsDir: true,
			Size:  info.Size(),
		}, nil
	}

	// For files, return basic information
	return &FileInfo{
		Name:      info.Name(),
		Path:      filepath.Join(filePath, info.Name()),
		Size:      info.Size(),
		Extension: filepath.Ext(info.Name()),
	}, nil
}

// FileContent represents file information along with its content
type FileContent struct {
	*FileInfo
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

// ReadFileContent reads the content of a file
func (a *App) ReadFileContent(filePath string) (*FileContent, error) {
	// Check if file exists
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return &FileContent{
			Error: fmt.Sprintf("file does not exist: %s", filePath),
		}, nil
	}
	if err != nil {
		return &FileContent{
			Error: fmt.Sprintf("failed to get file info: %v", err),
		}, nil
	}

	// Don't try to read directories
	if info.IsDir() {
		return &FileContent{
			Error: "cannot read content of a directory",
		}, nil
	}

	// Check file size (prevent reading very large files)
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if info.Size() > maxFileSize {
		return &FileContent{
			Error: fmt.Sprintf("file is too large (%.2f MB), maximum size is 10 MB", float64(info.Size())/1024/1024),
		}, nil
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return &FileContent{
			Error: fmt.Sprintf("failed to read file: %v", err),
		}, nil
	}

	// Create file info
	fileInfo := FileInfo{
		Name:      info.Name(),
		Path:      filePath,
		Size:      info.Size(),
		Extension: filepath.Ext(info.Name()),
	}

	return &FileContent{
		FileInfo: &fileInfo,
		Content:  string(content),
	}, nil
}

// GetFilesInDirectory returns files in a specific directory
func (a *App) GetFilesInDirectory(dirPath string) ([]FileInfo, error) {
	files := []FileInfo{}

	// Read directory contents
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue // Skip files that can't be accessed
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext == "" && !entry.IsDir() {
			ext = "file"
		}

		fileInfo := FileInfo{
			Name:      info.Name(),
			Path:      filepath.Join(dirPath, info.Name()),
			IsDir:     entry.IsDir(),
			Size:      info.Size(),
			Extension: ext,
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// PickFile opens native file dialog and returns full path
func (a *App) PickFile() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Pick media file",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Media files",
				Pattern:     "*.mp3;*.wav;*.m4a;*.mp4;*.mov;*.avi",
			},
		},
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// APIKeyResult represents the result of API key operations
type APIKeyResult struct {
	Success     bool   `json:"success"`
	Message     string `json:"message,omitempty"`
	APIKey      string `json:"apiKey,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

// ProgressReport represents progress update during processing
type ProgressReport struct {
	Percentage int    `json:"percentage"`
	Stage      string `json:"stage"`
	Details    string `json:"details"`
	IsComplete bool   `json:"isComplete"`
}

// GetOpenAIAPIKey retrieves the current OpenAI API key from database
func (a *App) GetOpenAIAPIKey() (*APIKeyResult, error) {
	apiKey, err := internal.GetOpenAIAPIKeyFromDB(a.cfg)
	if err != nil {
		if errors.Is(err, internal.ErrNoKey) {
			return &APIKeyResult{
				Success: false,
				Message: "No API key found in database",
			}, nil
		}
		return &APIKeyResult{
			Success: false,
			Message: fmt.Sprintf("Failed to retrieve API key: %s", err),
		}, nil
	}

	return &APIKeyResult{
		Success:     true,
		APIKey:      apiKey,
		LastUpdated: "Recently updated",
	}, nil
}

// SaveOpenAIAPIKey saves a new OpenAI API key to database
func (a *App) SaveOpenAIAPIKey(apiKey string) (*APIKeyResult, error) {
	if apiKey == "" {
		return &APIKeyResult{
			Message: "API key cannot be empty",
		}, nil
	}

	// Basic validation for OpenAI API key
	if !strings.HasPrefix(apiKey, "sk-") {
		return &APIKeyResult{
			Message: "Invalid API key format. OpenAI API keys should start with 'sk-'",
		}, nil
	}

	fmt.Printf("Saving API key: %s\n", apiKey)
	// Save to config for current session
	tmpCfg := *a.cfg
	tmpCfg.OpenAIAPIKey = apiKey

	aiClient := client.New(&tmpCfg)
	err := aiClient.IsValidAPIKeysProvided()
	if err != nil {
		return &APIKeyResult{
			Message: err.Error(),
		}, nil
	}

	// Save to database
	err = internal.InsertOpenAIAPIKey(&tmpCfg)
	if err != nil {
		return &APIKeyResult{
			Success: false,
			Message: fmt.Sprintf("Failed to save API key: %s", err),
		}, nil
	}

	a.cfg = &tmpCfg
	a.aiClient = aiClient

	return &APIKeyResult{
		Success:     true,
		Message:     "API key saved successfully",
		APIKey:      apiKey,
		LastUpdated: "Just now",
	}, nil
}

// DeleteOpenAIAPIKey removes the OpenAI API key from database
func (a *App) DeleteOpenAIAPIKey() (*APIKeyResult, error) {
	err := internal.DeleteOpenAIAPIKey(a.cfg)
	if err != nil {
		return &APIKeyResult{
			Success: false,
			Message: fmt.Sprintf("Failed to delete API key: %s", err),
		}, nil
	}

	// Clear from current config
	a.cfg.OpenAIAPIKey = ""

	return &APIKeyResult{
		Success: true,
		Message: "API key deleted successfully",
	}, nil
}

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

	err := a.audioRecorder.StartRecording()
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

	err := a.audioRecorder.StopRecording()
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
	filePath := filepath.Join(a.cfg.DefaultDir, filename)

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

// SaveAndProcessRecording saves the recording and immediately processes it for transcription
func (a *App) SaveAndProcessRecording(filename string) (*TranscriptionResult, error) {
	// First save the recording
	saveResult, err := a.SaveRecording(filename)
	if err != nil {
		return nil, err
	}

	if !saveResult.Success {
		return &TranscriptionResult{
			Success: false,
			Message: saveResult.Message,
		}, nil
	}

	// Then process the saved file for transcription
	return a.ProcessFileForTranscription(saveResult.FilePath, false)
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
