// Copyright 2020-2021 Dolthub, Inc.
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

package driver

import (
	"database/sql/driver"

	"github.com/dolthub/vitess/go/vt/proto/query"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
)

// Rows is an iterator over an executed query's results.
type Rows struct {
	rows    sql.RowIter
	options *Options
	ctx     *sql.Context
	cols    sql.Schema
}

// Columns returns the names of the columns. The number of
// columns of the result is inferred from the length of the
// slice. If a particular column name isn't known, an empty
// string should be returned for that entry.
func (r *Rows) Columns() []string {
	names := make([]string, len(r.cols))
	for i, col := range r.cols {
		names[i] = col.Name
	}
	return names
}

// ColumnTypeDatabaseTypeName returns the MySQL type name of a column, in the
// form MySQL's own database/sql driver reports it, such as "BIGINT",
// "UNSIGNED INT", "DATETIME" or "JSON", without length or precision.
func (r *Rows) ColumnTypeDatabaseTypeName(index int) string {
	if name, ok := columnTypeNames[r.cols[index].Type.Type()]; ok {
		return name
	}
	return ""
}

// columnTypeNames maps each wire type to its MySQL type name.
var columnTypeNames = map[query.Type]string{
	query.Type_INT8:      "TINYINT",
	query.Type_UINT8:     "UNSIGNED TINYINT",
	query.Type_INT16:     "SMALLINT",
	query.Type_UINT16:    "UNSIGNED SMALLINT",
	query.Type_INT24:     "MEDIUMINT",
	query.Type_UINT24:    "UNSIGNED MEDIUMINT",
	query.Type_INT32:     "INT",
	query.Type_UINT32:    "UNSIGNED INT",
	query.Type_INT64:     "BIGINT",
	query.Type_UINT64:    "UNSIGNED BIGINT",
	query.Type_FLOAT32:   "FLOAT",
	query.Type_FLOAT64:   "DOUBLE",
	query.Type_DECIMAL:   "DECIMAL",
	query.Type_YEAR:      "YEAR",
	query.Type_DATE:      "DATE",
	query.Type_TIME:      "TIME",
	query.Type_DATETIME:  "DATETIME",
	query.Type_TIMESTAMP: "TIMESTAMP",
	query.Type_CHAR:      "CHAR",
	query.Type_VARCHAR:   "VARCHAR",
	query.Type_TEXT:      "TEXT",
	query.Type_BINARY:    "BINARY",
	query.Type_VARBINARY: "VARBINARY",
	query.Type_BLOB:      "BLOB",
	query.Type_BIT:       "BIT",
	query.Type_ENUM:      "ENUM",
	query.Type_SET:       "SET",
	query.Type_GEOMETRY:  "GEOMETRY",
	query.Type_JSON:      "JSON",
	query.Type_NULL_TYPE: "NULL",
}

// ColumnTypePrecisionScale returns a column's precision and scale as MySQL's
// own database/sql driver reports them: the digits and the digits after the
// point of a DECIMAL, and the fractional second digits, as both, of a
// DATETIME, TIMESTAMP or TIME. Other columns report false.
func (r *Rows) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {
	switch typ := r.cols[index].Type.(type) {
	case sql.DecimalType:
		return int64(typ.Precision()), int64(typ.Scale()), true
	case sql.DatetimeType:
		return int64(typ.Precision()), int64(typ.Precision()), true
	case types.TimeType:
		return int64(typ.Precision()), int64(typ.Precision()), true
	}
	return 0, 0, false
}

// Close closes the rows iterator.
func (r *Rows) Close() error {
	return r.rows.Close(r.ctx)
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide.
//
// Next should return io.EOF when there are no more rows.
//
// The dest should not be written to outside of Next. Care
// should be taken when closing Rows not to modify
// a buffer held in dest.
func (r *Rows) Next(dest []driver.Value) error {
again:
	row, err := r.rows.Next(r.ctx)
	if err != nil {
		return err
	}
	if len(row) == 0 {
		return nil
	}

	if _, ok := row[0].(types.OkResult); ok {
		// skip OK results
		goto again
	}

	for i := range row {
		dest[i] = r.convert(i, row[i])
	}
	return nil
}

// convert returns v as one of the value types database/sql accepts from a
// driver. An ENUM, SET, DECIMAL or TIME, which the engine holds as an index,
// bits, *apd.Decimal or microseconds, becomes the text MySQL sends for it.
func (r *Rows) convert(col int, v driver.Value) any {
	switch r.cols[col].Type.Type() {
	case query.Type_NULL_TYPE:
		return nil

	case query.Type_ENUM, query.Type_SET, query.Type_DECIMAL, query.Type_TIME:
		if v == nil {
			return nil
		}
		if text, err := r.cols[col].Type.SQL(r.ctx, nil, v); err == nil {
			return text.ToString()
		}

	case query.Type_INT8, query.Type_INT16, query.Type_INT24, query.Type_INT32, query.Type_INT64,
		query.Type_UINT8, query.Type_UINT16, query.Type_UINT24, query.Type_UINT32, query.Type_UINT64,
		query.Type_FLOAT32, query.Type_FLOAT64:
		switch val := v.(type) {
		case int:
			return int64(val)
		case int8:
			return int64(val)
		case int16:
			return int64(val)
		case int32:
			return int64(val)
		case int64:
			return val
		case uint:
			return uint64(val)
		case uint8:
			return uint64(val)
		case uint16:
			return uint64(val)
		case uint32:
			return uint64(val)
		case uint64:
			return val
		case float32:
			return float64(val)
		case float64:
			return val
		case string:
			return val
		}

	case query.Type_JSON:
		var asObj, asStr bool
		if r.options != nil {
			switch r.options.JSON {
			case ScanAsObject:
				asObj = true
			case ScanAsString:
				asStr = true
			case ScanAsBytes:
				// nothing to do
			default: // unknown or ScanAsStored
				return v
			}
		}

		sqlValue, _, err := r.cols[col].Type.Convert(r.ctx, v)
		if err != nil {
			break
		}
		doc, ok := sqlValue.(types.JSONDocument)
		if !ok {
			break
		}
		if asObj {
			return doc.Val
		}

		str, err := doc.JSONString()
		if err != nil {
			break
		}

		if asStr {
			return str
		}
		return []byte(str)
	}
	return v
}

var (
	_ driver.RowsColumnTypeDatabaseTypeName = (*Rows)(nil)
	_ driver.RowsColumnTypePrecisionScale   = (*Rows)(nil)
)
