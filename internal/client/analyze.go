package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

	AnalyzeMockInterviewRequest struct {
		CV             string
		VacancyInfo    string
		Specialization string
		Level          string
		Meta           string
		Questions      []string
		Answers        []string
	}

	AnalyzeMockInterviewResponse struct {
		CandidateSummary    string `json:"candidate_summary"`
		EvaluationLevel     string `json:"evaluation_level"`
		QuestionsEvaluation []struct {
			Question         string  `json:"question"`
			Answer           string  `json:"answer"`
			Accuracy         float64 `json:"accuracy"`
			Assessment       string  `json:"assessment"`
			ReasonUnanswered string  `json:"reason_unanswered"`
			WhatWasExpected  string  `json:"what_was_expected"`
		} `json:"questions_evaluation"`
		FinalScore struct {
			AverageAccuracy float64 `json:"average_accuracy"`
			Verdict         string  `json:"verdict"`
			VerdictReason   string  `json:"verdict_reason"`
		} `json:"final_score"`
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

func (c *Client) AnalyzeMockInterview(ctx context.Context, req AnalyzeMockInterviewRequest) (out AnalyzeMockInterviewResponse, err error) {

	now := time.Now()
	log.Printf("[i] Analyzing mock interview")

	res, err := c.cl.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(fmt.Sprintf(promptAnalyzeMockInterview,
				req.CV, req.VacancyInfo, req.Specialization, req.Level,
				req.Meta, strings.Join(req.Questions, "\n"), strings.Join(req.Answers, "\n"))),
		},
		Model: c.cfg.GPTClassifyQuestionsModel,
	})
	if err != nil {
		return out, fmt.Errorf("failed to analyze mock interview: %w", err)
	}

	if err = json.Unmarshal([]byte(res.Choices[0].Message.Content), &out); err != nil {
		return out, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Printf("[i] Finished analyzing mock interview, seconds spent: %v", time.Since(now).Seconds())
	return out, nil
}
