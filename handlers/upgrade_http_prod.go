//go:build !dev

package handlers

import (
	"net/http"
	"os"
)

func upgradeToHttps(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// When behind a proxy (PS_BEHIND_PROXY is set), trust that the proxy
		// handles TLS termination. The proxy sets X-Forwarded-Proto, so we
		// shouldn't try to upgrade the scheme ourselves.
		if os.Getenv("PS_BEHIND_PROXY") != "" {
			h.ServeHTTP(w, r)
			return
		}
		// If client is connecting over plaintext HTTP, upgrade to HTTPS.
		if r.Header.Get("X-Forwarded-Proto") == "http" {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			return
		}
		h.ServeHTTP(w, r)
	})
}
