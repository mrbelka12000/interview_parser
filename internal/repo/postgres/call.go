package postgres

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/mrbelka12000/interview_parser/internal/models"
)

type CallRepo struct{}

func NewCallRepo() *CallRepo {
	return &CallRepo{}
}

// Create creates a new call record
func (r *CallRepo) Create(call *models.Call) (uint64, error) {
	now := time.Now()
	call.CreatedAt = now
	call.UpdatedAt = now

	if err := GetDB().Create(call).Error; err != nil {
		return 0, fmt.Errorf("failed to create call: %w", err)
	}

	return call.ID, nil
}

// Get retrieves a call by ID
func (r *CallRepo) Get(id uint64) (*models.Call, error) {
	var call models.Call
	if err := GetDB().First(&call, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no call found with id: %d", id)
		}
		return nil, fmt.Errorf("failed to retrieve call: %w", err)
	}

	return &call, nil
}

// GetAll retrieves all calls with optional pagination
func (r *CallRepo) GetAll(limit, offset int) ([]models.Call, error) {
	query := GetDB().Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	var calls []models.Call
	if err := query.Find(&calls).Error; err != nil {
		return nil, fmt.Errorf("failed to query calls: %w", err)
	}

	return calls, nil
}

// Update updates an existing call
func (r *CallRepo) Update(call *models.Call) error {
	now := time.Now()
	call.UpdatedAt = now

	result := GetDB().Model(call).Updates(map[string]interface{}{
		"transcript": call.Transcript,
		"analysis":   call.Analysis,
		"updated_at": now,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to update call: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no call found with id: %d", call.ID)
	}

	return nil
}

// Delete deletes a call by ID
func (r *CallRepo) Delete(id uint64) error {
	result := GetDB().Delete(&models.Call{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete call: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no call found with id: %d", id)
	}

	return nil
}

// GetByDateRange retrieves calls within a date range
func (r *CallRepo) GetByDateRange(dateFrom, dateTo time.Time) ([]models.Call, error) {
	var calls []models.Call
	if err := GetDB().Where("created_at >= ? AND created_at <= ?", dateFrom, dateTo).
		Order("created_at DESC").
		Find(&calls).Error; err != nil {
		return nil, fmt.Errorf("failed to query calls by date range: %w", err)
	}

	return calls, nil
}
