package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"

	"codeberg.org/dankstuff/danklyrics/pkg/models"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"
)

// Http is the dank lyrics finding client that makes a call to api.danklyrics.com to find the lyrics.
type Http struct {
	providers        string
	apiAddress       string
	providersHeaders http.Header
}

func NewHttp(c Config) (*Http, error) {
	if len(c.Providers) == 0 {
		return nil, errors.New("must specify at least one lyrics provider")
	}

	client := &Http{
		providersHeaders: make(http.Header),
	}

	client.providers = ""
	for i, p := range c.Providers {
		client.providers += "providers=" + string(p)
		if i < len(c.Providers)-1 {
			client.providers += "&"
		}
	}
	for p, auth := range c.ProvidersAuth {
		maps.Copy(client.providersHeaders, auth.HttpHeaders(p))
	}

	if c.ApiAddress == "" {
		client.apiAddress = "https://api.danklyrics.com"
	} else {
		client.apiAddress = c.ApiAddress
	}

	return client, nil
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
	maps.Copy(req.Header, c.providersHeaders)

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
