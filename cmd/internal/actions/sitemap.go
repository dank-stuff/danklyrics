package actions

import (
	"log"
	"time"
)

type Sitemap interface {
	GetLyricsEntries() ([]SitemapUrl, error)
	StoreLyricsesEntries(entries []SitemapUrl) error
	AddLyricsEntry(entry SitemapUrl) error
}

type SitemapUrl struct {
	PublicId string `json:"public_id"`
	AddedAt  string `json:"added_at"`
}

func (a *Actions) LoadLyricsPublicIds() error {
	log.Println("loading lyrics public ids...")
	lyrices, err := a.repo.GetLyricses(0)
	if err != nil {
		return err
	}

	entries := make([]SitemapUrl, 0, len(lyrices))
	for _, lyrics := range lyrices {
		entries = append(entries, SitemapUrl{
			PublicId: lyrics.PublicId,
			AddedAt:  lyrics.CreatedAt.Format(time.RFC3339),
		})
	}

	err = a.sitemap.StoreLyricsesEntries(entries)
	if err != nil {
		return err
	}

	log.Println("done loading lyrics public ids!")
	return nil
}

func (a *Actions) GetSitemap() ([]SitemapUrl, error) {
	sitemapEntries, err := a.sitemap.GetLyricsEntries()
	if err != nil {
		return nil, err
	}

	return sitemapEntries, nil
}
