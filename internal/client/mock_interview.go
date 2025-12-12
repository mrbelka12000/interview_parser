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
	MockInterviewRequest struct {
		CV             string
		VacancyInfo    string
		Specialization string
		Level          string
		Meta           string
		QuestionsCount int
	}

	MockInterviewResponse struct {
		VacancySummary     string              `json:"vacancy_summary"`
		GeneratedQuestions []GeneratedQuestion `json:"generated_questions"`
	}
	GeneratedQuestion struct {
		Category string `json:"category"`
		Question string `json:"question"`
		WhyAsked string `json:"why_asked"`
	}
)

func (c *Client) GetMockInterviewQuestions(ctx context.Context, req MockInterviewRequest) (out MockInterviewResponse, err error) {
	now := time.Now()
	log.Printf("[i] Generating mock interview")

	res, err := c.cl.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(fmt.Sprintf(promptGenerateMockInterview, req.CV, req.VacancyInfo, req.Specialization, req.Level, req.Meta, req.QuestionsCount)),
		},
		Model: c.cfg.GPTGenerateQuestionsModel,
	})
	if err != nil {
		return out, fmt.Errorf("failed to generate mock interview: %w", err)
	}

	if err = json.Unmarshal([]byte(res.Choices[0].Message.Content), &out); err != nil {
		return out, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Printf("[i] Successfully generated mock interview, spent: %v", time.Since(now).Seconds())

	return out, nil
}
