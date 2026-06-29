package main

import (
	"log"
	"net/http"

	"codeberg.org/dankstuff/danklyrics/internal/actions"
	"codeberg.org/dankstuff/danklyrics/internal/config"
	"codeberg.org/dankstuff/danklyrics/internal/handlers/api"
	"codeberg.org/dankstuff/danklyrics/internal/jwt"
	"codeberg.org/dankstuff/danklyrics/internal/mailer"
	"codeberg.org/dankstuff/danklyrics/internal/mariadb"
	"codeberg.org/dankstuff/danklyrics/internal/sitemap"
)

var (
	usecases *actions.Actions
)

func init() {
	repo, err := mariadb.New()
	if err != nil {
		panic(err)
	}

	err = mariadb.Migrate()
	if err != nil {
		panic(err)
	}

	mailUtil := mailer.New()
	jwtUtil := jwt.New[actions.TokenPayload]()
	sm := sitemap.New()
	usecases = actions.New(repo, mailUtil, jwtUtil, sm)

	err = usecases.LoadLyricsPublicIds()
	if err != nil {
		panic(err)
	}
}

func main() {
	apiHandler := http.NewServeMux()

	lyricsApi := api.NewLyricsFinderApi(usecases)
	dankLyricsApi := api.NewDankLyricsApi(usecases)
	authApi := api.NewAuthApi(usecases)
	sitemapApi := api.NewSitemapApi(usecases)

	apiHandler.HandleFunc("/", lyricsApi.HandleIndex)
	apiHandler.HandleFunc("GET /providers", lyricsApi.HandleListProviders)
	apiHandler.HandleFunc("GET /lyrics", lyricsApi.HandleGetSongLyrics)
	apiHandler.HandleFunc("POST /dank/lyrics", dankLyricsApi.HandleSubmitSongLyrics)
	apiHandler.HandleFunc("GET /dank/lyrics", dankLyricsApi.HandleGetSongLyrics)
	apiHandler.HandleFunc("POST /auth", authApi.HandleAuth)
	apiHandler.HandleFunc("POST /auth/confirm", authApi.HandleConfirmAuth)
	apiHandler.HandleFunc("GET /sitemap-kurwa", sitemapApi.HandleGetSitemapEntries)

	log.Printf("Starting web server at port %s", config.Env().ApiPort)
	log.Fatalln(http.ListenAndServe(":"+config.Env().ApiPort, apiHandler))
}
