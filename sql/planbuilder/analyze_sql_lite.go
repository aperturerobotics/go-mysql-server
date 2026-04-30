//go:build sql_lite

package planbuilder

import ast "github.com/dolthub/vitess/go/vt/sqlparser"

import "github.com/dolthub/go-mysql-server/sql"

func (b *Builder) buildAnalyze(inScope *scope, _ *ast.Analyze, _ string) (outScope *scope) {
	b.handleErr(sql.ErrUnsupportedFeature.New("analyze"))
	return inScope
}
