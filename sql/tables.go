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

package sql

import (
	"fmt"
)

// Table is a SQL table.
type Table interface {
	Nameable
	fmt.Stringer
	// Schema returns the table's schema.
	Schema() Schema
	// Collation returns the table's collation.
	Collation() CollationID
	// Partitions returns the table's partitions in an iterator.
	Partitions(*Context) (PartitionIter, error)
	// PartitionRows returns the rows in the given partition, which was returned by Partitions.
	PartitionRows(*Context, Partition) (RowIter, error)
}

// TableFunction is a node that is generated by a function and can be used as a table factor in many SQL queries.
type TableFunction interface {
	Node
	Expressioner
	Databaser
	Nameable

	// NewInstance calls the table function with the arguments provided, producing a Node
	NewInstance(ctx *Context, db Database, args []Expression) (Node, error)
}

// CatalogTableFunction is a table function that can be used as a table factor in many SQL queries.
type CatalogTableFunction interface {
	TableFunction

	// WithCatalog returns a new instance of the table function with the given catalog
	WithCatalog(c Catalog) (TableFunction, error)
}

// TemporaryTable allows tables to declare that they are temporary (created by CREATE TEMPORARY TABLE).
// Only used for validation of certain DDL operations -- in almost all respects TemporaryTables are indistinguishable
// from persisted tables to the engine.
type TemporaryTable interface {
	// IsTemporary should return true if the table is temporary to the session
	IsTemporary() bool
}

// TableWrapper is a node that wraps the real table. This is needed because wrappers cannot implement some methods the
// table may implement. This interface is used in analysis and planning and is not expected to be implemented by
// integrators.
type TableWrapper interface {
	// Underlying returns the underlying table.
	Underlying() Table
}

// MutableTableWrapper is a TableWrapper that can change its underlying table.
type MutableTableWrapper interface {
	TableWrapper
	WithUnderlying(Table) Table
}

// FilteredTable is a table that can filter its result rows from RowIter using filter expressions that would otherwise
// be applied by a separate Filter node.
type FilteredTable interface {
	Table
	// Filters returns the filter expressions that have been applied to this table.
	Filters() []Expression
	// HandledFilters returns the subset of the filter expressions given that this table can apply.
	HandledFilters(filters []Expression) []Expression
	// WithFilters returns a table with the given filter expressions applied.
	WithFilters(ctx *Context, filters []Expression) Table
}

// CommentedTable is a table that has a comment on it.
type CommentedTable interface {
	Table
	// Comment returns the table's optional comment.
	Comment() string
}

// ProjectedTable is a table that can return only a subset of its columns from RowIter. This provides a very large
// efficiency gain during table scans. Tables that implement this interface must return only the projected columns
// in future calls to Schema.
type ProjectedTable interface {
	Table
	// WithProjections returns a version of this table with only the subset of columns named. Calls to Schema must
	// only include these columns. A zero-length slice of column names is valid and indicates that rows from this table
	// should be spooled, but no columns should be returned. A nil slice will never be provided.
	WithProjections(colNames []string) Table
	// Projections returns the names of the column projections applied to this table, or nil if no projection is applied
	// and all columns of the schema will be returned.
	Projections() []string
}

// IndexAddressable is a table that can be scanned through a primary index
type IndexAddressable interface {
	// IndexedAccess returns a table that can perform scans constrained to
	// an IndexLookup on the index given, or nil if the index cannot support
	// the lookup expression.
	IndexedAccess(lookup IndexLookup) IndexedTable
	// GetIndexes returns an array of this table's Indexes
	GetIndexes(ctx *Context) ([]Index, error)
	// PreciseMatch returns whether an indexed access can substitute for filters
	PreciseMatch() bool
}

// IndexRequired tables cannot be executed without index lookups on certain
// columns. Join planning uses this interface to maintain plan correctness
// for these nodes
type IndexRequired interface {
	IndexAddressable
	// RequiredPredicates returns a list of columns that need IndexedTableAccess
	RequiredPredicates() []string
}

