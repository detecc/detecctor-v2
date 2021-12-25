package tracing

import (
	"context"
	"github.com/detecc/detecctor-v2/internal/model/configuration"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const (
	NoAuth = ""
	TLS    = "tls"
)

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func connectToBackend(ctx context.Context, address, authType string) (*grpc.ClientConn, error) {
	log.Debug("Connecting to OpenTelemetry backend with GRPC..")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())

	switch authType {
	case NoAuth:
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		break
	case TLS:
		break
	}

	dialContext, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	return grpc.DialContext(dialContext, address, opts...)
}

func createNewExporter(ctx context.Context, config configuration.OTel) (*otlptrace.Exporter, error) {
	conn, err := connectToBackend(ctx, config.Address, config.AuthType)
	handleErr(err, "failed to create gRPC connection to collector")

	// Set up a trace exporter
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

func InitTelemetry(res *resource.Resource, config configuration.OTel) func() {
	log.Info("Initializing OpenTelemetry")
	ctx := context.Background()

	// Create trace exporter
	exporter, err := createNewExporter(ctx, config)
	handleErr(err, "failed to create trace exporter")

	// Create tracer provider
	bsp := sdktrace.NewBatchSpanProcessor(
		exporter,
		sdktrace.WithBatchTimeout(time.Second*5),
		sdktrace.WithMaxQueueSize(1000),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// Set the global provider
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		handleErr(tracerProvider.Shutdown(ctx), "failed to shutdown TracerProvider")
	}
}
