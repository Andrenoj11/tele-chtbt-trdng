package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"crypto-telegram-bot/internal/telegram"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("WEBHOOK_URL is required (must be HTTPS public)")
	}

	secret := os.Getenv("TELEGRAM_SECRET_TOKEN") // optional

	b := telegram.NewBot(token)

	// set webhook once on startup
	if err := b.SetWebhook(webhookURL, secret); err != nil {
		log.Fatal("setWebhook failed:", err)
	}
	log.Println("webhook set:", webhookURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/telegram/webhook", telegram.WebhookHandler(b, secret))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      20 * time.Second,
	}

	log.Println("listening on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
