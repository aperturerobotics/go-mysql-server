// Copyright 2026 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sql

import (
	"fmt"
	"time"

	"github.com/cockroachdb/apd/v3"
)

// ValueKind is the concrete Go value family returned by Type.Convert.
type ValueKind byte

const (
	ValueKindUnknown ValueKind = iota
	ValueKindNull
	ValueKindBool
	ValueKindInt8
	ValueKindInt16
	ValueKindInt32
	ValueKindInt64
	ValueKindUint8
	ValueKindUint16
	ValueKindUint32
	ValueKindUint64
	ValueKindFloat32
	ValueKindFloat64
	ValueKindString
	ValueKindBytes
	ValueKindDecimal
	ValueKindTime
	ValueKindTimespan
	ValueKindJSONWrapper
	ValueKindJSONDocument
	ValueKindGeometry
	ValueKindPoint
	ValueKindLineString
	ValueKindPolygon
	ValueKindMultiPoint
	ValueKindMultiLineString
	ValueKindMultiPolygon
	ValueKindGeomColl
	ValueKindTuple
)

func (k ValueKind) String() string {
	switch k {
	case ValueKindUnknown:
		return "<unknown>"
	case ValueKindNull:
		return "<nil>"
	case ValueKindBool:
		return "bool"
	case ValueKindInt8:
		return "int8"
	case ValueKindInt16:
		return "int16"
	case ValueKindInt32:
		return "int32"
	case ValueKindInt64:
		return "int64"
	case ValueKindUint8:
		return "uint8"
	case ValueKindUint16:
		return "uint16"
	case ValueKindUint32:
		return "uint32"
	case ValueKindUint64:
		return "uint64"
	case ValueKindFloat32:
		return "float32"
	case ValueKindFloat64:
		return "float64"
	case ValueKindString:
		return "string"
	case ValueKindBytes:
		return "[]byte"
	case ValueKindDecimal:
		return "*apd.Decimal"
	case ValueKindTime:
		return "time.Time"
	case ValueKindTimespan:
		return "types.Timespan"
	case ValueKindJSONWrapper:
		return "sql.JSONWrapper"
	case ValueKindJSONDocument:
		return "types.JSONDocument"
	case ValueKindGeometry:
		return "types.GeometryValue"
	case ValueKindPoint:
		return "types.Point"
	case ValueKindLineString:
		return "types.LineString"
	case ValueKindPolygon:
		return "types.Polygon"
	case ValueKindMultiPoint:
		return "types.MultiPoint"
	case ValueKindMultiLineString:
		return "types.MultiLineString"
	case ValueKindMultiPolygon:
		return "types.MultiPolygon"
	case ValueKindGeomColl:
		return "types.GeomColl"
	case ValueKindTuple:
		return "[]interface{}"
	default:
		return "<invalid>"
	}
}

// Matches reports whether k is assignable to expected under the old
// value-validation contract.
func (k ValueKind) Matches(expected ValueKind) bool {
	if k == expected || expected == ValueKindUnknown {
		return true
	}
	switch expected {
	case ValueKindJSONWrapper:
		return k == ValueKindJSONDocument
	case ValueKindGeometry:
		return k.isGeometry()
	default:
		return false
	}
}

func (k ValueKind) isGeometry() bool {
	switch k {
	case ValueKindPoint, ValueKindLineString, ValueKindPolygon, ValueKindMultiPoint,
		ValueKindMultiLineString, ValueKindMultiPolygon, ValueKindGeomColl:
		return true
	default:
		return false
	}
}

// ValueKindProvider is implemented by values whose concrete package owns their
// SQL value-kind mapping.
type ValueKindProvider interface {
	ValueKind() ValueKind
}

// ValueKindFor returns the ValueKind for the generic type parameter.
func ValueKindFor[T any]() ValueKind {
	var zero T
	return ValueKindOf(zero)
}

// ValueKindOf returns the ValueKind for v.
func ValueKindOf(v any) ValueKind {
	switch val := v.(type) {
	case nil:
		return ValueKindNull
	case ValueKindProvider:
		return val.ValueKind()
	case JSONWrapper:
		return ValueKindJSONWrapper
	case bool:
		return ValueKindBool
	case int8:
		return ValueKindInt8
	case int16:
		return ValueKindInt16
	case int32:
		return ValueKindInt32
	case int64:
		return ValueKindInt64
	case uint8:
		return ValueKindUint8
	case uint16:
		return ValueKindUint16
	case uint32:
		return ValueKindUint32
	case uint64:
		return ValueKindUint64
	case float32:
		return ValueKindFloat32
	case float64:
		return ValueKindFloat64
	case string:
		return ValueKindString
	case []byte:
		return ValueKindBytes
	case *apd.Decimal:
		return ValueKindDecimal
	case time.Time:
		return ValueKindTime
	case []any:
		return ValueKindTuple
	default:
		return ValueKindUnknown
	}
}

func TypeName(v any) string {
	return fmt.Sprintf("%T", v)
}
