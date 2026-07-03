package apis

import (
	"encoding/json"
	"net/http"
	"strings"

	"codeberg.org/dankstuff/danklyrics/cmd/internal/actions"
	"codeberg.org/dankstuff/danklyrics/pkg/models"
)

type dankLyricsApi struct {
	usecases *actions.Actions
}

func NewDankLyricsApi(usecases *actions.Actions) *dankLyricsApi {
	return &dankLyricsApi{
		usecases: usecases,
	}
}

func (d *dankLyricsApi) HandleGetSongLyrics(w http.ResponseWriter, r *http.Request) {
	artistName := r.URL.Query().Get("artist")
	albumTitle := r.URL.Query().Get("album")
	songTile := r.URL.Query().Get("song")

	w.Header().Set("Content-Type", "application/json")

	if songTile == "" && albumTitle == "" && artistName == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:  "Missing all query parameters `artist`, `album` and `song`",
			DocsLink: docsLink,
		})
		return
	}

	if songTile == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:  "Missing required query parameter `song`",
			DocsLink: docsLink,
		})
		return
	}

	lyricses, err := d.usecases.FindLyrics(actions.FindLyricsParams{
		SongTitle:  songTile,
		ArtistName: artistName,
		AlbumTitle: albumTitle,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:         "Something went wrong",
			SuggestedAction: "Check the docs, or contact admin (baraa@dankstuff.net)",
			DocsLink:        docsLink,
		})
		return
	}

	if len(lyricses) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "No results were found",
		})
		return
	}

	_ = json.NewEncoder(w).Encode(lyricses)
}

func (l *dankLyricsApi) HandleSubmitSongLyrics(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Header["Authorization"]
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "No no no, you need to authenticate first!",
		})
		return
	}

	var lyrics struct {
		SongName   string `json:"song_name"`
		ArtistName string `json:"artist_name"`
		AlbumName  string `json:"album_name"`
		Plain      string `json:"plain_lyrics"`
	}
	err := json.NewDecoder(r.Body).Decode(&lyrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "invalid request body",
		})
		return
	}

	if lyrics.SongName == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "missing required field `song_name`",
		})
		return
	}
	if lyrics.ArtistName == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "missing required field `artist_name`",
		})
		return
	}
	if lyrics.AlbumName == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "missing required field `album_name`",
		})
		return
	}
	if lyrics.Plain == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "missing required field `plain_lyrics`",
		})
		return
	}

	err = l.usecases.CreateLyricsRequest(token[0], models.Lyrics{
		SongName:   lyrics.SongName,
		ArtistName: lyrics.ArtistName,
		AlbumName:  lyrics.AlbumName,
		Parts:      strings.Split(lyrics.Plain, "\n"),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "Something went wrong",
		})
		return
	}
}
