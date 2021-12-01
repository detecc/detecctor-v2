package main

import (
	"context"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/internal/config"
	"github.com/detecc/detecctor-v2/internal/logging"
	"github.com/detecc/detecctor-v2/internal/mqtt"
	"github.com/detecc/detecctor-v2/service/management"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	logging.SetupLogger(false)
	log.Info("Starting the management service...")

	// Create exit handlers
	ctx, cancel := context.WithCancel(context.Background())
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, os.Interrupt)

	// Get flags for service
	config.GetFlags(config.ManagementService)

	// Get configuration for management service
	managementConfig := config.GetManagementServiceConfiguration()

	// Connect to Database
	database.Connect(managementConfig.Database)

	// Start listening to the topics
	mqttClient := mqtt.NewMqttClient(managementConfig.MqttBroker)
	mqttClient.Subscribe(management.ChatAuth, management.ChatAuthHandler)
	mqttClient.Subscribe(management.ChatSetLang, management.SetLanguageHandler)
	mqttClient.Subscribe(management.ClientRegister, management.ClientRegisterHandler)
	mqttClient.Subscribe(management.ClientHeartbeat, management.ClientRegisterHandler)

Loop:
	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping the management service..")
			mqttClient.Disconnect()
			break Loop
		case <-quitChannel:
			cancel()
			break
		}
	}
}
