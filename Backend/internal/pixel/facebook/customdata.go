package facebook

import (
	"fmt"
	"strings"
)

// BuildCustomData extracts Meta custom_data from collect props (Plan/25 §5, Plan/26 §3).
func BuildCustomData(eventName string, props map[string]any) map[string]any {
	out := map[string]any{}
	if nested, ok := props["custom_data"].(map[string]any); ok {
		for k, v := range nested {
			out[k] = v
		}
	}
	keys := []string{
		"value", "currency", "order_id", "content_ids", "contents",
		"content_type", "content_name", "content_category", "num_items",
		"search_string", "predicted_ltv", "status",
	}
	for _, k := range keys {
		if v, ok := props[k]; ok && v != nil {
			out[k] = v
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// ValidateCustomData enforces Meta requirements for standard events.
func ValidateCustomData(eventName string, custom map[string]any) error {
	if custom == nil {
		custom = map[string]any{}
	}
	switch strings.TrimSpace(eventName) {
	case "Purchase":
		if !hasValue(custom) {
			return fmt.Errorf("Purchase: custom_data.value wajib")
		}
		if cur, _ := custom["currency"].(string); strings.TrimSpace(cur) == "" {
			return fmt.Errorf("Purchase: custom_data.currency wajib (ISO 4217)")
		}
	}
	return nil
}

func hasValue(custom map[string]any) bool {
	v, ok := custom["value"]
	if !ok || v == nil {
		return false
	}
	switch n := v.(type) {
	case float64:
		return n > 0 || n == 0
	case int:
		return true
	case int64:
		return true
	default:
		return fmt.Sprint(v) != "" && fmt.Sprint(v) != "0"
	}
}
