package mariadb

import (
	"fmt"

	"codeberg.org/dankstuff/danklyrics/internal/models"
)

func Migrate() error {
	dbConn, err := dbConnector()
	if err != nil {
		return err
	}

	err = dbConn.Debug().AutoMigrate(
		new(models.Lyrics),
		new(models.LyricsPart),
		new(models.LyricsSyncedPart),
		new(models.LyricsRequest),
		new(models.LyricsRequestPart),
		new(models.LyricsRequestSyncedPart),
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

	for _, columnName := range []string{
		"song_title", "artist_name", "album_title",
	} {
		err = dbConn.Exec(fmt.Sprintf("ALTER TABLE lyrics CHANGE COLUMN `%s` `%s` TEXT CHARACTER SET 'utf8' COLLATE 'utf8_general_ci';", columnName, columnName)).Error
		if err != nil {
			return err
		}
	}

	return nil
}
