//go:build sql_lite

package types

import (
	"context"
	"database/sql/driver"
	"io"
	"maps"
	"slices"

	"github.com/dolthub/vitess/go/sqltypes"
	"github.com/dolthub/vitess/go/vt/proto/query"

	"github.com/dolthub/go-mysql-server/sql"
)

const MaxJsonFieldByteLength = int64(1024) * int64(1024) * int64(1024)

var JSON sql.Type = JsonType{}

type JsonType struct{}

type JSONBytes interface {
	sql.JSONWrapper
	GetBytes(ctx context.Context) ([]byte, error)
}

type SearchableJSON interface {
	sql.JSONWrapper
	Lookup(ctx context.Context, path string) (sql.JSONWrapper, error)
}

type ComparableJSON interface {
	sql.JSONWrapper
	Compare(ctx context.Context, other interface{}) (int, error)
	JsonType(ctx context.Context) (string, error)
}

type JsonObject = map[string]interface{}
type JsonArray = []interface{}

type MutableJSON interface {
	sql.JSONWrapper
	Insert(context.Context, string, sql.JSONWrapper) (MutableJSON, bool, error)
	Remove(context.Context, string) (MutableJSON, bool, error)
	Set(context.Context, string, sql.JSONWrapper) (MutableJSON, bool, error)
	Replace(context.Context, string, sql.JSONWrapper) (MutableJSON, bool, error)
	ArrayInsert(context.Context, string, sql.JSONWrapper) (MutableJSON, bool, error)
	ArrayAppend(context.Context, string, sql.JSONWrapper) (MutableJSON, bool, error)
}

type JSONDocument struct {
	Val interface{}
}

func (JsonType) Compare(context.Context, interface{}, interface{}) (int, error) {
	return 0, sql.ErrUnsupportedFeature.New("json")
}

func (JsonType) Convert(context.Context, interface{}) (interface{}, sql.ConvertInRange, error) {
	return nil, sql.InRange, sql.ErrUnsupportedFeature.New("json")
}

func (JsonType) Equals(other sql.Type) bool {
	_, ok := other.(JsonType)
	return ok
}

func (JsonType) MaxTextResponseByteLength(*sql.Context) uint32 {
	return uint32(MaxJsonFieldByteLength) - 1
}

func (JsonType) Promote() sql.Type { return JSON }

func (JsonType) SQL(*sql.Context, []byte, interface{}) (sqltypes.Value, error) {
	return sqltypes.NULL, sql.ErrUnsupportedFeature.New("json")
}

func (JsonType) String() string { return "json" }

func (JsonType) Type() query.Type { return sqltypes.TypeJSON }

func (JsonType) ValueKind() sql.ValueKind { return sql.ValueKindJSONDocument }

func (JsonType) Zero() interface{} { return nil }

func (JsonType) CollationCoercibility(*sql.Context) (sql.CollationID, byte) {
	return sql.Collation_binary, 5
}

func (doc JSONDocument) ToInterface(context.Context) (interface{}, error) { return doc.Val, nil }

func (doc JSONDocument) Clone(context.Context) sql.JSONWrapper { return doc }

func (doc JSONDocument) Compare(context.Context, interface{}) (int, error) {
	return 0, sql.ErrUnsupportedFeature.New("json")
}

func (doc JSONDocument) JsonType(context.Context) (string, error) {
	return "", sql.ErrUnsupportedFeature.New("json")
}

func (doc JSONDocument) Lookup(context.Context, string) (sql.JSONWrapper, error) {
	return nil, sql.ErrUnsupportedFeature.New("json")
}

func (doc JSONDocument) Value() (driver.Value, error) {
	return nil, sql.ErrUnsupportedFeature.New("json")
}

func (doc JSONDocument) String() string { return "json" }

func (doc JSONDocument) ValueKind() sql.ValueKind { return sql.ValueKindJSONDocument }

func (doc JSONDocument) JSONString() (string, error) {
	return "", sql.ErrUnsupportedFeature.New("json")
}

func (doc JSONDocument) Insert(context.Context, string, sql.JSONWrapper) (MutableJSON, bool, error) {
	return nil, false, sql.ErrUnsupportedFeature.New("json")
}

func (doc JSONDocument) Remove(context.Context, string) (MutableJSON, bool, error) {
	return nil, false, sql.ErrUnsupportedFeature.New("json")
}

func (doc JSONDocument) Set(context.Context, string, sql.JSONWrapper) (MutableJSON, bool, error) {
	return nil, false, sql.ErrUnsupportedFeature.New("json")
}

func (doc JSONDocument) Replace(context.Context, string, sql.JSONWrapper) (MutableJSON, bool, error) {
	return nil, false, sql.ErrUnsupportedFeature.New("json")
}

func (doc JSONDocument) ArrayInsert(context.Context, string, sql.JSONWrapper) (MutableJSON, bool, error) {
	return nil, false, sql.ErrUnsupportedFeature.New("json")
}

func (doc JSONDocument) ArrayAppend(context.Context, string, sql.JSONWrapper) (MutableJSON, bool, error) {
	return nil, false, sql.ErrUnsupportedFeature.New("json")
}

func JsonToMySqlString(context.Context, sql.JSONWrapper) (string, error) {
	return "", sql.ErrUnsupportedFeature.New("json")
}

func JsonToMySqlBytes(context.Context, sql.JSONWrapper) ([]byte, error) {
	return nil, sql.ErrUnsupportedFeature.New("json")
}

func MarshallJson(context.Context, sql.JSONWrapper) ([]byte, error) {
	return nil, sql.ErrUnsupportedFeature.New("json")
}

func MarshallJsonValue(interface{}) ([]byte, error) {
	return nil, sql.ErrUnsupportedFeature.New("json")
}

func CompareJSON(context.Context, interface{}, interface{}) (int, error) {
	return 0, sql.ErrUnsupportedFeature.New("json")
}

func LookupJSONValue(context.Context, sql.JSONWrapper, string) (sql.JSONWrapper, error) {
	return nil, sql.ErrUnsupportedFeature.New("json")
}

func ConcatenateJSONValues(*sql.Context, ...sql.JSONWrapper) (sql.JSONWrapper, error) {
	return nil, sql.ErrUnsupportedFeature.New("json")
}

func ContainsJSON(interface{}, interface{}) (bool, error) {
	return false, sql.ErrUnsupportedFeature.New("json")
}

func DeepCopyJson(v interface{}) interface{} { return v }

func MustJSON(string) JSONDocument { return JSONDocument{} }

func NewLazyJSONDocument([]byte) sql.JSONWrapper { return JSONDocument{} }

type JSONIter struct {
	doc  *JsonObject
	keys []string
	idx  int
}

func NewJSONIter(json JsonObject) JSONIter {
	json = maps.Clone(json)
	keys := slices.Sorted(maps.Keys(json))
	return JSONIter{
		doc:  &json,
		keys: keys,
		idx:  0,
	}
}

func (iter *JSONIter) Next() (key string, value interface{}, err error) {
	if iter.idx >= len(iter.keys) {
		return "", nil, io.EOF
	}
	key = iter.keys[iter.idx]
	iter.idx++
	value = (*iter.doc)[key]
	return key, value, nil
}

func (iter *JSONIter) HasNext() bool {
	return iter.idx < len(iter.keys)
}

var _ sql.CollationCoercible = JsonType{}
var _ sql.JSONWrapper = JSONDocument{}
var _ MutableJSON = JSONDocument{}
var _ SearchableJSON = JSONDocument{}
