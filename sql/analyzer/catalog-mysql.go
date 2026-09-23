//go:build tinygo || (js && sql_lite)

package analyzer

type catalogMySQLDb interface {
	Enabled() bool
	AddRootAccount()
}

func catalogMySQLDbEnabled(db catalogMySQLDb) bool {
	return db != nil && db.Enabled()
}
