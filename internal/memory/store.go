package memory

import "errors"

var (
	ErrInvalidEpisode   = errors.New("invalid episode")
	ErrInvalidLimit     = errors.New("limit must be greater than zero")
	ErrEpisodeNotFound  = errors.New("episode not found")
	ErrProjectMismatch  = errors.New("episode belongs to another project")
	ErrDuplicateEpisode = errors.New("episode already exists")
)

type Store interface {
	SaveEpisode(Episode) error
	GetEpisode(projectID, id string) (Episode, error)
	ListRecentEpisodes(projectID string, limit int) ([]Episode, error)
	ListEpisodesByTargetKind(projectID, kind string) ([]Episode, error)
}
