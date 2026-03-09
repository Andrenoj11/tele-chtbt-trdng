package bot

import (
	"fmt"
	"strings"

	"crypto-telegram-bot/pkg/binance"
	"crypto-telegram-bot/pkg/strategy"
)

func HandleIncomingTextTelegram(chatID int64, text string) string {
	p := ParseMessage(text)

	if p.Command == "HELP" {
		return helpText()
	}

	if p.Command == "PRICE" {
		return handlePrice(p.Symbol)
	}

	if p.Command == "CHART" {
		return handleChart(p.Symbol, p.Interval)
	}

	// ANALYZE (default)
	klines, err := binance.GetKlines(p.Symbol, p.Interval, 200)
	if err != nil {
		return fmt.Sprintf("Error fetching market data for %s %s: %v", p.Symbol, p.Interval, err)
	}

	price, err := binance.GetPrice(p.Symbol)
	if err != nil {
		return fmt.Sprintf("Error fetching price for %s: %v", p.Symbol, err)
	}

	closes := make([]float64, 0, len(klines))
	for _, k := range klines {
		closes = append(closes, k.Close)
	}

	sig := strategy.BuildSignal(p.Symbol, p.Interval, closes, price)
	return FormatReply(sig)
}

func handlePrice(symbol string) string {
	t, err := binance.GetTicker24h(symbol)
	if err != nil {
		// fallback to last price only
		price, pErr := binance.GetPrice(symbol)
		if pErr != nil {
			return fmt.Sprintf("Error fetching price for %s: %v", symbol, err)
		}
		return fmt.Sprintf("%s\nPrice: %.10f", symbol, price)
	}

	// Price from ticker is string, keep as-is
	return fmt.Sprintf(
		"%s\nPrice: %s\n24h Change: %s%%",
		symbol,
		t.LastPrice,
		t.PriceChangePercent,
	)
}

func handleChart(symbol, interval string) string {
	klines, err := binance.GetKlines(symbol, interval, 200)
	if err != nil {
		return fmt.Sprintf("Error fetching klines for %s %s: %v", symbol, interval, err)
	}

	price, err := binance.GetPrice(symbol)
	if err != nil {
		return fmt.Sprintf("Error fetching price for %s: %v", symbol, err)
	}

	closes := make([]float64, 0, len(klines))
	for _, k := range klines {
		closes = append(closes, k.Close)
	}

	ema20 := strategy.EMA(closes, 20)
	ema50 := strategy.EMA(closes, 50)
	rsi14 := strategy.RSI(closes, 14)

	lastEma20 := ema20[len(ema20)-1]
	lastEma50 := ema50[len(ema50)-1]

	trend := "NEUTRAL"
	if lastEma20 > lastEma50 {
		trend = "UPTREND"
	} else if lastEma20 < lastEma50 {
		trend = "DOWNTREND"
	}

	lastRsiText := "n/a"
	if len(rsi14) > 0 {
		lastRsiText = fmt.Sprintf("%.1f", rsi14[len(rsi14)-1])
	}

	sig := strategy.BuildSignal(symbol, interval, closes, price)

	var b strings.Builder
	fmt.Fprintf(&b, "%s | TF: %s\n", symbol, interval)
	fmt.Fprintf(&b, "Price: %.10f\n", price)
	fmt.Fprintf(&b, "EMA20: %.10f\n", lastEma20)
	fmt.Fprintf(&b, "EMA50: %.10f\n", lastEma50)
	fmt.Fprintf(&b, "RSI14: %s\n", lastRsiText)
	fmt.Fprintf(&b, "Trend: %s\n\n", trend)
	fmt.Fprintf(&b, "Signal: %s (%s)", sig.Action, sig.Confidence)
	return b.String()
}

func helpText() string {
	return `Commands:
- /analyze BTCUSDT 15m
- /analyze ETHUSDT 1h
- /price BTCUSDT
- /chart BTCUSDT 15m
- /start or /help

Notes:
- Suggest-only (BUY/SELL/WAIT), not financial advice.
- Uses Binance SPOT market data (data-api.binance.vision).`
}
