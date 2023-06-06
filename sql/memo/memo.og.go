// Code generated by optgen; DO NOT EDIT.

package memo

import (
	"fmt"
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression"
	"github.com/dolthub/go-mysql-server/sql/plan"
)

type CrossJoin struct {
	*JoinBase
}

var _ RelExpr = (*CrossJoin)(nil)
var _ JoinRel = (*CrossJoin)(nil)

func (r *CrossJoin) String() string {
	return FormatExpr(r)
}

func (r *CrossJoin) JoinPrivate() *JoinBase {
	return r.JoinBase
}

type InnerJoin struct {
	*JoinBase
}

var _ RelExpr = (*InnerJoin)(nil)
var _ JoinRel = (*InnerJoin)(nil)

func (r *InnerJoin) String() string {
	return FormatExpr(r)
}

func (r *InnerJoin) JoinPrivate() *JoinBase {
	return r.JoinBase
}

type LeftJoin struct {
	*JoinBase
}

var _ RelExpr = (*LeftJoin)(nil)
var _ JoinRel = (*LeftJoin)(nil)

func (r *LeftJoin) String() string {
	return FormatExpr(r)
}

func (r *LeftJoin) JoinPrivate() *JoinBase {
	return r.JoinBase
}

type SemiJoin struct {
	*JoinBase
}

var _ RelExpr = (*SemiJoin)(nil)
var _ JoinRel = (*SemiJoin)(nil)

func (r *SemiJoin) String() string {
	return FormatExpr(r)
}

func (r *SemiJoin) JoinPrivate() *JoinBase {
	return r.JoinBase
}

type AntiJoin struct {
	*JoinBase
}

var _ RelExpr = (*AntiJoin)(nil)
var _ JoinRel = (*AntiJoin)(nil)

func (r *AntiJoin) String() string {
	return FormatExpr(r)
}

func (r *AntiJoin) JoinPrivate() *JoinBase {
	return r.JoinBase
}

type LookupJoin struct {
	*JoinBase
	Lookup *Lookup
}

var _ RelExpr = (*LookupJoin)(nil)
var _ JoinRel = (*LookupJoin)(nil)

func (r *LookupJoin) String() string {
	return FormatExpr(r)
}

func (r *LookupJoin) JoinPrivate() *JoinBase {
	return r.JoinBase
}

type ConcatJoin struct {
	*JoinBase
	Concat []*Lookup
}

var _ RelExpr = (*ConcatJoin)(nil)
var _ JoinRel = (*ConcatJoin)(nil)

func (r *ConcatJoin) String() string {
	return FormatExpr(r)
}

func (r *ConcatJoin) JoinPrivate() *JoinBase {
	return r.JoinBase
}

type HashJoin struct {
	*JoinBase
	RightAttrs []*ExprGroup
	LeftAttrs  []*ExprGroup
}

var _ RelExpr = (*HashJoin)(nil)
var _ JoinRel = (*HashJoin)(nil)

func (r *HashJoin) String() string {
	return FormatExpr(r)
}

func (r *HashJoin) JoinPrivate() *JoinBase {
	return r.JoinBase
}

type MergeJoin struct {
	*JoinBase
	InnerScan *IndexScan
	OuterScan *IndexScan
	SwapCmp   bool
}

var _ RelExpr = (*MergeJoin)(nil)
var _ JoinRel = (*MergeJoin)(nil)

func (r *MergeJoin) String() string {
	return FormatExpr(r)
}

func (r *MergeJoin) JoinPrivate() *JoinBase {
	return r.JoinBase
}

type FullOuterJoin struct {
	*JoinBase
}

var _ RelExpr = (*FullOuterJoin)(nil)
var _ JoinRel = (*FullOuterJoin)(nil)

func (r *FullOuterJoin) String() string {
	return FormatExpr(r)
}

func (r *FullOuterJoin) JoinPrivate() *JoinBase {
	return r.JoinBase
}

type TableScan struct {
	*relBase
	Table *plan.ResolvedTable
}

var _ RelExpr = (*TableScan)(nil)
var _ SourceRel = (*TableScan)(nil)

func (r *TableScan) String() string {
	return FormatExpr(r)
}

func (r *TableScan) Name() string {
	return strings.ToLower(r.Table.Name())
}

func (r *TableScan) TableId() TableId {
	return TableIdForSource(r.g.Id)
}

func (r *TableScan) OutputCols() sql.Schema {
	return r.Table.Schema()
}

func (r *TableScan) Children() []*ExprGroup {
	return nil
}

