package store

import "time"

const (
	AuthTypeAPIToken      = "api_token"
	AuthTypeGlobalAPIKey  = "global_api_key"
	ConnectorHealthy      = "healthy"
	ConnectorDown         = "down"
	ConnectorUnknown      = "unknown"
	ProjectTypeAdmin      = "admin"
)

type Credentials struct {
	ID               int16
	AuthType         string
	AccountID        string
	AccountEmail     string
	LastValidatedAt  *time.Time
	ValidationError  string
	UpdatedAt        time.Time
	TokenMasked      string
}

type CredentialsInput struct {
	AuthType     string `json:"auth_type"`
	APIToken     string `json:"api_token"`
	GlobalAPIKey string `json:"global_api_key"`
	AccountEmail string `json:"account_email"`
	AccountID    string `json:"account_id"`
}

type EnvEntry struct {
	Key       string
	Value     string
	IsPublic  bool
	UpdatedAt time.Time
}

type TunnelConfig struct {
	ID               int16
	TunnelID         string
	TunnelName       string
	InstallCommand   string
	IsActive         bool
	ConnectorStatus  string
	LastSeenAt       *time.Time
	OriginURL        string
	UpdatedAt        time.Time
}

type TunnelRoute struct {
	ID         int64
	Hostname   string
	PathPrefix string
	ServiceURL string
	IsEnabled  bool
}

type PagesProject struct {
	ID               int64
	ProjectType      string
	ProjectID        string
	ProjectName      string
	ProductionBranch string
	RootDirectory    string
	BuildCommand     string
	DeployCommand    string
	PagesURL         string
	DeployStatus     string
	LastDeployAt     *time.Time
	EnvSyncedAt      *time.Time
	UpdatedAt        time.Time
}

type APILog struct {
	ID           int64
	Method       string
	Path         string
	StatusCode   int
	Success      bool
	ErrorMessage string
	DurationMs   int
	CreatedAt    time.Time
}
