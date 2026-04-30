//go:build !tinygo

package sql

import "runtime"

func usedProcessMemory() uint64 {
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	return s.HeapInuse + s.StackInuse
}
