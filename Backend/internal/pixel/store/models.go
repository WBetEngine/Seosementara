package store

import (
	"encoding/json"
	"time"
)

type HubSettings struct {
	ID               int64
	TrackingHostname string
	DefaultMode      string
	ConsentRequired  bool
	CollectPath      string
	ScriptVersion    string
}

type FacebookConfig struct {
	ID                  int64
	Name                string
	Scope               string
	ManagedDomainID     *int64
	IsActive            bool
	ModeOverride        *string
	PixelID             string
	BusinessID          string
	CAPIEnabled         bool
	BrowserPixelEnabled bool
	TestEventCode       string
	CredentialsID       *int64
	CredentialName      string
	ValidationStatus    string
	LastValidatedAt     *time.Time
	UpdatedAt           time.Time
}

type FacebookSetupInput struct {
	Name                string `json:"name"`
	PixelID             string `json:"pixel_id"`
	BusinessID          string `json:"business_id"`
	CAPIAccessToken     string `json:"capi_access_token"`
	CredentialName      string `json:"credential_name"`
	TestEventCode       string `json:"test_event_code"`
	CAPIEnabled         bool   `json:"capi_enabled"`
	BrowserPixelEnabled bool   `json:"browser_pixel_enabled"`
	Scope               string `json:"scope"`
}

type CollectInput struct {
	Event           string
	EventID         string
	URL             string
	SiteKey         string
	ManagedDomainID *int64
	FBP             string
	FBC             string
	FBCLID          string
	ClientIP        string
	UserAgent       string
	Email           string
	Phone           string
	PhoneCountry    string // default 62 (ID)
	FirstName       string
	LastName        string
	ExternalID      string
	Country         string
	Props           map[string]any
}

type PixelEvent struct {
	ID              int64
	EventName       string
	EventID         string
	Status          string
	PixelConfigID   *int64
	ManagedDomainID *int64
	ErrorMessage    string
	PlatformEventID string
	CreatedAt       time.Time
	Payload         json.RawMessage
}

type Diagnostics struct {
	PendingCount    int64
	Failed24h       int64
	Sent24h         int64
	Received24h     int64
	FailureRatePct  float64
	LastError       string
	ConnectionState string
}

type DomainAssignment struct {
	ID              int64
	PixelConfigID   int64
	ManagedDomainID int64
	DomainHostname  string
	IsActive        bool
	DeployedAt      *time.Time
}
