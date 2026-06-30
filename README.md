<div align="center">
  <a href="https://danklyrics.com" target="_blank"><img src="https://danklyrics.com/static/favicon.png" width="150" /></a>

<h1>DankLyrics</h1>
  <p>
    <strong>A lyrics finder with the legendary cs1.6 theme.</strong>
  </p>
  <p>
    <a href="https://goreportcard.com/report/codeberg.org/dankstuff/danklyrics"><img alt="rex-deployment" src="https://goreportcard.com/badge/codeberg.org/dankstuff/danklyrics"/></a>
    <a href="https://godoc.org/codeberg.org/dankstuff/danklyrics"><img alt="rex-deployment" src="https://godoc.org/codeberg.org/dankstuff/danklyrics?status.png"/></a>
    <a href="https://codeberg.org/dankstuff/danklyrics/actions/workflows/rex-deploy.yml"><img alt="rex-deployment" src="https://codeberg.org/dankstuff/danklyrics/actions/workflows/rex-deploy.yml/badge.svg"/></a>
  </p>
</div>

## About

**DankLyrics:** A lyrics finder API, Website and Go package!

# Go Package Docs

DankLyrics provides a Go package, since the project is written in Go lol.

Here's a sample usage, it's pretty straight forward, as the client only has one
method :)

```go
package main

import (
	"codeberg.org/dankstuff/danklyrics/pkg/client"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"
)

func main() {
	lyricser, err := client.NewHttp(client.Config{
        // available providers are the following.
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
        SongName: "sos",
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
```

# REST API Docs

_Rest API is available at [api.danklyrics.com](https://api.danklyrics.com)_

- **`GET /`**:

_Displays this message_

```
refer to (https://codeberg.org/dankstuff/danklyrics) for API docs!
```

- **`GET /providers`**:

_Returns a list of the current supported lyrics providers_

```json
[
  {
    "name": "string: Name of the provider",
    "id": "string: id to specify the provider to use in /lyrics"
  }
]
```

- **`GET /lyrics`**:

_Finds lyrics for a song using the specified providers_

Query parameters

| name        | required              | description                                                                        |
| ----------- | --------------------- | ---------------------------------------------------------------------------------- |
| `providers` | required (at least 1) | to specify which lyrics provider(s) to use, list is fetched from `GET /providers`. |
| `song`      | required              | song's name to search for.                                                         |
| `artist`    | optional              | artist's name to search for, if the song's name isn't enough.                      |
| `album`     | optional              | album's name to search for, if the song's name isn't enough.                       |

```json
{
  "parts": ["string lyrics parts of the song"],
  "synced": { "time": "lyrics part" }
}
```

- **`GET /dank/lyrics`**:

_Find lyrics from DankLyrics' database, equivalent to using the Go client with
`provider.Dank` set_

Query parameters

| name     | required | description                                                   |
| -------- | -------- | ------------------------------------------------------------- |
| `song`   | required | song's name to search for.                                    |
| `artist` | optional | artist's name to search for, if the song's name isn't enough. |
| `album`  | optional | album's name to search for, if the song's name isn't enough.  |

```json
[
  {
    "song_name": "string: represents the song's name",
    "artist_name": "string: represents the song artist's name",
    "album_name": "string: represents the song album's name",
    "parts": ["string lyrics parts of the song"],
    "synced": { "time": "lyrics part" }
  }
]
```

---

A
[DankStuff <img height="16" width="16" src="https://dankstuff.net/assets/favicon.ico" />](https://dankstuff.net)
product!

Made with 🧉 by [Lord Baraa](https://lordbaraa.net).
