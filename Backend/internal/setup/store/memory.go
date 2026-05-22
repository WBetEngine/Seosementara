package store

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/WBetEngine/Seosementara/Backend/internal/crypto"
)

type MemoryStore struct {
	mu          sync.RWMutex
	cred        *struct {
		AuthType     string
		Secret       []byte
		AccountID    string
		AccountEmail string
		Validated    *time.Time
		Err          string
	}
	env         map[string]EnvEntry
	tunnel      TunnelConfig
	routes      []TunnelRoute
	pages       map[string]PagesProject
	logs        []APILog
	nextRouteID int64
}

func NewMemoryStore() *MemoryStore {
	m := &MemoryStore{
		env: map[string]EnvEntry{
			"PRIMARY_DOMAIN": {Key: "PRIMARY_DOMAIN", Value: "seosementara.org", IsPublic: true},
			"APEX_URL":       {Key: "APEX_URL", Value: "https://seosementara.org", IsPublic: true},
			"API_BASE_URL":   {Key: "API_BASE_URL", Value: "https://seosementara.org", IsPublic: true},
			"ADMIN_BASE_PATH": {Key: "ADMIN_BASE_PATH", Value: "/admin", IsPublic: true},
			"ENVIRONMENT":    {Key: "ENVIRONMENT", Value: "production", IsPublic: true},
		},
		tunnel: TunnelConfig{
			ID: 1, TunnelName: "seosementara-api", OriginURL: "http://127.0.0.1:8080",
			ConnectorStatus: ConnectorUnknown,
		},
		pages: map[string]PagesProject{
			ProjectTypeAdmin: {
				ProjectType: ProjectTypeAdmin, ProjectName: "seosementara",
				ProductionBranch: "main", RootDirectory: "Frontend-admin",
				BuildCommand: "npm ci", DeployCommand: "npx wrangler deploy",
			},
		},
	}
	m.routes = []TunnelRoute{
		{ID: 1, Hostname: "seosementara.org", PathPrefix: "/api", ServiceURL: m.tunnel.OriginURL, IsEnabled: true},
	}
	m.nextRouteID = 2
	return m
}

func (m *MemoryStore) GetCredentials(ctx context.Context) (*Credentials, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.cred == nil {
		return nil, nil
	}
	return &Credentials{
		AuthType: m.cred.AuthType, AccountID: m.cred.AccountID, AccountEmail: m.cred.AccountEmail,
		LastValidatedAt: m.cred.Validated, ValidationError: m.cred.Err, TokenMasked: "••••••••",
	}, nil
}

func (m *MemoryStore) GetCredentialsSecret(ctx context.Context, encKey []byte) (string, string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.cred == nil {
		return "", "", nil
	}
	return m.cred.AuthType, string(m.cred.Secret), nil
}

func (m *MemoryStore) SaveCredentials(ctx context.Context, authType string, secret []byte, accountID, accountEmail string, encKey []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	copySecret := make([]byte, len(secret))
	copy(copySecret, secret)
	m.cred = &struct {
		AuthType     string
		Secret       []byte
		AccountID    string
		AccountEmail string
		Validated    *time.Time
		Err          string
	}{authType, copySecret, accountID, accountEmail, nil, ""}
	_ = encKey
	return nil
}

func (m *MemoryStore) UpdateCredentialsValidation(ctx context.Context, accountID string, ok bool, errMsg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cred == nil {
		return nil
	}
	if accountID != "" {
		m.cred.AccountID = accountID
	}
	if ok {
		t := time.Now()
		m.cred.Validated = &t
		m.cred.Err = ""
	} else {
		m.cred.Err = errMsg
	}
	return nil
}

func (m *MemoryStore) ListEnv(ctx context.Context) ([]EnvEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]EnvEntry, 0, len(m.env))
	for _, e := range m.env {
		out = append(out, e)
	}
	return out, nil
}

