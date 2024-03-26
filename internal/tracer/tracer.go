package tracer

import (
	"io"
	"log/slog"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	"github.com/arslanovdi/logistic-package-api/internal/config"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// NewTracer - returns new tracer.
func NewTracer(cfg *config.Config) (io.Closer, error) {

	log := slog.With("func", "tracer.NewTracer")

	cfgTracer := &jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.Service,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: cfg.Jaeger.Host + cfg.Jaeger.Port,
		},
	}
	tracer, closer, err := cfgTracer.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		log.Error("failed init jaeger", slog.String("error", err.Error()))

		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	log.Info("Traces started")

	return closer, nil
}
