package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mrbelka12000/interview_parser/internal/models"
)

type CallRepo struct{}

func NewCallRepo() *CallRepo {
	return &CallRepo{}
}

// Create creates a new call record
func (r *CallRepo) Create(call *models.Call) (uint64, error) {
	now := time.Now()
	query := `
	INSERT INTO calls (transcript, analysis, created_at, updated_at) 
	VALUES (?, ?, ?, ?)
	`

	var (
		analysisJSON []byte
		err          error
	)
	if call.Analysis != nil {
		analysisJSON = call.Analysis
	} else {
		analysisJSON = []byte("null")
	}

	result, err := db.Exec(query, call.Transcript, analysisJSON, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to insert call: %w", err)
	}

	callID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get call ID: %w", err)
	}

	call.ID = uint64(callID)
	call.CreatedAt = now
	call.UpdatedAt = now

	return call.ID, nil
}

// Get retrieves a call by ID
func (r *CallRepo) Get(id uint64) (*models.Call, error) {
	query := `
	SELECT id, transcript, analysis, created_at, updated_at 
	FROM calls 
	WHERE id = ?
	`

	var call models.Call
	var analysisJSON []byte

	err := db.QueryRow(query, id).Scan(&call.ID, &call.Transcript, &analysisJSON, &call.CreatedAt, &call.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no call found with id: %d", id)
		}
		return nil, fmt.Errorf("failed to retrieve call: %w", err)
	}

	return &call, nil
}

// GetAll retrieves all calls with optional pagination
func (r *CallRepo) GetAll(limit, offset int) ([]models.Call, error) {
	query := `
	SELECT id, transcript, analysis, created_at, updated_at 
	FROM calls 
	ORDER BY created_at DESC
	`

	args := []interface{}{}
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}
	if offset > 0 {
		query += " OFFSET ?"
		args = append(args, offset)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query calls: %w", err)
	}
	defer rows.Close()

	var calls []models.Call
	for rows.Next() {
		var call models.Call
		var analysisJSON []byte

		err := rows.Scan(&call.ID, &call.Transcript, &analysisJSON, &call.CreatedAt, &call.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan call row: %w", err)
		}

		// Parse JSON analysis
		if analysisJSON != nil {
			call.Analysis = json.RawMessage(analysisJSON)
		}

		calls = append(calls, call)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate call rows: %w", err)
	}

	return calls, nil
}

// Update updates an existing call
func (r *CallRepo) Update(call *models.Call) error {
	now := time.Now()
	query := `
	UPDATE calls 
	SET transcript = ?, analysis = ?, updated_at = ? 
	WHERE id = ?
	`

	var analysisJSON []byte
	if call.Analysis != nil {
		analysisJSON = call.Analysis
	} else {
		analysisJSON = []byte("null")
	}

	result, err := db.Exec(query, call.Transcript, analysisJSON, now, call.ID)
	if err != nil {
		return fmt.Errorf("failed to update call: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no call found with id: %d", call.ID)
	}

	call.UpdatedAt = now
	return nil
}

// Delete deletes a call by ID
func (r *CallRepo) Delete(id uint64) error {
	query := `DELETE FROM calls WHERE id = ?`
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete call: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no call found with id: %d", id)
	}

	return nil
}

// GetByDateRange retrieves calls within a date range
func (r *CallRepo) GetByDateRange(dateFrom, dateTo time.Time) ([]models.Call, error) {
	query := `
	SELECT id, transcript, analysis, created_at, updated_at 
	FROM calls 
	WHERE created_at >= ? AND created_at <= ?
	ORDER BY created_at DESC
	`

	rows, err := db.Query(query, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to query calls by date range: %w", err)
	}
	defer rows.Close()

	var calls []models.Call
	for rows.Next() {
		var call models.Call
		var analysisJSON []byte

		err := rows.Scan(&call.ID, &call.Transcript, &analysisJSON, &call.CreatedAt, &call.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan call row: %w", err)
		}

		// Parse JSON analysis
		if analysisJSON != nil {
			call.Analysis = json.RawMessage(analysisJSON)
		}

		calls = append(calls, call)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate call rows: %w", err)
	}

	return calls, nil
}
