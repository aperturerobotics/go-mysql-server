//go:build !sql_lite

package planbuilder

import (
	ast "github.com/dolthub/vitess/go/vt/sqlparser"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression/function/json"
)

func (b *Builder) buildJSONExtract(operator string, l sql.Expression, r sql.Expression) (sql.Expression, error) {
	jsonExtract, err := json.NewJSONExtract(b.ctx, l, r)
	if err != nil {
		return nil, err
	}
	if operator == ast.JSONUnquoteExtractOp {
		return json.NewJSONUnquote(b.ctx, jsonExtract), nil
	}
	return jsonExtract, nil
}
