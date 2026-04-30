//go:build !sql_lite

package analyzer

import (
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression/function/spatial"
)

func isSpatialIndexExpression(e sql.Expression) bool {
	switch e.(type) {
	case *spatial.Intersects, *spatial.Within, *spatial.STEquals:
		return true
	default:
		return false
	}
}
