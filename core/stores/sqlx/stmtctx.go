// +build !prod

package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/tal-tech/go-zero/core/trace/alijaeger"
)

func execWithCtx(ctx context.Context, conn sessionConn, q string, args ...interface{}) (sql.Result, error) {
	stmt, err := format(q, args...)
	span, _ := opentracing.StartSpanFromContext(ctx, "db")
	defer span.Finish()
	ext.DBStatement.Set(span, q)
	span.SetTag(alijaeger.TagStmtArgs, stmt)

	result, err := conn.Exec(q, args...)
	if err != nil {
		logSqlError(stmt, err)
	}

	return result, err
}

func execStmtWithCtx(ctx context.Context, conn stmtConn, args ...interface{}) (sql.Result, error) {
	stmt := fmt.Sprint(args...)
	span, _ := opentracing.StartSpanFromContext(ctx, "db")
	defer span.Finish()
	ext.DBStatement.Set(span, stmt)

	result, err := conn.Exec(args...)
	if err != nil {
		logSqlError(stmt, err)
	}
	return result, err
}

func queryWithCtx(ctx context.Context, conn sessionConn, scanner func(*sql.Rows) error, q string, args ...interface{}) error {
	stmt := fmt.Sprint(args...)
	span, _ := opentracing.StartSpanFromContext(ctx, "db")
	defer span.Finish()
	ext.DBStatement.Set(span, q)
	span.SetTag(alijaeger.TagStmtArgs, stmt)

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
	span, _ := opentracing.StartSpanFromContext(ctx, "db")
	defer span.Finish()
	ext.DBStatement.Set(span, stmt)

	rows, err := conn.Query(args...)
	if err != nil {
		logSqlError(stmt, err)
		return err
	}
	defer rows.Close()

	return scanner(rows)
}
