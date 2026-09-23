// Copyright 2021 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package function

import (
	"fmt"
	"math"
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
)

// Format function returns a result of NumValue rounded to NumDecimalPlaces as a string.
type Format struct {
	NumValue         sql.Expression
	NumDecimalPlaces sql.Expression
	Locale           sql.Expression
}

var _ sql.FunctionExpression = (*Format)(nil)
var _ sql.CollationCoercible = (*Format)(nil)

// NewFormat returns a new Format expression.
func NewFormat(ctx *sql.Context, args ...sql.Expression) (sql.Expression, error) {
	var numValue, numDecimalPlaces, locale sql.Expression
	switch len(args) {
	case 2:
		numValue = args[0]
		numDecimalPlaces = args[1]
		locale = nil
	case 3:
		numValue = args[0]
		numDecimalPlaces = args[1]
		locale = args[2]
	default:
		return nil, sql.ErrInvalidArgumentNumber.New("FORMAT", "2 or 3", len(args))
	}
	return &Format{numValue, numDecimalPlaces, locale}, nil
}

// FunctionName implements sql.FunctionExpression
func (f *Format) FunctionName() string {
	return "format"
}

// Description implements sql.FunctionExpression
func (f *Format) Description() string {
	return "returns a number formatted to specified number of decimal places."
}

// Type implements the Expression interface.
func (f *Format) Type(ctx *sql.Context) sql.Type { return types.LongText }

// CollationCoercibility implements the interface sql.CollationCoercible.
func (*Format) CollationCoercibility(ctx *sql.Context) (collation sql.CollationID, coercibility byte) {
	return ctx.GetCollation(), 4
}

// IsNullable implements the Expression interface.
func (f *Format) IsNullable(ctx *sql.Context) bool {
	return f.NumValue.IsNullable(ctx) || f.NumDecimalPlaces.IsNullable(ctx) || (f.Locale != nil && f.Locale.IsNullable(ctx))
}

func (f *Format) String() string {
	if f.Locale == nil {
		return fmt.Sprintf("%s(%s,%s)", f.FunctionName(), f.NumValue, f.NumDecimalPlaces)
	}
	return fmt.Sprintf("%s(%s,%s,%s)", f.FunctionName(), f.NumValue, f.NumDecimalPlaces, f.Locale)
}

// Eval implements the Expression interface.
func (f *Format) Eval(ctx *sql.Context, row sql.Row) (any, error) {
	numVal, err := f.NumValue.Eval(ctx, row)
	if err != nil {
		return nil, err
	}

	numVal, _, err = types.Float64.Convert(ctx, numVal)
	if err != nil {
		ctx.Warn(1292, "Truncated incorrect DOUBLE value: %s", numVal)
	}
	if numVal == nil {
		return nil, nil
	}
	numValue := numVal.(float64)

	numDP, err := f.NumDecimalPlaces.Eval(ctx, row)
	if err != nil {
		return nil, err
	}
	numDP, _, err = types.Float64.Convert(ctx, numDP)
	if err != nil {
		ctx.Warn(1292, "Truncated incorrect DOUBLE value: %s", numDP)
	}
	if numDP == nil {
		return nil, nil
	}
	numDecimalPlaces := numDP.(float64)
	numDecimalPlaces = math.Round(numDecimalPlaces)

	// MySQL clamps numDecimalPlaces in [0, 30]
	numDecimalPlaces = max(0, numDecimalPlaces)
	numDecimalPlaces = min(30, numDecimalPlaces)

	var localeValue any = "en_US"
	if f.Locale != nil {
		loc, lErr := f.Locale.Eval(ctx, row)
		if lErr != nil {
			return nil, lErr
		}

		localeValue = loc
	}

	// Shift the requested decimal place into the integer part before rounding.
	// Floating-point precision can differ from MySQL at large scales.
	roundedValue := math.Round(numValue*math.Pow(10.0, numDecimalPlaces)) / math.Pow(10.0, numDecimalPlaces)

	// Values rounded to zero omit the negative sign.
	var whole int64
	var fractionStr string
	var negative string
	if roundedValue != 0 {
		res := types.DecimalFromFloat64(roundedValue)
		whole = types.DecimalTruncatedIntPart(res)
		if whole == 0 && res.Negative {
			negative = "-"
		}

		str := res.Text('f')
		_, after, ok := strings.Cut(str, ".")
		if !ok {
			fractionStr = ""
		} else {
			fractionStr = after
		}
	}

	return formatGrouped(ctx, localeValue, negative, whole, fractionStr, int(numDecimalPlaces))
}

// Resolved implements the Expression interface.
func (f *Format) Resolved() bool {
	if f.Locale == nil {
		return f.NumValue.Resolved() && f.NumDecimalPlaces.Resolved()
	}
	return f.NumValue.Resolved() && f.NumDecimalPlaces.Resolved() && f.Locale.Resolved()
}

// Children implements the Expression interface.
func (f *Format) Children() []sql.Expression {
	if f.Locale == nil {
		return []sql.Expression{f.NumValue, f.NumDecimalPlaces}
	}
	return []sql.Expression{f.NumValue, f.NumDecimalPlaces, f.Locale}
}

// WithChildren implements the Expression interface.
func (f *Format) WithChildren(ctx *sql.Context, children ...sql.Expression) (sql.Expression, error) {
	if (len(children) == 2 && f.Locale == nil) || (len(children) == 3 && f.Locale != nil) {
		return NewFormat(ctx, children...)
	}
	return nil, sql.ErrInvalidChildrenNumber.New(f, len(children), 2)
}
