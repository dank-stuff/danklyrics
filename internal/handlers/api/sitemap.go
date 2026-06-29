package api

import (
	"encoding/json"
	"net/http"

	"codeberg.org/dankstuff/danklyrics/internal/actions"
)

type sitemapApi struct {
	usecases *actions.Actions
}

func NewSitemapApi(usecases *actions.Actions) *sitemapApi {
	return &sitemapApi{
		usecases: usecases,
	}
}

func (s *sitemapApi) HandleGetSitemapEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sitemapEntries, err := s.usecases.GetSitemap()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "Something went wrong!",
		})
		return
	}

	_ = json.NewEncoder(w).Encode(sitemapEntries)
}
