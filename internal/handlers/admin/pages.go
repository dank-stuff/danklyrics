package admin

import (
	"embed"
	"io"
	"log"
	"net/http"

	"codeberg.org/dankstuff/danklyrics/internal/actions"
	website "codeberg.org/dankstuff/danklyrics/website/admin"
)

var publicFiles embed.FS

func init() {
	publicFiles = website.FS()
}

type pages struct {
	usecases *actions.Actions
}

func NewAdminPages(usecases *actions.Actions) *pages {
	return &pages{
		usecases: usecases,
	}
}

func (p *pages) HandleIndex(w http.ResponseWriter, r *http.Request) {
	content, err := publicFiles.Open("index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, _ = io.Copy(w, content)
	_ = content.Close()
}

func (p *pages) HandleRobots(w http.ResponseWriter, r *http.Request) {
	content, err := publicFiles.Open("robots.txt")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = io.Copy(w, content)
	_ = content.Close()
}

func (p *pages) StaticFiles() http.Handler {
	return http.StripPrefix("/static", http.FileServer(http.FS(publicFiles)))
}
