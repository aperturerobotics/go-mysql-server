//go:build sql_lite

package rowexec

import "github.com/dolthub/go-mysql-server/sql"

func buildJSONTableData(*sql.Context, interface{}, string) (interface{}, error) {
	return nil, sql.ErrUnsupportedFeature.New("json")
}
