package finder

import (
	"codeberg.org/dankstuff/danklyrics/pkg/errors"
	"codeberg.org/dankstuff/danklyrics/pkg/models"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"
)

// Service finds lyrics for a song using the enabled providers.
type Service struct {
	providers []provider.Service
}

// New creates a new [Service] with the selected configs.
func New(providers []provider.Service) (*Service, error) {
	if len(providers) == 0 {
		return nil, &errors.ErrMissingProvider{}
	}

	return &Service{
		providers: providers,
	}, nil
}

// GetSongLyrics search for song's lyrics using the enabled providers list,
// where using a provider depends on the provider's order in that list.
//
// returns [Lyrics] and an occurring [error]
func (l *Service) GetSongLyrics(params provider.SearchParams) (models.Lyrics, error) {
	for _, provider := range l.providers {
		lyrics, err := provider.GetSongLyrics(params)
		if err == nil {
			return lyrics, nil
		}
	}

	return models.Lyrics{}, &errors.ErrNotFound{}
}
