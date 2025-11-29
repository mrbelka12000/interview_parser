package interview_parser

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mrbelka12000/interview_parser/config"
)

// getDuration returns media duration in seconds using ffprobe.
func getDuration(mediaPath string) (float64, error) {
	cmd := exec.Command(
		"ffprobe",
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
func SplitIntoChunks(cfg *config.Config) ([]string, error) {
	file, err := os.Open(cfg.InputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	duration, err := getDuration(cfg.InputPath)
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

	base := strings.TrimSuffix(filepath.Base(cfg.InputPath), filepath.Ext(cfg.InputPath))

	for start < duration {
		outPath := filepath.Join(cfg.ChunksDir, fmt.Sprintf("%s_chunk_%03d.m4a", base, idx))

		cmd := exec.Command(
			"ffmpeg",
			"-loglevel", "error",
			"-y",
			"-ss", fmt.Sprintf("%.2f", start),
			"-t", strconv.Itoa(cfg.ChunkSeconds),
			"-i", cfg.InputPath,
			"-vn",
			"-acodec", "aac",
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

func LoadChunks(cfg *config.Config) ([]string, error) {
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
