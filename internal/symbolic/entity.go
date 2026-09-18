package symbolic

import "time"

// Entity is a named object of knowledge.
type Entity struct {
	ID        string
	ProjectID string
	Label     string
	Source    SourceRef
	CreatedAt time.Time
}
