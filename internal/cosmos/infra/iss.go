package infra

import (
	"context"
	"strconv"
	"time"

	"smoker/internal/cosmos/domain"
	"smoker/pkg/httpx"
	"smoker/pkg/smoke"
)

// ISSSource — поточна підсупутникова точка МКС (open-notify). Без ключа.
// wheretheiss.at був стабільно ~12-14s (не влазив у 7s таймаут провайдера),
// open-notify відповідає за ~0.5s. HTTP (не HTTPS) — у сервісу кривий cert,
// дані публічні, секретів нема.
type ISSSource struct {
	hc *httpx.Client
}

func NewISSSource(hc *httpx.Client) *ISSSource { return &ISSSource{hc: hc} }

func (s *ISSSource) Name() string { return "open-notify/iss" }

// open-notify віддає координати рядками у вкладеному об'єкті iss_position.
type issResponse struct {
	Position struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"iss_position"`
}

func (s *ISSSource) Evaluate(ctx context.Context, loc smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	const url = "http://api.open-notify.org/iss-now.json"
	var r issResponse
	if err := s.hc.GetJSON(ctx, url, &r); err != nil {
		return nil, err
	}
	lat, err := strconv.ParseFloat(r.Position.Latitude, 64)
	if err != nil {
		return nil, err
	}
	lon, err := strconv.ParseFloat(r.Position.Longitude, 64)
	if err != nil {
		return nil, err
	}
	dist := haversineKm(loc.Lat, loc.Lon, lat, lon)
	return []smoke.Condition{stamp(domain.ISSCondition(dist))}, nil
}
