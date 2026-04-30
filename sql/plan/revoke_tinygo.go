//go:build tinygo

package plan

import (
	"fmt"
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
)

type Revoke struct {
	MySQLDb           sql.Database
	PrivilegeLevel    PrivilegeLevel
	Privileges        []Privilege
	Users             []UserName
	ObjectType        ObjectType
	IgnoreUnknownUser bool
}

func (n *Revoke) Schema(ctx *sql.Context) sql.Schema { return types.OkResultSchema }
func (n *Revoke) IsReadOnly() bool                   { return false }
func (n *Revoke) String() string {
	users := make([]string, len(n.Users))
	for i, user := range n.Users {
		users[i] = user.String("")
	}
	return fmt.Sprintf("Revoke(On: %s, From: %s)", n.PrivilegeLevel.String(), strings.Join(users, ", "))
}
func (n *Revoke) Database() sql.Database { return n.MySQLDb }
func (n *Revoke) WithDatabase(db sql.Database) (sql.Node, error) {
	nn := *n
	nn.MySQLDb = db
	return &nn, nil
}
func (n *Revoke) Resolved() bool {
	_, ok := n.MySQLDb.(sql.UnresolvedDatabase)
	return !ok
}
func (n *Revoke) Children() []sql.Node { return nil }
func (n *Revoke) WithChildren(ctx *sql.Context, children ...sql.Node) (sql.Node, error) {
	return NillaryWithChildren(n, children...)
}
func (n *Revoke) CheckAuth(ctx *sql.Context, opChecker sql.PrivilegedOperationChecker) bool {
	return false
}
func (*Revoke) CollationCoercibility(ctx *sql.Context) (sql.CollationID, byte) {
	return sql.Collation_binary, 7
}

type RevokeRole struct {
	MySQLDb           sql.Database
	Roles             []UserName
	TargetUsers       []UserName
	IfExists          bool
	IgnoreUnknownUser bool
}

func NewRevokeRole(roles []UserName, users []UserName, ifExists, ignoreUnknownUser bool) *RevokeRole {
	return &RevokeRole{Roles: roles, TargetUsers: users, IfExists: ifExists, IgnoreUnknownUser: ignoreUnknownUser, MySQLDb: sql.UnresolvedDatabase("mysql")}
}
func (n *RevokeRole) Schema(ctx *sql.Context) sql.Schema { return types.OkResultSchema }
func (n *RevokeRole) String() string                     { return "RevokeRole" }
func (n *RevokeRole) Database() sql.Database             { return n.MySQLDb }
func (n *RevokeRole) WithDatabase(db sql.Database) (sql.Node, error) {
	nn := *n
	nn.MySQLDb = db
	return &nn, nil
}
func (n *RevokeRole) Resolved() bool {
	_, ok := n.MySQLDb.(sql.UnresolvedDatabase)
	return !ok
}
func (n *RevokeRole) IsReadOnly() bool     { return false }
func (n *RevokeRole) Children() []sql.Node { return nil }
func (n *RevokeRole) WithChildren(ctx *sql.Context, children ...sql.Node) (sql.Node, error) {
	return NillaryWithChildren(n, children...)
}
func (n *RevokeRole) CheckAuth(ctx *sql.Context, opChecker sql.PrivilegedOperationChecker) bool {
	return false
}
func (*RevokeRole) CollationCoercibility(ctx *sql.Context) (sql.CollationID, byte) {
	return sql.Collation_binary, 7
}

type RevokeProxy struct {
	On                UserName
	From              []UserName
	IfExists          bool
	ignoreUnknownUser bool
}

func NewRevokeProxy(on UserName, from []UserName, ifExists, ignoreUnknownUser bool) *RevokeProxy {
	return &RevokeProxy{On: on, From: from, IfExists: ifExists, ignoreUnknownUser: ignoreUnknownUser}
}
func (n *RevokeProxy) Schema(ctx *sql.Context) sql.Schema { return types.OkResultSchema }
func (n *RevokeProxy) String() string                     { return "RevokeProxy" }
func (n *RevokeProxy) Resolved() bool                     { return true }
func (n *RevokeProxy) IsReadOnly() bool                   { return false }
func (n *RevokeProxy) Children() []sql.Node               { return nil }
func (n *RevokeProxy) WithChildren(ctx *sql.Context, children ...sql.Node) (sql.Node, error) {
	return NillaryWithChildren(n, children...)
}
func (n *RevokeProxy) CheckAuth(ctx *sql.Context, opChecker sql.PrivilegedOperationChecker) bool {
	return false
}
func (*RevokeProxy) CollationCoercibility(ctx *sql.Context) (sql.CollationID, byte) {
	return sql.Collation_binary, 7
}
