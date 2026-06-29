package actions

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"time"

	"codeberg.org/dankstuff/danklyrics/internal/models"
)

type AuthenticateAdminParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticateAdminPayload struct {
	SessionToken string `json:"session_token"`
}

func sha512Hash(str string) string {
	hasher := sha512.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (a *Actions) AuthenticateAdmin(params AuthenticateAdminParams) (AuthenticateAdminPayload, error) {
	admin, err := a.repo.GetAdminByUsername(params.Username)
	if err != nil {
		return AuthenticateAdminPayload{}, err
	}

	adminPasswordHashed := admin.Password
	adminPasswordRequestHashed := sha512Hash(params.Password)

	if adminPasswordHashed != adminPasswordRequestHashed {
		return AuthenticateAdminPayload{}, errors.New("unauthorized")
	}

	sessionToken, err := a.jwt.Sign(TokenPayload{
		Email: admin.Username,
	}, JwtAdminToken, time.Hour*2)
	if err != nil {
		return AuthenticateAdminPayload{}, err
	}

	return AuthenticateAdminPayload{
		SessionToken: sessionToken,
	}, nil
}

type LyricsRequest struct {
	Id         uint              `json:"id"`
	SongTitle  string            `json:"song_title"`
	ArtistName string            `json:"artist_name"`
	AlbumTitle string            `json:"album_title"`
	Parts      []string          `json:"lyrics_plain,omitempty"`
	Synced     map[string]string `json:"lyrics_synced,omitempty"`
}

func (a *Actions) ListLyricsRequests(adminToken string) ([]LyricsRequest, error) {
	err := a.jwt.Validate(adminToken, JwtAdminToken)
	if err != nil {
		return nil, err
	}

	lrs, err := a.repo.GetLyricsRequests()
	if err != nil {
		return nil, err
	}

	lyricsRequests := make([]LyricsRequest, 0, len(lrs))
	for _, l := range lrs {
		lyricsRequests = append(lyricsRequests, LyricsRequest{
			Id:         l.Id,
			SongTitle:  l.SongTitle,
			ArtistName: l.ArtistName,
			AlbumTitle: l.AlbumTitle,
		})
	}

	return lyricsRequests, nil
}

func (a *Actions) GetLyricsRequest(adminToken string, id uint) (LyricsRequest, error) {
	err := a.jwt.Validate(adminToken, JwtAdminToken)
	if err != nil {
		return LyricsRequest{}, err
	}

	intLyrics, err := a.repo.GetLyricsRequestById(id)
	if err != nil {
		return LyricsRequest{}, err
	}

	return LyricsRequest{
		SongTitle:  intLyrics.SongTitle,
		ArtistName: intLyrics.ArtistName,
		AlbumTitle: intLyrics.AlbumTitle,
		Parts:      intLyrics.LyricsPlain,
		Synced:     intLyrics.LyricsSynced,
	}, nil
}

func (a *Actions) ApproveLyricsRequest(adminToken string, id uint) error {
	err := a.jwt.Validate(adminToken, JwtAdminToken)
	if err != nil {
		return err
	}

	lyricsRequest, err := a.repo.GetLyricsRequestById(id)
	if err != nil {
		return err
	}

	lyrics, err := a.repo.CreateLyrics(models.Lyrics{
		SongTitle:    lyricsRequest.SongTitle,
		ArtistName:   lyricsRequest.ArtistName,
		AlbumTitle:   lyricsRequest.AlbumTitle,
		LyricsPlain:  lyricsRequest.LyricsPlain,
		LyricsSynced: lyricsRequest.LyricsSynced,
	})
	if err != nil {
		return err
	}

	_ = a.sitemap.AddLyricsEntry(SitemapUrl{
		PublicId: lyrics.PublicId,
		AddedAt:  lyrics.CreatedAt.Format(time.RFC3339),
	})

	err = a.repo.DeleteLyricsRequest(id)
	if err != nil {
		return err
	}

	err = a.mailer.SendLyricsApprovedEmail(lyrics, lyricsRequest.RequesterEmail)
	if err != nil {
		return err
	}

	return nil
}

func (a *Actions) RejectLyricsRequest(adminToken string, id uint, reason string) error {
	err := a.jwt.Validate(adminToken, JwtAdminToken)
	if err != nil {
		return err
	}

	lyricsRequest, err := a.repo.GetLyricsRequestById(id)
	if err != nil {
		return err
	}

	err = a.repo.DeleteLyricsRequest(id)
	if err != nil {
		return err
	}

	err = a.mailer.SendLyricsRejectedEmail(reason, lyricsRequest.RequesterEmail)
	if err != nil {
		return err
	}

	return nil
}
