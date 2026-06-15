// Package domain — бізнес-правила контексту "Час і числа": усе, що виводиться
// з самого моменту часу без жодних API. Чисті функції, повністю тестовані.
package domain

import (
	"fmt"
	"time"

	"smoker/pkg/smoke"
)

// Category — назва bounded context.
const Category = "Час і числа"

type Condition = smoke.Condition

const (
	Favorable   = smoke.Favorable
	Neutral     = smoke.Neutral
	Unfavorable = smoke.Unfavorable
)

// DayProgressCondition — скільки відсотків доби вже минуло (у локальному часі).
// 42% — окрема пасхалка для шанувальників Дугласа Адамса.
func DayProgressCondition(local time.Time) Condition {
	secs := local.Hour()*3600 + local.Minute()*60 + local.Second()
	pct := float64(secs) / 86400 * 100

	switch {
	case pct >= 41.5 && pct < 42.5:
		return Condition{Code: "day_progress", Verdict: Favorable, Score: 2,
			Headline: "🔢 Минуло рівно 42% доби — відповідь на головне питання Всесвіту, час діяти"}
	case pct >= 50 && pct < 51:
		return Condition{Code: "day_progress", Verdict: Favorable, Score: 1,
			Headline: "⏳ Рівно половина доби позаду — символічний рубіж варто відзначити перекуром"}
	default:
		return Condition{Code: "day_progress", Verdict: Neutral, Score: 0,
			Headline: fmt.Sprintf("⏳ Минуло %.0f%% доби — кожна секунда наближає до заслуженої паузи", pct)}
	}
}

// YearProgressCondition — скільки відсотків року минуло.
func YearProgressCondition(local time.Time) Condition {
	yearStart := time.Date(local.Year(), 1, 1, 0, 0, 0, 0, local.Location())
	yearEnd := yearStart.AddDate(1, 0, 0)
	pct := float64(local.Sub(yearStart)) / float64(yearEnd.Sub(yearStart)) * 100
	return Condition{Code: "year_progress", Verdict: Neutral, Score: 0,
		Headline: fmt.Sprintf("📅 Рік пройдено на %.0f%% — час летить, цигарка чекати не буде", pct)}
}

// SpecialClockCondition — "красиві" значення годинника (13:37, 11:11, ...).
func SpecialClockCondition(local time.Time) Condition {
	h, m := local.Hour(), local.Minute()
	hm := fmt.Sprintf("%02d:%02d", h, m)
	switch hm {
	case "13:37":
		return Condition{Code: "special_clock", Verdict: Favorable, Score: 2,
			Headline: "🕐 Час 13:37 — елітний leet-час, справжні профі курять саме зараз"}
	case "11:11":
		return Condition{Code: "special_clock", Verdict: Favorable, Score: 2,
			Headline: "🕚 11:11 — загадуйте бажання і виходьте, Всесвіт слухає"}
	case "12:34":
		return Condition{Code: "special_clock", Verdict: Favorable, Score: 1,
			Headline: "🕧 12:34 — цифри по порядку, ідеальний момент для порядку в голові"}
	case "00:00":
		return Condition{Code: "special_clock", Verdict: Favorable, Score: 1,
			Headline: "🕛 Опівніч рівно — містична година, дим розчиняється в темряві красиво"}
	}
	if h == m {
		return Condition{Code: "special_clock", Verdict: Favorable, Score: 1,
			Headline: fmt.Sprintf("🕰️ %s — година дорівнює хвилині, рідкісна симетрія, гріх не вийти", hm)}
	}
	return Condition{Code: "special_clock", Verdict: Neutral, Score: 0,
		Headline: fmt.Sprintf("🕰️ Зараз %s — звичайний час, але хіба перекур колись буває не на часі?", hm)}
}
