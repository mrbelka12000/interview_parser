package models

import (
	"encoding/json"
	"time"
)

type (
	Call struct {
		ID         uint64         `json:"id" gorm:"primaryKey" db:"id"`
		Transcript string         `json:"transcript" gorm:"not null" db:"transcript"`
		Analysis   json.RawMessage `json:"analysis" db:"analysis"`
		CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime" db:"created_at"`
		UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime" db:"updated_at"`
	}
)
