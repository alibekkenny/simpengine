package romanticevent

import (
	"net/http/httptest"
	"testing"
)

func TestResolveVisitorID(t *testing.T) {
	if id, set := ResolveVisitorID("cookie-abc", "1.2.3.4", "UA"); id != "cookie-abc" || set {
		t.Fatalf("cookie present: got (%q,%v), want (cookie-abc,false)", id, set)
	}
	id1, set1 := ResolveVisitorID("", "1.2.3.4", "UA")
	id2, _ := ResolveVisitorID("", "1.2.3.4", "UA")
	if !set1 {
		t.Fatal("empty cookie: want setCookie=true")
	}
	if id1 != id2 || len(id1) != 64 {
		t.Fatalf("fingerprint not deterministic 64-hex: %q vs %q", id1, id2)
	}
	if other, _ := ResolveVisitorID("", "9.9.9.9", "UA"); other == id1 {
		t.Fatal("different IP should produce different fingerprint")
	}
}

func TestParseUserAgent(t *testing.T) {
	iphone := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"
	dev, os, br := ParseUserAgent(iphone)
	if dev == "" || os != "iOS" || br != "Safari" {
		t.Fatalf("iphone: got (%q,%q,%q)", dev, os, br)
	}
	desktop := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36"
	dev2, _, br2 := ParseUserAgent(desktop)
	if dev2 != "Desktop" || br2 != "Chrome" {
		t.Fatalf("desktop: got device=%q browser=%q", dev2, br2)
	}
}

func TestIsBotUserAgent(t *testing.T) {
	cases := map[string]bool{
		"":                        true,
		"TelegramBot (like TwitterBot)": true,
		"WhatsApp/2.23":           true,
		"Mozilla/5.0 (Windows NT 10.0) Chrome/120.0 Safari/537.36": false,
	}
	for ua, want := range cases {
		if got := IsBotUserAgent(ua); got != want {
			t.Errorf("IsBotUserAgent(%q)=%v want %v", ua, got, want)
		}
	}
}

func TestClientIP(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")
	if ip := ClientIP(r); ip != "1.2.3.4" {
		t.Fatalf("xff: got %q", ip)
	}
	r2 := httptest.NewRequest("GET", "/", nil)
	r2.RemoteAddr = "9.9.9.9:1234"
	if ip := ClientIP(r2); ip != "9.9.9.9" {
		t.Fatalf("remoteaddr: got %q", ip)
	}
}
