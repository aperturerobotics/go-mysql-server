package enginetest_test

import (
	"os"
	"testing"

	"github.com/dolthub/go-mysql-server/sql/variables"
)

// TestMain initializes status variables for tests that open components directly.
// These tests must not depend on an earlier test having constructed an engine.
func TestMain(m *testing.M) {
	variables.InitStatusVariables()
	os.Exit(m.Run())
}
