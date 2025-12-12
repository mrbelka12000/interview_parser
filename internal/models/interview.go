package models

import "time"

type (
	AnalyzeInterview struct {
		ID        int64     `json:"id" db:"id"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	}
	AnalyzeInterviewWithQA struct {
		ID        int64            `json:"id" db:"id"`
		QA        []QuestionAnswer `json:"qa" db:"qa"`
		CreatedAt time.Time        `json:"created_at" db:"created_at"`
		UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
	}
	QuestionAnswer struct {
		ID               int64     `json:"id" db:"id"`
		InterviewID      int64     `json:"interview_id" db:"interview_id"`
		Question         string    `json:"question" db:"question"`
		FullAnswer       string    `json:"full_answer" db:"full_answer"`
		Accuracy         float64   `json:"accuracy" db:"accuracy"`
		ReasonUnanswered string    `json:"reason_unanswered" db:"reason_unanswered"`
		CreatedAt        time.Time `json:"created_at" db:"created_at"`
		UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	}

	// GetInterviewsFilters represents filters for querying analytics
	GetInterviewsFilters struct {
		DateFrom *time.Time `json:"dateFrom,omitempty"`
		DateTo   *time.Time `json:"dateTo,omitempty"`
	}
)
