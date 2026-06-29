package web

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"

	"codeberg.org/dankstuff/danklyrics/internal/actions"
	"codeberg.org/dankstuff/danklyrics/internal/handlers/web/templates"
	"codeberg.org/dankstuff/danklyrics/pkg/models"
	website "codeberg.org/dankstuff/danklyrics/website/user"

	"github.com/a-h/templ"
)

//go:embed sitemap_template.xml
var sitemapTemplate embed.FS

var publicFiles embed.FS

func init() {
	publicFiles = website.FS()
}

type pages struct {
	usecases *actions.Actions
}

func NewPages(usecases *actions.Actions) *pages {
	return &pages{
		usecases: usecases,
	}
}

func (p *pages) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".go") || strings.Contains(r.URL.Path, "admin") {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
		return
	}
	if strings.HasSuffix(r.URL.Path, "favicon.ico") {
		f, err := publicFiles.Open("favicon.ico")
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "image/x-icon")
		io.Copy(w, f)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err := templates.Layout(templates.PageProps{
		PageId: templates.FindLyricsPage,
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
	err := templates.Layout(templates.PageProps{
		PageId: templates.AboutPage,
		Title:  "About",
	}, templates.About()).Render(r.Context(), w)
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
		_ = templates.Layout(templates.PageProps{
			PageId: templates.SubmitLyricsPage,
			Title:  "Submit Lyrics",
		}, templates.SubmitLyricsAuth()).Render(r.Context(), w)
		return
	}

	if err := p.usecases.ConfirmAuth(sessionToken.Value); err != nil {
		_ = templates.Layout(templates.PageProps{
			PageId: templates.SubmitLyricsPage,
			Title:  "Submit Lyrics",
		}, templates.SubmitLyricsAuth()).Render(r.Context(), w)
		return
	}

	err = templates.Layout(templates.PageProps{
		PageId: templates.SubmitLyricsPage,
		Title:  "Submit Lyrics",
	}, templates.SubmitLyrics()).Render(r.Context(), w)
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
	err := templates.About().Render(r.Context(), w)
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
		_ = templates.SubmitLyricsAuth().Render(r.Context(), w)
		return
	}

	if err := p.usecases.ConfirmAuth(sessionToken.Value); err != nil {
		_ = templates.SubmitLyricsAuth().Render(r.Context(), w)
		return
	}

	err = templates.SubmitLyrics().Render(r.Context(), w)
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
		_ = templates.Layout(templates.PageProps{
			PageId: templates.LyricsPage,
			Title:  "Not found",
		}, templates.SingleLyrics(models.Lyrics{})).Render(r.Context(), w)
		return
	}

	err = templates.Layout(templates.PageProps{
		PageId:      templates.LyricsPage,
		Title:       fmt.Sprintf("%s, %s Lyrics", lyrics.ArtistName, lyrics.SongName),
		Description: fmt.Sprintf("%s by %s from the album %s", lyrics.SongName, lyrics.ArtistName, lyrics.AlbumName),
		Url:         "https://danklyrics.com/lyrics/" + lyrics.PublicId,
		Audio: templates.AudioProps{
			Album:     lyrics.AlbumName,
			Musician:  lyrics.ArtistName,
			SongTitle: lyrics.SongName,
		},
	}, templates.SingleLyrics(lyrics)).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (p *pages) HandleSitemap(w http.ResponseWriter, r *http.Request) {
	sitemapEntries, err := makeApiGetRequest[[]actions.SitemapUrl]("/sitemap-kurwa", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Cache-Control", "max-age=300")
	w.Header().Set("Content-Type", "application/xml")

	t := template.Must(template.ParseFS(sitemapTemplate, "sitemap_template.xml"))
	err = t.Execute(w, sitemapEntries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *pages) HandleRobots(w http.ResponseWriter, r *http.Request) {
	robotsFile, _ := publicFiles.Open("robots.txt")
	w.Header().Set("Content-Type", "text/plain")
	_, _ = io.Copy(w, robotsFile)
}

func (p *pages) StaticFiles() http.Handler {
	return http.StripPrefix("/static", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, ".go") || strings.Contains(r.URL.Path, "admin") {
			http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
			return
		}

		http.FileServer(http.FS(publicFiles)).ServeHTTP(w, r)
	}))
}
