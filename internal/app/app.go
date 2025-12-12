package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	audiocapture "github.com/mrbelka12000/interview_parser/internal/audio_capture"
	"github.com/mrbelka12000/interview_parser/internal/client"
	"github.com/mrbelka12000/interview_parser/internal/config"
	"github.com/mrbelka12000/interview_parser/internal/delivery/ws"
	"github.com/mrbelka12000/interview_parser/internal/parser"
	"github.com/mrbelka12000/interview_parser/internal/service"
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
	audioRecorder *audiocapture.AudioCapturer
	service       *service.Service
}

// NewApp creates a new App application struct
func NewApp(cfg *config.Config) *App {
	audioRecorder, err := audiocapture.NewAudioCapturer(cfg.AudioSampleRate, cfg.AudioChannels, cfg.AudioBitrate)
	if err != nil {
		log.Println(fmt.Sprintf("Error creating audio recorder %v", err))
	}

	return &App{
		cfg:           cfg,
		parser:        parser.NewParser(cfg),
		aiClient:      client.New(cfg),
		audioRecorder: audioRecorder,
		service:       service.New(),
	}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	sync.OnceFunc(func() {
		if err := ws.RunServer(a.cfg); err != nil {
			log.Println(fmt.Sprintf("Error starting WS server: %v", err))
		}
	})()
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

// GetWebSocketURL returns the WebSocket server URL for the mock interview
func (a *App) GetWebSocketURL() string {
	return fmt.Sprintf("ws://localhost:%d/ws", a.cfg.WSServerPort)
}
