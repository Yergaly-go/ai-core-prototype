package core

import "time"

type LoopStatus string

const (
	LoopOpen   LoopStatus = "open"
	LoopClosed LoopStatus = "closed"
)

type Loop struct {
	ID             string
	ProjectID      string
	PreviousLoopID string
	Hypothesis     string
	Test           string
	Status         LoopStatus
	CreatedAt      time.Time
	ClosedAt       *time.Time
}
