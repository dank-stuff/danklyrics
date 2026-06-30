package main

import (
	"fmt"

	"codeberg.org/dankstuff/danklyrics/pkg/lrclib"
)

func main() {
	client := lrclib.NewClient()

	song, err := client.Songs.Get("4809799")
	if err != nil {
		panic(err)
	}

	lyrics := song.Lyrics()
	fmt.Println("plain", lyrics.String())
	fmt.Println("parts", lyrics.Parts)
	fmt.Println("synced", lyrics.Synced)
}