type Values struct {
	*relBase
	Table *plan.ValueDerivedTable
}

var _ RelExpr = (*Values)(nil)
var _ SourceRel = (*Values)(nil)

func (r *Values) String() string {
	return FormatExpr(r)
}

func (r *Values) Name() string {
	return strings.ToLower(r.Table.Name())
}

func (r *Values) TableId() TableId {
	return TableIdForSource(r.g.Id)
}

func (r *Values) OutputCols() sql.Schema {
	return r.Table.Schema()
}

func (r *Values) Children() []*ExprGroup {
	return nil
}

type TableAlias struct {
	*relBase
	Table *plan.TableAlias
}

var _ RelExpr = (*TableAlias)(nil)
var _ SourceRel = (*TableAlias)(nil)

func (r *TableAlias) String() string {
	return FormatExpr(r)
}

func (r *TableAlias) Name() string {
	return strings.ToLower(r.Table.Name())
}

func (r *TableAlias) TableId() TableId {
	return TableIdForSource(r.g.Id)
}

func (r *TableAlias) OutputCols() sql.Schema {
	return r.Table.Schema()
}

func (r *TableAlias) Children() []*ExprGroup {
	return nil
}

type RecursiveTable struct {
	*relBase
	Table *plan.RecursiveTable
}

var _ RelExpr = (*RecursiveTable)(nil)
var _ SourceRel = (*RecursiveTable)(nil)

func (r *RecursiveTable) String() string {
	return FormatExpr(r)
}

func (r *RecursiveTable) Name() string {
	return strings.ToLower(r.Table.Name())
}

func (r *RecursiveTable) TableId() TableId {
	return TableIdForSource(r.g.Id)
}

func (r *RecursiveTable) OutputCols() sql.Schema {
	return r.Table.Schema()
}

func (r *RecursiveTable) Children() []*ExprGroup {
	return nil
}

type RecursiveCte struct {
	*relBase
	Table *plan.RecursiveCte
}

var _ RelExpr = (*RecursiveCte)(nil)
var _ SourceRel = (*RecursiveCte)(nil)

func (r *RecursiveCte) String() string {
	return FormatExpr(r)
}

func (r *RecursiveCte) Name() string {
	return strings.ToLower(r.Table.Name())
}

func (r *RecursiveCte) TableId() TableId {
	return TableIdForSource(r.g.Id)
}

func (r *RecursiveCte) OutputCols() sql.Schema {
	return r.Table.Schema()
}

func (r *RecursiveCte) Children() []*ExprGroup {
	return nil
}

type SubqueryAlias struct {
	*relBase
	Table *plan.SubqueryAlias
}

var _ RelExpr = (*SubqueryAlias)(nil)
var _ SourceRel = (*SubqueryAlias)(nil)

func (r *SubqueryAlias) String() string {
	return FormatExpr(r)
}

func (r *SubqueryAlias) Name() string {
	return strings.ToLower(r.Table.Name())
}

func (r *SubqueryAlias) TableId() TableId {
	return TableIdForSource(r.g.Id)
}

func (r *SubqueryAlias) OutputCols() sql.Schema {
	return r.Table.Schema()
}

func (r *SubqueryAlias) Children() []*ExprGroup {
	return nil
}

type Max1Row struct {
	*relBase
	Table sql.NameableNode
}

var _ RelExpr = (*Max1Row)(nil)
var _ SourceRel = (*Max1Row)(nil)

func (r *Max1Row) String() string {
	return FormatExpr(r)
}

func (r *Max1Row) Name() string {
	return strings.ToLower(r.Table.Name())
}

func (r *Max1Row) TableId() TableId {
	return TableIdForSource(r.g.Id)
}

func (r *Max1Row) OutputCols() sql.Schema {
	return r.Table.Schema()
}

func (r *Max1Row) Children() []*ExprGroup {
	return nil
}

type TableFunc struct {
	*relBase
	Table sql.TableFunction
}

var _ RelExpr = (*TableFunc)(nil)
var _ SourceRel = (*TableFunc)(nil)

func (r *TableFunc) String() string {
	return FormatExpr(r)
}

func (r *TableFunc) Name() string {
	return strings.ToLower(r.Table.Name())
}

func (r *TableFunc) TableId() TableId {
	return TableIdForSource(r.g.Id)
}

func (r *TableFunc) OutputCols() sql.Schema {
	return r.Table.Schema()
}

func (r *TableFunc) Children() []*ExprGroup {
	return nil
}

