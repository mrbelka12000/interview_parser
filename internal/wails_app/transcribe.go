package wails_app

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

func (a *App) transcribeFile(filePath string) (string, error) {
	a.sendProgress(15, "Splitting into chunks...", "Dividing file into manageable segments...")
	chunks, err := a.parser.SplitIntoChunks(a.cfg, filePath)
	if err != nil {
		return "", fmt.Errorf("failed to split into chunks: %w", err)
	}

	if len(chunks) == 0 {
		return "", fmt.Errorf("no chunks to process")
	}

	// Step 2: Transcribe chunks using the provided parser logic
	a.sendProgress(25, "Transcribing audio...", "Converting speech to text using AI...")

	var (
		wg            sync.WaitGroup
		mx            sync.Mutex
		workers       = make(chan struct{}, a.cfg.ParallelWorkers)
		collectedText = make([]string, len(chunks))
		completed     = 0
	)

	for i, chunk := range chunks {
		wg.Add(1)
		workers <- struct{}{}
		go func(ind int, chunkVar string) {
			defer func() {
				<-workers
				wg.Done()
			}()

			textFromChunk, err := a.aiClient.Transcribe(a.ctx, chunkVar)
			if err != nil {
				log.Printf("Error transcribing chunk %d: %v", ind, err)
				return
			}

			mx.Lock()
			collectedText[ind] = textFromChunk
			completed++
			progress := 25 + int(float64(completed)/float64(len(chunks))*35) // 25% to 60%
			a.sendProgress(progress, "Processing chunks...", fmt.Sprintf("Processed %d/%d audio segments...", completed, len(chunks)))
			mx.Unlock()

		}(i, chunk)
	}

	wg.Wait()

	return strings.Join(collectedText, ""), nil
}
