// Package domain — core domain застосунку: SmokeAdvisor агрегує умови від усіх
// провайдерів у єдиний вердикт. Чиста логіка без мережі та proto.
package domain

import (
	"sort"

	"smoker/pkg/smoke"
)

// maxReasons — скільки причин показуємо в повідомленні.
const maxReasons = 4

// minReasons — мінімум причин, які намагаємось показати.
const minReasons = 3

// Decide — серце домену: підсумовує ваги умов і вирішує GO/WAIT.
// Поріг: TotalScore >= 0 → GO. Причини сортуються так, щоб найрелевантніші до
// рішення (сприятливі для GO / несприятливі для WAIT) йшли першими.
func Decide(conditions []smoke.Condition) smoke.Recommendation {
	total := 0
	for _, c := range conditions {
		total += c.Score
	}

	decision := smoke.Wait
	if total >= 0 {
		decision = smoke.Go
	}

	return smoke.Recommendation{
		Decision:   decision,
		TotalScore: total,
		Confidence: confidence(total),
		Reasons:    pickReasons(conditions, decision),
	}
}

// pickReasons обирає найвиразніші причини під рішення.
func pickReasons(conditions []smoke.Condition, decision smoke.Decision) []smoke.Condition {
	sorted := make([]smoke.Condition, len(conditions))
	copy(sorted, conditions)

	// "Вага релевантності": для GO вище — краще, для WAIT навпаки.
	weight := func(c smoke.Condition) int {
		if decision == smoke.Go {
			return c.Score
		}
		return -c.Score
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		return weight(sorted[i]) > weight(sorted[j])
	})

	n := maxReasons
	if n > len(sorted) {
		n = len(sorted)
	}
	// Якщо релевантних мало — все одно віддаємо хоча б minReasons (із нейтральних).
	if n < minReasons && len(sorted) >= minReasons {
		n = minReasons
	}
	return sorted[:n]
}

// confidence — людський опис впевненості за модулем підсумку.
func confidence(total int) string {
	switch a := abs(total); {
	case a >= 5:
		return "Космос наполягає"
	case a >= 3:
		return "Схоже на правду"
	case a >= 1:
		return "Є нюанси, але загалом так"
	default:
		return "На межі, вирішувати вам"
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
