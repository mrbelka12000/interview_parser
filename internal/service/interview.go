package service

import (
	"fmt"

	"github.com/mrbelka12000/interview_parser/internal/models"
)

// SaveInterview creates a new interview with its question answers
func (s *Service) SaveInterview(interview *models.AnalyzeInterviewWithQA) error {
	if interview == nil {
		return fmt.Errorf("interview cannot be nil")
	}
	if len(interview.QA) == 0 {
		return fmt.Errorf("question answers list cannot be empty")
	}

	// Validate question answers
	for i, qa := range interview.QA {
		if qa.Question == "" {
			return fmt.Errorf("question at index %d cannot be empty", i)
		}
		if qa.Accuracy < 0 || qa.Accuracy > 100 {
			return fmt.Errorf("accuracy at index %d must be between 0 and 100", i)
		}
	}

	return s.interviewRepo.Save(interview)
}

// UpdateInterview updates an interview and its question answers
func (s *Service) UpdateInterview(interview *models.AnalyzeInterview, qaList []models.QuestionAnswer) error {
	if interview == nil {
		return fmt.Errorf("interview cannot be nil")
	}
	if interview.ID == 0 {
		return fmt.Errorf("invalid interview ID: %d", interview.ID)
	}
	if len(qaList) == 0 {
		return fmt.Errorf("question answers list cannot be empty")
	}

	// Validate question answers
	for i, qa := range qaList {
		if qa.Question == "" {
			return fmt.Errorf("question at index %d cannot be empty", i)
		}
		if qa.Accuracy < 0 || qa.Accuracy > 100 {
			return fmt.Errorf("accuracy at index %d must be between 0 and 100", i)
		}
	}

	return s.interviewRepo.Update(interview, qaList)
}

// DeleteInterview deletes an interview and its question answers
func (s *Service) DeleteInterview(id uint64) error {
	if id == 0 {
		return fmt.Errorf("invalid interview ID: %d", id)
	}

	return s.interviewRepo.Delete(id)
}

// GetInterview combines interview and question answers into a single struct
func (s *Service) GetInterview(id uint64) (*models.AnalyzeInterviewWithQA, error) {
	interview, qaList, err := s.interviewRepo.Get(id)
	if err != nil {
		return nil, err
	}

	return &models.AnalyzeInterviewWithQA{
		ID:        interview.ID,
		QA:        qaList,
		CreatedAt: interview.CreatedAt,
		UpdatedAt: interview.UpdatedAt,
	}, nil
}

// GetAllInterviews retrieves all interviews with their question answers combined
func (s *Service) GetAllInterviews(filters *models.GetInterviewsFilters) ([]models.AnalyzeInterviewWithQA, error) {
	interviews, qaLists, err := s.interviewRepo.GetAll(filters)
	if err != nil {
		return nil, err
	}

	var result []models.AnalyzeInterviewWithQA
	for i, interview := range interviews {
		result = append(result, models.AnalyzeInterviewWithQA{
			ID:        interview.ID,
			QA:        qaLists[i],
			CreatedAt: interview.CreatedAt,
			UpdatedAt: interview.UpdatedAt,
		})
	}

	return result, nil
}
