//go:build !tinygo

package analyzer

import (
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/mysql_db"
)

func newCatalogMySQLDb() catalogMySQLDb {
	return mysql_db.CreateEmptyMySQLDb()
}

func newCatalogDatabaseProvider(db catalogMySQLDb, pro sql.DatabaseProvider, auth sql.AuthorizationHandler) sql.DatabaseProvider {
	return mysql_db.NewPrivilegedDatabaseProvider(db.(*mysql_db.MySQLDb), pro, auth)
}
