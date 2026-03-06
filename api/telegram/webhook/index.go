package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"crypto-telegram-bot/pkg/bot"
	"crypto-telegram-bot/pkg/telegram"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// allow GET for quick test
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
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
		if tok != "" {
			b := telegram.NewBot(tok)
			_ = b.SendMessage(chatID, reply)
		}
	}

	w.WriteHeader(http.StatusOK)
}
