package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/WBetEngine/Seosementara/Backend/internal/pixel/service"
	"github.com/WBetEngine/Seosementara/Backend/internal/pixel/store"
	"github.com/google/uuid"
)

type CollectHandler struct {
	FB *service.FacebookService
}

type collectBody struct {
	Event        string         `json:"event"`
	EventID      string         `json:"event_id"`
	URL          string         `json:"url"`
	SiteKey      string         `json:"site_key"`
	DomainID     *int64         `json:"managed_domain_id"`
	FBP          string         `json:"fbp"`
	FBC          string         `json:"fbc"`
	FBCLID       string         `json:"fbclid"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	PhoneCountry string         `json:"phone_country"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	ExternalID   string         `json:"external_id"`
	Country      string         `json:"country"`
	LeadID       string         `json:"lead_id"`
	Props        map[string]any `json:"props"`
}

func (h *CollectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body collectBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	ip := clientIP(r)
	ua := r.UserAgent()
	if body.URL == "" {
		body.URL = r.Header.Get("Referer")
	}
	props := body.Props
	if props == nil {
		props = map[string]any{}
	}
	if body.LeadID != "" {
		props["lead_id"] = body.LeadID
	}

	in := store.CollectInput{
		Event:           body.Event,
		EventID:         body.EventID,
		URL:             body.URL,
		SiteKey:         body.SiteKey,
		ManagedDomainID: body.DomainID,
		FBP:             body.FBP,
		FBC:             body.FBC,
		FBCLID:          body.FBCLID,
		ClientIP:        ip,
		UserAgent:       ua,
		Email:           body.Email,
		Phone:           body.Phone,
		PhoneCountry:    body.PhoneCountry,
		FirstName:       body.FirstName,
		LastName:        body.LastName,
		ExternalID:      body.ExternalID,
		Country:         body.Country,
		Props:           props,
	}
	if in.EventID == "" {
		in.EventID = uuid.New().String()
	}
	id, err := h.FB.IngestCollect(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "queued_id": id, "event_id": in.EventID})
}

func clientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		if i := strings.Index(ip, ","); i > 0 {
			ip = strings.TrimSpace(ip[:i])
		}
		return ip
	}
	host := r.RemoteAddr
	if i := strings.LastIndex(host, ":"); i > 0 {
		return host[:i]
	}
	return host
}
