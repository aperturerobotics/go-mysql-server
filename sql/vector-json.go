//go:build !sql_lite

package sql

import (
	"encoding/json"
	"fmt"
)

func convertJSONStringToVector(s string) ([]float32, error) {
	var val any
	err := json.Unmarshal([]byte(s), &val)
	if err != nil {
		return nil, fmt.Errorf("can't convert JSON to vector: %w", err)
	}
	return convertJsonInterfaceToVector(val)
}
