package memory

import "time"

type Episode struct {
	ID         string
	ProjectID  string
	Target     RecordRef
	OccurredAt time.Time
	CreatedAt  time.Time
}
