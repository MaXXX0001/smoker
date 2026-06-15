// Package domain — бізнес-правила контексту "Природа та космос": як сирі числа
// зовнішніх API перетворюються на смішні Condition. Без мережі, чисті функції —
// усе тестується юніт-тестами.
package domain

import (
	"fmt"
	"math"
	"time"
)

// synodicMonth — середня тривалість місячного циклу (нов→нов), діб.
const synodicMonth = 29.530588853

// refNewMoon — відома нова місяць (UTC) як точка відліку.
var refNewMoon = time.Date(2000, time.January, 6, 18, 14, 0, 0, time.UTC)

// MoonAgeDays — вік місяця у добах [0, synodicMonth).
func MoonAgeDays(t time.Time) float64 {
	days := t.UTC().Sub(refNewMoon).Hours() / 24
	age := math.Mod(days, synodicMonth)
	if age < 0 {
		age += synodicMonth
	}
	return age
}

// MoonIllumination — частка освітленого диска [0,1].
func MoonIllumination(t time.Time) float64 {
	age := MoonAgeDays(t)
	return (1 - math.Cos(2*math.Pi*age/synodicMonth)) / 2
}

// MoonPhaseName — українська назва фази за віком.
func MoonPhaseName(t time.Time) string {
	age := MoonAgeDays(t)
	switch {
	case age < 1.85 || age >= synodicMonth-1.85:
		return "новий місяць"
	case age < 5.54:
		return "молодий місяць"
	case age < 9.23:
		return "перша чверть"
	case age < 12.91:
		return "прибуваючий місяць"
	case age < 16.61:
		return "повня"
	case age < 20.30:
		return "спадаючий місяць"
	case age < 23.99:
		return "остання чверть"
	default:
		return "старий місяць"
	}
}

// MoonCondition — фаза місяця як умова перекуру.
func MoonCondition(t time.Time) Condition {
	illum := MoonIllumination(t)
	phase := MoonPhaseName(t)
	pct := int(math.Round(illum * 100))

	switch {
	case phase == "повня":
		return Condition{
			Code: "moon_phase", Verdict: Favorable, Score: 2,
			Headline: fmt.Sprintf("🌕 Повня (%d%% освітлення) — час перевертнів і перекурів, інстинкт кличе надвір", pct),
		}
	case illum >= 0.5:
		return Condition{
			Code: "moon_phase", Verdict: Favorable, Score: 1,
			Headline: fmt.Sprintf("🌖 Місяць світить на %d%% (%s) — енергетика сприяє виходу", pct, phase),
		}
	case phase == "новий місяць":
		return Condition{
			Code: "moon_phase", Verdict: Neutral, Score: 0,
			Headline: "🌑 Новий місяць — небо порожнє, можна заповнити паузу димком",
		}
	default:
		return Condition{
			Code: "moon_phase", Verdict: Neutral, Score: 0,
			Headline: fmt.Sprintf("🌒 %s, освітлення %d%% — космос спостерігає без оцінок", phase, pct),
		}
	}
}
