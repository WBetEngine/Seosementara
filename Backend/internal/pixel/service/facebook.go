package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/WBetEngine/Seosementara/Backend/internal/pixel/facebook"
	"github.com/WBetEngine/Seosementara/Backend/internal/pixel/store"
	"github.com/google/uuid"
)

type FacebookService struct {
	Store  store.PixelStore
	CAPI   *facebook.CAPIClient
	EncKey []byte
}

func NewFacebookService(st store.PixelStore, capi *facebook.CAPIClient, encKey []byte) *FacebookService {
	return &FacebookService{Store: st, CAPI: capi, EncKey: encKey}
}

func (s *FacebookService) GetSetup(ctx context.Context) (*store.FacebookConfig, error) {
	return s.Store.GetFacebookConfig(ctx)
}

func (s *FacebookService) SaveSetup(ctx context.Context, in store.FacebookSetupInput) (*store.FacebookConfig, error) {
	if in.PixelID == "" {
		return nil, fmt.Errorf("pixel_id wajib")
	}
	if in.Scope == "" {
		in.Scope = "global"
	}
	return s.Store.SaveFacebookSetup(ctx, in, s.EncKey)
}

func (s *FacebookService) TestConnection(ctx context.Context) (map[string]any, error) {
	cfg, err := s.Store.GetFacebookConfig(ctx)
	if err != nil || cfg == nil || cfg.CredentialsID == nil {
		return nil, fmt.Errorf("setup Facebook belum lengkap atau token kosong")
	}
	token, err := s.Store.GetCredentialSecret(ctx, *cfg.CredentialsID, s.EncKey)
	if err != nil {
		return nil, err
	}
	res, err := s.sendTestEvent(ctx, cfg, token, "PageView")
	if err != nil {
		_ = s.Store.UpdateCredentialValidation(ctx, *cfg.CredentialsID, "error")
		return nil, err
	}
	_ = s.Store.UpdateCredentialValidation(ctx, *cfg.CredentialsID, "connected")
	return map[string]any{
		"events_received": res.EventsReceived,
		"messages":        res.Messages,
		"fbtrace_id":      res.FBTraceID,
	}, nil
}

func (s *FacebookService) SendTestEvent(ctx context.Context, eventName string) (map[string]any, error) {
	cfg, err := s.Store.GetFacebookConfig(ctx)
	if err != nil || cfg == nil || cfg.CredentialsID == nil {
		return nil, fmt.Errorf("setup belum ada")
	}
	token, err := s.Store.GetCredentialSecret(ctx, *cfg.CredentialsID, s.EncKey)
	if err != nil {
		return nil, err
	}
	if eventName == "" {
		eventName = "PageView"
	}
	res, err := s.sendTestEvent(ctx, cfg, token, eventName)
	if err != nil {
		return nil, err
	}
	return map[string]any{"ok": true, "events_received": res.EventsReceived}, nil
}

func (s *FacebookService) sendTestEvent(ctx context.Context, cfg *store.FacebookConfig, token, eventName string) (*facebook.EventsResponse, error) {
	props := map[string]any{
		"url":        "https://seosementara.org/admin/pixel/facebook/",
		"client_ip":  "127.0.0.1",
		"user_agent": "Seosementara-PixelHub/1.0 (CAPI test)",
		"fbp":        facebook.EnsureFBP(""),
	}
	facebook.EnrichCollectProps(props, "", "", "", "", "test-connection", "", "", facebook.DefaultPhoneCountry)

	ev, err := facebook.BuildServerEvent(eventName, uuid.New().String(), time.Now().Unix(), props)
	if err != nil {
		return nil, err
	}
	payload := facebook.EventsPayload{Data: []facebook.ServerEvent{ev}}
	if cfg.TestEventCode != "" {
		payload.TestEventCode = cfg.TestEventCode
	}
	return s.CAPI.SendEvents(ctx, cfg.PixelID, token, payload)
}

