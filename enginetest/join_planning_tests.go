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

package enginetest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/enginetest/scriptgen/setup"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/plan"
	"github.com/dolthub/go-mysql-server/sql/transform"
)

type JoinPlanTest struct {
	q       string
	types   []plan.JoinType
	indexes []string
	exp     []sql.Row
	order   []string
	skipOld bool
}

var JoinPlanningTests = []struct {
	name  string
	setup []string
	tests []JoinPlanTest
}{
	{
		name: "merge join unary index",
		setup: []string{
			"CREATE table xy (x int primary key, y int, index y_idx(y));",
			"create table rs (r int primary key, s int, index s_idx(s));",
			"CREATE table uv (u int primary key, v int);",
			"CREATE table ab (a int primary key, b int);",
			"insert into xy values (1,0), (2,1), (0,2), (3,3);",
			"insert into rs values (0,0), (1,0), (2,0), (4,4), (5,4);",
			"insert into uv values (0,1), (1,1), (2,2), (3,2);",
			"insert into ab values (0,2), (1,2), (2,2), (3,1);",
			"update information_schema.statistics set cardinality = 1000 where table_name in ('ab', 'rs', 'xy', 'uv');",
		},
		tests: []JoinPlanTest{
			{
				q:     "select u,a,y from uv join (select /*+ JOIN_ORDER(ab, xy) */ * from ab join xy on y = a) r on u = r.a order by 1",
				types: []plan.JoinType{plan.JoinTypeLookup, plan.JoinTypeMerge},
				exp:   []sql.Row{{0, 0, 0}, {1, 1, 1}, {2, 2, 2}, {3, 3, 3}},
			},
			{
				q:     "select /*+ JOIN_ORDER(ab, xy) */ * from ab join xy on y = a order by 1, 3",
				types: []plan.JoinType{plan.JoinTypeMerge},
				exp:   []sql.Row{{0, 2, 1, 0}, {1, 2, 2, 1}, {2, 2, 0, 2}, {3, 1, 3, 3}},
			},
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs left join xy on y = s order by 1, 3",
				types: []plan.JoinType{plan.JoinTypeLeftOuterMerge},
				exp:   []sql.Row{{0, 0, 1, 0}, {1, 0, 1, 0}, {2, 0, 1, 0}, {4, 4, nil, nil}, {5, 4, nil, nil}},
			},
			{
				// extra join condition does not filter left-only rows
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs left join xy on y = s and y+s = 0 order by 1, 3",
				types: []plan.JoinType{plan.JoinTypeLeftOuterMerge},
				exp:   []sql.Row{{0, 0, 1, 0}, {1, 0, 1, 0}, {2, 0, 1, 0}, {4, 4, nil, nil}, {5, 4, nil, nil}},
			},
			{
				// extra join condition does not filter left-only rows
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs left join xy on y+2 = s and s-y = 2 order by 1, 3",
				types: []plan.JoinType{plan.JoinTypeLeftOuterMerge},
				exp:   []sql.Row{{0, 0, nil, nil}, {1, 0, nil, nil}, {2, 0, nil, nil}, {4, 4, 0, 2}, {5, 4, 0, 2}},
			},
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on y = r order by 1, 3",
				types: []plan.JoinType{plan.JoinTypeMerge},
				exp:   []sql.Row{{0, 0, 1, 0}, {1, 0, 2, 1}, {2, 0, 0, 2}},
			},
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on r = y order by 1, 3",
				types: []plan.JoinType{plan.JoinTypeMerge},
				exp:   []sql.Row{{0, 0, 1, 0}, {1, 0, 2, 1}, {2, 0, 0, 2}},
			},
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on y = s order by 1, 3",
				types: []plan.JoinType{plan.JoinTypeMerge},
				exp:   []sql.Row{{0, 0, 1, 0}, {1, 0, 1, 0}, {2, 0, 1, 0}},
			},
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on y = s and y = r order by 1, 3",
				types: []plan.JoinType{plan.JoinTypeMerge},
				exp:   []sql.Row{{0, 0, 1, 0}},
			},
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on y+2 = s order by 1, 3",
				types: []plan.JoinType{plan.JoinTypeMerge},
				exp:   []sql.Row{{4, 4, 0, 2}, {5, 4, 0, 2}},
			},
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on y = s-1 order by 1, 3",
				types: []plan.JoinType{plan.JoinTypeHash},
				exp:   []sql.Row{{4, 4, 3, 3}, {5, 4, 3, 3}},
			},
			//{
			// TODO: cannot hash join on compound expressions
			//	q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on y = mod(s,2) order by 1, 3",
			//	types: []plan.JoinType{plan.JoinTypeInner},
			//	exp:   []sql.Row{{0,0,1,0},{0, 0, 1, 0},{2,0,1,0},{4,4,1,0}},
			//},
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on 2 = s+y order by 1, 3",
				types: []plan.JoinType{plan.JoinTypeInner},
				exp:   []sql.Row{{0, 0, 0, 2}, {1, 0, 0, 2}, {2, 0, 0, 2}},
			},
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on y > s+2 order by 1, 3",
				types: []plan.JoinType{plan.JoinTypeInner},
				exp:   []sql.Row{{0, 0, 3, 3}, {1, 0, 3, 3}, {2, 0, 3, 3}},
			},
		},
	},
	{
		name: "merge join multi match",
		setup: []string{
			"CREATE table xy (x int primary key, y int, index y_idx(y));",
			"create table rs (r int primary key, s int, index s_idx(s));",
			"insert into xy values (1,0), (2,1), (0,8), (3,7), (5,4), (4,0);",
			"insert into rs values (0,0),(2,3),(3,0), (4,8), (5,4);",
			"update information_schema.statistics set cardinality = 1000 where table_name in ('rs', 'xy');",
		},
		tests: []JoinPlanTest{
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on y = s order by 1,3",
				types: []plan.JoinType{plan.JoinTypeMerge},
				exp:   []sql.Row{{0, 0, 1, 0}, {0, 0, 4, 0}, {3, 0, 1, 0}, {3, 0, 4, 0}, {4, 8, 0, 8}, {5, 4, 5, 4}},
			},
		},
	},
	{
		name: "merge join zero rows",
		setup: []string{
			"CREATE table xy (x int primary key, y int, index y_idx(y));",
			"create table rs (r int primary key, s int, index s_idx(s));",
			"insert into xy values (1,0);",
			"update information_schema.statistics set cardinality = 10 where table_name = 'xy';",
			"update information_schema.statistics set cardinality = 1000000000 where table_name = 'rs';",
		},
		tests: []JoinPlanTest{
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on y = s order by 1,3",
				types: []plan.JoinType{plan.JoinTypeMerge},
				exp:   []sql.Row{},
			},
		},
	},
	{
		name: "merge join multi arity",
		setup: []string{
			"CREATE table xy (x int primary key, y int, index yx_idx(y,x));",
			"create table rs (r int primary key, s int, index s_idx(s));",
			"insert into xy values (1,0), (2,1), (0,8), (3,7), (5,4), (4,0);",
			"insert into rs values (0,0),(2,3),(3,0), (4,8), (5,4);",
			"update information_schema.statistics set cardinality = 1000 where table_name in ('xy', 'rs');",
		},
		tests: []JoinPlanTest{
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on y = s order by 1,3",
				types: []plan.JoinType{plan.JoinTypeMerge},
				exp:   []sql.Row{{0, 0, 1, 0}, {0, 0, 4, 0}, {3, 0, 1, 0}, {3, 0, 4, 0}, {4, 8, 0, 8}, {5, 4, 5, 4}},
			},
		},
	},
	{
		name: "merge join keyless index",
		setup: []string{
			"CREATE table xy (x int, y int, index yx_idx(y,x));",
			"create table rs (r int, s int, index s_idx(s));",
			"insert into xy values (1,0), (2,1), (0,8), (3,7), (5,4), (4,0);",
			"insert into rs values (0,0),(2,3),(3,0), (4,8), (5,4);",
			"update information_schema.statistics set cardinality = 1000 where table_name in ('xy', 'rs');",
		},
		tests: []JoinPlanTest{
			{
				q:     "select /*+ JOIN_ORDER(rs, xy) */ * from rs join xy on y = s order by 1,3",
				types: []plan.JoinType{plan.JoinTypeMerge},
				exp:   []sql.Row{{0, 0, 1, 0}, {0, 0, 4, 0}, {3, 0, 1, 0}, {3, 0, 4, 0}, {4, 8, 0, 8}, {5, 4, 5, 4}},
			},
		},
	},
	{
		name: "partial [lookup] join tests",
		setup: []string{
			"CREATE table xy (x int primary key, y int);",
			"create table rs (r int primary key, s int);",
			"CREATE table uv (u int primary key, v int);",
			"CREATE table ab (a int primary key, b int);",
			"insert into xy values (1,0), (2,1), (0,2), (3,3);",
			"insert into rs values (0,0), (1,0), (2,0), (4,4);",
			"insert into uv values (0,1), (1,1), (2,2), (3,2);",
			"insert into ab values (0,2), (1,2), (2,2), (3,1);",
			"update information_schema.statistics set cardinality = 100 where table_name in ('xy', 'rs', 'uv', 'ab');",
		},
		tests: []JoinPlanTest{
			{
				q:     "select /*+ LOOKUP_JOIN(xy,scalarSubq0) */ * from xy where x in (select b from ab where a in (0,1,2));",
				types: []plan.JoinType{plan.JoinTypeLookup},
				exp:   []sql.Row{{2, 1}},
			},
			{
				// TODO: RIGHT_SEMI_JOIN tuple equalities
				q:     "select /*+ LOOKUP_JOIN(xy,scalarSubq0) */ * from xy where (x,y) in (select b,a from ab where a in (0,1,2));",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{2, 1}},
			},
			{
				q:     "select /*+ LOOKUP_JOIN(xy,scalarSubq0) */ * from xy where x in (select a from ab);",
				types: []plan.JoinType{plan.JoinTypeLookup},
				exp:   []sql.Row{{2, 1}, {1, 0}, {0, 2}, {3, 3}},
			},
			{
				q:     "select /*+ LOOKUP_JOIN(xy,scalarSubq0) */ * from xy where x in (select a from ab where a in (1,2));",
				types: []plan.JoinType{plan.JoinTypeLookup},
				exp:   []sql.Row{{2, 1}, {1, 0}},
			},
			{
				q:     "select /*+ LOOKUP_JOIN(xy,scalarSubq0)  */* from xy where x in (select a from ab);",
				types: []plan.JoinType{plan.JoinTypeLookup},
				exp:   []sql.Row{{2, 1}, {1, 0}, {0, 2}, {3, 3}},
			},
			{
				q:     "select /*+ LOOKUP_JOIN(xy,ab) MERGE_JOIN(ab,uv) JOIN_ORDER(ab,uv,xy) */ * from xy where EXISTS (select 1 from ab join uv on a = u where x = a);",
				types: []plan.JoinType{plan.JoinTypeLookup, plan.JoinTypeMerge},
				exp:   []sql.Row{{2, 1}, {1, 0}, {0, 2}, {3, 3}},
			},
			{
				q:     "select * from xy where y+1 not in (select u from uv);",
				types: []plan.JoinType{plan.JoinTypeLeftOuterHashExcludeNulls},
				exp:   []sql.Row{{3, 3}},
			},
			{
				q:     "select * from xy where x not in (select u from uv where u not in (select a from ab where a not in (select r from rs where r = 1))) order by 1;",
				types: []plan.JoinType{plan.JoinTypeLeftOuterHashExcludeNulls, plan.JoinTypeLeftOuterHashExcludeNulls, plan.JoinTypeAntiLookup},
				exp:   []sql.Row{{0, 2}, {2, 1}, {3, 3}},
			},
			{
				q:     "select * from xy where x != (select r from rs where r = 1) order by 1;",
				types: []plan.JoinType{plan.JoinTypeLeftOuterHashExcludeNulls},
				exp:   []sql.Row{{0, 2}, {2, 1}, {3, 3}},
			},
			{
				// anti join will be cross-join-right, be passed non-nil parent row
				q:     "select x,a from ab, (select * from xy where x != (select r from rs where r = 1) order by 1) sq where x = 2 and b = 2 order by 1,2;",
				types: []plan.JoinType{plan.JoinTypeCrossHash, plan.JoinTypeLeftOuterHashExcludeNulls},
				exp:   []sql.Row{{2, 0}, {2, 1}, {2, 2}},
			},
			{
				// scope and parent row are non-nil
				q: `
select * from uv where u > (
  select x from ab, (
    select x from xy where x != (
      select r from rs where r = 1
    ) order by 1
  ) sq
  order by 1 limit 1
)
order by 1;`,
				types: []plan.JoinType{plan.JoinTypeSemi, plan.JoinTypeCrossHash, plan.JoinTypeLeftOuterHashExcludeNulls},
				exp:   []sql.Row{{1, 1}, {2, 2}, {3, 2}},
			},
			{
				// cast prevents scope merging
				q:     "select * from xy where x != (select cast(r as signed) from rs where r = 1) order by 1;",
				types: []plan.JoinType{plan.JoinTypeLeftOuterHashExcludeNulls},
				exp:   []sql.Row{{0, 2}, {2, 1}, {3, 3}},
			},
			{
				// order by will be discarded
				q:     "select * from xy where x != (select r from rs where r = 1 order by 1) order by 1;",
				types: []plan.JoinType{plan.JoinTypeLeftOuterHashExcludeNulls},
				exp:   []sql.Row{{0, 2}, {2, 1}, {3, 3}},
			},
			{
				// limit prevents scope merging
				q:     "select * from xy where x != (select r from rs where r = 1 limit 1) order by 1;",
				types: []plan.JoinType{plan.JoinTypeLeftOuterHashExcludeNulls},
				exp:   []sql.Row{{0, 2}, {2, 1}, {3, 3}},
			},
			{
				q:     "select * from xy where y-1 in (select u from uv) order by 1;",
				types: []plan.JoinType{plan.JoinTypeHash},
				exp:   []sql.Row{{0, 2}, {2, 1}, {3, 3}},
			},
			{
				// semi join will be right-side, be passed non-nil parent row
				q:     "select x,a from ab, (select * from xy where x = (select r from rs where r = 1) order by 1) sq order by 1,2",
				types: []plan.JoinType{plan.JoinTypeCrossHash, plan.JoinTypeLookup},
				exp:   []sql.Row{{1, 0}, {1, 1}, {1, 2}, {1, 3}},
			},
			//{
			// scope and parent row are non-nil
			// TODO: subquery alias unable to track parent row from a different scope
			//				q: `
			//select * from uv where u > (
			//  select x from ab, (
			//    select x from xy where x = (
			//      select r from rs where r = 1
			//    ) order by 1
			//  ) sq
			//  order by 1 limit 1
			//)
			//order by 1;`,
			//types: []plan.JoinType{plan.JoinTypeCrossHash, plan.JoinTypeLookup},
			//exp:   []sql.Row{{2, 2}, {3, 2}},
			//},
			{
				q:     "select * from xy where y-1 in (select cast(u as signed) from uv) order by 1;",
				types: []plan.JoinType{plan.JoinTypeHash},
				exp:   []sql.Row{{0, 2}, {2, 1}, {3, 3}},
			},
			{
				q:     "select * from xy where y-1 in (select u from uv order by 1) order by 1;",
				types: []plan.JoinType{plan.JoinTypeHash},
				exp:   []sql.Row{{0, 2}, {2, 1}, {3, 3}},
			},
			{
				q:     "select * from xy where y-1 in (select u from uv order by 1 limit 1) order by 1;",
				types: []plan.JoinType{plan.JoinTypeHash},
				exp:   []sql.Row{{2, 1}},
			},
			{
				q:     "select * from xy where x in (select u from uv join ab on u = a and a = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeHash, plan.JoinTypeLookup},
				exp:   []sql.Row{{2, 1}},
			},
			{
				q:     "select * from xy where x = (select u from uv join ab on u = a and a = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeLookup, plan.JoinTypeLookup},
				exp:   []sql.Row{{2, 1}},
			},
			{
				// group by doesn't transform
				q:     "select * from xy where y-1 in (select u from uv group by v having v = 2 order by 1) order by 1;",
				types: []plan.JoinType{plan.JoinTypeHash},
				exp:   []sql.Row{{3, 3}},
			},
			{
				// window doesn't transform
				q:     "select * from xy where y-1 in (select row_number() over (order by v) from uv) order by 1;",
				types: []plan.JoinType{plan.JoinTypeHash},
				exp:   []sql.Row{{0, 2}, {3, 3}},
			},
		},
	},
	{
		name: "empty join tests",
		setup: []string{
			"CREATE table xy (x int primary key, y int);",
			"CREATE table uv (u int primary key, v int);",
			"insert into xy values (1,0), (2,1), (0,2), (3,3);",
			"insert into uv values (0,1), (1,1), (2,2), (3,2);",
		},
		tests: []JoinPlanTest{
			{
				q:     "select * from xy where y-1 = (select u from uv limit 1 offset 5);",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{},
			},
			{
				q:     "select * from xy where x != (select u from uv limit 1 offset 5);",
				types: []plan.JoinType{plan.JoinTypeAnti},
				exp:   []sql.Row{},
			},
		},
	},
	{
		name: "unnest with scope filters",
		setup: []string{
			"create table ab (a int primary key, b int);",
			"create table rs (r int primary key, s int);",
			"CREATE table xy (x int primary key, y int);",
			"CREATE table uv (u int primary key, v int);",
			"insert into ab values (0,2), (1,2), (2,2), (3,1);",
			"insert into rs values (0,0), (1,0), (2,0), (4,4), (5,4);",
			"insert into xy values (1,0), (2,1), (0,2), (3,3);",
			"insert into uv values (0,1), (1,1), (2,2), (3,2);",
		},
		tests: []JoinPlanTest{
			{
				q: `
SELECT x
FROM xy 
WHERE EXISTS (SELECT count(v) AS count_1 
FROM uv 
WHERE y = v and v = 1 GROUP BY v
HAVING count(v) >= 1)`,
				types:   []plan.JoinType{},
				exp:     []sql.Row{{2}},
				skipOld: true,
			},
			{
				q:     "select * from xy where y-1 = (select u from uv where v = 2 order by 1 limit 1);",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{3, 3}},
			},
			{
				q:     "select * from xy where x != (select u from uv where v = 2 order by 1 limit 1) order by 1;",
				types: []plan.JoinType{plan.JoinTypeAnti},
				exp:   []sql.Row{{0, 2}, {1, 0}, {3, 3}},
			},
			{
				q:     "select * from xy where x != (select distinct u from uv where v = 2 order by 1 limit 1) order by 1;",
				types: []plan.JoinType{plan.JoinTypeAnti},
				exp:   []sql.Row{{0, 2}, {1, 0}, {3, 3}},
			},
			{
				q:     "select * from xy where (x,y+1) = (select u,v from uv where v = 2 order by 1 limit 1) order by 1;",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{2, 1}},
			},
			{
				q:     "select * from xy where x in (select cnt from (select count(u) as cnt from uv group by v having cnt > 0) sq) order by 1,2;",
				types: []plan.JoinType{plan.JoinTypeHash},
				exp:   []sql.Row{{2, 1}},
			},
			{
				q: `SELECT * FROM xy WHERE (
      				EXISTS (SELECT * FROM xy Alias1 WHERE Alias1.x = (xy.x + 1))
      				AND EXISTS (SELECT * FROM uv Alias2 WHERE Alias2.u = (xy.x + 2)));`,
				// These should both be JoinTypeSemiLookup, but for https://github.com/dolthub/go-mysql-server/issues/1893
				types: []plan.JoinType{plan.JoinTypeSemiLookup, plan.JoinTypeMerge},
				exp:   []sql.Row{{0, 2}, {1, 0}},
			},
			{
				q: `SELECT * FROM xy WHERE (
      				EXISTS (SELECT * FROM xy Alias1 WHERE Alias1.x = (xy.x + 1))
      				AND EXISTS (SELECT * FROM uv Alias1 WHERE Alias1.u = (xy.x + 2)));`,
				// These should both be JoinTypeSemiLookup, but for https://github.com/dolthub/go-mysql-server/issues/1893
				types: []plan.JoinType{plan.JoinTypeSemiLookup, plan.JoinTypeMerge},
				exp:   []sql.Row{{0, 2}, {1, 0}},
			},
			{
				q: `SELECT *
FROM ab A0
WHERE EXISTS (
    SELECT U0.a
    FROM
    (
        ab U0
        LEFT OUTER JOIN
        rs U1
        ON (U0.a = U1.s)
    )
    WHERE (U1.s IS NULL AND U0.a = A0.a)
);`,
				types: []plan.JoinType{plan.JoinTypeLookup, plan.JoinTypeLeftOuterHash},
				exp: []sql.Row{
					{1, 2},
					{2, 2},
					{3, 1},
				},
			},
			{
				q:     `select * from xy where exists (select * from uv) and x = 0`,
				types: []plan.JoinType{plan.JoinTypeCrossHash},
				exp:   []sql.Row{{0, 2}},
			},
			{
				q: `
select x from xy where
  not exists (select a from ab where a = x and a = 1) and
  not exists (select a from ab where a = x and a = 2)`,
				types: []plan.JoinType{plan.JoinTypeAntiLookup, plan.JoinTypeAntiLookup},
				exp:   []sql.Row{{0}, {3}},
			},
			{
				q: `
select * from xy where x in (
    with recursive tree(s) AS (
        SELECT 1
    )
    SELECT u FROM uv, tree where u = s
)`,
				types: []plan.JoinType{plan.JoinTypeHash, plan.JoinTypeHash},
				exp:   []sql.Row{{1, 0}},
			},
			{
				q: `
SELECT *
FROM xy
  WHERE
    EXISTS (
    SELECT 1
    FROM ab
    WHERE
      xy.x = ab.a AND
      EXISTS (
        SELECT 1
        FROM uv
        WHERE
          ab.a = uv.v
    )
  )`,
				types: []plan.JoinType{plan.JoinTypeLookup, plan.JoinTypeLookup},
				exp:   []sql.Row{{1, 0}, {2, 1}},
			},
			{
				q:     `select * from xy where exists (select * from uv join ab on u = a)`,
				types: []plan.JoinType{plan.JoinTypeCrossHash, plan.JoinTypeMerge},
				exp:   []sql.Row{{0, 2}, {1, 0}, {2, 1}, {3, 3}},
			},
		},
	},
	{
		name: "unnest non-equality comparisons",
		setup: []string{
			"CREATE table xy (x int primary key, y int);",
			"CREATE table uv (u int primary key, v int);",
			"insert into xy values (1,0), (2,1), (0,2), (3,3);",
			"insert into uv values (0,1), (1,1), (2,2), (3,2);",
		},
		tests: []JoinPlanTest{
			{
				q:     "select * from xy where y >= (select u from uv where u = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{0, 2}, {3, 3}},
			},
			{
				q:     "select * from xy where x <= (select u from uv where u = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{0, 2}, {1, 0}, {2, 1}},
			},
			{
				q:     "select * from xy where x < (select u from uv where u = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{0, 2}, {1, 0}},
			},
			{
				q:     "select * from xy where x > (select u from uv where u = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{3, 3}},
			},
			{
				q:     "select * from uv where v <=> (select u from uv where u = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{2, 2}, {3, 2}},
			},
		},
	},
	{
		name: "unnest twice-nested subquery",
		setup: []string{
			"CREATE table xy (x int primary key, y int);",
			"CREATE table uv (u int primary key, v int);",
			"insert into xy values (1,0), (2,1), (0,2), (3,3);",
			"insert into uv values (0,1), (1,1), (2,2), (3,2);",
		},
		tests: []JoinPlanTest{
			{
				q:     "select * from xy where x in (select * from (select 1) r where x = 1);",
				types: []plan.JoinType{plan.JoinTypeHash},
				exp:   []sql.Row{{1, 0}},
			},
			{
				q:     "select * from xy where x in (select 1 where 1 in (select 1 where 1 in (select 1 where x != 2)) and x = 1);",
				types: []plan.JoinType{plan.JoinTypeHash, plan.JoinTypeHash, plan.JoinTypeHash},
				exp:   []sql.Row{{1, 0}},
			},
			{
				q:     "select * from xy where x in (select * from (select 1 where 1 in (select 1 where x != 2)) r where x = 1);",
				types: []plan.JoinType{plan.JoinTypeHash, plan.JoinTypeHash},
				exp:   []sql.Row{{1, 0}},
			},
			{
				q:     "select * from xy where x in (select * from (select 1) r);",
				types: []plan.JoinType{plan.JoinTypeHash},
				exp:   []sql.Row{{1, 0}},
			},
			{
				q: `
with recursive rec(x) as (select 1 union select 1)
select * from xy where x in (
  select * from rec
);`,
				types: []plan.JoinType{plan.JoinTypeHash},
				exp:   []sql.Row{{1, 0}},
			},
			{
				q: `
with recursive rec(x) as (
  select 1
  union
  select rec.x from rec join xy on rec.x = xy.y
)
select * from uv
where u in (select * from rec);`,
				types:   []plan.JoinType{plan.JoinTypeHash, plan.JoinTypeHash},
				exp:     []sql.Row{{1, 1}},
				skipOld: true,
			},
			{
				q:     "select x+1 as newX, y from xy having y in (select x from xy where newX=1)",
				types: []plan.JoinType{},
				exp:   []sql.Row{{1, 2}},
			},
			{
				q:     "select x, x+1 as newX from xy having x in (select * from (select 1 where 1 in (select 1 where newX != 1)) r where x = 1);",
				types: []plan.JoinType{},
				exp:   []sql.Row{{1, 2}},
			},
			{
				q:   "select * from uv where not exists (select * from xy where u = 1)",
				exp: []sql.Row{{0, 1}, {2, 2}, {3, 2}},
			},
			{
				q:   "select * from uv where not exists (select * from xy where not exists (select * from xy where u = 1))",
				exp: []sql.Row{{1, 1}},
			},
			{
				q:   "select * from uv where not exists (select * from xy where not exists (select * from xy where u = 1 or v = 2))",
				exp: []sql.Row{{1, 1}, {2, 2}, {3, 2}},
			},
			{
				q:   "select * from uv where not exists (select * from xy where v = 1 and not exists (select * from xy where u = 1))",
				exp: []sql.Row{{1, 1}, {2, 2}, {3, 2}},
			},
			{
				q:   "select * from uv where not exists (select * from xy where not exists (select * from xy where not(u = 1)))",
				exp: []sql.Row{{0, 1}, {2, 2}, {3, 2}},
			},
		},
	},
	{
		name: "convert semi to inner join",
		setup: []string{
			"CREATE table xy (x int, y int, primary key(x,y));",
			"CREATE table uv (u int primary key, v int);",
			"CREATE table ab (a int primary key, b int);",
			"insert into xy values (1,0), (2,1), (0,2), (3,3);",
			"insert into uv values (0,1), (1,1), (2,2), (3,2);",
			"insert into ab values (0,2), (1,2), (2,2), (3,1);",
			"update information_schema.statistics set cardinality = 100 where table_name in ('xy', 'ab', 'uv') and table_schema = 'mydb';",
		},
		tests: []JoinPlanTest{
			{
				q:     "select * from xy where x in (select u from uv join ab on u = a and a = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeHash, plan.JoinTypeLookup},
				exp:   []sql.Row{{2, 1}},
			},
			{
				q: `select x from xy where x in (
	select (select u from uv where u = sq.a)
    from (select a from ab) sq);`,
				types: []plan.JoinType{plan.JoinTypeHash},
				exp:   []sql.Row{{0}, {1}, {2}, {3}},
			},
			{
				q:     "select * from xy where y >= (select u from uv where u = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{0, 2}, {3, 3}},
			},
			{
				q:     "select * from xy where x <= (select u from uv where u = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{0, 2}, {1, 0}, {2, 1}},
			},
			{
				q:     "select * from xy where x < (select u from uv where u = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{0, 2}, {1, 0}},
			},
			{
				q:     "select * from xy where x > (select u from uv where u = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{3, 3}},
			},
			{
				q:     "select * from uv where v <=> (select u from uv where u = 2) order by 1;",
				types: []plan.JoinType{plan.JoinTypeSemi},
				exp:   []sql.Row{{2, 2}, {3, 2}},
			},
		},
	},
	{
		name: "convert anti to left join",
		setup: []string{
			"CREATE table xy (x int, y int, primary key(x,y));",
			"CREATE table uv (u int primary key, v int);",
			"create table empty_tbl (a int, b int);",
			"insert into xy values (1,0), (2,1), (0,2), (3,3);",
			"insert into uv values (0,1), (1,1), (2,2), (3,2);",
		},
		// write a bunch of left joins and make sure they are converted to anti joins
		tests: []JoinPlanTest{
			{
				q:     "select /*+ HASH_JOIN(xy,scalarSubq0) */ * from xy where x not in (select a from empty_tbl) order by x",
				types: []plan.JoinType{plan.JoinTypeLeftOuterHashExcludeNulls},
				exp: []sql.Row{
					{0, 2},
					{1, 0},
					{2, 1},
					{3, 3},
				},
			},
			{
				q:     "select /*+ HASH_JOIN(xy,scalarSubq0) */ * from xy where x not in (select v from uv) order by x",
				types: []plan.JoinType{plan.JoinTypeLeftOuterHashExcludeNulls},
				exp: []sql.Row{
					{0, 2},
					{3, 3},
				},
			},
			{
				q:     "select /*+ HASH_JOIN(xy,scalarSubq0) */ * from xy where x not in (select v from uv where u = 2) order by x",
				types: []plan.JoinType{plan.JoinTypeLeftOuterHashExcludeNulls},
				exp: []sql.Row{
					{0, 2},
					{1, 0},
					{3, 3},
				},
			},
			{
				q:     "select /*+ HASH_JOIN(xy,scalarSubq0) */ * from xy where x != (select v from uv where u = 2) order by x",
				types: []plan.JoinType{plan.JoinTypeLeftOuterHashExcludeNulls},
				exp: []sql.Row{
					{0, 2},
					{1, 0},
					{3, 3},
				},
			},
			{
				q:     "select * from xy where not exists (select * from empty_tbl) order by x",
				types: []plan.JoinType{plan.JoinTypeLeftOuter},
				exp: []sql.Row{
					{0, 2},
					{1, 0},
					{2, 1},
					{3, 3},
				},
			},
			{
				q:     "select /*+ MERGE_JOIN(xy,scalarSubq0) */ * from xy where x not in (select u from uv WHERE u = 2) order by x",
				types: []plan.JoinType{plan.JoinTypeAntiLookup},
				exp: []sql.Row{
					{0, 2},
					{1, 0},
					{3, 3},
				},
			},
			{
				q:     "select /*+ LEFT_OUTER_LOOKUP_JOIN(xy,scalarSubq0) */ * from xy where x not in (select u from uv WHERE u = 2) order by x",
				types: []plan.JoinType{plan.JoinTypeLeftOuterLookup},
				exp: []sql.Row{
					{0, 2},
					{1, 0},
					{3, 3},
				},
			},
		},
	},
	{
		name: "join concat tests",
		setup: []string{
			"CREATE table xy (x int primary key, y int);",
			"CREATE table uv (u int primary key, v int);",
			"insert into xy values (1,0), (2,1), (0,2), (3,3);",
			"insert into uv values (0,1), (1,1), (2,2), (3,2);",
			"update information_schema.statistics set cardinality = 100 where table_name in ('xy', 'uv');",
		},
		tests: []JoinPlanTest{
			{
				q:     "select x, u from xy inner join uv on u+1 = x OR u+2 = x OR u+3 = x;",
				types: []plan.JoinType{plan.JoinTypeLookup},
				exp:   []sql.Row{{3, 0}, {2, 0}, {1, 0}, {3, 1}, {2, 1}, {3, 2}},
			},
		},
	},
	{
		name: "join order hint",
		setup: []string{
			"CREATE table xy (x int primary key, y int);",
			"CREATE table uv (u int primary key, v int);",
			"insert into xy values (1,0), (2,1), (0,2), (3,3);",
			"insert into uv values (0,1), (1,1), (2,2), (3,2);",
			"update information_schema.statistics set cardinality = 100 where table_name in ('xy', 'uv');",
		},
		tests: []JoinPlanTest{
			{
				q:     "select /*+ JOIN_ORDER(b, c, a) */ 1 from xy a join xy b on a.x+3 = b.x join xy c on a.x+3 = c.x and a.x+3 = b.x",
				order: []string{"b", "c", "a"},
			},
			{
				q:     "select /*+ JOIN_ORDER(a, c, b) */ 1 from xy a join xy b on a.x+3 = b.x join xy c on a.x+3 = c.x and a.x+3 = b.x",
				order: []string{"a", "c", "b"},
			},
			{
				q:     "select /*+ JOIN_ORDER(a,c,b) */ 1 from xy a join xy b on a.x+3 = b.x WHERE EXISTS (select 1 from uv c where c.u = a.x+2)",
				order: []string{"a", "c", "b"},
			},
			{
				q:     "select /*+ JOIN_ORDER(b,c,a) */ 1 from xy a join xy b on a.x+3 = b.x WHERE EXISTS (select 1 from uv c where c.u = a.x+2)",
				order: []string{"b", "c", "a"},
			},
			{
				q:     "select /*+ JOIN_ORDER(b,scalarSubq0,a) */ 1 from xy a join xy b on a.x+3 = b.x WHERE a.x in (select u from uv c)",
				order: []string{"b", "scalarSubq0", "a"},
			},
		},
	},
	{
		name: "join op hint",
		setup: []string{
			"CREATE table xy (x int primary key, y int);",
			"CREATE table uv (u int primary key, v int, key(v));",
			"insert into xy values (1,0), (2,1), (0,2), (3,3);",
			"insert into uv values (0,1), (1,1), (2,2), (3,2);",
		},
		tests: []JoinPlanTest{
			{
				q:     "select /*+ LOOKUP_JOIN(xy,uv) */ 1 from xy join uv on x = u",
				types: []plan.JoinType{plan.JoinTypeLookup},
			},
			{
				q:     "select /*+ MERGE_JOIN(xy,uv) */ 1 from xy join uv on x = u",
				types: []plan.JoinType{plan.JoinTypeMerge},
			},
			{
				q:     "select /*+ INNER_JOIN(xy,uv) */ 1 from xy join uv on x = u",
				types: []plan.JoinType{plan.JoinTypeInner},
			},
			{
				q:     "select /*+ HASH_JOIN(xy,uv) */ 1 from xy join uv on x = u",
				types: []plan.JoinType{plan.JoinTypeHash},
			},
			{
				q:     "select /*+ JOIN_ORDER(a,b,c) HASH_JOIN(a,b) HASH_JOIN(b,c) */ 1 from xy a join uv b on a.x = b.u join xy c on b.u = c.x",
				types: []plan.JoinType{plan.JoinTypeHash, plan.JoinTypeHash},
				order: []string{"a", "b", "c"},
			},
			{
				q:     "select /*+ LOOKUP_JOIN(b,a) HASH_JOIN(b,c) */ 1 from xy a join uv b on a.x = b.u join xy c on b.u = c.x",
				types: []plan.JoinType{plan.JoinTypeLookup, plan.JoinTypeHash},
			},
			{
				q:     "select /*+ LOOKUP_JOIN(b,a) MERGE_JOIN(b,c) */ 1 from xy a join uv b on a.x = b.u join xy c on b.u = c.x",
				types: []plan.JoinType{plan.JoinTypeLookup, plan.JoinTypeMerge},
			},
			{
				q:     "select /*+ JOIN_ORDER(b,c,a) LOOKUP_JOIN(b,a) MERGE_JOIN(b,c) */ 1 from xy a join uv b on a.x = b.u join xy c on b.u = c.x",
				types: []plan.JoinType{plan.JoinTypeLookup, plan.JoinTypeMerge},
				order: []string{"b", "c", "a"},
			},
			{
				q:     "select /*+ JOIN_ORDER(a,b,c) LOOKUP_JOIN(b,a) HASH_JOIN(b,c) */ 1 from xy a join uv b on a.x = b.u join xy c on b.u = c.x",
				types: []plan.JoinType{plan.JoinTypeHash, plan.JoinTypeLookup},
				order: []string{"a", "b", "c"},
			},
			{
				q:     "select /*+ JOIN_ORDER(c,a,b) MERGE_JOIN(a,b) HASH_JOIN(b,c) */ 1 from xy a join uv b on a.x = b.u join xy c on b.u = c.x",
				types: []plan.JoinType{plan.JoinTypeHash, plan.JoinTypeMerge},
				order: []string{"c", "a", "b"},
			},
			{
				q: `
select /*+ JOIN_ORDER(d,c,b,a) MERGE_JOIN(d,c) MERGE_JOIN(b,a) INNER_JOIN(c,a)*/ 1
from xy a
join uv b on a.x = b.u
join xy c on a.x = c.x
join uv d on d.u = c.x`,
				types: []plan.JoinType{plan.JoinTypeInner, plan.JoinTypeMerge, plan.JoinTypeMerge},
				order: []string{"d", "c", "b", "a"},
			},
			{
				q: `
select /*+ JOIN_ORDER(a,b,c,d) LOOKUP_JOIN(d,c) MERGE_JOIN(b,a) HASH_JOIN(c,a)*/ 1
from xy a
join uv b on a.x = b.u
join xy c on a.x = c.x
join uv d on d.u = c.x`,
				types: []plan.JoinType{plan.JoinTypeLookup, plan.JoinTypeHash, plan.JoinTypeMerge},
				order: []string{"a", "b", "c", "d"},
			},
			{
				q: "select /*+ LOOKUP_JOIN(xy,scalarSubq0) */ 1 from xy where x not in (select u from uv)",
				// This should be JoinTypeSemiLookup, but for https://github.com/dolthub/go-mysql-server/issues/1894
				types: []plan.JoinType{plan.JoinTypeAntiLookup},
			},
			{
				q:     "select /*+ ANTI_JOIN(xy,scalarSubq0) */ 1 from xy where x not in (select u from uv)",
				types: []plan.JoinType{plan.JoinTypeAnti},
			},
			{
				q:     "select /*+ LOOKUP_JOIN(xy,scalarSubq0) */ 1 from xy where x in (select u from uv)",
				types: []plan.JoinType{plan.JoinTypeLookup},
			},
			{
				q:     "select /*+ SEMI_JOIN(xy,scalarSubq0) */ 1 from xy where x in (select u from uv)",
				types: []plan.JoinType{plan.JoinTypeSemi},
			},
			{
				q:     "select /*+ LOOKUP_JOIN(s,scalarSubq0) */ 1 from xy s where x in (select u from uv)",
				types: []plan.JoinType{plan.JoinTypeLookup},
			},
			{
				q:     "select /*+ SEMI_JOIN(s,scalarSubq0) */ 1 from xy s where x in (select u from uv)",
				types: []plan.JoinType{plan.JoinTypeSemi},
			},
		},
	},
	{
		// This is a regression test for https://github.com/dolthub/go-mysql-server/pull/1889.
		// We should always prefer a more specific index over a less specific index for lookups.
		name: "lookup join multiple indexes",
		setup: []string{
			"create table lhs (a int, b int, c int);",
			"create table rhs (a int, b int, c int, d int, index a_idx(a), index abcd_idx(a,b,c,d));",
			"insert into lhs values (0, 0, 0), (0, 0, 1), (0, 1, 1), (1, 1, 1);",
			"insert into rhs values " +
				"(0, 0, 0, 0)," +
				"(0, 0, 0, 1)," +
				"(0, 0, 1, 0)," +
				"(0, 0, 1, 1)," +
				"(0, 1, 0, 0)," +
				"(0, 1, 0, 1)," +
				"(0, 1, 1, 0)," +
				"(0, 1, 1, 1)," +
				"(1, 0, 0, 0)," +
				"(1, 0, 0, 1)," +
				"(1, 0, 1, 0)," +
				"(1, 0, 1, 1)," +
				"(1, 1, 0, 0)," +
				"(1, 1, 0, 1)," +
				"(1, 1, 1, 0)," +
				"(1, 1, 1, 1);",
		},
		tests: []JoinPlanTest{
			{
				q:       "select rhs.* from lhs left join rhs on lhs.a = rhs.a and lhs.b = rhs.b and lhs.c = rhs.c",
				types:   []plan.JoinType{plan.JoinTypeLeftOuterLookup},
				indexes: []string{"abcd_idx"},
				exp: []sql.Row{
					{0, 0, 0, 0},
					{0, 0, 0, 1},
					{0, 0, 1, 0},
					{0, 0, 1, 1},
					{0, 1, 1, 0},
					{0, 1, 1, 1},
					{1, 1, 1, 0},
					{1, 1, 1, 1},
				},
			},
		},
	},
	{
		name: "primary key range join",
		setup: []string{
			"create table vals (val int primary key)",
			"create table ranges (min int primary key, max int, unique key(min,max))",
			"insert into vals values (0), (1), (2), (3), (4), (5), (6)",
			"insert into ranges values (0,2), (1,3), (2,4), (3,5), (4,6)",
		},
		tests: []JoinPlanTest{
			{
				q:     "select * from vals join ranges on val between min and max",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{0, 0, 2},
					{1, 0, 2},
					{1, 1, 3},
					{2, 0, 2},
					{2, 1, 3},
					{2, 2, 4},
					{3, 1, 3},
					{3, 2, 4},
					{3, 3, 5},
					{4, 2, 4},
					{4, 3, 5},
					{4, 4, 6},
					{5, 3, 5},
					{5, 4, 6},
					{6, 4, 6},
				},
			},
			{
				q:     "select * from vals join ranges on val > min and val < max",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{1, 0, 2},
					{2, 1, 3},
					{3, 2, 4},
					{4, 3, 5},
					{5, 4, 6},
				},
			},
			{
				q:     "select * from vals join ranges on val >= min and val < max",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{0, 0, 2},
					{1, 0, 2},
					{1, 1, 3},
					{2, 1, 3},
					{2, 2, 4},
					{3, 2, 4},
					{3, 3, 5},
					{4, 3, 5},
					{4, 4, 6},
					{5, 4, 6},
				},
			},
			{
				q:     "select * from vals join ranges on val > min and val <= max",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{1, 0, 2},
					{2, 0, 2},
					{2, 1, 3},
					{3, 1, 3},
					{3, 2, 4},
					{4, 2, 4},
					{4, 3, 5},
					{5, 3, 5},
					{5, 4, 6},
					{6, 4, 6},
				},
			},
			{
				q:     "select * from vals join ranges on val >= min and val <= max",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{0, 0, 2},
					{1, 0, 2},
					{1, 1, 3},
					{2, 0, 2},
					{2, 1, 3},
					{2, 2, 4},
					{3, 1, 3},
					{3, 2, 4},
					{3, 3, 5},
					{4, 2, 4},
					{4, 3, 5},
					{4, 4, 6},
					{5, 3, 5},
					{5, 4, 6},
					{6, 4, 6},
				},
			},
		},
	},
	{
		name: "keyless range join",
		setup: []string{
			"create table vals (val int)",
			"create table ranges (min int, max int)",
			"insert into vals values (0), (1), (2), (3), (4), (5), (6)",
			"insert into ranges values (0,2), (1,3), (2,4), (3,5), (4,6)",
		},
		tests: []JoinPlanTest{
			{
				q:     "select * from vals join ranges on val between min and max",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{0, 0, 2},
					{1, 0, 2},
					{1, 1, 3},
					{2, 0, 2},
					{2, 1, 3},
					{2, 2, 4},
					{3, 1, 3},
					{3, 2, 4},
					{3, 3, 5},
					{4, 2, 4},
					{4, 3, 5},
					{4, 4, 6},
					{5, 3, 5},
					{5, 4, 6},
					{6, 4, 6},
				},
			},
			{
				q:     "select * from vals join ranges on val > min and val < max",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{1, 0, 2},
					{2, 1, 3},
					{3, 2, 4},
					{4, 3, 5},
					{5, 4, 6},
				},
			},
			{
				q:     "select * from vals join ranges on min < val and max > val",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{1, 0, 2},
					{2, 1, 3},
					{3, 2, 4},
					{4, 3, 5},
					{5, 4, 6},
				},
			},
			{
				q:     "select * from vals join ranges on val >= min and val < max",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{0, 0, 2},
					{1, 0, 2},
					{1, 1, 3},
					{2, 1, 3},
					{2, 2, 4},
					{3, 2, 4},
					{3, 3, 5},
					{4, 3, 5},
					{4, 4, 6},
					{5, 4, 6},
				},
			},
			{
				q:     "select * from vals join ranges on val > min and val <= max",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{1, 0, 2},
					{2, 0, 2},
					{2, 1, 3},
					{3, 1, 3},
					{3, 2, 4},
					{4, 2, 4},
					{4, 3, 5},
					{5, 3, 5},
					{5, 4, 6},
					{6, 4, 6},
				},
			},
			{
				q:     "select * from vals join ranges on val >= min and val <= max",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{0, 0, 2},
					{1, 0, 2},
					{1, 1, 3},
					{2, 0, 2},
					{2, 1, 3},
					{2, 2, 4},
					{3, 1, 3},
					{3, 2, 4},
					{3, 3, 5},
					{4, 2, 4},
					{4, 3, 5},
					{4, 4, 6},
					{5, 3, 5},
					{5, 4, 6},
					{6, 4, 6},
				},
			},
			{
				q:     "select * from vals left join ranges on val > min and val < max",
				types: []plan.JoinType{plan.JoinTypeLeftOuterRangeHeap},
				exp: []sql.Row{
					{0, nil, nil},
					{1, 0, 2},
					{2, 1, 3},
					{3, 2, 4},
					{4, 3, 5},
					{5, 4, 6},
					{6, nil, nil},
				},
			},
			{
				q:     "select * from ranges l join ranges r on l.min > r.min and l.min < r.max",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{1, 3, 0, 2},
					{2, 4, 1, 3},
					{3, 5, 2, 4},
					{4, 6, 3, 5},
				},
			},
			{
				q:     "select * from vals left join ranges r1 on val > r1.min and val < r1.max left join ranges r2 on r1.min > r2.min and r1.min < r2.max",
				types: []plan.JoinType{plan.JoinTypeLeftOuterRangeHeap, plan.JoinTypeLeftOuterRangeHeap},
				exp: []sql.Row{
					{0, nil, nil, nil, nil},
					{1, 0, 2, nil, nil},
					{2, 1, 3, 0, 2},
					{3, 2, 4, 1, 3},
					{4, 3, 5, 2, 4},
					{5, 4, 6, 3, 5},
					{6, nil, nil, nil, nil},
				},
			},
			{
				q:     "select * from (select vals.val * 2 as val from vals) as newVals join (select ranges.min * 2 as min, ranges.max * 2 as max from ranges) as newRanges on val > min and val < max;",
				types: []plan.JoinType{plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{2, 0, 4},
					{4, 2, 6},
					{6, 4, 8},
					{8, 6, 10},
					{10, 8, 12},
				},
			},
			{
				// This tests that the RangeHeapJoin node functions correctly even if its rows are iterated over multiple times.
				q:     "select * from (select 1 union select 2) as l left join (select * from vals join ranges on val > min and val < max) as r on max = max",
				types: []plan.JoinType{plan.JoinTypeLeftOuter, plan.JoinTypeRangeHeap},
				exp: []sql.Row{
					{1, 1, 0, 2},
					{1, 2, 1, 3},
					{1, 3, 2, 4},
					{1, 4, 3, 5},
					{1, 5, 4, 6},
					{2, 1, 0, 2},
					{2, 2, 1, 3},
					{2, 3, 2, 4},
					{2, 4, 3, 5},
					{2, 5, 4, 6},
				},
			},
			{
				q:     "select * from vals left join (select * from ranges where 0) as newRanges on val > min and val < max;",
				types: []plan.JoinType{plan.JoinTypeLeftOuterRangeHeap},
				exp: []sql.Row{
					{0, nil, nil},
					{1, nil, nil},
					{2, nil, nil},
					{3, nil, nil},
					{4, nil, nil},
					{5, nil, nil},
					{6, nil, nil},
				},
			},
		},
	},
	{
		name: "range join vs good lookup join regression test",
		setup: []string{
			"create table vals (val int, filter1 int, filter2 int, filter3 int)",
			"create table ranges (min int, max int, filter1 int, filter2 int, filter3 int, key filters (filter1, filter2, filter3))",
			"insert into vals values (0, 0, 0, 0), " +
				"(1, 0, 0, 0), " +
				"(2, 0, 0, 0), " +
				"(3, 0, 0, 0), " +
				"(4, 0, 0, 0), " +
				"(5, 0, 0, 0), " +
				"(6, 0, 0, 0), " +
				"(0, 0, 0, 1), " +
				"(1, 0, 0, 1), " +
				"(2, 0, 0, 1), " +
				"(3, 0, 0, 1), " +
				"(4, 0, 0, 1), " +
				"(5, 0, 0, 1), " +
				"(6, 0, 0, 1), " +
				"(0, 0, 1, 0), " +
				"(1, 0, 1, 0), " +
				"(2, 0, 1, 0), " +
				"(3, 0, 1, 0), " +
				"(4, 0, 1, 0), " +
				"(5, 0, 1, 0), " +
				"(6, 0, 1, 0), " +
				"(0, 0, 1, 1), " +
				"(1, 0, 1, 1), " +
				"(2, 0, 1, 1), " +
				"(3, 0, 1, 1), " +
				"(4, 0, 1, 1), " +
				"(5, 0, 1, 1), " +
				"(6, 0, 1, 1), " +
				"(0, 1, 0, 0), " +
				"(1, 1, 0, 0), " +
				"(2, 1, 0, 0), " +
				"(3, 1, 0, 0), " +
				"(4, 1, 0, 0), " +
				"(5, 1, 0, 0), " +
				"(6, 1, 0, 0), " +
				"(0, 1, 0, 1), " +
				"(1, 1, 0, 1), " +
				"(2, 1, 0, 1), " +
				"(3, 1, 0, 1), " +
				"(4, 1, 0, 1), " +
				"(5, 1, 0, 1), " +
				"(6, 1, 0, 1), " +
				"(0, 1, 1, 0), " +
				"(1, 1, 1, 0), " +
				"(2, 1, 1, 0), " +
				"(3, 1, 1, 0), " +
				"(4, 1, 1, 0), " +
				"(5, 1, 1, 0), " +
				"(6, 1, 1, 0), " +
				"(0, 1, 1, 1), " +
				"(1, 1, 1, 1), " +
				"(2, 1, 1, 1), " +
				"(3, 1, 1, 1), " +
				"(4, 1, 1, 1), " +
				"(5, 1, 1, 1), " +
				"(6, 1, 1, 1);",
			"insert into ranges values " +
				"(0, 2, 0, 0, 0), " +
				"(1, 3, 0, 0, 0), " +
				"(2, 4, 0, 0, 0), " +
				"(3, 5, 0, 0, 0), " +
				"(4, 6, 0, 0, 0), " +
				"(0, 2, 0, 0, 1), " +
				"(1, 3, 0, 0, 1), " +
				"(2, 4, 0, 0, 1), " +
				"(3, 5, 0, 0, 1), " +
				"(4, 6, 0, 0, 1), " +
				"(0, 2, 0, 1, 0), " +
				"(1, 3, 0, 1, 0), " +
				"(2, 4, 0, 1, 0), " +
				"(3, 5, 0, 1, 0), " +
				"(4, 6, 0, 1, 0), " +
				"(0, 2, 0, 1, 1), " +
				"(1, 3, 0, 1, 1), " +
				"(2, 4, 0, 1, 1), " +
				"(3, 5, 0, 1, 1), " +
				"(4, 6, 0, 1, 1), " +
				"(0, 2, 1, 0, 0), " +
				"(1, 3, 1, 0, 0), " +
				"(2, 4, 1, 0, 0), " +
				"(3, 5, 1, 0, 0), " +
				"(4, 6, 1, 0, 0), " +
				"(0, 2, 1, 0, 1), " +
				"(1, 3, 1, 0, 1), " +
				"(2, 4, 1, 0, 1), " +
				"(3, 5, 1, 0, 1), " +
				"(4, 6, 1, 0, 1), " +
				"(0, 2, 1, 1, 0), " +
				"(1, 3, 1, 1, 0), " +
				"(2, 4, 1, 1, 0), " +
				"(3, 5, 1, 1, 0), " +
				"(4, 6, 1, 1, 0), " +
				"(0, 2, 1, 1, 1), " +
				"(1, 3, 1, 1, 1), " +
				"(2, 4, 1, 1, 1), " +
				"(3, 5, 1, 1, 1), " +
				"(4, 6, 1, 1, 1); ",
		},
		tests: []JoinPlanTest{
			{
				// Test that a RangeHeapJoin won't be chosen over a LookupJoin with a multiple-column index.
				q:     "select val, min, max, vals.filter1, vals.filter2, vals.filter3 from vals join ranges on val > min and val < max and vals.filter1 = ranges.filter1 and vals.filter2 = ranges.filter2 and vals.filter3 = ranges.filter3",
				types: []plan.JoinType{plan.JoinTypeLookup},
				exp: []sql.Row{
					{1, 0, 2, 0, 0, 0},
					{2, 1, 3, 0, 0, 0},
					{3, 2, 4, 0, 0, 0},
					{4, 3, 5, 0, 0, 0},
					{5, 4, 6, 0, 0, 0},
					{1, 0, 2, 0, 0, 1},
					{2, 1, 3, 0, 0, 1},
					{3, 2, 4, 0, 0, 1},
					{4, 3, 5, 0, 0, 1},
					{5, 4, 6, 0, 0, 1},
					{1, 0, 2, 0, 1, 0},
					{2, 1, 3, 0, 1, 0},
					{3, 2, 4, 0, 1, 0},
					{4, 3, 5, 0, 1, 0},
					{5, 4, 6, 0, 1, 0},
					{1, 0, 2, 0, 1, 1},
					{2, 1, 3, 0, 1, 1},
					{3, 2, 4, 0, 1, 1},
					{4, 3, 5, 0, 1, 1},
					{5, 4, 6, 0, 1, 1},
					{1, 0, 2, 1, 0, 0},
					{2, 1, 3, 1, 0, 0},
					{3, 2, 4, 1, 0, 0},
					{4, 3, 5, 1, 0, 0},
					{5, 4, 6, 1, 0, 0},
					{1, 0, 2, 1, 0, 1},
					{2, 1, 3, 1, 0, 1},
					{3, 2, 4, 1, 0, 1},
					{4, 3, 5, 1, 0, 1},
					{5, 4, 6, 1, 0, 1},
					{1, 0, 2, 1, 1, 0},
					{2, 1, 3, 1, 1, 0},
					{3, 2, 4, 1, 1, 0},
					{4, 3, 5, 1, 1, 0},
					{5, 4, 6, 1, 1, 0},
					{1, 0, 2, 1, 1, 1},
					{2, 1, 3, 1, 1, 1},
					{3, 2, 4, 1, 1, 1},
					{4, 3, 5, 1, 1, 1},
					{5, 4, 6, 1, 1, 1},
				},
			},
		},
	},
}

