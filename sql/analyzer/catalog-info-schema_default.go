//go:build !tinygo

package analyzer

import (
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
)

func newCatalogInfoSchema() sql.Database {
	return information_schema.NewInformationSchemaDatabase()
}
