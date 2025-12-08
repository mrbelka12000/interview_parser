package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	ctx    context.Context
	cfg    *config.Config
	parser *parser.Parser
}

// NewApp creates a new App application struct
func NewApp(cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		parser: parser.NewParser(cfg, client.New(cfg)),
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
	var files []FileInfo

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
