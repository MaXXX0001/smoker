package infra

import (
	_ "embed"
	"strings"
)

// Контент chaos тримаємо в data/*.txt (по рядку на фразу), а не в коді:
// правити набори може будь-хто без дотику до логіки, а go:embed вшиває їх
// у бінар на етапі компіляції — жодних файлів у рантаймі.

//go:embed data/jokes.txt
var jokesRaw string

//go:embed data/catfacts.txt
var catFactsRaw string

var (
	jokes    = nonEmptyLines(jokesRaw)
	catFacts = nonEmptyLines(catFactsRaw)
)

// nonEmptyLines розбиває вбудований текст на рядки, відкидаючи порожні
// та зайві пробіли по краях.
func nonEmptyLines(s string) []string {
	raw := strings.Split(s, "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		if line = strings.TrimSpace(line); line != "" {
			out = append(out, line)
		}
	}
	return out
}
