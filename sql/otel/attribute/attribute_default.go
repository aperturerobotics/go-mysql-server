//go:build !tinygo && !js

package attribute

import (
	"fmt"

	otelattribute "go.opentelemetry.io/otel/attribute"
)

type KeyValue = otelattribute.KeyValue

func Bool(key string, value bool) KeyValue {
	return otelattribute.Bool(key, value)
}

func Int(key string, value int) KeyValue {
	return otelattribute.Int(key, value)
}

func Int64(key string, value int64) KeyValue {
	return otelattribute.Int64(key, value)
}

func String(key string, value string) KeyValue {
	return otelattribute.String(key, value)
}

func Stringer(key string, value fmt.Stringer) KeyValue {
	return otelattribute.Stringer(key, value)
}
