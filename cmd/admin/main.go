package main

import (
	"log"
	"net/http"

	"codeberg.org/dankstuff/danklyrics/internal/actions"
	"codeberg.org/dankstuff/danklyrics/internal/config"
	"codeberg.org/dankstuff/danklyrics/internal/handlers/admin"
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
}

func main() {
	pages := admin.NewAdminPages(usecases)
	pagesHandler := http.NewServeMux()
	pagesHandler.Handle("/static/", pages.StaticFiles())
	pagesHandler.HandleFunc("/", pages.HandleIndex)
	pagesHandler.HandleFunc("/robots.txt", pages.HandleIndex)

	apis := admin.NewAdminApi(usecases)
	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("GET /lyrics/requests", apis.HandleListLyricsRequests)
	apisHandler.HandleFunc("GET /lyrics/request/{id}", apis.HandleGetLyricsRequest)
	apisHandler.HandleFunc("POST /lyrics/request/approve/{id}", apis.HandleApproveLyricsRequest)
	apisHandler.HandleFunc("POST /lyrics/request/reject/{id}", apis.HandleRejectLyricsRequest)
	apisHandler.HandleFunc("POST /auth", apis.HandleAuthenticate)

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", pagesHandler)
	applicationHandler.Handle("/api/", http.StripPrefix("/api", apisHandler))

	log.Printf("Starting web server at port %s", config.Env().AdminPort)
	log.Fatalln(http.ListenAndServe(":"+config.Env().AdminPort, applicationHandler))
}
