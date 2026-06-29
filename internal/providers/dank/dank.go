package dank

import (
	"codeberg.org/dankstuff/danklyrics/pkg/models"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type dankProvider struct {
}

func New() provider.Service {
	return &dankProvider{}
}

func (d *dankProvider) GetSongLyrics(s provider.SearchParams) (models.Lyrics, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"https://api.danklyrics.com/dank/lyrics?song=%s&artist=%s&album=%s",
			url.QueryEscape(s.SongName), url.QueryEscape(s.ArtistName), url.QueryEscape(s.AlbumName),
		),
		http.NoBody)
	if err != nil {
		return models.Lyrics{}, err
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return models.Lyrics{}, err
	}

	var lyrics []models.Lyrics
	err = json.NewDecoder(resp.Body).Decode(&lyrics)
	if err != nil {
		return models.Lyrics{}, err
	}
	_ = resp.Body.Close()

	if len(lyrics) == 0 {
		return models.Lyrics{}, errors.New("no results were found")
	}

	return lyrics[0], nil
}
