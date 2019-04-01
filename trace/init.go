package trace

import (
	"io"
	"time"

	gCfg "github.com/g-airport/go-infra/config"

	jCfg "github.com/uber/jaeger-client-go/config"
	jLog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

// this method will initial  global tracer
// if no call global tracer default  no-op tracer
// init config from key:value
func Init(serviceName string) (io.Closer, error) {
	cfg := &jCfg.Configuration{}
	if err := gCfg.Get("trace").Scan(cfg); err != nil {
		return nil, err
	}
	return InitByConfig(serviceName, *cfg)
}

func InitByConfig(serviceName string, cfg jCfg.Configuration) (io.Closer, error) {
	cfg.Sampler.SamplingRefreshInterval = 5 * time.Second
	cfg.Reporter.BufferFlushInterval = 5 * time.Second

	// TODO logging & metric framework
	jLogger := jLog.StdLogger
	jMetricsFactory := metrics.NullFactory
	return cfg.InitGlobalTracer(serviceName, jCfg.Logger(jLogger), jCfg.Metrics(jMetricsFactory))
}
