package main

import (
	"log"
	"net/http"
	"regexp"

	"codeberg.org/dankstuff/danklyrics/internal/actions"
	"codeberg.org/dankstuff/danklyrics/internal/config"
	"codeberg.org/dankstuff/danklyrics/internal/handlers/web"
	"codeberg.org/dankstuff/danklyrics/internal/jwt"
	"codeberg.org/dankstuff/danklyrics/internal/mailer"
	"codeberg.org/dankstuff/danklyrics/internal/mariadb"
	"codeberg.org/dankstuff/danklyrics/internal/sitemap"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	mjson "github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

var (
	minifyer *minify.M
	usecases *actions.Actions
)

func init() {
	minifyer = minify.New()
	minifyer.AddFunc("text/css", css.Minify)
	minifyer.AddFunc("text/html", html.Minify)
	minifyer.AddFunc("image/svg+xml", svg.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]json$"), mjson.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	repo, err := mariadb.New()
	if err != nil {
		panic(err)
	}

	mailUtil := mailer.New()
	jwtUtil := jwt.New[actions.TokenPayload]()
	sm := sitemap.New()
	usecases = actions.New(repo, mailUtil, jwtUtil, sm)
}

func main() {
	pages := web.NewPages(usecases)
	pagesHandler := http.NewServeMux()
	pagesHandler.HandleFunc("/", pages.HandleIndex)
	pagesHandler.HandleFunc("/about", pages.HandleAbout)
	pagesHandler.HandleFunc("/lyrics/{id}", pages.HandleLyrics)
	pagesHandler.HandleFunc("/lyrics/submit", pages.HandleSubmitLyrics)
	pagesHandler.HandleFunc("/tab/about", pages.HandleAboutTab)
	pagesHandler.HandleFunc("/tab/lyrics/submit", pages.HandleSubmitLyricsTab)
	pagesHandler.HandleFunc("/sitemap.xml", pages.HandleSitemap)
	pagesHandler.HandleFunc("/robots.txt", pages.HandleRobots)
	pagesHandler.Handle("/static/", pages.StaticFiles())

	apis := web.NewApi(usecases)
	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("GET /lyrics", apis.HandleGetSongLyrics)
	apisHandler.HandleFunc("POST /lyrics", apis.HandleSubmitLyrics)
	apisHandler.HandleFunc("POST /auth", apis.HandleAuthSubmitLyrics)
	apisHandler.HandleFunc("GET /auth/confirm", apis.HandleConfirmAuthSubmitLyrics)

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", minifyer.Middleware(pagesHandler))
	applicationHandler.Handle("/api/", http.StripPrefix("/api", apisHandler))

	log.Printf("Starting web server at port %s", config.Env().Port)
	log.Fatalln(http.ListenAndServe(":"+config.Env().Port, applicationHandler))
}
