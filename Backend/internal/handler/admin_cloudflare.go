package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/WBetEngine/Seosementara/Backend/internal/setup/service"
	"github.com/WBetEngine/Seosementara/Backend/internal/setup/store"
	"github.com/go-chi/chi/v5"
)

type AdminCloudflare struct {
	Svc          *service.SetupService
	PartialsDir  string
}

func (h *AdminCloudflare) RoutesHTML(r chi.Router) {
	r.Get("/", h.PageShell)
	r.Get("/koneksi", h.PageKoneksi)
	r.Get("/domain", h.PageDomain)
	r.Get("/tunnel", h.PageTunnel)
	r.Get("/pages", h.PagePages)
	r.Get("/dns", h.PageDNS)
}

func (h *AdminCloudflare) RoutesAPI(r chi.Router) {
	r.Get("/credentials", h.APIGetCredentials)
	r.Put("/credentials", h.APIPutCredentials)
	r.Post("/credentials/test", h.APITestCredentials)

	r.Get("/domain-env", h.APIGetDomainEnv)
	r.Put("/domain-env", h.APIPutDomainEnv)
	r.Post("/domain-env/sync-pages", h.APISyncPagesEnv)

	r.Get("/tunnel", h.APIGetTunnel)
	r.Post("/tunnel", h.APICreateTunnel)
	r.Post("/tunnel/routes/apply", h.APIApplyRoutes)
	r.Get("/tunnel/status", h.APITunnelStatus)
	r.Put("/tunnel/route", h.APIPutRoute)

	r.Get("/pages", h.APIGetPages)
	r.Put("/pages", h.APIPutPages)
	r.Post("/pages/deploy", h.APIPagesDeploy)

	r.Post("/dns/apply", h.APIDNSApply)
	r.Get("/logs", h.APIGetLogs)
}

func (h *AdminCloudflare) PageShell(w http.ResponseWriter, r *http.Request) {
	h.servePartial(w, "settings-cloudflare.html", nil)
}

func (h *AdminCloudflare) PageKoneksi(w http.ResponseWriter, r *http.Request) {
	v, _ := h.Svc.GetCredentialsView(r.Context())
	h.servePartial(w, "settings-cf-koneksi.html", map[string]any{"View": v})
}

func (h *AdminCloudflare) PageDomain(w http.ResponseWriter, r *http.Request) {
	env, _ := h.Svc.ListEnv(r.Context())
	h.servePartial(w, "settings-cf-domain.html", map[string]any{"Env": env})
}

func (h *AdminCloudflare) PageTunnel(w http.ResponseWriter, r *http.Request) {
	tv, _ := h.Svc.GetTunnelView(r.Context())
	h.servePartial(w, "settings-cf-tunnel.html", map[string]any{"Tunnel": tv})
}

func (h *AdminCloudflare) PagePages(w http.ResponseWriter, r *http.Request) {
	pp, _ := h.Svc.GetPagesAdmin(r.Context())
	h.servePartial(w, "settings-cf-pages.html", map[string]any{"Pages": pp})
}

func (h *AdminCloudflare) PageDNS(w http.ResponseWriter, r *http.Request) {
	h.servePartial(w, "settings-cf-dns.html", nil)
}

func (h *AdminCloudflare) servePartial(w http.ResponseWriter, file string, data any) {
	path := filepath.Join(h.PartialsDir, file)
	if _, err := os.Stat(path); err != nil {
		http.Error(w, "partial not found: "+file, http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = t.Execute(w, data)
}

func (h *AdminCloudflare) APIGetCredentials(w http.ResponseWriter, r *http.Request) {
	v, err := h.Svc.GetCredentialsView(r.Context())
	writeJSON(w, v, err)
}

func (h *AdminCloudflare) APIPutCredentials(w http.ResponseWriter, r *http.Request) {
	var in store.CredentialsInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	test := r.URL.Query().Get("test") == "1" || r.URL.Query().Get("test") == "true"
	v, err := h.Svc.SaveCredentials(r.Context(), in, test)
	writeJSON(w, v, err)
}

func (h *AdminCloudflare) APITestCredentials(w http.ResponseWriter, r *http.Request) {
	err := h.Svc.TestCredentials(r.Context())
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, nil)
		return
	}
	v, _ := h.Svc.GetCredentialsView(r.Context())
	writeJSON(w, map[string]any{"ok": true, "credentials": v}, nil)
}

