//go:build !tinygo && (!js || !sql_lite)

package analyzer

import (
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/mysql_db"
)

type catalogMySQLDb = *mysql_db.MySQLDb

func catalogMySQLDbEnabled(db catalogMySQLDb) bool {
	return db != nil && db.Enabled()
}

func newCatalogMySQLDb() catalogMySQLDb {
	return mysql_db.CreateEmptyMySQLDb()
}

func newCatalogDatabaseProvider(db catalogMySQLDb, pro sql.DatabaseProvider, auth sql.AuthorizationHandler) sql.DatabaseProvider {
	return mysql_db.NewPrivilegedDatabaseProvider(db, pro, auth)
}
