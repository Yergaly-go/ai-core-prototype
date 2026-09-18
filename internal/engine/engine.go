package engine

import (
	"strconv"
	"sync"
	"time"

	"example.com/aicore/internal/core"
)

// Engine applies deterministic state transitions. It does not call models or
// external services.
type Engine struct {
	store Store
	mu    sync.Mutex
	next  uint64
}

func New(store Store) *Engine {
	return &Engine{store: store}
}

type NextLoopSpec struct {
	Hypothesis string
	Test       string
}

type ProjectState struct {
	Project   core.Project
	Loops     []core.Loop
	Evidence  []core.Evidence
	Decisions []core.Decision
}

func (e *Engine) CreateProject(name string) (core.Project, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	project := core.Project{
		ID:        e.newID("project"),
		Name:      name,
		Status:    core.ProjectActive,
		CreatedAt: time.Now().UTC(),
	}
	if err := e.store.CreateProject(project); err != nil {
		return core.Project{}, err
	}
	return project, nil
}

// StartLoop creates the first loop of a project. Successor loops are created
// only by Decide, so their chain relationship cannot be supplied incorrectly.
func (e *Engine) StartLoop(projectID, hypothesis, test string) (core.Loop, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, err := e.store.GetProject(projectID); err != nil {
		return core.Loop{}, err
	}
	loops, err := e.store.ListLoops(projectID)
	if err != nil {
		return core.Loop{}, err
	}
	if len(loops) != 0 {
		return core.Loop{}, ErrInitialLoopExists
	}
	return e.createLoop(projectID, "", hypothesis, test)
}

func (e *Engine) AddEvidence(projectID, loopID, content, source string) (core.Evidence, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, err := e.store.GetProject(projectID); err != nil {
		return core.Evidence{}, err
	}
	loop, err := e.store.GetLoop(loopID)
	if err != nil {
		return core.Evidence{}, err
	}
	if loop.ProjectID != projectID {
		return core.Evidence{}, ErrLoopProjectMismatch
	}
	if loop.Status != core.LoopOpen {
		return core.Evidence{}, ErrLoopClosed
	}
	evidence := core.Evidence{ID: e.newID("evidence"), ProjectID: projectID, LoopID: loopID, Content: content, Source: source, CreatedAt: time.Now().UTC()}
	if err := e.store.AppendEvidence(evidence); err != nil {
		return core.Evidence{}, err
	}
	return evidence, nil
}

func (e *Engine) Decide(projectID, loopID, result, reason, nextAction string, next *NextLoopSpec) (core.Decision, *core.Loop, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, err := e.store.GetProject(projectID); err != nil {
		return core.Decision{}, nil, err
	}
	loop, err := e.store.GetLoop(loopID)
	if err != nil {
		return core.Decision{}, nil, err
	}
	if loop.ProjectID != projectID {
		return core.Decision{}, nil, ErrLoopProjectMismatch
	}
	if loop.Status != core.LoopOpen {
		return core.Decision{}, nil, ErrLoopClosed
	}

	decision := core.Decision{ID: e.newID("decision"), ProjectID: projectID, LoopID: loopID, Result: result, Reason: reason, NextAction: nextAction, CreatedAt: time.Now().UTC()}
	if err := e.store.CreateDecision(decision); err != nil {
		return core.Decision{}, nil, err
	}
	closedAt := time.Now().UTC()
	loop.Status = core.LoopClosed
	loop.ClosedAt = &closedAt
	if err := e.store.UpdateLoop(loop); err != nil {
		return core.Decision{}, nil, err
	}
	if next == nil {
		return decision, nil, nil
	}
	nextLoop, err := e.createLoop(projectID, loop.ID, next.Hypothesis, next.Test)
	if err != nil {
		return core.Decision{}, nil, err
	}
	return decision, &nextLoop, nil
}

func (e *Engine) GetProjectState(projectID string) (ProjectState, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	project, err := e.store.GetProject(projectID)
	if err != nil {
		return ProjectState{}, err
	}
	loops, err := e.store.ListLoops(projectID)
	if err != nil {
		return ProjectState{}, err
	}
	evidence, err := e.store.ListEvidence(projectID)
	if err != nil {
		return ProjectState{}, err
	}
	decisions, err := e.store.ListDecisions(projectID)
	if err != nil {
		return ProjectState{}, err
	}
	return ProjectState{Project: project, Loops: loops, Evidence: evidence, Decisions: decisions}, nil
}

func (e *Engine) createLoop(projectID, previousLoopID, hypothesis, test string) (core.Loop, error) {
	loop := core.Loop{ID: e.newID("loop"), ProjectID: projectID, PreviousLoopID: previousLoopID, Hypothesis: hypothesis, Test: test, Status: core.LoopOpen, CreatedAt: time.Now().UTC()}
	if err := e.store.CreateLoop(loop); err != nil {
		return core.Loop{}, err
	}
	return loop, nil
}

func (e *Engine) newID(kind string) string {
	e.next++
	return kind + "-" + strconv.FormatUint(e.next, 10)
}
