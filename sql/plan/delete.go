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

package plan

import (
	"fmt"
	"gopkg.in/src-d/go-errors.v1"

	"github.com/dolthub/go-mysql-server/sql"
)

var ErrDeleteFromNotSupported = errors.NewKind("table doesn't support DELETE FROM")

// DeleteFrom is a node describing a deletion from some table.
type DeleteFrom struct {
	UnaryNode
	Targets []sql.Node
}

var _ sql.Databaseable = (*DeleteFrom)(nil)

// NewDeleteFrom creates a DeleteFrom node.
func NewDeleteFrom(n sql.Node, targets []sql.Node) *DeleteFrom {
	return &DeleteFrom{
		UnaryNode: UnaryNode{n},
		Targets:   targets,
	}
}

func GetDeletable(node sql.Node) (sql.DeletableTable, error) {
	switch node := node.(type) {
	case sql.DeletableTable:
		return node, nil
	case *IndexedTableAccess:
		return GetDeletable(node.ResolvedTable)
	case *ResolvedTable:
		return getDeletableTable(node.Table)
	case *SubqueryAlias:
		return nil, ErrDeleteFromNotSupported.New()
	case *TriggerExecutor:
		return GetDeletable(node.Left())
	case sql.TableWrapper:
		return getDeletableTable(node.Underlying())
	}
	if len(node.Children()) > 1 {
		return nil, ErrDeleteFromNotSupported.New()
	}
	for _, child := range node.Children() {
		deleter, _ := GetDeletable(child)
		if deleter != nil {
			return deleter, nil
		}
	}
	return nil, ErrDeleteFromNotSupported.New()
}

func getDeletableTable(t sql.Table) (sql.DeletableTable, error) {
	switch t := t.(type) {
	case sql.DeletableTable:
		return t, nil
	case sql.TableWrapper:
		return getDeletableTable(t.Underlying())
	default:
		return nil, ErrDeleteFromNotSupported.New()
	}
}

func deleteDatabaseHelper(node sql.Node) string {
	switch node := node.(type) {
	case sql.DeletableTable:
		return ""
	case *IndexedTableAccess:
		return deleteDatabaseHelper(node.ResolvedTable)
	case *ResolvedTable:
		return node.Database.Name()
	case *UnresolvedTable:
		return node.Database()
	}

	for _, child := range node.Children() {
		return deleteDatabaseHelper(child)
	}

	return ""
}

// WithTargets returns a new instance of this DeleteFrom, with the specified |targets|.
func (p *DeleteFrom) WithTargets(targets []sql.Node) sql.Node {
	newDeleteFrom := *p
	newDeleteFrom.Targets = targets
	return &newDeleteFrom
}

// Resolved implements the sql.Resolvable interface.
func (p *DeleteFrom) Resolved() bool {
	if p.Child.Resolved() == false {
		return false
	}

	if len(p.Targets) > 0 {
		for _, target := range p.Targets {
			if target.Resolved() == false {
				return false
			}
		}
	}

	return true
}

func (p *DeleteFrom) Database() string {
	return deleteDatabaseHelper(p.Child)
}

// RowIter implements the Node interface.
func (p *DeleteFrom) RowIter(ctx *sql.Context, row sql.Row) (sql.RowIter, error) {
	// If an empty table is passed in (potentially from a bad filter) return an empty row iter.
	// Note: emptyTable could also implement sql.DetetableTable
	if _, ok := p.Child.(*emptyTable); ok {
		return sql.RowsToRowIter(), nil
	}

	iter, err := p.Child.RowIter(ctx, row)
	if err != nil {
		return nil, err
	}

	if len(p.Targets) == 0 {
		deletable, err := GetDeletable(p.Child)
		if err != nil {
			return nil, err
		}
		return newDeleteIter(iter, deletable.Schema(), deleterStruct{deletable.Deleter(ctx), 0, len(deletable.Schema())}), nil
	} else {
		// TODO: Validate table wasn't specified twice? validate no multi-db?
		deleterStructs := make([]deleterStruct, len(p.Targets))
		for i, target := range p.Targets {
			deletable, err := GetDeletable(target)
			if err != nil {
				return nil, err
			}
			deleter := deletable.Deleter(ctx)
			start, end := findPosition(p.Child.Schema(), deletable.Name())
			deleterStructs[i] = deleterStruct{deleter, int(start), int(end)}
		}
		return newDeleteIter(iter, p.Child.Schema(), deleterStructs...), nil

	}

}

