// Package infra — зовнішні адаптери gateway: геокодер і клієнт до orchestrator.
package infra

import (
	"context"
	"fmt"
	"net/url"

	"smoker/pkg/httpx"
)

// Place — результат геокодування назви міста.
type Place struct {
	Name string
	Lat  float64
	Lon  float64
	TZ   string
}

// Geocoder перетворює назву міста на координати+TZ через Open-Meteo Geocoding
// API (без ключа).
type Geocoder struct {
	hc   *httpx.Client
	lang string
}

func NewGeocoder(hc *httpx.Client, lang string) *Geocoder {
	return &Geocoder{hc: hc, lang: lang}
}

type geoResp struct {
	Results []struct {
		Name      string  `json:"name"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timezone  string  `json:"timezone"`
	} `json:"results"`
}

// Lookup шукає перше місто за назвою.
func (g *Geocoder) Lookup(ctx context.Context, name string) (Place, error) {
	u := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=%s&format=json",
		url.QueryEscape(name), g.lang)
	var r geoResp
	if err := g.hc.GetJSON(ctx, u, &r); err != nil {
		return Place{}, err
	}
	if len(r.Results) == 0 {
		return Place{}, fmt.Errorf("місто %q не знайдено", name)
	}
	res := r.Results[0]
	label := res.Name
	if res.Country != "" {
		label = res.Name + ", " + res.Country
	}
	return Place{Name: label, Lat: res.Latitude, Lon: res.Longitude, TZ: res.Timezone}, nil
}
