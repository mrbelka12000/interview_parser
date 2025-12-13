package postgres

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/mrbelka12000/interview_parser/internal/models"
)

type InterviewRepo struct{}

func NewInterviewRepo() *InterviewRepo {
	return &InterviewRepo{}
}

// Save creates a new interview and its question answers
func (r *InterviewRepo) Save(interview *models.AnalyzeInterviewWithQA) error {
	return GetDB().Transaction(func(tx *gorm.DB) error {
		// Create interview
		interviewModel := &models.AnalyzeInterview{
			CreatedAt: interview.CreatedAt,
			UpdatedAt: interview.UpdatedAt,
		}

		if err := tx.Create(interviewModel).Error; err != nil {
			return fmt.Errorf("failed to create interview: %w", err)
		}

		// Create question answers
		for i := range interview.QA {
			qa := &interview.QA[i]
			qa.InterviewID = interviewModel.ID
			if err := tx.Create(qa).Error; err != nil {
				return fmt.Errorf("failed to create question answer: %w", err)
			}
		}

		// Update the interview with the generated ID
		interview.ID = interviewModel.ID
		return nil
	})
}

// Get retrieves an interview with its question answers by ID
func (r *InterviewRepo) Get(id uint64) (*models.AnalyzeInterview, []models.QuestionAnswer, error) {
	var interview models.AnalyzeInterview
	if err := GetDB().First(&interview, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("no interview found with id: %d", id)
		}
		return nil, nil, fmt.Errorf("failed to retrieve interview: %w", err)
	}

	var qaList []models.QuestionAnswer
	if err := GetDB().Where("interview_id = ?", id).Order("id").Find(&qaList).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve question answers: %w", err)
	}

	return &interview, qaList, nil
}

// GetAll retrieves all interviews with their question answers
func (r *InterviewRepo) GetAll(filters *models.GetInterviewsFilters) ([]models.AnalyzeInterview, [][]models.QuestionAnswer, error) {
	query := GetDB().Model(&models.AnalyzeInterview{})

	if filters != nil {
		if filters.DateFrom != nil {
			query = query.Where("created_at >= ?", *filters.DateFrom)
		}
		if filters.DateTo != nil {
			query = query.Where("created_at <= ?", *filters.DateTo)
		}
	}

	var interviews []models.AnalyzeInterview
	if err := query.Order("created_at DESC").Find(&interviews).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to query interviews: %w", err)
	}

	var allQALists [][]models.QuestionAnswer
	for _, interview := range interviews {
		qaList, err := r.getQuestionAnswersByInterviewID(interview.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get question answers for interview %d: %w", interview.ID, err)
		}
		allQALists = append(allQALists, qaList)
	}

	return interviews, allQALists, nil
}

// Update updates an interview and its question answers
func (r *InterviewRepo) Update(interview *models.AnalyzeInterview, qaList []models.QuestionAnswer) error {
	return GetDB().Transaction(func(tx *gorm.DB) error {
		// Update interview
		if err := tx.Model(interview).Update("updated_at", interview.UpdatedAt).Error; err != nil {
			return fmt.Errorf("failed to update interview: %w", err)
		}

		// Delete existing question answers
		if err := tx.Where("interview_id = ?", interview.ID).Delete(&models.QuestionAnswer{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing question answers: %w", err)
		}

		// Create new question answers
		for i := range qaList {
			qa := &qaList[i]
			qa.InterviewID = interview.ID
			if err := tx.Create(qa).Error; err != nil {
				return fmt.Errorf("failed to create question answer: %w", err)
			}
		}

		return nil
	})
}

// Delete deletes an interview and its question answers
func (r *InterviewRepo) Delete(id uint64) error {
	result := GetDB().Delete(&models.AnalyzeInterview{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete interview: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no interview found with id: %d", id)
	}

	return nil
}

// Helper method to get question answers by interview ID
func (r *InterviewRepo) getQuestionAnswersByInterviewID(interviewID uint64) ([]models.QuestionAnswer, error) {
	var qaList []models.QuestionAnswer
	if err := GetDB().Where("interview_id = ?", interviewID).Order("id").Find(&qaList).Error; err != nil {
		return nil, fmt.Errorf("failed to query question answers: %w", err)
	}
	return qaList, nil
}
