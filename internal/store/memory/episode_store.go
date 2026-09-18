package memory

import (
	"sort"
	"sync"

	memorycore "example.com/aicore/internal/memory"
)

type EpisodeStore struct {
	mu       sync.RWMutex
	episodes map[string]memorycore.Episode
}

func NewEpisodeStore() *EpisodeStore {
	return &EpisodeStore{episodes: make(map[string]memorycore.Episode)}
}

func (s *EpisodeStore) SaveEpisode(episode memorycore.Episode) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.episodes[episode.ID]; ok {
		return memorycore.ErrDuplicateEpisode
	}
	s.episodes[episode.ID] = episode
	return nil
}

func (s *EpisodeStore) GetEpisode(projectID, id string) (memorycore.Episode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	episode, ok := s.episodes[id]
	if !ok {
		return memorycore.Episode{}, memorycore.ErrEpisodeNotFound
	}
	if episode.ProjectID != projectID {
		return memorycore.Episode{}, memorycore.ErrProjectMismatch
	}
	return episode, nil
}

func (s *EpisodeStore) ListRecentEpisodes(projectID string, limit int) ([]memorycore.Episode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := s.list(projectID, "")
	if limit < len(result) {
		result = result[:limit]
	}
	return result, nil
}

func (s *EpisodeStore) ListEpisodesByTargetKind(projectID, kind string) ([]memorycore.Episode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.list(projectID, kind), nil
}

func (s *EpisodeStore) list(projectID, kind string) []memorycore.Episode {
	result := make([]memorycore.Episode, 0)
	for _, episode := range s.episodes {
		if episode.ProjectID == projectID && (kind == "" || episode.Target.Kind == kind) {
			result = append(result, episode)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if !result[i].OccurredAt.Equal(result[j].OccurredAt) {
			return result[i].OccurredAt.After(result[j].OccurredAt)
		}
		if !result[i].CreatedAt.Equal(result[j].CreatedAt) {
			return result[i].CreatedAt.After(result[j].CreatedAt)
		}
		return result[i].ID < result[j].ID
	})
	return result
}

var _ memorycore.Store = (*EpisodeStore)(nil)
