package trace

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"runtime"
)

func InitJaegerTracer(serviceName string, addr string, samplerType string, samplerParam float64) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  samplerType,
			Param: samplerParam,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			//BufferFlushInterval: 1 * time.Second,
		},
	}

	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer, err
}

func NewGinContextAndSpan(c *gin.Context) (context.Context, opentracing.Span) {
	var sp opentracing.Span
	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request.Header),
	)
	if err != nil {
		// If for whatever reason we can't join, go ahead an start a new root span.
		sp = opentracing.StartSpan(c.Request.RequestURI)
	} else {
		sp = opentracing.StartSpan(c.Request.RequestURI, opentracing.ChildOf(wireContext))
	}
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, sp)
	return ctx, sp
}

//内部函数链路追踪
func NewInnerSpan(name string, ctx context.Context) opentracing.Span {
	parent := opentracing.SpanFromContext(ctx)
	if parent == nil {
		return nil
	}
	tracer := opentracing.GlobalTracer()
	if tracer == nil {
		return nil
	}
	return tracer.StartSpan(name, opentracing.ChildOf(parent.Context()))
}

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func AddSpanError(span opentracing.Span, err error) {
	if span == nil {
		return
	}
	span.SetTag(string(ext.Error), true)
	span.LogKV(otlog.Error(err))
}
