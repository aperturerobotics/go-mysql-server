//go:build !tinygo

package analyzer

import (
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
)

type informationSchemaColumnsTable interface {
	AllColumns(ctx *sql.Context) (sql.Schema, error)
	WithColumnDefaults(columnDefaults []sql.Expression) (sql.Table, error)
}

func getInformationSchemaColumnsTable(table sql.Table) (informationSchemaColumnsTable, bool) {
	ct, ok := table.(*information_schema.ColumnsTable)
	return ct, ok
}
