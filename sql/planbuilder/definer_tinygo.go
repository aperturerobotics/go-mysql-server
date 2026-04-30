//go:build tinygo

package planbuilder

import "github.com/dolthub/go-mysql-server/sql"

func (b *Builder) mockDefiner(privileges ...sql.PrivilegeType) func() {
	return func() {}
}
