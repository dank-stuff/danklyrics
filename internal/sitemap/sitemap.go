package sitemap

import (
	"sync"

	"codeberg.org/dankstuff/danklyrics/internal/actions"
)

type sitemap struct {
	mu      sync.Mutex
	entries []actions.SitemapUrl
}

func New() *sitemap {
	return &sitemap{
		mu:      sync.Mutex{},
		entries: make([]actions.SitemapUrl, 0),
	}
}

func (s *sitemap) GetLyricsEntries() ([]actions.SitemapUrl, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.entries, nil
}

func (s *sitemap) StoreLyricsesEntries(entries []actions.SitemapUrl) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = entries

	return nil
}

func (s *sitemap) AddLyricsEntry(entry actions.SitemapUrl) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, entry)

	return nil
}
