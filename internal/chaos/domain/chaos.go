// Package domain — бізнес-правила контексту "Випадковість заради сміху": бере
// випадкові факти/жарти та подає їх як псевдонаукове обґрунтування перекуру.
package domain

import (
	"fmt"

	"smoker/pkg/smoke"
)

const Category = "Випадковість заради сміху"

type Condition = smoke.Condition

const (
	Favorable   = smoke.Favorable
	Neutral     = smoke.Neutral
	Unfavorable = smoke.Unfavorable
)

// DiceCondition — кидок d6 вирішує долю.
func DiceCondition(roll int) Condition {
	switch {
	case roll == 6:
		return Condition{Code: "dice", Verdict: Favorable, Score: 2,
			Headline: "🎲 Випало 6! Кубик долі однозначний — негайно на перекур"}
	case roll >= 4:
		return Condition{Code: "dice", Verdict: Favorable, Score: 1,
			Headline: fmt.Sprintf("🎲 Кубик показав %d — фортуна радше «за»", roll)}
	case roll == 1:
		return Condition{Code: "dice", Verdict: Neutral, Score: 0,
			Headline: "🎲 Випала 1 — кубик скептичний, але хто його слухає"}
	default:
		return Condition{Code: "dice", Verdict: Neutral, Score: 0,
			Headline: fmt.Sprintf("🎲 Кубик видав %d — ні так ні сяк, вирішуйте серцем", roll)}
	}
}

// JokeCondition — dad joke як "офіційне обґрунтування".
func JokeCondition(joke string) Condition {
	return Condition{Code: "dad_joke", Verdict: Favorable, Score: 1,
		Headline: fmt.Sprintf("🃏 Наукове обґрунтування дня: «%s» — після такого треба вийти провітритись", joke)}
}

// OracleCondition — вердикт «оракула» yesno.wtf: так/ні на питання «чи час?».
func OracleCondition(yes bool) Condition {
	if yes {
		return Condition{Code: "oracle", Verdict: Favorable, Score: 1,
			Headline: "🔮 Оракул на питання «час курити?» видав ТАК — доля не проти, вставайте"}
	}
	return Condition{Code: "oracle", Verdict: Unfavorable, Score: -1,
		Headline: "🔮 Оракул на питання «час курити?» видав НІ — але з нами не сперечаються, вирішуйте самі"}
}

// CatFactCondition — коти знають толк у відпочинку (спін уже вшитий у сам факт).
func CatFactCondition(fact string) Condition {
	return Condition{Code: "cat_fact", Verdict: Neutral, Score: 0,
		Headline: fmt.Sprintf("🐱 Котофакт: %s", fact)}
}
