package app

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// onlyAdmin — middleware: пропускає команду лише від дозволених Telegram user ID.
// Якщо AdminIDs порожній — пускає всіх, але логує ID відправника (щоб легко
// дізнатись свій і потім звузити доступ через ADMIN_IDS).
func (a *App) onlyAdmin(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, u *models.Update) {
		if u.Message == nil || u.Message.From == nil {
			return
		}
		from := u.Message.From
		if len(a.AdminIDs) == 0 {
			a.Log.Info("команда (allowlist порожній — пускаю всіх)",
				"user_id", from.ID, "username", from.Username)
			next(ctx, b, u)
			return
		}
		if !a.isAdmin(from.ID) {
			a.Log.Debug("команда відхилена: не в allowlist",
				"user_id", from.ID, "username", from.Username)
			return
		}
		next(ctx, b, u)
	}
}

func (a *App) isAdmin(id int64) bool {
	for _, x := range a.AdminIDs {
		if x == id {
			return true
		}
	}
	return false
}
