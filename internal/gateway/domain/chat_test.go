package domain

import (
	"testing"
	"time"
)

func baseChat() Chat {
	return Chat{
		ChatID: 1, Lat: 50.45, Lon: 30.52, TZ: "Europe/Kyiv",
		IntervalMinutes: 90, Enabled: true,
	}
}

func TestDueRespectsInterval(t *testing.T) {
	c := baseChat()
	// 12:00 за Києвом — в робочих годинах.
	now := time.Date(2026, 6, 15, 9, 0, 0, 0, time.UTC) // 12:00 Kyiv (UTC+3)
	c.LastSent = now.Add(-30 * time.Minute)
	if c.Due(now, 9, 19) {
		t.Fatal("30 хв < 90 хв інтервалу — не має слати")
	}
	c.LastSent = now.Add(-2 * time.Hour)
	if !c.Due(now, 9, 19) {
		t.Fatal("2 год > 90 хв — має слати")
	}
}

func TestDueOutsideWorkingHours(t *testing.T) {
	c := baseChat()
	c.LastSent = time.Time{}
	// 03:00 UTC = 06:00 Kyiv — поза 9..19.
	now := time.Date(2026, 6, 15, 3, 0, 0, 0, time.UTC)
	if c.Due(now, 9, 19) {
		t.Fatal("06:00 поза робочими годинами — не має слати")
	}
}

func TestDueDisabled(t *testing.T) {
	c := baseChat()
	c.Enabled = false
	c.LastSent = time.Time{}
	now := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)
	if c.Due(now, 9, 19) {
		t.Fatal("вимкнений чат не має слати")
	}
}
