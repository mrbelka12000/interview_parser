package wails_app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mrbelka12000/interview_parser/internal/models"
)

// TranscriptionResult represents the result of transcription and analysis
type TranscriptionResult struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	TranscriptPath string `json:"transcriptPath,omitempty"`
	AnalysisPath   string `json:"analysisPath,omitempty"`
}

func (a *App) SaveInterviewAPI(interview *models.AnalyzeInterviewWithQA) (int64, error) {
	if err := a.service.SaveInterview(interview); err != nil {
		return 0, err
	}

	return 0, nil
}

// GetAllInterviewsAPI retrieves all interviews with optional date filters
func (a *App) GetAllInterviewsAPI(dateFrom, dateTo string) ([]models.AnalyzeInterviewWithQA, error) {
	filters := &models.GetInterviewsFilters{}

	if dateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters.DateFrom = &parsed
		}
	}

	if dateTo != "" {
		if parsed, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters.DateTo = &parsed
		}
	}

	return a.service.GetAllInterviews(filters)
}

// GetInterviewAPI retrieves a specific interview by ID
func (a *App) GetInterviewAPI(id uint64) (*models.AnalyzeInterviewWithQA, error) {
	return a.service.GetInterview(id)
}

// DeleteInterviewAPI deletes an interview and its question answers
func (a *App) DeleteInterviewAPI(id uint64) error {
	return a.service.DeleteInterview(id)
}

// UpdateInterviewAPI updates an existing interview and its question answers
func (a *App) UpdateInterviewAPI(interview *models.AnalyzeInterview, qaList []models.QuestionAnswer) error {
	return a.service.UpdateInterview(interview, qaList)
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
	return a.ProcessFileForTranscription(saveResult.FilePath)
}

// ProcessFileForTranscription handles file upload and processing using the parser logic
func (a *App) ProcessFileForTranscription(filePath string) (*TranscriptionResult, error) {
	fmt.Printf("Processing file %s\n", filePath)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &TranscriptionResult{
			Success: false,
			Message: fmt.Sprintf("file does not exist: %s", filePath),
		}, nil
	}

	apiKey, err := a.service.GetAPIKey()
	if err != nil || apiKey == "" {
		return &TranscriptionResult{
			Success: false,
			Message: "No API Key provided",
		}, nil
	}

	dir, err := os.ReadDir(a.cfg.DefaultDir)
	if err != nil {
		return &TranscriptionResult{
			Success: false,
			Message: fmt.Sprintf("failed to read working dir: %s", err),
		}, nil
	}

	// Generate unique output paths for this transcription
	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	transcriptPath := filepath.Join(a.cfg.DefaultTranscriptDir, fmt.Sprintf("%s_transcript_%v.txt", baseName, len(dir)))

	transcript, err := a.transcribeFile(filePath)
	if err != nil {
		return &TranscriptionResult{
			Message: err.Error(),
		}, nil
	}

	transcript = a.parser.FormatText(transcript)

	// Step 4: Save transcript
	a.sendProgress(75, "Saving transcript...", "Writing transcript file...")
	if err = a.parser.SaveTranscript(transcriptPath, transcript); err != nil {
		return &TranscriptionResult{
			Success: false,
			Message: fmt.Sprintf("failed to save transcript: %s", err),
		}, nil
	}

	err = a.analyzeInterview(transcript)
	if err != nil {
		return &TranscriptionResult{
			Message: err.Error(),
		}, nil
	}

	// Final progress
	a.sendProgress(100, "Complete!", "Processing finished successfully!")

	return &TranscriptionResult{
		Success:        true,
		Message:        "File processed successfully",
		TranscriptPath: transcriptPath,
	}, nil
}
