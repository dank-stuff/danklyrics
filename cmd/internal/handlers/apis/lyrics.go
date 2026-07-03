package apis

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"codeberg.org/dankstuff/danklyrics/cmd/internal/actions"
	staticuser "codeberg.org/dankstuff/danklyrics/cmd/website/static/user"
	"codeberg.org/dankstuff/danklyrics/pkg/client"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"
)

type lyricsFinderApi struct {
	usecases *actions.Actions
}

func NewLyricsFinderApi(usecases *actions.Actions) *lyricsFinderApi {
	return &lyricsFinderApi{
		usecases: usecases,
	}
}

func (l *lyricsFinderApi) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "favicon.ico") {
		f, err := staticuser.FS().Open("favicon.ico")
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "image/x-icon")
		io.Copy(w, f)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("refer to (" + docsLink + ") for API docs!"))
}

func (l *lyricsFinderApi) HandleListProviders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode([]map[string]string{
		{"name": "DankLyrics", "id": "dank"},
		{"name": "LyricFind", "id": "lrc"},
		{"name": "Genius", "id": "genius"},
	})
}

func (l *lyricsFinderApi) HandleGetSongLyrics(w http.ResponseWriter, r *http.Request) {
	providers := r.URL.Query()["providers"]

	searchQuery, okSearchQuery := r.URL.Query()["q"]
	artistName, okArtist := r.URL.Query()["artist"]
	albumName, okAlbum := r.URL.Query()["album"]
	songName, okSong := r.URL.Query()["song"]

	w.Header().Set("Content-Type", "application/json")

	if len(providers) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:         "You must specify at least one provider",
			SuggestedAction: "Check the `GET /providers` endpoint.",
			DocsLink:        docsLink,
		})
		return
	}

	if !okArtist && !okAlbum && !okSong && !okSearchQuery {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:  "Missing all query parameters `artist`, `album` and `song` or just `q`",
			DocsLink: docsLink,
		})
		return
	}

	if !okSong && !okSearchQuery {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:  "Missing required query parameter `song` or `q`",
			DocsLink: docsLink,
		})
		return
	}

	searchInput := actions.FindLyricsProviderParams{}
	if okSong {
		searchInput.SongTitle = songName[0]
	}
	if okAlbum {
		searchInput.AlbumTitle = albumName[0]
	}
	if okArtist {
		searchInput.ArtistName = artistName[0]
	}
	if okSearchQuery {
		searchInput.Query = searchQuery[0]
	}

	providersConfig := make([]provider.Name, 0, len(providers))
	providersAuth := make(map[provider.Name]provider.Auth)
	for _, p := range providers {
		if p == string(provider.Dank) {
			continue
		}
		providersConfig = append(providersConfig, provider.Name(p))
		providersAuth[provider.Name(p)] = provider.AuthFromHttpHeaders(provider.Name(p), r.Header)
	}

	var err error
	searchInput.Lyricser, err = client.New(client.Config{
		Providers:     providersConfig,
		ProvidersAuth: providersAuth,
	})

	lyricses, err := l.usecases.FindLyricsProvider(searchInput)
	if err != nil {
		log.Println("oppsie doopsie some shit happened", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "No results were found",
		})
		return
	}

	_ = json.NewEncoder(w).Encode(lyricses[0])
}
