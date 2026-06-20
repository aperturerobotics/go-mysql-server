//go:build js

package function

import (
	"strconv"
	"strings"
)

// numFormat describes the grouping and separator style for a locale.
type numFormat struct {
	group   string
	decimal string
	indian  bool
}

// enUS is the Western default: comma grouping, period decimal, groups of three.
var enUS = numFormat{group: ",", decimal: ".", indian: false}

// localeFormats maps MySQL FORMAT() locales that differ from the en_US default
// to their grouping and decimal separators. Western comma/period locales and
// any unrecognized locale use enUS. This is the reflect-free browser variant of
// the FORMAT() locale formatting; the native build (format_default.go) uses
// golang.org/x/text for full CLDR coverage, which the GoScript browser closure
// excludes to keep reflect out of the dependency graph. Locale codes whose
// x/text output is already documented as non-MySQL-conforming fall back to
// en_US in the browser.
var localeFormats = map[string]numFormat{
	"da_DK": {group: ".", decimal: ",", indian: false},
	"de_BE": {group: ".", decimal: ",", indian: false},
	"de_DE": {group: ".", decimal: ",", indian: false},
	"de_LU": {group: ".", decimal: ",", indian: false},
	"es_AR": {group: ".", decimal: ",", indian: false},
	"fo_FO": {group: ".", decimal: ",", indian: false},
	"id_ID": {group: ".", decimal: ",", indian: false},
	"is_IS": {group: ".", decimal: ",", indian: false},
	"ro_RO": {group: ".", decimal: ",", indian: false},
	"tr_TR": {group: ".", decimal: ",", indian: false},
	"vi_VN": {group: ".", decimal: ",", indian: false},

	"en_IN": {group: ",", decimal: ".", indian: true},
	"ta_IN": {group: ",", decimal: ".", indian: true},
	"te_IN": {group: ",", decimal: ".", indian: true},
}

// formatGrouped renders the whole and fractional parts of a FORMAT() result
// using a reflect-free locale table. An empty or unrecognized locale uses the
// en_US default.
func formatGrouped(localeStr, negative string, whole int64, fractionStr string, numDecimalPlaces int) string {
	nf, ok := localeFormats[localeStr]
	if !ok {
		nf = enUS
	}

	sign := ""
	u := whole
	if u < 0 {
		sign = "-"
		u = -u
	}
	formattedWhole := sign + groupDigits(strconv.FormatInt(u, 10), nf.group, nf.indian)
	if numDecimalPlaces == 0 {
		return negative + formattedWhole
	}

	if len(fractionStr) < numDecimalPlaces {
		fractionStr += strings.Repeat("0", numDecimalPlaces-len(fractionStr))
	}
	return negative + formattedWhole + nf.decimal + fractionStr
}

// groupDigits inserts sep into an unsigned base-10 digit string. Western
// grouping splits into groups of three from the right; Indian grouping keeps a
// final group of three preceded by groups of two.
func groupDigits(digits, sep string, indian bool) string {
	n := len(digits)
	if n <= 3 || sep == "" {
		return digits
	}

	if !indian {
		var b strings.Builder
		first := n % 3
		if first > 0 {
			b.WriteString(digits[:first])
		}
		for i := first; i < n; i += 3 {
			if b.Len() > 0 {
				b.WriteString(sep)
			}
			b.WriteString(digits[i : i+3])
		}
		return b.String()
	}

	last3 := digits[n-3:]
	rest := digits[:n-3]
	m := len(rest)
	var b strings.Builder
	first := m % 2
	if first > 0 {
		b.WriteString(rest[:first])
	}
	for i := first; i < m; i += 2 {
		if b.Len() > 0 {
			b.WriteString(sep)
		}
		b.WriteString(rest[i : i+2])
	}
	b.WriteString(sep)
	b.WriteString(last3)
	return b.String()
}
