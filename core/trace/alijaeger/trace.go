package alijaeger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
	"os"
)

const (
	EndpointKey = "Jaeger"
)

type AliJaeger struct {
	closer io.Closer
	tracer opentracing.Tracer
}

func NewAliJaeger(service, endpoint string) *AliJaeger {
	envEndpoint := os.Getenv(EndpointKey)
	if len(envEndpoint) > 0 {
		endpoint = envEndpoint
	}

	if len(endpoint) < 1 {
		return &AliJaeger{}
	}

	sender := transport.NewHTTPTransport(endpoint)
	tracer, closer := jaeger.NewTracer(service,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender))

	opentracing.SetGlobalTracer(tracer)

	return &AliJaeger{
		tracer: tracer,
		closer: closer,
	}
}

func (aj *AliJaeger) SafeStop() error {
	return aj.closer.Close()
}

func (aj *AliJaeger) Trace(ctx context.Context, name string, fn func(ctx context.Context, span opentracing.Span) error) error {
	span, c := opentracing.StartSpanFromContext(ctx, name)
	defer span.Finish()
	return fn(c, span)
}

func (aj *AliJaeger) TraceMysql(ctx context.Context, sql, name string, fn func(ctx context.Context, span opentracing.Span) error) error {
	span, c := opentracing.StartSpanFromContext(ctx, sql)
	ext.DBStatement.Set(span, sql)
	ext.DBType.Set(span, "mysql")
	defer span.Finish()
	return fn(c, span)
}

func (aj *AliJaeger) TraceRedis(ctx context.Context, name string, fn func(ctx context.Context, span opentracing.Span) error) error {
	span, c := opentracing.StartSpanFromContext(ctx, name)
	defer span.Finish()
	return fn(c, span)
}

// AliTracingHandler http中间件
func AliTracingHandler(next http.Handler) http.Handler {
	tracer := opentracing.GlobalTracer()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var span opentracing.Span
		ctx, err := tracer.Extract(opentracing.HTTPHeaders, r.Header)
		if err != nil {
			span = opentracing.StartSpan(r.RequestURI)
		} else {
			span = opentracing.StartSpan(r.RequestURI, opentracing.ChildOf(ctx))
		}

		defer span.Finish()
		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, r.RequestURI)

		//span.Tracer().Inject(ctx, opentracing.HTTPHeaders, r.Header)
		next.ServeHTTP(w, r.WithContext(opentracing.ContextWithSpan(r.Context(), span)))
	})
}

// AliUnaryTracingInterceptor rpc中间件
func AliUnaryTracingInterceptor() grpc.UnaryServerInterceptor {
	tracer := opentracing.GlobalTracer()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		sc, err := tracer.Extract(opentracing.HTTPHeaders, metadataReaderWriter{md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			resp, err = handler(ctx, req)
		} else {
			span := tracer.StartSpan(info.FullMethod, ext.RPCServerOption(sc), gRPCComponentTag)
			defer span.Finish()

			resp, err = handler(opentracing.ContextWithSpan(ctx, span), req)
			if err != nil {
				span.SetTag("error", err)
			}
		}
		return
	}
}
