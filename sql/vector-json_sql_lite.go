//go:build sql_lite

package sql

func convertJSONStringToVector(string) ([]float32, error) {
	return nil, ErrUnsupportedFeature.New("json vector conversion")
}
