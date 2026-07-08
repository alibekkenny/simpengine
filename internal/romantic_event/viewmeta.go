package romanticevent

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"strings"

	"github.com/mileusna/useragent"
)

// ViewMeta is request-derived data needed to record one public open.
type ViewMeta struct {
	VisitorID string
	Device    string
	OS        string
	Browser   string
	IP        string
	IsOwner   bool
	IsBot     bool
}

const visitorCookieName = "visitor_id"

// ClientIP returns the first hop of X-Forwarded-For, else the RemoteAddr host.
func ClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// ResolveVisitorID uses the cookie value if present; otherwise it returns a
// deterministic sha256(ip|ua) fingerprint and setCookie=true.
func ResolveVisitorID(cookieVal, ip, ua string) (id string, setCookie bool) {
	if strings.TrimSpace(cookieVal) != "" {
		return cookieVal, false
	}
	sum := sha256.Sum256([]byte(ip + "|" + ua))
	return hex.EncodeToString(sum[:]), true
}

// ParseUserAgent extracts a friendly device / OS / browser from a UA string.
func ParseUserAgent(ua string) (device, os, browser string) {
	p := useragent.Parse(ua)
	device = p.Device
	if device == "" {
		switch {
		case p.Mobile:
			device = "Mobile"
		case p.Tablet:
			device = "Tablet"
		default:
			device = "Desktop"
		}
	}
	return device, p.OS, p.Name
}

var botUASubstrings = []string{
	"telegrambot", "whatsapp", "twitterbot", "facebookexternalhit",
	"slackbot", "discordbot", "linkedinbot", "googlebot", "bingbot",
	"applebot", "redditbot", "bot", "crawler", "spider",
}

// IsBotUserAgent reports whether the UA looks like a crawler / link-preview fetcher.
func IsBotUserAgent(ua string) bool {
	if strings.TrimSpace(ua) == "" {
		return true
	}
	if useragent.Parse(ua).Bot {
		return true
	}
	low := strings.ToLower(ua)
	for _, s := range botUASubstrings {
		if strings.Contains(low, s) {
			return true
		}
	}
	return false
}
