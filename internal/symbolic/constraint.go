package symbolic

import "time"

type ConstraintOperator string

const (
	Must    ConstraintOperator = "MUST"
	MustNot ConstraintOperator = "MUST_NOT"
)

// Constraint is a normative condition. It has no execution semantics here.
type Constraint struct {
	ID         string
	ProjectID  string
	SubjectID  string
	Operator   ConstraintOperator
	Predicate  string
	Value      string
	Source     SourceRef
	Confidence *float64
	CreatedAt  time.Time
}
