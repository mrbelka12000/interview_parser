package service

import (
	"github.com/mrbelka12000/interview_parser/internal/repo"
)

type (
	Service struct {
		apiKeyRepo    repo.ApiKeyRepository
		interviewRepo repo.InterviewRepository
		callRepo      repo.CallRepository
	}
)

func New(apiKeyRepo repo.ApiKeyRepository, interviewRepo repo.InterviewRepository, callRepo repo.CallRepository) *Service {
	return &Service{
		apiKeyRepo:    apiKeyRepo,
		interviewRepo: interviewRepo,
		callRepo:      callRepo,
	}
}
