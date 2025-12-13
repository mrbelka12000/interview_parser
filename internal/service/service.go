package service

import (
	"github.com/mrbelka12000/interview_parser/internal/repo/postgres"
)

type (
	Service struct {
		apiKeyRepo    *postgres.ApiKeyRepo
		interviewRepo *postgres.InterviewRepo
		callRepo      *postgres.CallRepo
	}
)

func New() *Service {
	return &Service{
		apiKeyRepo:    postgres.NewApiKeyRepo(),
		interviewRepo: postgres.NewInterviewRepo(),
		callRepo:      postgres.NewCallRepo(),
	}
}
