//go:build !sql_lite

package analyzer

import (
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/sql"
)

func newStatsProvider() sql.StatsProvider {
	return memory.NewStatsProv()
}
