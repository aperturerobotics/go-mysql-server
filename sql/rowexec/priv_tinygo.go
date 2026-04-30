//go:build tinygo

package rowexec

import (
	"github.com/pkg/errors"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/plan"
)

var errTinyGoPrivilegesUnsupported = errors.New("mysql privilege statements are unsupported in TinyGo")

func (b *BaseBuilder) buildFlushPrivileges(ctx *sql.Context, n *plan.FlushPrivileges, row sql.Row) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}

func (b *BaseBuilder) buildDropUser(ctx *sql.Context, n *plan.DropUser, row sql.Row) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}

func (b *BaseBuilder) buildRevokeRole(ctx *sql.Context, n *plan.RevokeRole, row sql.Row) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}

func (b *BaseBuilder) buildDropRole(ctx *sql.Context, n *plan.DropRole, row sql.Row) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}

func (b *BaseBuilder) buildRevokeProxy(ctx *sql.Context, n *plan.RevokeProxy, row sql.Row) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}

func (b *BaseBuilder) buildGrantRole(ctx *sql.Context, n *plan.GrantRole, row sql.Row) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}

func (b *BaseBuilder) buildGrantProxy(ctx *sql.Context, n *plan.GrantProxy, row sql.Row) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}

func (b *BaseBuilder) buildRenameUser(ctx *sql.Context, n *plan.RenameUser, row sql.Row) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}

func (b *BaseBuilder) buildRevoke(ctx *sql.Context, n *plan.Revoke, row sql.Row) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}

func (b *BaseBuilder) buildGrant(ctx *sql.Context, n *plan.Grant, row sql.Row) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}

func (b *BaseBuilder) buildCreateRole(ctx *sql.Context, n *plan.CreateRole, row sql.Row) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}
