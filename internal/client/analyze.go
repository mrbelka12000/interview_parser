package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/openai/openai-go"
)

type (
	AnalyzeResponse struct {
		Questions []Question `json:"questions"`
	}
	Question struct {
		Question         string  `json:"question"`
		FullAnswer       string  `json:"full_answer"`
		Accuracy         float64 `json:"accuracy"`
		Questioner       string  `json:"questioner"`
		Answerer         string  `json:"answerer"`
		ReasonUnanswered string  `json:"reason"`
	}
	
	CallAnalyzeResponse struct {
		AnalysisText string `json:"analysis_text"`
	}
)

func (c *Client) AnalyzeTranscript(ctx context.Context, text string) (out AnalyzeResponse, err error) {
	now := time.Now()
	log.Printf("[i] Analyzing transcript")

	res, err := c.cl.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(getAnalyzePrompt(c.cfg.Language)),
			openai.UserMessage(fmt.Sprintf(getTranscriptPrompt(c.cfg.Language), text)),
		},
		Model: c.cfg.GPTClassifyQuestionsModel,
	})
	if err != nil {
		return out, fmt.Errorf("failed to analyze transcript: %w", err)
	}

	if err = json.Unmarshal([]byte(res.Choices[0].Message.Content), &out); err != nil {
		return out, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Printf("[i] Finished analyzing transcript, seconds spent: %v", time.Since(now).Seconds())
	return out, nil
}

func (c *Client) AnalyzeCall(ctx context.Context, text string) (out CallAnalyzeResponse, err error) {
	now := time.Now()
	log.Printf("[i] Analyzing call transcript")

	res, err := c.cl.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(promptCallAnalyze),
			openai.UserMessage(text),
		},
		Model: c.cfg.GPTClassifyQuestionsModel,
	})
	if err != nil {
		return out, fmt.Errorf("failed to analyze call: %w", err)
	}

	out.AnalysisText = res.Choices[0].Message.Content

	log.Printf("[i] Finished analyzing call, seconds spent: %v", time.Since(now).Seconds())
	return out, nil
}
