package lrclib

import (
	"net/url"
	"strconv"

	"codeberg.org/dankstuff/danklyrics/pkg/errors"
	"codeberg.org/dankstuff/danklyrics/pkg/models"
)

// Song represents how a song looks like coming from lrclib.net/api
type Song struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	TrackName    string  `json:"trackName"`
	ArtistName   string  `json:"artistName"`
	AlbumName    string  `json:"albumName"`
	Duration     float32 `json:"duration"`
	Instrumental bool    `json:"instrumental"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`

	lyrics *models.Lyrics
}

type SongsService struct {
	gClient *apiClient[Song]
}

func (s *SongsService) Get(id string) (Song, error) {
	s.gClient.appendToPath("get/" + id)

	res, err := s.gClient.callEndpoint()
	if err != nil {
		return Song{}, err
	}

	return res, nil
}

type GetSongParams struct {
	// TrackName is the title of the track.
	TrackName string
	// ArtistName is the name of the artist who performed the track.
	ArtistName string
	// AlbumName is the name of the album to which the track belongs.
	AlbumName string
	// Duration is the track's duration in seconds.
	Duration int `json:"duration"`
}

func (p GetSongParams) Validate() error {
	if p.TrackName == "" {
		return &errors.ErrInvalidParams{
			ParamName: "track_name",
		}
	}
	if p.ArtistName == "" {
		return &errors.ErrInvalidParams{
			ParamName: "artist_name",
		}
	}

	return nil
}

func (p GetSongParams) UrlParams() url.Values {
	q := url.Values{}

	if p.TrackName != "" {
		q.Set("track_name", p.TrackName)
	}
	if p.ArtistName != "" {
		q.Set("artist_name", p.ArtistName)
	}
	if p.AlbumName != "" {
		q.Set("artist_name", p.ArtistName)
	}
	if p.Duration != 0 {
		q.Set("duration", strconv.Itoa(p.Duration))
	}

	return q
}

func (s *SongsService) GetByParams(params GetSongParams) (Song, error) {
	err := params.Validate()
	if err != nil {
		return Song{}, err
	}

	s.gClient.appendToPath("get")

	urlParams := params.UrlParams()
	for key, value := range urlParams {
		if len(value) == 0 {
			continue
		}
		s.gClient.setQueryParam(key, value[0])
	}

	res, err := s.gClient.callEndpoint()
	if err != nil {
		return Song{}, err
	}

	return res, nil
}
