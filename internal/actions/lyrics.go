package actions

import (
	"errors"
	"time"

	intmodels "codeberg.org/dankstuff/danklyrics/internal/models"
	"codeberg.org/dankstuff/danklyrics/pkg/models"
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

func (a *Actions) GetLyricsBySongTitle(title string) ([]models.Lyrics, error) {
	findParams := FindLyricsParams{
		SongTitle: title,
	}
	intLyricses, err := a.repo.FindLyricsAll(findParams)
	if err != nil {
		return nil, err
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

func (a *Actions) GetLyricsBySongTitleAndArtistName(title, artistName string) ([]models.Lyrics, error) {
	findParams := FindLyricsParams{
		SongTitle:  title,
		ArtistName: artistName,
	}
	intLyricses, err := a.repo.FindLyricsAll(findParams)
	if err != nil {
		return nil, err
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

func (a *Actions) GetLyricsBySongTitleAndAlbumTitle(title, albumTitle string) ([]models.Lyrics, error) {
	findParams := FindLyricsParams{
		SongTitle:  title,
		AlbumTitle: albumTitle,
	}
	intLyricses, err := a.repo.FindLyricsAll(findParams)
	if err != nil {
		return nil, err
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

func (a *Actions) GetLyricsBySongTitleArtistNameAndAlbumTitle(title, artistName, albumTitle string) ([]models.Lyrics, error) {
	findParams := FindLyricsParams{
		SongTitle:  title,
		ArtistName: artistName,
		AlbumTitle: albumTitle,
	}
	intLyricses, err := a.repo.FindLyricsAll(findParams)
	if err != nil {
		return nil, err
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
		AddedAt:  newLyrics.CreatedAt.Format(time.RFC3339),
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
