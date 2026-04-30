//go:build !sql_lite

package iters

import "github.com/dolthub/jsonpath"

func jsonPathLookup(obj any, path string) (any, error) {
	return jsonpath.JsonPathLookup(obj, path)
}
