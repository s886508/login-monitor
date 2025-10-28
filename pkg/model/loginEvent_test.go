package model

import (
	"time"
)

func TestLoginEventIsValid() {
	// Case 1: Invalid user id
	e := &LoginEvent{Timestamp: time.Now(), Success: true}
	asssert.False(t, e.IsValid())

	// Case 2: Invalid timestamp
	e := &LoginEvent{UserID: "A", Success: true}
	asssert.False(t, e.IsValid())
}