// IndexAddressableTable is a table that can be accessed through an index
type IndexAddressableTable interface {
	Table
	IndexAddressable
}

// IndexedTable is a table with an index chosen for range scans
type IndexedTable interface {
	Table
	// LookupPartitions returns partitions scanned by the given IndexLookup
	LookupPartitions(*Context, IndexLookup) (PartitionIter, error)
}

// IndexAlterableTable represents a table that supports index modification operations.
type IndexAlterableTable interface {
	Table
	// CreateIndex creates an index for this table, using the provided parameters.
	// Returns an error if the index name already exists, or an index with the same columns already exists.
	CreateIndex(ctx *Context, indexDef IndexDef) error
	// DropIndex removes an index from this table, if it exists.
	// Returns an error if the removal failed or the index does not exist.
	DropIndex(ctx *Context, indexName string) error
	// RenameIndex renames an existing index to another name that is not already taken by another index on this table.
	RenameIndex(ctx *Context, fromIndexName string, toIndexName string) error
}

// IndexBuildingTable is an optional extension to IndexAlterableTable that supports the engine's assistance in building
// a newly created index, or rebuilding an existing one. This interface is non-optional for tables that wish to create
// indexes on virtual columns, as the engine must provide a value for these columns.
type IndexBuildingTable interface {
	IndexAlterableTable
	// ShouldBuildIndex returns whether the given index should be build via BuildIndex. Some indexes require building,
	// in which case this method is not called.
	ShouldBuildIndex(ctx *Context, indexDef IndexDef) (bool, error)
	// BuildIndex returns a RowInserter for that will be passed all existing rows of the table. The returned RowInserter
	// should use the rows provided to populate the newly created index given by the definition. When |Close| is called
	// on the RowInserter, the index should be fully populated and available for further use in the session.
	BuildIndex(ctx *Context, indexDef IndexDef) (RowInserter, error)
}

// ForeignKeyTable is a table that declares foreign key constraints, and can be referenced by other tables' foreign
// key constraints.
type ForeignKeyTable interface {
	IndexAddressableTable
	// CreateIndexForForeignKey creates an index for this table, using the provided parameters. Indexes created through
	// this function are specifically ones generated for use with a foreign key. Returns an error if the index name
	// already exists, or an index on the same columns already exists.
	CreateIndexForForeignKey(ctx *Context, indexDef IndexDef) error
	// GetDeclaredForeignKeys returns the foreign key constraints that are declared by this table.
	GetDeclaredForeignKeys(ctx *Context) ([]ForeignKeyConstraint, error)
	// GetReferencedForeignKeys returns the foreign key constraints that are referenced by this table.
	GetReferencedForeignKeys(ctx *Context) ([]ForeignKeyConstraint, error)
	// AddForeignKey adds the given foreign key constraint to the table. Returns an error if the foreign key name
	// already exists on any other table within the database.
	AddForeignKey(ctx *Context, fk ForeignKeyConstraint) error
	// DropForeignKey removes a foreign key from the table.
	DropForeignKey(ctx *Context, fkName string) error
	// UpdateForeignKey updates the given foreign key constraint. May range from updated table names to setting the
	// IsResolved boolean.
	UpdateForeignKey(ctx *Context, fkName string, fk ForeignKeyConstraint) error
	// GetForeignKeyEditor returns a ForeignKeyEditor for this table.
	GetForeignKeyEditor(ctx *Context) ForeignKeyEditor
}

// ForeignKeyEditor is a TableEditor that is addressable via IndexLookup.
type ForeignKeyEditor interface {
	TableEditor
	IndexAddressable
}

// ReferenceChecker is usually an IndexAddressableTable that does key
// lookups for existence checks. Indicating that the engine is performing
// a reference check lets the integrator avoid expensive deserialization
// steps.
type ReferenceChecker interface {
	SetReferenceCheck() error
}

