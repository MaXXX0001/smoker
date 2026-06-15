package domain

import (
	"math"
	"testing"
	"time"
)

func TestMoonAgeAtReferenceIsNew(t *testing.T) {
	// На момент відлікової нової місяця вік ~0.
	age := MoonAgeDays(refNewMoon)
	if age > 0.5 && age < synodicMonth-0.5 {
		t.Fatalf("очікував вік біля 0, отримав %.2f", age)
	}
}

func TestMoonIlluminationRange(t *testing.T) {
	// Освітлення завжди в [0,1] на довільних датах.
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 400; i++ {
		ts := start.Add(time.Duration(i) * 6 * time.Hour)
		il := MoonIllumination(ts)
		if il < -1e-9 || il > 1+1e-9 {
			t.Fatalf("освітлення поза діапазоном на %v: %.4f", ts, il)
		}
	}
}

func TestFullMoonAround14DaysAfterNew(t *testing.T) {
	// Через ~14.77 діб після нової — повня, освітлення близьке до 1.
	full := refNewMoon.Add(time.Duration(synodicMonth / 2 * 24 * float64(time.Hour)))
	il := MoonIllumination(full)
	if math.Abs(il-1) > 0.05 {
		t.Fatalf("очікував майже повню, освітлення %.3f", il)
	}
	if MoonPhaseName(full) != "повня" {
		t.Fatalf("очікував 'повня', отримав %q", MoonPhaseName(full))
	}
}

func TestMoonConditionAlwaysHasHeadline(t *testing.T) {
	c := MoonCondition(time.Now())
	if c.Code != "moon_phase" || c.Headline == "" {
		t.Fatalf("неповна умова місяця: %+v", c)
	}
}
