package bot

import (
	"fmt"

	"crypto-telegram-bot/pkg/binance"
	"crypto-telegram-bot/pkg/strategy"
)

func HandleIncomingTextTelegram(chatID int64, text string) string {
	p := ParseMessage(text)

	if p.Command == "HELP" {
		return helpText()
	}

	klines, err := binance.GetKlines(p.Symbol, p.Interval, 200)
	if err != nil {
		return fmt.Sprintf("Error fetching market data for %s %s: %v", p.Symbol, p.Interval, err)
	}

	price, err := binance.GetMarkPrice(p.Symbol)
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

func helpText() string {
	return `Commands:
- /analyze BTCUSDT 15m
- /analyze ETHUSDT 1h
- HELP

Notes:
- Suggest-only (BUY/SELL/WAIT), not financial advice.
- Uses Binance USDⓈ-M Futures market data (demo by default).`
}