type EmptyTable struct {
	*relBase
	Table *plan.EmptyTable
}

var _ RelExpr = (*EmptyTable)(nil)
var _ SourceRel = (*EmptyTable)(nil)

func (r *EmptyTable) String() string {
	return FormatExpr(r)
}

func (r *EmptyTable) Name() string {
	return strings.ToLower(r.Table.Name())
}

func (r *EmptyTable) TableId() TableId {
	return TableIdForSource(r.g.Id)
}

func (r *EmptyTable) OutputCols() sql.Schema {
	return r.Table.Schema()
}

func (r *EmptyTable) Children() []*ExprGroup {
	return nil
}

type Project struct {
	*relBase
	Child       *ExprGroup
	Projections []*ExprGroup
}

var _ RelExpr = (*Project)(nil)

func (r *Project) String() string {
	return FormatExpr(r)
}

func (r *Project) Children() []*ExprGroup {
	return []*ExprGroup{r.Child}
}

func (r *Project) outputCols() sql.Schema {
	var s = make(sql.Schema, len(r.Projections))
	for i, e := range r.Projections {
		s[i] = ScalarToSqlCol(e)
	}
	return s
}

type Distinct struct {
	*relBase
	Child *ExprGroup
}

var _ RelExpr = (*Distinct)(nil)

func (r *Distinct) String() string {
	return FormatExpr(r)
}

func (r *Distinct) Children() []*ExprGroup {
	return []*ExprGroup{r.Child}
}

func (r *Distinct) outputCols() sql.Schema {
	return r.Child.RelProps.OutputCols()
}

type Filter struct {
	*relBase
	Child   *ExprGroup
	Filters []*ExprGroup
}

var _ RelExpr = (*Filter)(nil)

func (r *Filter) String() string {
	return FormatExpr(r)
}

func (r *Filter) Children() []*ExprGroup {
	return []*ExprGroup{r.Child}
}

func (r *Filter) outputCols() sql.Schema {
	return r.Child.RelProps.OutputCols()
}

type Equal struct {
	*scalarBase
	Left  *ExprGroup
	Right *ExprGroup
}

var _ ScalarExpr = (*Equal)(nil)

func (r *Equal) ExprId() ScalarExprId {
	return ScalarExprEqual
}

func (r *Equal) String() string {
	return FormatExpr(r)
}

func (r *Equal) Children() []*ExprGroup {
	return []*ExprGroup{r.Left, r.Right}
}

type Literal struct {
	*scalarBase
	Val interface{}
	Typ sql.Type
}

var _ ScalarExpr = (*Literal)(nil)

func (r *Literal) ExprId() ScalarExprId {
	return ScalarExprLiteral
}

func (r *Literal) String() string {
	return FormatExpr(r)
}

func (r *Literal) Children() []*ExprGroup {
	return nil
}

type ColRef struct {
	*scalarBase
	Col   sql.ColumnId
	Table GroupId
	Gf    *expression.GetField
}

var _ ScalarExpr = (*ColRef)(nil)

func (r *ColRef) ExprId() ScalarExprId {
	return ScalarExprColRef
}

func (r *ColRef) String() string {
	return FormatExpr(r)
}

func (r *ColRef) Children() []*ExprGroup {
	return nil
}

type Not struct {
	*scalarBase
	Child *ExprGroup
}

var _ ScalarExpr = (*Not)(nil)

func (r *Not) ExprId() ScalarExprId {
	return ScalarExprNot
}

func (r *Not) String() string {
	return FormatExpr(r)
}

func (r *Not) Children() []*ExprGroup {
	return []*ExprGroup{r.Child}
}

func (r *Not) outputCols() sql.Schema {
	return r.Child.RelProps.OutputCols()
}

type Or struct {
	*scalarBase
	Left  *ExprGroup
	Right *ExprGroup
}

var _ ScalarExpr = (*Or)(nil)

func (r *Or) ExprId() ScalarExprId {
	return ScalarExprOr
}

func (r *Or) String() string {
	return FormatExpr(r)
}

func (r *Or) Children() []*ExprGroup {
	return []*ExprGroup{r.Left, r.Right}
}

type And struct {
	*scalarBase
	Left  *ExprGroup
	Right *ExprGroup
}

var _ ScalarExpr = (*And)(nil)

func (r *And) ExprId() ScalarExprId {
	return ScalarExprAnd
}

func (r *And) String() string {
	return FormatExpr(r)
}

func (r *And) Children() []*ExprGroup {
	return []*ExprGroup{r.Left, r.Right}
}

