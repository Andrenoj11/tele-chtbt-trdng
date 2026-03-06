package telegram

import (
	"encoding/json"
	"net/http"

	"crypto-telegram-bot/pkg/bot"
)

func WebhookHandler(b *Bot, expectedSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Telegram will send this header if you set secret_token
		if expectedSecret != "" {
			got := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
			if got != expectedSecret {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

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
