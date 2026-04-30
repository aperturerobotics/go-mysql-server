//go:build tinygo

package analyzer

import (
	"github.com/dolthub/go-mysql-server/sql"
)

type informationSchemaColumnsTable interface {
	AllColumns(ctx *sql.Context) (sql.Schema, error)
	WithColumnDefaults(columnDefaults []sql.Expression) (sql.Table, error)
}

func getInformationSchemaColumnsTable(table sql.Table) (informationSchemaColumnsTable, bool) {
	return nil, false
}
