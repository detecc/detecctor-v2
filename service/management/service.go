package management

import (
	"context"
	"fmt"
	"github.com/detecc/detecctor-v2/database/repositories"
	cmd "github.com/detecc/detecctor-v2/internal/model/command"
	replyBuilder "github.com/detecc/detecctor-v2/internal/reply"
	"github.com/detecc/detecctor-v2/pkg/cache"
	. "github.com/detecc/detecctor-v2/pkg/mqtt"
	"github.com/detecc/detecctor-v2/pkg/observability/tracing"
	p "github.com/detecc/detecctor-v2/pkg/payload"
	"github.com/detecc/detecctor-v2/service/management/auth"
	"github.com/mitchellh/mapstructure"
	goCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

const (
	ChatAuth             = Topic("chat/+/auth")
	ChatDeAuth           = Topic("chat/+/deauth")
	ChatSetLang          = Topic("chat/+/lang/set")
	ClientHeartbeat      = Topic("client/+/heartbeat")
	ClientRegister       = Topic("client/+/register")
	ClientPluginResponse = Topic("client/+/plugin/+/reply")
	ExecutePlugin        = Topic("plugin/cmd/+/execute")
)

type Service struct {
	chatRepository   repositories.ChatRepository
	clientRepository repositories.ClientRepository
	cache            *goCache.Cache
	tracer           *tracing.Tracer
}

func NewManagementService(chatRepository repositories.ChatRepository, clientRepository repositories.ClientRepository) *Service {
	return &Service{
		chatRepository:   chatRepository,
		clientRepository: clientRepository,
		cache:            cache.NewCache(),
	}
}

func (s *Service) ChatAuthHandler() MessageHandler {
	return func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
		var (
			chatId         = topicIds[0]
			token          string
			builder        = replyBuilder.NewReplyBuilder().TypeMessage().ForChat(chatId)
			chatReply      = fmt.Sprintf("chat/%s/notify", chatId)
			chatReplyTopic = Topic(chatReply)
			command        cmd.Command
			ctx            = context.Background()
		)

		err = mapstructure.Decode(payload, &command)
		if err != nil {
			log.Errorf("Error decoding the map: %v", err)
			return
		}

		log.Debugf("Authorizing a chat: %s", chatId)

		// Check if chat already authorized
		if s.chatRepository.IsChatAuthorized(ctx, chatId) {
			message := builder.WithContent("Chat already authorized").Build()
			_ = client.Publish(chatReplyTopic, message)
			return
		}

		// Check if the token is in the cache and if it matches the provided token
		if command.Args != nil && len(command.Args) >= 1 {
			token = command.Args[0]
		}

		cachedTokenId := fmt.Sprintf("auth-token-%s", chatId)
		cachedToken, isFound := s.cache.Get(cachedTokenId)
		if isFound && cachedToken.(string) == token {
			replyMessage := builder.WithContent("Chat successfully authorized.")

			err := s.chatRepository.AuthorizeChat(ctx, chatId)
			if err != nil {
				log.WithFields(log.Fields{
					"chatId": chatId,
					"token":  token,
				}).WithError(err).Errorf("Error authorizing chat")
				builder.WithContent(fmt.Sprintf("Error authorizing chat: %v", err))
			} else {
				s.cache.Delete(cachedTokenId)
			}

			_ = client.Publish(chatReplyTopic, replyMessage.Build())
			return
		}

		if !isFound && token == "" {
			// Generate a token
			auth.GenerateChatAuthenticationToken(s.cache, chatId)
			message := builder.WithContent("Authentication token generated. Check with the admin for the token").Build()
			_ = client.Publish(chatReplyTopic, message)
		}
	}
}

func (s *Service) DeAuthHandler() MessageHandler {
	return func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
		var (
			chatId  = topicIds[0]
			logInfo = log.WithFields(log.Fields{
				"chatId":   chatId,
				"language": "lang",
			})
			builder    = replyBuilder.NewReplyBuilder().TypeMessage().ForChat(chatId)
			replyTopic = fmt.Sprintf("Revoked the chat %s.", chatId)
			ctx        = context.Background()
		)

		// Revoke chat
		builder.WithContent(fmt.Sprintf("Revoked the chat %s.", chatId))
		err = s.chatRepository.RevokeChatAuthorization(ctx, chatId)
		if err != nil {
			builder.WithContent(fmt.Sprintf("Could not remove chat authorization: %v", err))
			logInfo.WithError(err).Errorf("Error revoking the authorization")
		}

		_ = client.Publish(Topic(replyTopic), builder.Build())
	}
}
func (s *Service) SetLanguageHandler() MessageHandler {
	return func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
		var (
			chatId  = topicIds[0]
			logInfo = log.WithFields(log.Fields{
				"chatId":   chatId,
				"language": "lang",
			})
			builder               = replyBuilder.NewReplyBuilder().TypeMessage().ForChat(chatId)
			chatNotification      = fmt.Sprintf("chat/%s/notify", chatId)
			chatNotificationTopic = Topic(chatNotification)
			command               cmd.Command
			ctx                   = context.Background()
		)

		err = mapstructure.Decode(payload, &command)
		if err != nil {
			log.Errorf("Error decoding the map: %v", err)
			return
		}

		builder.WithContent(fmt.Sprintf("Successfully set the language to: %s.", err))
		err = s.chatRepository.SetLanguage(ctx, chatId, command.Args[0])
		if err != nil {
			logInfo.Errorf("Error updating the language: %v", err)
			builder.WithContent(fmt.Sprintf("An error occured while changing the language: %v.", err))
		}

		client.Publish(chatNotificationTopic, builder.Build())
	}
}

