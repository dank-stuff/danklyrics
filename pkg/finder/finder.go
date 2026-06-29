package finder

import (
	"codeberg.org/dankstuff/danklyrics/pkg/models"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"
	"errors"
)

// Service finds lyrics for a song using the enabled providers.
type Service struct {
	providers []provider.Service
}

// New creates a new [Service] with the selected configs.
func New(providers []provider.Service) (*Service, error) {
	if len(providers) == 0 {
		return nil, errors.New("must specify at least one lyrics provider")
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
	lyrics, err := l.providers[0].GetSongLyrics(params)
	if err != nil {
		if len(l.providers) <= 1 {
			return models.Lyrics{}, err
		}

		for _, providerFn := range l.providers[1:] {
			lyrics, err = providerFn.GetSongLyrics(params)
			if err == nil {
				break
			}
		}
	}

	return lyrics, nil
}
