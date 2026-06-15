package app

import (
	"context"
	"time"

	"github.com/go-telegram/bot"
)

// RunScheduler — фоновий цикл: щохвилини перевіряє активні чати і шле тим, кому
// час (за їхнім інтервалом і робочими годинами). Блокується до скасування ctx.
func (a *App) RunScheduler(ctx context.Context, b *bot.Bot) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	a.Log.Info("планувальник запущено",
		"work_hours", []int{a.Def.WorkStart, a.Def.WorkEnd})

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.tick(ctx, b, time.Now())
		}
	}
}

func (a *App) tick(ctx context.Context, b *bot.Bot, now time.Time) {
	chats, err := a.Store.ListEnabled(ctx)
	if err != nil {
		a.Log.Warn("список чатів", "err", err)
		return
	}
	for _, chat := range chats {
		if !chat.Due(now, a.Def.WorkStart, a.Def.WorkEnd) {
			continue
		}
		loc := chat.Location()
		if !chat.HasLocation() {
			loc.Lat, loc.Lon, loc.Name, loc.TZ = a.Def.Lat, a.Def.Lon, a.Def.Place, a.Def.TZ
		}

		advice, err := a.Advisor.Recommend(ctx, loc, now)
		if err != nil {
			a.Log.Warn("advisor у планувальнику", "chat", chat.ChatID, "err", err)
			continue
		}
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chat.ChatID, Text: advice.Message}); err != nil {
			a.Log.Warn("надсилання за розкладом", "chat", chat.ChatID, "err", err)
			continue
		}
		if err := a.Store.MarkSent(ctx, chat.ChatID, now); err != nil {
			a.Log.Warn("mark sent", "chat", chat.ChatID, "err", err)
		}
		a.Log.Info("надіслано за розкладом", "chat", chat.ChatID, "go", advice.GoNow)
	}
}
