package lrclib

import (
	"net/url"

	"codeberg.org/dankstuff/danklyrics/pkg/errors"
)

type SearchService struct {
	gClient *apiClient[[]Song]
}

// SearchParams represents parameters for searching music tracks.
// It allows for general keyword search or specific searches within
// track title, artist name, or album name.
type SearchParams struct {
	// Query searches for keywords present in any of the track's title,
	// artist name, or album name fields.
	Query string
	// TrackName searches for keywords specifically within the track's title.
	TrackName string
	// ArtistName searches for keywords specifically within the track's artist name.
	ArtistName string
	// AlbumName searches for keywords specifically within the track's album name.
	AlbumName string
	// Limit limits the search results count.
	Limit int
}

func (p SearchParams) Validate() error {
	if p.Query == "" && p.TrackName == "" {
		return &errors.ErrInvalidParams{
			ParamName: "q or track_name",
		}
	}

	if p.Query != "" && p.TrackName != "" {
		return &errors.ErrInvalidParams{
			ParamName: "q has a higher priority than track_name",
		}
	}

	return nil
}

func (p SearchParams) UrlParams() url.Values {
	q := url.Values{}

	if p.Query != "" {
		q.Set("q", p.Query)
	}
	if p.TrackName != "" {
		q.Set("track_name", p.TrackName)
	}
	if p.ArtistName != "" {
		q.Set("artist_name", p.ArtistName)
	}
	if p.AlbumName != "" {
		q.Set("artist_name", p.ArtistName)
	}

	return q
}

func (s *SearchService) Get(params SearchParams) ([]Song, error) {
	err := params.Validate()
	if err != nil {
		return nil, err
	}

	urlParams := params.UrlParams()
	for key, value := range urlParams {
		if len(value) == 0 {
			continue
		}
		s.gClient.setQueryParam(key, value[0])
	}

	res, err := s.gClient.callEndpoint()
	if err != nil {
		return nil, err
	}

	if params.Limit > 0 && len(res) > params.Limit {
		res = res[:params.Limit]
	}

	return res, nil
}
