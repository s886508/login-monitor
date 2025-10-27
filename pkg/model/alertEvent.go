package model

type Alert struct {
	UserID      string       `json:"user_id"`
	FailedCount int          `json:"failed_count"`
	TimeWindow  string       `json:"time_window"`
	Events      []LoginEvent `json:"events"`
}
