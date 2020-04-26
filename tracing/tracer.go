package tracing

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/whatisfaker/zaptrace/log"
)

func NewTracer(serviceName string, logger *log.Factory) (opentracing.Tracer, io.Closer, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, nil, err
	}
	// cfg.Sampler.Type = "const"
	// cfg.Sampler.Param = 1
	cfg.ServiceName = serviceName
	jaegerLogger := jaegerLoggerImp{logger.Normal()}
	return cfg.NewTracer(
		config.Logger(jaegerLogger),
	)
}
