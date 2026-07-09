// Command gateway — публікаційний контекст: Telegram-бот, реєстрація чатів,
// планувальник нагадувань. Єдиний сервіс, що дивиться у зовнішній світ.
package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"smoker/internal/gateway/app"
	"smoker/internal/gateway/infra"
	"smoker/internal/gateway/store"
	"smoker/pkg/env"
	"smoker/pkg/httpx"

	"github.com/go-telegram/bot"
)

// parseIDs розбирає "123,456" у список Telegram user ID; порожнє = nil (усі).
func parseIDs(s string) []int64 {
	var out []int64
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p == "" {
			continue
		}
		if id, err := strconv.ParseInt(p, 10, 64); err == nil {
			out = append(out, id)
		}
	}
	return out
}

func main() {
	log := env.Logger("gateway")

	token := env.String("TELEGRAM_TOKEN", "")
	if token == "" {
		log.Error("TELEGRAM_TOKEN не заданий")
		os.Exit(1)
	}
	dbPath := env.String("DB_PATH", "smoker.db")
	orchAddr := env.String("ORCHESTRATOR_ADDR", "localhost:9100")
	adminIDs := parseIDs(env.String("ADMIN_IDS", ""))

	defaults := app.Defaults{
		Lat:       env.Float("DEFAULT_LAT", 50.4501),
		Lon:       env.Float("DEFAULT_LON", 30.5234),
		Place:     env.String("DEFAULT_PLACE", "Київ"),
		TZ:        env.String("DEFAULT_TZ", "Europe/Kyiv"),
		Interval:  int(env.Float("DEFAULT_INTERVAL_MIN", 90)),
		WorkStart: int(env.Float("WORK_START_HOUR", 9)),
		WorkEnd:   int(env.Float("WORK_END_HOUR", 19)),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	st, err := store.Open(ctx, dbPath)
	if err != nil {
		log.Error("БД недоступна", "err", err)
		os.Exit(1)
	}
	defer st.Close()
	log.Info("БД підключено")

	advisor, err := infra.DialAdvisor(orchAddr)
	if err != nil {
		log.Error("orchestrator недоступний", "err", err)
		os.Exit(1)
	}
	defer advisor.Close()

	geo := infra.NewGeocoder(httpx.New(), env.String("WIKI_LANG", "uk"))

	application := &app.App{
		Store:    st,
		Geo:      geo,
		Advisor:  advisor,
		Log:      log,
		Def:      defaults,
		AdminIDs: adminIDs,
	}

	b, err := bot.New(token)
	if err != nil {
		log.Error("не вдалось створити бота", "err", err)
		os.Exit(1)
	}
	application.Register(b)

	go application.RunScheduler(ctx, b)

	log.Info("бот запущено (long polling)")
	b.Start(ctx) // блокується до скасування ctx
	log.Info("бот зупинено")
}
