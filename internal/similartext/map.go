//go:build !sql_lite

package similartext

import "reflect"

// FindFromMap does the same as Find but taking a map instead
// of a string array as first argument.
func FindFromMap(names any, src string) string {
	rnames := reflect.ValueOf(names)
	if rnames.Kind() != reflect.Map {
		panic("Implementation error: non map used as first argument " +
			"to FindFromMap")
	}

	t := rnames.Type()
	if t.Key().Kind() != reflect.String {
		panic("Implementation error: non string key for map used as " +
			"first argument to FindFromMap")
	}

	var namesList []string
	for _, kv := range rnames.MapKeys() {
		namesList = append(namesList, kv.String())
	}

	return Find(namesList, src)
}
