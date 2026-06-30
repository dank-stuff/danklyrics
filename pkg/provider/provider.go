package provider

import (
	"fmt"
	"net/http"

	"codeberg.org/dankstuff/danklyrics/pkg/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// SearchParams holds the search criteria to find a song from a provider.
type SearchParams struct {
	SongName   string
	ArtistName string
	AlbumName  string
	Query      string
}

// Service fetches lyrics for the given song in the search params.
type Service interface {
	// GetSongLyrics searches for a song's lyrics using the given search params.
	GetSongLyrics(s SearchParams) (models.Lyrics, error)
	// GetSongsLyrics same as [GetSongLyrics] but returns all the songs in the search results with their lyrics.
	// GetSongsLyrics(s SearchParams) ([]models.Song, error)
}

// Name represents lyrics finding providers to choose from when doing a lyrics search.
type Name string

const (
	// Dank pass this to [GetSongLyrics] to use DankLyrics as a lyrics provider.
	Dank Name = "dank"
	// LyricFind pass this to [GetSongLyrics] to use LyricFind as a lyrics provider.
	LyricFind Name = "lrc"
	// Genius pass this to [GetSongLyrics] to use Genius as a lyrics provider.
	Genius Name = "genius"
)

// Auth represents a provider's needed auth credentials
type Auth struct {
	ClientId     string
	ClientSecret string
}

func clientIdHttpHeaderKey(name Name) string {
	return fmt.Sprintf(
		"X-%s-Auth-Client-Id",
		cases.
			Title(language.English).
			String(string(name)),
	)
}

func clientSecretHttpHeaderKey(name Name) string {
	return fmt.Sprintf(
		"X-%s-Auth-Client-Secret",
		cases.
			Title(language.English).
			String(string(name)),
	)
}

// HttpHeaders returns a map containing the requirered http headers to authenticate with a provider.
func (a Auth) HttpHeaders(name Name) http.Header {
	return http.Header{
		clientIdHttpHeaderKey(name):     []string{a.ClientId},
		clientSecretHttpHeaderKey(name): []string{a.ClientSecret},
	}
}

func AuthFromHttpHeaders(name Name, headers http.Header) Auth {
	a := Auth{}
	if clientId, ok := headers[clientIdHttpHeaderKey(name)]; ok {
		a.ClientId = clientId[0]
	}
	if clientSecret, ok := headers[clientSecretHttpHeaderKey(name)]; ok {
		a.ClientSecret = clientSecret[0]
	}

	return a
}
