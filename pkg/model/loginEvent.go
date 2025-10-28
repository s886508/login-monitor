package model

import "time"

type LoginEvent struct {
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
}

func (e *LoginEvent) IsValid() bool {
	if len(e.UserID) == 0 {
		return false
	}

	if e.Timestamp.IsZero() {
		return false
	}
	return true
}
