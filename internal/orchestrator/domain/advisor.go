// Package domain — core domain застосунку: SmokeAdvisor агрегує умови від усіх
// провайдерів у єдиний вердикт. Чиста логіка без мережі та proto.
package domain

import (
	"math/rand"

	"smoker/pkg/smoke"
)

// maxReasons — скільки причин показуємо в повідомленні.
const maxReasons = 4

// Decide — серце домену: підсумовує ваги умов і вирішує GO/WAIT.
// Поріг: TotalScore >= 0 → GO. Причини для показу обираються випадково —
// щоб повідомлення щоразу було різне (вердикт рахується з усіх умов).
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
		Reasons:    pickReasons(conditions),
	}
}

// pickReasons бере до maxReasons ВИПАДКОВИХ умов — щоб причини щоразу різнились.
func pickReasons(conditions []smoke.Condition) []smoke.Condition {
	shuffled := make([]smoke.Condition, len(conditions))
	copy(shuffled, conditions)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	n := maxReasons
	if n > len(shuffled) {
		n = len(shuffled)
	}
	return shuffled[:n]
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
