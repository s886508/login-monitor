package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoginEventIsValid(t *testing.T) {
	// Case 1: Invalid user id
	e := &LoginEvent{Timestamp: time.Now(), Success: true}
	assert.False(t, e.IsValid())

	// Case 2: Invalid timestamp
	e = &LoginEvent{UserID: "A", Success: true}
	assert.False(t, e.IsValid())
}
