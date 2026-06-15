// Package infra — джерела контексту chaos: випадкові жарт-API + локальний кубик.
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

// DadJokeSource — icanhazdadjoke.com (вимагає Accept: application/json).
type DadJokeSource struct{ hc *httpx.Client }

func NewDadJokeSource() *DadJokeSource {
	return &DadJokeSource{hc: httpx.New(httpx.WithHeader("Accept", "application/json"))}
}

func (s *DadJokeSource) Name() string { return "icanhazdadjoke" }

type dadJoke struct {
	Joke string `json:"joke"`
}

func (s *DadJokeSource) Evaluate(ctx context.Context, _ smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	var r dadJoke
	if err := s.hc.GetJSON(ctx, "https://icanhazdadjoke.com/", &r); err != nil {
		return nil, err
	}
	return []smoke.Condition{stamp(domain.JokeCondition(r.Joke))}, nil
}

// ChuckSource — api.chucknorris.io.
type ChuckSource struct{ hc *httpx.Client }

func NewChuckSource(hc *httpx.Client) *ChuckSource { return &ChuckSource{hc: hc} }

func (s *ChuckSource) Name() string { return "chucknorris.io" }

type chuckFact struct {
	Value string `json:"value"`
}

func (s *ChuckSource) Evaluate(ctx context.Context, _ smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	var r chuckFact
	if err := s.hc.GetJSON(ctx, "https://api.chucknorris.io/jokes/random", &r); err != nil {
		return nil, err
	}
	return []smoke.Condition{stamp(domain.ChuckCondition(r.Value))}, nil
}

// CatFactSource — catfact.ninja.
type CatFactSource struct{ hc *httpx.Client }

func NewCatFactSource(hc *httpx.Client) *CatFactSource { return &CatFactSource{hc: hc} }

func (s *CatFactSource) Name() string { return "catfact.ninja" }

type catFact struct {
	Fact string `json:"fact"`
}

func (s *CatFactSource) Evaluate(ctx context.Context, _ smoke.Location, _ time.Time) ([]smoke.Condition, error) {
	var r catFact
	if err := s.hc.GetJSON(ctx, "https://catfact.ninja/fact", &r); err != nil {
		return nil, err
	}
	return []smoke.Condition{stamp(domain.CatFactCondition(r.Fact))}, nil
}
