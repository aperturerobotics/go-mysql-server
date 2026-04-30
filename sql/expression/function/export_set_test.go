// Copyright 2020-2024 Dolthub, Inc.
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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression"
	"github.com/dolthub/go-mysql-server/sql/types"
)

func TestExportSet(t *testing.T) {
	testCases := []struct {
		name     string
		args     []any
		expected any
		err      bool
	}{
		// MySQL documentation examples
		{"mysql example 1", []any{5, "Y", "N", ",", 4}, "Y,N,Y,N", false},
		{"mysql example 2", []any{6, "1", "0", ",", 10}, "0,1,1,0,0,0,0,0,0,0", false},

		// Basic functionality tests
		{"zero value", []any{0, "1", "0", ",", 4}, "0,0,0,0", false},
		{"all bits set", []any{15, "1", "0", ",", 4}, "1,1,1,1", false},
		{"single bit", []any{1, "T", "F", ",", 3}, "T,F,F", false},
		{"single bit position 2", []any{2, "T", "F", ",", 3}, "F,T,F", false},
		{"single bit position 3", []any{4, "T", "F", ",", 3}, "F,F,T", false},

		// Different separators
		{"pipe separator", []any{5, "1", "0", "|", 4}, "1|0|1|0", false},
		{"space separator", []any{5, "1", "0", " ", 4}, "1 0 1 0", false},
		{"empty separator", []any{5, "1", "0", "", 4}, "1010", false},
		{"no separator specified", []any{5, "1", "0"}, "1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0", false},

		// Different on/off strings
		{"word strings", []any{5, "ON", "OFF", ",", 4}, "ON,OFF,ON,OFF", false},
		{"empty on string", []any{5, "", "0", ",", 4}, ",0,,0", false},
		{"empty off string", []any{5, "1", "", ",", 4}, "1,,1,", false},

		// Number of bits tests
		{"no number of bits specified", []any{5, "1", "0"}, "1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0", false},
		{"1 bit", []any{5, "1", "0", ",", 1}, "1", false},
		{"8 bits", []any{255, "1", "0", ",", 8}, "1,1,1,1,1,1,1,1", false},
		{"large number of bits", []any{5, "1", "0", ",", 100}, "1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0", false},
		{"negative number of bits", []any{5, "1", "0", ",", -5}, "1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0", false},

		// Large numbers
		{"large number", []any{4294967295, "1", "0", ",", 32}, "1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1", false},
		{"powers of 2", []any{1024, "1", "0", ",", 12}, "0,0,0,0,0,0,0,0,0,0,1,0", false},

		// NULL handling
		{"null bits", []any{nil, "1", "0", ",", 4}, nil, false},
		{"null on", []any{5, nil, "0", ",", 4}, nil, false},
		{"null off", []any{5, "1", nil, ",", 4}, nil, false},
		{"null separator", []any{5, "1", "0", nil, 4}, nil, false},
		{"null number of bits", []any{5, "1", "0", ",", nil}, nil, false},

		// Type conversion
		{"string number", []any{"5", "1", "0", ",", 4}, "1,0,1,0", false},
		{"float number", []any{5.7, "1", "0", ",", 4}, "0,1,1,0", false},
		{"negative number", []any{-1, "1", "0", ",", 4}, "1,1,1,1", false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			require := require.New(t)
			ctx := sql.NewEmptyContext()

			// Convert test args to expressions
			args := make([]sql.Expression, len(tt.args))
			for i, arg := range tt.args {
				if arg == nil {
					args[i] = expression.NewLiteral(nil, types.Null)
				} else {
					switch v := arg.(type) {
					case int:
						args[i] = expression.NewLiteral(int64(v), types.Int64)
					case string:
						args[i] = expression.NewLiteral(v, types.LongText)
					default:
						args[i] = expression.NewLiteral(v, types.LongText)
					}
				}
			}

			f, err := NewExportSet(ctx, args...)
			require.NoError(err)

			v, err := f.Eval(ctx, nil)
			if tt.err {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tt.expected, v)
			}
		})
	}
}

func TestExportSetArguments(t *testing.T) {
	require := require.New(t)
	ctx := sql.NewEmptyContext()

	// Test invalid number of arguments
	_, err := NewExportSet(ctx)
	require.Error(err)

	_, err = NewExportSet(ctx, expression.NewLiteral(1, types.Int64))
	require.Error(err)

	_, err = NewExportSet(ctx, expression.NewLiteral(1, types.Int64), expression.NewLiteral("1", types.Text))
	require.Error(err)

	// Test too many arguments
	args := make([]sql.Expression, 6)
	for i := range args {
		args[i] = expression.NewLiteral(1, types.Int64)
	}
	_, err = NewExportSet(ctx, args...)
	require.Error(err)

	// Test valid argument counts
	validArgs := [][]sql.Expression{
		{expression.NewLiteral(1, types.Int64), expression.NewLiteral("1", types.Text), expression.NewLiteral("0", types.Text)},
		{expression.NewLiteral(1, types.Int64), expression.NewLiteral("1", types.Text), expression.NewLiteral("0", types.Text), expression.NewLiteral(",", types.Text)},
		{expression.NewLiteral(1, types.Int64), expression.NewLiteral("1", types.Text), expression.NewLiteral("0", types.Text), expression.NewLiteral(",", types.Text), expression.NewLiteral(4, types.Int64)},
	}

	for _, args := range validArgs {
		_, err := NewExportSet(ctx, args...)
		require.NoError(err)
	}
}