// CheckTable is a table that declares check constraints.
type CheckTable interface {
	Table
	// GetChecks returns the check constraints on this table.
	GetChecks(ctx *Context) ([]CheckDefinition, error)
}

// CheckAlterableTable represents a table that supports check constraints.
type CheckAlterableTable interface {
	Table
	// CreateCheck creates an check constraint for this table, using the provided parameters.
	// Returns an error if the constraint name already exists.
	CreateCheck(ctx *Context, check *CheckDefinition) error
	// DropCheck removes a check constraint from the database.
	DropCheck(ctx *Context, chName string) error
}

// CollationAlterableTable represents a table that supports the alteration of its collation.
type CollationAlterableTable interface {
	Table
	// ModifyStoredCollation modifies the default collation that is set on the table, along with converting all columns
	// to the given collation (ALTER TABLE ... CONVERT TO CHARACTER SET).
	ModifyStoredCollation(ctx *Context, collation CollationID) error
	// ModifyDefaultCollation modifies the default collation that is set on the table (ALTER TABLE ... COLLATE).
	ModifyDefaultCollation(ctx *Context, collation CollationID) error
}

// PrimaryKeyTable is a table with a primary key.
type PrimaryKeyTable interface {
	// PrimaryKeySchema returns this table's PrimaryKeySchema
	PrimaryKeySchema() PrimaryKeySchema
}

// PrimaryKeyAlterableTable represents a table that supports primary key changes.
type PrimaryKeyAlterableTable interface {
	Table
	// CreatePrimaryKey creates a primary key for this table, using the provided parameters.
	// Returns an error if the new primary key set is not compatible with the current table data.
	CreatePrimaryKey(ctx *Context, columns []IndexColumn) error
	// DropPrimaryKey drops a primary key on a table. Returns an error if that table does not have a key.
	DropPrimaryKey(ctx *Context) error
}

// EditOpenerCloser is the base interface for table editors, and deals with statement boundaries.
type EditOpenerCloser interface {
	// StatementBegin is called before the first operation of a statement. Integrators should mark the state of the data
	// in some way that it may be returned to in the case of an error.
	StatementBegin(ctx *Context)
	// DiscardChanges is called if a statement encounters an error, and all current changes since the statement beginning
	// should be discarded.
	DiscardChanges(ctx *Context, errorEncountered error) error
	// StatementComplete is called after the last operation of the statement, indicating that it has successfully completed.
	// The mark set in StatementBegin may be removed, and a new one should be created on the next StatementBegin.
	StatementComplete(ctx *Context) error
}

// InsertableTable is a table that can process insertion of new rows.
type InsertableTable interface {
	Table
	// Inserter returns an Inserter for this table. The Inserter will get one call to Insert() for each row to be
	// inserted, and will end with a call to Close() to finalize the insert operation.
	Inserter(*Context) RowInserter
}

// RowInserter is an insert cursor that can insert one or more values to a table.
type RowInserter interface {
	EditOpenerCloser
	// Insert inserts the row given, returning an error if it cannot. Insert will be called once for each row to process
	// for the insert operation, which may involve many rows. After all rows in an operation have been processed, Close
	// is called.
	Insert(*Context, Row) error
	// Close finalizes the insert operation, persisting its result.
	Closer
}

// DeletableTable is a table that can delete rows.
type DeletableTable interface {
	Table
	// Deleter returns a RowDeleter for this table. The RowDeleter will get one call to Delete for each row to be deleted,
	// and will end with a call to Close() to finalize the delete operation.
	Deleter(*Context) RowDeleter
}

// RowDeleter is a delete cursor that can delete one or more rows from a table.
type RowDeleter interface {
	EditOpenerCloser
	// Delete deletes the given row. Returns ErrDeleteRowNotFound if the row was not found. Delete will be called once for
	// each row to process for the delete operation, which may involve many rows. After all rows have been processed,
	// Close is called.
	Delete(*Context, Row) error
	// Closer finalizes the delete operation, persisting the result.
	Closer
}

