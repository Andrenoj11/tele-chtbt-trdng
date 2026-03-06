package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Bot struct {
	token      string
	httpClient *http.Client
}

func NewBot(token string) *Bot {
	return &Bot{
		token: token,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (b *Bot) apiURL(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.token, method)
}

func (b *Bot) SetWebhook(webhookURL string, secretToken string) error {
	body := map[string]any{
		"url": webhookURL,
	}
	if secretToken != "" {
		body["secret_token"] = secretToken
	}

	bb, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, b.apiURL("setWebhook"), bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")

	res, err := b.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("setWebhook status=%d", res.StatusCode)
	}
	return nil
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	body := map[string]any{
		"chat_id": chatID,
		"text":    text,
	}

	bb, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, b.apiURL("sendMessage"), bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")

	res, err := b.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("sendMessage status=%d", res.StatusCode)
	}
	return nil
}
