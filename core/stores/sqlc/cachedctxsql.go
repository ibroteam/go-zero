package sqlc

import (
	"context"
	"database/sql"
	"time"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type (
	// ExecCtxFn defines the sql exec method.
	ExecCtxFn func(ctx context.Context, conn sqlx.SqlConnCtx) (sql.Result, error)

	// IndexQueryCtxFn defines the query method that based on unique indexes.
	IndexQueryCtxFn func(ctx context.Context, conn sqlx.SqlConnCtx, v interface{}) (interface{}, error)

	// PrimaryQueryCtxFn defines the query method that based on primary keys.
	PrimaryQueryCtxFn func(ctx context.Context, conn sqlx.SqlConnCtx, v, primary interface{}) error

	// QueryCtxFn defines the query method.
	QueryCtxFn func(ctx context.Context, conn sqlx.SqlConnCtx, v interface{}) error

	// A CachedConnCtx is a DB connection with cache capability.
	CachedConnCtx struct {
		db    sqlx.SqlConnCtx
		cache cache.Cache
	}
)

// NewNodeConn returns a CachedConn with a redis node cache.
func NewNodeConnCtx(db sqlx.SqlConnCtx, rds *redis.Redis, opts ...cache.Option) CachedConnCtx {
	return CachedConnCtx{
		db:    db,
		cache: cache.NewNode(rds, exclusiveCalls, stats, sql.ErrNoRows, opts...),
	}
}

// NewConn returns a CachedConn with a redis cluster cache.
func NewConnCtx(db sqlx.SqlConnCtx, c cache.CacheConf, opts ...cache.Option) CachedConnCtx {
	return CachedConnCtx{
		db:    db,
		cache: cache.New(c, exclusiveCalls, stats, sql.ErrNoRows, opts...),
	}
}

// DelCache deletes cache with keys.
func (cc CachedConnCtx) DelCache(keys ...string) error {
	return cc.cache.Del(keys...)
}

// GetCache unmarshals cache with given key into v.
func (cc CachedConnCtx) GetCache(key string, v interface{}) error {
	return cc.cache.Get(key, v)
}

// Exec runs given exec on given keys, and returns execution result.
func (cc CachedConnCtx) Exec(ctx context.Context, exec ExecCtxFn, keys ...string) (sql.Result, error) {
	res, err := exec(ctx, cc.db)
	if err != nil {
		return nil, err
	}

	if err := cc.DelCache(keys...); err != nil {
		return nil, err
	}

	return res, nil
}

// ExecNoCache runs exec with given sql statement, without affecting cache.
func (cc CachedConnCtx) ExecNoCache(ctx context.Context, q string, args ...interface{}) (sql.Result, error) {
	return cc.db.Exec(ctx, q, args...)
}

// QueryRow unmarshals into v with given key and query func.
func (cc CachedConnCtx) QueryRow(ctx context.Context, v interface{}, key string, query QueryCtxFn) error {
	return cc.cache.Take(v, key, func(v interface{}) error {
		return query(ctx, cc.db, v)
	})
}

// QueryRowIndex unmarshals into v with given key.
func (cc CachedConnCtx) QueryRowIndex(ctx context.Context, v interface{}, key string, keyer func(primary interface{}) string,
	indexQuery IndexQueryCtxFn, primaryQuery PrimaryQueryCtxFn) error {
	var primaryKey interface{}
	var found bool

	if err := cc.cache.TakeWithExpire(&primaryKey, key, func(val interface{}, expire time.Duration) (err error) {
		primaryKey, err = indexQuery(ctx, cc.db, v)
		if err != nil {
			return
		}

		found = true
		return cc.cache.SetWithExpire(keyer(primaryKey), v, expire+cacheSafeGapBetweenIndexAndPrimary)
	}); err != nil {
		return err
	}

	if found {
		return nil
	}

	return cc.cache.Take(v, keyer(primaryKey), func(v interface{}) error {
		return primaryQuery(ctx, cc.db, v, primaryKey)
	})
}

// QueryRowNoCache unmarshals into v with given statement.
func (cc CachedConnCtx) QueryRowNoCache(ctx context.Context, v interface{}, q string, args ...interface{}) error {
	return cc.db.QueryRow(ctx, v, q, args...)
}

// QueryRowsNoCache unmarshals into v with given statement.
// It doesn't use cache, because it might cause consistency problem.
func (cc CachedConnCtx) QueryRowsNoCache(ctx context.Context, v interface{}, q string, args ...interface{}) error {
	return cc.db.QueryRows(ctx, v, q, args...)
}

// SetCache sets v into cache with given key.
func (cc CachedConnCtx) SetCache(key string, v interface{}) error {
	return cc.cache.Set(key, v)
}

// Transact runs given fn in transaction mode.
func (cc CachedConnCtx) Transact(ctx context.Context, fn func(context.Context, sqlx.SessionCtx) error) error {
	return cc.db.Transact(ctx, fn)
}
