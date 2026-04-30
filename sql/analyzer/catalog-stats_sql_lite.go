//go:build sql_lite

package analyzer

import (
	"github.com/dolthub/go-mysql-server/sql"
)

type liteStatsProvider struct{}

func newStatsProvider() sql.StatsProvider {
	return liteStatsProvider{}
}

func (liteStatsProvider) GetTableStats(*sql.Context, string, sql.Table) ([]sql.Statistic, error) {
	return nil, nil
}

func (liteStatsProvider) AnalyzeTable(*sql.Context, sql.Table, string) error {
	return nil
}

func (liteStatsProvider) SetStats(*sql.Context, sql.Statistic) error {
	return nil
}

func (liteStatsProvider) GetStats(*sql.Context, sql.StatQualifier, []string) (sql.Statistic, bool) {
	return nil, false
}

func (liteStatsProvider) DropStats(*sql.Context, sql.StatQualifier, []string) error {
	return nil
}

func (liteStatsProvider) DropDbStats(*sql.Context, string, bool) error {
	return nil
}

func (liteStatsProvider) RowCount(*sql.Context, string, sql.Table) (uint64, error) {
	return 0, nil
}

func (liteStatsProvider) DataLength(*sql.Context, string, sql.Table) (uint64, error) {
	return 0, nil
}

var _ sql.StatsProvider = liteStatsProvider{}