type InTuple struct {
	*scalarBase
	Left  *ExprGroup
	Right *ExprGroup
}

var _ ScalarExpr = (*InTuple)(nil)

func (r *InTuple) ExprId() ScalarExprId {
	return ScalarExprInTuple
}

func (r *InTuple) String() string {
	return FormatExpr(r)
}

func (r *InTuple) Children() []*ExprGroup {
	return []*ExprGroup{r.Left, r.Right}
}

type Lt struct {
	*scalarBase
	Left  *ExprGroup
	Right *ExprGroup
}

var _ ScalarExpr = (*Lt)(nil)

func (r *Lt) ExprId() ScalarExprId {
	return ScalarExprLt
}

func (r *Lt) String() string {
	return FormatExpr(r)
}

func (r *Lt) Children() []*ExprGroup {
	return []*ExprGroup{r.Left, r.Right}
}

type Leq struct {
	*scalarBase
	Left  *ExprGroup
	Right *ExprGroup
}

var _ ScalarExpr = (*Leq)(nil)

func (r *Leq) ExprId() ScalarExprId {
	return ScalarExprLeq
}

func (r *Leq) String() string {
	return FormatExpr(r)
}

func (r *Leq) Children() []*ExprGroup {
	return []*ExprGroup{r.Left, r.Right}
}

type Gt struct {
	*scalarBase
	Left  *ExprGroup
	Right *ExprGroup
}

var _ ScalarExpr = (*Gt)(nil)

func (r *Gt) ExprId() ScalarExprId {
	return ScalarExprGt
}

func (r *Gt) String() string {
	return FormatExpr(r)
}

func (r *Gt) Children() []*ExprGroup {
	return []*ExprGroup{r.Left, r.Right}
}

type Geq struct {
	*scalarBase
	Left  *ExprGroup
	Right *ExprGroup
}

var _ ScalarExpr = (*Geq)(nil)

func (r *Geq) ExprId() ScalarExprId {
	return ScalarExprGeq
}

func (r *Geq) String() string {
	return FormatExpr(r)
}

func (r *Geq) Children() []*ExprGroup {
	return []*ExprGroup{r.Left, r.Right}
}

type NullSafeEq struct {
	*scalarBase
	Left  *ExprGroup
	Right *ExprGroup
}

var _ ScalarExpr = (*NullSafeEq)(nil)

func (r *NullSafeEq) ExprId() ScalarExprId {
	return ScalarExprNullSafeEq
}

func (r *NullSafeEq) String() string {
	return FormatExpr(r)
}

func (r *NullSafeEq) Children() []*ExprGroup {
	return []*ExprGroup{r.Left, r.Right}
}

type Regexp struct {
	*scalarBase
	Left  *ExprGroup
	Right *ExprGroup
}

var _ ScalarExpr = (*Regexp)(nil)

func (r *Regexp) ExprId() ScalarExprId {
	return ScalarExprRegexp
}

func (r *Regexp) String() string {
	return FormatExpr(r)
}

func (r *Regexp) Children() []*ExprGroup {
	return []*ExprGroup{r.Left, r.Right}
}

type Arithmetic struct {
	*scalarBase
	Left  *ExprGroup
	Right *ExprGroup
	Op    ArithType
}

var _ ScalarExpr = (*Arithmetic)(nil)

func (r *Arithmetic) ExprId() ScalarExprId {
	return ScalarExprArithmetic
}

func (r *Arithmetic) String() string {
	return FormatExpr(r)
}

func (r *Arithmetic) Children() []*ExprGroup {
	return []*ExprGroup{r.Left, r.Right}
}

type Bindvar struct {
	*scalarBase
	Name string
	Typ  sql.Type
}

var _ ScalarExpr = (*Bindvar)(nil)

func (r *Bindvar) ExprId() ScalarExprId {
	return ScalarExprBindvar
}

func (r *Bindvar) String() string {
	return FormatExpr(r)
}

func (r *Bindvar) Children() []*ExprGroup {
	return nil
}

type IsNull struct {
	*scalarBase
	Child *ExprGroup
}

var _ ScalarExpr = (*IsNull)(nil)

func (r *IsNull) ExprId() ScalarExprId {
	return ScalarExprIsNull
}

func (r *IsNull) String() string {
	return FormatExpr(r)
}

func (r *IsNull) Children() []*ExprGroup {
	return []*ExprGroup{r.Child}
}

func (r *IsNull) outputCols() sql.Schema {
	return r.Child.RelProps.OutputCols()
}

