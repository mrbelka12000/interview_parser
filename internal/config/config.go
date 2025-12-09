package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type Config struct {
	ChunkSeconds              int    `mapstructure:"chunk_seconds"`
	GPTTranscribeModel        string `mapstructure:"gpt_transcribe_model"`
	GPTClassifyQuestionsModel string `mapstructure:"gpt_classify_questions_model"`
	DBPath                    string `mapstructure:"-"`
	LoadChunks                bool   `mapstructure:"load_chunks"`
	ChunksDir                 string `mapstructure:"-"`
	ParallelWorkers           int    `mapstructure:"parallel_workers"`
	OpenAIAPIKey              string `mapstructure:"openai_api_key"`
	Language                  string `mapstructure:"language"`
	DefaultDir                string `mapstructure:"-"`
	DefaultTranscriptDir      string `mapstructure:"-"`
	DefaultAnalyzeDir         string `mapstructure:"-"`
}

const (
	defaultDirName                   = ".interview_parser"
	defaultTranscriptDir             = "transcripts"
	defaultAnalyzeDir                = "analyzes"
	defaultChunksDir                 = "output/chunks"
	defaultChunksSeconds             = 100
	defaultOutputName                = "analytics.md"
	defaultGPTClassifyQuestionsModel = "o3"
	defaultGPTTranscribeModels       = "gpt-4o-transcribe"
	defaultLoadChunks                = false
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

	pflag.Usage = func() {
		fmt.Println("Interview Parser - A tool for transcribing and analyzing interview recordings")
		fmt.Println("")
		fmt.Println("DESCRIPTION:")
		fmt.Println("  Interview Parser processes audio/video interview recordings to transcribe")
		fmt.Println("  them and analyze the questions asked by interviewers. It extracts all meaningful")
		fmt.Println("  questions from the interviewer and classifies them as answered or unanswered")
		fmt.Println("  by the candidate, providing detailed analytics about the interview.")
		fmt.Println("")
		fmt.Println("USAGE:")
		fmt.Println("  interview-parser [flags]")
		fmt.Println("")
		fmt.Println("EXAMPLES:")
		fmt.Println("  # Process a new interview recording")
		fmt.Println("  interview-parser --input interview.mov")
		fmt.Println("")
		fmt.Println("  # Use custom chunk size and output file")
		fmt.Println("  interview-parser --input interview.mov --chunk_seconds 120 --output analysis.md")
		fmt.Println("")
		fmt.Println("  # Load previously created chunks instead of creating new ones")
		fmt.Println("  interview-parser --load_chunks")
		fmt.Println("")
		fmt.Println("  # Use different AI models")
		fmt.Println("  interview-parser --input interview.mov --gpt_transcribe_model whisper-1 --gpt_classify_questions_model gpt-4")
		fmt.Println("")
		fmt.Println("  # Specify language for analysis (ru or en)")
		fmt.Println("  interview-parser --input interview.mov --language en")
		fmt.Println("")
		fmt.Println("FLAGS:")
		pflag.PrintDefaults()
		fmt.Println("")
		fmt.Println("ENVIRONMENT VARIABLES:")
		fmt.Println("  OPENAI_API_KEY    OpenAI API key (can also be set via --openai_api_key flag)")
		fmt.Println("")
		fmt.Println("REQUIREMENTS:")
		fmt.Println("  - ffmpeg and ffprobe must be installed and available in PATH")
		fmt.Println("  - Valid OpenAI API key with access to the specified models")
		fmt.Println("")
		fmt.Println("OUTPUT:")
		fmt.Println("  The tool generates two files:")
		fmt.Println("  1. Transcript file (saved to ~/.interview_parser/transcript.txt)")
		fmt.Println("  2. Analysis file (saved to specified output path, default: analytics.md)")
		fmt.Println("")
		os.Exit(0)
	}

	pflag.String("input", "", "Path to input audio/video file to transcribe and analyze")
	pflag.String("output", defaultOutputName, "Path to output file for analysis results")
	pflag.Int("chunk_seconds", defaultChunksSeconds, "Duration in seconds for each audio chunk during processing")
	pflag.String("gpt_transcribe_model", defaultGPTTranscribeModels, "OpenAI model to use for audio transcription")
	pflag.String("gpt_classify_questions_model", defaultGPTClassifyQuestionsModel, "OpenAI model to use for question analysis")
	pflag.Bool("load_chunks", defaultLoadChunks, "Load previously created chunks instead of creating new ones")
	pflag.Int("parallel_workers", runtime.NumCPU(), "Number of parallel workers for processing")
	pflag.String("openai_api_key", "", "OpenAI API key (can also be set via OPENAI_API_KEY environment variable)")
	pflag.String("language", "ru", "Language to use for analysis results")

	pflag.Parse()

	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		fmt.Printf("Error parsing flags: %s\n", err)
		return nil
	}

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("Error unmarshalling to config: %s\n", err)
		return nil
	}

	if err := gotenv.Load(".env"); err == nil {
		cfg.OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
	}

	cfg.DBPath = filepath.Join(defaultDir, "local.db")
	cfg.ChunksDir = filepath.Join(defaultDir, defaultChunksDir)
	cfg.DefaultDir = defaultDir
	cfg.DefaultTranscriptDir = filepath.Join(defaultDir, defaultTranscriptDir)
	cfg.DefaultAnalyzeDir = filepath.Join(defaultDir, defaultAnalyzeDir)

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

	fmt.Printf("DBPath: %s\n", cfg.DBPath)

	return &cfg
}
