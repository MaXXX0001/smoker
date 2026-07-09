// Package infra — джерела контексту chaos: локальні укр-набори + оракул yesno.wtf.
package infra

import (
	"context"
	"math/rand"
	"time"

	"smoker/internal/chaos/domain"
	"smoker/pkg/httpx"
	"smoker/pkg/smoke"
)

func stamp(c domain.Condition) smoke.Condition {
	c.Category = domain.Category
	return c
}

// DiceSource — локальний кидок d6. Ніколи не падає.
type DiceSource struct{}

func NewDiceSource() *DiceSource { return &DiceSource{} }

func (s *DiceSource) Name() string { return "local/dice" }

func (s *DiceSource) Evaluate(_ context.Context, _ smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	return []smoke.Condition{stamp(domain.DiceCondition(rand.Intn(6) + 1))}, nil
}

// JokeSource — випадкове «наукове обґрунтування» з локального набору.
type JokeSource struct{}

func NewJokeSource() *JokeSource { return &JokeSource{} }

func (s *JokeSource) Name() string { return "local/jokes" }

func (s *JokeSource) Evaluate(_ context.Context, _ smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	return []smoke.Condition{stamp(domain.JokeCondition(jokes[rand.Intn(len(jokes))]))}, nil
}

// CatFactSource — випадковий котофакт з локального набору.
type CatFactSource struct{}

func NewCatFactSource() *CatFactSource { return &CatFactSource{} }

func (s *CatFactSource) Name() string { return "local/catfacts" }

func (s *CatFactSource) Evaluate(_ context.Context, _ smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	return []smoke.Condition{stamp(domain.CatFactCondition(catFacts[rand.Intn(len(catFacts))]))}, nil
}

// OracleSource — yesno.wtf: мовно-нейтральний вердикт "yes"/"no".
type OracleSource struct{ hc *httpx.Client }

func NewOracleSource(hc *httpx.Client) *OracleSource { return &OracleSource{hc: hc} }

func (s *OracleSource) Name() string { return "yesno.wtf" }

type oracleAnswer struct {
	Answer string `json:"answer"`
}

func (s *OracleSource) Evaluate(ctx context.Context, _ smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	var r oracleAnswer
	if err := s.hc.GetJSON(ctx, "https://yesno.wtf/api", &r); err != nil {
		return nil, err
	}
	return []smoke.Condition{stamp(domain.OracleCondition(r.Answer == "yes"))}, nil
}
