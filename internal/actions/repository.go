package actions

import "codeberg.org/dankstuff/danklyrics/internal/models"

type Repository interface {
	CreateLyrics(l models.Lyrics) (models.Lyrics, error)
	GetLyricsByPublicId(id string) (models.Lyrics, error)
	GetLyricsBySongTitle(Title string) ([]models.Lyrics, error)
	GetLyricsBySongTitleAndArtistName(songTitle, artistName string) ([]models.Lyrics, error)
	GetLyricsBySongAndAlbumTitle(songTitle, albumTitle string) ([]models.Lyrics, error)
	GetLyricsBySongTitleArtistNameAndAlbumTitle(songTitle, artistName, albumTitle string) ([]models.Lyrics, error)

	GetLyricses(page int) ([]models.Lyrics, error)

	CreateLyricsRequest(l models.LyricsRequest) (models.LyricsRequest, error)
	DeleteLyricsRequest(id uint) error
	GetLyricsRequestById(id uint) (models.LyricsRequest, error)
	GetLyricsRequests() ([]models.LyricsRequest, error)
	GetAdminByUsername(username string) (models.Admin, error)
}
