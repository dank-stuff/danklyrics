package actions

import "codeberg.org/dankstuff/danklyrics/cmd/internal/models"

type Mailer interface {
	SendVerificationEmail(token, email string) error
	SendLyricsApprovedEmail(lyrics models.Lyrics, email string) error
	SendLyricsRejectedEmail(reason, email string) error
}
