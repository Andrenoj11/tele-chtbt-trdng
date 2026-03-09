package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Ticker24h struct {
	Symbol             string `json:"symbol"`
	LastPrice          string `json:"lastPrice"`
	PriceChangePercent string `json:"priceChangePercent"`
}

func GetTicker24h(symbol string) (Ticker24h, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/24hr?symbol=%s", apiBase(), symbol)

	client := &http.Client{Timeout: 20 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return Ticker24h{}, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return Ticker24h{}, fmt.Errorf("ticker24h status=%d", res.StatusCode)
	}

	var t Ticker24h
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return Ticker24h{}, err
	}
	return t, nil
}
