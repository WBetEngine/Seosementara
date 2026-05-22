package handler

import (
	"net/http"
	"os"
	"path/filepath"
)

func ServeTrackScript(staticDir string) http.HandlerFunc {
	path := filepath.Join(staticDir, "js", "sseo-track.js")
	data, err := os.ReadFile(path)
	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			http.Error(w, "script not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		_, _ = w.Write(data)
	}
}
