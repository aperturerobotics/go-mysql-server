//go:build !tinygo

package planbuilder

import (
	"context"
	"runtime/trace"
)

type traceRegion interface {
	End()
}

func startTraceRegion(ctx context.Context, name string) traceRegion {
	return trace.StartRegion(ctx, name)
}
