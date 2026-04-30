//go:build tinygo

package mysql_db

import "github.com/dolthub/vitess/go/mysql"

func tlsCipherSuiteName(_ *mysql.Conn) string {
	return ""
}
