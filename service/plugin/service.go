package plugin

import (
	"context"
	"fmt"
	"github.com/detecc/detecctor-v2/database/repositories"
	"github.com/detecc/detecctor-v2/internal/command/logs"
	cmd "github.com/detecc/detecctor-v2/internal/model/command"
	"github.com/detecc/detecctor-v2/pkg/cache"
	"github.com/detecc/detecctor-v2/pkg/mqtt"
	"github.com/detecc/detecctor-v2/pkg/observability/tracing"
	. "github.com/detecc/detecctor-v2/pkg/payload"
	"github.com/detecc/detecctor-v2/service/plugin/plugin"
	"github.com/mitchellh/mapstructure"
	goCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	ExecutionTopic = mqtt.Topic("cmd/+/execute")
	ResponseTopic  = mqtt.Topic("cmd/+/execute/response")
)

type Service struct {
	pluginManager plugin.Manager
	logRepository repositories.LogRepository
	cache         *goCache.Cache
	tracer        *tracing.Tracer
}

func NewPluginService(manager plugin.Manager, logRepository repositories.LogRepository) *Service {
	return &Service{
		pluginManager: manager,
		logRepository: logRepository,
		cache:         cache.NewCache(),
	}
}

func (s *Service) ExecutionHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
		var (
			pluginName = topicIds[0]
			logInfo    = log.WithFields(log.Fields{
				"pluginName": pluginName,
				"payload":    payload,
			})
			command = cmd.Command{}
			ctx     = context.Background()
		)

		logInfo.Debug("Received cmd execution request")

		err = mapstructure.Decode(payload, &command)
		if err != nil {
			logInfo.Errorf("Error decoding the map: %v", err)
			return
		}

		// Get a cmd if it exists
		mPlugin, err := s.pluginManager.GetPlugin(pluginName)
		if err != nil {
			logInfo.WithError(err).Errorf("Plugin doesnt exist")
			s.logRepository.UpdateCommandLogWithId(ctx, command.MessageId, logs.WithErrors(err))
			// todo notify the chat
			return
		}

		execCtx, cancel := context.WithTimeout(ctx, time.Second*30)

		// Execute any Middleware before executing the actual cmd
		middlewareExecErr := executeMiddleware(execCtx, mPlugin.GetMetadata())

		// Execute the cmd with given arguments
		logInfo.Debug("Executing the cmd")
		payloads, err := mPlugin.Execute(execCtx, command.Args...)

		s.logRepository.UpdateCommandLogWithId(ctx, command.MessageId, logs.WithPayloads(payloads...), logs.WithErrors(err, middlewareExecErr))
		if err != nil {
			logInfo.WithError(err).Errorf("Error executing the cmd for message %s", command.MessageId)
			// todo notify the chat
			cancel()
			return
		}

		cancel()

		switch mPlugin.GetMetadata().Type {
		case plugin.ServerClient:
			logInfo.Debugf("Sending payloads to the clients: %v", payloads)

			for _, response := range payloads {
				// Send the data to the client(s)
				GeneratePayloadId(s.cache, &response, command.ChatId)
				clientTopic := fmt.Sprintf("client/%s/cmd/%s/execute", response.ServiceNodeKey, pluginName)
				client.Publish(mqtt.Topic(clientTopic), response)
			}
			break
		case plugin.ServerOnly:
			break
		}
	}
}
func (s *Service) ResponseHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
		var (
			pluginName = topicIds[0]
			logInfo    = log.WithFields(log.Fields{
				"payload": payload,
				"topic":   topicIds[0],
			})
			nPayload = Payload{}
			chatId   = ""
			ctx      = context.Background()
		)

		logInfo.Debug("Received client cmd execution response")

		err = mapstructure.Decode(payload, &nPayload)
		if err != nil {
			logInfo.WithError(err).Errorf("Error decoding the map")
			return
		}

		// Get the chatId
		var replyTopic mqtt.Topic
		chat, isFound := s.cache.Get(nPayload.Id)
		if isFound {
			chatId = chat.(string)
			replyTopic = mqtt.Topic(fmt.Sprintf("chat/%s/notify", chatId))
		}

		// Get the cmd from the manager
		mPlugin, err := s.pluginManager.GetPlugin(pluginName)
		if err != nil {
			logInfo.Errorf("Plugin doesn't exist: %v", err)
			client.Publish(replyTopic, "")
			return
		}

		execCtx, cancel := context.WithTimeout(ctx, time.Second*30)

		// Process the client response
		logInfo.Info("Processing the response..")
		reply, err := mPlugin.Response(execCtx, nPayload)

		s.logRepository.AddCommandResponse(ctx, nPayload.Id, logs.WithResponse(reply))
		if err != nil {
			cancel()
			logInfo.WithError(err).Errorf("The client response could not be processed")
			client.Publish(replyTopic, "")
			return
		}

		cancel()

		// Forward the message to the cmd service
		logInfo.Debugf("Sending the response to the cmd service: %v", reply)

		client.Publish(replyTopic, reply)
	}
}
