//go:build sql_lite

package memory

import "github.com/dolthub/go-mysql-server/sql"

// ExternalStoredProcedures is empty in the lite profile, which omits reflected calls.
var ExternalStoredProcedures []sql.ExternalStoredProcedureDetails
