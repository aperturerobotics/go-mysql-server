//go:build sql_lite

package analyzer

import "github.com/dolthub/go-mysql-server/sql"

func isSpatialIndexExpression(sql.Expression) bool {
	return false
}
