package repo

import (
	"time"

	"github.com/mrbelka12000/interview_parser/internal/models"
)

// ApiKeyRepository defines interface for API key operations
type ApiKeyRepository interface {
	GetOpenAIAPIKeyFromDB() (string, error)
	InsertOpenAIAPIKey(openAIAPIKey string) error
	DeleteOpenAIAPIKey() error
}

// InterviewRepository defines interface for interview operations
type InterviewRepository interface {
	Save(interview *models.AnalyzeInterviewWithQA) error
	Get(id uint64) (*models.AnalyzeInterview, []models.QuestionAnswer, error)
	GetAll(filters *models.GetInterviewsFilters) ([]models.AnalyzeInterview, [][]models.QuestionAnswer, error)
	Update(interview *models.AnalyzeInterview, qaList []models.QuestionAnswer) error
	Delete(id uint64) error
}

// CallRepository defines interface for call operations
type CallRepository interface {
	Create(call *models.Call) (uint64, error)
	Get(id uint64) (*models.Call, error)
	GetAll(limit, offset int) ([]models.Call, error)
	Update(call *models.Call) error
	Delete(id uint64) error
	GetByDateRange(dateFrom, dateTo time.Time) ([]models.Call, error)
}
