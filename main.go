package main

import (
	"embed"
	"errors"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/mrbelka12000/interview_parser/internal/app"
	"github.com/mrbelka12000/interview_parser/internal/config"
	"github.com/mrbelka12000/interview_parser/internal/repo"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cfg := config.ParseConfig()
	if cfg == nil {
		fmt.Println("config is nil")
		os.Exit(1)
	}

	if err := repo.InitDB(cfg); err != nil {
		fmt.Println("db init error", err)
		os.Exit(1)
	}

	apiKey, err := repo.NewApiKeyRepo().GetOpenAIAPIKeyFromDB()
	if err != nil && !errors.Is(err, repo.ErrNoKey) {
		fmt.Printf("error getting open AI api key: %v\n", err)
		os.Exit(1)
	}
	cfg.OpenAIAPIKey = apiKey

	app := app.NewApp(cfg)

	// Create application with options
	err = wails.Run(&options.App{
		Title: "interview_parser_app",
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Fullscreen:       true,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}

/*
TODO:
dashboard with metrics
Integrate audio parser
*/
