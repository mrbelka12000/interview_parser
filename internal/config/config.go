package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	ChunkSeconds              int
	GPTTranscribeModel        string
	GPTClassifyQuestionsModel string
	DBPath                    string
	LoadChunks                bool
	ChunksDir                 string
	ParallelWorkers           int
	OpenAIAPIKey              string
	Language                  string
	DefaultDir                string
	DefaultTranscriptDir      string
	DefaultAnalyzeDir         string
	DefaultOutputName         string
	DefaultAnalyzeCallDir     string

	AudioSampleRate uint32
	AudioChannels   uint32
	AudioBitrate    uint16
}

const (
	defaultDirName                   = ".interview_parser"
	defaultTranscriptDir             = "transcripts"
	defaultAnalyzeDir                = "analyzes"
	defaultAnalyzeCallDir            = "calls"
	defaultChunksDir                 = "output/chunks"
	defaultChunksSeconds             = 100
	defaultOutputName                = "analytics.md"
	defaultGPTClassifyQuestionsModel = "o3"
	defaultGPTTranscribeModels       = "gpt-4o-transcribe"
	defaultLoadChunks                = false
	defaultAudioSampleRate           = 48000
	defaultAudioChannels             = 2
	defaultAudioBitrate              = 16
)

func ParseConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %s\n", err)
		return nil
	}

	defaultDir := filepath.Join(home, defaultDirName)
	if err := os.Mkdir(defaultDir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			fmt.Printf("Failed to create default directory: %s\n", err)
		}
	}

	// Initialize config with default values directly from constants
	cfg := &Config{
		ChunkSeconds:              defaultChunksSeconds,
		GPTTranscribeModel:        defaultGPTTranscribeModels,
		GPTClassifyQuestionsModel: defaultGPTClassifyQuestionsModel,
		LoadChunks:                defaultLoadChunks,
		ParallelWorkers:           runtime.NumCPU(),
		Language:                  "ru",
		DefaultOutputName:         defaultOutputName,
		AudioSampleRate:           defaultAudioSampleRate,
		AudioChannels:             defaultAudioChannels,
		AudioBitrate:              defaultAudioBitrate,
	}

	// Set derived paths
	cfg.DBPath = filepath.Join(defaultDir, "local.db")
	cfg.ChunksDir = filepath.Join(defaultDir, defaultChunksDir)
	cfg.DefaultDir = defaultDir
	cfg.DefaultTranscriptDir = filepath.Join(defaultDir, defaultTranscriptDir)
	cfg.DefaultAnalyzeDir = filepath.Join(defaultDir, defaultAnalyzeDir)
	cfg.DefaultAnalyzeCallDir = filepath.Join(defaultDir, defaultAnalyzeCallDir)

	if err := os.Mkdir(cfg.DefaultTranscriptDir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			fmt.Printf("Failed to create default transcript directory: %s\n", err)
		}
	}
	if err := os.Mkdir(cfg.DefaultAnalyzeDir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			fmt.Printf("Failed to create default analyze directory: %s\n", err)
		}
	}
	if err := os.Mkdir(cfg.DefaultAnalyzeCallDir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			fmt.Printf("Failed to create default analyze calls directory: %s\n", err)
		}
	}

	fmt.Printf("DBPath: %s\n", cfg.DBPath)

	return cfg
}
