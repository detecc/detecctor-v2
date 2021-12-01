package main

import (
	"context"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/internal/config"
	"github.com/detecc/detecctor-v2/internal/logging"
	"github.com/detecc/detecctor-v2/internal/mqtt"
	plugin2 "github.com/detecc/detecctor-v2/service/plugin"
	"github.com/detecc/detecctor-v2/service/plugin/plugin"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	logging.SetupLogger(false)
	log.Info("Starting the plugin service..")

	// Create exit handlers
	ctx, cancel := context.WithCancel(context.Background())
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, os.Interrupt)

	// Get the flags for the service
	config.GetFlags(config.PluginService)

	// Get the service configuration
	pluginConfig := config.GetPluginServiceConfiguration()

	// Connect to the Database
	database.Connect(pluginConfig.Database)

	// Load plugins into the manager
	plugin.GetPluginManager().LoadPlugins(pluginConfig.PluginConfiguration)

	// Start listening to the MQTT client
	mqttClient := mqtt.NewMqttClient(pluginConfig.MqttBroker)
	mqttClient.Subscribe(plugin2.ExecutionTopic, plugin2.ExecutionHandler)
	mqttClient.Subscribe(plugin2.ResponseTopic, plugin2.ResponseHandler)

Loop:
	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping the plugin service..")
			mqttClient.Disconnect()
			break Loop
		case <-quitChannel:
			cancel()
			break
		}
	}
}