func (s *Service) ClientRegisterHandler() MessageHandler {
	return func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
		var (
			clientId      = topicIds[0]
			clientPayload p.Payload
			logInfo       = log.WithFields(log.Fields{
				"payload":  clientPayload,
				"clientId": clientId,
			})
			ctx = context.Background()
		)

		err = mapstructure.Decode(payload, &clientPayload)
		if err != nil {
			logInfo.Errorf("Error decoding the map: %v", err)
			return
		}

		// When the client registers, create a new database entry if it doesn't exist yet
		_, err = s.clientRepository.CreateClientIfNotExists(ctx, clientId, "", "")
		if err != nil {
			logInfo.WithError(err).Warn("Unable to create a new client")
		}

		if clientPayload.Data == nil {
			err = fmt.Errorf("payload data is empty; client %s cannot be authorized", clientId)
			logInfo.WithError(err).Warn("Unable to create a client")
			clientPayload.SetError(err)
			return
		}

		//todo
		// Check if the secret is the same
		if clientPayload.Data.(string) == "" {
			return
		}

		// Try to authorize the client
		err = s.clientRepository.AuthorizeClient(ctx, clientId, clientPayload.ServiceNodeKey)
		if err != nil {
			logInfo.Errorf("Error updating the client authorization status: %v", err)
		}

		clientPayload.Success = true

		// Reply with the authorization status
		//client.Publish("", clientPayload)
	}
}

func (s Service) PluginExecutionHandler() MessageHandler {
	return func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
		var (
			pluginName = topicIds[0]
			logInfo    = log.WithFields(log.Fields{
				"pluginName": pluginName,
			})
			command               cmd.Command
			pluginExecute         = fmt.Sprintf("cmd/%s/execute", pluginName)
			chatNotification      = fmt.Sprintf("chat/%s/notify", command.ChatId)
			pluginExecuteTopic    = Topic(pluginExecute)
			chatNotificationTopic = Topic(chatNotification)
			ctx                   = context.Background()
		)

		err = mapstructure.Decode(payload, &command)
		if err != nil {
			logInfo.WithError(err).Errorf("Error decoding the payload")
			return
		}

		if s.chatRepository.IsChatAuthorized(ctx, command.ChatId) {
			// Forward execution to the cmd service
			client.Publish(pluginExecuteTopic, payload)
		} else {
			// Tell the user the chat is not authorized
			message := replyBuilder.NewReplyBuilder().ForChat(command.ChatId).TypeMessage().
				WithContent(fmt.Sprintf("Chat not authorized to issue the command")).Build()
			client.Publish(chatNotificationTopic, message)
		}
	}
}

func (s Service) ClientPluginRegisterResponseHandler() MessageHandler {
	return func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
		var (
			clientId            = topicIds[0]
			pluginName          = topicIds[1]
			pluginResponseTopic = fmt.Sprintf("cmd/%s/execute/response", pluginName)
			ctx                 = context.Background()
		)

		// Check if the client is authorized
		if s.clientRepository.IsClientAuthorized(ctx, clientId) && s.clientRepository.IsOnline(ctx, clientId) {
			client.Publish(Topic(pluginResponseTopic), payload)
		}
	}
}

func (s *Service) NotificationServiceRegisterHandler() MessageHandler {
	return func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
		//Todo
	}
}

func (s *Service) HeartbeatHandler() MessageHandler {
	return func(client Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
		var (
			clientId = topicIds[0]
			ctx      = context.Background()
		)

		// Update the lastOnline timestamp
		err = s.clientRepository.UpdateLastOnline(ctx, clientId)
		if err != nil {
			log.WithError(err).Errorf("Couldn't update last online timestamp")
		}
	}
}
