package infra

import (
	"context"
	"time"

	"smoker/internal/cosmos/domain"
	"smoker/pkg/smoke"
)

// MoonSource — фаза місяця. Єдине локальне джерело cosmos: рахується астрономічно,
// без мережі, тому ніколи не падає.
type MoonSource struct{}

func NewMoonSource() *MoonSource { return &MoonSource{} }

func (s *MoonSource) Name() string { return "local/moon-phase" }

func (s *MoonSource) Evaluate(_ context.Context, _ smoke.Location, t time.Time) ([]smoke.Condition, error) {
	return []smoke.Condition{stamp(domain.MoonCondition(t))}, nil
}
