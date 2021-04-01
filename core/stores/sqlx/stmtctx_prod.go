// +build prod

package sqlx

import (
	"context"
	"database/sql"
	"fmt"
)

func execWithCtx(ctx context.Context, conn sessionConn, q string, args ...interface{}) (sql.Result, error) {
	stmt := formatForPrint(q, args)

	result, err := conn.Exec(q, args...)
	if err != nil {
		logSqlError(stmt, err)
	}

	return result, err
}

func execStmtWithCtx(ctx context.Context, conn stmtConn, args ...interface{}) (sql.Result, error) {
	stmt := fmt.Sprint(args...)

	result, err := conn.Exec(args...)
	if err != nil {
		logSqlError(stmt, err)
	}
	return result, err
}

func queryWithCtx(ctx context.Context, conn sessionConn, scanner func(*sql.Rows) error, q string, args ...interface{}) error {
	stmt := fmt.Sprint(args...)

	rows, err := conn.Query(q, args...)
	if err != nil {
		logSqlError(stmt, err)
		return err
	}
	defer rows.Close()
	return scanner(rows)
}

func queryStmtWithCtx(ctx context.Context, conn stmtConn, scanner func(*sql.Rows) error, args ...interface{}) error {
	stmt := fmt.Sprint(args...)

	rows, err := conn.Query(args...)
	if err != nil {
		logSqlError(stmt, err)
		return err
	}
	defer rows.Close()

	return scanner(rows)
}
