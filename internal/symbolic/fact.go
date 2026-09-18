package symbolic

import "time"

// Fact is an assertion from an entity to a literal value.
type Fact struct {
	ID         string
	ProjectID  string
	SubjectID  string
	Predicate  string
	Value      string
	Source     SourceRef
	Confidence *float64
	CreatedAt  time.Time
}
