package main

import (
	"fmt"
	"os"

	"codeberg.org/dankstuff/danklyrics/pkg/client"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"
)

func main() {
	lyricser, err := client.NewHttp(client.Config{
		Providers: []provider.Name{provider.Dank, provider.LyricFind, provider.Genius},
		ProvidersAuth: map[provider.Name]provider.Auth{
			provider.Genius: provider.Auth{
				ClientId:     os.Getenv("GENIUS_CLIENT_ID"),
				ClientSecret: os.Getenv("GENIUS_CLIENT_SECRET"),
			},
		},
	})
	if err != nil {
		panic(err)
	}

	searchInput := provider.SearchParams{
		SongName:   "sos",
		ArtistName: "abba",
	}
	lyrics, err := lyricser.GetSongLyrics(searchInput)
	if err != nil {
		panic(err)
	}

	fmt.Println(lyrics.String())
	fmt.Println(lyrics.Parts)
	fmt.Println(lyrics.Synced)
}
