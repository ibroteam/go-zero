package sqlx

import (
	"context"
	"database/sql"
	"fmt"
)

type (
	beginnableCtx func(*sql.DB) (transCtx, error)

	transCtx interface {
		SessionCtx
		Commit() error
		Rollback() error
	}

	txSessionCtx struct {
		*sql.Tx
	}
)

func (t txSessionCtx) Exec(ctx context.Context, q string, args ...interface{}) (sql.Result, error) {
	return execWithCtx(ctx, t.Tx, q, args...)
}

func (t txSessionCtx) Prepare(q string) (StmtSessionCtx, error) {
	stmt, err := t.Tx.Prepare(q)
	if err != nil {
		return nil, err
	}

	return statementCtx{
		stmt: stmt,
	}, nil
}

func (t txSessionCtx) QueryRow(ctx context.Context, v interface{}, q string, args ...interface{}) error {
	return queryWithCtx(ctx, t.Tx, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, true)
	}, q, args...)
}

func (t txSessionCtx) QueryRowPartial(ctx context.Context, v interface{}, q string, args ...interface{}) error {
	return queryWithCtx(ctx, t.Tx, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, false)
	}, q, args...)
}

func (t txSessionCtx) QueryRows(ctx context.Context, v interface{}, q string, args ...interface{}) error {
	return queryWithCtx(ctx, t.Tx, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, true)
	}, q, args...)
}

func (t txSessionCtx) QueryRowsPartial(ctx context.Context, v interface{}, q string, args ...interface{}) error {
	return queryWithCtx(ctx, t.Tx, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, q, args...)
}

func beginCtx(db *sql.DB) (transCtx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	return txSessionCtx{
		Tx: tx,
	}, nil
}

func transactCtx(ctx context.Context, db *commonSqlConnCtx, b beginnableCtx, fn func(ctx context.Context, sess SessionCtx) error) (err error) {
	conn, err := getSqlConn(db.driverName, db.datasource)
	if err != nil {
		logInstanceError(db.datasource, err)
		return err
	}

	return transactOnConnCtx(ctx, conn, b, fn)
}

func transactOnConnCtx(ctx context.Context, conn *sql.DB, b beginnableCtx, fn func(ctx context.Context, sess SessionCtx) error) (err error) {
	var tx transCtx
	tx, err = b(conn)
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("recover from %#v, rollback failed: %s", p, e)
			} else {
				err = fmt.Errorf("recoveer from %#v", p)
			}
		} else if err != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("transaction failed: %s, rollback failed: %s", err, e)
			}
		} else {
			err = tx.Commit()
		}
	}()

	return fn(ctx, tx)
}
