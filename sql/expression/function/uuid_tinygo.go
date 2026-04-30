//go:build tinygo

package function

import (
	"fmt"

	"github.com/dolthub/vitess/go/sqltypes"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression"
	"github.com/dolthub/go-mysql-server/sql/types"
)

type UUIDFunc struct{}

var _ sql.FunctionExpression = (*UUIDFunc)(nil)
var _ sql.CollationCoercible = (*UUIDFunc)(nil)

func NewUUIDFunc(ctx *sql.Context) sql.Expression {
	return &UUIDFunc{}
}

func (u *UUIDFunc) FunctionName() string { return "uuid" }
func (u *UUIDFunc) Description() string  { return "returns a Universal Unique Identifier." }
func (u *UUIDFunc) String() string       { return "uuid()" }
func (u *UUIDFunc) Type(ctx *sql.Context) sql.Type {
	return types.MustCreateStringWithDefaults(sqltypes.VarChar, 36)
}
func (u *UUIDFunc) CollationCoercibility(ctx *sql.Context) (sql.CollationID, byte) {
	return sql.Collation_utf8mb3_general_ci, 4
}
func (u *UUIDFunc) Eval(ctx *sql.Context, row sql.Row) (interface{}, error) { return nil, nil }
func (u *UUIDFunc) WithChildren(ctx *sql.Context, children ...sql.Expression) (sql.Expression, error) {
	if len(children) != 0 {
		return nil, sql.ErrInvalidChildrenNumber.New(u, len(children), 0)
	}
	return &UUIDFunc{}, nil
}
func (u *UUIDFunc) Resolved() bool                   { return true }
func (u *UUIDFunc) Children() []sql.Expression       { return nil }
func (u *UUIDFunc) IsNullable(ctx *sql.Context) bool { return true }
func (u *UUIDFunc) IsNonDeterministic() bool         { return true }

type UUIDShortFunc struct{}

var _ sql.FunctionExpression = (*UUIDShortFunc)(nil)
var _ sql.CollationCoercible = (*UUIDShortFunc)(nil)

func NewUUIDShortFunc(ctx *sql.Context) sql.Expression  { return &UUIDShortFunc{} }
func (u *UUIDShortFunc) FunctionName() string           { return "uuid_short" }
func (u *UUIDShortFunc) Description() string            { return "returns a short universal identifier." }
func (u *UUIDShortFunc) String() string                 { return "UUID_SHORT()" }
func (u *UUIDShortFunc) Type(ctx *sql.Context) sql.Type { return types.Uint64 }
func (u *UUIDShortFunc) CollationCoercibility(ctx *sql.Context) (sql.CollationID, byte) {
	return sql.Collation_binary, 5
}
func (u *UUIDShortFunc) Eval(ctx *sql.Context, row sql.Row) (interface{}, error) { return nil, nil }
func (u *UUIDShortFunc) WithChildren(ctx *sql.Context, children ...sql.Expression) (sql.Expression, error) {
	if len(children) != 0 {
		return nil, sql.ErrInvalidChildrenNumber.New(u, len(children), 0)
	}
	return &UUIDShortFunc{}, nil
}
func (u *UUIDShortFunc) Resolved() bool                   { return true }
func (u *UUIDShortFunc) Children() []sql.Expression       { return nil }
func (u *UUIDShortFunc) IsNullable(ctx *sql.Context) bool { return true }
func (u *UUIDShortFunc) IsNonDeterministic() bool         { return true }

type IsUUID struct {
	expression.UnaryExpressionStub
}

var _ sql.FunctionExpression = (*IsUUID)(nil)
var _ sql.CollationCoercible = (*IsUUID)(nil)