type Tuple struct {
	*scalarBase
	Values []*ExprGroup
}

var _ ScalarExpr = (*Tuple)(nil)

func (r *Tuple) ExprId() ScalarExprId {
	return ScalarExprTuple
}

func (r *Tuple) String() string {
	return FormatExpr(r)
}

func (r *Tuple) Children() []*ExprGroup {
	return nil
}

type Hidden struct {
	*scalarBase
	E      sql.Expression
	Cols   sql.ColSet
	Tables sql.FastIntSet
}

var _ ScalarExpr = (*Hidden)(nil)

func (r *Hidden) ExprId() ScalarExprId {
	return ScalarExprHidden
}

func (r *Hidden) String() string {
	return FormatExpr(r)
}

func (r *Hidden) Children() []*ExprGroup {
	return nil
}

func FormatExpr(r exprType) string {
	switch r := r.(type) {
	case *CrossJoin:
		return fmt.Sprintf("crossjoin %d %d", r.Left.Id, r.Right.Id)
	case *InnerJoin:
		return fmt.Sprintf("innerjoin %d %d", r.Left.Id, r.Right.Id)
	case *LeftJoin:
		return fmt.Sprintf("leftjoin %d %d", r.Left.Id, r.Right.Id)
	case *SemiJoin:
		return fmt.Sprintf("semijoin %d %d", r.Left.Id, r.Right.Id)
	case *AntiJoin:
		return fmt.Sprintf("antijoin %d %d", r.Left.Id, r.Right.Id)
	case *LookupJoin:
		return fmt.Sprintf("lookupjoin %d %d", r.Left.Id, r.Right.Id)
	case *ConcatJoin:
		return fmt.Sprintf("concatjoin %d %d", r.Left.Id, r.Right.Id)
	case *HashJoin:
		return fmt.Sprintf("hashjoin %d %d", r.Left.Id, r.Right.Id)
	case *MergeJoin:
		return fmt.Sprintf("mergejoin %d %d", r.Left.Id, r.Right.Id)
	case *FullOuterJoin:
		return fmt.Sprintf("fullouterjoin %d %d", r.Left.Id, r.Right.Id)
	case *TableScan:
		return fmt.Sprintf("tablescan: %s", r.Name())
	case *Values:
		return fmt.Sprintf("values: %s", r.Name())
	case *TableAlias:
		return fmt.Sprintf("tablealias: %s", r.Name())
	case *RecursiveTable:
		return fmt.Sprintf("recursivetable: %s", r.Name())
	case *RecursiveCte:
		return fmt.Sprintf("recursivecte: %s", r.Name())
	case *SubqueryAlias:
		return fmt.Sprintf("subqueryalias: %s", r.Name())
	case *Max1Row:
		return fmt.Sprintf("max1row: %s", r.Name())
	case *TableFunc:
		return fmt.Sprintf("tablefunc: %s", r.Name())
	case *EmptyTable:
		return fmt.Sprintf("emptytable: %s", r.Name())
	case *Project:
		return fmt.Sprintf("project: %d", r.Child.Id)
	case *Distinct:
		return fmt.Sprintf("distinct: %d", r.Child.Id)
	case *Filter:
		return fmt.Sprintf("filter: %d", r.Child.Id)
	case *Equal:
		return fmt.Sprintf("equal %d %d", r.Left.Id, r.Right.Id)
	case *Literal:
		return fmt.Sprintf("literal: %v %s", r.Val, r.Typ)
	case *ColRef:
		return fmt.Sprintf("colref: '%s.%s'", r.Gf.Table(), r.Gf.Name())
	case *Not:
		return fmt.Sprintf("not: %d", r.Child.Id)
	case *Or:
		return fmt.Sprintf("or %d %d", r.Left.Id, r.Right.Id)
	case *And:
		return fmt.Sprintf("and %d %d", r.Left.Id, r.Right.Id)
	case *InTuple:
		return fmt.Sprintf("intuple %d %d", r.Left.Id, r.Right.Id)
	case *Lt:
		return fmt.Sprintf("lt %d %d", r.Left.Id, r.Right.Id)
	case *Leq:
		return fmt.Sprintf("leq %d %d", r.Left.Id, r.Right.Id)
	case *Gt:
		return fmt.Sprintf("gt %d %d", r.Left.Id, r.Right.Id)
	case *Geq:
		return fmt.Sprintf("geq %d %d", r.Left.Id, r.Right.Id)
	case *NullSafeEq:
		return fmt.Sprintf("nullsafeeq %d %d", r.Left.Id, r.Right.Id)
	case *Regexp:
		return fmt.Sprintf("regexp %d %d", r.Left.Id, r.Right.Id)
	case *Arithmetic:
		return fmt.Sprintf("arithmetic %d %d", r.Left.Id, r.Right.Id)
	case *Bindvar:
		return fmt.Sprintf("bindvar: %s", r.Name)
	case *IsNull:
		return fmt.Sprintf("isnull: %d", r.Child.Id)
	case *Tuple:
		vals := make([]string, len(r.Values))
		for i, v := range r.Values {
			vals[i] = fmt.Sprintf("%d", v.Id)
		}
		return fmt.Sprintf("tuple: %s", strings.Join(vals, " "))
	case *Hidden:
		return fmt.Sprintf("hidden: %s", r.E)
	default:
		panic(fmt.Sprintf("unknown RelExpr type: %T", r))
	}
}

