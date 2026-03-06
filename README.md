# crypto-telegram-bot (Go) — Telegram + Binance Futures Signal Bot

MVP Telegram bot in Go that:
- Polls Telegram `getUpdates` (long polling)
- Parses user text (e.g. `ANALYZE BTCUSDT 15m`)
- Fetches Binance USDⓈ-M Futures market data (klines + mark price)
- Computes EMA/RSI and returns BUY/SELL/WAIT suggestions

## Env

Create `.env` (or set env vars in your shell):

```env
TELEGRAM_BOT_TOKEN=123456:ABCDEF...

# Binance Futures (public endpoints don't require keys, but kept for later trading)
BINANCE_API_KEY=
BINANCE_API_SECRET=
BINANCE_FAPI_BASE=https://demo-fapi.binance.com  # demo first; production: https://fapi.binance.com
```

## Run

```bash
go mod tidy
go run .
```

## Commands

- `/analyze BTCUSDT 15m`
- `/analyze ETHUSDT 1h`
- `HELP`

## Notes
- This bot is suggest-only (no auto trading).
- Start with demo futures endpoints.
