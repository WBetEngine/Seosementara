package facebook

import (
	"fmt"
	"time"
)

// FormatFBC builds Meta click cookie value: fb.1.{unix}.{fbclid}.
func FormatFBC(fbclid string, clickTime int64) string {
	fbclid = trimSpace(fbclid)
	if fbclid == "" {
		return ""
	}
	if clickTime <= 0 {
		clickTime = time.Now().Unix()
	}
	return fmt.Sprintf("fb.1.%d.%s", clickTime, fbclid)
}

// ResolveFBC prefers existing _fbc cookie; else builds from fbclid query param.
func ResolveFBC(fbc, fbclid string, clickTime int64) string {
	fbc = trimSpace(fbc)
	if fbc != "" {
		return fbc
	}
	return FormatFBC(fbclid, clickTime)
}

func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}
