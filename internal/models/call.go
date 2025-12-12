package models

import (
	"encoding/json"
	"time"
)

type (
	Call struct {
		ID         int64
		Transcript string
		Analysis   json.RawMessage
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}
)
