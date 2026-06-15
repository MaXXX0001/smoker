package infra

import (
	"context"
	"fmt"
	"time"

	"smoker/internal/chronos/domain"
	"smoker/pkg/httpx"
	"smoker/pkg/smoke"
)

// HolidaySource — Nager.Date: найближчі державні свята країни. Без ключа.
// Код країни береться з конфігу (Nager оперує ISO-кодами, а не координатами).
type HolidaySource struct {
	hc      *httpx.Client
	country string
}

func NewHolidaySource(hc *httpx.Client, country string) *HolidaySource {
	return &HolidaySource{hc: hc, country: country}
}

func (s *HolidaySource) Name() string { return "nager/holidays" }

type nagerHoliday struct {
	Date      string `json:"date"`
	LocalName string `json:"localName"`
	Name      string `json:"name"`
}

func (s *HolidaySource) Evaluate(ctx context.Context, _ smoke.Location, t time.Time) ([]smoke.Condition, error) {
	url := fmt.Sprintf("https://date.nager.at/api/v3/NextPublicHolidays/%s", s.country)
	var hs []nagerHoliday
	if err := s.hc.GetJSON(ctx, url, &hs); err != nil {
		return nil, err
	}
	if len(hs) == 0 {
		return nil, fmt.Errorf("nager: порожній список свят")
	}
	next := hs[0]
	d, err := time.Parse("2006-01-02", next.Date)
	if err != nil {
		return nil, fmt.Errorf("nager: дата %q: %w", next.Date, err)
	}
	days := int(d.Sub(t.Truncate(24*time.Hour)).Hours() / 24)
	name := next.LocalName
	if name == "" {
		name = next.Name
	}
	return []smoke.Condition{stamp(domain.HolidayCondition(name, days))}, nil
}
