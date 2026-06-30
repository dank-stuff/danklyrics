package genius

import (
	"fmt"
	"strings"

	"codeberg.org/dankstuff/danklyrics/pkg/errors"
	"codeberg.org/dankstuff/danklyrics/pkg/models"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"

	"codeberg.org/dankstuff/danklyrics/pkg/gonius"
)

type geniusProvider struct {
	client *gonius.Client
}

func New(clientId, clientSecret string) provider.Service {
	return &geniusProvider{
		client: gonius.NewClient(clientId, clientSecret),
	}
}

func (g *geniusProvider) GetSongLyrics(s provider.SearchParams) (models.Lyrics, error) {
	query := new(strings.Builder)
	if s.ArtistName != "" {
		fmt.Fprintf(query, "%s ", s.ArtistName)
	}
	query.WriteString(s.SongName)
	if s.Query != "" {
		query.Reset()
		query.WriteString(s.Query)
	}

	hits, err := g.client.Search.Get(query.String())
	if err != nil {
		return models.Lyrics{}, err
	}
	if len(hits) == 0 {
		return models.Lyrics{}, &errors.ErrNotFound{}
	}

	var bestMatch = hits[0]
	for _, hit := range hits {
		if hit.Result == nil {
			continue
		}
		if strings.Contains(
			strings.ToLower(hit.Result.FullTitle), strings.ToLower(s.SongName)) &&
			hit.Result.PrimaryArtist != nil &&
			strings.Contains(
				strings.ToLower(hit.Result.PrimaryArtist.Name), strings.ToLower(s.ArtistName)) {
			bestMatch = hit
		}
	}

	lyrics, err := g.client.Lyrics.FindForSong(bestMatch.Result.URL)
	if err != nil {
		return models.Lyrics{}, err
	}

	return models.Lyrics{
		SongName:   bestMatch.Result.Title,
		ArtistName: bestMatch.Result.PrimaryArtist.Name,
		Parts:      lyrics.Parts,
	}, nil
}
