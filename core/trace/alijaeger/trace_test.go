package alijaeger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"testing"
	"time"
)

const endpoint = "<endpoint>"

func TestAliJaeger(t *testing.T) {
	aj := NewAliJaeger("test", endpoint)
	aj.Trace(context.Background(), "test", func(ctx context.Context, span opentracing.Span) error {
		sp, _ := opentracing.StartSpanFromContext(ctx, "db")
		sp.SetTag("db.statement", "select * from broadcast where 1 limit 1")
		defer sp.Finish()
		time.Sleep(time.Second * 2)
		return nil
	})
	aj.SafeStop()
}