func buildRelExpr(b *ExecBuilder, r RelExpr, input sql.Schema, children ...sql.Node) (sql.Node, error) {
	var result sql.Node
	var err error

	switch r := r.(type) {
	case *CrossJoin:
		result, err = b.buildCrossJoin(r, input, children...)
	case *InnerJoin:
		result, err = b.buildInnerJoin(r, input, children...)
	case *LeftJoin:
		result, err = b.buildLeftJoin(r, input, children...)
	case *SemiJoin:
		result, err = b.buildSemiJoin(r, input, children...)
	case *AntiJoin:
		result, err = b.buildAntiJoin(r, input, children...)
	case *LookupJoin:
		result, err = b.buildLookupJoin(r, input, children...)
	case *ConcatJoin:
		result, err = b.buildConcatJoin(r, input, children...)
	case *HashJoin:
		result, err = b.buildHashJoin(r, input, children...)
	case *MergeJoin:
		result, err = b.buildMergeJoin(r, input, children...)
	case *FullOuterJoin:
		result, err = b.buildFullOuterJoin(r, input, children...)
	case *TableScan:
		result, err = b.buildTableScan(r, input, children...)
	case *Values:
		result, err = b.buildValues(r, input, children...)
	case *TableAlias:
		result, err = b.buildTableAlias(r, input, children...)
	case *RecursiveTable:
		result, err = b.buildRecursiveTable(r, input, children...)
	case *RecursiveCte:
		result, err = b.buildRecursiveCte(r, input, children...)
	case *SubqueryAlias:
		result, err = b.buildSubqueryAlias(r, input, children...)
	case *Max1Row:
		result, err = b.buildMax1Row(r, input, children...)
	case *TableFunc:
		result, err = b.buildTableFunc(r, input, children...)
	case *EmptyTable:
		result, err = b.buildEmptyTable(r, input, children...)
	case *Project:
		result, err = b.buildProject(r, input, children...)
	case *Filter:
		result, err = b.buildFilter(r, input, children...)
	default:
		panic(fmt.Sprintf("unknown RelExpr type: %T", r))
	}

	if err != nil {
		return nil, err
	}

	result, err = r.Group().finalize(result, input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func buildScalarExpr(b *ExecBuilder, r ScalarExpr, sch sql.Schema) (sql.Expression, error) {
	switch r := r.(type) {
	case *Equal:
		return b.buildEqual(r, sch)
	case *Literal:
		return b.buildLiteral(r, sch)
	case *ColRef:
		return b.buildColRef(r, sch)
	case *Not:
		return b.buildNot(r, sch)
	case *Or:
		return b.buildOr(r, sch)
	case *And:
		return b.buildAnd(r, sch)
	case *InTuple:
		return b.buildInTuple(r, sch)
	case *Lt:
		return b.buildLt(r, sch)
	case *Leq:
		return b.buildLeq(r, sch)
	case *Gt:
		return b.buildGt(r, sch)
	case *Geq:
		return b.buildGeq(r, sch)
	case *NullSafeEq:
		return b.buildNullSafeEq(r, sch)
	case *Regexp:
		return b.buildRegexp(r, sch)
	case *Arithmetic:
		return b.buildArithmetic(r, sch)
	case *Bindvar:
		return b.buildBindvar(r, sch)
	case *IsNull:
		return b.buildIsNull(r, sch)
	case *Tuple:
		return b.buildTuple(r, sch)
	case *Hidden:
		return b.buildHidden(r, sch)
	default:
		panic(fmt.Sprintf("unknown ScalarExpr type: %T", r))
	}
}
