package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/openai/openai-go"

	"github.com/mrbelka12000/interview_parser/internal/models"
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

	CallResponse struct {
		MeetingAnalysis struct {
			KeyTopics []string `json:"key_topics"`
			Tasks     []struct {
				Title    string  `json:"title"`
				Assignee string  `json:"assignee"`
				Deadline *string `json:"deadline"`
			} `json:"tasks"`
			OpenQuestionsAndBlockers []string `json:"open_questions_and_blockers"`
			NextSteps                []string `json:"next_steps"`
		} `json:"meeting_analysis"`
	}
)

func (c *Client) AnalyzeTranscript(ctx context.Context, text string) (out models.AnalyzeInterviewWithQA, err error) {
	now := time.Now()
	log.Printf("[i] Analyzing transcript")

	var resp AnalyzeResponse

	res, err := c.cl.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(promptAnalyze),
			openai.UserMessage(fmt.Sprintf(transcriptHeader, text)),
		},
		Model: c.cfg.GPTClassifyQuestionsModel,
	})
	if err != nil {
		return out, fmt.Errorf("failed to analyze transcript: %w", err)
	}

	if err = json.Unmarshal([]byte(res.Choices[0].Message.Content), &resp); err != nil {
		return out, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	for _, q := range resp.Questions {
		out.QA = append(out.QA, models.QuestionAnswer{
			Question:         q.Question,
			FullAnswer:       q.FullAnswer,
			Accuracy:         q.Accuracy,
			ReasonUnanswered: q.ReasonUnanswered,
		})
	}

	log.Printf("[i] Finished analyzing transcript, seconds spent: %v", time.Since(now).Seconds())
	return out, nil
}

func (c *Client) AnalyzeCall(ctx context.Context, transcript string) (out *models.Call, err error) {
	now := time.Now()
	log.Printf("[i] Analyzing call transcript")

	res, err := c.cl.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(promptCallAnalyze),
			openai.UserMessage(transcript),
		},
		Model: c.cfg.GPTClassifyQuestionsModel,
	})
	if err != nil {
		return out, fmt.Errorf("failed to analyze call: %w", err)
	}

	var resp CallResponse
	if err = json.Unmarshal([]byte(res.Choices[0].Message.Content), &resp); err != nil {
		return out, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Printf("[i] Finished analyzing call, seconds spent: %v", time.Since(now).Seconds())

	body, err := json.Marshal(resp.MeetingAnalysis)
	if err != nil {
		return out, fmt.Errorf("failed to marshal response: %w", err)
	}

	return &models.Call{
		Transcript: transcript,
		Analysis:   body,
	}, nil
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
