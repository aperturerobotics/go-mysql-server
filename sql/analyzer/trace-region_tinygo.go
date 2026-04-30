//go:build tinygo

package analyzer

import "context"

type traceRegion struct{}

func (traceRegion) End() {}

func startTraceRegion(_ context.Context, _ string) traceRegion {
	return traceRegion{}
}
