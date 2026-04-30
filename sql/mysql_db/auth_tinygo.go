//go:build tinygo

package mysql_db

import (
	"net"

	"github.com/dolthub/vitess/go/mysql"
)

// DefaultAuthMethod specifies the default MySQL auth protocol.
const DefaultAuthMethod = mysql.MysqlNativePassword

type authServer struct{}

func newAuthServer(db *MySQLDb) *authServer {
	return &authServer{}
}

func (as *authServer) AuthMethods() []mysql.AuthMethod {
	return nil
}

func (as *authServer) DefaultAuthMethodDescription() mysql.AuthMethodDescription {
	return DefaultAuthMethod
}

func extractHostAddress(addr net.Addr) (string, error) {
	if addr == nil {
		return "", nil
	}
	return addr.String(), nil
}

func validateMysqlNativePassword(authResponse, salt []byte, mysqlNativePassword string) bool {
	return false
}
