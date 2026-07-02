package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type LyricsRequest struct {
	Id uint `gorm:"primaryKey;autoIncrement"`

	SongTitle      string `gorm:"index"`
	ArtistName     string
	AlbumTitle     string
	RequesterEmail string

	LyricsPlain  []string          `gorm:"-"`
	LyricsSynced map[string]string `gorm:"-"`

	CreatedAt time.Time `gorm:"index"`
}

func (l *LyricsRequest) AfterCreate(tx *gorm.DB) error {
	errs := make([]error, 0, len(l.LyricsPlain)+len(l.LyricsSynced))
	for _, part := range l.LyricsPlain {
		lp := &LyricsRequestPart{
			LyricsRequestId: l.Id,
			Text:            part,
		}

		err := tx.Model(new(LyricsRequestPart)).Create(lp).Error
		errs = append(errs, err)
	}

	for tm, part := range l.LyricsSynced {
		lp := &LyricsRequestSyncedPart{
			LyricsRequestId: l.Id,
			Text:            part,
			Time:            tm,
		}

		err := tx.Model(new(LyricsRequestSyncedPart)).Create(lp).Error
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

type LyricsRequestPart struct {
	LyricsRequestId uint
	LyricsRequest   LyricsRequest
	Text            string
}

type LyricsRequestSyncedPart struct {
	LyricsRequestId uint
	LyricsRequest   LyricsRequest
	Time            string
	Text            string
}
