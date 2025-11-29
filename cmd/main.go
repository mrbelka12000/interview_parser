package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/mrbelka12000/interview_parser"
	"github.com/mrbelka12000/interview_parser/client"
	"github.com/mrbelka12000/interview_parser/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("Received Ctrl+C, shutting down...")
		cancel()
	}()

	cfg := config.ParseConfig()
	if cfg == nil {
		fmt.Println("config is nil")
		os.Exit(1)
	}

	if cfg.OpenAIAPIKey == "" {
		log.Println("OpenAI API key is empty, getting from DB...")
		apiKey, err := interview_parser.GetOpenAIAPIKeyFromDB(cfg)
		if err != nil {
			if errors.Is(err, interview_parser.ErrNoKey) {
				log.Fatal("No OPENAI API key found in DB")
			}
			fmt.Println(err)
			os.Exit(1)
		}

		log.Println("OpenAI API key found, using it...")
		cfg.OpenAIAPIKey = apiKey
	} else {
		apiKey, err := interview_parser.GetOpenAIAPIKeyFromDB(cfg)
		if err != nil {
			if !errors.Is(err, interview_parser.ErrNoKey) {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if apiKey == cfg.OpenAIAPIKey {
			fmt.Println("provided key exists, skipping creating new one...")
			goto skipInsert
		}

		err = interview_parser.InsertOpenAIAPIKey(cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
skipInsert:

	openAIClient := client.New(cfg)
	if err := openAIClient.IsValidAPIKeysProvided(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var (
		chunks []string
		err    error
	)
	if !cfg.LoadChunks {
		chunks, err = interview_parser.SplitIntoChunks(cfg)
	} else {
		chunks, err = interview_parser.LoadChunks(cfg)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		wg            sync.WaitGroup
		mx            sync.Mutex
		workers       = make(chan struct{}, cfg.ParallelWorkers)
		collectedText = make([]string, len(chunks))
	)

	for i, chunk := range chunks {
		wg.Add(1)
		workers <- struct{}{}
		go func(ind int) {
			defer func() {
				<-workers
				wg.Done()
			}()

			textFromChunk, err := openAIClient.Transcribe(ctx, chunk)
			if err != nil {
				fmt.Println(err)
				return
			}

			mx.Lock()
			collectedText[ind] = textFromChunk
			mx.Unlock()
		}(i)
	}

	wg.Wait()
	close(workers)
	if errors.Is(ctx.Err(), context.Canceled) {
		log.Println("[i] cancelled by signal, skip analyze")
		return
	}

	transcript := strings.Join(collectedText, "")
	transcript = interview_parser.FormatText(transcript)

	if err = interview_parser.SaveTranscript(cfg.TranscriptPath, transcript); err != nil {
		fmt.Printf("[i] Failed to save transcript: %s\n", err)
	}

	analyzeResp, err := openAIClient.AnalyzeTranscript(ctx, transcript)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = interview_parser.SaveAnalyzeResponse(cfg.OutputPath, analyzeResp); err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("[i] Successfully saved and analyzed transcript: %s", cfg.OutputPath)
}
