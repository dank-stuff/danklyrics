package main

import (
	"encoding/json"
	"fmt"

	"codeberg.org/dankstuff/danklyrics/pkg/lrclib"
)

func main() {
	client := lrclib.NewClient()

	song, err := client.Songs.Get("4809799")
	if err != nil {
		panic(err)
	}

	songJson, err := json.MarshalIndent(song, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf("song: %s\n\n", songJson)

	songByParams, err := client.Songs.GetByParams(lrclib.GetSongParams{
		ArtistName: "Abba",
		TrackName:  "waterloo",
	})
	if err != nil {
		panic(err)
	}

	songByParamsJson, err := json.MarshalIndent(songByParams, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf("song by params: %s\n", songByParamsJson)
}
