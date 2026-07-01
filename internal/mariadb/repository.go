package mariadb

import (
	"fmt"
	"strings"

	"codeberg.org/dankstuff/danklyrics/internal/actions"
	"codeberg.org/dankstuff/danklyrics/internal/models"

	"gorm.io/gorm"
)

type repository struct {
	client *gorm.DB
}

func New() (*repository, error) {
	conn, err := dbConnector()
	if err != nil {
		return nil, err
	}

	return &repository{
		client: conn,
	}, nil
}

func (r *repository) CreateLyrics(lyrics models.Lyrics) (models.Lyrics, error) {
	lyrics.PublicId = actions.Slugify(fmt.Sprintf("%s-%s", lyrics.ArtistName, lyrics.SongTitle))

	err := tryWrapDbError(
		r.client.
			Model(new(models.Lyrics)).
			Create(&lyrics).
			Error,
	)
	if err != nil {
		return models.Lyrics{}, err
	}

	return lyrics, nil
}

func (r *repository) getLyricsParts(lyricsId uint) (plain []string, synced map[string]string, err error) {
	parts := make([]models.LyricsPart, 0)
	err = r.client.
		Model(new(models.LyricsPart)).
		Where("lyrics_id = ?", lyricsId).
		Find(&parts).
		Error
	if err != nil {
		return nil, nil, err
	}

	plain = make([]string, 0, len(parts))
	for _, part := range parts {
		plain = append(plain, part.Text)
	}

	syncedParts := make([]models.LyricsSyncedPart, 0)
	_ = r.client.
		Model(new(models.LyricsSyncedPart)).
		Where("lyrics_id = ?", lyricsId).
		Find(&syncedParts).
		Error

	synced = make(map[string]string, 0)
	for _, part := range syncedParts {
		synced[part.Time] = part.Text
	}

	return
}

func (r *repository) GetLyricsByPublicId(id string) (models.Lyrics, error) {
	var lyrics models.Lyrics

	err := tryWrapDbError(
		r.client.
			Model(new(models.Lyrics)).
			First(&lyrics, "public_id = ?", id).
			Error,
	)
	if err != nil {
		return models.Lyrics{}, err
	}

	parts, synced, err := r.getLyricsParts(lyrics.Id)
	if err != nil {
		return models.Lyrics{}, err
	}
	lyrics.LyricsPlain = parts
	lyrics.LyricsSynced = synced

	return lyrics, nil
}

