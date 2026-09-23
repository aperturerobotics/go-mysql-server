//go:build !js

package function

import (
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
)

// formatGrouped renders the whole and fractional parts of a FORMAT() result
// using locale-aware grouping and decimal separators from golang.org/x/text.
// An empty or unparseable locale falls back to English. The browser build uses
// a reflect-free variant (format_js.go) to keep x/text out of the GoScript
// dependency graph.
func formatGrouped(ctx *sql.Context, localeValue any, negative string, whole int64, fractionStr string, numDecimalPlaces int) (string, error) {
	// Preserve native locale conversion and MySQL's warning for invalid values.
	localeValue, _, err := types.Text.Convert(ctx, localeValue)
	if err != nil {
		return "", err
	}

	locale := language.English
	if localeValue == nil {
		ctx.Warn(1649, "Unknown Locale: 'NULL'")
	} else {
		locale, err = language.Parse(localeValue.(string))
		if err != nil {
			ctx.Warn(1649, "Unknown Locale: %s", localeValue)
			locale = language.English
		}
	}

	p := message.NewPrinter(locale)
	formattedWhole := p.Sprintf("%v", number.Decimal(whole))
	if numDecimalPlaces == 0 {
		return negative + formattedWhole, nil
	}

	decimalChar := p.Sprintf("%v", number.Decimal(1.5))
	if len(fractionStr) < numDecimalPlaces {
		fractionStr += strings.Repeat("0", numDecimalPlaces-len(fractionStr))
	}

	return negative + formattedWhole + decimalChar[1:2] + fractionStr, nil
}
