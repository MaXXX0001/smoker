package domain

import "fmt"

// HolidayCondition — наскільки близько найближче свято.
func HolidayCondition(name string, daysUntil int) Condition {
	switch {
	case daysUntil <= 0:
		return Condition{Code: "holiday", Verdict: Favorable, Score: 2,
			Headline: fmt.Sprintf("🎉 Сьогодні ж %s! Святковий перекур обов'язковий", name)}
	case daysUntil <= 3:
		return Condition{Code: "holiday", Verdict: Favorable, Score: 1,
			Headline: fmt.Sprintf("🎊 До свята «%s» лишилось %d дн. — треба ж тренувати святковий настрій", name, daysUntil)}
	default:
		return Condition{Code: "holiday", Verdict: Neutral, Score: 0,
			Headline: fmt.Sprintf("📆 До найближчого свята («%s») аж %d дн. — наблизьмо його перекуром", name, daysUntil)}
	}
}

// OnThisDayCondition — історична подія "N років тому сьогодні".
func OnThisDayCondition(yearsAgo int, text string) Condition {
	return Condition{Code: "on_this_day", Verdict: Neutral, Score: 0,
		Headline: fmt.Sprintf("📜 Цього дня %d р. тому: %s. Історія робила паузу — зробіть і ви", yearsAgo, text)}
}
