package memory

import (
	"sync"

	"example.com/aicore/internal/symbolic"
)

// SymbolicStore is an in-memory, concurrency-safe symbolic.Store adapter.
type SymbolicStore struct {
	mu          sync.RWMutex
	entities    map[string]symbolic.Entity
	facts       map[string]symbolic.Fact
	relations   map[string]symbolic.Relation
	constraints map[string]symbolic.Constraint
}

func NewSymbolicStore() *SymbolicStore {
	return &SymbolicStore{entities: make(map[string]symbolic.Entity), facts: make(map[string]symbolic.Fact), relations: make(map[string]symbolic.Relation), constraints: make(map[string]symbolic.Constraint)}
}

func (s *SymbolicStore) SaveEntity(entity symbolic.Entity) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entities[entity.ID]; ok {
		return symbolic.ErrDuplicateRecord
	}
	s.entities[entity.ID] = entity
	return nil
}
func (s *SymbolicStore) GetEntity(projectID, id string) (symbolic.Entity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entity, ok := s.entities[id]
	if !ok {
		return symbolic.Entity{}, symbolic.ErrEntityNotFound
	}
	if entity.ProjectID != projectID {
		return symbolic.Entity{}, symbolic.ErrProjectMismatch
	}
	return entity, nil
}
func (s *SymbolicStore) ListEntities(projectID string) ([]symbolic.Entity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]symbolic.Entity, 0)
	for _, item := range s.entities {
		if item.ProjectID == projectID {
			result = append(result, item)
		}
	}
	return result, nil
}
func (s *SymbolicStore) SaveFact(fact symbolic.Fact) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.facts[fact.ID]; ok {
		return symbolic.ErrDuplicateRecord
	}
	s.facts[fact.ID] = fact
	return nil
}
func (s *SymbolicStore) ListFacts(projectID string) ([]symbolic.Fact, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]symbolic.Fact, 0)
	for _, item := range s.facts {
		if item.ProjectID == projectID {
			result = append(result, item)
		}
	}
	return result, nil
}
func (s *SymbolicStore) SaveRelation(relation symbolic.Relation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.relations[relation.ID]; ok {
		return symbolic.ErrDuplicateRecord
	}
	s.relations[relation.ID] = relation
	return nil
}
func (s *SymbolicStore) ListRelations(projectID string) ([]symbolic.Relation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]symbolic.Relation, 0)
	for _, item := range s.relations {
		if item.ProjectID == projectID {
			result = append(result, item)
		}
	}
	return result, nil
}
func (s *SymbolicStore) SaveConstraint(constraint symbolic.Constraint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.constraints[constraint.ID]; ok {
		return symbolic.ErrDuplicateRecord
	}
	s.constraints[constraint.ID] = constraint
	return nil
}
func (s *SymbolicStore) ListConstraints(projectID string) ([]symbolic.Constraint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]symbolic.Constraint, 0)
	for _, item := range s.constraints {
		if item.ProjectID == projectID {
			result = append(result, item)
		}
	}
	return result, nil
}

var _ symbolic.Store = (*SymbolicStore)(nil)
