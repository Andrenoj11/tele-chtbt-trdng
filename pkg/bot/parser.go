package bot

import (
	"regexp"
	"strings"
)

var (
	reSymbol   = regexp.MustCompile(`^[A-Z0-9]{3,}USDT$`)
	reInterval = regexp.MustCompile(`^[0-9]+(m|h|d|w)$`)
)

type Parsed struct {
	Symbol   string
	Interval string
	Command  string // ANALYZE / HELP / PRICE / CHART
}

func ParseMessage(msg string) Parsed {
	raw := strings.TrimSpace(msg)

	// if user types anything not starting with '/', treat as HELP
	if raw == "" || !strings.HasPrefix(raw, "/") {
		return Parsed{
			Command:  "HELP",
			Symbol:   "BTCUSDT",
			Interval: "15m",
		}
	}

	parts := strings.Fields(strings.ToUpper(raw))

	out := Parsed{Command: "ANALYZE", Symbol: "", Interval: "15m"}

	// command routing
	if len(parts) > 0 {
		switch parts[0] {
		case "/HELP", "/START":
			out.Command = "HELP"
			return out
		case "/PRICE":
			out.Command = "PRICE"
			parts = parts[1:]
		case "/CHART":
			out.Command = "CHART"
			parts = parts[1:]
		case "/ANALYZE":
			out.Command = "ANALYZE"
			parts = parts[1:]
		default:
			// unknown command -> show help
			out.Command = "HELP"
			return out
		}
	}

	// parse args (symbol, interval)
	for _, p := range parts {
		if reSymbol.MatchString(p) {
			out.Symbol = p
		}
		if reInterval.MatchString(strings.ToLower(p)) {
			out.Interval = strings.ToLower(p)
		}
	}

	// defaults
	if out.Symbol == "" {
		out.Symbol = "BTCUSDT"
	}

	// for /price, interval not required; keep default 15m (harmless)

	return out
}
