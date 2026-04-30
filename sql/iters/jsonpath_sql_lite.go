//go:build sql_lite

package iters

import "github.com/dolthub/go-mysql-server/sql"

func jsonPathLookup(any, string) (any, error) {
	return nil, sql.ErrUnsupportedFeature.New("json")
}
