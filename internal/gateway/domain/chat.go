// Package domain — доменна модель контексту gateway: зареєстрований чат і
// правила, коли йому час слати нагадування. Без залежностей на БД/Telegram.
package domain

import (
	"time"

	"smoker/pkg/smoke"
)

// Chat — агрегат: налаштування групи, куди бот пише.
type Chat struct {
	ChatID          int64
	Lat             float64
	Lon             float64
	PlaceName       string
	TZ              string
	IntervalMinutes int
	Enabled         bool
	LastSent        time.Time
}

// Location повертає доменну локацію чату для запиту порад.
func (c Chat) Location() smoke.Location {
	return smoke.Location{Lat: c.Lat, Lon: c.Lon, Name: c.PlaceName, TZ: c.TZ}
}

// HasLocation — чи задано координати (0,0 вважаємо "не задано").
func (c Chat) HasLocation() bool {
	return c.Lat != 0 || c.Lon != 0
}

// Due вирішує, чи час слати цьому чату нагадування зараз.
// Правила: чат увімкнено, ми в робочих годинах (за TZ чату) і від останнього
// повідомлення минув заданий інтервал.
func (c Chat) Due(now time.Time, workStartHour, workEndHour int) bool {
	if !c.Enabled || c.IntervalMinutes <= 0 {
		return false
	}
	local := now.UTC()
	if c.TZ != "" {
		if z, err := time.LoadLocation(c.TZ); err == nil {
			local = now.In(z)
		}
	}
	if h := local.Hour(); h < workStartHour || h >= workEndHour {
		return false
	}
	return now.Sub(c.LastSent) >= time.Duration(c.IntervalMinutes)*time.Minute
}
