// Package engine coordinates valid transitions of the project loop.
package engine

import (
	"errors"

	"example.com/aicore/internal/core"
)

var (
	ErrProjectNotFound     = errors.New("project not found")
	ErrLoopNotFound        = errors.New("loop not found")
	ErrLoopProjectMismatch = errors.New("loop does not belong to project")
	ErrLoopClosed          = errors.New("loop is closed")
	ErrInitialLoopExists   = errors.New("project already has a loop chain")
	ErrDuplicateRecord     = errors.New("record already exists")
)

// Store is the sole persistence boundary for Loop 1.
type Store interface {
	CreateProject(core.Project) error
	GetProject(id string) (core.Project, error)
	CreateLoop(core.Loop) error
	GetLoop(id string) (core.Loop, error)
	UpdateLoop(core.Loop) error
	AppendEvidence(core.Evidence) error
	CreateDecision(core.Decision) error
	ListLoops(projectID string) ([]core.Loop, error)
	ListEvidence(projectID string) ([]core.Evidence, error)
	ListDecisions(projectID string) ([]core.Decision, error)
}
