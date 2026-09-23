//go:build !sql_lite

package plan

import (
	"strconv"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression"
)

// ExtendVariadic returns a new procedure that has the variadic parameter extended to match the CALL's parameter count.
func (p *Procedure) ExtendVariadic(ctx *sql.Context, length int) *Procedure {
	if !p.HasVariadicParameter() {
		return p
	}
	np := *p
	body := p.ExternalProc.(*ExternalProcedure)
	newBody := *body
	np.ExternalProc = &newBody

	newParamDefinitions := make([]ProcedureParam, length)
	newParams := make([]*expression.ProcedureParam, length)
	if length < len(p.Params) {
		newParamDefinitions = p.Params[:len(p.Params)-1]
		newParams = body.Params[:len(body.Params)-1]
	} else {
		for i := range p.Params {
			newParamDefinitions[i] = p.Params[i]
			newParams[i] = body.Params[i]
		}
		if length >= len(p.Params) {
			variadicParam := p.Params[len(p.Params)-1]
			for i := len(p.Params); i < length; i++ {
				paramName := "A" + strconv.FormatInt(int64(i), 10)
				newParamDefinitions[i] = ProcedureParam{
					Direction: variadicParam.Direction,
					Name:      paramName,
					Type:      variadicParam.Type,
					Variadic:  variadicParam.Variadic,
				}
				newParams[i] = expression.NewProcedureParam(paramName, variadicParam.Type)
			}
		}
	}

	newBody.ParamDefinitions = newParamDefinitions
	newBody.Params = newParams
	np.Params = newParamDefinitions
	return &np
}

// IsExternal returns whether the stored procedure is external.
func (p *Procedure) IsExternal() bool {
	if _, ok := p.ExternalProc.(*ExternalProcedure); ok {
		return true
	}
	return false
}
