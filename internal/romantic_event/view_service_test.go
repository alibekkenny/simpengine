package romanticevent

import (
	"strings"
	"testing"
	"time"
)

func TestShouldCountView(t *testing.T) {
	cases := []struct {
		name string
		meta ViewMeta
		want bool
	}{
		{"target", ViewMeta{}, true},
		{"owner", ViewMeta{IsOwner: true}, false},
		{"bot", ViewMeta{IsBot: true}, false},
		{"owner-and-bot", ViewMeta{IsOwner: true, IsBot: true}, false},
	}
	for _, c := range cases {
		if got := shouldCountView(c.meta); got != c.want {
			t.Errorf("%s: got %v want %v", c.name, got, c.want)
		}
	}
}

func TestBuildViewPingMessage(t *testing.T) {
	when := time.Date(2026, 7, 8, 14, 32, 0, 0, time.UTC)
	msg := buildViewPingMessage("Diana", "Coffee, maybe?", "iPhone", "iOS", "Safari", when, 3)
	for _, want := range []string{"Diana", "Coffee, maybe?", "iPhone", "Safari", "14:32", "#3"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message %q missing %q", msg, want)
		}
	}
	// HTML-escaped and falls back to "Someone" without a name.
	esc := buildViewPingMessage("", "A & B <x>", "", "", "", when, 1)
	if !strings.Contains(esc, "Someone") || !strings.Contains(esc, "A &amp; B &lt;x&gt;") {
		t.Errorf("escaping/fallback wrong: %q", esc)
	}
}
