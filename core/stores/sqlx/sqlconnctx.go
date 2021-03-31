package sqlx

import (
	"context"
	"database/sql"

	"github.com/tal-tech/go-zero/core/breaker"
)

type (
	// Session stands for raw connections or transaction sessions
	SessionCtx interface {
		Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		Prepare(query string) (StmtSessionCtx, error)
		QueryRow(ctx context.Context, v interface{}, query string, args ...interface{}) error
		QueryRowPartial(ctx context.Context, v interface{}, query string, args ...interface{}) error
		QueryRows(ctx context.Context, v interface{}, query string, args ...interface{}) error
		QueryRowsPartial(ctx context.Context, v interface{}, query string, args ...interface{}) error
	}

	// SqlConn only stands for raw connections, so Transact method can be called.
	SqlConnCtx interface {
		SessionCtx
		Transact(context.Context, func(ctx context.Context, session SessionCtx) error) error
	}

	// SqlOption defines the method to customize a sql connection.
	SqlOptionCtx func(*commonSqlConnCtx)

	// StmtSession interface represents a session that can be used to execute statements.
	StmtSessionCtx interface {
		Close() error
		Exec(ctx context.Context, args ...interface{}) (sql.Result, error)
		QueryRow(ctx context.Context, v interface{}, args ...interface{}) error
		QueryRowPartial(ctx context.Context, v interface{}, args ...interface{}) error
		QueryRows(ctx context.Context, v interface{}, args ...interface{}) error
		QueryRowsPartial(ctx context.Context, v interface{}, args ...interface{}) error
	}

	// thread-safe
	// Because CORBA doesn't support PREPARE, so we need to combine the
	// query arguments into one string and do underlying query without arguments
	commonSqlConnCtx struct {
		driverName string
		datasource string
		beginTx    beginnableCtx
		brk        breaker.Breaker
		accept     func(error) bool
	}

	statementCtx struct {
		stmt *sql.Stmt
	}
)

// NewSqlConn returns a SqlConn with given driver name and datasource.
func NewSqlConnCtx(driverName, datasource string, opts ...SqlOptionCtx) SqlConnCtx {
	conn := &commonSqlConnCtx{
		driverName: driverName,
		datasource: datasource,
		beginTx:    beginCtx,
		brk:        breaker.NewBreaker(),
	}
	for _, opt := range opts {
		opt(conn)
	}

	return conn
}

func (db *commonSqlConnCtx) Exec(ctx context.Context, q string, args ...interface{}) (result sql.Result, err error) {
	err = db.brk.DoWithAcceptable(func() error {
		var conn *sql.DB
		conn, err = getSqlConn(db.driverName, db.datasource)
		if err != nil {
			logInstanceError(db.datasource, err)
			return err
		}

		result, err = execWithCtx(ctx, conn, q, args...)
		return err
	}, db.acceptable)

	return
}

func (db *commonSqlConnCtx) Prepare(query string) (stmt StmtSessionCtx, err error) {
	err = db.brk.DoWithAcceptable(func() error {
		var conn *sql.DB
		conn, err = getSqlConn(db.driverName, db.datasource)
		if err != nil {
			logInstanceError(db.datasource, err)
			return err
		}

		st, err := conn.Prepare(query)
		if err != nil {
			return err
		}

		stmt = statementCtx{
			stmt: st,
		}
		return nil
	}, db.acceptable)

	return
}

func (db *commonSqlConnCtx) QueryRow(ctx context.Context, v interface{}, q string, args ...interface{}) error {
	return db.queryRows(ctx, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, true)
	}, q, args...)
}

func (db *commonSqlConnCtx) QueryRowPartial(ctx context.Context, v interface{}, q string, args ...interface{}) error {
	return db.queryRows(ctx, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, false)
	}, q, args...)
}

func (db *commonSqlConnCtx) QueryRows(ctx context.Context, v interface{}, q string, args ...interface{}) error {
	return db.queryRows(ctx, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, true)
	}, q, args...)
}

func (db *commonSqlConnCtx) QueryRowsPartial(ctx context.Context, v interface{}, q string, args ...interface{}) error {
	return db.queryRows(ctx, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, q, args...)
}

func (db *commonSqlConnCtx) Transact(ctx context.Context, fn func(ctx context.Context, sessionCtx SessionCtx) error) error {
	return db.brk.DoWithAcceptable(func() error {
		return transactCtx(ctx, db, db.beginTx, fn)
	}, db.acceptable)
}

func (db *commonSqlConnCtx) acceptable(err error) bool {
	ok := err == nil || err == sql.ErrNoRows || err == sql.ErrTxDone
	if db.accept == nil {
		return ok
	}

	return ok || db.accept(err)
}

func (db *commonSqlConnCtx) queryRows(ctx context.Context, scanner func(*sql.Rows) error, q string, args ...interface{}) error {
	var qErr error
	return db.brk.DoWithAcceptable(func() error {
		conn, err := getSqlConn(db.driverName, db.datasource)
		if err != nil {
			logInstanceError(db.datasource, err)
			return err
		}

		return queryWithCtx(ctx, conn, func(rows *sql.Rows) error {
			qErr = scanner(rows)
			return qErr
		}, q, args...)
	}, func(err error) bool {
		return qErr == err || db.acceptable(err)
	})
}

func (s statementCtx) Close() error {
	return s.stmt.Close()
}

func (s statementCtx) Exec(ctx context.Context, args ...interface{}) (sql.Result, error) {
	return execStmtWithCtx(ctx, s.stmt, args...)
}

func (s statementCtx) QueryRow(ctx context.Context, v interface{}, args ...interface{}) error {
	return queryStmtWithCtx(ctx, s.stmt, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, true)
	}, args...)
}

func (s statementCtx) QueryRowPartial(ctx context.Context, v interface{}, args ...interface{}) error {
	return queryStmtWithCtx(ctx, s.stmt, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, false)
	}, args...)
}

func (s statementCtx) QueryRows(ctx context.Context, v interface{}, args ...interface{}) error {
	return queryStmtWithCtx(ctx, s.stmt, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, true)
	}, args...)
}

func (s statementCtx) QueryRowsPartial(ctx context.Context, v interface{}, args ...interface{}) error {
	return queryStmtWithCtx(ctx, s.stmt, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, args...)
}
