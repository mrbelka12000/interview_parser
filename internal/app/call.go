package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// CallAnalysisResult represents the result of call analysis
type CallAnalysisResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	AnalysisPath string `json:"analysisPath,omitempty"`
}

// SaveAndProcessRecordingForCall saves the recording and analyzes it as a call
func (a *App) SaveAndProcessRecordingForCall(filename string) (*CallAnalysisResult, error) {
	// First save the recording
	saveResult, err := a.SaveRecording(filename)
	if err != nil {
		return nil, err
	}

	if !saveResult.Success {
		return &CallAnalysisResult{
			Success: false,
			Message: saveResult.Message,
		}, nil
	}

	// Then process the saved file for call analysis
	return a.ProcessFileForCallAnalysis(saveResult.FilePath)
}

// ProcessFileForCallAnalysis processes file for call analysis
func (a *App) ProcessFileForCallAnalysis(filePath string) (*CallAnalysisResult, error) {
	fmt.Printf("Processing file for call analysis %s\n", filePath)

	defer os.Remove(filePath)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &CallAnalysisResult{
			Success: false,
			Message: fmt.Sprintf("file does not exist: %s", filePath),
		}, nil
	}

	apiKey, err := a.service.GetAPIKey()
	if err != nil || apiKey == "" {
		return &CallAnalysisResult{
			Success: false,
			Message: "No API Key provided",
		}, nil
	}

	dir, err := os.ReadDir(a.cfg.DefaultDir)
	if err != nil {
		return &CallAnalysisResult{
			Success: false,
			Message: fmt.Sprintf("failed to read working dir: %s", err),
		}, nil
	}

	// Generate unique output paths for this analysis
	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	analysisCallPath := filepath.Join(a.cfg.DefaultAnalyzeCallDir, fmt.Sprintf("%s_call_analysis_%v.md", baseName, len(dir)))

	transcript, err := a.transcribeFile(filePath)
	if err != nil {
		return &CallAnalysisResult{
			Message: err.Error(),
		}, nil
	}

	if errors.Is(a.ctx.Err(), context.Canceled) {
		log.Println("[i] cancelled by signal, skip analyze")
		return &CallAnalysisResult{
			Success: false,
			Message: "processing was cancelled",
		}, nil
	}

	err = a.analyzeCall(transcript)
	if err != nil {
		return &CallAnalysisResult{
			Message: err.Error(),
		}, nil
	}

	// Final progress
	a.sendProgress(100, "Complete!", "Call analysis finished successfully!")

	return &CallAnalysisResult{
		Success:      true,
		Message:      "File processed successfully",
		AnalysisPath: analysisCallPath,
	}, nil
}
