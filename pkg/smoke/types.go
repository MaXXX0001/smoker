// Package smoke — shared kernel (спільне ядро) у термінах DDD: мова, якою
// розмовляють усі bounded contexts. Тут чисті доменні типи без залежностей на
// proto чи мережу. Мапінг у/з gRPC лежить поряд у proto.go.
package smoke

// Verdict — внесок умови у рішення.
type Verdict int

const (
	Neutral     Verdict = iota // просто цікаво
	Favorable                  // сприяє перекуру
	Unfavorable                // проти
)

// Location — value object: де ми оцінюємо умови.
type Location struct {
	Lat  float64
	Lon  float64
	Name string
	TZ   string
}

// Condition — одна оцінена умова (value object). Незмінна за задумом:
// провайдери створюють її та віддають, ніхто не мутує.
type Condition struct {
	Code     string  // "uv_index"
	Category string  // "Природа та космос"
	Verdict  Verdict //
	Score    int     // вага у підсумку, типово -3..+3
	Headline string  // смішне пояснення українською
}

// Decision — підсумок core domain.
type Decision int

const (
	Wait Decision = iota
	Go
)

func (d Decision) String() string {
	if d == Go {
		return "GO"
	}
	return "WAIT"
}

// Recommendation — агрегат core domain: рішення + причини, які лягають у текст.
type Recommendation struct {
	Decision   Decision
	TotalScore int
	Confidence string
	Reasons    []Condition
}