// TruncateableTable is a table that can process the deletion of all rows either via a TRUNCATE TABLE statement or a
// DELETE statement without a WHERE clause. This is usually much faster that deleting rows one at a time.
type TruncateableTable interface {
	Table
	// Truncate removes all rows from the table. If the table also implements DeletableTable and it is determined that
	// truncate would be equivalent to a DELETE which spans the entire table, then this function will be called instead.
	// Returns the number of rows that were removed.
	Truncate(*Context) (int, error)
}

// AutoIncrementTable is a table that supports AUTO_INCREMENT. Getter and Setter methods access the table's
// AUTO_INCREMENT sequence. These methods should only be used for tables with and AUTO_INCREMENT column in their schema.
type AutoIncrementTable interface {
	Table
	// PeekNextAutoIncrementValue returns the next AUTO_INCREMENT value without incrementing the current
	// auto_increment counter.
	PeekNextAutoIncrementValue(ctx *Context) (uint64, error)
	// GetNextAutoIncrementValue gets the next AUTO_INCREMENT value. In the case that a table with an autoincrement
	// column is passed in a row with the autoinc column failed, the next auto increment value must
	// update its internal state accordingly and use the insert val at runtime.
	// Implementations are responsible for updating their state to provide the correct values.
	GetNextAutoIncrementValue(ctx *Context, insertVal interface{}) (uint64, error)
	// AutoIncrementSetter returns an AutoIncrementSetter.
	AutoIncrementSetter(*Context) AutoIncrementSetter
}

// AutoIncrementSetter provides support for altering a table's
// AUTO_INCREMENT sequence, eg 'ALTER TABLE t AUTO_INCREMENT = 10;'
type AutoIncrementSetter interface {
	// SetAutoIncrementValue sets a new AUTO_INCREMENT value.
	SetAutoIncrementValue(*Context, uint64) error
	// Closer finalizes the set operation, persisting the result.
	Closer
}

// ReplaceableTable allows rows to be replaced through a Delete (if applicable) then Insert.
type ReplaceableTable interface {
	Table
	// Replacer returns a RowReplacer for this table. The RowReplacer will have Insert and optionally Delete called once
	// for each row, followed by a call to Close() when all rows have been processed.
	Replacer(ctx *Context) RowReplacer
}

// RowReplacer is a combination of RowDeleter and RowInserter.
type RowReplacer interface {
	EditOpenerCloser
	RowInserter
	RowDeleter
}

// UpdatableTable is a table that can process updates of existing rows via update statements.
type UpdatableTable interface {
	Table
	// Updater returns a RowUpdater for this table. The RowUpdater will have Update called once for each row to be
	// updated, followed by a call to Close() when all rows have been processed.
	Updater(ctx *Context) RowUpdater
}

// RowUpdater is an update cursor that can update one or more rows in a table.
type RowUpdater interface {
	EditOpenerCloser
	// Update the given row. Provides both the old and new rows.
	Update(ctx *Context, old Row, new Row) error
	// Closer finalizes the update operation, persisting the result.
	Closer
}

// TableEditor is the combination of interfaces that allow any table edit operation:
// i.e. INSERT, UPDATE, DELETE, REPLACE
type TableEditor interface {
	RowReplacer
	RowUpdater
}

// RewritableTable is an extension to Table that makes it simpler for integrators to adapt to schema changes that must
// rewrite every row of the table. In this case, rows are streamed from the existing table in the old schema,
// transformed / updated appropriately, and written with the new format.
type RewritableTable interface {
	Table
	AlterableTable

	// ShouldRewriteTable returns whether this table should be rewritten because of a schema change. The old and new
	// versions of the schema and modified column are provided. For some operations, one or both of |oldColumn| or
	// |newColumn| may be nil.
	// The engine may decide to rewrite tables regardless in some cases, such as when a new non-nullable column is added.
	ShouldRewriteTable(ctx *Context, oldSchema, newSchema PrimaryKeySchema, oldColumn, newColumn *Column) bool

	// RewriteInserter returns a RowInserter for the new schema. Rows from the current table, with the old schema, will
	// be streamed from the table and passed to this RowInserter. Implementor tables must still return rows in the
	// current schema until the rewrite operation completes. |Close| will be called on RowInserter when all rows have
	// been inserted.
	RewriteInserter(ctx *Context, oldSchema, newSchema PrimaryKeySchema, oldColumn, newColumn *Column, idxCols []IndexColumn) (RowInserter, error)
}

