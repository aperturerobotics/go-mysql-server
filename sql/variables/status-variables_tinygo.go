//go:build tinygo

package variables

import (
	"sync/atomic"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
)

type globalStatusVariables struct {
	varVals map[string]sql.StatusVarValue
}

var _ sql.StatusVariableRegistry = (*globalStatusVariables)(nil)

func (g *globalStatusVariables) NewSessionMap() map[string]sql.StatusVarValue {
	return map[string]sql.StatusVarValue{}
}

func (g *globalStatusVariables) NewGlobalMap() map[string]sql.StatusVarValue {
	globalMap := make(map[string]sql.StatusVarValue, len(g.varVals))
	for k, v := range g.varVals {
		globalMap[k] = v.Copy()
	}
	return globalMap
}

func (g *globalStatusVariables) GetGlobal(name string) (sql.StatusVariable, any, bool) {
	v, ok := g.varVals[name]
	if !ok {
		return nil, nil, false
	}
	return v.Variable(), v.Value(), true
}

func (g *globalStatusVariables) SetGlobal(name string, val any) error {
	v, ok := g.varVals[name]
	if !ok {
		return sql.ErrUnknownStatusVariable.New(name)
	}
	return v.Set(val)
}

func (g *globalStatusVariables) IncrementGlobal(name string, val int) {
	v, ok := g.varVals[name]
	if !ok {
		return
	}
	v.Increment(uint64(val))
}

func InitStatusVariables() {
	globalVars := &globalStatusVariables{
		varVals: make(map[string]sql.StatusVarValue, len(statusVars)),
	}
	for _, statusVar := range statusVars {
		switch statusVar.GetDefault().(type) {
		case atomic.Uint64:
			globalVars.varVals[statusVar.GetName()] = &sql.MutableStatusVarValue{
				Var: statusVar,
				Val: &atomic.Uint64{},
			}
		default:
			globalVars.varVals[statusVar.GetName()] = &sql.ImmutableStatusVarValue{
				Var: statusVar,
				Val: statusVar.GetDefault(),
			}
		}
	}
	sql.StatusVariables = globalVars
}

func init() {
	InitStatusVariables()
}

var statusVars = map[string]sql.StatusVariable{
	"Bytes_received": &sql.MySQLStatusVariable{
		Name:    "Bytes_received",
		Scope:   sql.StatusVariableScope_Global,
		Type:    types.NewSystemIntType("Bytes_received", 0, 0, false),
		Default: atomic.Uint64{},
	},
	"Bytes_sent": &sql.MySQLStatusVariable{
		Name:    "Bytes_sent",
		Scope:   sql.StatusVariableScope_Global,
		Type:    types.NewSystemIntType("Bytes_sent", 0, 0, false),
		Default: atomic.Uint64{},
	},
}
