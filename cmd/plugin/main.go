package main

import (
	"context"
	"fmt"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/internal/config"
	"github.com/detecc/detecctor-v2/pkg/mqtt"
	"github.com/detecc/detecctor-v2/pkg/observability"
	"github.com/detecc/detecctor-v2/pkg/observability/logging"
	"github.com/detecc/detecctor-v2/pkg/observability/metrics"
	"github.com/detecc/detecctor-v2/pkg/observability/tracing"
	pluginService "github.com/detecc/detecctor-v2/service/plugin"
	"github.com/detecc/detecctor-v2/service/plugin/plugin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"os"
	"os/signal"
)

const (
	debugFlag   = "debug"
	tracingFlag = "tracing"
	metricsFlag = "metrics"

	logPath = "/var/log/detecctor/plugin.log"
)

var (
	configurationFilePathFlag string
	databaseType              string

	isDebug       = false
	enableMetrics = false
	enableTracing = false

	rootCmd = &cobra.Command{
		Use:   "detecctor-plugin",
		Short: "Plugin service",
		Long:  ``,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		quitChannel = make(chan os.Signal)
		res, err    = resource.New(ctx,
			resource.WithAttributes(
				semconv.ServiceNameKey.String("plugin-service"),
				semconv.ServiceVersionKey.String("v0.1.0"),
				observability.DatabaseTypeLabel(databaseType),
			),
			resource.WithOS(),
			resource.WithOSType(),
			resource.WithHost(),
		)
		prom *metrics.Prometheus
		//tracer      *tracing.Tracer
	)
	// Create exit handlers
	signal.Notify(quitChannel, os.Interrupt)

	log.Info("Starting the plugin service..")

	// Get the service configuration
	pluginConfig := config.GetPluginServiceConfiguration(configurationFilePathFlag)

	// Setup logging
	logging.Setup(log.StandardLogger(), pluginConfig.Observability.Logging, logPath, isDebug)

	if enableTracing && err == nil {
		f := tracing.InitTelemetry(res, pluginConfig.Observability.OTel)
		defer f()
		//tracer = tracing.NewTracer(log.StandardLogger(), nil)
	}

	if enableMetrics && err == nil {
		prom = metrics.FromConfiguration(res, pluginConfig.Observability.Metrics)
		go prom.Start(pluginConfig.Observability.Metrics.Address)
	}

	// Connect to the Database
	database.Connect(database.Database(databaseType), pluginConfig.Database)

	// Load plugins into the manager
	pluginManager := plugin.GetPluginManager()
	pluginManager.LoadPlugins(pluginConfig.PluginConfiguration)

	// Create the service
	service := pluginService.NewPluginService(pluginManager, database.GetLogRepository())

	// Start listening to the MQTT client
	mqttClient := mqtt.NewMqttClient(pluginConfig.MqttBroker)
	mqttClient.Subscribe(pluginService.ExecutionTopic, service.ExecutionHandler())
	mqttClient.Subscribe(pluginService.ResponseTopic, service.ResponseHandler())

Loop:
	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping the cmd service..")
			mqttClient.Disconnect()
			cancel()
			break Loop
		case <-quitChannel:
			cancel()
			break
		}
	}
}

func main() {
	var (
		workingDirectory, _   = os.Getwd()
		defaultConfigFileName = fmt.Sprintf("%s/plugin-service.%s", workingDirectory, "yaml")
	)

	// Set flags
	rootCmd.PersistentFlags().StringVar(&configurationFilePathFlag, "config", defaultConfigFileName, "config file path")
	rootCmd.PersistentFlags().StringVar(&databaseType, "database", "Mongo", "database type (supported: Mongo)")
	rootCmd.PersistentFlags().BoolVarP(&isDebug, debugFlag, "d", false, "debug mode")
	rootCmd.PersistentFlags().BoolVarP(&enableMetrics, metricsFlag, "m", false, "enable metrics for Prometheus")
	rootCmd.PersistentFlags().BoolVarP(&enableTracing, tracingFlag, "t", false, "enable tracing using OpenTelemetry")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Cannot run the service")
	}
}
