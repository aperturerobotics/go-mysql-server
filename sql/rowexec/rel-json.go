//go:build !sql_lite

package rowexec

import (
	"github.com/dolthub/jsonpath"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression/function/json"
)

func buildJSONTableData(ctx *sql.Context, data any, path string) (any, error) {
	jsonData, err := json.GetJSONFromWrapperOrCoercibleString(ctx, data, "json_table", 1)
	if err != nil {
		return nil, err
	}
	return jsonpath.JsonPathLookup(jsonData, path)
}
