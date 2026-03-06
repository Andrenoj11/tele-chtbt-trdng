package strategy

import "fmt"

type Signal struct {
	Symbol     string
	Interval   string
	Price      float64
	Action     string // BUY/SELL/WAIT
	Confidence string // LOW/MEDIUM/HIGH
	Reasons    []string
	RiskNotes  []string
}

func BuildSignal(symbol, interval string, closes []float64, lastPrice float64) Signal {
	ema20 := EMA(closes, 20)
	ema50 := EMA(closes, 50)
	rsi14 := RSI(closes, 14)

	lastEma20 := ema20[len(ema20)-1]
	lastEma50 := ema50[len(ema50)-1]

	var lastRsi float64
	hasRsi := len(rsi14) > 0
	if hasRsi {
		lastRsi = rsi14[len(rsi14)-1]
	}

	reasons := []string{}
	risk := []string{
		"High volatility; use small size.",
		"Use stop-loss and max daily loss limit.",
		"This signal uses SPOT market data (data-api.binance.vision).",
	}

	uptrend := lastEma20 > lastEma50
	downtrend := lastEma20 < lastEma50

	if uptrend {
		reasons = append(reasons, "EMA20 > EMA50 (uptrend bias)")
	}
	if downtrend {
		reasons = append(reasons, "EMA20 < EMA50 (downtrend bias)")
	}
	if hasRsi {
		reasons = append(reasons, fmt.Sprintf("RSI(14) ≈ %.1f", lastRsi))
	}

	action := "WAIT"
	conf := "LOW"

	oversold := hasRsi && lastRsi <= 40
	overbought := hasRsi && lastRsi >= 60

	if uptrend && oversold {
		action, conf = "BUY", "HIGH"
		reasons = append(reasons, "Uptrend + RSI oversold → pullback entry signal")
	} else if downtrend && overbought {
		action, conf = "SELL", "HIGH"
		reasons = append(reasons, "Downtrend + RSI overbought → pullback sell/exit signal")
	} else if uptrend {
		action, conf = "WAIT", "MEDIUM"
		reasons = append(reasons, "Trend up, but no oversold confirmation")
	} else if downtrend {
		action, conf = "WAIT", "MEDIUM"
		reasons = append(reasons, "Trend down, but no overbought confirmation")
	} else {
		action, conf = "WAIT", "LOW"
		reasons = append(reasons, "No clear trend edge")
	}

	return Signal{
		Symbol:     symbol,
		Interval:   interval,
		Price:      lastPrice,
		Action:     action,
		Confidence: conf,
		Reasons:    reasons,
		RiskNotes:  risk,
	}
}
