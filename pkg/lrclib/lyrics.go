package lrclib

import (
	"regexp"
	"strings"

	"codeberg.org/dankstuff/danklyrics/pkg/models"
)

var (
	sectionsPattern = regexp.MustCompile(`(\[.*\])`)
	timesPattern    = regexp.MustCompile(`(\[.*\]) `)
)

func (s *Song) Lyrics() *models.Lyrics {
	if s.lyrics == nil {
		s.lyrics = newLyrics(s)
	}

	return s.lyrics
}

func newLyrics(song *Song) *models.Lyrics {
	lyricsParts := strings.Split(song.PlainLyrics, "\n")
	fixedLyrics := make([]string, 0, len(lyricsParts))
	for _, part := range lyricsParts {
		if part == "" {
			continue
		}
		fixedLyrics = append(fixedLyrics, part)
	}

	var syncedParts map[string]string

	if song.SyncedLyrics != "" {
		syncedParts = make(map[string]string)
		syncedLyricsParts := strings.Split(song.SyncedLyrics, "\n")
		fixedSyncedLyrics := make([]string, 0, len(syncedLyricsParts))
		for _, part := range syncedLyricsParts {
			if part == "" {
				continue
			}
			fixedSyncedLyrics = append(fixedSyncedLyrics, part)
		}
		for _, part := range fixedSyncedLyrics {
			matches := timesPattern.FindSubmatch([]byte(part))
			if len(matches) > 0 {
				timeStampMarker := string(matches[0])
				nonTimestampPart := strings.TrimSpace(part[len(timeStampMarker):])
				justTime := strings.Trim(strings.TrimSpace(timeStampMarker), "[]")

				syncedParts[justTime] = nonTimestampPart
			}
		}
	}

	return &models.Lyrics{
		SongName:   song.Name,
		ArtistName: song.ArtistName,
		AlbumName:  song.AlbumName,
		Parts:      lyricsParts,
		Synced:     syncedParts,
	}
}
