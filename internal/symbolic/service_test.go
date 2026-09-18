package symbolic_test

import (
	"errors"
	"reflect"
	"testing"

	"example.com/aicore/internal/store/memory"
	"example.com/aicore/internal/symbolic"
)

func TestSymbolicCreationAndProvenance(t *testing.T) {
	s := symbolic.NewService(memory.NewSymbolicStore())
	source := symbolic.SourceRef{Kind: "evidence", ID: "evidence-7"}
	project := "project-a"
	first, err := s.RegisterEntity(project, "Project", source)
	if err != nil {
		t.Fatal(err)
	}
	second, err := s.RegisterEntity(project, "Founder", source)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.AddFact(project, first.ID, "HAS_STAGE", "validation", source, nil); err != nil {
		t.Fatal(err)
	}
	zero, one := 0.0, 1.0
	relation, err := s.AddRelation(project, first.ID, "HAS_OWNER", second.ID, source, &zero)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(relation.Source, source) {
		t.Fatalf("source changed: %#v", relation.Source)
	}
	if _, err := s.AddConstraint(project, first.ID, symbolic.MustNot, "EXPOSE", "secret", source, &one); err != nil {
		t.Fatal(err)
	}
	if facts, _ := s.ListFacts(project); len(facts) != 1 || facts[0].Confidence != nil {
		t.Fatal("literal fact or nil confidence not retained")
	}
}

func TestSymbolicValidationAndProjectIsolation(t *testing.T) {
	s := symbolic.NewService(memory.NewSymbolicStore())
	source := symbolic.SourceRef{Kind: "user", ID: "user-1"}
	a, err := s.RegisterEntity("project-a", "A", source)
	if err != nil {
		t.Fatal(err)
	}
	b, err := s.RegisterEntity("project-b", "B", source)
	if err != nil {
		t.Fatal(err)
	}
	negative, above := -0.1, 1.1
	if _, err := s.AddFact("project-a", a.ID, "P", "v", source, &negative); !errors.Is(err, symbolic.ErrConfidenceOutOfRange) {
		t.Fatalf("negative confidence: %v", err)
	}
	if _, err := s.AddFact("project-a", a.ID, "P", "v", source, &above); !errors.Is(err, symbolic.ErrConfidenceOutOfRange) {
		t.Fatalf("above-one confidence: %v", err)
	}
	if _, err := s.AddFact("project-a", "unknown", "P", "v", source, nil); !errors.Is(err, symbolic.ErrEntityNotFound) {
		t.Fatalf("unknown subject: %v", err)
	}
	if _, err := s.AddFact("project-a", b.ID, "P", "v", source, nil); !errors.Is(err, symbolic.ErrProjectMismatch) {
		t.Fatalf("foreign subject: %v", err)
	}
	fact, err := s.AddFact("project-a", a.ID, "P", a.ID, source, nil)
	if err != nil {
		t.Fatalf("entity-like literal value: %v", err)
	}
	if fact.Value != a.ID {
		t.Fatalf("fact value = %q, want literal %q", fact.Value, a.ID)
	}
	if relations, _ := s.ListRelations("project-a"); len(relations) != 0 {
		t.Fatalf("literal fact unexpectedly created relations: %#v", relations)
	}
	if _, err := s.AddRelation("project-a", a.ID, "LINKS", b.ID, source, nil); !errors.Is(err, symbolic.ErrProjectMismatch) {
		t.Fatalf("foreign relation: %v", err)
	}
	if _, err := s.AddConstraint("project-a", b.ID, symbolic.Must, "DO", "x", source, nil); !errors.Is(err, symbolic.ErrProjectMismatch) {
		t.Fatalf("foreign constraint: %v", err)
	}
	if entities, _ := s.ListEntities("project-a"); len(entities) != 1 || entities[0].ID != a.ID {
		t.Fatalf("project-scoped list leaked: %#v", entities)
	}
}
