//go:build tinygo

package rowexec

import (
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/plan"
)

func (b *BaseBuilder) buildShowGrantsImpl(ctx *sql.Context, n *plan.ShowGrants) (sql.RowIter, error) {
	return nil, errTinyGoPrivilegesUnsupported
}
