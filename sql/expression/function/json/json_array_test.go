// Copyright 2022 Dolthub, Inc.
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

package json

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression"
	"github.com/dolthub/go-mysql-server/sql/types"
)

func TestJSONArray(t *testing.T) {
	ctx := sql.NewEmptyContext()
	f0, err := NewJSONArray(ctx)
	require.NoError(t, err)

	f1, err := NewJSONArray(
		ctx,
		expression.NewGetField(0, types.LongText, "arg1", false),
	)
	require.NoError(t, err)

	f2, err := NewJSONArray(
		ctx,
		expression.NewGetField(0, types.LongText, "arg1", false),
		expression.NewGetField(1, types.LongText, "arg2", false),
	)
	require.NoError(t, err)

	f3, err := NewJSONArray(
		ctx,
		expression.NewGetField(0, types.LongText, "arg1", false),
		expression.NewGetField(1, types.LongText, "arg2", false),
		expression.NewGetField(2, types.LongText, "arg3", false),
	)
	require.NoError(t, err)

	f4, err := NewJSONArray(
		ctx,
		expression.NewGetField(0, types.LongText, "arg1", false),
		expression.NewGetField(1, types.LongText, "arg2", false),
		expression.NewGetField(2, types.LongText, "arg3", false),
		expression.NewGetField(3, types.LongText, "arg4", false),
	)
	require.NoError(t, err)

	testCases := []struct {
		f        sql.Expression
		row      sql.Row
		expected any
		err      error
	}{
		{f0, sql.Row{}, types.JSONDocument{Val: []any{}}, nil},
		{f1, sql.Row{[]any{1, 2}}, types.JSONDocument{Val: []any{[]any{1, 2}}}, nil},
		{f2, sql.Row{[]any{1, 2}, "second item"}, types.JSONDocument{Val: []any{[]any{1, 2}, "second item"}}, nil},
		{f2, sql.Row{[]any{1, 2}, map[string]any{"name": "x"}}, types.JSONDocument{Val: []any{[]any{1, 2}, map[string]any{"name": "x"}}}, nil},
		{f2, sql.Row{map[string]any{"name": "x"}, map[string]any{"id": 47}}, types.JSONDocument{Val: []any{map[string]any{"name": "x"}, map[string]any{"id": 47}}}, nil},
		{f3, sql.Row{"foo", -44, "b"}, types.JSONDocument{Val: []any{"foo", -44, "b"}}, nil},
		{f4, sql.Row{100, true, nil, "four"}, types.JSONDocument{Val: []any{100, true, nil, "four"}}, nil},
		{f4, sql.Row{100.44, `{"name":null,"id":{"number":998,"type":"A"}}`, nil, `four`},
			types.JSONDocument{Val: []any{100.44, "{\"name\":null,\"id\":{\"number\":998,\"type\":\"A\"}}", nil, "four"}}, nil},
	}

	for _, tt := range testCases {
		t.Run(tt.f.String(), func(t *testing.T) {
			require := require.New(t)
			result, err := tt.f.Eval(sql.NewEmptyContext(), tt.row)
			if tt.err == nil {
				require.NoError(err)
			} else {
				require.Equal(err.Error(), tt.err.Error())
			}

			require.Equal(tt.expected, result)
		})
	}
}
