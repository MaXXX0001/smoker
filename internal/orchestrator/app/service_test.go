package app

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"smoker/pkg/smoke"
)

// fakeProvider — тестовий дублер app.Provider.
type fakeProvider struct {
	name  string
	conds []smoke.Condition
	err   error
}

func (f fakeProvider) Name() string { return f.name }
func (f fakeProvider) Evaluate(context.Context, smoke.Location, time.Time) ([]smoke.Condition, error) {
	return f.conds, f.err
}

func quietLog() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestRecommendAggregatesAcrossProviders(t *testing.T) {
	svc := NewService(quietLog(),
		fakeProvider{name: "cosmos", conds: []smoke.Condition{
			{Code: "uv", Verdict: smoke.Favorable, Score: 2, Headline: "сонце"},
		}},
		fakeProvider{name: "chaos", conds: []smoke.Condition{
			{Code: "dice", Verdict: smoke.Favorable, Score: 1, Headline: "шістка"},
		}},
	)
	res := svc.Recommend(context.Background(), smoke.Location{Name: "Київ"}, time.Now())

	if res.Recommendation.Decision != smoke.Go {
		t.Fatalf("сума +3 → GO, отримав %v", res.Recommendation.Decision)
	}
	if res.Recommendation.TotalScore != 3 {
		t.Fatalf("очікував score 3, отримав %d", res.Recommendation.TotalScore)
	}
	if !strings.Contains(res.Message, "Київ") || !strings.Contains(res.Message, "ПЕРЕКУР") {
		t.Fatalf("повідомлення без очікуваного тексту:\n%s", res.Message)
	}
}

// Ключова перевірка: якщо один провайдер упав — вердикт рахується з решти.
func TestRecommendDegradesOnProviderError(t *testing.T) {
	svc := NewService(quietLog(),
		fakeProvider{name: "cosmos", err: errors.New("API лежить")},
		fakeProvider{name: "chronos", conds: []smoke.Condition{
			{Code: "clock", Verdict: smoke.Favorable, Score: 1, Headline: "13:37"},
		}},
	)
	res := svc.Recommend(context.Background(), smoke.Location{Name: "Львів"}, time.Now())

	if res.Recommendation.TotalScore != 1 {
		t.Fatalf("очікував score 1 від уцілілого провайдера, отримав %d", res.Recommendation.TotalScore)
	}
	if len(res.Recommendation.Reasons) == 0 {
		t.Fatal("мали б лишитись причини попри падіння одного провайдера")
	}
}
