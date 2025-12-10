package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

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

// ProgressReport represents progress update during processing
type ProgressReport struct {
	Percentage int    `json:"percentage"`
	Stage      string `json:"stage"`
	Details    string `json:"details"`
	IsComplete bool   `json:"isComplete"`
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

// GetInterviewAnalyticsAPI retrieves analytics for a specific interview
func (a *App) GetInterviewAnalyticsAPI(interviewPath string) (*InterviewAnalytics, error) {
	return a.GetInterviewAnalytics(interviewPath)
}

// GetAllInterviewAnalyticsAPI retrieves all interview analytics with optional filters
func (a *App) GetAllInterviewAnalyticsAPI(dateFrom, dateTo string, minAccuracy, maxAccuracy float64) ([]InterviewAnalytics, error) {
	filters := &AnalyticsFilters{}
	
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
	
	if minAccuracy > 0 {
		filters.MinAccuracy = &minAccuracy
	}
	
	if maxAccuracy > 0 {
		filters.MaxAccuracy = &maxAccuracy
	}
	
	return a.GetAllInterviewAnalytics(filters)
}

// GetGlobalAnalyticsAPI calculates aggregated statistics across all interviews
func (a *App) GetGlobalAnalyticsAPI(dateFrom, dateTo string, minAccuracy, maxAccuracy float64) (*GlobalAnalytics, error) {
	filters := &AnalyticsFilters{}
	
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
	
	if minAccuracy > 0 {
		filters.MinAccuracy = &minAccuracy
	}
	
	if maxAccuracy > 0 {
		filters.MaxAccuracy = &maxAccuracy
	}
	
	return a.GetGlobalAnalytics(filters)
}
