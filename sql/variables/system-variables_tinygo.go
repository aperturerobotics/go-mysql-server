//go:build tinygo

package variables

import (
	"maps"
	"strings"
	"sync"
	"time"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
)

// ServerStartUpTime is needed by uptime status variable.
var ServerStartUpTime = time.Now()

type globalSystemVariables struct {
	mutex      *sync.RWMutex
	sysVarVals map[string]sql.SystemVarValue
}

var _ sql.SystemVariableRegistry = (*globalSystemVariables)(nil)

func (sv *globalSystemVariables) AddSystemVariables(sysVars []sql.SystemVariable) {
	sv.mutex.Lock()
	defer sv.mutex.Unlock()
	for _, sysVar := range sysVars {
		if sysVar == nil {
			continue
		}
		name := strings.ToLower(sysVar.GetName())
		systemVars[name] = sysVar
		sv.sysVarVals[name] = sql.SystemVarValue{Var: sysVar, Val: sysVar.GetDefault()}
	}
}

func (sv *globalSystemVariables) AssignValues(vals map[string]any) error {
	ctx := sql.NewEmptyContext()
	sv.mutex.Lock()
	defer sv.mutex.Unlock()
	for name, val := range vals {
		name = strings.ToLower(name)
		sysVar, ok := systemVars[name]
		if !ok {
			return sql.ErrUnknownSystemVariable.New(name)
		}
		svv, err := sysVar.InitValue(ctx, val, true)
		if err != nil {
			return err
		}
		sv.sysVarVals[name] = svv
	}
	return nil
}

func (sv *globalSystemVariables) NewSessionMap() map[string]sql.SystemVarValue {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()
	sessionVals := make(map[string]sql.SystemVarValue, len(sv.sysVarVals))
	maps.Copy(sessionVals, sv.sysVarVals)
	return sessionVals
}

func (sv *globalSystemVariables) GetGlobal(name string) (sql.SystemVariable, any, bool) {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()
	name = strings.ToLower(name)
	sysVal, ok := sv.sysVarVals[name]
	if !ok {
		return nil, nil, false
	}
	return sysVal.Var, sysVal.Val, true
}

func (sv *globalSystemVariables) SetGlobal(ctx *sql.Context, name string, val any) error {
	sv.mutex.Lock()
	defer sv.mutex.Unlock()
	name = strings.ToLower(name)
	sysVar, ok := systemVars[name]
	if !ok {
		return sql.ErrUnknownSystemVariable.New(name)
	}
	svv, err := sysVar.SetValue(ctx, val, true)
	if err != nil {
		return err
	}
	sv.sysVarVals[name] = svv
	return nil
}

func (sv *globalSystemVariables) GetAllGlobalVariables() map[string]any {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()
	m := make(map[string]any, len(sv.sysVarVals))
	for k, v := range sv.sysVarVals {
		m[k] = v.Val
	}
	return m
}

func InitSystemVariables() {
	out := &globalSystemVariables{
		mutex:      &sync.RWMutex{},
		sysVarVals: make(map[string]sql.SystemVarValue, len(systemVars)),
	}
	for _, sysVar := range systemVars {
		out.sysVarVals[sysVar.GetName()] = sql.SystemVarValue{
			Var: sysVar,
			Val: sysVar.GetDefault(),
		}
	}
	sql.SystemVariables = out
}

func init() {
	InitSystemVariables()
}

var systemVars = map[string]sql.SystemVariable{
	"local_infile": &sql.MysqlSystemVariable{
		Name:    "local_infile",
		Scope:   sql.GetMysqlScope(sql.SystemVariableScope_Both),
		Dynamic: true,
		Type:    types.NewSystemBoolType("local_infile"),
		Default: int8(0),
	},
	"regexp_buffer_size": &sql.MysqlSystemVariable{
		Name:    "regexp_buffer_size",
		Scope:   sql.GetMysqlScope(sql.SystemVariableScope_Both),
		Dynamic: true,
		Type:    types.NewSystemIntType("regexp_buffer_size", 0, 0, false),
		Default: int64(8000000),
	},
	"secure_file_priv": &sql.MysqlSystemVariable{
		Name:    "secure_file_priv",
		Scope:   sql.GetMysqlScope(sql.SystemVariableScope_Both),
		Dynamic: true,
		Type:    types.NewSystemStringType("secure_file_priv"),
		Default: "",
	},
	"server_id": &sql.MysqlSystemVariable{
		Name:    "server_id",
		Scope:   sql.GetMysqlScope(sql.SystemVariableScope_Both),
		Dynamic: true,
		Type:    types.NewSystemUintType("server_id", 0, 0),
		Default: uint64(1),
	},
	"sql_mode": &sql.MysqlSystemVariable{
		Name:    "sql_mode",
		Scope:   sql.GetMysqlScope(sql.SystemVariableScope_Both),
		Dynamic: true,
		Type:    types.NewSystemStringType("sql_mode"),
		Default: "",
	},
}
