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

func TestPickReasonsOrdersByRelevance(t *testing.T) {
	// Для GO найсприятливіша умова має йти першою.
	rec := Decide([]smoke.Condition{
		cond(0, smoke.Neutral, "neutral"),
		cond(2, smoke.Favorable, "best"),
		cond(1, smoke.Favorable, "good"),
	})
	if rec.Reasons[0].Headline != "best" {
		t.Fatalf("перша причина для GO має бути 'best', отримав %q", rec.Reasons[0].Headline)
	}
}

func TestPickReasonsForWaitSurfacesNegatives(t *testing.T) {
	rec := Decide([]smoke.Condition{
		cond(1, smoke.Favorable, "sun"),
		cond(-2, smoke.Unfavorable, "storm"),
		cond(0, smoke.Neutral, "meh"),
	})
	if rec.Reasons[0].Headline != "storm" {
		t.Fatalf("для WAIT першою має бути 'storm', отримав %q", rec.Reasons[0].Headline)
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
