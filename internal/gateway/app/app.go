// Package app — застосунковий шар gateway: команди Telegram і планувальник.
package app

import (
	"log/slog"
	"smoker/internal/gateway/infra"
	"smoker/internal/gateway/store"
)

// Defaults — дефолтні налаштування для нових чатів і вікно робочих годин.
type Defaults struct {
	Lat       float64
	Lon       float64
	Place     string
	TZ        string
	Interval  int // хвилини
	WorkStart int // година (локальна)
	WorkEnd   int
}

// App — спільні залежності хендлерів і планувальника.
type App struct {
	Store    *store.Store
	Geo      *infra.Geocoder
	Advisor  *infra.AdvisorClient
	Log      *slog.Logger
	Def      Defaults
	AdminIDs []int64 // хто може керувати ботом; порожньо = всі
}
