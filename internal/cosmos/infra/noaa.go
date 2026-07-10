package infra

import (
	"context"
	"fmt"
	"time"

	"smoker/internal/cosmos/domain"
	"smoker/pkg/httpx"
	"smoker/pkg/smoke"
)

// KpSource — NOAA SWPC planetary K-index (геомагнітна активність). Без ключа.
// Відповідь: масив об'єктів {time_tag, Kp, ...}; останній — найсвіжіший.
type KpSource struct {
	hc *httpx.Client
}

func NewKpSource(hc *httpx.Client) *KpSource { return &KpSource{hc: hc} }

func (s *KpSource) Name() string { return "noaa/planetary-k-index" }

type kpRow struct {
	Kp float64 `json:"Kp"`
}

func (s *KpSource) Evaluate(ctx context.Context, _ smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	const url = "https://services.swpc.noaa.gov/products/noaa-planetary-k-index.json"
	var rows []kpRow
	if err := s.hc.GetJSON(ctx, url, &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("noaa: порожня відповідь")
	}
	kp := rows[len(rows)-1].Kp
	return []smoke.Condition{stamp(domain.KpCondition(kp))}, nil
}
