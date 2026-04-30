//go:build tinygo

package server

import "github.com/dolthub/go-mysql-server/sql/mysql"

func getAuthServer(v interface{}) mysql.AuthServer {
	return nil
}
