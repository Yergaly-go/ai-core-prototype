package symbolic

import "time"

// Relation is a directed assertion between two entities.
type Relation struct {
	ID         string
	ProjectID  string
	SubjectID  string
	Predicate  string
	ObjectID   string
	Source     SourceRef
	Confidence *float64
	CreatedAt  time.Time
}
