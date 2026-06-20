//go:build !js

package function

import (
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

// formatGrouped renders the whole and fractional parts of a FORMAT() result
// using locale-aware grouping and decimal separators from golang.org/x/text.
// An empty or unparseable locale falls back to English. The browser build uses
// a reflect-free variant (format_js.go) to keep x/text out of the GoScript
// dependency graph.
func formatGrouped(localeStr, negative string, whole int64, fractionStr string, numDecimalPlaces int) string {
	locale := language.English
	if localeStr != "" {
		if l, err := language.Parse(localeStr); err == nil {
			locale = l
		}
	}

	p := message.NewPrinter(locale)
	formattedWhole := p.Sprintf("%v", number.Decimal(whole))
	if numDecimalPlaces == 0 {
		return negative + formattedWhole
	}

	decimalChar := p.Sprintf("%v", number.Decimal(1.5))
	if len(fractionStr) < numDecimalPlaces {
		fractionStr += strings.Repeat("0", numDecimalPlaces-len(fractionStr))
	}

	return negative + formattedWhole + decimalChar[1:2] + fractionStr
}
