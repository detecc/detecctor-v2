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
	"github.com/detecc/detecctor-v2/service/notification/bot"
	"github.com/detecc/detecctor-v2/service/notification/proxy"
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

	logPath = "/var/log/detecctor/notification.log"
)

var (
	configurationFilePath string
	databaseType          string

	isDebug       = false
	enableMetrics = false
	enableTracing = false

	rootCmd = &cobra.Command{
		Use:   "detecctor-notification",
		Short: "Notification service",
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
				semconv.ServiceNameKey.String("notification-service"),
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

	log.Info("Starting the cmd service..")

	// Get configuration for the service
	notificationConfig := config.GetNotificationServiceConfiguration(configurationFilePath)

	// Setup logging
	logging.Setup(log.StandardLogger(), notificationConfig.Observability.Logging, logPath, isDebug)

	if enableTracing && err == nil {
		f := tracing.InitTelemetry(res, notificationConfig.Observability.OTel)
		defer f()
		//tracer = tracing.NewTracer(log.StandardLogger(), nil)
	}

	if enableMetrics && err == nil {
		prom = metrics.FromConfiguration(res, notificationConfig.Observability.Metrics)
		go prom.Start(notificationConfig.Observability.Metrics.Address)
	}

	// Connect to Database
	database.Connect(database.Database(databaseType), notificationConfig.Database)

	// Create a new bot for cmd
	nBot := bot.NewBot(notificationConfig.Bot)

	// Create a communication layer client
	mqttClient := mqtt.NewMqttClient(notificationConfig.MqttBroker)

	// Start listening to incoming requests from both the bot and the service communication layer
	botProxy := proxy.NewProxy(
		nBot,
		mqttClient,
		database.GetMessageRepository(),
		database.GetStatistics(),
		database.GetChatRepository(),
		database.GetClientRepository(),
		database.GetLogRepository(),
	)
	botProxy.Start(ctx, notificationConfig.Bot.Token)

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
		defaultConfigFileName = fmt.Sprintf("%s/notification-service.%s", workingDirectory, "yaml")
	)

	// Set flags
	rootCmd.PersistentFlags().StringVar(&configurationFilePath, "config", defaultConfigFileName, "config file path")
	rootCmd.PersistentFlags().StringVar(&databaseType, "database", "Mongo", "database type (supported: Mongo)")
	rootCmd.PersistentFlags().BoolVarP(&isDebug, debugFlag, "d", false, "debug mode")
	rootCmd.PersistentFlags().BoolVarP(&enableMetrics, metricsFlag, "m", false, "enable metrics for Prometheus")
	rootCmd.PersistentFlags().BoolVarP(&enableTracing, tracingFlag, "t", false, "enable tracing using OpenTelemetry")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Cannot run the service")
	}
}