// TODO: Rename and document
type deleterStruct struct {
	deleter     sql.RowDeleter
	schemaStart int
	schemaEnd   int
}

func findPosition(schema sql.Schema, name string) (uint, uint) {
	foundStart := false
	var start uint
	for i, col := range schema {
		// casing?
		if col.Source == name {
			if !foundStart {
				start = uint(i)
				foundStart = true
			}
		} else {
			if foundStart {
				return start, uint(i)
			}
		}
	}
	if foundStart {
		return start, uint(len(schema))
	}

	// TODO: Turn this into an error
	panic(fmt.Sprintf("didn't find table %s", name))
}

type deleteIter struct {
	deleters  []deleterStruct
	schema    sql.Schema
	childIter sql.RowIter
	closed    bool
}

func (d *deleteIter) Next(ctx *sql.Context) (sql.Row, error) {
	row, err := d.childIter.Next(ctx)
	if err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Delete iterator receives rows from its source, which could be a table, a
	// subquery(?), a join statement, etc. Any filtering is done as a layer between the source
	// and the delete node. The delete node simply reads in all of it's source and calls delete
	// on the deletable interface. (If no explicit target is given... the delete is limited to
	// having a single table as it's source and that is implicitly the target.) The delete interface
	// accepts the full row to delete, so it is critical that the delete node extract the
	// correct columns from the context row!

	// Now that we are supporting multiple targets (actually... multiple sources even more, since
	// they generate the extra columns in the child row!) we need to give the delete iterator
	// a set of: DeletableTable and the position in the row from where to pull the PK to delete
	// from that table.

	// TODO: Add subquery test cases
	// TODO: re-read MySQL docs for any other edge cases
	//       https://dev.mysql.com/doc/refman/8.0/en/delete.html

	// Reduce the row to the length of the schema. The length can differ when some update values come from an outer
	// scope, which will be the first N values in the row.
	// TODO: handle this in the analyzer instead?
	fullSchemaLength := len(d.schema)
	rowLength := len(row)
	for _, deleter := range d.deleters {
		schemaLength := deleter.schemaEnd - deleter.schemaStart
		subSlice := row
		if schemaLength < rowLength {
			subSlice = row[(rowLength - fullSchemaLength + deleter.schemaStart):(rowLength - fullSchemaLength + deleter.schemaEnd)]
		}
		err = deleter.deleter.Delete(ctx, subSlice)
		if err != nil {
			return nil, err
		}
	}

	return row, nil
}

func (d *deleteIter) Close(ctx *sql.Context) error {
	if !d.closed {
		d.closed = true
		for _, deleter := range d.deleters {
			// TODO: collect errs?
			if err := deleter.deleter.Close(ctx); err != nil {
				return err
			}
		}
		return d.childIter.Close(ctx)
	}
	return nil
}

func newDeleteIter(childIter sql.RowIter, schema sql.Schema, deleters ...deleterStruct) sql.RowIter {
	openerClosers := make([]sql.EditOpenerCloser, len(deleters))
	for i, ds := range deleters {
		openerClosers[i] = ds.deleter
	}
	return NewTableEditorIter(&deleteIter{
		deleters:  deleters,
		childIter: childIter,
		schema:    schema,
	}, openerClosers...)
}

// WithChildren implements the Node interface.
func (p *DeleteFrom) WithChildren(children ...sql.Node) (sql.Node, error) {
	if len(children) != 1 {
		return nil, sql.ErrInvalidChildrenNumber.New(p, len(children), 1)
	}
	return NewDeleteFrom(children[0], p.Targets), nil
}

// CheckPrivileges implements the interface sql.Node.
func (p *DeleteFrom) CheckPrivileges(ctx *sql.Context, opChecker sql.PrivilegedOperationChecker) bool {
	//TODO: If column values are retrieved then the SELECT privilege is required
	// For example: "DELETE FROM table WHERE z > 0"
	// We would need SELECT privileges on the "z" column as it's retrieving values
	return opChecker.UserHasPrivileges(ctx,
		sql.NewPrivilegedOperation(p.Database(), getTableName(p.Child), "", sql.PrivilegeType_Delete))
}

func (p DeleteFrom) String() string {
	pr := sql.NewTreePrinter()
	_ = pr.WriteNode("Delete")
	_ = pr.WriteChildren(p.Child.String())
	return pr.String()
}

func (p DeleteFrom) DebugString() string {
	pr := sql.NewTreePrinter()
	_ = pr.WriteNode("Delete")
	_ = pr.WriteChildren(sql.DebugString(p.Child))
	return pr.String()
}
