//go:build sql_lite

package similartext

// FindFromMap does the same as Find but taking a map instead
// of a string array as first argument.
func FindFromMap(_ any, _ string) string {
	return ""
}
