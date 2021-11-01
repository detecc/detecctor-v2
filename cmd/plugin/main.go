package main

import (
	"context"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/internal/config"
	"github.com/detecc/detecctor-v2/internal/mqtt"
	plugin2 "github.com/detecc/detecctor-v2/service/plugin"
	"github.com/detecc/detecctor-v2/service/plugin/plugin"
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
	log.Info("Starting the plugin service..")

	// Create exit handlers
	ctx, cancel := context.WithCancel(context.Background())
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, os.Interrupt)

	// Get the flags for the service
	config.GetFlags(config.PluginService)

	// Get the service configuration
	pluginConfig := config.GetPluginServiceConfiguration()

	// Connect to the Mongo
	database.Connect(pluginConfig.Mongo)

	// Load plugins into the manager
	plugin.GetPluginManager().LoadPlugins(pluginConfig.PluginConfiguration)

	//Start listening to the MQTT client
	mqttClient := mqtt.NewMqttClient(pluginConfig.MqttBroker)
	//mqttClient.Subscribe(plugin2.Metadata, plugin2.MetadataHandler)
	mqttClient.Subscribe(plugin2.ExecutionTopic, plugin2.ExecutionHandler)
	mqttClient.Subscribe(plugin2.ResponseTopic, plugin2.ResponseHandler)

	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping the plugin service..")
			break
		case <-quitChannel:
			cancel()
		}
	}
}
