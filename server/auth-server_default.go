//go:build !tinygo

package server

import "github.com/dolthub/go-mysql-server/sql/mysql"

func getAuthServer(v any) mysql.AuthServer {
	return v.(mysql.AuthServer)
}
