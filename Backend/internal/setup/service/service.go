package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/WBetEngine/Seosementara/Backend/internal/cloudflare"
	"github.com/WBetEngine/Seosementara/Backend/internal/setup/store"
)

type SetupService struct {
	Store  store.SetupStore
	EncKey []byte
}

func NewSetupService(st store.SetupStore, encKey []byte) *SetupService {
	return &SetupService{Store: st, EncKey: encKey}
}

func (s *SetupService) cfClient(ctx context.Context) (*cloudflare.Client, error) {
	authType, secret, err := s.Store.GetCredentialsSecret(ctx, s.EncKey)
	if err != nil {
		return nil, err
	}
	if secret == "" {
		return nil, errors.New("cloudflare credentials belum dikonfigurasi")
	}
	c := cloudflare.NewClient("")
	c.OnAudit = func(method, path string, status int, ok bool, errMsg string, ms int) {
		_ = s.Store.InsertAPILog(ctx, method, path, status, ok, errMsg, ms)
	}
	if authType == store.AuthTypeGlobalAPIKey {
		key, email, err := store.DecodeGlobalSecret([]byte(secret))
		if err != nil {
			return nil, err
		}
		return c.WithGlobalKey(email, key), nil
	}
	return c.WithAPIToken(secret), nil
}

func (s *SetupService) accountID(ctx context.Context) (string, error) {
	cred, err := s.Store.GetCredentials(ctx)
	if err != nil {
		return "", err
	}
	if cred == nil || cred.AccountID == "" {
		return "", errors.New("account_id belum diisi — test koneksi dulu")
	}
	return cred.AccountID, nil
}

// --- Credentials ---

type CredentialsView struct {
	AuthType        string     `json:"auth_type"`
	AccountID       string     `json:"account_id"`
	AccountEmail    string     `json:"account_email"`
	TokenMasked     string     `json:"token_masked"`
	LastValidatedAt *time.Time `json:"last_validated_at"`
	ValidationError string     `json:"validation_error"`
	Configured      bool       `json:"configured"`
}

func (s *SetupService) GetCredentialsView(ctx context.Context) (*CredentialsView, error) {
	c, err := s.Store.GetCredentials(ctx)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return &CredentialsView{AuthType: store.AuthTypeAPIToken}, nil
	}
	return &CredentialsView{
		AuthType: c.AuthType, AccountID: c.AccountID, AccountEmail: c.AccountEmail,
		TokenMasked: "cfpat_••••••••", LastValidatedAt: c.LastValidatedAt,
		ValidationError: c.ValidationError, Configured: true,
	}, nil
}

func (s *SetupService) SaveCredentials(ctx context.Context, in store.CredentialsInput, test bool) (*CredentialsView, error) {
	var secret []byte
	var authType = in.AuthType
	if authType == "" {
		authType = store.AuthTypeAPIToken
	}
	switch authType {
	case store.AuthTypeAPIToken:
		if in.APIToken == "" {
			return nil, errors.New("api_token wajib")
		}
		secret = []byte(in.APIToken)
	case store.AuthTypeGlobalAPIKey:
		if in.GlobalAPIKey == "" || in.AccountEmail == "" {
			return nil, errors.New("global_api_key dan account_email wajib")
		}
		b, err := store.EncodeGlobalSecret(in.GlobalAPIKey, in.AccountEmail)
		if err != nil {
			return nil, err
		}
		secret = b
	default:
		return nil, errors.New("auth_type tidak valid")
	}
	if err := s.Store.SaveCredentials(ctx, authType, secret, in.AccountID, in.AccountEmail, s.EncKey); err != nil {
		return nil, err
	}
	if test {
		if err := s.TestCredentials(ctx); err != nil {
			_ = s.Store.UpdateCredentialsValidation(ctx, in.AccountID, false, err.Error())
			return s.GetCredentialsView(ctx)
		}
	}
	return s.GetCredentialsView(ctx)
}

func (s *SetupService) TestCredentials(ctx context.Context) error {
	cf, err := s.cfClient(ctx)
	if err != nil {
		return err
	}
	if err := cf.VerifyToken(ctx); err != nil {
		// global key: verify via accounts list
		if cf.AuthType == "global_api_key" || cf.GlobalKey != "" {
			accts, err2 := cf.ListAccounts(ctx)
			if err2 != nil {
				_ = s.Store.UpdateCredentialsValidation(ctx, "", false, err2.Error())
				return err2
			}
			if len(accts) == 0 {
				return errors.New("tidak ada account cloudflare")
			}
			_ = s.Store.UpdateCredentialsValidation(ctx, accts[0].ID, true, "")
			return nil
		}
		_ = s.Store.UpdateCredentialsValidation(ctx, "", false, err.Error())
		return err
	}
	accts, err := cf.ListAccounts(ctx)
	if err != nil {
		_ = s.Store.UpdateCredentialsValidation(ctx, "", false, err.Error())
		return err
	}
	aid := ""
	if len(accts) > 0 {
		aid = accts[0].ID
	}
	_ = s.Store.UpdateCredentialsValidation(ctx, aid, true, "")
	return nil
}

