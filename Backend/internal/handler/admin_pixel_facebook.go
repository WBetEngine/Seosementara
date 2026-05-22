package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/WBetEngine/Seosementara/Backend/internal/pixel/service"
	"github.com/WBetEngine/Seosementara/Backend/internal/pixel/store"
	"github.com/go-chi/chi/v5"
)

type AdminPixelFacebook struct {
	Svc       *service.FacebookService
	Templates *template.Template
}

func (h *AdminPixelFacebook) Routes(r chi.Router) {
	r.Get("/", h.PageOverview)
	r.Get("/setup", h.PageSetup)
	r.Get("/connection", h.PageConnection)
	r.Get("/domains", h.PageDomains)
	r.Get("/diagnostics", h.PageDiagnostics)
	r.Get("/events", h.PageEvents)
	r.Get("/analytics", h.PageAnalytics)

	r.Get("/api/setup", h.APIGetSetup)
	r.Post("/api/setup", h.APIPostSetup)
	r.Post("/api/test-connection", h.APITestConnection)
	r.Post("/api/test-event", h.APITestEvent)
	r.Get("/api/diagnostics", h.APIDiagnostics)
	r.Get("/api/events", h.APIEvents)
	r.Post("/api/domains/assign", h.APIAssignDomain)
}

func (h *AdminPixelFacebook) PageOverview(w http.ResponseWriter, r *http.Request) {
	h.render(w, "pixel/facebook/overview.html", h.pageData(r, "overview"))
}

func (h *AdminPixelFacebook) PageSetup(w http.ResponseWriter, r *http.Request) {
	cfg, _ := h.Svc.GetSetup(r.Context())
	h.render(w, "pixel/facebook/setup.html", map[string]any{
		"Tab": "setup", "Config": cfg,
		"EventsManagerURL": eventsManagerURL(cfg),
	})
}

func (h *AdminPixelFacebook) PageConnection(w http.ResponseWriter, r *http.Request) {
	cfg, _ := h.Svc.GetSetup(r.Context())
	diag, _ := h.Svc.Diagnostics(r.Context())
	h.render(w, "pixel/facebook/connection.html", map[string]any{
		"Tab": "connection", "Config": cfg, "Diag": diag,
		"EventsManagerURL": eventsManagerURL(cfg),
	})
}

func (h *AdminPixelFacebook) PageDomains(w http.ResponseWriter, r *http.Request) {
	domains, _ := h.Svc.ListDomains(r.Context())
	h.render(w, "pixel/facebook/domains.html", map[string]any{
		"Tab": "domains", "Assignments": domains,
	})
}

func (h *AdminPixelFacebook) PageDiagnostics(w http.ResponseWriter, r *http.Request) {
	diag, _ := h.Svc.Diagnostics(r.Context())
	h.render(w, "pixel/facebook/diagnostics.html", map[string]any{"Tab": "diagnostics", "Diag": diag})
}

func (h *AdminPixelFacebook) PageEvents(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	events, total, _ := h.Svc.ListEvents(r.Context(), status, 1, 50)
	h.render(w, "pixel/facebook/events.html", map[string]any{
		"Tab": "events", "Events": events, "Total": total, "StatusFilter": status,
	})
}

func (h *AdminPixelFacebook) PageAnalytics(w http.ResponseWriter, r *http.Request) {
	diag, _ := h.Svc.Diagnostics(r.Context())
	h.render(w, "pixel/facebook/analytics.html", map[string]any{"Tab": "analytics", "Diag": diag})
}

func (h *AdminPixelFacebook) APIGetSetup(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.Svc.GetSetup(r.Context())
	writeJSON(w, cfg, err)
}

func (h *AdminPixelFacebook) APIPostSetup(w http.ResponseWriter, r *http.Request) {
	var in store.FacebookSetupInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cfg, err := h.Svc.SaveSetup(r.Context(), in)
	writeJSON(w, cfg, err)
}

func (h *AdminPixelFacebook) APITestConnection(w http.ResponseWriter, r *http.Request) {
	res, err := h.Svc.TestConnection(r.Context())
	writeJSON(w, res, err)
}

func (h *AdminPixelFacebook) APITestEvent(w http.ResponseWriter, r *http.Request) {
	var body struct{ EventName string `json:"event_name"` }
	_ = json.NewDecoder(r.Body).Decode(&body)
	res, err := h.Svc.SendTestEvent(r.Context(), body.EventName)
	writeJSON(w, res, err)
}

func (h *AdminPixelFacebook) APIDiagnostics(w http.ResponseWriter, r *http.Request) {
	d, err := h.Svc.Diagnostics(r.Context())
	writeJSON(w, d, err)
}

func (h *AdminPixelFacebook) APIEvents(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	events, total, err := h.Svc.ListEvents(r.Context(), r.URL.Query().Get("status"), page, 50)
	writeJSON(w, map[string]any{"items": events, "total": total}, err)
}

func (h *AdminPixelFacebook) APIAssignDomain(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DomainID int64  `json:"managed_domain_id"`
		Hostname string `json:"hostname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.Svc.AssignDomain(r.Context(), body.DomainID, body.Hostname)
	writeJSON(w, map[string]any{"ok": err == nil}, err)
}

func (h *AdminPixelFacebook) pageData(r *http.Request, tab string) map[string]any {
	diag, _ := h.Svc.Diagnostics(r.Context())
	cfg, _ := h.Svc.GetSetup(r.Context())
	return map[string]any{"Tab": tab, "Diag": diag, "Config": cfg}
}

func (h *AdminPixelFacebook) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if h.Templates == nil {
		http.Error(w, "templates not loaded", http.StatusInternalServerError)
		return
	}
	if err := h.Templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func eventsManagerURL(cfg *store.FacebookConfig) string {
	if cfg == nil || cfg.PixelID == "" {
		return "https://business.facebook.com/events_manager2"
	}
	return "https://business.facebook.com/events_manager2/list/pixel/" + cfg.PixelID
}

func writeJSON(w http.ResponseWriter, v any, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(v)
}

func LoadTemplates(dir string) (*template.Template, error) {
	order := []string{
		filepath.Join(dir, "pixel/facebook/partials/*.html"),
		filepath.Join(dir, "pixel/facebook/*.html"),
	}
	var tpl *template.Template
	for _, pattern := range order {
		m, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		if len(m) == 0 {
			continue
		}
		if tpl == nil {
			tpl, err = template.ParseFiles(m...)
		} else {
			tpl, err = tpl.ParseFiles(m...)
		}
		if err != nil {
			return nil, err
		}
	}
	return tpl, nil
}