func TestJoinPlanning(t *testing.T, harness Harness) {
	for _, tt := range JoinPlanningTests {
		t.Run(tt.name, func(t *testing.T) {
			harness.Setup([]setup.SetupScript{setup.MydbData[0], tt.setup})
			e := mustNewEngine(t, harness)
			defer e.Close()
			for _, tt := range tt.tests {
				if tt.types != nil {
					evalJoinTypeTest(t, harness, e, tt)
				}
				if tt.indexes != nil {
					evalIndexTest(t, harness, e, tt.q, tt.indexes, tt.skipOld)
				}
				if tt.exp != nil {
					evalJoinCorrectness(t, harness, e, tt.q, tt.q, tt.exp, tt.skipOld)
				}
				if tt.order != nil {
					evalJoinOrder(t, harness, e, tt.q, tt.order, tt.skipOld)
				}
			}
		})
	}
}

func evalJoinTypeTest(t *testing.T, harness Harness, e *sqle.Engine, tt JoinPlanTest) {
	t.Run(tt.q+" join types", func(t *testing.T) {
		if tt.skipOld {
			t.Skip()
		}

		ctx := NewContext(harness)
		ctx = ctx.WithQuery(tt.q)

		a, err := e.AnalyzeQuery(ctx, tt.q)
		require.NoError(t, err)

		jts := collectJoinTypes(a)
		var exp []string
		for _, t := range tt.types {
			exp = append(exp, t.String())
		}
		var cmp []string
		for _, t := range jts {
			cmp = append(cmp, t.String())
		}
		require.Equal(t, exp, cmp, fmt.Sprintf("unexpected plan:\n%s", sql.DebugString(a)))
	})
}

