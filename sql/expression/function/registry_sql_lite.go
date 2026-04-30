//go:build sql_lite

package function

import (
	"gopkg.in/src-d/go-errors.v1"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression/function/aggregation"
)

// ErrFunctionAlreadyRegistered is thrown when a function is already registered
var ErrFunctionAlreadyRegistered = errors.NewKind("function '%s' is already registered")

// BuiltIns is the reduced set of built-in functions for the sql_lite profile.
var BuiltIns = []sql.Function{
	sql.Function1{Name: "abs", Fn: NewAbsVal},
	sql.Function1{Name: "avg", Fn: func(ctx *sql.Context, e sql.Expression) sql.Expression { return aggregation.NewAvg(e) }},
	sql.FunctionN{Name: "coalesce", Fn: NewCoalesce},
	sql.FunctionN{Name: "concat", Fn: NewConcat},
	sql.Function1{Name: "count", Fn: func(ctx *sql.Context, e sql.Expression) sql.Expression { return aggregation.NewCount(e) }},
	sql.Function2{Name: "ifnull", Fn: NewIfNull},
	sql.Function1{Name: "length", Fn: NewLength},
	sql.Function1{Name: "lower", Fn: NewLower},
	sql.Function1{Name: "max", Fn: func(ctx *sql.Context, e sql.Expression) sql.Expression { return aggregation.NewMax(e) }},
	sql.Function1{Name: "min", Fn: func(ctx *sql.Context, e sql.Expression) sql.Expression { return aggregation.NewMin(e) }},
	sql.FunctionN{Name: "now", Fn: NewNow},
	sql.Function1{Name: "sum", Fn: func(ctx *sql.Context, e sql.Expression) sql.Expression { return aggregation.NewSum(e) }},
	sql.FunctionN{Name: "substr", Fn: NewSubstring},
	sql.FunctionN{Name: "substring", Fn: NewSubstring},
	sql.Function1{Name: "upper", Fn: NewUpper},
}

func GetLockingFuncs(ctx *sql.Context, ls *sql.LockSubsystem) []sql.Function {
	return nil
}

// Registry is used to register functions
type Registry map[string]sql.Function

var _ sql.FunctionProvider = Registry{}

// NewRegistry creates a new Registry.
func NewRegistry() Registry {
	fr := make(Registry)
	fr.mustRegister(BuiltIns...)
	return fr
}

// Register registers functions, returning an error if it's already registered
func (r Registry) Register(fn ...sql.Function) error {
	for _, f := range fn {
		if _, ok := r[f.FunctionName()]; ok {
			return ErrFunctionAlreadyRegistered.New(f.FunctionName())
		}
		r[f.FunctionName()] = f
	}
	return nil
}

// Function implements sql.FunctionProvider
func (r Registry) Function(ctx *sql.Context, name string) (sql.Function, bool) {
	if fn, ok := r[name]; ok {
		return fn, true
	}
	return nil, false
}

func (r Registry) mustRegister(fn ...sql.Function) {
	if err := r.Register(fn...); err != nil {
		panic(err)
	}
}
