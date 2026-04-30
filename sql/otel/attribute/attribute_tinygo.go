//go:build tinygo || js

package attribute

import "fmt"

type KeyValue struct{}

func Bool(_ string, _ bool) KeyValue {
	return KeyValue{}
}

func Int(_ string, _ int) KeyValue {
	return KeyValue{}
}

func Int64(_ string, _ int64) KeyValue {
	return KeyValue{}
}

func String(_ string, _ string) KeyValue {
	return KeyValue{}
}

func Stringer(_ string, _ fmt.Stringer) KeyValue {
	return KeyValue{}
}
