//go:build tinygo

package plan

import (
	"fmt"
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
)

type Grant struct {
	MySQLDb         sql.Database
	Catalog         sql.Catalog
	As              *GrantUserAssumption
	PrivilegeLevel  PrivilegeLevel
	Privileges      []Privilege
	Users           []UserName
	ObjectType      ObjectType
	WithGrantOption bool
}

func (n *Grant) Schema(ctx *sql.Context) sql.Schema { return types.OkResultSchema }
func (n *Grant) IsReadOnly() bool                   { return false }
func (n *Grant) String() string {
	users := make([]string, len(n.Users))
	for i, user := range n.Users {
		users[i] = user.String("")
	}
	return fmt.Sprintf("Grant(On: %s, To: %s)", n.PrivilegeLevel.String(), strings.Join(users, ", "))
}
func (n *Grant) Database() sql.Database { return n.MySQLDb }
func (n *Grant) WithDatabase(db sql.Database) (sql.Node, error) {
	nn := *n
	nn.MySQLDb = db
	return &nn, nil
}
func (n *Grant) Resolved() bool {
	_, ok := n.MySQLDb.(sql.UnresolvedDatabase)
	return !ok
}
func (n *Grant) Children() []sql.Node { return nil }
func (n *Grant) WithChildren(ctx *sql.Context, children ...sql.Node) (sql.Node, error) {
	return NillaryWithChildren(n, children...)
}
func (n *Grant) CheckAuth(ctx *sql.Context, opChecker sql.PrivilegedOperationChecker) bool {
	return false
}
func (*Grant) CollationCoercibility(ctx *sql.Context) (sql.CollationID, byte) {
	return sql.Collation_binary, 7
}

type GrantRole struct {
	MySQLDb         sql.Database
	Roles           []UserName
	TargetUsers     []UserName
	WithAdminOption bool
}

func NewGrantRole(roles []UserName, users []UserName, withAdmin bool) *GrantRole {
	return &GrantRole{Roles: roles, TargetUsers: users, WithAdminOption: withAdmin, MySQLDb: sql.UnresolvedDatabase("mysql")}
}
func (n *GrantRole) Schema(ctx *sql.Context) sql.Schema { return types.OkResultSchema }
func (n *GrantRole) String() string                     { return "GrantRole" }
func (n *GrantRole) Database() sql.Database             { return n.MySQLDb }
func (n *GrantRole) WithDatabase(db sql.Database) (sql.Node, error) {
	nn := *n
	nn.MySQLDb = db
	return &nn, nil
}
func (n *GrantRole) Resolved() bool {
	_, ok := n.MySQLDb.(sql.UnresolvedDatabase)
	return !ok
}
func (n *GrantRole) Children() []sql.Node { return nil }
func (n *GrantRole) IsReadOnly() bool     { return false }
func (n *GrantRole) WithChildren(ctx *sql.Context, children ...sql.Node) (sql.Node, error) {
	return NillaryWithChildren(n, children...)
}
func (n *GrantRole) CheckAuth(ctx *sql.Context, opChecker sql.PrivilegedOperationChecker) bool {
	return false
}
func (*GrantRole) CollationCoercibility(ctx *sql.Context) (sql.CollationID, byte) {
	return sql.Collation_binary, 7
}

type GrantProxy struct {
	On              UserName
	To              []UserName
	WithGrantOption bool
}

func NewGrantProxy(on UserName, to []UserName, withGrant bool) *GrantProxy {
	return &GrantProxy{On: on, To: to, WithGrantOption: withGrant}
}
func (n *GrantProxy) Schema(ctx *sql.Context) sql.Schema { return types.OkResultSchema }
func (n *GrantProxy) String() string                     { return "GrantProxy" }
func (n *GrantProxy) Resolved() bool                     { return true }
func (n *GrantProxy) IsReadOnly() bool                   { return false }
func (n *GrantProxy) Children() []sql.Node               { return nil }
func (n *GrantProxy) WithChildren(ctx *sql.Context, children ...sql.Node) (sql.Node, error) {
	return NillaryWithChildren(n, children...)
}
func (n *GrantProxy) CheckAuth(ctx *sql.Context, opChecker sql.PrivilegedOperationChecker) bool {
	return false
}
func (*GrantProxy) CollationCoercibility(ctx *sql.Context) (sql.CollationID, byte) {
	return sql.Collation_binary, 7
}
