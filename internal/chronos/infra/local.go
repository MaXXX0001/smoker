// Package infra — джерела контексту chronos. Локальне джерело рахує час без
// мережі; решта ходить у Nager.Date та Wikimedia.
package infra

import (
	"context"
	"time"

	"smoker/internal/chronos/domain"
	"smoker/pkg/smoke"
)

func stamp(c domain.Condition) smoke.Condition {
	c.Category = domain.Category
	return c
}

// localTime повертає момент у часовому поясі місця (із фолбеком на UTC).
func localTime(loc smoke.Location, t time.Time) time.Time {
	if loc.TZ != "" {
		if z, err := time.LoadLocation(loc.TZ); err == nil {
			return t.In(z)
		}
	}
	return t.UTC()
}

// LocalSource — усі суто-часові умови. Ніколи не падає.
type LocalSource struct{}

func NewLocalSource() *LocalSource { return &LocalSource{} }

func (s *LocalSource) Name() string { return "local/time" }

func (s *LocalSource) Evaluate(_ context.Context, loc smoke.Location, t time.Time) ([]smoke.Condition, error) {
	lt := localTime(loc, t)
	return []smoke.Condition{
		stamp(domain.DayProgressCondition(lt)),
		stamp(domain.YearProgressCondition(lt)),
		stamp(domain.SpecialClockCondition(lt)),
		stamp(domain.MinuteNumerologyCondition(lt.Minute())),
	}, nil
}
