//go:build tinygo

package rowexec

import (
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/plan"
)

func (b *BaseBuilder) buildAlterUserImpl(ctx *sql.Context, a *plan.AlterUser) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}

func (b *BaseBuilder) buildCreateUserImpl(ctx *sql.Context, n *plan.CreateUser) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}
