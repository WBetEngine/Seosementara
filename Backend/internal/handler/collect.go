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
	Event   string         `json:"event"`
	EventID string         `json:"event_id"`
	URL     string         `json:"url"`
	SiteKey string         `json:"site_key"`
	DomainID *int64        `json:"managed_domain_id"`
	FBP     string         `json:"fbp"`
	FBC     string         `json:"fbc"`
	Email   string         `json:"email"`
	Phone   string         `json:"phone"`
	Props   map[string]any `json:"props"`
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
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	props := body.Props
	if props == nil {
		props = map[string]any{}
	}
	props["url"] = body.URL
	props["client_ip"] = ip
	props["user_agent"] = r.UserAgent()
	props["fbp"] = body.FBP
	props["fbc"] = body.FBC
	if body.Email != "" {
		props["email_hash"] = body.Email
	}
	if body.Phone != "" {
		props["phone_hash"] = body.Phone
	}

	in := store.CollectInput{
		Event: body.Event, EventID: body.EventID, URL: body.URL,
		SiteKey: body.SiteKey, ManagedDomainID: body.DomainID,
		FBP: body.FBP, FBC: body.FBC, ClientIP: ip, UserAgent: r.UserAgent(),
		Props: props,
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
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "queued_id": id})
}
