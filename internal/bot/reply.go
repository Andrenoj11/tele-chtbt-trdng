package bot

import (
	"fmt"
	"strings"

	"crypto-telegram-bot/internal/strategy"
)

func FormatReply(s strategy.Signal) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Symbol: %s | TF: %s\n", s.Symbol, s.Interval)
	fmt.Fprintf(&b, "Price: %.4f\n", s.Price)
	fmt.Fprintf(&b, "Action: %s (%s)\n\n", s.Action, s.Confidence)

	b.WriteString("Reasons:\n")
	for _, r := range s.Reasons {
		b.WriteString("- " + r + "\n")
	}

	b.WriteString("\nRisk notes:\n")
	for _, r := range s.RiskNotes {
		b.WriteString("- " + r + "\n")
	}
	return b.String()
}
