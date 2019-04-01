package trace

import (
	"context"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/server"
	"github.com/opentracing/opentracing-go"
)

//Usage:
//service := micro.NewService()
//
//service.Init()
//service.Init(
//    micro.WrapCall(ClientWrapper())
//)

func ServerWrapper() server.HandlerWrapper {
	return func(f server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			span := ExtractSpanFromContext(ctx, opentracing.GlobalTracer(), "server."+req.Method())
			span.SetTag("service", req.Service())
			span.SetTag("method", req.Method())
			defer span.Finish()
			return f(opentracing.ContextWithSpan(ctx, span), req, rsp)
		}
	}
}

func ClientWrapper() client.CallWrapper {
	return func(f client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			span, ctx := opentracing.StartSpanFromContext(ctx, "client."+req.Method())
			span.SetTag("address", node.Address)
			span.SetTag("service", req.Service())
			span.SetTag("method", req.Method())
			defer span.Finish()
			ctx, _ = InjectSpanToContext(ctx, span)
			return f(ctx, &registry.Node{Address: node.Address}, req, rsp, opts)
		}
	}
}
