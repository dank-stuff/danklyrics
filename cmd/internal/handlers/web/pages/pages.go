package pages

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"codeberg.org/dankstuff/danklyrics/cmd/internal/actions"
	"codeberg.org/dankstuff/danklyrics/cmd/website/layouts"
	"codeberg.org/dankstuff/danklyrics/cmd/website/partials"
	"codeberg.org/dankstuff/danklyrics/cmd/website/types"
	"codeberg.org/dankstuff/danklyrics/pkg/models"

	"github.com/a-h/templ"
)

type pages struct {
	usecases *actions.Actions
}

func New(usecases *actions.Actions) *pages {
	return &pages{
		usecases: usecases,
	}
}

func (p *pages) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".go") || strings.Contains(r.URL.Path, "admin") {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err := layouts.Default(types.PageProps{
		PageId: types.FindLyricsPage,
		Title:  "Find lyrics",
	}, templ.NopComponent).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (p *pages) HandleAbout(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".go") || strings.Contains(r.URL.Path, "admin") {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err := layouts.Default(types.PageProps{
		PageId: types.AboutPage,
		Title:  "About",
	}, partials.About()).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (p *pages) HandleSubmitLyrics(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".go") || strings.Contains(r.URL.Path, "admin") {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	sessionToken, err := r.Cookie("token")
	if err != nil || sessionToken == nil {
		_ = layouts.Default(types.PageProps{
			PageId: types.SubmitLyricsPage,
			Title:  "Submit Lyrics",
		}, partials.SubmitLyricsAuth()).Render(r.Context(), w)
		return
	}

	if err := p.usecases.ConfirmAuth(sessionToken.Value); err != nil {
		_ = layouts.Default(types.PageProps{
			PageId: types.SubmitLyricsPage,
			Title:  "Submit Lyrics",
		}, partials.SubmitLyricsAuth()).Render(r.Context(), w)
		return
	}

	err = layouts.Default(types.PageProps{
		PageId: types.SubmitLyricsPage,
		Title:  "Submit Lyrics",
	}, partials.SubmitLyrics()).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (p *pages) HandleAboutTab(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".go") || strings.Contains(r.URL.Path, "admin") {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err := partials.About().Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (p *pages) HandleSubmitLyricsTab(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".go") || strings.Contains(r.URL.Path, "admin") {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	sessionToken, err := r.Cookie("token")
	if err != nil || sessionToken == nil {
		_ = partials.SubmitLyricsAuth().Render(r.Context(), w)
		return
	}

	if err := p.usecases.ConfirmAuth(sessionToken.Value); err != nil {
		_ = partials.SubmitLyricsAuth().Render(r.Context(), w)
		return
	}

	err = partials.SubmitLyrics().Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (p *pages) HandleLyrics(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".go") || strings.Contains(r.URL.Path, "/admin") {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	lyricsSlug := r.PathValue("id")
	lyrics, err := p.usecases.GetLyricsByPublicId(lyricsSlug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = layouts.Default(types.PageProps{
			PageId: types.LyricsPage,
			Title:  "Not found",
		}, partials.SingleLyrics(models.Lyrics{})).Render(r.Context(), w)
		return
	}

	err = layouts.Default(types.PageProps{
		PageId:      types.LyricsPage,
		Title:       fmt.Sprintf("%s, %s Lyrics", lyrics.ArtistName, lyrics.SongName),
		Description: fmt.Sprintf("%s by %s from the album %s", lyrics.SongName, lyrics.ArtistName, lyrics.AlbumName),
		Url:         "https://danklyrics.com/lyrics/" + lyrics.PublicId,
		Audio: types.AudioProps{
			Album:     lyrics.AlbumName,
			Musician:  lyrics.ArtistName,
			SongTitle: lyrics.SongName,
		},
	}, partials.SingleLyrics(lyrics)).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
