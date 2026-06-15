package infra

import (
	"context"
	"time"

	"smoker/internal/cosmos/domain"
	"smoker/pkg/httpx"
	"smoker/pkg/smoke"
)

// ISSSource — поточна підсупутникова точка МКС (wheretheiss.at). Без ключа.
// Рахуємо відстань по поверхні від місця до точки під МКС.
type ISSSource struct {
	hc *httpx.Client
}

func NewISSSource(hc *httpx.Client) *ISSSource { return &ISSSource{hc: hc} }

func (s *ISSSource) Name() string { return "wheretheiss/iss" }

type issPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (s *ISSSource) Evaluate(ctx context.Context, loc smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	const url = "https://api.wheretheiss.at/v1/satellites/25544"
	var p issPosition
	if err := s.hc.GetJSON(ctx, url, &p); err != nil {
		return nil, err
	}
	dist := haversineKm(loc.Lat, loc.Lon, p.Latitude, p.Longitude)
	return []smoke.Condition{stamp(domain.ISSCondition(dist))}, nil
}