func evalIndexTest(t *testing.T, harness Harness, e *sqle.Engine, q string, indexes []string, skip bool) {
	t.Run(q+" join indexes", func(t *testing.T) {
		if skip {
			t.Skip()
		}

		ctx := NewContext(harness)
		ctx = ctx.WithQuery(q)

		a, err := e.AnalyzeQuery(ctx, q)
		require.NoError(t, err)

		idxs := collectIndexes(a)
		var exp []string
		for _, i := range indexes {
			exp = append(exp, i)
		}
		var cmp []string
		for _, i := range idxs {
			cmp = append(cmp, i.ID())
		}
		require.Equal(t, exp, cmp, fmt.Sprintf("unexpected plan:\n%s", sql.DebugString(a)))
	})
}

func evalJoinCorrectness(t *testing.T, harness Harness, e *sqle.Engine, name, q string, exp []sql.Row, skipOld bool) {
	t.Run(name, func(t *testing.T) {
		ctx := NewContext(harness)
		ctx = ctx.WithQuery(q)

		sch, iter, err := e.QueryWithBindings(ctx, q, nil)
		require.NoError(t, err, "Unexpected error for query %s: %s", q, err)

		rows, err := sql.RowIterToRows(ctx, sch, iter)
		require.NoError(t, err, "Unexpected error for query %s: %s", q, err)

		if exp != nil {
			checkResults(t, exp, nil, sch, rows, q)
		}

		require.Equal(t, 0, ctx.Memory.NumCaches())
		validateEngine(t, ctx, harness, e)
	})
}