func (s *FacebookService) IngestCollect(ctx context.Context, in store.CollectInput) (int64, error) {
	cfg, _ := s.Store.GetFacebookConfig(ctx)
	var cfgID *int64
	if cfg != nil {
		cfgID = &cfg.ID
	}
	name := mapCanonicalToFacebook(in.Event)
	if in.Props == nil {
		in.Props = map[string]any{}
	}
	in.Props["url"] = in.URL
	in.Props["client_ip"] = in.ClientIP
	in.Props["user_agent"] = in.UserAgent
	in.Props["fbp"] = in.FBP
	in.Props["fbc"] = in.FBC
	if in.FBCLID != "" {
		in.Props["fbclid"] = in.FBCLID
	}
	facebook.EnrichCollectProps(in.Props, in.Email, in.Phone, in.FirstName, in.LastName, in.ExternalID, in.Country, in.FBCLID, in.PhoneCountry)

	payload, _ := json.Marshal(in.Props)
	ev := store.PixelEvent{
		EventName:       name,
		EventID:         in.EventID,
		PixelConfigID:   cfgID,
		ManagedDomainID: in.ManagedDomainID,
		Payload:         payload,
	}
	if ev.EventID == "" {
		ev.EventID = uuid.New().String()
	}
	id, err := s.Store.EnqueueEvent(ctx, ev)
	if err != nil {
		return 0, err
	}
	if cfg != nil {
		_ = s.Store.IncrementDailyStat(ctx, cfg.ID, "received")
	}
	return id, nil
}

func mapCanonicalToFacebook(e string) string {
	switch strings.ToLower(e) {
	case "page_view", "pageview":
		return "PageView"
	case "lead":
		return "Lead"
	case "purchase":
		return "Purchase"
	case "view_content", "viewcontent":
		return "ViewContent"
	case "add_to_cart", "addtocart":
		return "AddToCart"
	case "initiate_checkout":
		return "InitiateCheckout"
	case "complete_registration":
		return "CompleteRegistration"
	case "search":
		return "Search"
	case "contact":
		return "Contact"
	case "click":
		return "ViewContent"
	default:
		if e == "" {
			return "PageView"
		}
		// PascalCase custom names as-is
		return e
	}
}

func (s *FacebookService) DispatchPending(ctx context.Context, batch int) (sent, failed int, err error) {
	cfg, err := s.Store.GetFacebookConfig(ctx)
	if err != nil || cfg == nil || !cfg.CAPIEnabled || cfg.CredentialsID == nil {
		return 0, 0, nil
	}
	token, err := s.Store.GetCredentialSecret(ctx, *cfg.CredentialsID, s.EncKey)
	if err != nil {
		return 0, 0, err
	}

	pending, err := s.Store.ListPendingEvents(ctx, batch)
	if err != nil {
		return 0, 0, err
	}

	for _, pe := range pending {
		var props map[string]any
		_ = json.Unmarshal(pe.Payload, &props)

		ev, buildErr := facebook.BuildServerEvent(pe.EventName, pe.EventID, pe.CreatedAt.Unix(), props)
		if buildErr != nil {
			_ = s.Store.MarkEventFailed(ctx, pe.ID, buildErr.Error())
			_ = s.Store.IncrementDailyStat(ctx, cfg.ID, "failed")
			failed++
			continue
		}

		payload := facebook.EventsPayload{Data: []facebook.ServerEvent{ev}}
		if cfg.TestEventCode != "" {
			payload.TestEventCode = cfg.TestEventCode
		}
		res, sendErr := s.CAPI.SendEvents(ctx, cfg.PixelID, token, payload)
		if sendErr != nil {
			_ = s.Store.MarkEventFailed(ctx, pe.ID, sendErr.Error())
			_ = s.Store.IncrementDailyStat(ctx, cfg.ID, "failed")
			failed++
			continue
		}
		_ = s.Store.MarkEventSent(ctx, pe.ID, res.FBTraceID)
		_ = s.Store.IncrementDailyStat(ctx, cfg.ID, "sent")
		sent++
	}
	return sent, failed, nil
}

func (s *FacebookService) Diagnostics(ctx context.Context) (*store.Diagnostics, error) {
	return s.Store.GetDiagnostics(ctx)
}

func (s *FacebookService) ListEvents(ctx context.Context, status string, page, limit int) ([]store.PixelEvent, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	return s.Store.ListEvents(ctx, status, limit, offset)
}

func (s *FacebookService) ListDomains(ctx context.Context) ([]store.DomainAssignment, error) {
	cfg, err := s.Store.GetFacebookConfig(ctx)
	if err != nil || cfg == nil {
		return nil, err
	}
	return s.Store.ListDomainAssignments(ctx, cfg.ID)
}

func (s *FacebookService) AssignDomain(ctx context.Context, domainID int64, hostname string) error {
	cfg, err := s.Store.GetFacebookConfig(ctx)
	if err != nil || cfg == nil {
		return fmt.Errorf("simpan setup Facebook terlebih dahulu")
	}
	return s.Store.AssignDomain(ctx, cfg.ID, domainID, hostname)
}
