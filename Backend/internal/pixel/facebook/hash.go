package facebook

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// HashPII normalizes then SHA256-hashes for Meta user_data (em, ph).
func HashPII(value string) string {
	v := strings.TrimSpace(strings.ToLower(value))
	if v == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(v))
	return hex.EncodeToString(sum[:])
}
