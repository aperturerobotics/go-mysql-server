//go:build tinygo

package server

import "github.com/dolthub/vitess/go/mysql"

func getAuthServer(v interface{}) mysql.AuthServer {
	return nil
}
