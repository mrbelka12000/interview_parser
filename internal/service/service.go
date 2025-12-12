package service

import "github.com/mrbelka12000/interview_parser/internal/repo"

type Service struct {
	apiKeyRepo    *repo.ApiKeyRepo
	interviewRepo *repo.InterviewRepo
	callRepo      *repo.CallRepo
}

func New() *Service {
	return &Service{
		apiKeyRepo:    repo.NewApiKeyRepo(),
		interviewRepo: repo.NewInterviewRepo(),
		callRepo:      repo.NewCallRepo(),
	}
}
