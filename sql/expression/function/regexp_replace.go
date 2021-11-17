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
	"regexp"
	"strings"
	"sync/atomic"

	"gopkg.in/src-d/go-errors.v1"

	"github.com/dolthub/go-mysql-server/sql"
)

// RegexpReplace implements the REGEXP_REPLACE function.
// https://dev.mysql.com/doc/refman/8.0/en/regexp.html#function_regexp-replace
type RegexpReplace struct {
	args []sql.Expression

	cachedVal  atomic.Value
	re         *regexp.Regexp
	compileErr error
}

var _ sql.FunctionExpression = (*RegexpReplace)(nil)

// NewRegexpReplace creates a new RegexpReplace expression.
func NewRegexpReplace(args ...sql.Expression) (sql.Expression, error) {
	if len(args) < 3 || len(args) > 6 {
		return nil, sql.ErrInvalidArgumentNumber.New("regexp_replace", "3,4,5 or 6", len(args))
	}

	return &RegexpReplace{args: args}, nil
}

// FunctionName implements sql.FunctionExpression
func (r *RegexpReplace) FunctionName() string {
	return "regexp_replace"
}

// Type implements the sql.Expression interface.
func (r *RegexpReplace) Type() sql.Type { return sql.LongText }

// IsNullable implements the sql.Expression interface.
func (r *RegexpReplace) IsNullable() bool { return true }

// Children implements the sql.Expression interface.
func (r *RegexpReplace) Children() []sql.Expression {
	return r.args
}

// Resolved implements the sql.Expression interface.
func (r *RegexpReplace) Resolved() bool {
	for _, arg := range r.args {
		if !arg.Resolved() {
			return false
		}
	}
	return true
}

// WithChildren implements the sql.Expression interface.
func (r *RegexpReplace) WithChildren(children ...sql.Expression) (sql.Expression, error) {
	if len(children) != len(r.args) {
		return nil, sql.ErrInvalidChildrenNumber.New(r, len(children), len(r.args))
	}
	return NewRegexpReplace(children...)
}

func (r *RegexpReplace) String() string {
	var args []string
	for _, e := range r.args {
		args = append(args, e.String())
	}
	return fmt.Sprintf("regexp_replace(%s)", strings.Join(args, ", "))
}

func (r *RegexpReplace) compile(ctx *sql.Context, row sql.Row) {
	pattern := r.args[1]
	var flags sql.Expression = nil
	if len(r.args) == 6 {
		flags = r.args[5]
	}
	r.re, r.compileErr = compileRegex(ctx, pattern, flags, r.FunctionName(), row)
}

// Eval implements the sql.Expression interface.
func (r *RegexpReplace) Eval(ctx *sql.Context, row sql.Row) (interface{}, error) {
	span, ctx := ctx.Span("function.RegexpReplace")
	defer span.Finish()

	cached := r.cachedVal.Load()
	if cached != nil {
		return cached, nil
	}

	// Evaluate string value
	str, err := r.args[0].Eval(ctx, row)
	if err != nil {
		return nil, err
	}
	if str == nil {
		return nil, nil
	}
	str, err = sql.LongText.Convert(str)
	if err != nil {
		return nil, err
	}

	// Convert to string
	_str := str.(string)

	// Create regex, should handle null pattern and null flags
	r.compile(ctx, row) // TODO: Should I write my own instead of using regex_like
	if r.compileErr != nil {
		return "compile bad", r.compileErr
	}
	if r.re == nil {
		return nil, nil
	}

	// Evaluate ReplaceStr
	replaceStr, err := r.args[2].Eval(ctx, row)
	if err != nil {
		return nil, err
	}
	if replaceStr == nil {
		return nil, nil
	}
	replaceStr, err = sql.LongText.Convert(replaceStr)
	if err != nil {
		return nil, err
	}

	// Convert to string
	_replaceStr := replaceStr.(string)

	// Do nothing if str is empty
	if len(_str) == 0 {
		return _str, nil
	}

	// Check if position argument was provided
	var pos interface{} = nil
	if len(r.args) >= 4 {
		pos, err = r.args[3].Eval(ctx, row)
		if err != nil {
			return nil, err
		}
		if pos == nil {
			return nil, nil
		}
	} else {
		// Default position is 1, if argument not provided
		pos = 1 // TODO: use constant
	}

	// Convert to int32
	pos, err = sql.Int32.Convert(pos)
	if err != nil {
		return nil, err
	}
	// Convert to int
	_pos := int(pos.(int32))

	// Non-positive position throws incorrect parameter
	if _pos <= 0 {
		return nil, ErrInvalidArgument.New(r.FunctionName(), fmt.Sprintf("%d", _pos))
	}

	// Handle out of bounds
	if _pos > len(_str) {
		return nil, errors.NewKind("Index out of bounds for regular expression search.").New()
	}

	// Check if Occurrence argument was provided
	var occ interface{} = nil
	if len(r.args) >= 5 {
		occ, err = r.args[4].Eval(ctx, row)
		if err != nil {
			return nil, err
		}
		if occ == nil {
			return nil, nil
		}
	} else {
		// Default occurrence is 0 (replace all occurrences)
		occ = 0 // TODO: use constant
	}

	// Convert occurrence to int32
	occ, err = sql.Int32.Convert(occ)
	if err != nil {
		return nil, err
	}
	_occ := int(occ.(int32))

	// MySQL interprets negative occurrences as first for some reason
	if _occ < 0 {
		_occ = 1
	} else if _occ == 0 {
		// Replace everything
		return _str[:_pos-1] + r.re.ReplaceAllString(_str[_pos-1:], _replaceStr), nil
	}

	// Split string into prefix and suffix
	prefix := _str[:_pos-1]
	suffix := _str[_pos-1:]

	// Extract all matches
	matches := r.re.FindAllString(suffix, -1)
	indexes := r.re.FindAllStringIndex(suffix, -1)

	// No matches, return original string
	if len(matches) == 0 {
		return _str, nil
	}

	// If there aren't enough occurrences
	if _occ > len(matches) {
		return _str, nil
	}

	// Replace only the nth occurrence
	matches[_occ-1] = _replaceStr

	// Initialize result string
	res := prefix                 // attach prefix
	res += suffix[:indexes[0][0]] // attach text before first match
	res += matches[0]             // attach first match

	// Recombine rest of matches
	for i := 1; i < len(matches); i++ {
		// Attach text before match
		res += suffix[indexes[i-1][1]:indexes[i][0]] // end of prev to start of curr match
		// Attach match
		res += matches[i]
	}

	// Append text after last match
	res += suffix[indexes[len(indexes)-1][1]:]

	return res, nil
}
