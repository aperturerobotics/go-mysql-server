//go:build !sql_lite

package analyzer

import (
	"fmt"
	"reflect"

	"github.com/dolthub/go-mysql-server/sql"
)

// operatorOf retains native support for integrators with custom Operator types.
func operatorOf(e sql.Expression) (string, bool) {
	if o, ok := e.(interface{ Operator() string }); ok {
		return o.Operator(), true
	}
	m := reflect.ValueOf(e).MethodByName("Operator")
	if !m.IsValid() || m.Type().NumIn() != 0 || m.Type().NumOut() != 1 {
		return "", false
	}
	return fmt.Sprint(m.Call(nil)[0].Interface()), true
}
