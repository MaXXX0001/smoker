package domain

import (
	"testing"
	"time"
)

// mkTime — допоміжний конструктор локального часу для тестів.
func mkTime(h, m int) time.Time {
	return time.Date(2026, 6, 15, h, m, 0, 0, time.UTC)
}

func TestDayProgress42(t *testing.T) {
	// 42% доби ≈ 10:04:48.
	c := DayProgressCondition(mkTime(10, 5))
	if c.Score != 2 {
		t.Fatalf("очікував пасхалку 42%%, отримав %+v", c)
	}
}

func TestDayProgressGeneric(t *testing.T) {
	c := DayProgressCondition(mkTime(8, 0))
	if c.Verdict != Neutral || c.Headline == "" {
		t.Fatalf("звичайний час: %+v", c)
	}
}

func TestSpecialClockSymmetry(t *testing.T) {
	c := SpecialClockCondition(mkTime(15, 15))
	if c.Verdict != Favorable {
		t.Fatalf("15:15 симетрія має сприяти: %+v", c)
	}
}
