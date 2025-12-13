package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/mrbelka12000/interview_parser/internal/config"
	wailsapp "github.com/mrbelka12000/interview_parser/internal/wails_app"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cfg := config.ParseConfig()
	if cfg == nil {
		fmt.Println("config is nil")
		os.Exit(1)
	}

	if cfg.ENV == config.ENVLocal {
		app := wailsapp.NewApp(cfg)
		// Create application with options
		err := wails.Run(&options.App{
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

	fmt.Println(cfg.ServiceConfig.ENV, "starting")
}
