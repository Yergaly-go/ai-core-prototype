package core

import "time"

type Evidence struct {
	ID        string
	ProjectID string
	LoopID    string
	Content   string
	Source    string
	CreatedAt time.Time
}