func (m *MemoryStore) UpsertEnv(ctx context.Context, key, value string, isPublic bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.env[key] = EnvEntry{Key: key, Value: value, IsPublic: isPublic, UpdatedAt: time.Now()}
	return nil
}

func (m *MemoryStore) GetEnvMap(ctx context.Context) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]string, len(m.env))
	for k, e := range m.env {
		out[k] = e.Value
	}
	return out, nil
}

func (m *MemoryStore) GetTunnelConfig(ctx context.Context) (*TunnelConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t := m.tunnel
	return &t, nil
}

func (m *MemoryStore) SaveTunnelConfig(ctx context.Context, tunnelID, name, installCmd, origin string, active bool, status string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tunnel.TunnelID = tunnelID
	m.tunnel.TunnelName = name
	m.tunnel.InstallCommand = installCmd
	m.tunnel.OriginURL = origin
	m.tunnel.IsActive = active
	m.tunnel.ConnectorStatus = status
	return nil
}

func (m *MemoryStore) UpdateTunnelStatus(ctx context.Context, status string, lastSeen *time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tunnel.ConnectorStatus = status
	m.tunnel.LastSeenAt = lastSeen
	return nil
}

func (m *MemoryStore) ListTunnelRoutes(ctx context.Context) ([]TunnelRoute, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]TunnelRoute, len(m.routes))
	copy(out, m.routes)
	return out, nil
}

func (m *MemoryStore) UpsertTunnelRoute(ctx context.Context, r TunnelRoute) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r.ID > 0 {
		for i := range m.routes {
			if m.routes[i].ID == r.ID {
				m.routes[i] = r
				return r.ID, nil
			}
		}
	}
	r.ID = m.nextRouteID
	m.nextRouteID++
	m.routes = append(m.routes, r)
	return r.ID, nil
}

func (m *MemoryStore) DeleteTunnelRoute(ctx context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, r := range m.routes {
		if r.ID == id {
			m.routes = append(m.routes[:i], m.routes[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MemoryStore) GetPagesProject(ctx context.Context, projectType string) (*PagesProject, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if p, ok := m.pages[projectType]; ok {
		pp := p
		return &pp, nil
	}
	return nil, nil
}

func (m *MemoryStore) SavePagesProject(ctx context.Context, pp PagesProject) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pages[pp.ProjectType] = pp
	return nil
}

func (m *MemoryStore) UpdatePagesDeploy(ctx context.Context, projectType, status string, deployedAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	p := m.pages[projectType]
	p.DeployStatus = status
	p.LastDeployAt = &deployedAt
	m.pages[projectType] = p
	return nil
}

func (m *MemoryStore) InsertAPILog(ctx context.Context, method, path string, status int, success bool, errMsg string, durationMs int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append([]APILog{{
		ID: int64(len(m.logs) + 1), Method: method, Path: path, StatusCode: status,
		Success: success, ErrorMessage: errMsg, DurationMs: durationMs, CreatedAt: time.Now(),
	}}, m.logs...)
	if len(m.logs) > 200 {
		m.logs = m.logs[:200]
	}
	return nil
}

func (m *MemoryStore) ListAPILogs(ctx context.Context, limit int) ([]APILog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if limit <= 0 || limit > len(m.logs) {
		limit = len(m.logs)
	}
	return m.logs[:limit], nil
}

// EncodeGlobalSecret JSON for global api key + email
func EncodeGlobalSecret(key, email string) ([]byte, error) {
	return json.Marshal(map[string]string{"key": key, "email": email})
}

func DecodeGlobalSecret(b []byte) (key, email string, err error) {
	var m map[string]string
	if err = json.Unmarshal(b, &m); err != nil {
		return "", "", err
	}
	return m["key"], m["email"], nil
}

// EncryptSecretForStore encrypts for postgres memory path with key
func EncryptSecretForStore(encKey []byte, plain []byte) ([]byte, []byte, error) {
	return crypto.EncryptAESGCM(encKey, plain)
}
