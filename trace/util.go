package trace

import (
	"context"
	"github.com/micro/go-micro/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"
)

//---------------------------------------------
// gPRC via context share span

func InjectSpanToContext(ctx context.Context, span opentracing.Span) (context.Context, error) {
	return InjectSpanContextToContext(ctx, span.Tracer(), span.Context())
}

func InjectSpanContextToContext(ctx context.Context, tracer opentracing.Tracer, spanContext opentracing.SpanContext) (context.Context, error) {
	md := make(metadata.Metadata)
	if err := tracer.Inject(spanContext, opentracing.TextMap, opentracing.TextMapCarrier(md)); err != nil {
		return ctx, err
	}
	nCtx := metadata.NewContext(ctx, md)
	return nCtx, nil
}

func ExtractSpanFromContext(ctx context.Context, tracer opentracing.Tracer, operationName string) opentracing.Span {
	wireContext, err := ExtractSpanContextFromContext(ctx, tracer)
	var sp opentracing.Span
	if err != nil {
		sp = tracer.StartSpan(operationName)
	} else {
		sp = tracer.StartSpan(operationName, opentracing.ChildOf(wireContext))
	}
	return sp
}

func ExtractSpanContextFromContext(ctx context.Context, tracer opentracing.Tracer) (opentracing.SpanContext, error) {
	md, _ := metadata.FromContext(ctx)
	return tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
}

//---------------------------------------------
// mq via header share span
// bridge mq.Table & TextMapReader,TextMapWriter (rabbit mq)

type MQTextMapCarrier amqp.Table

func (c MQTextMapCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, v := range c {
		if sv, ok := v.(string); ok {
			if err := handler(k, sv); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c MQTextMapCarrier) Set(key, val string) {
	c[key] = val
}

func InjectSpanToMQHeader(header amqp.Table, span opentracing.Span) error {
	return InjectSpanContextToMQHeader(header, span.Tracer(), span.Context())
}

func InjectSpanContextToMQHeader(header amqp.Table, tracer opentracing.Tracer, spanContext opentracing.SpanContext) error {
	c := MQTextMapCarrier(header)
	if err := tracer.Inject(spanContext, opentracing.TextMap, c); err != nil {
		return err
	}
	return nil
}

func ExtractSpanFromMQHeader(header amqp.Table, tracer opentracing.Tracer, name string) (sp opentracing.Span) {
	wireContext, err := ExtractSpanContextFromMQHeader(header, tracer)
	if err != nil {
		sp = tracer.StartSpan(name)
	} else {
		sp = tracer.StartSpan(name, opentracing.ChildOf(wireContext))
	}
	return
}

func ExtractSpanContextFromMQHeader(header amqp.Table, tracer opentracing.Tracer) (opentracing.SpanContext, error) {
	c := MQTextMapCarrier(header)
	return tracer.Extract(opentracing.TextMap, c)
}
