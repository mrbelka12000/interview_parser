package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type (
	Config struct {
		ServiceConfig
		WSConfig
		GPTConfig
		TranscribeConfig
		LocalConfig
		DBConfig
		AudioConfig
	}

	ServiceConfig struct {
		ServiceName     string `env:"SERVICE_NAME"`
		ParallelWorkers int    `env:"PARALLEL_WORKERS, default=6"`
	}

	WSConfig struct {
		WSServerPort int `env:"WS_SERVER_PORT, default=35044"`
	}

	GPTConfig struct {
		GPTTranscribeModel        string `env:"GPT_TRANSCRIBE_MODEL,required"`
		GPTClassifyQuestionsModel string `env:"GPT_CLASSIFY_QUESTIONS_MODEL,required"`
		GPTGenerateQuestionsModel string `env:"GPT_GENERATE_QUESTIONS_MODEL,required"`
		OpenAIAPIKey              string `env:"OPENAI_API_KEY,required"`
	}

	TranscribeConfig struct {
		ChunkSeconds int `env:"CHUNK_SECONDS, default=100"`
	}

	LocalConfig struct {
		DefaultDir            string
		DefaultTranscriptDir  string
		DefaultAnalyzeDir     string
		DefaultAnalyzeCallDir string
		ChunksDir             string
	}

	DBConfig struct {
		Path  string `env:"DB_PATH"`
		PGURL string `env:"PG_URL"`
	}

	AudioConfig struct {
		AudioSampleRate uint32 `env:"AUDIO_SAMPLE_RATE, default=48000"`
		AudioChannels   uint32 `env:"AUDIO_CHANNELS, default=2"`
		AudioBitrate    uint16 `env:"AUDIO_BITRATE, default=16"`
	}
)

const (
	defaultDirName                   = ".interview_parser"
	defaultTranscriptDir             = "transcripts"
	defaultAnalyzeDir                = "analyzes"
	defaultAnalyzeCallDir            = "calls"
	defaultChunksDir                 = "output/chunks"
	defaultChunksSeconds             = 200
	defaultGPTClassifyQuestionsModel = "o3"
	defaultGPTGenerateQuestionsModel = "gpt-4.1"
	defaultGPTTranscribeModels       = "gpt-4o-transcribe"
	defaultLoadChunks                = false
	defaultAudioSampleRate           = 48000
	defaultAudioChannels             = 2
	defaultAudioBitrate              = 16
	defaultWSServerPort              = 35044

	productionEnv = "PRODUCTION"
	localEnv      = "LOCAL"
)

func ParseConfig() *Config {
	godotenv.Load()

	env := os.Getenv("ENV")
	if env == productionEnv {
		log.Println("[I] Production environment variable detected")
		cfg := &Config{}

		err := envconfig.Process(context.Background(), &cfg)
		if err != nil {
			log.Printf("[I] Error parsing env vars: %v\n", err)
			return nil
		}

		return cfg
	}

	log.Println("[I] Using local environment variable")

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

	// Initialize config with default values using nested structs
	cfg := &Config{
		WSConfig: WSConfig{
			WSServerPort: defaultWSServerPort,
		},
		GPTConfig: GPTConfig{
			GPTTranscribeModel:        defaultGPTTranscribeModels,
			GPTClassifyQuestionsModel: defaultGPTClassifyQuestionsModel,
			GPTGenerateQuestionsModel: defaultGPTGenerateQuestionsModel,
			OpenAIAPIKey:              os.Getenv("OPENAI_API_KEY"),
		},
		TranscribeConfig: TranscribeConfig{
			ChunkSeconds: defaultChunksSeconds,
		},
		LocalConfig: LocalConfig{
			DefaultDir:            defaultDir,
			DefaultTranscriptDir:  filepath.Join(defaultDir, defaultTranscriptDir),
			DefaultAnalyzeDir:     filepath.Join(defaultDir, defaultAnalyzeDir),
			DefaultAnalyzeCallDir: filepath.Join(defaultDir, defaultAnalyzeCallDir),
			ChunksDir:             filepath.Join(defaultDir, defaultChunksDir),
		},
		DBConfig: DBConfig{
			Path: filepath.Join(defaultDir, "local.db"),
		},
		AudioConfig: AudioConfig{
			AudioSampleRate: defaultAudioSampleRate,
			AudioChannels:   defaultAudioChannels,
			AudioBitrate:    defaultAudioBitrate,
		},
	}

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

	fmt.Printf("DBPath: %s\n", cfg.DBConfig.Path)

	return cfg
}
