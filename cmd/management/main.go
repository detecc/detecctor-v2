package main

import (
	"context"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/internal/config"
	"github.com/detecc/detecctor-v2/internal/mqtt"
	"github.com/detecc/detecctor-v2/service/management"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func setupLogger(isProduction bool) {
	if isProduction {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	setupLogger(false)
	log.Info("Starting the management service...")

	// Create exit handlers
	ctx, cancel := context.WithCancel(context.Background())
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, os.Interrupt)

	// Get flags for service
	config.GetFlags(config.ManagementService)

	// Get configuration for management service
	managementConfig := config.GetManagementServiceConfiguration()

	// Connect to Mongo
	database.Connect(managementConfig.Mongo)

	// Start listening to the topics
	mqttClient := mqtt.NewMqttClient(managementConfig.MqttBroker)
	mqttClient.Subscribe(management.ChatAuth, management.ChatAuthHandler)
	mqttClient.Subscribe(management.ChatSetLang, management.SetLanguageHandler)
	mqttClient.Subscribe(management.ClientRegister, management.ClientRegisterHandler)
	mqttClient.Subscribe(management.ClientHeartbeat, management.ClientRegisterHandler)

	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping the management service..")
			break
		case <-quitChannel:
			cancel()
		}
	}
}
