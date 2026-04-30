//go:build sql_lite

package plan

import "github.com/dolthub/go-mysql-server/sql"

// CheckSubStatement permits lite procedures, which cannot contain native callbacks.
func (p *Procedure) CheckSubStatement(*sql.Context) error { return nil }

// ExtendVariadic returns the unchanged procedure because lite procedures are not variadic.
func (p *Procedure) ExtendVariadic(*sql.Context, int) *Procedure {
	return p
}

// IsExternal reports false because the lite profile omits external procedures.
func (p *Procedure) IsExternal() bool {
	return false
}
