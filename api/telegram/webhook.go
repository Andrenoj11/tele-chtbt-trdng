package main

import (
	"encoding/json"
	"net/http"
	"os"

	"crypto-telegram-bot/pkg/bot"
	"crypto-telegram-bot/pkg/telegram"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("TELEGRAM_SECRET_TOKEN")
	if secret != "" {
		got := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
		if got != secret {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	var u telegram.Update
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, text, ok := u.FirstText()
	if ok {
		reply := bot.HandleIncomingTextTelegram(chatID, text)

		tok := os.Getenv("TELEGRAM_BOT_TOKEN")
		b := telegram.NewBot(tok)
		_ = b.SendMessage(chatID, reply)
	}

	w.WriteHeader(http.StatusOK)
}
