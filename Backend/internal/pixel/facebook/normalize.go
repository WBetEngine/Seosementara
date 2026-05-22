package facebook

import (
	"regexp"
	"strings"
	"unicode"
)

const DefaultPhoneCountry = "62"

var nonDigit = regexp.MustCompile(`[^0-9]+`)

// NormalizeEmail trims and lowercases before hashing.
func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

// NormalizePhone returns digits-only E.164-style without plus (default country Indonesia 62).
func NormalizePhone(phone, defaultCountry string) string {
	if defaultCountry == "" {
		defaultCountry = DefaultPhoneCountry
	}
	d := nonDigit.ReplaceAllString(phone, "")
	if d == "" {
		return ""
	}
	if strings.HasPrefix(d, "0") {
		d = defaultCountry + strings.TrimPrefix(d, "0")
	}
	if !strings.HasPrefix(d, defaultCountry) {
		d = defaultCountry + d
	}
	return d
}

// NormalizeName lowercases and strips punctuation for fn/ln hashing.
func NormalizeName(name string) string {
	var b strings.Builder
	for _, r := range strings.TrimSpace(strings.ToLower(name)) {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// NormalizeCountry returns ISO 3166-1 alpha-2 lowercase.
func NormalizeCountry(country string) string {
	c := strings.TrimSpace(strings.ToLower(country))
	if len(c) == 2 {
		return c
	}
	return ""
}

// HashEmail hashes normalized email for Meta em[].
func HashEmail(email string) string {
	return hashNormalized(NormalizeEmail(email))
}

// HashPhone hashes normalized phone for Meta ph[].
func HashPhone(phone, defaultCountry string) string {
	return hashNormalized(NormalizePhone(phone, defaultCountry))
}

// HashName hashes normalized name for fn/ln[].
func HashName(name string) string {
	return hashNormalized(NormalizeName(name))
}

// HashCountry hashes normalized country code.
func HashCountry(country string) string {
	return hashNormalized(NormalizeCountry(country))
}

// HashExternalID hashes stable user id (recommended by Meta).
func HashExternalID(id string) string {
	return hashNormalized(strings.TrimSpace(id))
}

func hashNormalized(v string) string {
	if v == "" {
		return ""
	}
	sum := sha256Sum(v)
	return sum
}
