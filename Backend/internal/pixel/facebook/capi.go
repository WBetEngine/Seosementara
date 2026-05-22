package facebook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const graphAPIVersion = "v21.0"

// CAPIClient sends events to Meta Conversions API.
type CAPIClient struct {
	HTTP     *http.Client
	BaseURL  string
}

func NewCAPIClient() *CAPIClient {
	return &CAPIClient{
		HTTP: &http.Client{Timeout: 15 * time.Second},
		BaseURL: fmt.Sprintf("https://graph.facebook.com/%s", graphAPIVersion),
	}
}

type UserData struct {
	ClientIPAddress string `json:"client_ip_address,omitempty"`
	ClientUserAgent string `json:"client_user_agent,omitempty"`
	FBP             string `json:"fbp,omitempty"`
	FBC             string `json:"fbc,omitempty"`
	EmailHash       string `json:"em,omitempty"`
	PhoneHash       string `json:"ph,omitempty"`
}

type ServerEvent struct {
	EventName      string   `json:"event_name"`
	EventTime      int64    `json:"event_time"`
	EventID        string   `json:"event_id"`
	ActionSource   string   `json:"action_source"`
	EventSourceURL string   `json:"event_source_url,omitempty"`
	UserData       UserData `json:"user_data"`
	CustomData     map[string]any `json:"custom_data,omitempty"`
}

type EventsPayload struct {
	Data         []ServerEvent `json:"data"`
	TestEventCode string       `json:"test_event_code,omitempty"`
}

type EventsResponse struct {
	EventsReceived int      `json:"events_received"`
	Messages       []string `json:"messages"`
	FBTraceID      string   `json:"fbtrace_id"`
}

func (c *CAPIClient) SendEvents(ctx context.Context, pixelID, accessToken string, payload EventsPayload) (*EventsResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s/events?access_token=%s", c.BaseURL, pixelID, accessToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("meta CAPI HTTP %d: %s", resp.StatusCode, string(raw))
	}

	var out EventsResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w (body: %s)", err, string(raw))
	}
	return &out, nil
}

// DebugToken validates access token (basic health check).
func (c *CAPIClient) DebugToken(ctx context.Context, inputToken, appAccessToken string) (map[string]any, error) {
	url := fmt.Sprintf("%s/debug_token?input_token=%s&access_token=%s",
		c.BaseURL, inputToken, appAccessToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("debug_token HTTP %d: %s", resp.StatusCode, string(raw))
	}
	var out struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
