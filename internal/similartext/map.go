//go:build !sql_lite

package similartext

// FindFromMap does the same as Find but taking a map instead
// of a string array as first argument.
func FindFromMap[V any](names map[string]V, src string) string {
	var namesList []string
	for name := range names {
		namesList = append(namesList, name)
	}

	return Find(namesList, src)
}
