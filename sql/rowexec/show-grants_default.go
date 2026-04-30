//go:build !tinygo

package rowexec

import (
	"fmt"
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/mysql_db"
	"github.com/dolthub/go-mysql-server/sql/plan"
)

func (b *BaseBuilder) buildShowGrantsImpl(ctx *sql.Context, n *plan.ShowGrants) (sql.RowIter, error) {
	mysqlDb, ok := n.MySQLDb.(*mysql_db.MySQLDb)
	if !ok {
		return nil, sql.ErrDatabaseNotFound.New("mysql")
	}
	if n.For == nil || n.CurrentUser {
		client := ctx.Session.Client()
		n.For = &plan.UserName{Name: client.User, Host: client.Address}
	}

	reader := mysqlDb.Reader()
	defer reader.Close()

	user := mysqlDb.GetUser(reader, n.For.Name, n.For.Host, false)
	if user == nil {
		return nil, sql.ErrShowGrantsUserDoesNotExist.New(n.For.Name, n.For.Host)
	}

	var rows []sql.Row
	userStr := user.UserHostToString("`")
	privStr := generatePrivStrings("*", "*", userStr, user.PrivilegeSet.ToSlice())
	rows = append(rows, sql.Row{privStr})

	for _, db := range user.PrivilegeSet.GetDatabases() {
		dbStr := fmt.Sprintf("`%s`", db.Name())
		if privStr = generatePrivStrings(dbStr, "*", userStr, db.ToSlice()); len(privStr) != 0 {
			rows = append(rows, sql.Row{privStr})
		}
		for _, tbl := range db.GetTables() {
			tblStr := fmt.Sprintf("`%s`", tbl.Name())
			privStr = generatePrivStrings(dbStr, tblStr, userStr, tbl.ToSlice())
			rows = append(rows, sql.Row{privStr})
		}
		for _, routine := range db.GetRoutines() {
			quotedRoutine := fmt.Sprintf("`%s`", routine.RoutineName())
			privStr = generateRoutinePrivStrings(dbStr, quotedRoutine, routine.RoutineType(), userStr, routine.ToSlice())
			rows = append(rows, sql.Row{privStr})
		}
	}

	sb := strings.Builder{}
	roleEdges := reader.GetToUserRoleEdges(mysql_db.RoleEdgesToKey{ToHost: user.Host, ToUser: user.User})
	for i, roleEdge := range roleEdges {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(roleEdge.FromString("`"))
	}
	if sb.Len() > 0 {
		rows = append(rows, sql.Row{fmt.Sprintf("GRANT %s TO %s", sb.String(), user.UserHostToString("`"))})
	}

	sb.Reset()
	for i, dynamicPrivWithWgo := range user.PrivilegeSet.ToSliceDynamic(true) {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(dynamicPrivWithWgo)
	}
	if sb.Len() > 0 {
		rows = append(rows, sql.Row{fmt.Sprintf("GRANT %s ON *.* TO %s WITH GRANT OPTION", sb.String(), user.UserHostToString("`"))})
	}

	sb.Reset()
	for i, dynamicPrivWithoutWgo := range user.PrivilegeSet.ToSliceDynamic(false) {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(dynamicPrivWithoutWgo)
	}
	if sb.Len() > 0 {
		rows = append(rows, sql.Row{fmt.Sprintf("GRANT %s ON *.* TO %s", sb.String(), user.UserHostToString("`"))})
	}
	return sql.RowsToRowIter(rows...), nil
}
