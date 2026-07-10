package domain

import (
	"testing"

	"smoker/pkg/smoke"
)

func cond(score int, v smoke.Verdict, head string) smoke.Condition {
	return smoke.Condition{Code: "x", Category: "c", Verdict: v, Score: score, Headline: head}
}

func TestDecideGoWhenPositive(t *testing.T) {
	rec := Decide([]smoke.Condition{
		cond(2, smoke.Favorable, "a"),
		cond(1, smoke.Favorable, "b"),
		cond(-1, smoke.Unfavorable, "c"),
	})
	if rec.Decision != smoke.Go {
		t.Fatalf("сума +2 → має бути GO, отримав %v", rec.Decision)
	}
	if rec.TotalScore != 2 {
		t.Fatalf("сума має бути 2, отримав %d", rec.TotalScore)
	}
}

func TestDecideWaitWhenNegative(t *testing.T) {
	rec := Decide([]smoke.Condition{
		cond(-2, smoke.Unfavorable, "storm"),
		cond(-1, smoke.Unfavorable, "pollen"),
		cond(1, smoke.Favorable, "sun"),
	})
	if rec.Decision != smoke.Wait {
		t.Fatalf("сума -2 → має бути WAIT, отримав %v", rec.Decision)
	}
}

func TestDecideZeroIsGo(t *testing.T) {
	rec := Decide([]smoke.Condition{cond(0, smoke.Neutral, "meh")})
	if rec.Decision != smoke.Go {
		t.Fatal("нуль має схилятись до GO")
	}
}

func TestPickReasonsAreFromInput(t *testing.T) {
	// Причини обираються випадково, але завжди мають бути підмножиною вхідних умов.
	in := []smoke.Condition{
		cond(0, smoke.Neutral, "neutral"),
		cond(2, smoke.Favorable, "best"),
		cond(1, smoke.Favorable, "good"),
	}
	valid := map[string]bool{"neutral": true, "best": true, "good": true}
	rec := Decide(in)
	if len(rec.Reasons) != len(in) { // 3 < maxReasons → повертаємо всі
		t.Fatalf("очікував %d причин, отримав %d", len(in), len(rec.Reasons))
	}
	for _, r := range rec.Reasons {
		if !valid[r.Headline] {
			t.Fatalf("причина %q не з вхідного набору", r.Headline)
		}
	}
}

func TestPickReasonsNoDuplicates(t *testing.T) {
	var cs []smoke.Condition
	for i := 0; i < 10; i++ {
		cs = append(cs, cond(1, smoke.Favorable, string(rune('a'+i))))
	}
	rec := Decide(cs)
	seen := map[string]bool{}
	for _, r := range rec.Reasons {
		if seen[r.Headline] {
			t.Fatalf("причина %q повторюється", r.Headline)
		}
		seen[r.Headline] = true
	}
}

func TestReasonsCappedAtMax(t *testing.T) {
	var cs []smoke.Condition
	for i := 0; i < 10; i++ {
		cs = append(cs, cond(1, smoke.Favorable, "x"))
	}
	if got := len(Decide(cs).Reasons); got != maxReasons {
		t.Fatalf("очікував %d причин, отримав %d", maxReasons, got)
	}
}
