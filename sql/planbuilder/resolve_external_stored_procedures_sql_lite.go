//go:build sql_lite

package planbuilder

import (
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/plan"
)

func resolveExternalStoredProcedure(sql.ExternalStoredProcedureDetails) (*plan.Procedure, error) {
	return nil, sql.ErrUnsupportedFeature.New("external stored procedures")
}
