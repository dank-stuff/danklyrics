package lyricfind

import (
	"codeberg.org/dankstuff/danklyrics/pkg/errors"
	"codeberg.org/dankstuff/danklyrics/pkg/models"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"

	"codeberg.org/dankstuff/danklyrics/pkg/lrclib"
)

type lyricFindProvider struct {
	client *lrclib.Client
}

func New() provider.Service {
	return &lyricFindProvider{
		client: lrclib.NewClient(),
	}
}

func (l *lyricFindProvider) GetSongLyrics(s provider.SearchParams) (models.Lyrics, error) {
	var lrcSearch lrclib.SearchParams
	if s.Query != "" {
		lrcSearch.Query = s.Query
	} else {
		lrcSearch = lrclib.SearchParams{
			TrackName:  s.SongName,
			ArtistName: s.ArtistName,
			AlbumName:  s.AlbumName,
			Limit:      0,
		}
	}

	hits, err := l.client.Search.Get(lrcSearch)
	if err != nil {
		return models.Lyrics{}, err
	}

	if len(hits) == 0 {
		return models.Lyrics{}, &errors.ErrNotFound{}
	}

	return *hits[0].Lyrics(), nil
}
