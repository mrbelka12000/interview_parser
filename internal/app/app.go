package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
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

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %v", err)
	}

	// Read directory contents
	entries, err := os.ReadDir(cwd)
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
			Path:      filepath.Join(cwd, info.Name()),
			IsDir:     entry.IsDir(),
			Size:      info.Size(),
			Extension: ext,
		}

		files = append(files, fileInfo)
	}

	return files, nil
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
