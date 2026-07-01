package actions

import "codeberg.org/dankstuff/danklyrics/internal/models"

type FindLyricsParams struct {
	SongTitle  string
	ArtistName string
	AlbumTitle string
}

type Repository interface {
	CreateLyrics(l models.Lyrics) (models.Lyrics, error)
	GetLyricsByPublicId(id string) (models.Lyrics, error)

	FindLyricsExact(search FindLyricsParams) ([]models.Lyrics, error)
	FindLyricsAll(search FindLyricsParams) ([]models.Lyrics, error)

	GetLyricses(page int) ([]models.Lyrics, error)

	CreateLyricsRequest(l models.LyricsRequest) (models.LyricsRequest, error)
	DeleteLyricsRequest(id uint) error
	GetLyricsRequestById(id uint) (models.LyricsRequest, error)
	GetLyricsRequests() ([]models.LyricsRequest, error)
	GetAdminByUsername(username string) (models.Admin, error)
}
