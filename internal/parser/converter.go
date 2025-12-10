package parser

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mrbelka12000/interview_parser/internal/config"
)

// getDuration returns media duration in seconds using ffprobe.
func getDuration(mediaPath string) (float64, error) {
	cmd := exec.Command(
		ffprobePath(),
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		mediaPath,
	)

	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// trim пробелы/переводы строк
	s := strings.TrimSpace(string(out))
	dur, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return dur, nil
}

// SplitIntoChunks splits audio file into N-second chunks using ffmpeg.
// Produces .m4a chunks inside output/chunks.
func (p *Parser) SplitIntoChunks(cfg *config.Config, inputPath string) ([]string, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	duration, err := getDuration(inputPath)
	if err != nil {
		return nil, err
	}
	log.Printf("[i] Media duration: %.2fs\n", duration)

	if err := os.MkdirAll(cfg.ChunksDir, 0o755); err != nil {
		return nil, err
	}

	var (
		chunkPaths []string
		start      float64
		idx        int
	)

	ext, codec := outputFormatFor(inputPath)
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))

	for start < duration {
		outPath := filepath.Join(cfg.ChunksDir, fmt.Sprintf("%s_chunk_%03d%s", base, idx, ext))

		cmd := exec.Command(
			ffmpegPath(),
			"-loglevel", "error",
			"-y",
			"-ss", fmt.Sprintf("%.2f", start),
			"-t", strconv.Itoa(cfg.ChunkSeconds),
			"-i", inputPath,
			"-vn",
			"-acodec", codec,
			outPath,
		)

		log.Printf("[i] Creating chunk %d: start=%.2fs → %s\n", idx, start, filepath.Base(outPath))

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return nil, err
		}

		chunkPaths = append(chunkPaths, outPath)
		idx++
		start += float64(cfg.ChunkSeconds)
	}

	log.Printf("[i] Total chunks: %d\n", len(chunkPaths))
	return chunkPaths, nil
}

func (p *Parser) LoadChunks(cfg *config.Config) ([]string, error) {
	dir, err := os.ReadDir(cfg.ChunksDir)
	if err != nil {
		return nil, fmt.Errorf("read dir error: %w", err)
	}

	chunkPaths := make([]string, 0, len(dir))
	for _, de := range dir {
		chunkPaths = append(chunkPaths, filepath.Join(cfg.ChunksDir, de.Name()))
	}

	return chunkPaths, nil
}

func bundledBinary(name string) string {
	exe, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(exe)
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path // bundled version
		}
	}
	return name // fallback on PATH for dev mode
}

func ffprobePath() string { return bundledBinary("ffprobe") }
func ffmpegPath() string  { return bundledBinary("ffmpeg") }

func outputFormatFor(inputPath string) (ext, codec string) {
	switch strings.ToLower(filepath.Ext(inputPath)) {
	case ".wav":
		return ".wav", "pcm_s16le"
	default:
		return ".m4a", "aac"
	}
}
