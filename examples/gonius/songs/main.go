package main

import (
	"encoding/json"
	"fmt"
	"os"

	"codeberg.org/dankstuff/danklyrics/pkg/gonius"
)

func main() {
	client := gonius.NewClient(os.Getenv("GENIUS_CLIENT_ID"), os.Getenv("GENIUS_CLIENT_SECRET"))
	song, err := client.Songs.Get("3130730")
	if err != nil {
		panic(err)
	}

	fmt.Printf("song.URL: %v\n", song.Description)

	songJson, err := json.MarshalIndent(song, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf("song: %s\n", songJson)

}
