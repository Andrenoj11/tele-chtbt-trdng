package telegram

import (
	"encoding/json"
	"net/http"

	"crypto-telegram-bot/pkg/bot"
)

func WebhookHandler(b *Bot, expectedSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u Update
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chatID, text, ok := u.FirstText()
		if ok {
			reply := bot.HandleIncomingTextTelegram(chatID, text)
			_ = b.SendMessage(chatID, reply)
		}

		// always ACK Telegram quickly
		w.WriteHeader(http.StatusOK)
	}
}
