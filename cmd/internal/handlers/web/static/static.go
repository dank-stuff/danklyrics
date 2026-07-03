package static

import (
	"embed"
	"net/http"
	"strings"
	"text/template"
	"time"

	"codeberg.org/dankstuff/danklyrics/cmd/internal/actions"
	static "codeberg.org/dankstuff/danklyrics/cmd/website/static/user"
	"github.com/tdewolff/minify/v2"
)

//go:embed sitemap_template.xml
var sitemapTemplate embed.FS

var publicFiles embed.FS

func init() {
	publicFiles = static.FS()
}

type staticHandler struct {
	usecases *actions.Actions
	minifyer *minify.M
}

func New(usecases *actions.Actions, minifyer *minify.M) *staticHandler {
	return &staticHandler{
		usecases: usecases,
		minifyer: minifyer,
	}
}

func (s *staticHandler) HandleRobots(w http.ResponseWriter, r *http.Request) {
	robotsFile, _ := publicFiles.ReadFile("robots.txt")
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write(robotsFile)
}

type sitemapValue struct {
	PublicId string
	LastMod  string
}

func (s *staticHandler) HandleSitemap(w http.ResponseWriter, r *http.Request) {
	sitemapEntries, err := s.usecases.GetSitemap()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=15")
	w.Header().Set("Content-Type", "application/xml")

	sitemapValues := make([]sitemapValue, 0, len(sitemapEntries))
	for _, entry := range sitemapEntries {
		sitemapValues = append(sitemapValues, sitemapValue{
			PublicId: entry.PublicId,
			LastMod:  entry.LastMod.UTC().Format(time.RFC3339),
		})
	}

	t := template.Must(template.ParseFS(sitemapTemplate, "sitemap_template.xml"))
	err = t.Execute(w, sitemapValues)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *staticHandler) HandleFavicon(w http.ResponseWriter, r *http.Request) {
	faviconFile, _ := publicFiles.ReadFile("favicon.ico")
	w.Header().Set("Content-Type", "image/x-icon")
	_, _ = w.Write(faviconFile)
}

func (s *staticHandler) AssetsHandler() http.Handler {
	return http.StripPrefix("/static", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, ".go") || strings.Contains(r.URL.Path, "admin") {
			http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
			return
		}

		w.Header().Set("Cache-Control", "public, max-age=7200, stale-while-revalidate=5")

		s.minifyer.Middleware(http.FileServer(http.FS(publicFiles))).
			ServeHTTP(w, r)
	}))

}
