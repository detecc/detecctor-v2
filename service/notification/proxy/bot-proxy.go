package proxy

import (
	"context"
	"fmt"
	"github.com/detecc/detecctor-v2/database/repositories"
	"github.com/detecc/detecctor-v2/internal/command/logs"
	"github.com/detecc/detecctor-v2/internal/model/command"
	"github.com/detecc/detecctor-v2/internal/model/reply"
	replyBuilder "github.com/detecc/detecctor-v2/internal/reply"
	"github.com/detecc/detecctor-v2/pkg/mqtt"
	"github.com/detecc/detecctor-v2/service/notification"
	"github.com/detecc/detecctor-v2/service/notification/bot"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

const (
	ReceiveNotificationTopic = mqtt.Topic("chat/+/notify")
)

type (
	Proxy interface {
		Start(ctx context.Context, token string)
		listenForMessages(ctx context.Context)
		cleanup()
		processMessage(message notification.ProxyMessage)
	}

	// ProxyImpl acts as a proxy between the Bot implementation and the cmd service.
	// It performs all the database operations and handles the bot communication.
	ProxyImpl struct {
		once              sync.Once
		bot               bot.Bot
		client            mqtt.Client
		replyBuilder      *replyBuilder.Builder
		messageRepository repositories.MessageRepository
		statistics        repositories.Statistics
		chatRepository    repositories.ChatRepository
		clientRepository  repositories.ClientRepository
		logRepository     repositories.LogRepository
	}
)

//NewProxy creates a new proxy for a bot and
func NewProxy(
	bot bot.Bot,
	client mqtt.Client,
	messageRepository repositories.MessageRepository,
	statistics repositories.Statistics,
	chatRepository repositories.ChatRepository,
	clientRepository repositories.ClientRepository,
	logRepository repositories.LogRepository,
) Proxy {
	return &ProxyImpl{
		bot:               bot,
		client:            client,
		once:              sync.Once{},
		replyBuilder:      replyBuilder.NewReplyBuilder(),
		messageRepository: messageRepository,
		statistics:        statistics,
		chatRepository:    chatRepository,
		clientRepository:  clientRepository,
		logRepository:     logRepository,
	}
}

// Start the bot and listening for messages from the bot.
func (proxy *ProxyImpl) Start(ctx context.Context, token string) {
	log.Info("Starting the proxy...")
	proxy.once.Do(func() {
		handler := func(client mqtt.Client, topicIds []string, payloadId uint16, payload interface{}, err error) {
			replyMsg := reply.Reply{}

			err = mapstructure.Decode(payload, &replyMsg)
			if err != nil {
				log.WithError(err).Errorf("Error decoding the message")
				return
			}

			// If the message is translatable, translate it before replying to the bot
			if replyMsg.ReplyType == reply.TranslatableMessage {
				replyMessage, err := TranslateReplyMessage(proxy.chatRepository, replyMsg.ChatId, replyMsg.Content)
				if err != nil {
					log.WithFields(log.Fields{
						"chatId":  replyMsg.ChatId,
						"content": replyMsg.Content,
					}).WithError(err).Errorf("Unable to translate the reply message")
					return
				}

				replyMsg.ReplyType = reply.TypeMessage
				replyMsg.Content = replyMessage
			}

			proxy.bot.ReplyToChat(replyMsg)
		}

		proxy.client.Subscribe(ReceiveNotificationTopic, handler)

		proxy.bot.Start(token)
		go proxy.listenForMessages(ctx)
		go proxy.bot.ListenToChannels(ctx)
	})
}

// cleanup cleans up the bot and other components
func (proxy *ProxyImpl) cleanup() {
	log.Debug("Cleaning up the proxy & bot")
}

// listenForMessages Listens and to messages from the Bot implementation and processes the message.
func (proxy *ProxyImpl) listenForMessages(ctx context.Context) {
	log.Info("Starting to listen for messages..")

Listener:
	for {
		select {
		case message := <-proxy.bot.GetMessageChannel():
			// Ignore any non-Message Updates
			proxy.processMessage(message)
			break
		case <-ctx.Done():
			log.Info("Stopping the proxy listener..")
			proxy.cleanup()
			break Listener
		}
	}
}

// processMessage processes an incoming bot message.
// Adds the necessary information to the database and forwards the command to an appropriate service.
// If an error occurs during processing, the proxy will send a response with the error to the bot.
func (proxy *ProxyImpl) processMessage(message notification.ProxyMessage) {
	ctx := context.Background()
	logInfo := log.WithFields(log.Fields{
		"chatId":    message.ChatId,
		"messageId": message.MessageId,
		"message":   message,
	})
	logInfo.Debug("Processing a message")

	err := proxy.chatRepository.AddChatIfDoesntExist(ctx, message.ChatId, message.Username)
	if err != nil {
		logInfo.WithError(err).Warnf("Error adding the chat")
	}

	// Add a new message to database
	_, err = proxy.messageRepository.NewMessage(ctx, message.ChatId, message.MessageId, message.Message)
	if err != nil {
		logInfo.WithError(err).Errorf("Error adding a message")
		return
	}

	// Update last message id
	err = proxy.statistics.UpdateLastMessageId(ctx, message.MessageId)
	if err != nil {
		logInfo.WithError(err).Errorf("Error updating the lastMessageId")
		return
	}

	cmd, err := parseCommand(message.Message, message.ChatId, message.MessageId)
	if err != nil {
		logInfo.WithError(err).Errorf("Error parsing a message as a command")

		// If the command is invalid, notify the user
		replyMessage := proxy.replyBuilder.TypeMessage().ForChat(message.ChatId).WithContent(fmt.Sprintf("%s is not a command", message)).Build()
		proxy.bot.ReplyToChat(replyMessage)
		return
	}

	// Check if command has the / prefix and remove it because of mqtt topics
	pluginName := cmd.Name
	if strings.HasPrefix(pluginName, "/") {
		pluginName = strings.TrimPrefix(pluginName, "/")
	}

	// Default publish topic
	publishTopic := fmt.Sprintf("cmd/cmd/%s/execute", pluginName)

	switch pluginName {
	case command.AuthCommand, command.AuthShortCommand:
		publishTopic = fmt.Sprintf("chat/%s/auth", cmd.ChatId)
		break
	case command.DeAuthCommand:
		publishTopic = fmt.Sprintf("chat/%s/deauth", cmd.ChatId)
		break
	case command.SubscribeShortCommand, command.SubscribeCommand:
		publishTopic = fmt.Sprintf("chat/%s/subscribe", cmd.ChatId)
		break
	case command.UnsubscribeShortCommand, command.UnsubscribeCommand:
		publishTopic = fmt.Sprintf("chat/%s/unsubscribe", cmd.ChatId)
		break
	}

	// Send the command to the designated topic that will handle the action
	err = proxy.client.Publish(mqtt.Topic(publishTopic), cmd)
	if err != nil {
		logInfo.WithError(err).Errorf("Error sending the command to the topic %s", publishTopic)
	}

	proxy.logRepository.AddCommandLog(ctx, cmd, logs.WithErrors(err))
}
