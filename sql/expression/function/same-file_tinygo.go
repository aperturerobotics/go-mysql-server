//go:build tinygo

package function

import "os"

func sameFile(a, b os.FileInfo) bool {
	return a.Name() == b.Name()
}
