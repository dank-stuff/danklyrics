# lrclibgo

**lrclibgo** is a Go client to access the [LyricFind API](https://lrclib.net/docs).

# Roadmap

- [x] Search
- [x] Get song
- [x] Lyrics
- [x] Lyrics timing
- [ ] Publish lyrics

# Usage

```go
package main

import (
	"encoding/json"
	"fmt"

	"codeberg.org/dankstuff/danklyrics/pkg/lrclib"
)

func main() {
	client := lrclibgo.NewClient()
	results, err := client.Search.Get(lrclibgo.SearchParams{
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
```
