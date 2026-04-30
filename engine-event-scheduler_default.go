//go:build !tinygo

package sqle

import (
	"github.com/dolthub/go-mysql-server/eventscheduler"
	"github.com/dolthub/go-mysql-server/sql"
)

// InitializeEventScheduler initializes the EventScheduler for the engine.
func (e *Engine) InitializeEventScheduler(ctxGetterFunc func() (*sql.Context, error), status eventscheduler.SchedulerStatus, period int) error {
	var err error
	e.EventScheduler, err = eventscheduler.InitEventScheduler(e.Analyzer, e.BackgroundThreads, ctxGetterFunc, status, e.executeEvent, period)
	if err != nil {
		return err
	}
	return nil
}
