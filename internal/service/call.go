package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mrbelka12000/interview_parser/internal/models"
)

// SaveCall creates a new call with transcript and optional analysis
func (s *Service) SaveCall(obj *models.Call) error {
	if obj.Transcript == "" {
		return fmt.Errorf("transcript cannot be empty")
	}

	_, err := s.callRepo.Create(obj)
	if err != nil {
		return fmt.Errorf("failed to create call: %w", err)
	}

	return nil
}

// GetCall retrieves a call by ID
func (s *Service) GetCall(id int64) (*models.Call, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid call ID: %d", id)
	}

	call, err := s.callRepo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get call: %w", err)
	}

	return call, nil
}

// GetAllCalls retrieves all calls with optional pagination
func (s *Service) GetAllCalls(limit, offset int) ([]models.Call, error) {
	if limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative: %d", limit)
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative: %d", offset)
	}

	calls, err := s.callRepo.GetAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all calls: %w", err)
	}

	return calls, nil
}

// UpdateCall updates an existing call
func (s *Service) UpdateCall(id int64, transcript string, analysis interface{}) (*models.Call, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid call ID: %d", id)
	}
	if transcript == "" {
		return nil, fmt.Errorf("transcript cannot be empty")
	}

	// Get existing call first
	existingCall, err := s.callRepo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing call: %w", err)
	}

	// Update fields
	existingCall.Transcript = transcript

	// Convert analysis to JSON if provided
	if analysis != nil {
		analysisJSON, err := json.Marshal(analysis)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal analysis: %w", err)
		}
		existingCall.Analysis = analysisJSON
	} else {
		existingCall.Analysis = nil
	}

	err = s.callRepo.Update(existingCall)
	if err != nil {
		return nil, fmt.Errorf("failed to update call: %w", err)
	}

	return existingCall, nil
}

// DeleteCall deletes a call by ID
func (s *Service) DeleteCall(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid call ID: %d", id)
	}

	err := s.callRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete call: %w", err)
	}

	return nil
}

// GetCallsByDateRange retrieves calls within a specified date range
func (s *Service) GetCallsByDateRange(dateFrom, dateTo time.Time) ([]models.Call, error) {
	if dateFrom.After(dateTo) {
		return nil, fmt.Errorf("dateFrom cannot be after dateTo")
	}

	calls, err := s.callRepo.GetByDateRange(dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get calls by date range: %w", err)
	}

	return calls, nil
}

// UpdateCallAnalysis updates only the analysis field of a call
func (s *Service) UpdateCallAnalysis(id int64, analysis interface{}) error {
	if id <= 0 {
		return fmt.Errorf("invalid call ID: %d", id)
	}

	// Get existing call first
	existingCall, err := s.callRepo.Get(id)
	if err != nil {
		return fmt.Errorf("failed to get existing call: %w", err)
	}

	// Convert analysis to JSON if provided
	if analysis != nil {
		analysisJSON, err := json.Marshal(analysis)
		if err != nil {
			return fmt.Errorf("failed to marshal analysis: %w", err)
		}
		existingCall.Analysis = analysisJSON
	} else {
		existingCall.Analysis = nil
	}

	err = s.callRepo.Update(existingCall)
	if err != nil {
		return fmt.Errorf("failed to update call analysis: %w", err)
	}

	return nil
}
