package domain

import (
	"fmt"
	"strings"

	"smoker/pkg/smoke"
)

// Compose будує готовий до відправки текст українською (звичайний текст, без
// Markdown — у причинах трапляється довільний текст жартів, який зламав би
// розмітку). Gateway шле його як є.
func Compose(loc smoke.Location, rec smoke.Recommendation) string {
	var b strings.Builder

	if rec.Decision == smoke.Go {
		b.WriteString("🚬 ЧАС НА ПЕРЕКУР!\n")
	} else {
		b.WriteString("🤔 Краще перечекати з перекуром.\n")
	}

	place := loc.Name
	if place == "" {
		place = "ваша локація"
	}
	b.WriteString(fmt.Sprintf("📍 %s · %s\n\n", place, rec.Confidence))

	b.WriteString("Чому саме зараз:\n")
	for _, c := range rec.Reasons {
		b.WriteString("• ")
		b.WriteString(c.Headline)
		b.WriteString("\n")
	}

	b.WriteString(fmt.Sprintf("\n⚖️ Підсумковий індекс перекуру: %+d", rec.TotalScore))
	return b.String()
}
