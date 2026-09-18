package symbolic

import (
	"strconv"
	"sync"
	"time"
)

// Service validates and appends symbolic records. It has no reasoning or
// lifecycle responsibilities.
type Service struct {
	store Store
	mu    sync.Mutex
	next  uint64
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) RegisterEntity(projectID, label string, source SourceRef) (Entity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if projectID == "" || label == "" || !source.valid() {
		return Entity{}, ErrInvalidRecord
	}
	entity := Entity{ID: s.newID("entity"), ProjectID: projectID, Label: label, Source: source, CreatedAt: time.Now().UTC()}
	if err := s.store.SaveEntity(entity); err != nil {
		return Entity{}, err
	}
	return entity, nil
}

func (s *Service) AddFact(projectID, subjectID, predicate, value string, source SourceRef, confidence *float64) (Fact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if projectID == "" || subjectID == "" || predicate == "" || value == "" || !source.valid() {
		return Fact{}, ErrInvalidRecord
	}
	if err := validateConfidence(confidence); err != nil {
		return Fact{}, err
	}
	if err := s.requireEntity(projectID, subjectID); err != nil {
		return Fact{}, err
	}
	fact := Fact{ID: s.newID("fact"), ProjectID: projectID, SubjectID: subjectID, Predicate: predicate, Value: value, Source: source, Confidence: confidence, CreatedAt: time.Now().UTC()}
	if err := s.store.SaveFact(fact); err != nil {
		return Fact{}, err
	}
	return fact, nil
}

func (s *Service) AddRelation(projectID, subjectID, predicate, objectID string, source SourceRef, confidence *float64) (Relation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if projectID == "" || subjectID == "" || predicate == "" || objectID == "" || !source.valid() {
		return Relation{}, ErrInvalidRecord
	}
	if err := validateConfidence(confidence); err != nil {
		return Relation{}, err
	}
	if err := s.requireEntity(projectID, subjectID); err != nil {
		return Relation{}, err
	}
	if err := s.requireEntity(projectID, objectID); err != nil {
		return Relation{}, err
	}
	relation := Relation{ID: s.newID("relation"), ProjectID: projectID, SubjectID: subjectID, Predicate: predicate, ObjectID: objectID, Source: source, Confidence: confidence, CreatedAt: time.Now().UTC()}
	if err := s.store.SaveRelation(relation); err != nil {
		return Relation{}, err
	}
	return relation, nil
}

func (s *Service) AddConstraint(projectID, subjectID string, operator ConstraintOperator, predicate, value string, source SourceRef, confidence *float64) (Constraint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if projectID == "" || subjectID == "" || predicate == "" || value == "" || !source.valid() || (operator != Must && operator != MustNot) {
		return Constraint{}, ErrInvalidRecord
	}
	if err := validateConfidence(confidence); err != nil {
		return Constraint{}, err
	}
	if err := s.requireEntity(projectID, subjectID); err != nil {
		return Constraint{}, err
	}
	constraint := Constraint{ID: s.newID("constraint"), ProjectID: projectID, SubjectID: subjectID, Operator: operator, Predicate: predicate, Value: value, Source: source, Confidence: confidence, CreatedAt: time.Now().UTC()}
	if err := s.store.SaveConstraint(constraint); err != nil {
		return Constraint{}, err
	}
	return constraint, nil
}

func (s *Service) ListEntities(projectID string) ([]Entity, error) {
	return s.store.ListEntities(projectID)
}
func (s *Service) ListFacts(projectID string) ([]Fact, error) { return s.store.ListFacts(projectID) }
func (s *Service) ListRelations(projectID string) ([]Relation, error) {
	return s.store.ListRelations(projectID)
}
func (s *Service) ListConstraints(projectID string) ([]Constraint, error) {
	return s.store.ListConstraints(projectID)
}

func (s *Service) requireEntity(projectID, id string) error {
	_, err := s.store.GetEntity(projectID, id)
	return err
}

func (s *Service) newID(kind string) string {
	s.next++
	return kind + "-" + strconv.FormatUint(s.next, 10)
}

func validateConfidence(confidence *float64) error {
	if confidence != nil && (*confidence < 0 || *confidence > 1) {
		return ErrConfidenceOutOfRange
	}
	return nil
}