// --- Domain env ---

func (s *SetupService) ListEnv(ctx context.Context) ([]store.EnvEntry, error) {
	return s.Store.ListEnv(ctx)
}

func (s *SetupService) SaveEnvBulk(ctx context.Context, entries map[string]string) error {
	secretKeys := map[string]bool{
		"DATABASE_URL": true, "SESSION_SECRET": true, "MASTER_ENCRYPTION_KEY": true,
	}
	for k, v := range entries {
		if err := s.Store.UpsertEnv(ctx, k, v, !secretKeys[k]); err != nil {
			return err
		}
	}
	return nil
}

func (s *SetupService) SyncPagesEnv(ctx context.Context) error {
	cf, err := s.cfClient(ctx)
	if err != nil {
		return err
	}
	aid, err := s.accountID(ctx)
	if err != nil {
		return err
	}
	pp, err := s.Store.GetPagesProject(ctx, store.ProjectTypeAdmin)
	if err != nil || pp == nil {
		return errors.New("pages project admin belum dikonfigurasi")
	}
	env, err := s.Store.GetEnvMap(ctx)
	if err != nil {
		return err
	}
	for k, v := range env {
		if v == "" {
			continue
		}
		if err := cf.UpsertPagesEnv(ctx, aid, pp.ProjectName, k, v, false); err != nil {
			return fmt.Errorf("sync %s: %w", k, err)
		}
	}
	return nil
}

// --- Tunnel ---

type TunnelView struct {
	Config store.TunnelConfig `json:"config"`
	Routes []store.TunnelRoute `json:"routes"`
}

func (s *SetupService) GetTunnelView(ctx context.Context) (*TunnelView, error) {
	cfg, err := s.Store.GetTunnelConfig(ctx)
	if err != nil {
		return nil, err
	}
	routes, err := s.Store.ListTunnelRoutes(ctx)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &store.TunnelConfig{TunnelName: "seosementara-api", OriginURL: "http://127.0.0.1:8080"}
	}
	return &TunnelView{Config: *cfg, Routes: routes}, nil
}

func (s *SetupService) CreateTunnel(ctx context.Context, name, origin string) (*TunnelView, error) {
	cf, err := s.cfClient(ctx)
	if err != nil {
		return nil, err
	}
	aid, err := s.accountID(ctx)
	if err != nil {
		return nil, err
	}
	if name == "" {
		name = "seosementara-api"
	}
	if origin == "" {
		origin = "http://127.0.0.1:8080"
	}
	tun, err := cf.CreateTunnel(ctx, aid, name)
	if err != nil {
		return nil, err
	}
	token, err := cf.GetTunnelToken(ctx, aid, tun.ID)
	if err != nil {
		return nil, err
	}
	install := fmt.Sprintf("sudo cloudflared service install %s\nsudo systemctl enable cloudflared\nsudo systemctl start cloudflared", token)
	_ = s.Store.SaveTunnelConfig(ctx, tun.ID, name, install, origin, true, store.ConnectorUnknown)
	return s.GetTunnelView(ctx)
}

func (s *SetupService) RefreshTunnelStatus(ctx context.Context) error {
	cf, err := s.cfClient(ctx)
	if err != nil {
		return err
	}
	aid, err := s.accountID(ctx)
	if err != nil {
		return err
	}
	cfg, err := s.Store.GetTunnelConfig(ctx)
	if err != nil || cfg == nil || cfg.TunnelID == "" {
		return errors.New("tunnel belum dibuat")
	}
	conns, err := cf.TunnelConnections(ctx, aid, cfg.TunnelID)
	if err != nil {
		_ = s.Store.UpdateTunnelStatus(ctx, store.ConnectorDown, nil)
		return err
	}
	status := store.ConnectorDown
	var last *time.Time
	if len(conns) > 0 {
		status = store.ConnectorHealthy
		t := time.Now()
		last = &t
	}
	_ = s.Store.UpdateTunnelStatus(ctx, status, last)
	return nil
}

