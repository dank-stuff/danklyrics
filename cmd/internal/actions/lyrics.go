package actions

import (
	"errors"

	intmodels "codeberg.org/dankstuff/danklyrics/cmd/internal/models"
	"codeberg.org/dankstuff/danklyrics/pkg/models"
	"codeberg.org/dankstuff/danklyrics/pkg/provider"
)

func (a *Actions) GetLyricsByPublicId(id string) (models.Lyrics, error) {
	intLyrics, err := a.repo.GetLyricsByPublicId(id)
	if err != nil {
		return models.Lyrics{}, err
	}

	return models.Lyrics{
		PublicId:   intLyrics.PublicId,
		SongName:   intLyrics.SongTitle,
		ArtistName: intLyrics.ArtistName,
		AlbumName:  intLyrics.AlbumTitle,
		Parts:      intLyrics.LyricsPlain,
		Synced:     intLyrics.LyricsSynced,
	}, nil
}

type FindLyricsParams struct {
	SongTitle  string
	ArtistName string
	AlbumTitle string
}

func (a *Actions) FindLyrics(search FindLyricsParams) ([]models.Lyrics, error) {
	findParams := FindLyricsArgs{
		SongTitle:  search.SongTitle,
		ArtistName: search.ArtistName,
		AlbumTitle: search.AlbumTitle,
	}

	intLyricses, err := a.repo.FindLyricsExact(findParams)
	if err != nil {
		intLyricses, err = a.repo.FindLyricsAll(findParams)
		if err != nil {
			return nil, err
		}
	}

	lyricses := make([]models.Lyrics, 0, len(intLyricses))
	for _, intLyrics := range intLyricses {
		lyricses = append(lyricses, models.Lyrics{
			SongName:   intLyrics.SongTitle,
			ArtistName: intLyrics.ArtistName,
			AlbumName:  intLyrics.AlbumTitle,
			Parts:      intLyrics.LyricsPlain,
			Synced:     intLyrics.LyricsSynced,
		})
	}

	return lyricses, nil
}

type FindLyricsProviderParams struct {
	SongTitle  string
	ArtistName string
	AlbumTitle string
	Query      string
	Lyricser   provider.Service
}

func (a *Actions) FindLyricsProvider(params FindLyricsProviderParams) ([]models.Lyrics, error) {
	lyricses, err := a.FindLyrics(FindLyricsParams{
		SongTitle:  params.SongTitle,
		ArtistName: params.ArtistName,
		AlbumTitle: params.AlbumTitle,
	})

	if len(lyricses) > 0 {
		return lyricses, nil
	}

	if params.Lyricser == nil {
		return nil, &ErrNoResultsFound{}
	}

	lyrics, err := params.Lyricser.GetSongLyrics(provider.SearchParams{
		SongName:   params.SongTitle,
		ArtistName: params.ArtistName,
		AlbumName:  params.AlbumTitle,
		Query:      params.Query,
	})
	if err != nil {
		return nil, err
	}
	if len(lyricses) == 0 && len(lyrics.Parts) > 0 {
		_, _ = a.CreateLyrics(lyrics)
	}

	return []models.Lyrics{lyrics}, nil
}

func (a *Actions) CreateLyrics(l models.Lyrics) (models.Lyrics, error) {
	if l.SongName == "" {
		return models.Lyrics{}, errors.New("missing song name")
	}

	intLyrics := intmodels.Lyrics{
		SongTitle:    l.SongName,
		ArtistName:   l.ArtistName,
		AlbumTitle:   l.AlbumName,
		LyricsPlain:  l.Parts,
		LyricsSynced: l.Synced,
	}

	newLyrics, err := a.repo.CreateLyrics(intLyrics)
	if err != nil {
		return models.Lyrics{}, err
	}

	_ = a.sitemap.AddLyricsEntry(SitemapUrl{
		PublicId: newLyrics.PublicId,
		LastMod:  newLyrics.CreatedAt,
	})

	return models.Lyrics{}, nil
}

func (a *Actions) CreateLyricsRequest(token string, l models.Lyrics) error {
	tokenDecoded, err := a.jwt.Decode(token, JwtSessionToken)
	if err != nil {
		return err
	}

	if l.SongName == "" {
		return errors.New("missing song name")
	}

	intLyrics := intmodels.LyricsRequest{
		SongTitle:      l.SongName,
		ArtistName:     l.ArtistName,
		AlbumTitle:     l.AlbumName,
		LyricsPlain:    l.Parts,
		LyricsSynced:   l.Synced,
		RequesterEmail: tokenDecoded.Payload.Email,
	}

	_, err = a.repo.CreateLyricsRequest(intLyrics)
	if err != nil {
		return err
	}

	return nil
}
