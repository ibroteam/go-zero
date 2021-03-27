package alijaeger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"testing"
	"time"
)

const endpoint = "http://tracing-analysis-dc-sz.aliyuncs.com/adapt_jcbj57tc59@49deac3fa650936_jcbj57tc59@53df7ad2afe8301/api/traces"

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
