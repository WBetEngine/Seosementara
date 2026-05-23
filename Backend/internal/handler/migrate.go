package handler

import (
	"net/http"

	"github.com/WBetEngine/Seosementara/Backend/internal/migrate"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MigrateHandler struct {
	Pool *pgxpool.Pool
	Dir  string
}

func (h *MigrateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.Pool == nil {
		http.Error(w, "database not configured", http.StatusServiceUnavailable)
		return
	}
	dir := h.Dir
	if dir == "" {
		dir = "migrations"
	}
	if err := migrate.Up(r.Context(), h.Pool, dir); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, nil)
		return
	}
	writeJSON(w, map[string]bool{"ok": true}, nil)
}
