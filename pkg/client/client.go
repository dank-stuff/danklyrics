package client

import (
	"codeberg.org/dankstuff/danklyrics/internal/providers/dank"
	"codeberg.org/dankstuff/danklyrics/internal/providers/lyricfind"
	"codeberg.org/dankstuff/danklyrics/pkg/finder"
	"codeberg.org/dankstuff/danklyrics/pkg/models"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"
)

// Local is the dank lyrics finding client that uses [finder.Service] to find lyrics using the enabled providers.
type Local struct {
	finder *finder.Service
}

// Config holds the configs needed to initialize [Local] or [Http] clients.
type Config struct {
	Providers []provider.Name
	// ApiAddress only used by [Http] client, setting its value for [Local] client won't destroy the world, but it's pointless.
	// defaults to (https://api.danklyrics.com)
	ApiAddress string
}

// New initializes a new [Local] instance with the given configs.
func New(c Config) (*Local, error) {
	providers := make([]provider.Service, 0, len(c.Providers))

	for _, providerName := range c.Providers {
		switch providerName {
		case provider.LyricFind:
			providers = append(providers, lyricfind.New())
		case provider.Dank:
			providers = append(providers, dank.New())
		}
	}

	finder, err := finder.New(providers)
	if err != nil {
		return nil, err
	}

	return &Local{
		finder: finder,
	}, nil
}

// GetSongLyrics search for song's lyrics using the enabled providers list,
// where using a provider depends on the provider's order in that list.
//
// returns [Lyrics] and an occurring [error]
func (c *Local) GetSongLyrics(s provider.SearchParams) (models.Lyrics, error) {
	return c.finder.GetSongLyrics(s)
}