func (h *AdminCloudflare) APIGetDomainEnv(w http.ResponseWriter, r *http.Request) {
	env, err := h.Svc.ListEnv(r.Context())
	writeJSON(w, env, err)
}

func (h *AdminCloudflare) APIPutDomainEnv(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, nil, h.Svc.SaveEnvBulk(r.Context(), body))
}

func (h *AdminCloudflare) APISyncPagesEnv(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]bool{"ok": true}, h.Svc.SyncPagesEnv(r.Context()))
}

func (h *AdminCloudflare) APIGetTunnel(w http.ResponseWriter, r *http.Request) {
	tv, err := h.Svc.GetTunnelView(r.Context())
	writeJSON(w, tv, err)
}

func (h *AdminCloudflare) APICreateTunnel(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name      string `json:"name"`
		OriginURL string `json:"origin_url"`
	}
	_ = json.NewDecoder(r.Body).Decode(&in)
	tv, err := h.Svc.CreateTunnel(r.Context(), in.Name, in.OriginURL)
	writeJSON(w, tv, err)
}

func (h *AdminCloudflare) APIApplyRoutes(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]bool{"ok": true}, h.Svc.ApplyTunnelRoutes(r.Context()))
}

func (h *AdminCloudflare) APITunnelStatus(w http.ResponseWriter, r *http.Request) {
	err := h.Svc.RefreshTunnelStatus(r.Context())
	tv, _ := h.Svc.GetTunnelView(r.Context())
	writeJSON(w, map[string]any{"ok": err == nil, "error": errString(err), "tunnel": tv}, nil)
}

func (h *AdminCloudflare) APIPutRoute(w http.ResponseWriter, r *http.Request) {
	var rte store.TunnelRoute
	if err := json.NewDecoder(r.Body).Decode(&rte); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]bool{"ok": true}, h.Svc.SaveTunnelRoute(r.Context(), rte))
}

func (h *AdminCloudflare) APIGetPages(w http.ResponseWriter, r *http.Request) {
	pp, err := h.Svc.GetPagesAdmin(r.Context())
	writeJSON(w, pp, err)
}

func (h *AdminCloudflare) APIPutPages(w http.ResponseWriter, r *http.Request) {
	var pp store.PagesProject
	if err := json.NewDecoder(r.Body).Decode(&pp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]bool{"ok": true}, h.Svc.SavePagesAdmin(r.Context(), pp))
}

func (h *AdminCloudflare) APIPagesDeploy(w http.ResponseWriter, r *http.Request) {
	err := h.Svc.TriggerPagesDeploy(r.Context())
	writeJSON(w, map[string]any{"ok": err == nil, "error": errString(err)}, nil)
}

func (h *AdminCloudflare) APIDNSApply(w http.ResponseWriter, r *http.Request) {
	res, err := h.Svc.ApplyDNS(r.Context())
	writeJSON(w, res, err)
}

func (h *AdminCloudflare) APIGetLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.Svc.ListLogs(r.Context(), 50)
	writeJSON(w, logs, err)
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// MountCloudflare registers HTMX HTML (settings) + JSON API (setup alias).
func MountCloudflare(r chi.Router, h *AdminCloudflare, guard func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(guard)
		r.Route("/api/admin/settings/cloudflare", h.RoutesHTML)
		r.Route("/api/admin/setup/cloudflare", h.RoutesAPI)
	})
}
