package facebook

import (
	"fmt"
	"strings"
	"time"
)

// UserData matches Meta Conversions API customer information (Plan/25 §3).
type UserData struct {
	ClientIPAddress string   `json:"client_ip_address,omitempty"`
	ClientUserAgent string   `json:"client_user_agent,omitempty"`
	FBP             string   `json:"fbp,omitempty"`
	FBC             string   `json:"fbc,omitempty"`
	EM              []string `json:"em,omitempty"`
	PH              []string `json:"ph,omitempty"`
	FN              []string `json:"fn,omitempty"`
	LN              []string `json:"ln,omitempty"`
	ExternalID      []string `json:"external_id,omitempty"`
	Country         []string `json:"country,omitempty"`
	LeadID          string   `json:"lead_id,omitempty"`
}

// HasValidCustomerInfo returns true if at least one Meta-recognized matching parameter is present.
func (u UserData) HasValidCustomerInfo() bool {
	return len(u.EM) > 0 || len(u.PH) > 0 || len(u.FN) > 0 || len(u.LN) > 0 ||
		len(u.ExternalID) > 0 || len(u.Country) > 0 ||
		u.FBP != "" || u.FBC != "" || u.LeadID != ""
}

// QualityTier estimates EMQ tier per Plan/25 §7 (ingest diagnostics).
func (u UserData) QualityTier(hasEventSourceURL bool) string {
	if !hasEventSourceURL || u.ClientUserAgent == "" {
		return "D"
	}
	if !u.HasValidCustomerInfo() {
		return "D"
	}
	hasContact := len(u.EM) > 0 || len(u.PH) > 0
	hasClick := u.FBC != ""
	hasDedup := u.FBP != "" && len(u.ExternalID) > 0
	if hasContact && u.FBP != "" && hasDedup {
		return "A"
	}
	if hasContact && u.FBP != "" {
		return "B"
	}
	if u.FBP != "" || hasContact {
		return "C"
	}
	if hasClick {
		return "C"
	}
	return "D"
}

// BuildUserDataFromProps builds user_data from queued collect payload (hashed fields stored at ingest).
func BuildUserDataFromProps(props map[string]any) UserData {
	ip, _ := props["client_ip"].(string)
	ua, _ := props["user_agent"].(string)
	fbp, _ := props["fbp"].(string)
	fbc, _ := props["fbc"].(string)
	fbclid, _ := props["fbclid"].(string)
	leadID, _ := props["lead_id"].(string)

	clickTime := int64(0)
	if t, ok := props["fbc_click_time"].(float64); ok {
		clickTime = int64(t)
	}

	ud := UserData{
		ClientIPAddress: ip,
		ClientUserAgent: ua,
		FBP:             strings.TrimSpace(fbp),
		FBC:             ResolveFBC(fbc, fbclid, clickTime),
		LeadID:          strings.TrimSpace(leadID),
	}

	ud.EM = readHashedField(props, "em", "email", HashEmail)
	ud.PH = readHashedField(props, "ph", "phone", func(s string) string { return HashPhone(s, DefaultPhoneCountry) })
	ud.FN = readHashedField(props, "fn", "first_name", HashName)
	ud.LN = readHashedField(props, "ln", "last_name", HashName)
	ud.ExternalID = readHashedField(props, "external_id", "external_id", HashExternalID)
	ud.Country = readHashedField(props, "country", "country", HashCountry)

	return ud
}

func readHashedField(props map[string]any, hashKey, plainKey string, hashFn func(string) string) []string {
	var out []string
	if v, ok := props[hashKey]; ok {
		out = mergeStringSlices(out, v, nil)
	}
	// legacy ingest keys
	if hashKey == "em" {
		if v, ok := props["email_hash"]; ok {
			out = mergeStringSlices(out, v, nil)
		}
	}
	if hashKey == "ph" {
		if v, ok := props["phone_hash"]; ok {
			out = mergeStringSlices(out, v, nil)
		}
	}
	if plain, ok := props[plainKey].(string); ok && plain != "" {
		out = appendHashed(out, hashFn, plain)
	}
	return out
}

// EnrichCollectProps hashes PII and normalizes tracking fields for queue storage (Plan/25 §12).
func EnrichCollectProps(props map[string]any, email, phone, firstName, lastName, externalID, country, fbclid string, defaultPhoneCountry string) map[string]any {
	if props == nil {
		props = map[string]any{}
	}
	if defaultPhoneCountry == "" {
		defaultPhoneCountry = DefaultPhoneCountry
	}
	if email != "" {
		props["em"] = hashField(HashEmail, email)
	}
	if phone != "" {
		props["ph"] = hashField(func(s string) string { return HashPhone(s, defaultPhoneCountry) }, phone)
	}
	if firstName != "" {
		props["fn"] = hashField(HashName, firstName)
	}
	if lastName != "" {
		props["ln"] = hashField(HashName, lastName)
	}
	if externalID != "" {
		props["external_id"] = hashField(HashExternalID, externalID)
	}
	if country != "" {
		props["country"] = hashField(HashCountry, country)
	}
	if fbclid != "" {
		props["fbclid"] = strings.TrimSpace(fbclid)
		if _, ok := props["fbc_click_time"]; !ok {
			props["fbc_click_time"] = float64(time.Now().Unix())
		}
	}
	if fbc := ResolveFBC(str(props["fbc"]), fbclid, int64(float64Prop(props["fbc_click_time"]))); fbc != "" {
		props["fbc"] = fbc
	}
	if fbp := EnsureFBP(str(props["fbp"])); fbp != "" {
		props["fbp"] = fbp
	}
	return props
}

func str(v any) string {
	s, _ := v.(string)
	return strings.TrimSpace(s)
}

func float64Prop(v any) int64 {
	switch t := v.(type) {
	case float64:
		return int64(t)
	case int64:
		return t
	case int:
		return int64(t)
	default:
		return time.Now().Unix()
	}
}

// BuildServerEvent assembles a CAPI server event from queued props (Plan/25–26).
func BuildServerEvent(eventName, eventID string, eventTime int64, props map[string]any) (ServerEvent, error) {
	url, _ := props["url"].(string)
	if url == "" {
		return ServerEvent{}, fmt.Errorf("event_source_url wajib untuk website")
	}
	ua, _ := props["user_agent"].(string)
	if ua == "" {
		return ServerEvent{}, fmt.Errorf("client_user_agent wajib untuk website")
	}

	ud := BuildUserDataFromProps(props)
	if !ud.HasValidCustomerInfo() {
		return ServerEvent{}, fmt.Errorf("minimal satu parameter customer info valid (em, ph, fbp, fbc, external_id, dll.)")
	}

	custom := BuildCustomData(eventName, props)
	if err := ValidateCustomData(eventName, custom); err != nil {
		return ServerEvent{}, err
	}

	if eventID == "" {
		return ServerEvent{}, fmt.Errorf("event_id wajib")
	}
	if eventTime <= 0 {
		eventTime = time.Now().Unix()
	}

	return ServerEvent{
		EventName:      eventName,
		EventTime:      eventTime,
		EventID:        eventID,
		ActionSource:   "website",
		EventSourceURL: url,
		UserData:       ud,
		CustomData:     custom,
	}, nil
}
