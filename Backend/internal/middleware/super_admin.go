package middleware

import (
	"net/http"
	"os"
	"strings"
)

// RequireSuperAdmin protects Cloudflare setup APIs.
// Set SUPER_ADMIN_TOKEN env; client sends Authorization: Bearer <token>
// or X-Super-Admin-Token header. HTMX: hx-headers in admin index.
func RequireSuperAdmin(next http.Handler) http.Handler {
	token := os.Getenv("SUPER_ADMIN_TOKEN")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token == "" {
			if os.Getenv("ALLOW_SETUP_WITHOUT_AUTH") == "true" {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "SUPER_ADMIN_TOKEN not configured", http.StatusServiceUnavailable)
			return
		}
		got := ""
		if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
			got = strings.TrimPrefix(h, "Bearer ")
		}
		if got == "" {
			got = r.Header.Get("X-Super-Admin-Token")
		}
		if got != token {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