func NewIsUUID(ctx *sql.Context, arg sql.Expression) sql.Expression {
	return &IsUUID{expression.UnaryExpressionStub{Child: arg}}
}
func (u *IsUUID) FunctionName() string           { return "is_uuid" }
func (u *IsUUID) Description() string            { return "returns whether argument is a valid UUID." }
func (u *IsUUID) String() string                 { return fmt.Sprintf("%s(%s)", u.FunctionName(), u.Child) }
func (u *IsUUID) Type(ctx *sql.Context) sql.Type { return types.Boolean }
func (u *IsUUID) CollationCoercibility(ctx *sql.Context) (sql.CollationID, byte) {
	return sql.Collation_binary, 5
}
func (u *IsUUID) Eval(ctx *sql.Context, row sql.Row) (interface{}, error) { return nil, nil }
func (u *IsUUID) WithChildren(ctx *sql.Context, children ...sql.Expression) (sql.Expression, error) {
	if len(children) != 1 {
		return nil, sql.ErrInvalidChildrenNumber.New(u, len(children), 1)
	}
	return NewIsUUID(ctx, children[0]), nil
}

type UUIDToBin struct {
	children []sql.Expression
}

var _ sql.FunctionExpression = (*UUIDToBin)(nil)
var _ sql.CollationCoercible = (*UUIDToBin)(nil)

func NewUUIDToBin(ctx *sql.Context, args ...sql.Expression) (sql.Expression, error) {
	return &UUIDToBin{children: args}, nil
}
func (u *UUIDToBin) FunctionName() string { return "uuid_to_bin" }
func (u *UUIDToBin) Description() string  { return "converts a string UUID to a binary UUID." }
func (u *UUIDToBin) String() string       { return "UUID_TO_BIN()" }
func (u *UUIDToBin) Type(ctx *sql.Context) sql.Type {
	return types.MustCreateBinary(sqltypes.Binary, 16)
}
func (u *UUIDToBin) CollationCoercibility(ctx *sql.Context) (sql.CollationID, byte) {
	return sql.Collation_binary, 5
}
func (u *UUIDToBin) Eval(ctx *sql.Context, row sql.Row) (interface{}, error) { return nil, nil }
func (u *UUIDToBin) WithChildren(ctx *sql.Context, children ...sql.Expression) (sql.Expression, error) {
	return &UUIDToBin{children: children}, nil
}
func (u *UUIDToBin) Resolved() bool {
	for _, child := range u.children {
		if !child.Resolved() {
			return false
		}
	}
	return true
}
func (u *UUIDToBin) Children() []sql.Expression       { return u.children }
func (u *UUIDToBin) IsNullable(ctx *sql.Context) bool { return true }

type BinToUUID struct {
	children []sql.Expression
}

var _ sql.FunctionExpression = (*BinToUUID)(nil)
var _ sql.CollationCoercible = (*BinToUUID)(nil)

func NewBinToUUID(ctx *sql.Context, args ...sql.Expression) (sql.Expression, error) {
	return &BinToUUID{children: args}, nil
}
func (u *BinToUUID) FunctionName() string { return "bin_to_uuid" }
func (u *BinToUUID) Description() string  { return "converts a binary UUID to a string UUID." }
func (u *BinToUUID) String() string       { return "BIN_TO_UUID()" }
func (u *BinToUUID) Type(ctx *sql.Context) sql.Type {
	return types.MustCreateStringWithDefaults(sqltypes.VarChar, 36)
}
func (u *BinToUUID) CollationCoercibility(ctx *sql.Context) (sql.CollationID, byte) {
	return sql.Collation_binary, 5
}
func (u *BinToUUID) Eval(ctx *sql.Context, row sql.Row) (interface{}, error) { return nil, nil }
func (u *BinToUUID) WithChildren(ctx *sql.Context, children ...sql.Expression) (sql.Expression, error) {
	return &BinToUUID{children: children}, nil
}
func (u *BinToUUID) Resolved() bool {
	for _, child := range u.children {
		if !child.Resolved() {
			return false
		}
	}
	return true
}
func (u *BinToUUID) Children() []sql.Expression       { return u.children }
func (u *BinToUUID) IsNullable(ctx *sql.Context) bool { return true }
