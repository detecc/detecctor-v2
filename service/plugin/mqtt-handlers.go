package plugin

import (
	"context"
	"fmt"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/internal/mqtt"
	pl "github.com/detecc/detecctor-v2/internal/payload"
	command2 "github.com/detecc/detecctor-v2/model/command"
	"github.com/detecc/detecctor-v2/model/command/logs"
	. "github.com/detecc/detecctor-v2/model/payload"
	"github.com/detecc/detecctor-v2/service/plugin/plugin"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

const (
	ExecutionTopic = mqtt.Topic("plugin/+/execute")
	ResponseTopic  = mqtt.Topic("plugin/+/execute/response")
)

var ExecutionHandler = func(client mqtt.Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
	pluginName := topicIds[0]

	logInfo := log.WithFields(log.Fields{
		"pluginName": pluginName,
		"payload":    payload,
	})
	logInfo.Debug("Received plugin execution request")

	command := command2.Command{}
	err = mapstructure.Decode(payload, &command)
	if err != nil {
		logInfo.Errorf("Error decoding the map: %v", err)
		return
	}

	// Get a plugin if it exists
	mPlugin, err := plugin.GetPluginManager().GetPlugin(pluginName)
	if err != nil {
		logInfo.Errorf("Plugin doesnt exist: %v", err)
		database.GetLogRepository().UpdateCommandLogWithId(nil, command.MessageId, logs.WithErrors(err))
		return
	}

	// Execute any middleware before executing the actual plugin
	middlewareExecErr := executeMiddleware(context.Background(), mPlugin.GetMetadata())

	// Execute the plugin with given arguments
	logInfo.Debug("Executing the plugin")
	payloads, err := mPlugin.Execute()
	database.GetLogRepository().UpdateCommandLogWithId(nil, command.MessageId, logs.WithPayloads(payloads...), logs.WithErrors(err, middlewareExecErr))
	if err != nil {
		logInfo.Errorf("Error executing the plugin for message %s: %v", command.MessageId, err)
		return
	}

	switch mPlugin.GetMetadata().Type {
	case plugin.ServerClient:
		logInfo.Debugf("Sending payloads to the clients: %v", payloads)

		for _, response := range payloads {
			// Send the data to the client(s)
			pl.GeneratePayloadId(&response, command.ChatId)
			client.Publish("client/{id}/plugin/{pluginName}/execute", response)
		}
		break
	case plugin.ServerOnly:
		break
	}
}

var ResponseHandler = func(client mqtt.Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
	logInfo := log.WithFields(log.Fields{
		"payload": payload,
		"topic":   topicIds[0],
	})

	pluginName := topicIds[0]
	logInfo.Debug("Received client plugin execution response")

	nPayload := payload.(Payload)

	// Get the plugin from the manager
	mPlugin, err := plugin.GetPluginManager().GetPlugin(pluginName)
	if err != nil {
		logInfo.Errorf("Plugin doesn't exist: %v", err)
		return
	}

	// Process the client response
	logInfo.Info("Processing the response..")
	reply, err := mPlugin.Response(nPayload)
	database.GetLogRepository().AddCommandResponse(nil, nPayload.Id, logs.WithResponse(reply))
	if err != nil {
		logInfo.Errorf("The client response could not be processed: %v", err)
	}

	// Forward the message to the notification service
	logInfo.Debugf("Sending the response to the notification service: %v", reply)

	client.Publish(mqtt.Topic(fmt.Sprintf("chat/%s/notify", reply.ChatId)), reply)
}
