//go:build tinygo

package sql

import "context"

type traceRegion struct{}

func (traceRegion) End() {}

func startTraceRegion(_ context.Context, _ string) traceRegion {
	return traceRegion{}
}
