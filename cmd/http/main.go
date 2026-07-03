package main

import (
	"log"
	"net/http"
	"regexp"

	"codeberg.org/dankstuff/danklyrics/cmd/internal/actions"
	"codeberg.org/dankstuff/danklyrics/cmd/internal/config"
	"codeberg.org/dankstuff/danklyrics/cmd/internal/handlers/apis"
	"codeberg.org/dankstuff/danklyrics/cmd/internal/handlers/middlewares/contenttype"
	"codeberg.org/dankstuff/danklyrics/cmd/internal/handlers/middlewares/ismobile"
	"codeberg.org/dankstuff/danklyrics/cmd/internal/handlers/middlewares/version"
	webapis "codeberg.org/dankstuff/danklyrics/cmd/internal/handlers/web/apis"
	"codeberg.org/dankstuff/danklyrics/cmd/internal/handlers/web/pages"
	"codeberg.org/dankstuff/danklyrics/cmd/internal/handlers/web/static"
	"codeberg.org/dankstuff/danklyrics/cmd/internal/jwt"
	"codeberg.org/dankstuff/danklyrics/cmd/internal/mailer"
	"codeberg.org/dankstuff/danklyrics/cmd/internal/mariadb"
	"codeberg.org/dankstuff/danklyrics/cmd/internal/sitemap"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
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
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	repo, err := mariadb.New()
	if err != nil {
		log.Panicln(err)
	}

	err = mariadb.Migrate()
	if err != nil {
		log.Panicln(err)
	}

	mailUtil := mailer.New()
	jwtUtil := jwt.New[actions.TokenPayload]()
	sm := sitemap.New()
	usecases = actions.New(repo, mailUtil, jwtUtil, sm)

	err = usecases.LoadLyricsPublicIds()
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	///
	/// REST APIS
	///

	lyricsApi := apis.NewLyricsFinderApi(usecases)
	dankLyricsApi := apis.NewDankLyricsApi(usecases)

	v1ApiHandler := http.NewServeMux()
	v1ApiHandler.HandleFunc("/", lyricsApi.HandleIndex)
	v1ApiHandler.HandleFunc("GET /providers", lyricsApi.HandleListProviders)
	v1ApiHandler.HandleFunc("GET /lyrics", lyricsApi.HandleGetSongLyrics)
	v1ApiHandler.HandleFunc("GET /dank/lyrics", dankLyricsApi.HandleGetSongLyrics)

	///
	/// PAGES
	///

	pages := pages.New(usecases)

	pagesHandler := http.NewServeMux()
	pagesHandler.HandleFunc("/", pages.HandleIndex)
	pagesHandler.HandleFunc("/about", pages.HandleAbout)
	pagesHandler.HandleFunc("/lyrics/{id}", pages.HandleLyrics)
	pagesHandler.HandleFunc("/lyrics/submit", pages.HandleSubmitLyrics)
	pagesHandler.HandleFunc("/tab/about", pages.HandleAboutTab)
	pagesHandler.HandleFunc("/tab/lyrics/submit", pages.HandleSubmitLyricsTab)

	///
	/// WEB APIS
	///

	webApis := webapis.New(usecases)

	webApisHandler := http.NewServeMux()
	webApisHandler.HandleFunc("GET /lyrics", webApis.HandleGetSongLyrics)
	webApisHandler.HandleFunc("POST /lyrics", webApis.HandleSubmitLyrics)
	webApisHandler.HandleFunc("POST /auth", webApis.HandleAuthSubmitLyrics)
	webApisHandler.HandleFunc("GET /auth/confirm", webApis.HandleConfirmAuthSubmitLyrics)

	///
	/// APPLICATION HANDLER
	///

	applicationHandler := http.NewServeMux()

	applicationHandler.Handle("/", version.Handler(ismobile.Handler(contenttype.Html(pagesHandler))))
	applicationHandler.Handle("/api/json/", contenttype.Json(http.StripPrefix("/api/json", v1ApiHandler)))
	applicationHandler.Handle("/api/web/", ismobile.Handler(http.StripPrefix("/api/web", webApisHandler)))

	///
	/// STATICS HANDLER
	///

	statics := static.New(usecases, minifyer)
	applicationHandler.HandleFunc("/robots.txt", statics.HandleRobots)
	applicationHandler.HandleFunc("/sitemap.xml", statics.HandleSitemap)
	applicationHandler.HandleFunc("/favicon.ico", statics.HandleFavicon)
	applicationHandler.Handle("/static/", (statics.AssetsHandler()))

	log.Printf("Starting web server at port %s", config.Env().Port)
	log.Fatalln(http.ListenAndServe(":"+config.Env().Port, applicationHandler))
}
