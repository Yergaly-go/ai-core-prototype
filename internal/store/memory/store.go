// Package memory provides the in-memory Store implementation for Loop 1.
package memory

import (
	"sync"

	"example.com/aicore/internal/core"
	"example.com/aicore/internal/engine"
)

type Store struct {
	mu        sync.RWMutex
	projects  map[string]core.Project
	loops     map[string]core.Loop
	evidence  map[string]core.Evidence
	decisions map[string]core.Decision
}

func New() *Store {
	return &Store{projects: make(map[string]core.Project), loops: make(map[string]core.Loop), evidence: make(map[string]core.Evidence), decisions: make(map[string]core.Decision)}
}

func (s *Store) CreateProject(project core.Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.projects[project.ID]; ok {
		return engine.ErrDuplicateRecord
	}
	s.projects[project.ID] = project
	return nil
}

func (s *Store) GetProject(id string) (core.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	project, ok := s.projects[id]
	if !ok {
		return core.Project{}, engine.ErrProjectNotFound
	}
	return project, nil
}

func (s *Store) CreateLoop(loop core.Loop) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.loops[loop.ID]; ok {
		return engine.ErrDuplicateRecord
	}
	s.loops[loop.ID] = loop
	return nil
}

func (s *Store) GetLoop(id string) (core.Loop, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	loop, ok := s.loops[id]
	if !ok {
		return core.Loop{}, engine.ErrLoopNotFound
	}
	return loop, nil
}

func (s *Store) UpdateLoop(loop core.Loop) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.loops[loop.ID]; !ok {
		return engine.ErrLoopNotFound
	}
	s.loops[loop.ID] = loop
	return nil
}

func (s *Store) AppendEvidence(evidence core.Evidence) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.evidence[evidence.ID]; ok {
		return engine.ErrDuplicateRecord
	}
	s.evidence[evidence.ID] = evidence
	return nil
}

func (s *Store) CreateDecision(decision core.Decision) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.decisions[decision.ID]; ok {
		return engine.ErrDuplicateRecord
	}
	s.decisions[decision.ID] = decision
	return nil
}

func (s *Store) ListLoops(projectID string) ([]core.Loop, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]core.Loop, 0)
	for _, item := range s.loops {
		if item.ProjectID == projectID {
			result = append(result, item)
		}
	}
	return result, nil
}
func (s *Store) ListEvidence(projectID string) ([]core.Evidence, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]core.Evidence, 0)
	for _, item := range s.evidence {
		if item.ProjectID == projectID {
			result = append(result, item)
		}
	}
	return result, nil
}
func (s *Store) ListDecisions(projectID string) ([]core.Decision, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]core.Decision, 0)
	for _, item := range s.decisions {
		if item.ProjectID == projectID {
			result = append(result, item)
		}
	}
	return result, nil
}

var _ engine.Store = (*Store)(nil)
