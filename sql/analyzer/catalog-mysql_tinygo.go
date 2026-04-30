//go:build tinygo

package analyzer

import "github.com/dolthub/go-mysql-server/sql"

func newCatalogMySQLDb() catalogMySQLDb {
	return nil
}

func newCatalogDatabaseProvider(db catalogMySQLDb, pro sql.DatabaseProvider, auth sql.AuthorizationHandler) sql.DatabaseProvider {
	return pro
}
