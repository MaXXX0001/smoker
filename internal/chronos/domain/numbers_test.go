package domain

import "testing"

func TestIsPrime(t *testing.T) {
	primes := map[int]bool{0: false, 1: false, 2: true, 3: true, 4: false, 17: true, 51: false, 59: true}
	for n, want := range primes {
		if isPrime(n) != want {
			t.Errorf("isPrime(%d)=%v, want %v", n, isPrime(n), want)
		}
	}
}

func TestIsFibonacci(t *testing.T) {
	fib := map[int]bool{0: true, 1: true, 2: true, 3: true, 5: true, 8: true, 13: true, 21: true, 34: true, 55: true, 4: false, 50: false}
	for n, want := range fib {
		if isFibonacci(n) != want {
			t.Errorf("isFibonacci(%d)=%v, want %v", n, isFibonacci(n), want)
		}
	}
}

func TestMinuteNumerologyAlwaysHeadline(t *testing.T) {
	for m := 0; m < 60; m++ {
		c := MinuteNumerologyCondition(m)
		if c.Headline == "" || c.Code != "minute_numerology" {
			t.Fatalf("хвилина %d: неповна умова %+v", m, c)
		}
	}
}

func TestSpecialClockLeet(t *testing.T) {
	c := SpecialClockCondition(mkTime(13, 37))
	if c.Score != 2 || c.Verdict != Favorable {
		t.Fatalf("13:37 має бути топовим: %+v", c)
	}
}
