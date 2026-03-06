package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Kline struct {
	OpenTime  int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	CloseTime int64
}

func apiBase() string {
	if v := os.Getenv("BINANCE_API_BASE"); v != "" {
		return v
	}
	// default to data-api (works in your network)
	return "https://data-api.binance.vision"
}

func GetKlines(symbol, interval string, limit int) ([]Kline, error) {
	url := fmt.Sprintf("%s/api/v3/klines?symbol=%s&interval=%s&limit=%d", apiBase(), symbol, interval, limit)

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("klines status=%d", res.StatusCode)
	}

	var raw [][]any
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	out := make([]Kline, 0, len(raw))
	for _, r := range raw {
		out = append(out, Kline{
			OpenTime:  toInt64(r[0]),
			Open:      toFloat(r[1]),
			High:      toFloat(r[2]),
			Low:       toFloat(r[3]),
			Close:     toFloat(r[4]),
			Volume:    toFloat(r[5]),
			CloseTime: toInt64(r[6]),
		})
	}
	return out, nil
}

func GetMarkPrice(symbol string) (float64, error) {
	// spot price
	url := fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", apiBase(), symbol)

	client := &http.Client{Timeout: 20 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return 0, fmt.Errorf("price status=%d", res.StatusCode)
	}

	var data struct {
		Price string `json:"price"`
	}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(data.Price, 64)
}

func toFloat(v any) float64 {
	switch t := v.(type) {
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case float64:
		return t
	default:
		return 0
	}
}

func toInt64(v any) int64 {
	switch t := v.(type) {
	case float64:
		return int64(t)
	case string:
		i, _ := strconv.ParseInt(t, 10, 64)
		return i
	default:
		return 0
	}
}
