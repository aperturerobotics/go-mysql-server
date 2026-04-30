//go:build sql_lite

package planbuilder

import "github.com/dolthub/go-mysql-server/sql"

func (b *Builder) buildJSONExtract(string, sql.Expression, sql.Expression) (sql.Expression, error) {
	return nil, sql.ErrUnsupportedFeature.New("json")
}
