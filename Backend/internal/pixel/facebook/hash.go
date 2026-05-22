package facebook

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// HashPII is an alias for HashEmail-style normalization+hash (backward compatible).
func HashPII(value string) string {
	return HashEmail(value)
}

func sha256Sum(v string) string {
	sum := sha256.Sum256([]byte(v))
	return hex.EncodeToString(sum[:])
}

// hashField returns a single-element slice for Meta array fields, or nil if empty.
func hashField(fn func(string) string, raw string) []string {
	h := fn(raw)
	if h == "" {
		return nil
	}
	return []string{h}
}

// appendHashed adds hash to slice if raw non-empty and not already present.
func appendHashed(dst []string, fn func(string) string, raw string) []string {
	h := fn(raw)
	if h == "" {
		return dst
	}
	for _, x := range dst {
		if x == h {
			return dst
		}
	}
	return append(dst, h)
}

// mergeStringSlices merges string slices from props (pre-hashed arrays or plain strings).
func mergeStringSlices(existing []string, from any, hashFn func(string) string) []string {
	switch v := from.(type) {
	case string:
		if len(existing) > 0 && hashFn == nil {
			return existing
		}
		if hashFn != nil {
			return appendHashed(existing, hashFn, v)
		}
		if v != "" {
			return appendHashed(existing, func(s string) string { return strings.TrimSpace(s) }, v)
		}
	case []any:
		for _, item := range v {
			existing = mergeStringSlices(existing, item, hashFn)
		}
	case []string:
		for _, item := range v {
			existing = mergeStringSlices(existing, item, hashFn)
		}
	}
	return existing
}
