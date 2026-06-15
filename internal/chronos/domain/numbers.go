package domain

import "fmt"

// isPrime — просте число (для невеликих хвилин достатньо наївної перевірки).
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}

// isFibonacci — чи є n числом Фібоначчі (n у [0,59]).
func isFibonacci(n int) bool {
	a, b := 0, 1
	for a < n {
		a, b = b, a+b
	}
	return a == n
}

// MinuteNumerologyCondition — нумерологія поточної хвилини: просте чи Фібоначчі.
func MinuteNumerologyCondition(minute int) Condition {
	switch {
	case isFibonacci(minute) && minute >= 2:
		return Condition{Code: "minute_numerology", Verdict: Favorable, Score: 1,
			Headline: fmt.Sprintf("🌀 Хвилина %d — число Фібоначчі, навіть природа натякає вийти", minute)}
	case isPrime(minute):
		return Condition{Code: "minute_numerology", Verdict: Favorable, Score: 1,
			Headline: fmt.Sprintf("🔱 Хвилина %d — просте число, неподільне, як ваше право на перекур", minute)}
	default:
		return Condition{Code: "minute_numerology", Verdict: Neutral, Score: 0,
			Headline: fmt.Sprintf("🔢 Хвилина %d — складене число, ділиться на турботи, треба розвіятись", minute)}
	}
}
