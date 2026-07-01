package mariadb

import (
	"codeberg.org/dankstuff/danklyrics/internal/models"
)

func Migrate() error {
	dbConn, err := dbConnector()
	if err != nil {
		return err
	}

	err = dbConn.Debug().AutoMigrate(
		new(models.LyricsPart),
		new(models.LyricsSyncedPart),
		new(models.Lyrics),
		new(models.LyricsRequestPart),
		new(models.LyricsRequestSyncedPart),
		new(models.LyricsRequest),
		new(models.Admin),
	)
	if err != nil {
		return err
	}

	for _, tableName := range []string{
		"lyrics", "lyrics_parts", "lyrics_synced_parts",
		"lyrics_requests", "lyrics_request_parts", "lyrics_request_synced_parts",
	} {
		err = dbConn.Exec("ALTER TABLE " + tableName + " CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error
		if err != nil {
			return err
		}
	}

	return nil
}
