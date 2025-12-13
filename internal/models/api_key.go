package models

import (
	"time"
)

type (
	APIKey struct {
		ID        uint64    `json:"id" gorm:"primaryKey" db:"id"`
		APIKey    string    `json:"api_key" gorm:"not null" db:"api_key"`
		CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" db:"created_at"`
	}
)