func collectJoinTypes(n sql.Node) []plan.JoinType {
	var types []plan.JoinType
	transform.Inspect(n, func(n sql.Node) bool {
		if n == nil {
			return true
		}
		j, ok := n.(*plan.JoinNode)
		if ok {
			types = append(types, j.Op)
		}

		if ex, ok := n.(sql.Expressioner); ok {
			for _, e := range ex.Expressions() {
				transform.InspectExpr(e, func(e sql.Expression) bool {
					sq, ok := e.(*plan.Subquery)
					if !ok {
						return false
					}
					types = append(types, collectJoinTypes(sq.Query)...)
					return false
				})
			}
		}
		return true
	})
	return types
}

func collectIndexes(n sql.Node) []sql.Index {
	var indexes []sql.Index
	transform.Inspect(n, func(n sql.Node) bool {
		if n == nil {
			return true
		}
		access, ok := n.(*plan.IndexedTableAccess)
		if ok {
			indexes = append(indexes, access.Index())
			return true
		}

		if ex, ok := n.(sql.Expressioner); ok {
			for _, e := range ex.Expressions() {
				transform.InspectExpr(e, func(e sql.Expression) bool {
					sq, ok := e.(*plan.Subquery)
					if !ok {
						return false
					}
					indexes = append(indexes, collectIndexes(sq.Query)...)
					return false
				})
			}
		}
		return true
	})
	return indexes
}

func evalJoinOrder(t *testing.T, harness Harness, e *sqle.Engine, q string, exp []string, skipOld bool) {
	t.Run(q+" join order", func(t *testing.T) {
		ctx := NewContext(harness)
		ctx = ctx.WithQuery(q)

		a, err := e.AnalyzeQuery(ctx, q)
		require.NoError(t, err)

		cmp := collectJoinOrder(a)
		require.Equal(t, exp, cmp, fmt.Sprintf("expected order '%s' found '%s'\ndetail:\n%s", strings.Join(exp, ","), strings.Join(cmp, ","), sql.DebugString(a)))
	})
}

func collectJoinOrder(n sql.Node) []string {
	order := []string{}

	switch n := n.(type) {
	case *plan.JoinNode:
		order = append(order, collectJoinOrder(n.Left())...)
		order = append(order, collectJoinOrder(n.Right())...)
	case *plan.TableAlias:
		order = append(order, n.Name())
	default:
		children := n.Children()
		for _, c := range children {
			order = append(order, collectJoinOrder(c)...)
		}
	}

	return order
}
