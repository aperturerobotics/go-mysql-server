//go:build tinygo

package rowexec

import (
	"math"
	"os"
)

func scannerMaxTokenSize() int {
	return math.MaxInt
}

func sameFile(a, b os.FileInfo) bool {
	return a.Name() == b.Name()
}
