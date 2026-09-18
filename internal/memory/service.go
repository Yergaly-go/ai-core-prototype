package memory

import (
	"strconv"
	"sync"
	"time"
)

type Service struct {
	store Store
	mu    sync.Mutex
	next  uint64
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) RecordEpisode(projectID string, target RecordRef, occurredAt time.Time) (Episode, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if projectID == "" || !target.valid() || occurredAt.IsZero() {
		return Episode{}, ErrInvalidEpisode
	}
	episode := Episode{
		ID:         s.newID(),
		ProjectID:  projectID,
		Target:     target,
		OccurredAt: occurredAt,
		CreatedAt:  time.Now().UTC(),
	}
	if err := s.store.SaveEpisode(episode); err != nil {
		return Episode{}, err
	}
	return episode, nil
}

func (s *Service) GetEpisode(projectID, id string) (Episode, error) {
	if projectID == "" || id == "" {
		return Episode{}, ErrInvalidEpisode
	}
	return s.store.GetEpisode(projectID, id)
}

func (s *Service) ListRecentEpisodes(projectID string, limit int) ([]Episode, error) {
	if projectID == "" {
		return nil, ErrInvalidEpisode
	}
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}
	return s.store.ListRecentEpisodes(projectID, limit)
}

func (s *Service) ListEpisodesByTargetKind(projectID, kind string) ([]Episode, error) {
	if projectID == "" || kind == "" {
		return nil, ErrInvalidEpisode
	}
	return s.store.ListEpisodesByTargetKind(projectID, kind)
}

func (s *Service) newID() string {
	s.next++
	return "episode-" + strconv.FormatUint(s.next, 10)
}
