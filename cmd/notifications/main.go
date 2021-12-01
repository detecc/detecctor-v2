package main

import (
	"context"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/internal/config"
	"github.com/detecc/detecctor-v2/internal/logging"
	"github.com/detecc/detecctor-v2/internal/mqtt"
	"github.com/detecc/detecctor-v2/service/notification/bot"
	"github.com/detecc/detecctor-v2/service/notification/proxy"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	logging.SetupLogger(false)
	log.Info("Starting notification service..")

	// Create exit handlers
	ctx, cancel := context.WithCancel(context.Background())
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, os.Interrupt)

	// Get service flags
	config.GetFlags(config.NotificationService)

	// Get configuration for the service
	notificationConfig := config.GetNotificationServiceConfiguration()

	// Connect to Database
	database.Connect(notificationConfig.Database)

	// Create a new bot for notifications
	nBot := bot.NewBot(notificationConfig.Bot)

	// Create a communication layer client
	mqttClient := mqtt.NewMqttClient(notificationConfig.MqttBroker)

	// Start listening to incoming requests from both the bot and the service communication layer
	botProxy := proxy.NewProxy(nBot, mqttClient)
	botProxy.Start(ctx)

Loop:
	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping the notification service..")
			break Loop
		case <-quitChannel:
			cancel()
			break
		}
	}
}
