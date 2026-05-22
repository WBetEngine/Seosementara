package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DefaultAPIBase = "https://api.cloudflare.com/client/v4"

type Client struct {
	BaseURL    string
	HTTP       *http.Client
	AuthType   string
	APIToken   string
	GlobalKey  string
	Email      string
	OnAudit    func(method, path string, status int, ok bool, errMsg string, ms int)
}

type apiResponse struct {
	Success bool            `json:"success"`
	Errors  []apiError      `json:"errors"`
	Result  json.RawMessage `json:"result"`
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewClient(base string) *Client {
	if base == "" {
		base = DefaultAPIBase
	}
	return &Client{
		BaseURL: base,
		HTTP:    &http.Client{Timeout: 45 * time.Second},
	}
}

func (c *Client) WithAPIToken(token string) *Client {
	c.AuthType = "api_token"
	c.APIToken = token
	return c
}

func (c *Client) WithGlobalKey(email, key string) *Client {
	c.AuthType = "global_api_key"
	c.Email = email
	c.GlobalKey = key
	return c
}

func (c *Client) do(ctx context.Context, method, path string, body any, out any) error {
	start := time.Now()
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, rdr)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.AuthType == "api_token" || c.APIToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIToken)
	} else {
		req.Header.Set("X-Auth-Email", c.Email)
		req.Header.Set("X-Auth-Key", c.GlobalKey)
	}
	resp, err := c.HTTP.Do(req)
	ms := int(time.Since(start).Milliseconds())
	if err != nil {
		if c.OnAudit != nil {
			c.OnAudit(method, path, 0, false, err.Error(), ms)
		}
		return err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var wrap apiResponse
	_ = json.Unmarshal(raw, &wrap)
	ok := resp.StatusCode >= 200 && resp.StatusCode < 300 && wrap.Success
	errMsg := ""
	if !ok {
		if len(wrap.Errors) > 0 {
			errMsg = wrap.Errors[0].Message
		} else {
			errMsg = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
	}
	if c.OnAudit != nil {
		c.OnAudit(method, path, resp.StatusCode, ok, errMsg, ms)
	}
	if !ok {
		if errMsg == "" {
			errMsg = string(raw)
		}
		return fmt.Errorf("cloudflare api: %s", errMsg)
	}
	if out != nil && len(wrap.Result) > 0 && string(wrap.Result) != "null" {
		return json.Unmarshal(wrap.Result, out)
	}
	return nil
}

func (c *Client) VerifyToken(ctx context.Context) error {
	var result struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	return c.do(ctx, http.MethodGet, "/user/tokens/verify", nil, &result)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
	var result []Account
	if err := c.do(ctx, http.MethodGet, "/accounts", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

type Tunnel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) CreateTunnel(ctx context.Context, accountID, name string) (*Tunnel, error) {
	var result Tunnel
	err := c.do(ctx, http.MethodPost, "/accounts/"+accountID+"/cfd_tunnel", map[string]string{
		"name": name,
		"config_src": "cloudflare",
	}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type TunnelToken struct {
	Token string `json:"token"`
}

func (c *Client) GetTunnelToken(ctx context.Context, accountID, tunnelID string) (string, error) {
	var result TunnelToken
	err := c.do(ctx, http.MethodGet, "/accounts/"+accountID+"/cfd_tunnel/"+tunnelID+"/token", nil, &result)
	if err != nil {
		return "", err
	}
	return result.Token, nil
}

type IngressRule struct {
	Hostname string `json:"hostname"`
	Path     string `json:"path,omitempty"`
	Service  string `json:"service"`
}

func (c *Client) ConfigureTunnel(ctx context.Context, accountID, tunnelID string, rules []IngressRule) error {
	body := map[string]any{
		"config": map[string]any{
			"ingress": rules,
		},
	}
	return c.do(ctx, http.MethodPut, "/accounts/"+accountID+"/cfd_tunnel/"+tunnelID+"/configurations", body, nil)
}

type TunnelConnection struct {
	ID       string `json:"id"`
	IsActive bool   `json:"is_active"`
}

func (c *Client) TunnelConnections(ctx context.Context, accountID, tunnelID string) ([]TunnelConnection, error) {
	var result []TunnelConnection
	err := c.do(ctx, http.MethodGet, "/accounts/"+accountID+"/cfd_tunnel/"+tunnelID+"/connections", nil, &result)
	return result, err
}

type PagesProject struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Subdomain         string `json:"subdomain"`
	ProductionBranch  string `json:"production_branch"`
}

func (c *Client) ListPagesProjects(ctx context.Context, accountID string) ([]PagesProject, error) {
	var result []PagesProject
	err := c.do(ctx, http.MethodGet, "/accounts/"+accountID+"/pages/projects", nil, &result)
	return result, err
}

func (c *Client) CreatePagesDeployment(ctx context.Context, accountID, projectName string) error {
	return c.do(ctx, http.MethodPost, "/accounts/"+accountID+"/pages/projects/"+projectName+"/deployments", map[string]any{}, nil)
}

func (c *Client) UpsertPagesEnv(ctx context.Context, accountID, projectName, key, value string, isSecret bool) error {
	t := "plain_text"
	if isSecret {
		t = "secret_text"
	}
	body := map[string]any{
		"name":  key,
		"value": value,
		"type":  t,
	}
	return c.do(ctx, http.MethodPost, "/accounts/"+accountID+"/pages/projects/"+projectName+"/environment/variables", body, nil)
}

type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
}

func (c *Client) ListDNSRecords(ctx context.Context, zoneID, name string) ([]DNSRecord, error) {
	path := fmt.Sprintf("/zones/%s/dns_records?name=%s", zoneID, name)
	var result []DNSRecord
	err := c.do(ctx, http.MethodGet, path, nil, &result)
	return result, err
}

func (c *Client) CreateDNSRecord(ctx context.Context, zoneID, rtype, name, content string, proxied bool) error {
	body := map[string]any{
		"type":    rtype,
		"name":    name,
		"content": content,
		"proxied": proxied,
	}
	return c.do(ctx, http.MethodPost, "/zones/"+zoneID+"/dns_records", body, nil)
}
