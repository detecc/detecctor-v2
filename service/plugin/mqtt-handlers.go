package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/internal/mqtt"
	"github.com/detecc/detecctor-v2/internal/payload"
	command2 "github.com/detecc/detecctor-v2/model/command"
	"github.com/detecc/detecctor-v2/model/command/logs"
	. "github.com/detecc/detecctor-v2/model/payload"
	"github.com/detecc/detecctor-v2/service/plugin/plugin"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const (
	ExecutionTopic = mqtt.Topic("plugin/+/execute")
	ResponseTopic  = mqtt.Topic("plugin/+/execute/response")
)

var ExecutionHandler = func(client mqtt2.Client, message mqtt2.Message) {
	log.WithFields(log.Fields{
		"topic":   message.Topic(),
		"payload": message.Payload(),
	}).Debug("Received plugin execution request")

	// parse the topic for ids
	pluginName, err := getPluginNameFromTopic(message.Topic())
	if err != nil {
		return
	}

	command := command2.Command{}
	err = json.Unmarshal(message.Payload(), &command)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"payload": message.Payload(),
			"topic":   message.Topic(),
		}).Error("Could not unmarshall the payload")
		return
	}

	// get a plugin if it exists
	mPlugin, err := plugin.GetPluginManager().GetPlugin(pluginName)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"pluginName": pluginName,
			"topic":      message.Topic(),
		}).Error("Plugin doesn't exist")

		database.UpdateCommandLogWithId(command.MessageId, logs.WithErrors(err))
		return
	}

	// execute any middleware before executing the actual plugin
	middlewareExecErr := executeMiddleware(context.Background(), mPlugin.GetMetadata())

	// execute the plugin with given arguments

	log.WithField("pluginName", pluginName).Debug("Executing the plugin")
	payloads, err := mPlugin.Execute()
	database.UpdateCommandLogWithId(command.MessageId, logs.WithPayloads(payloads...), logs.WithErrors(err, middlewareExecErr))
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"pluginName": pluginName,
			"messageId":  command.MessageId,
		}).Error("Error executing the plugin")
		return
	}

	switch mPlugin.GetMetadata().Type {
	case plugin.ServerClient:
		log.WithFields(log.Fields{
			"messageId":  command.MessageId,
			"pluginName": pluginName,
			"payloads":   payloads,
		}).Debug("Sending payloads to the clients")

		for _, response := range payloads {
			// send the data to the client(s)
			payload.GeneratePayloadId(&response, command.ChatId)
			client.Publish("client/{id}/plugin/{pluginName}/execute", 1, false, response)
		}
		break
	case plugin.ServerOnly:
		break
	}
}

var ResponseHandler = func(client mqtt2.Client, message mqtt2.Message) {
	log.WithFields(log.Fields{
		"topic":   message.Topic(),
		"payload": message.Payload(),
	}).Debug("Received client plugin execution response")

	// parse the topic for ids
	pluginName, err := getPluginNameFromTopic(message.Topic())
	if err != nil {
		return
	}

	payload := Payload{}
	err = json.Unmarshal(message.Payload(), &payload)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"payload": message.Payload(),
			"topic":   message.Topic(),
		}).Error("Could not unmarshall the payload")
		return
	}

	// get the plugin from the manager
	mPlugin, err := plugin.GetPluginManager().GetPlugin(pluginName)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"pluginName": pluginName,
			"topic":      message.Topic(),
		}).Error("Plugin doesn't exist")
		return
	}

	// process the client response
	log.Println("Processing the response..")
	reply, err := mPlugin.Response(payload)
	database.AddCommandResponse(payload.Id, logs.WithResponse(reply))
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"payloadId": payload,
			"reply":     reply,
			"topic":     message.Topic(),
		}).Error("The client response could not be processed")
	}

	// forward the message to the notification service
	log.WithFields(log.Fields{
		"pluginName": pluginName,
		"payloadId":  payload.Id,
		"topic":      message.Topic(),
		"reply":      reply,
	}).Debug("Sending the response to the notification service")
	client.Publish(fmt.Sprintf("chat/%s/notify", reply.ChatId), 1, false, reply)
}
