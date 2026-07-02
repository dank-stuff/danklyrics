package contenttype

import "net/http"

func Html(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}

// IsNoLayoutPage checks if the requested page requires a no reload or not.
func IsNoLayoutPage(r *http.Request) bool {
	noReload, exists := r.URL.Query()["no_layout"]
	return exists && noReload[0] == "true"
}
