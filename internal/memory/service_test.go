package memory_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	memorycore "example.com/aicore/internal/memory"
	memorystore "example.com/aicore/internal/store/memory"
)

func TestRecordEpisode(t *testing.T) {
	service := memorycore.NewService(memorystore.NewEpisodeStore())
	occurredAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	evidenceRef := memorycore.RecordRef{Kind: "evidence", ID: "evidence-7"}
	evidence, err := service.RecordEpisode("project-a", evidenceRef, occurredAt)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(evidence.Target, evidenceRef) {
		t.Fatalf("target changed: %#v", evidence.Target)
	}
	if evidence.CreatedAt.IsZero() {
		t.Fatal("created at was not generated")
	}
	decision, err := service.RecordEpisode("project-a", memorycore.RecordRef{Kind: "decision", ID: "decision-3"}, occurredAt.Add(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if decision.Target.Kind != "decision" || decision.Target.ID != "decision-3" {
		t.Fatalf("decision target = %#v", decision.Target)
	}
	if episodes, err := service.ListRecentEpisodes("project-a", 2); err != nil || len(episodes) != 2 {
		t.Fatalf("episodes = %#v, error = %v", episodes, err)
	}
}

func TestEpisodeValidationAndProjectIsolation(t *testing.T) {
	service := memorycore.NewService(memorystore.NewEpisodeStore())
	now := time.Now().UTC()
	invalid := []struct {
		project string
		target  memorycore.RecordRef
		at      time.Time
	}{
		{"", memorycore.RecordRef{Kind: "evidence", ID: "evidence-1"}, now},
		{"project-a", memorycore.RecordRef{ID: "evidence-1"}, now},
		{"project-a", memorycore.RecordRef{Kind: "evidence"}, now},
		{"project-a", memorycore.RecordRef{Kind: "evidence", ID: "evidence-1"}, time.Time{}},
	}
	for _, test := range invalid {
		if _, err := service.RecordEpisode(test.project, test.target, test.at); !errors.Is(err, memorycore.ErrInvalidEpisode) {
			t.Fatalf("invalid episode error = %v", err)
		}
	}
	episode, err := service.RecordEpisode("project-a", memorycore.RecordRef{Kind: "evidence", ID: "evidence-1"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.GetEpisode("project-b", episode.ID); !errors.Is(err, memorycore.ErrProjectMismatch) {
		t.Fatalf("foreign get error = %v", err)
	}
	if _, err := service.GetEpisode("", episode.ID); !errors.Is(err, memorycore.ErrInvalidEpisode) {
		t.Fatalf("empty project get error = %v", err)
	}
	if _, err := service.ListRecentEpisodes("project-a", 0); !errors.Is(err, memorycore.ErrInvalidLimit) {
		t.Fatalf("zero limit error = %v", err)
	}
	if _, err := service.ListRecentEpisodes("project-a", -1); !errors.Is(err, memorycore.ErrInvalidLimit) {
		t.Fatalf("negative limit error = %v", err)
	}
	if episodes, err := service.ListRecentEpisodes("project-b", 1); err != nil || len(episodes) != 0 {
		t.Fatalf("foreign project episodes = %#v, error = %v", episodes, err)
	}
}

func TestEpisodeOrderingAndKindFiltering(t *testing.T) {
	store := memorystore.NewEpisodeStore()
	service := memorycore.NewService(store)
	occurredAt := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, 1, 2, 1, 0, 0, 0, time.UTC)
	entries := []memorycore.Episode{
		{ID: "episode-c", ProjectID: "project-a", Target: memorycore.RecordRef{Kind: "evidence", ID: "evidence-c"}, OccurredAt: occurredAt, CreatedAt: createdAt},
		{ID: "episode-a", ProjectID: "project-a", Target: memorycore.RecordRef{Kind: "decision", ID: "decision-a"}, OccurredAt: occurredAt, CreatedAt: createdAt},
		{ID: "episode-b", ProjectID: "project-a", Target: memorycore.RecordRef{Kind: "evidence", ID: "evidence-b"}, OccurredAt: occurredAt, CreatedAt: createdAt.Add(time.Second)},
		{ID: "episode-old", ProjectID: "project-a", Target: memorycore.RecordRef{Kind: "evidence", ID: "evidence-old"}, OccurredAt: occurredAt.Add(-time.Second), CreatedAt: createdAt.Add(2 * time.Second)},
		{ID: "episode-other", ProjectID: "project-b", Target: memorycore.RecordRef{Kind: "evidence", ID: "evidence-other"}, OccurredAt: occurredAt.Add(time.Hour), CreatedAt: createdAt},
	}
	for _, entry := range entries {
		if err := store.SaveEpisode(entry); err != nil {
			t.Fatal(err)
		}
	}
	recent, err := service.ListRecentEpisodes("project-a", 10)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"episode-b", "episode-a", "episode-c", "episode-old"}
	if len(recent) != len(want) {
		t.Fatalf("recent length = %d, want %d", len(recent), len(want))
	}
	for index, id := range want {
		if recent[index].ID != id {
			t.Fatalf("recent[%d] = %q, want %q", index, recent[index].ID, id)
		}
	}
	evidence, err := service.ListEpisodesByTargetKind("project-a", "evidence")
	if err != nil {
		t.Fatal(err)
	}
	if len(evidence) != 3 {
		t.Fatalf("evidence episodes = %#v", evidence)
	}
	for _, episode := range evidence {
		if episode.ProjectID != "project-a" || episode.Target.Kind != "evidence" {
			t.Fatalf("kind filter leaked episode: %#v", episode)
		}
	}
}
