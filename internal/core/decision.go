package core

import "time"

type Decision struct {
	ID         string
	ProjectID  string
	LoopID     string
	Result     string
	Reason     string
	NextAction string
	CreatedAt  time.Time
}
