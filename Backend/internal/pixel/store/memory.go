package store

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/WBetEngine/Seosementara/Backend/internal/crypto"
)

// MemoryStore is a dev/test store when PostgreSQL is unavailable.
type MemoryStore struct {
	mu       sync.RWMutex
	hub      HubSettings
	config   *FacebookConfig
	secret   string
	events   []PixelEvent
	eventSeq int64
	assigns  []DomainAssignment
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		hub: HubSettings{
			ID: 1, TrackingHostname: "pelacak.seosementara.org",
			DefaultMode: "server_first", CollectPath: "/collect", ScriptVersion: "1",
		},
	}
}

func (m *MemoryStore) GetHubSettings(ctx context.Context) (*HubSettings, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	h := m.hub
	return &h, nil
}

func (m *MemoryStore) UpdateHubSettings(ctx context.Context, hostname, mode string, consent bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hub.TrackingHostname = hostname
	m.hub.DefaultMode = mode
	m.hub.ConsentRequired = consent
	return nil
}

func (m *MemoryStore) GetFacebookConfig(ctx context.Context) (*FacebookConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.config == nil {
		return nil, nil
	}
	c := *m.config
	return &c, nil
}

func (m *MemoryStore) SaveFacebookSetup(ctx context.Context, in FacebookSetupInput, encKey []byte) (*FacebookConfig, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if in.CAPIAccessToken != "" {
		ct, nonce, err := crypto.EncryptAESGCM(encKey, []byte(in.CAPIAccessToken))
		if err != nil {
			return nil, err
		}
		_ = ct
		_ = nonce
		m.secret = in.CAPIAccessToken
	}
	now := time.Now()
	m.config = &FacebookConfig{
		ID: 1, Name: in.Name, Scope: in.Scope, IsActive: true,
		PixelID: in.PixelID, BusinessID: in.BusinessID,
		CAPIEnabled: in.CAPIEnabled, BrowserPixelEnabled: in.BrowserPixelEnabled,
		TestEventCode: in.TestEventCode, CredentialName: in.CredentialName,
		ValidationStatus: "unknown", UpdatedAt: now,
	}
	return m.config, nil
}

func (m *MemoryStore) GetCredentialSecret(ctx context.Context, credID int64, encKey []byte) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.secret, nil
}

func (m *MemoryStore) UpdateCredentialValidation(ctx context.Context, credID int64, status string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.config != nil {
		m.config.ValidationStatus = status
		now := time.Now()
		m.config.LastValidatedAt = &now
	}
	return nil
}

func (m *MemoryStore) EnqueueEvent(ctx context.Context, ev PixelEvent) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.eventSeq++
	ev.ID = m.eventSeq
	ev.Status = "pending"
	ev.CreatedAt = time.Now()
	m.events = append(m.events, ev)
	return ev.ID, nil
}

func (m *MemoryStore) ListPendingEvents(ctx context.Context, limit int) ([]PixelEvent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []PixelEvent
	for _, e := range m.events {
		if e.Status == "pending" {
			out = append(out, e)
			if len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}

func (m *MemoryStore) MarkEventSent(ctx context.Context, id int64, platformEventID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.events {
		if m.events[i].ID == id {
			m.events[i].Status = "sent"
			m.events[i].PlatformEventID = platformEventID
			return nil
		}
	}
	return nil
}

func (m *MemoryStore) MarkEventFailed(ctx context.Context, id int64, errMsg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.events {
		if m.events[i].ID == id {
			m.events[i].Status = "failed"
			m.events[i].ErrorMessage = errMsg
			return nil
		}
	}
	return nil
}

func (m *MemoryStore) ListEvents(ctx context.Context, status string, limit, offset int) ([]PixelEvent, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var filtered []PixelEvent
	for i := len(m.events) - 1; i >= 0; i-- {
		e := m.events[i]
		if status == "" || e.Status == status {
			filtered = append(filtered, e)
		}
	}
	total := int64(len(filtered))
	if offset >= len(filtered) {
		return nil, total, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

func (m *MemoryStore) GetDiagnostics(ctx context.Context) (*Diagnostics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d := &Diagnostics{ConnectionState: "unknown"}
	if m.config != nil {
		d.ConnectionState = m.config.ValidationStatus
	}
	since := time.Now().Add(-24 * time.Hour)
	for _, e := range m.events {
		if e.CreatedAt.Before(since) {
			continue
		}
		d.Received24h++
		switch e.Status {
		case "sent":
			d.Sent24h++
		case "failed":
			d.Failed24h++
			if d.LastError == "" {
				d.LastError = e.ErrorMessage
			}
		case "pending":
			d.PendingCount++
		}
	}
	if d.Received24h > 0 {
		d.FailureRatePct = float64(d.Failed24h) / float64(d.Received24h) * 100
	}
	return d, nil
}

func (m *MemoryStore) IncrementDailyStat(ctx context.Context, configID int64, field string) error {
	return nil
}

func (m *MemoryStore) ListDomainAssignments(ctx context.Context, configID int64) ([]DomainAssignment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]DomainAssignment(nil), m.assigns...), nil
}

func (m *MemoryStore) AssignDomain(ctx context.Context, configID, domainID int64, hostname string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.assigns = append(m.assigns, DomainAssignment{
		ID: int64(len(m.assigns) + 1), PixelConfigID: configID,
		ManagedDomainID: domainID, DomainHostname: hostname, IsActive: true,
	})
	return nil
}

func (m *MemoryStore) UnassignDomain(ctx context.Context, configID, domainID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var next []DomainAssignment
	for _, a := range m.assigns {
		if a.PixelConfigID == configID && a.ManagedDomainID == domainID {
			continue
		}
		next = append(next, a)
	}
	m.assigns = next
	return nil
}

func (m *MemoryStore) SaveFacebookSetupPayload(secret string) {
	m.secret = secret
}

// helper for tests
func (m *MemoryStore) RawEvents() []PixelEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]PixelEvent(nil), m.events...)
}

var _ PixelStore = (*MemoryStore)(nil)

// Ensure JSON payloads work in memory enqueue
func marshalPayload(v map[string]any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
