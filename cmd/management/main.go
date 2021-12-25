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
	"github.com/detecc/detecctor-v2/service/management"
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

	logPath = "/var/log/detecctor-v2/management.log"
)

var (
	configurationFilePathFlag string
	databaseType              string

	isDebug       = false
	enableMetrics = false
	enableTracing = false

	rootCmd = &cobra.Command{
		Use:   "detecctor-mgmt",
		Short: "Management service",
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
				semconv.ServiceNameKey.String("management-service"),
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

	log.Info("Starting the cmd service...")

	// Get configuration for cmd service
	managementConfig := config.GetManagementServiceConfiguration(configurationFilePathFlag)

	// Setup logging for the service
	logging.Setup(log.StandardLogger(), managementConfig.Observability.Logging, logPath, isDebug)

	if enableTracing && err == nil {
		f := tracing.InitTelemetry(res, managementConfig.Observability.OTel)
		defer f()
		//tracer = tracing.NewTracer(log.StandardLogger(), nil)
	}

	if enableMetrics && err == nil {
		prom = metrics.FromConfiguration(res, managementConfig.Observability.Metrics)
		go prom.Start(managementConfig.Observability.Metrics.Address)
	}

	// Connect to Database
	database.Connect(database.Database(databaseType), managementConfig.Database)

	managementService := management.NewManagementService(database.GetChatRepository(), database.GetClientRepository())

	// Start listening to the topics
	mqttClient := mqtt.NewMqttClient(managementConfig.MqttBroker)
	mqttClient.Subscribe(management.ChatAuth, managementService.ChatAuthHandler())
	mqttClient.Subscribe(management.ChatDeAuth, managementService.DeAuthHandler())
	mqttClient.Subscribe(management.ChatSetLang, managementService.SetLanguageHandler())
	mqttClient.Subscribe(management.ExecutePlugin, managementService.PluginExecutionHandler())
	mqttClient.Subscribe(management.ClientRegister, managementService.ClientRegisterHandler())
	mqttClient.Subscribe(management.ClientHeartbeat, managementService.HeartbeatHandler())
	mqttClient.Subscribe(management.ClientPluginResponse, managementService.ClientPluginRegisterResponseHandler())

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
		defaultConfigFileName = fmt.Sprintf("%s/management-service.%s", workingDirectory, "yaml")
	)

	// Set flags
	rootCmd.PersistentFlags().StringVarP(&configurationFilePathFlag, "config", "c", defaultConfigFileName, "config file path")
	rootCmd.PersistentFlags().StringVar(&databaseType, "database", "Mongo", "database type (supported: Mongo)")
	rootCmd.PersistentFlags().BoolVarP(&isDebug, debugFlag, "d", false, "debug mode")
	rootCmd.PersistentFlags().BoolVarP(&enableMetrics, metricsFlag, "m", false, "enable metrics for Prometheus")
	rootCmd.PersistentFlags().BoolVarP(&enableTracing, tracingFlag, "t", false, "enable tracing using OpenTelemetry")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Cannot run the service")
	}
}
