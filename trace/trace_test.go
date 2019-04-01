package trace

import (
	"encoding/json"
	"fmt"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  "192.168.50.171:5775",
			//LocalAgentHostPort: localAgentHostPort,
		},
	}

	bytes, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bytes))

}
