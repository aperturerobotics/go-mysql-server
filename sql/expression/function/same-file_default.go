//go:build !tinygo

package function

import (
	"os"
)

func sameFile(a, b os.FileInfo) bool {
	return os.SameFile(a, b)
}
