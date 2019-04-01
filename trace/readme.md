## Embed Trace

#### config


```go
import "github.com/uber/jaeger-client-go/config"

cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:                    jaeger.SamplerTypeConst,
			Param:                   1,
			SamplingRefreshInterval: time.Second,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			LocalAgentHostPort:  "127.0.0.1:5775",
			BufferFlushInterval: time.Second,
		},
	}
    
```

#### init

```go
    InitByConfig(serviceName, cfg)
```


#### trace

```go
    span, ctx := opentracing.StartSpanFromContext(context.Background(), "name")
    span.Finish()
```