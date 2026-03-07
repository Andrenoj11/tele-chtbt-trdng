# crypto-telegram-bot (Go) — Telegram Signal Bot (SPOT data)

Telegram bot in Go that:
- Receives updates via **Telegram Webhook**
- Parses user commands (e.g. `/analyze BTCUSDT 15m`)
- Fetches **Binance SPOT** market data via `data-api.binance.vision`
- Computes EMA/RSI and returns **BUY / SELL / WAIT** suggestions

> Suggest-only (no auto trading). Not financial advice.

---

## Features
- `/analyze <SYMBOL> <TF>`: analyze pair + timeframe
- `/help` or `/start`: usage guide
- Common supported intervals: `1m 3m 5m 15m 30m 1h 2h 4h 6h 12h 1d 3d 1w`
  - Note: month `1M` needs a small parser tweak if you want it.

---

## Environment Variables

### Required
```env
TELEGRAM_BOT_TOKEN=123456:ABCDEF...
```

### Recommended (Market data base)
```env
# Use this base if api.binance.com / fapi.binance.com is blocked/reset in your network
BINANCE_API_BASE=https://data-api.binance.vision
```

### Optional (reserved for future trading features)
```env
BINANCE_API_KEY=
BINANCE_API_SECRET=
```

---

## Run locally (Webhook mode)

### 1) Install dependencies
```bash
go mod tidy
```

### 2) Run server
```bash
go run .
```

### 3) Expose webhook using ngrok
```bash
ngrok http 3000
```

Copy the HTTPS URL from ngrok, then set Telegram webhook to:
`https://<NGROK_DOMAIN>/api/telegram/webhook`

```bash
curl -s "https://api.telegram.org/bot<TELEGRAM_BOT_TOKEN>/setWebhook" \
  -d "url=https://<NGROK_DOMAIN>/api/telegram/webhook"
```

Check webhook status:
```bash
curl -s "https://api.telegram.org/bot<TELEGRAM_BOT_TOKEN>/getWebhookInfo"
```

---

## Commands & Examples

### Analyze
```text
/analyze BTCUSDT 15m
/analyze ETHUSDT 1h
/analyze SHIBUSDT 1m
/analyze PEPEUSDT 3m
```

### Help
```text
/help
/start
```

---

## Notes / Troubleshooting

### If bot doesn’t respond
1) Check webhook points to the correct URL:
```bash
curl -s "https://api.telegram.org/bot<TELEGRAM_BOT_TOKEN>/getWebhookInfo"
```
2) Ensure your webhook URL is publicly reachable (HTTPS) and returns `200 OK`.

### About the signal
Current strategy is a simple EMA(20/50) + RSI(14) “trend + pullback” style rule.
It will often return WAIT when RSI is not extreme (30/70).

---

## Disclaimer
This project is for learning. Signals can be wrong. Use demo/paper trading and risk limits.