// AlterableTable should be implemented by tables that can receive
// ALTER TABLE statements to modify their schemas.
type AlterableTable interface {
	Table
	UpdatableTable

	// AddColumn adds a column to this table as given. If non-nil, order
	// specifies where in the schema to add the column.
	AddColumn(ctx *Context, column *Column, order *ColumnOrder) error
	// DropColumn drops the column with the name given.
	DropColumn(ctx *Context, columnName string) error
	// ModifyColumn modifies the column with the name given, replacing
	// with the new column definition provided (which may include a name
	// change). If non-nil, order specifies where in the schema to move
	// the column.
	ModifyColumn(ctx *Context, columnName string, column *Column, order *ColumnOrder) error
}

// SchemaValidator is a database that performs schema compatibility checks
// for CREATE and ALTER TABLE statements.
type SchemaValidator interface {
	Database
	// ValidateSchema lets storage integrators validate whether they can
	// serialize a given schema.
	ValidateSchema(Schema) error
}

// UnresolvedTable is a Table that is either unresolved or deferred for until an asOf resolution.
// Used by the analyzer during planning, and is not expected to be implemented by integrators.
type UnresolvedTable interface {
	Nameable
	// Database returns the database, which may be unresolved
	Database() Database
	// WithAsOf returns a copy of this versioned table with its AsOf
	// field set to the given value. Analogous to WithChildren.
	WithAsOf(asOf Expression) (Node, error)
	// AsOf returns this table's asof expression.
	AsOf() Expression
}

// TableNode is an interface for nodes that are also tables. A node that implements this interface exposes all the
// information needed for filters on the table to be optimized into indexes. This is possible when the return value
// of `UnderlyingTable` is a table that implements `sql.IndexAddressable`
// For an example of how to use this interface to optimize a system table or table function, see memory.IntSequenceTable
type TableNode interface {
	Table
	Node
	CollationCoercible
	Databaser
	// UnderlyingTable returns the table that this node is wrapping, recursively unwrapping any further layers of
	// wrapping to get to the base sql.Table.
	UnderlyingTable() Table
}

// MutableTableNode is a TableNode that can update its underlying table. Different methods are provided to accommodate
// different use cases that require working with base-level tables v. wrappers on top of them. Some uses of these
// methods might require that return values that implement all the same subinterfaces as the wrapped table, e.g.
// IndexedTable, ProjectableTable, etc. Callers of these methods should verify that the MutableTableNode's
// new table respects this contract.
type MutableTableNode interface {
	TableNode
	// WithTable returns a new TableNode with the table given. If the MutableTableNode has a MutableTableWrapper, it must
	// re-wrap the table given with this wrapper.
	WithTable(Table) (MutableTableNode, error)
	// ReplaceTable replaces the table with the table given, with no re-wrapping semantics.
	ReplaceTable(table Table) (MutableTableNode, error)
	// WrappedTable returns the Table this node wraps, without unwinding any additional layers of wrapped tables.
	WrappedTable() Table
}

// IndexSearchable lets a node use custom logic to create
// *plan.IndexedTableAccess
type IndexSearchable interface {
	// SkipIndexCosting defers to an integrator for provide a suitable
	// index lookup.
	SkipIndexCosting() bool
	// LookupForExpressions returns an sql.IndexLookup for an expression
	// set.
	LookupForExpressions(*Context, []Expression) (IndexLookup, error)
}

// IndexSearchableTable is a Table supports custom index generation.
type IndexSearchableTable interface {
	IndexAddressableTable
	IndexSearchable
}
