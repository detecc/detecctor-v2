package metrics

import (
	"github.com/detecc/detecctor-v2/internal/model/configuration"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	cfg = prometheus.Config{}
)

type Prometheus struct {
	controller *controller.Controller
	cfg        prometheus.Config
	exporter   *prometheus.Exporter
	http       *gin.Engine
}

// NewPrometheus creates a new HTTP server with metrics endpoint
func NewPrometheus(controller *controller.Controller, config prometheus.Config, exporter *prometheus.Exporter, endpoint string) *Prometheus {
	// Configure metrics endpoint
	r := gin.New()
	r.GET(endpoint, prometheusHandler(exporter))

	return &Prometheus{
		controller: controller,
		cfg:        config,
		exporter:   exporter,
		http:       r,
	}
}

// Start Run the HTTP server with metrics
func (p *Prometheus) Start(address string) {
	err := p.http.Run(address)
	if err != nil {
		log.WithError(err).Errorf("Error exposing prometheus")
	}
}

func prometheusHandler(exporter *prometheus.Exporter) gin.HandlerFunc {
	return func(c *gin.Context) {
		exporter.ServeHTTP(c.Writer, c.Request)
	}
}

// FromConfiguration sets up and exposes prometheus metrics.
func FromConfiguration(res *resource.Resource, config configuration.Metrics) *Prometheus {
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(cfg.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithResource(res),
	)

	exporter, err := prometheus.New(cfg, c)
	if err != nil {
		log.WithError(err).Fatalf("Failed to initialize prometheus exporter")
	}

	global.SetMeterProvider(exporter.MeterProvider())

	return NewPrometheus(c, cfg, exporter, config.Endpoint)
}
