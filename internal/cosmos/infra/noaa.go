package infra

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"smoker/internal/cosmos/domain"
	"smoker/pkg/httpx"
	"smoker/pkg/smoke"
)

// KpSource — NOAA SWPC planetary K-index (геомагнітна активність). Без ключа.
// Відповідь: масив рядків, перший — заголовок, далі дані; останній — найсвіжіший.
type KpSource struct {
	hc *httpx.Client
}

func NewKpSource(hc *httpx.Client) *KpSource { return &KpSource{hc: hc} }

func (s *KpSource) Name() string { return "noaa/planetary-k-index" }

func (s *KpSource) Evaluate(ctx context.Context, _ smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	const url = "https://services.swpc.noaa.gov/products/noaa-planetary-k-index.json"
	var rows [][]string
	if err := s.hc.GetJSON(ctx, url, &rows); err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("noaa: замало рядків (%d)", len(rows))
	}
	last := rows[len(rows)-1]
	if len(last) < 3 {
		return nil, fmt.Errorf("noaa: несподіваний формат рядка")
	}
	kp, err := strconv.ParseFloat(last[2], 64) // колонка kp_fraction
	if err != nil {
		return nil, fmt.Errorf("noaa: парсинг Kp: %w", err)
	}
	return []smoke.Condition{stamp(domain.KpCondition(kp))}, nil
}
