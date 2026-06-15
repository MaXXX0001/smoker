package domain

import "testing"

func TestDiceSixIsBest(t *testing.T) {
	c := DiceCondition(6)
	if c.Verdict != Favorable || c.Score != 2 {
		t.Fatalf("шістка має бути топовою: %+v", c)
	}
}

func TestDiceRange(t *testing.T) {
	for r := 1; r <= 6; r++ {
		if DiceCondition(r).Headline == "" {
			t.Fatalf("кидок %d без headline", r)
		}
	}
}

func TestJokeWrapping(t *testing.T) {
	c := JokeCondition("чому програміст не виходить надвір")
	if c.Code != "dad_joke" || c.Verdict != Favorable {
		t.Fatalf("жарт має бути сприятливим обґрунтуванням: %+v", c)
	}
}
