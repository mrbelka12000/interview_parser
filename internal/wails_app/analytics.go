package wails_app

import (
	"time"

	"github.com/mrbelka12000/interview_parser/internal/models"
)

// GetInterviewAnalyticsAPI retrieves analytics for a specific interview
func (a *App) GetInterviewAnalyticsAPI(interviewID uint64) (models.InterviewAnalytics, error) {
	return a.service.GetInterviewAnalytics(interviewID)
}

// GetAllInterviewAnalyticsAPI retrieves all interview analytics with optional filters
func (a *App) GetAllInterviewAnalyticsAPI(dateFrom, dateTo string) ([]models.InterviewAnalytics, error) {
	filters := &models.GetInterviewsFilters{}

	if dateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters.DateFrom = &parsed
		}
	}

	if dateTo != "" {
		if parsed, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters.DateTo = &parsed
		}
	}

	return a.service.GetAllInterviewAnalytics(filters)
}

// GetGlobalAnalyticsAPI calculates aggregated statistics across all interviews
func (a *App) GetGlobalAnalyticsAPI(dateFrom, dateTo string) (*models.GlobalAnalytics, error) {
	filters := &models.GetInterviewsFilters{}

	if dateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters.DateFrom = &parsed
		}
	}

	if dateTo != "" {
		if parsed, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters.DateTo = &parsed
		}
	}

	return a.service.GetGlobalAnalytics(filters)
}
