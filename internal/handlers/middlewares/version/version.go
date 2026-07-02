package version

import (
	"context"
	"net/http"

	"codeberg.org/dankstuff/danklyrics/pkg/version"
)

const VersionKey = "version"

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), VersionKey, version.Version)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