func (s *SetupService) ApplyTunnelRoutes(ctx context.Context) error {
	cf, err := s.cfClient(ctx)
	if err != nil {
		return err
	}
	aid, err := s.accountID(ctx)
	if err != nil {
		return err
	}
	cfg, err := s.Store.GetTunnelConfig(ctx)
	if err != nil || cfg == nil || cfg.TunnelID == "" {
		return errors.New("tunnel belum dibuat")
	}
	routes, err := s.Store.ListTunnelRoutes(ctx)
	if err != nil {
		return err
	}
	var ingress []cloudflare.IngressRule
	for _, r := range routes {
		if !r.IsEnabled {
			continue
		}
		svc := r.ServiceURL
		if svc == "" {
			svc = cfg.OriginURL
		}
		ingress = append(ingress, cloudflare.IngressRule{
			Hostname: r.Hostname,
			Path:     r.PathPrefix,
			Service:  svc,
		})
	}
	if len(ingress) == 0 {
		return errors.New("tidak ada route aktif")
	}
	ingress = append(ingress, cloudflare.IngressRule{Service: "http_status:404"})
	return cf.ConfigureTunnel(ctx, aid, cfg.TunnelID, ingress)
}

func (s *SetupService) SaveTunnelRoute(ctx context.Context, r store.TunnelRoute) error {
	_, err := s.Store.UpsertTunnelRoute(ctx, r)
	return err
}

// --- Pages ---

func (s *SetupService) GetPagesAdmin(ctx context.Context) (*store.PagesProject, error) {
	return s.Store.GetPagesProject(ctx, store.ProjectTypeAdmin)
}

func (s *SetupService) SavePagesAdmin(ctx context.Context, pp store.PagesProject) error {
	pp.ProjectType = store.ProjectTypeAdmin
	return s.Store.SavePagesProject(ctx, pp)
}

func (s *SetupService) TriggerPagesDeploy(ctx context.Context) error {
	cf, err := s.cfClient(ctx)
	if err != nil {
		return err
	}
	aid, err := s.accountID(ctx)
	if err != nil {
		return err
	}
	pp, err := s.Store.GetPagesProject(ctx, store.ProjectTypeAdmin)
	if err != nil || pp == nil {
		return errors.New("pages admin tidak ada")
	}
	err = cf.CreatePagesDeployment(ctx, aid, pp.ProjectName)
	status := "success"
	if err != nil {
		status = "failed: " + err.Error()
	}
	_ = s.Store.UpdatePagesDeploy(ctx, store.ProjectTypeAdmin, status, time.Now())
	return err
}

// --- DNS ---

type DNSApplyResult struct {
	Created []string `json:"created"`
	Skipped []string `json:"skipped"`
	Errors  []string `json:"errors"`
}

func (s *SetupService) ApplyDNS(ctx context.Context) (*DNSApplyResult, error) {
	cf, err := s.cfClient(ctx)
	if err != nil {
		return nil, err
	}
	env, err := s.Store.GetEnvMap(ctx)
	if err != nil {
		return nil, err
	}
	zoneID := env["ZONE_ID"]
	if zoneID == "" {
		return nil, errors.New("ZONE_ID kosong — isi di tab Domain & env")
	}
	cfg, _ := s.Store.GetTunnelConfig(ctx)
	tunnelTarget := "{tunnel-id}.cfargotunnel.com"
	if cfg != nil && cfg.TunnelID != "" {
		tunnelTarget = cfg.TunnelID + ".cfargotunnel.com"
	}
	domain := env["PRIMARY_DOMAIN"]
	if domain == "" {
		domain = "seosementara.org"
	}
	res := &DNSApplyResult{}
	plans := []struct{ name, content string }{
		{domain, tunnelTarget},
		{"*." + domain, tunnelTarget},
		{"www." + domain, domain},
	}
	for _, p := range plans {
		existing, _ := cf.ListDNSRecords(ctx, zoneID, p.name)
		if len(existing) > 0 {
			res.Skipped = append(res.Skipped, p.name)
			continue
		}
		rtype := "CNAME"
		if p.name == domain && !strings.Contains(p.content, ".") {
			rtype = "A"
		}
		if err := cf.CreateDNSRecord(ctx, zoneID, rtype, p.name, p.content, true); err != nil {
			res.Errors = append(res.Errors, p.name+": "+err.Error())
			continue
		}
		res.Created = append(res.Created, p.name)
	}
	return res, nil
}

func (s *SetupService) ListLogs(ctx context.Context, limit int) ([]store.APILog, error) {
	return s.Store.ListAPILogs(ctx, limit)
}

// MaskToken for display
func MaskToken(t string) string {
	if len(t) <= 8 {
		return "••••"
	}
	return t[:4] + "••••" + t[len(t)-4:]
}

// ParseEnvJSON from PUT body
func ParseEnvJSON(raw json.RawMessage) (map[string]string, error) {
	var m map[string]string
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}
