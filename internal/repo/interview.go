package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mrbelka12000/interview_parser/internal/models"
)

type InterviewRepo struct{}

func NewInterviewRepo() *InterviewRepo {
	return &InterviewRepo{}
}

// Save creates a new interview and its question answers
func (r *InterviewRepo) Save(interview *models.AnalyzeInterviewWithQA) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert interview
	now := time.Now()
	query := `
	INSERT INTO analyze_interviews (created_at, updated_at) 
	VALUES (?, ?)
	`
	result, err := tx.Exec(query, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to insert interview: %w", err)
	}

	interviewID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get interview ID: %w", err)
	}

	// Insert question answers
	for i := range interview.QA {
		qa := &interview.QA[i]
		qa.InterviewID = interviewID
		qa.CreatedAt = now
		qa.UpdatedAt = now

		qaQuery := `
		INSERT INTO question_answers (interview_id, question, full_answer, accuracy, reason_unanswered, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		_, err = tx.Exec(qaQuery, qa.InterviewID, qa.Question, qa.FullAnswer, qa.Accuracy, qa.ReasonUnanswered, qa.CreatedAt, qa.UpdatedAt)
		if err != nil {
			return 0, fmt.Errorf("failed to insert question answer: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	interview.ID = interviewID
	interview.CreatedAt = now
	interview.UpdatedAt = now

	return interviewID, nil
}

// Get retrieves an interview with its question answers by ID
func (r *InterviewRepo) Get(id int64) (*models.AnalyzeInterview, []models.QuestionAnswer, error) {
	// Get interview
	query := `
	SELECT id, created_at, updated_at 
	FROM analyze_interviews 
	WHERE id = ?
	`
	var interview models.AnalyzeInterview
	err := db.QueryRow(query, id).Scan(&interview.ID, &interview.CreatedAt, &interview.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no interview found with id: %d", id)
		}
		return nil, nil, fmt.Errorf("failed to retrieve interview: %w", err)
	}

	// Get question answers
	qaQuery := `
	SELECT id, interview_id, question, full_answer, accuracy, reason_unanswered, created_at, updated_at
	FROM question_answers 
	WHERE interview_id = ?
	ORDER BY id
	`
	rows, err := db.Query(qaQuery, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query question answers: %w", err)
	}
	defer rows.Close()

	var qaList []models.QuestionAnswer
	for rows.Next() {
		var qa models.QuestionAnswer
		err := rows.Scan(&qa.ID, &qa.InterviewID, &qa.Question, &qa.FullAnswer, &qa.Accuracy, &qa.ReasonUnanswered, &qa.CreatedAt, &qa.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan question answer row: %w", err)
		}
		qaList = append(qaList, qa)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to iterate question answer rows: %w", err)
	}

	return &interview, qaList, nil
}

// GetAll retrieves all interviews with their question answers
func (r *InterviewRepo) GetAll(filters *models.GetInterviewsFilters) ([]models.AnalyzeInterview, [][]models.QuestionAnswer, error) {
	// Build query with filters
	query := `
	SELECT id, created_at, updated_at 
	FROM analyze_interviews 
	WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if filters != nil {
		if filters.DateFrom != nil {
			query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
			args = append(args, filters.DateFrom)
			argIndex++
		}
		if filters.DateTo != nil {
			query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
			args = append(args, filters.DateTo)
			argIndex++
		}
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query interviews: %w", err)
	}
	defer rows.Close()

	var interviews []models.AnalyzeInterview
	var interviewIDs []int64
	for rows.Next() {
		var interview models.AnalyzeInterview
		err := rows.Scan(&interview.ID, &interview.CreatedAt, &interview.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan interview row: %w", err)
		}
		interviews = append(interviews, interview)
		interviewIDs = append(interviewIDs, interview.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to iterate interview rows: %w", err)
	}

	// Get question answers for all interviews
	var allQALists [][]models.QuestionAnswer
	for _, interviewID := range interviewIDs {
		qaList, err := r.getQuestionAnswersByInterviewID(interviewID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get question answers for interview %d: %w", interviewID, err)
		}
		allQALists = append(allQALists, qaList)
	}

	return interviews, allQALists, nil
}

// Update updates an interview and its question answers
func (r *InterviewRepo) Update(interview *models.AnalyzeInterview, qaList []models.QuestionAnswer) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()

	// Update interview
	query := `
	UPDATE analyze_interviews 
	SET updated_at = ? 
	WHERE id = ?
	`
	_, err = tx.Exec(query, now, interview.ID)
	if err != nil {
		return fmt.Errorf("failed to update interview: %w", err)
	}

	// Delete existing question answers
	deleteQuery := `DELETE FROM question_answers WHERE interview_id = ?`
	_, err = tx.Exec(deleteQuery, interview.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing question answers: %w", err)
	}

	// Insert new question answers
	for i := range qaList {
		qa := &qaList[i]
		qa.InterviewID = interview.ID
		qa.CreatedAt = now
		qa.UpdatedAt = now

		qaQuery := `
		INSERT INTO question_answers (interview_id, question, full_answer, accuracy, reason_unanswered, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		_, err = tx.Exec(qaQuery, qa.InterviewID, qa.Question, qa.FullAnswer, qa.Accuracy, qa.ReasonUnanswered, qa.CreatedAt, qa.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert question answer: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	interview.UpdatedAt = now
	return nil
}

// Delete deletes an interview and its question answers
func (r *InterviewRepo) Delete(id int64) error {
	query := `DELETE FROM analyze_interviews WHERE id = ?`
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete interview: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no interview found with id: %d", id)
	}

	return nil
}

// Helper method to get question answers by interview ID
func (r *InterviewRepo) getQuestionAnswersByInterviewID(interviewID int64) ([]models.QuestionAnswer, error) {
	query := `
	SELECT id, interview_id, question, full_answer, accuracy, reason_unanswered, created_at, updated_at
	FROM question_answers 
	WHERE interview_id = ?
	ORDER BY id
	`
	rows, err := db.Query(query, interviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to query question answers: %w", err)
	}
	defer rows.Close()

	var qaList []models.QuestionAnswer
	for rows.Next() {
		var qa models.QuestionAnswer
		err := rows.Scan(&qa.ID, &qa.InterviewID, &qa.Question, &qa.FullAnswer, &qa.Accuracy, &qa.ReasonUnanswered, &qa.CreatedAt, &qa.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question answer row: %w", err)
		}
		qaList = append(qaList, qa)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate question answer rows: %w", err)
	}

	return qaList, nil
}
