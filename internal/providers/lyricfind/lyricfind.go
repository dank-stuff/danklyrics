package lyricfind

import (
	"errors"

	"codeberg.org/dankstuff/danklyrics/pkg/models"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"

	"github.com/mbaraa/lrclibgo"
)

type lyricFindProvider struct {
	client *lrclibgo.Client
}

func New() provider.Service {
	return &lyricFindProvider{
		client: lrclibgo.NewClient(),
	}
}

func (l *lyricFindProvider) GetSongLyrics(s provider.SearchParams) (models.Lyrics, error) {
	var lrcSearch lrclibgo.SearchParams
	if s.Query != "" {
		lrcSearch.Query = s.Query
	} else {
		lrcSearch = lrclibgo.SearchParams{
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
		return models.Lyrics{}, errors.New("no results were found")
	}

	lyrics := hits[0].Lyrics()
	return models.Lyrics{
		SongName:   hits[0].TrackName,
		ArtistName: hits[0].ArtistName,
		AlbumName:  hits[0].AlbumName,
		Parts:      lyrics.Parts(),
		Synced:     lyrics.Synced(),
	}, nil
}
