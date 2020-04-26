package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// QuickStartSpan 通过上下文快速创建Span
func QuickStartSpan(ctx context.Context, operationName string, normaltag opentracing.Tag, tags ...map[string]string) (context.Context, opentracing.Span) {
	return QuickStartSpanWithTracer(ctx, opentracing.GlobalTracer(), operationName, normaltag, tags...)
}

// QuickStartSpanWithTracer 通过上下文快速创建Span
func QuickStartSpanWithTracer(ctx context.Context, tr opentracing.Tracer, operationName string, normaltag opentracing.Tag, tags ...map[string]string) (context.Context, opentracing.Span) {
	var span opentracing.Span
	if sp := opentracing.SpanFromContext(ctx); sp != nil {
		span = tr.StartSpan(operationName, opentracing.ChildOf(sp.Context()))
	} else {
		span = opentracing.StartSpan(operationName)
	}
	if len(tags) > 0 {
		for k, v := range tags[0] {
			span.SetTag(k, v)
		}
	}
	normaltag.Set(span)
	ctx = opentracing.ContextWithSpan(ctx, span)
	return ctx, span
}
