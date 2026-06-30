package main

import (
	"fmt"
	"os"

	"codeberg.org/dankstuff/danklyrics/pkg/gonius"
)

func main() {
	client := gonius.NewClient(os.Getenv("GENIUS_CLIENT_ID"), os.Getenv("GENIUS_CLIENT_SECRET"))
	lyrics, err := client.Lyrics.FindForSong("https://genius.com/Adele-set-fire-to-the-rain-lyrics")
	if err != nil {
		panic(err)
	}

	fmt.Println(lyrics.String())
}
