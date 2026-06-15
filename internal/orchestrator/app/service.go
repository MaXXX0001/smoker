// Package app — застосунковий шар orchestrator: оркеструє виклики провайдерів
// (fan-out з graceful degradation) і застосовує core domain.
package app

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"smoker/internal/orchestrator/domain"
	"smoker/pkg/smoke"
)

// Provider — порт до одного provider-сервісу. Інтерфейс заради тестованості:
// у проді це gRPC-клієнт, у тестах — фейк.
type Provider interface {
	Name() string
	Evaluate(ctx context.Context, loc smoke.Location, t time.Time) ([]smoke.Condition, error)
}

// Service — use-case "порадити перекур".
type Service struct {
	providers []Provider
	log       *slog.Logger
	timeout   time.Duration
}

func NewService(log *slog.Logger, providers ...Provider) *Service {
	return &Service{providers: providers, log: log, timeout: 8 * time.Second}
}

// Result — вердикт + готовий текст.
type Result struct {
	Recommendation smoke.Recommendation
	Message        string
}

// Recommend збирає умови з усіх провайдерів і повертає рішення.
func (s *Service) Recommend(ctx context.Context, loc smoke.Location, t time.Time) Result {
	var (
		mu  sync.Mutex
		all []smoke.Condition
		wg  sync.WaitGroup
	)
	for _, p := range s.providers {
		wg.Add(1)
		go func(p Provider) {
			defer wg.Done()
			cctx, cancel := context.WithTimeout(ctx, s.timeout)
			defer cancel()
			conds, err := p.Evaluate(cctx, loc, t)
			if err != nil {
				s.log.Warn("провайдер недоступний", "provider", p.Name(), "err", err)
				return
			}
			mu.Lock()
			all = append(all, conds...)
			mu.Unlock()
		}(p)
	}
	wg.Wait()

	rec := domain.Decide(all)
	msg := domain.Compose(loc, rec)
	s.log.Info("рекомендація готова",
		"place", loc.Name, "decision", rec.Decision.String(),
		"score", rec.TotalScore, "conditions", len(all))
	return Result{Recommendation: rec, Message: msg}
}
