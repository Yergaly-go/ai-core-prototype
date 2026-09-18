package engine_test

import (
	"errors"
	"testing"

	"example.com/aicore/internal/core"
	"example.com/aicore/internal/engine"
	"example.com/aicore/internal/store/memory"
)

func TestCompleteLoopChain(t *testing.T) {
	e := engine.New(memory.New())
	project, err := e.CreateProject("MVP validation")
	if err != nil {
		t.Fatal(err)
	}
	loop1, err := e.StartLoop(project.ID, "Users need this", "Interview five users")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.AddEvidence(project.ID, loop1.ID, "4 of 5 confirmed", "interviews"); err != nil {
		t.Fatal(err)
	}
	decision, loop2, err := e.Decide(project.ID, loop1.ID, "continue", "signal is sufficient", "run pricing test", &engine.NextLoopSpec{Hypothesis: "Users will pay", Test: "Offer pre-order"})
	if err != nil {
		t.Fatal(err)
	}
	if decision.LoopID != loop1.ID {
		t.Fatalf("decision belongs to %q, want %q", decision.LoopID, loop1.ID)
	}
	if loop2 == nil {
		t.Fatal("next loop was not created")
	}
	if loop2.PreviousLoopID != loop1.ID {
		t.Fatalf("previous loop ID = %q, want %q", loop2.PreviousLoopID, loop1.ID)
	}
	state, err := e.GetProjectState(project.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Loops) != 2 || len(state.Evidence) != 1 || len(state.Decisions) != 1 {
		t.Fatalf("unexpected state sizes: loops=%d evidence=%d decisions=%d", len(state.Loops), len(state.Evidence), len(state.Decisions))
	}
	for _, loop := range state.Loops {
		if loop.ID == loop1.ID && (loop.Status != core.LoopClosed || loop.ClosedAt == nil) {
			t.Fatalf("loop 1 was not closed")
		}
	}
}

func TestInvalidTransitions(t *testing.T) {
	e := engine.New(memory.New())
	project1, err := e.CreateProject("one")
	if err != nil {
		t.Fatal(err)
	}
	project2, err := e.CreateProject("two")
	if err != nil {
		t.Fatal(err)
	}
	loop, err := e.StartLoop(project1.ID, "h", "t")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.AddEvidence(project1.ID, "unknown", "x", "test"); !errors.Is(err, engine.ErrLoopNotFound) {
		t.Fatalf("unknown loop error = %v", err)
	}
	if _, err := e.AddEvidence(project2.ID, loop.ID, "x", "test"); !errors.Is(err, engine.ErrLoopProjectMismatch) {
		t.Fatalf("foreign loop error = %v", err)
	}
	if _, _, err := e.Decide(project1.ID, loop.ID, "stop", "enough", "", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := e.AddEvidence(project1.ID, loop.ID, "late", "test"); !errors.Is(err, engine.ErrLoopClosed) {
		t.Fatalf("closed loop evidence error = %v", err)
	}
	if _, _, err := e.Decide(project1.ID, loop.ID, "again", "", "", nil); !errors.Is(err, engine.ErrLoopClosed) {
		t.Fatalf("second decision error = %v", err)
	}
}
