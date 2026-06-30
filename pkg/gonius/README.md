# gonius

**gonius** is a Go client to access the [Genius API](https://docs.genius.com/).

# Roadmap

- [x] Search
- [x] Get artist
- [x] Get artist's songs
- [x] Get annotation
- [x] Get song
- [x] Lyrics
- [ ] Pagination
- [ ] Account
- [ ] Find missing shit using [genius-lyrics](https://www.npmjs.com/package/genius-lyrics) as a reference
- [ ] Lyrics timing?

# Usage

```go
package main

import (
	"encoding/json"
	"fmt"

	"codeberg.org/dankstuff/danklyrics/pkg/gonius"
)

func main() {
	client := gonius.NewClient("top-secret-token-woo-scary")
	results, err := client.Search.Get("lana del rey jealous girl")
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		jsonn, _ := json.MarshalIndent(result, "", "\t")
		fmt.Println("search result", string(jsonn))
	}
}
```
