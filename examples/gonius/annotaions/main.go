package main

import (
	"encoding/json"
	"fmt"
	"os"

	"codeberg.org/dankstuff/danklyrics/pkg/gonius"
)

func main() {
	client := gonius.NewClient(os.Getenv("GENIUS_CLIENT_ID"), os.Getenv("GENIUS_CLIENT_SECRET"))
	annotation, err := client.Annotations.Get("10225840")
	if err != nil {
		panic(err)
	}

	annotationJson, err := json.MarshalIndent(annotation, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf("annotation: %s\n", annotationJson)
}
