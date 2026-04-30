//go:build !tinygo

package rowexec

import (
	"time"

	"github.com/dolthub/go-mysql-server/sql/mysql"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/mysql_db"
	"github.com/dolthub/go-mysql-server/sql/plan"
	"github.com/dolthub/go-mysql-server/sql/types"
)

func (b *BaseBuilder) buildAlterUserImpl(ctx *sql.Context, a *plan.AlterUser) (sql.RowIter, error) {
	mysqlDb, ok := a.MySQLDb.(*mysql_db.MySQLDb)
	if !ok {
		return nil, sql.ErrDatabaseNotFound.New("mysql")
	}
	editor := mysqlDb.Editor()
	defer editor.Close()

	user := a.User
	if user.UserName.Host == "" {
		user.UserName.Host = "%"
	}

	userPk := mysql_db.UserPrimaryKey{
		Host: user.UserName.Host,
		User: user.UserName.Name,
	}
	previousUserEntry, ok := editor.GetUser(userPk)
	if !ok {
		if a.IfExists {
			return sql.RowsToRowIter(sql.Row{types.NewOkResult(0)}), nil
		}
		return nil, sql.ErrUserAlterFailure.New(user.UserName.String("'"))
	}

	plugin := previousUserEntry.Plugin
	authString := previousUserEntry.AuthString
	if user.Auth1 != nil {
		plugin = user.Auth1.Plugin()
		var err error
		authString, err = user.Auth1.AuthString()
		if err != nil {
			return nil, err
		}
	}
	if plugin != string(mysql.MysqlNativePassword) && plugin != string(mysql.CachingSha2Password) {
		if err := mysqlDb.VerifyPlugin(plugin); err != nil {
			return nil, sql.ErrUserAlterFailure.New(err)
		}
	}

	previousUserEntry.Plugin = plugin
	previousUserEntry.AuthString = authString
	previousUserEntry.PasswordLastChanged = time.Now().UTC()
	editor.PutUser(previousUserEntry)

	if err := mysqlDb.Persist(ctx, editor); err != nil {
		return nil, err
	}

	return sql.RowsToRowIter(sql.Row{types.NewOkResult(0)}), nil
}

func (b *BaseBuilder) buildCreateUserImpl(ctx *sql.Context, n *plan.CreateUser) (sql.RowIter, error) {
	mysqlDb, ok := n.MySQLDb.(*mysql_db.MySQLDb)
	if !ok {
		return nil, sql.ErrDatabaseNotFound.New("mysql")
	}

	editor := mysqlDb.Editor()
	defer editor.Close()

	for _, user := range n.Users {
		if user.UserName.Host == "" {
			user.UserName.Host = "%"
		}

		userPk := mysql_db.UserPrimaryKey{
			Host: user.UserName.Host,
			User: user.UserName.Name,
		}
		_, ok := editor.GetUser(userPk)
		if ok {
			if n.IfNotExists {
				continue
			}
			return nil, sql.ErrUserCreationFailure.New(user.UserName.String("'"))
		}

		if len(user.UserName.Name) > 32 {
			return nil, sql.ErrUserNameTooLong.New(user.UserName.Name)
		}

		if len(user.UserName.Host) > 255 {
			return nil, sql.ErrUserHostTooLong.New(user.UserName.Host)
		}

		plugin := string(mysql_db.DefaultAuthMethod)
		authString := ""
		if user.Auth1 != nil {
			plugin = user.Auth1.Plugin()
			var err error
			authString, err = user.Auth1.AuthString()
			if err != nil {
				return nil, err
			}
		}
		if plugin != string(mysql.MysqlNativePassword) && plugin != string(mysql.CachingSha2Password) {
			if err := mysqlDb.VerifyPlugin(plugin); err != nil {
				return nil, sql.ErrUserCreationFailure.New(err)
			}
		}

		sslType, sslCipher, x509Issuer, x509Subject := parseTlsOptions(n.TLSOptions)
		editor.PutUser(&mysql_db.User{
			User:                user.UserName.Name,
			Host:                user.UserName.Host,
			PrivilegeSet:        mysql_db.NewPrivilegeSet(),
			Plugin:              plugin,
			AuthString:          authString,
			PasswordLastChanged: time.Now().UTC(),
			Locked:              false,
			Attributes:          nil,
			IsRole:              false,
			Identity:            user.Identity,
			SslType:             sslType,
			X509Issuer:          x509Issuer,
			X509Subject:         x509Subject,
			SslCipher:           sslCipher,
		})
	}
	if err := mysqlDb.Persist(ctx, editor); err != nil {
		return nil, err
	}
	return rowIterWithOkResultWithZeroRowsAffected(), nil
}
