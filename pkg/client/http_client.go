package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"codeberg.org/dankstuff/danklyrics/pkg/models"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"
)

// Http is the dank lyrics finding client that makes a call to api.danklyrics.com to find the lyrics.
type Http struct {
	providers  string
	apiAddress string
}

func NewHttp(c Config) (*Http, error) {
	if len(c.Providers) == 0 {
		return nil, errors.New("must specify at least one lyrics provider")
	}

	providersStr := ""
	for i, p := range c.Providers {
		providersStr += "providers=" + string(p)
		if i < len(c.Providers)-1 {
			providersStr += "&"
		}
	}

	if c.ApiAddress == "" {
		c.ApiAddress = "https://api.danklyrics.com"
	}

	return &Http{
		providers:  providersStr,
		apiAddress: c.ApiAddress,
	}, nil
}

// GetSongLyrics search for song's lyrics using the enabled providers list,
// where using a provider depends on the provider's order in that list.
//
// returns [Lyrics] and an occurring [error]
func (c *Http) GetSongLyrics(s provider.SearchParams) (models.Lyrics, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/lyrics?%s&q=%s&song=%s&artist=%s&album=%s",
			c.apiAddress, c.providers, url.QueryEscape(s.Query), url.QueryEscape(s.SongName), url.QueryEscape(s.ArtistName), url.QueryEscape(s.AlbumName),
		),
		http.NoBody)
	if err != nil {
		return models.Lyrics{}, err
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return models.Lyrics{}, err
	}

	var lyrics models.Lyrics
	err = json.NewDecoder(resp.Body).Decode(&lyrics)
	if err != nil {
		return models.Lyrics{}, err
	}
	_ = resp.Body.Close()

	return lyrics, nil
}
