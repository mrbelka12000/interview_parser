package app

import (
	"fmt"
	"log"
	"sync"

	"github.com/mrbelka12000/interview_parser/internal/client"
)

func (a *App) analyzeInterview(analyzePath, transcript string) error {
	a.sendProgress(78, "Starting analyzing...", "Analyzing transcripts...")
	var (
		wg        sync.WaitGroup
		mx        sync.Mutex
		workers   = make(chan struct{}, a.cfg.ParallelWorkers)
		completed = 0
	)

	transcriptBatches := a.parser.BatchTranscript(transcript)
	analyzeRespBatches := make([][]client.Question, len(transcriptBatches))

	for i, batch := range transcriptBatches {
		wg.Add(1)
		workers <- struct{}{}
		go func(ind int, b string) {
			defer func() {
				<-workers
				wg.Done()
			}()

			analyzeRespTmp, err := a.aiClient.AnalyzeTranscript(a.ctx, b)
			if err != nil {
				log.Printf("Error analyzing transcript %d: %v", ind, err)
				return
			}

			mx.Lock()
			analyzeRespBatches[ind] = analyzeRespTmp.Questions
			completed++
			progress := 85 + int(float64(completed)/float64(len(transcriptBatches))*10) // 85% to 95%

			a.sendProgress(progress, "Analyzing transcript...", fmt.Sprintf("Analyzed %d/%d tranascript segments...", completed, len(transcriptBatches)))
			mx.Unlock()

		}(i, batch)
	}

	wg.Wait()
	close(workers)

	var analyzeResp client.AnalyzeResponse
	for _, batch := range analyzeRespBatches {
		analyzeResp.Questions = append(analyzeResp.Questions, batch...)
	}

	// Step 6: Save analysis response
	a.sendProgress(97, "Saving analysis...", "Writing analysis file...")
	if err := a.parser.SaveAnalyzeResponse(analyzePath, analyzeResp); err != nil {
		return fmt.Errorf("failed to save analysis: %w", err)
	}

	return nil
}

func (a *App) analyzeCall(analysisCallPath, transcript string) error {
	a.sendProgress(75, "Analyzing call...", "Analyzing meeting content...")
	analyzeResp, err := a.aiClient.AnalyzeCall(a.ctx, transcript)
	if err != nil {
		return fmt.Errorf("faield to analyze call: %w", err)
	}

	// Step 5: Save analysis response
	a.sendProgress(90, "Saving analysis...", "Writing analysis file...")
	if err = a.parser.SaveCallAnalysis(analysisCallPath, analyzeResp.AnalysisText); err != nil {
		return fmt.Errorf("failed to save analysis: %w", err)
	}

	return nil
}
