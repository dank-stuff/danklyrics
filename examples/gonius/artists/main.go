package main

import (
	"encoding/json"
	"fmt"
	"os"

	"codeberg.org/dankstuff/danklyrics/pkg/gonius"
)

func main() {
	client := gonius.NewClient(os.Getenv("GENIUS_CLIENT_ID"), os.Getenv("GENIUS_CLIENT_SECRET"))
	artist, err := client.Artists.Get("15740")
	if err != nil {
		panic(err)
	}

	artistJson, err := json.MarshalIndent(artist, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf("artist: %s\n", artistJson)

	artistSongs, err := client.Artists.GetSongs("15740", gonius.ArtistSongsSortPopularity)
	if err != nil {
		panic(err)
	}

	artistSongsJson, err := json.MarshalIndent(artistSongs, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf("artist songs: %s\n", artistSongsJson)
}
