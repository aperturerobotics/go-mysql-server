//go:build tinygo || js

package mysql_db

import "github.com/dolthub/go-mysql-server/sql/mysql"

func tlsCipherSuiteName(_ *mysql.Conn) string {
	return ""
}
