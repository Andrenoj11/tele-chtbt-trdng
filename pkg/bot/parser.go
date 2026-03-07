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
	Command  string // ANALYZE / HELP
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

	if len(parts) > 0 && (parts[0] == "/HELP" || parts[0] == "/START") {
		out.Command = "HELP"
		return out
	}

	if len(parts) > 0 && parts[0] == "/ANALYZE" {
		parts = parts[1:]
	}

	for _, p := range parts {
		if reSymbol.MatchString(p) {
			out.Symbol = p
		}
		if reInterval.MatchString(strings.ToLower(p)) {
			out.Interval = strings.ToLower(p)
		}
	}

	if out.Symbol == "" {
		out.Symbol = "BTCUSDT"
	}
	return out
}

// func ParseMessage(msg string) Parsed {
// 	parts := strings.Fields(strings.ToUpper(strings.TrimSpace(msg)))

// 	out := Parsed{Command: "ANALYZE", Symbol: "", Interval: "15m"}

// 	if len(parts) > 0 && (parts[0] == "HELP" || parts[0] == "/HELP" || parts[0] == "START" || parts[0] == "/START") {
// 		out.Command = "HELP"
// 		return out
// 	}

// 	if len(parts) > 0 && (parts[0] == "ANALYZE" || parts[0] == "/ANALYZE") {
// 		parts = parts[1:]
// 	}

// 	for _, p := range parts {
// 		if reSymbol.MatchString(p) {
// 			out.Symbol = p
// 		}
// 		if reInterval.MatchString(strings.ToLower(p)) {
// 			out.Interval = strings.ToLower(p)
// 		}
// 	}

// 	if out.Symbol == "" && len(parts) > 0 && reSymbol.MatchString(parts[0]) {
// 		out.Symbol = parts[0]
// 	}
// 	if out.Symbol == "" {
// 		out.Symbol = "BTCUSDT"
// 	}
// 	return out
// }
