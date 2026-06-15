// Package infra — anti-corruption layer контексту cosmos: ходить у зовнішні API
// та перекладає їхні відповіді у доменні Condition. Кожен тип тут реалізує
// provider.Source.
package infra

import (
	"context"
	"time"

	"smoker/internal/cosmos/domain"
	"smoker/pkg/httpx"
	"smoker/pkg/smoke"
)

// stamp проставляє категорію bounded context на доменну умову.
func stamp(c domain.Condition) smoke.Condition {
	c.Category = domain.Category
	return c
}

// WeatherSource — Open-Meteo Forecast API: UV, вітер, тиск. Без ключа.
type WeatherSource struct {
	hc *httpx.Client
}

func NewWeatherSource(hc *httpx.Client) *WeatherSource { return &WeatherSource{hc: hc} }

func (s *WeatherSource) Name() string { return "open-meteo/weather" }

type omForecast struct {
	Current struct {
		UVIndex       float64 `json:"uv_index"`
		SurfacePress  float64 `json:"surface_pressure"`
		WindSpeed     float64 `json:"wind_speed_10m"`
		WindDirection float64 `json:"wind_direction_10m"`
	} `json:"current"`
}

func (s *WeatherSource) Evaluate(ctx context.Context, loc smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	url := "https://api.open-meteo.com/v1/forecast" +
		"?latitude=" + ftoa(loc.Lat) + "&longitude=" + ftoa(loc.Lon) +
		"&current=uv_index,surface_pressure,wind_speed_10m,wind_direction_10m" +
		"&wind_speed_unit=ms"

	var r omForecast
	if err := s.hc.GetJSON(ctx, url, &r); err != nil {
		return nil, err
	}
	return []smoke.Condition{
		stamp(domain.UVCondition(r.Current.UVIndex)),
		stamp(domain.WindCondition(r.Current.WindSpeed, r.Current.WindDirection)),
		stamp(domain.PressureCondition(r.Current.SurfacePress)),
	}, nil
}

// PollenSource — Open-Meteo Air-Quality API: пилок. Дані лише по Європі; поза
// зоною поля приходять null → беремо 0 і умова стає сприятливою.
type PollenSource struct {
	hc *httpx.Client
}

func NewPollenSource(hc *httpx.Client) *PollenSource { return &PollenSource{hc: hc} }

func (s *PollenSource) Name() string { return "open-meteo/air-quality" }

type omAirQuality struct {
	Current struct {
		Birch *float64 `json:"birch_pollen"`
		Grass *float64 `json:"grass_pollen"`
		Alder *float64 `json:"alder_pollen"`
	} `json:"current"`
}

func (s *PollenSource) Evaluate(ctx context.Context, loc smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	url := "https://air-quality-api.open-meteo.com/v1/air-quality" +
		"?latitude=" + ftoa(loc.Lat) + "&longitude=" + ftoa(loc.Lon) +
		"&current=birch_pollen,grass_pollen,alder_pollen"

	var r omAirQuality
	if err := s.hc.GetJSON(ctx, url, &r); err != nil {
		return nil, err
	}
	maxGrains := maxPtr(r.Current.Birch, r.Current.Grass, r.Current.Alder)
	return []smoke.Condition{stamp(domain.PollenCondition(maxGrains))}, nil
}

func maxPtr(vals ...*float64) float64 {
	m := 0.0
	for _, v := range vals {
		if v != nil && *v > m {
			m = *v
		}
	}
	return m
}
