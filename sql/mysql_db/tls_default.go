//go:build !tinygo

package mysql_db

import (
	"crypto/tls"

	"github.com/dolthub/go-mysql-server/sql/mysql"
)

func tlsCipherSuiteName(conn *mysql.Conn) string {
	tlsConn, ok := conn.Conn.(*tls.Conn)
	if !ok {
		return ""
	}
	return tls.CipherSuiteName(tlsConn.ConnectionState().CipherSuite)
}
