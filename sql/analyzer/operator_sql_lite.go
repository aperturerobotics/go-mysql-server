//go:build sql_lite

package analyzer

import "github.com/dolthub/go-mysql-server/sql"

func operatorOf(e sql.Expression) (string, bool) {
	o, ok := e.(interface{ Operator() string })
	if !ok {
		return "", false
	}
	return o.Operator(), true
}
