// Package provider — перевикористовуваний gRPC-адаптер для всіх provider-
// сервісів (cosmos / chronos / chaos). Доменна логіка кожного сервісу живе у
// наборі Source; цей пакет лише запускає їх паралельно, збирає Condition та
// реалізує спільний контракт ConditionProvider. Якщо джерело впало —
// пропускаємо його (graceful degradation), вердикт порахується з решти.
package provider

import (
	"context"
	"log/slog"
	"sync"
	"time"

	conditionv1 "smoker/pkg/proto/condition/v1"
	"smoker/pkg/smoke"
)

// Source — одне джерело умов усередині bounded context. Це і є точка
// розширення домену: додати нову смішну умову = додати Source.
type Source interface {
	// Name — для логів.
	Name() string
	// Evaluate повертає 0..N умов для місця/моменту. Помилка не фатальна.
	Evaluate(ctx context.Context, loc smoke.Location, t time.Time) ([]smoke.Condition, error)
}

// Server реалізує conditionv1.ConditionProviderServer поверх набору Source.
type Server struct {
	conditionv1.UnimplementedConditionProviderServer
	sources []Source
	log     *slog.Logger
	// timeout на одне джерело, щоб повільний API не тримав усю відповідь.
	timeout time.Duration
}

// NewServer збирає адаптер із джерел.
func NewServer(log *slog.Logger, sources ...Source) *Server {
	return &Server{sources: sources, log: log, timeout: 7 * time.Second}
}

// Evaluate — реалізація контракту: паралельний fan-out по джерелах.
func (s *Server) Evaluate(ctx context.Context, req *conditionv1.EvaluateRequest) (*conditionv1.EvaluateResponse, error) {
	loc := smoke.LocationFromProto(req.GetLocation())
	t := time.Unix(req.GetUnixTs(), 0)
	if req.GetUnixTs() == 0 {
		t = time.Now()
	}

	var (
		mu  sync.Mutex
		all []smoke.Condition
		wg  sync.WaitGroup
	)
	for _, src := range s.sources {
		wg.Add(1)
		go func(src Source) {
			defer wg.Done()
			cctx, cancel := context.WithTimeout(ctx, s.timeout)
			defer cancel()
			conds, err := src.Evaluate(cctx, loc, t)
			if err != nil {
				s.log.Warn("джерело пропущено", "source", src.Name(), "err", err)
				return
			}
			mu.Lock()
			all = append(all, conds...)
			mu.Unlock()
		}(src)
	}
	wg.Wait()

	return &conditionv1.EvaluateResponse{Conditions: smoke.ConditionsToProto(all)}, nil
}
