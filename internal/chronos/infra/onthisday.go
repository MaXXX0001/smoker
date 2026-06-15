package infra

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"smoker/internal/chronos/domain"
	"smoker/pkg/httpx"
	"smoker/pkg/smoke"
)

// OnThisDaySource — історичні події цього дня з Wikipedia REST feed (keyless).
// Мова налаштовується (uk/en).
type OnThisDaySource struct {
	hc   *httpx.Client
	lang string
}

func NewOnThisDaySource(hc *httpx.Client, lang string) *OnThisDaySource {
	return &OnThisDaySource{hc: hc, lang: lang}
}

func (s *OnThisDaySource) Name() string { return "wikipedia/onthisday" }

type wikiOnThisDay struct {
	Events []struct {
		Year int    `json:"year"`
		Text string `json:"text"`
	} `json:"events"`
}

func (s *OnThisDaySource) Evaluate(ctx context.Context, _ smoke.Location, t time.Time) ([]smoke.Condition, error) {
	url := fmt.Sprintf("https://%s.wikipedia.org/api/rest_v1/feed/onthisday/events/%02d/%02d",
		s.lang, int(t.Month()), t.Day())
	var r wikiOnThisDay
	if err := s.hc.GetJSON(ctx, url, &r); err != nil {
		return nil, err
	}
	if len(r.Events) == 0 {
		return nil, fmt.Errorf("wikipedia: подій не знайдено")
	}
	e := r.Events[rand.Intn(len(r.Events))]
	yearsAgo := t.Year() - e.Year
	if yearsAgo < 0 {
		yearsAgo = 0
	}
	return []smoke.Condition{stamp(domain.OnThisDayCondition(yearsAgo, e.Text))}, nil
}
