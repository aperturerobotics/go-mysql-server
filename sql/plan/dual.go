// Copyright 2022 Dolthub, Inc.
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

package plan

import (
	"io"
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
)

// DualTableName is empty string because no table with empty name can be created
const DualTableName = ""

var dualTableSchema = sql.Schema{}

type dualTable struct{}

// IsDualTable returns whether the given table is the "dual" table.
func IsDualTable(t sql.Table) bool {
	if t == nil {
		return false
	}
	return strings.ToLower(t.Name()) == DualTableName && IsDualSchema(t.Schema(sql.NewEmptyContext()))
}

// IsDualSchema returns whether the given schema is the "dual" table schema.
func IsDualSchema(s sql.Schema) bool {
	return s.Equals(dualTableSchema)
}

var dualTableInstance sql.Table = dualTable{}

// NewDualSqlTable creates a new Dual table.
func NewDualSqlTable() sql.Table {
	return dualTableInstance
}

func (dualTable) Name() string { return DualTableName }

func (dualTable) String() string { return "dual" }

func (dualTable) Schema(*sql.Context) sql.Schema { return dualTableSchema }

func (dualTable) Collation() sql.CollationID { return sql.Collation_Default }

func (dualTable) Partitions(*sql.Context) (sql.PartitionIter, error) {
	return &dualPartitionIter{}, nil
}

func (dualTable) PartitionRows(*sql.Context, sql.Partition) (sql.RowIter, error) {
	return &dualRowIter{}, nil
}

type dualPartition struct{}

func (dualPartition) Key() []byte { return nil }

type dualPartitionIter struct {
	done bool
}

func (i *dualPartitionIter) Next(*sql.Context) (sql.Partition, error) {
	if i.done {
		return nil, io.EOF
	}
	i.done = true
	return dualPartition{}, nil
}

func (*dualPartitionIter) Close(*sql.Context) error { return nil }

type dualRowIter struct {
	done bool
}

func (i *dualRowIter) Next(*sql.Context) (sql.Row, error) {
	if i.done {
		return nil, io.EOF
	}
	i.done = true
	return sql.Row{}, nil
}

func (*dualRowIter) Close(*sql.Context) error { return nil }

var _ sql.Table = dualTable{}
