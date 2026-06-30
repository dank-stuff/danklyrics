package main

import (
	"encoding/json"
	"fmt"

	"codeberg.org/dankstuff/danklyrics/pkg/lrclib"
)

func main() {
	client := lrclib.NewClient()
	results, err := client.Search.Get(lrclib.SearchParams{
		Query: "lana del rey jealous girl",
		Limit: 5,
	})
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		jsonn, _ := json.MarshalIndent(result, "", "\t")
		fmt.Println("search result", string(jsonn))
	}
}
