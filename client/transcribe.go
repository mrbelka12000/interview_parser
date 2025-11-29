package client

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go"
)

func (c *Client) Transcribe(ctx context.Context, chunkPath string) (string, error) {
	log.Printf("[i] Transcribing media from %s", chunkPath)

	f, err := os.Open(chunkPath)
	if err != nil {
		return "", fmt.Errorf("%v file open error: %w", chunkPath, err)
	}

	defer f.Close()

	res, err := c.cl.Audio.Transcriptions.New(ctx, openai.AudioTranscriptionNewParams{
		Model: c.cfg.GPTTranscribeModel,
		File:  f,
	})
	if err != nil {
		return "", fmt.Errorf("%v transcribe error: %w", chunkPath, err)
	}

	return res.Text, nil
}
