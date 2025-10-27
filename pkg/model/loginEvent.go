package model

import "time"

type LoginEvent struct {
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
}
