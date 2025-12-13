package models

import (
	"time"
)

type (
	AnalyzeInterview struct {
		ID        uint64    `json:"id" gorm:"primaryKey" db:"id"`
		CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" db:"updated_at"`
	}
	AnalyzeInterviewWithQA struct {
		ID        uint64           `json:"id" gorm:"primaryKey" db:"id"`
		QA        []QuestionAnswer `json:"qa" gorm:"foreignKey:InterviewID" db:"qa"`
		CreatedAt time.Time        `json:"created_at" gorm:"autoCreateTime" db:"created_at"`
		UpdatedAt time.Time        `json:"updated_at" gorm:"autoUpdateTime" db:"updated_at"`
	}
	QuestionAnswer struct {
		ID               uint64    `json:"id" gorm:"primaryKey" db:"id"`
		InterviewID      uint64    `json:"interview_id" gorm:"index;not null" db:"interview_id"`
		Question         string    `json:"question" gorm:"not null" db:"question"`
		FullAnswer       string    `json:"full_answer" db:"full_answer"`
		Accuracy         float64   `json:"accuracy" gorm:"not null" db:"accuracy"`
		ReasonUnanswered string    `json:"reason_unanswered" db:"reason_unanswered"`
		CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime" db:"created_at"`
		UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime" db:"updated_at"`
	}

	// GetInterviewsFilters represents filters for querying analytics
	GetInterviewsFilters struct {
		DateFrom *time.Time `json:"dateFrom,omitempty"`
		DateTo   *time.Time `json:"dateTo,omitempty"`
	}
)
