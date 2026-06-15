package domain

import (
	"fmt"
	"math"

	"smoker/pkg/smoke"
)

// Category — назва bounded context для всіх умов cosmos.
const Category = "Природа та космос"

// Аліаси на shared kernel, щоб домен читався у власних термінах.
type Condition = smoke.Condition

const (
	Favorable   = smoke.Favorable
	Neutral     = smoke.Neutral
	Unfavorable = smoke.Unfavorable
)

// UVCondition — UV-індекс як "безкоштовний солярій".
func UVCondition(uv float64) Condition {
	switch {
	case uv >= 6:
		return Condition{Code: "uv_index", Verdict: Favorable, Score: 2,
			Headline: fmt.Sprintf("☀️ UV-індекс %.0f — встигнеш засмагнути за одну цигарку, безкоштовний солярій", uv)}
	case uv >= 3:
		return Condition{Code: "uv_index", Verdict: Favorable, Score: 1,
			Headline: fmt.Sprintf("🌤️ UV-індекс %.0f — сонце є, вітамін D сам себе не виробить", uv)}
	case uv >= 1:
		return Condition{Code: "uv_index", Verdict: Neutral, Score: 0,
			Headline: fmt.Sprintf("🌥️ UV-індекс %.0f — мляве сонце, але промінь надії пробивається", uv)}
	default:
		return Condition{Code: "uv_index", Verdict: Neutral, Score: 0,
			Headline: "🌑 UV майже нуль — шкіра в безпеці, можна курити без сонцезахисту"}
	}
}

// windRose — напрямок у компасну стрілку (укр.).
func windRose(deg float64) string {
	dirs := []string{"північ", "північний схід", "схід", "південний схід",
		"південь", "південний захід", "захід", "північний захід"}
	i := int(math.Mod(deg/45+0.5, 8))
	if i < 0 {
		i += 8
	}
	return dirs[i]
}

// WindCondition — напрямок/сила вітру: куди понесе дим.
func WindCondition(speed, deg float64) Condition {
	rose := windRose(deg)
	switch {
	case speed >= 12:
		return Condition{Code: "wind", Verdict: Unfavorable, Score: -1,
			Headline: fmt.Sprintf("💨 Вітер %.0f м/с на %s — запальничку доведеться прикривати всім тілом", speed, rose)}
	case speed >= 3:
		return Condition{Code: "wind", Verdict: Favorable, Score: 1,
			Headline: fmt.Sprintf("💨 Легкий вітер %.0f м/с дме на %s — дим елегантно понесе геть від обличчя", speed, rose)}
	default:
		return Condition{Code: "wind", Verdict: Neutral, Score: 0,
			Headline: "🍃 Штиль — дим висітиме навколо вас урочистою аурою"}
	}
}

// PressureCondition — атмосферний тиск як "медичне показання".
func PressureCondition(hPa float64) Condition {
	switch {
	case hPa < 1005:
		return Condition{Code: "pressure", Verdict: Favorable, Score: 1,
			Headline: fmt.Sprintf("🌧️ Тиск %.0f гПа — падає, тіло вимагає компенсації нікотином, майже рецепт", hPa)}
	case hPa > 1028:
		return Condition{Code: "pressure", Verdict: Neutral, Score: 0,
			Headline: fmt.Sprintf("🗻 Тиск %.0f гПа — височенний, голова ясна, рішення про перекур зважене", hPa)}
	default:
		return Condition{Code: "pressure", Verdict: Neutral, Score: 0,
			Headline: fmt.Sprintf("🌡️ Тиск %.0f гПа — норма, організм не проти", hPa)}
	}
}

// KpCondition — геомагнітна активність.
func KpCondition(kp float64) Condition {
	switch {
	case kp >= 5:
		return Condition{Code: "geomagnetic_kp", Verdict: Unfavorable, Score: -2,
			Headline: fmt.Sprintf("🧲 Магнітна буря Kp=%.0f — мозок і так не варить, краще перечекати 20 хв (або вийти на зло)", kp)}
	case kp >= 3:
		return Condition{Code: "geomagnetic_kp", Verdict: Neutral, Score: 0,
			Headline: fmt.Sprintf("🧲 Магнітне поле нервує (Kp=%.0f) — компас бреше, а перекур — ні", kp)}
	default:
		return Condition{Code: "geomagnetic_kp", Verdict: Favorable, Score: 1,
			Headline: fmt.Sprintf("🧲 Геомагніт спокійний (Kp=%.0f) — космос не заперечує", kp)}
	}
}

// ISSCondition — наскільки МКС близько до вашого неба.
func ISSCondition(distanceKm float64) Condition {
	switch {
	case distanceKm < 1500:
		return Condition{Code: "iss_overhead", Verdict: Favorable, Score: 2,
			Headline: fmt.Sprintf("🛰️ МКС просто зараз за %.0f км над вами — космонавти дивляться, вийдіть гідно", distanceKm)}
	case distanceKm < 4000:
		return Condition{Code: "iss_overhead", Verdict: Favorable, Score: 1,
			Headline: fmt.Sprintf("🛰️ МКС наближається (%.0f км) — встигнете помахати з двору", distanceKm)}
	default:
		return Condition{Code: "iss_overhead", Verdict: Neutral, Score: 0,
			Headline: fmt.Sprintf("🛰️ МКС аж за %.0f км — далеко, але вона думає про вас", distanceKm)}
	}
}

// PollenCondition — рівень пилку (макс. з берези/трав/вільхи), частинок/м³.
func PollenCondition(grains float64) Condition {
	switch {
	case grains >= 50:
		return Condition{Code: "pollen", Verdict: Unfavorable, Score: -1,
			Headline: fmt.Sprintf("🤧 Пилок зашкалює (%.0f/м³) — виходьте гуртом, разом чхати веселіше", grains)}
	case grains >= 10:
		return Condition{Code: "pollen", Verdict: Neutral, Score: 0,
			Headline: fmt.Sprintf("🌼 Пилок помірний (%.0f/м³) — антигістамінне і вперед", grains)}
	default:
		return Condition{Code: "pollen", Verdict: Favorable, Score: 1,
			Headline: "🌸 Пилку майже нема — дихайте на повні груди (поки що)"}
	}
}
