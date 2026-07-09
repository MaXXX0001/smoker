package app

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const helpText = `🚬 Привіт! Я підкажу, коли вашій команді час на перекур — і завжди поясню, чому саме зараз (UV-індекс, фаза місяця, МКС над головою, магнітні бурі, прості числа на годиннику та інша наукова дичина).

Команди:
/smoke — перевірити прямо зараз
/setlocation <місто> — задати локацію (напр. /setlocation Київ)
/setschedule <хвилини|30m|2h> — як часто нагадувати
/stop — припинити нагадування у цьому чаті
/start — це повідомлення`

// Register чіпляє всі команди до бота.
func (a *App) Register(b *bot.Bot) {
	// Патерн — без провідного "/": MatchTypeCommand порівнює голе ім'я команди
	// (бібліотека зрізає слеш і "@botname" сама).
	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, a.handleStart)
	b.RegisterHandler(bot.HandlerTypeMessageText, "smoke", bot.MatchTypeCommand, a.handleSmoke)
	b.RegisterHandler(bot.HandlerTypeMessageText, "setlocation", bot.MatchTypeCommand, a.handleSetLocation)
	b.RegisterHandler(bot.HandlerTypeMessageText, "setschedule", bot.MatchTypeCommand, a.handleSetSchedule)
	b.RegisterHandler(bot.HandlerTypeMessageText, "stop", bot.MatchTypeCommand, a.handleStop)
}

func (a *App) reply(ctx context.Context, b *bot.Bot, chatID int64, text string) {
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: text}); err != nil {
		a.Log.Warn("не вдалось надіслати", "chat", chatID, "err", err)
	}
}

// ensure гарантує, що чат є в БД із дефолтами.
func (a *App) ensure(ctx context.Context, chatID int64) error {
	return a.Store.Ensure(ctx, chatID, a.Def.Lat, a.Def.Lon, a.Def.Place, a.Def.TZ, a.Def.Interval)
}

func (a *App) handleStart(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.Message == nil {
		return
	}
	if err := a.ensure(ctx, u.Message.Chat.ID); err != nil {
		a.Log.Warn("ensure", "err", err)
	}
	a.reply(ctx, b, u.Message.Chat.ID, helpText)
}

func (a *App) handleSmoke(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.Message == nil {
		return
	}
	chatID := u.Message.Chat.ID
	if err := a.ensure(ctx, chatID); err != nil {
		a.Log.Warn("ensure", "err", err)
	}
	chat, err := a.Store.Get(ctx, chatID)
	if err != nil {
		a.reply(ctx, b, chatID, "Щось пішло не так із базою 😬")
		return
	}
	loc := chat.Location()
	if !chat.HasLocation() {
		loc.Lat, loc.Lon, loc.Name, loc.TZ = a.Def.Lat, a.Def.Lon, a.Def.Place, a.Def.TZ
	}

	advice, err := a.Advisor.Recommend(ctx, loc, time.Now())
	if err != nil {
		a.Log.Warn("advisor", "err", err)
		a.reply(ctx, b, chatID, "Радник у відпустці (orchestrator недоступний) 🛠️")
		return
	}
	a.reply(ctx, b, chatID, advice.Message)
	_ = a.Store.MarkSent(ctx, chatID, time.Now())
}

func (a *App) handleSetLocation(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.Message == nil {
		return
	}
	chatID := u.Message.Chat.ID
	arg := commandArg(u.Message.Text)
	if arg == "" {
		a.reply(ctx, b, chatID, "Вкажіть місто: /setlocation Київ (або координати: /setlocation 50.45,30.52)")
		return
	}
	if err := a.ensure(ctx, chatID); err != nil {
		a.Log.Warn("ensure", "err", err)
	}

	// Спершу пробуємо як координати "lat,lon".
	if lat, lon, ok := parseCoords(arg); ok {
		if err := a.Store.SetLocation(ctx, chatID, lat, lon, fmt.Sprintf("%.3f,%.3f", lat, lon), a.Def.TZ); err != nil {
			a.reply(ctx, b, chatID, "Не зміг зберегти локацію 😬")
			return
		}
		a.reply(ctx, b, chatID, fmt.Sprintf("📍 Готово: %.3f, %.3f", lat, lon))
		return
	}

	place, err := a.Geo.Lookup(ctx, arg)
	if err != nil {
		a.reply(ctx, b, chatID, fmt.Sprintf("Не знайшов «%s» 🗺️ Спробуйте іншу назву або координати.", arg))
		return
	}
	if err := a.Store.SetLocation(ctx, chatID, place.Lat, place.Lon, place.Name, place.TZ); err != nil {
		a.reply(ctx, b, chatID, "Не зміг зберегти локацію 😬")
		return
	}
	a.reply(ctx, b, chatID, fmt.Sprintf("📍 Локацію встановлено: %s (часовий пояс %s)", place.Name, place.TZ))
}

func (a *App) handleSetSchedule(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.Message == nil {
		return
	}
	chatID := u.Message.Chat.ID
	minutes, err := parseInterval(commandArg(u.Message.Text))
	if err != nil {
		a.reply(ctx, b, chatID, "Формат: /setschedule 90  (або 30m, 2h). Мінімум 5 хв.")
		return
	}
	if err := a.ensure(ctx, chatID); err != nil {
		a.Log.Warn("ensure", "err", err)
	}
	if err := a.Store.SetInterval(ctx, chatID, minutes); err != nil {
		a.reply(ctx, b, chatID, "Не зміг зберегти розклад 😬")
		return
	}
	a.reply(ctx, b, chatID, fmt.Sprintf("⏰ Нагадуватиму кожні %d хв (у робочі години %d:00–%d:00).",
		minutes, a.Def.WorkStart, a.Def.WorkEnd))
}

func (a *App) handleStop(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.Message == nil {
		return
	}
	chatID := u.Message.Chat.ID
	if err := a.ensure(ctx, chatID); err != nil {
		a.Log.Warn("ensure", "err", err)
	}
	if err := a.Store.SetEnabled(ctx, chatID, false); err != nil {
		a.reply(ctx, b, chatID, "Не зміг вимкнути 😬")
		return
	}
	a.reply(ctx, b, chatID, "🛑 Окей, мовчу. Увімкнути знову — /setschedule або /start.")
}

// --- парсинг аргументів ---

// commandArg повертає текст після команди (прибирає "/cmd" і "/cmd@bot").
func commandArg(text string) string {
	parts := strings.SplitN(strings.TrimSpace(text), " ", 2)
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func parseCoords(s string) (lat, lon float64, ok bool) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return 0, 0, false
	}
	var err1, err2 error
	lat, err1 = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	lon, err2 = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return lat, lon, true
}

// parseInterval приймає "90", "30m", "2h" → хвилини (мінімум 5).
func parseInterval(s string) (int, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0, fmt.Errorf("порожньо")
	}
	var minutes int
	switch {
	case strings.HasSuffix(s, "h"):
		h, err := strconv.Atoi(strings.TrimSuffix(s, "h"))
		if err != nil {
			return 0, err
		}
		minutes = h * 60
	case strings.HasSuffix(s, "m"):
		m, err := strconv.Atoi(strings.TrimSuffix(s, "m"))
		if err != nil {
			return 0, err
		}
		minutes = m
	default:
		m, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		minutes = m
	}
	if minutes < 5 {
		return 0, fmt.Errorf("замало")
	}
	return minutes, nil
}
