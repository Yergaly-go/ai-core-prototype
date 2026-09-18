// Package symbolic provides canonical, project-scoped symbolic records.
package symbolic

import "errors"

var (
	ErrInvalidRecord        = errors.New("invalid symbolic record")
	ErrEntityNotFound       = errors.New("entity not found")
	ErrProjectMismatch      = errors.New("symbolic record belongs to another project")
	ErrConfidenceOutOfRange = errors.New("confidence must be between zero and one")
	ErrDuplicateRecord      = errors.New("symbolic record already exists")
)

// SourceRef identifies the origin of a symbolic record.
type SourceRef struct {
	Kind string
	ID   string
}

func (s SourceRef) valid() bool {
	return s.Kind != "" && s.ID != ""
}
