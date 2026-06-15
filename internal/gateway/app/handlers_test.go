package app

import "testing"

func TestCommandArg(t *testing.T) {
	cases := map[string]string{
		"/setlocation Київ":         "Київ",
		"/setlocation@bot Львів":    "Львів",
		"/smoke":                    "",
		"  /setschedule   2h  ":     "2h",
		"/setlocation 50.45, 30.52": "50.45, 30.52",
	}
	for in, want := range cases {
		if got := commandArg(in); got != want {
			t.Errorf("commandArg(%q)=%q, want %q", in, got, want)
		}
	}
}

func TestParseCoords(t *testing.T) {
	lat, lon, ok := parseCoords("50.45, 30.52")
	if !ok || lat != 50.45 || lon != 30.52 {
		t.Fatalf("очікував (50.45,30.52,true), отримав (%v,%v,%v)", lat, lon, ok)
	}
	if _, _, ok := parseCoords("Київ"); ok {
		t.Fatal("назва міста не має парситись як координати")
	}
}

func TestParseInterval(t *testing.T) {
	cases := []struct {
		in   string
		want int
		ok   bool
	}{
		{"90", 90, true},
		{"30m", 30, true},
		{"2h", 120, true},
		{"3", 0, false}, // < 5 хв
		{"abc", 0, false},
		{"", 0, false},
	}
	for _, c := range cases {
		got, err := parseInterval(c.in)
		if c.ok && (err != nil || got != c.want) {
			t.Errorf("parseInterval(%q)=(%d,%v), want %d", c.in, got, err, c.want)
		}
		if !c.ok && err == nil {
			t.Errorf("parseInterval(%q) мав повернути помилку", c.in)
		}
	}
}
