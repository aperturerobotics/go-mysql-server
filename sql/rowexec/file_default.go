//go:build !tinygo

package rowexec

import (
	"os"

	"github.com/dolthub/go-mysql-server/sql/types"
)

func scannerMaxTokenSize() int {
	return int(types.LongTextBlobMax)
}

func sameFile(a, b os.FileInfo) bool {
	return os.SameFile(a, b)
}