func (r *repository) FindLyricsExact(search actions.FindLyricsParams) ([]models.Lyrics, error) {
	lyricses := make([]models.Lyrics, 0)

	whereClause := "LOWER(song_title) LIKE LOWER(?) AND LOWER(artist_name) LIKE LOWER(?) AND LOWER(album_title) LIKE LOWER(?)"
	whereArgs := []any{
		likeArg(search.SongTitle),
		likeArg(search.ArtistName),
		likeArg(search.AlbumTitle),
	}

	err := tryWrapDbError(
		r.client.
			Model(new(models.Lyrics)).
			Preload("LyricsPlain").
			Preload("LyricsSynced").
			Where(whereClause, whereArgs...).
			Find(&lyricses).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return lyricses, nil
}

var innodbStopwords = map[string]struct{}{
	"a": struct{}{}, "about": struct{}{}, "an": struct{}{}, "are": struct{}{}, "as": struct{}{}, "at": struct{}{},
	"be": struct{}{}, "by": struct{}{}, "com": struct{}{}, "de": struct{}{}, "en": struct{}{}, "for": struct{}{},
	"from": struct{}{}, "how": struct{}{}, "i": struct{}{}, "in": struct{}{}, "is": struct{}{}, "it": struct{}{},
	"la": struct{}{}, "of": struct{}{}, "on": struct{}{}, "or": struct{}{}, "that": struct{}{}, "the": struct{}{},
	"this": struct{}{}, "to": struct{}{}, "was": struct{}{}, "what": struct{}{}, "when": struct{}{},
	"where": struct{}{}, "who": struct{}{}, "will": struct{}{}, "with": struct{}{}, "und": struct{}{}, "www": struct{}{},
}

func (r *repository) FindLyricsAll(search actions.FindLyricsParams) ([]models.Lyrics, error) {
	lyricses := make([]models.Lyrics, 0)

	searchWords := make([]string, 0)
	for _, word := range []string{
		search.SongTitle,
		search.AlbumTitle,
		search.ArtistName,
	} {
		word = strings.TrimSpace(word)
		if len(word) == 0 {
			continue
		}
		searchWords = append(searchWords, strings.Split(word, " ")...)
	}

	booleanQuery := new(strings.Builder)
	for _, word := range searchWords {
		word = strings.ToLower(word)

		if _, ok := innodbStopwords[word]; ok {
			continue
		}

		if len(word) < 3 {
			booleanQuery.WriteString(word + " ")
			continue
		}

		booleanQuery.WriteString("+")
		booleanQuery.WriteString(word)
		booleanQuery.WriteString("* ")
	}
	ftsBooleanQuery := strings.TrimSpace(booleanQuery.String())

	whereClause := "MATCH(song_title, artist_name, album_title) AGAINST(? IN BOOLEAN MODE)"
	whereArgs := []any{
		ftsBooleanQuery,
	}

	err := tryWrapDbError(
		r.client.
			Model(new(models.Lyrics)).
			Where(whereClause, whereArgs...).
			Find(&lyricses).
			Error,
	)
	if err != nil {
		return nil, err
	}

	for i := range lyricses {
		parts, synced, err := r.getLyricsParts(lyricses[i].Id)
		if err != nil {
			return nil, err
		}
		lyricses[i].LyricsPlain = parts
		lyricses[i].LyricsSynced = synced
	}

	return lyricses, nil
}

func (r *repository) GetLyricses(page int) ([]models.Lyrics, error) {
	lyricses := make([]models.Lyrics, 0)

	err := tryWrapDbError(
		r.client.
			Model(new(models.Lyrics)).
			Find(&lyricses).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return lyricses, nil
}

func (r *repository) CreateLyricsRequest(l models.LyricsRequest) (models.LyricsRequest, error) {
	err := tryWrapDbError(
		r.client.
			Model(new(models.LyricsRequest)).
			Create(&l).
			Error,
	)
	if err != nil {
		return models.LyricsRequest{}, err
	}

	return l, nil
}

func (r *repository) DeleteLyricsRequest(id uint) error {
	err := r.client.
		Exec("DELETE FROM lyrics_request_parts WHERE lyrics_request_id = ?", id).
		Error
	if err != nil {
		return err
	}

	_ = r.client.
		Exec("DELETE FROM lyrics_request_synced_parts WHERE lyrics_request_id = ?", id).
		Error

	return tryWrapDbError(
		r.client.
			Exec("DELETE FROM lyrics_requests WHERE id = ?", id).
			Error,
	)
}

func (r *repository) GetLyricsRequestById(id uint) (models.LyricsRequest, error) {
	var lyrics models.LyricsRequest

	err := tryWrapDbError(
		r.client.
			Model(new(models.LyricsRequest)).
			First(&lyrics, "id = ?", id).
			Error,
	)
	if err != nil {
		return models.LyricsRequest{}, err
	}

	parts := make([]models.LyricsRequestPart, 0)
	err = tryWrapDbError(
		r.client.
			Model(new(models.LyricsRequestPart)).
			Where("lyrics_request_id = ?", id).
			Find(&parts).
			Error,
	)
	if err != nil {
		return models.LyricsRequest{}, err
	}

	lyrics.LyricsPlain = make([]string, 0, len(parts))
	for _, part := range parts {
		lyrics.LyricsPlain = append(lyrics.LyricsPlain, part.Text)
	}

	synced := make([]models.LyricsRequestSyncedPart, 0)
	err = tryWrapDbError(
		r.client.
			Model(new(models.LyricsRequestSyncedPart)).
			Where("lyrics_request_id = ?", id).
			Find(&synced).
			Error,
	)
	if err != nil {
		return lyrics, nil
	}

	lyrics.LyricsSynced = make(map[string]string, 0)
	for _, part := range synced {
		lyrics.LyricsSynced[part.Time] = part.Text
	}

	return lyrics, nil
}

func (r *repository) GetLyricsRequests() ([]models.LyricsRequest, error) {
	lyricsRequests := make([]models.LyricsRequest, 0)

	err := tryWrapDbError(
		r.client.
			Model(new(models.LyricsRequest)).
			Find(&lyricsRequests).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return lyricsRequests, nil
}

func (r *repository) GetAdminByUsername(username string) (models.Admin, error) {
	var admin models.Admin

	err := tryWrapDbError(
		r.client.
			Model(new(models.Admin)).
			First(&admin, "username = ?", username).
			Error,
	)
	if err != nil {
		return models.Admin{}, err
	}

	return admin, nil
}

func likeArg(arg string) string {
	return fmt.Sprintf("%%%s%%", arg)
}
