package domain

import (
	"testing"

	"smoker/pkg/smoke"
)

func TestUVConditionVerdicts(t *testing.T) {
	cases := []struct {
		uv    float64
		want  smoke.Verdict
		score int
	}{
		{8, Favorable, 2},
		{4, Favorable, 1},
		{1.5, Neutral, 0},
		{0, Neutral, 0},
	}
	for _, c := range cases {
		got := UVCondition(c.uv)
		if got.Verdict != c.want || got.Score != c.score {
			t.Errorf("UV %.1f: отримав verdict=%v score=%d", c.uv, got.Verdict, got.Score)
		}
		if got.Headline == "" {
			t.Errorf("UV %.1f: порожній headline", c.uv)
		}
	}
}

func TestKpStormIsUnfavorable(t *testing.T) {
	c := KpCondition(6)
	if c.Verdict != Unfavorable || c.Score >= 0 {
		t.Fatalf("буря має бути несприятливою: %+v", c)
	}
}

func TestISSCloseIsFavorable(t *testing.T) {
	if ISSCondition(800).Verdict != Favorable {
		t.Fatal("МКС над головою має сприяти")
	}
	if ISSCondition(9000).Verdict != Neutral {
		t.Fatal("далека МКС має бути нейтральною")
	}
}

func TestWindRoseCardinal(t *testing.T) {
	if got := windRose(0); got != "північ" {
		t.Fatalf("0°: %s", got)
	}
	if got := windRose(90); got != "схід" {
		t.Fatalf("90°: %s", got)
	}
	if got := windRose(180); got != "південь" {
		t.Fatalf("180°: %s", got)
	}
}
