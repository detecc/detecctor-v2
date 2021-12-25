package bot

import (
	"context"
	"fmt"
	"github.com/detecc/detecctor-v2/database"
	"github.com/detecc/detecctor-v2/internal/model/reply"
	"github.com/detecc/detecctor-v2/service/notification"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// Telegram is a wrapper for the Telegram bot API.
type Telegram struct {
	botAPI         *telegram.BotAPI
	messageChannel chan notification.ProxyMessage
	messageBuilder *notification.MessageBuilder
}

func NewTelegramProvider(messageChannel chan notification.ProxyMessage) *Telegram {
	telegramProvider := new(Telegram)
	telegramProvider.messageChannel = messageChannel
	telegramProvider.messageBuilder = notification.NewMessageBuilder()

	return telegramProvider
}

// Start listening to the bot updates and the updates from the cmd service
func (t *Telegram) Start(token string) {
	telegramBot, err := telegram.NewBotAPI(token)
	if err != nil {
		log.WithError(err).Panic("Cannot start the telegram bot")
	}

	t.botAPI = telegramBot
}

//GetMessageChannel returns the channel for outgoing messages
func (t *Telegram) GetMessageChannel() <-chan notification.ProxyMessage {
	return t.messageChannel
}

// ListenToChannels listens for incoming data from telegram bot messages
func (t *Telegram) ListenToChannels(ctx context.Context) {
	log.Infof("Authorized on account %s", t.botAPI.Self.UserName)
	var (
		databaseCtx, cancel = context.WithTimeout(ctx, time.Second*10)
		message, err        = database.GetStatistics().GetStatistics(databaseCtx)
		lastMessageId       = 0
	)

	if err == nil {
		messageId, err := strconv.Atoi(message.LastMessageId)
		if err == nil {
			lastMessageId = messageId
		} else {
			log.WithError(err).Errorf("Error converting message id to int")
		}
	}

	cancel()

	u := telegram.NewUpdate(lastMessageId)
	u.Timeout = 60

	updates, err := t.botAPI.GetUpdatesChan(u)
	if err != nil {
		log.WithError(err).Errorf("Error receiving update channel from Telegram")
		return
	}

Listener:
	for {
		select {
		case update := <-updates:
			if update.Message == nil || update.Message.Entities == nil || len(*update.Message.Entities) == 0 {
				continue
			}

			for _, entity := range *update.Message.Entities {
				if entity.Type == "bot_command" {
					chatId := fmt.Sprintf("%d", update.Message.Chat.ID)
					messageId := fmt.Sprintf("%d", update.Message.MessageID)

					t.messageBuilder.WithId(chatId).FromUser(update.Message.Chat.UserName).WithMessage(messageId, update.Message.Text)
					t.messageChannel <- t.messageBuilder.Build()
				}
			}
			break
		case <-ctx.Done():
			log.Info("Stopping the Telegram bot..")
			break Listener
		}
	}
}

//ReplyToChat replies to a telegram chat
func (t *Telegram) ReplyToChat(replyMessage reply.Reply) {
	var (
		msg     telegram.Chattable
		logInfo = log.WithFields(log.Fields{
			"chatId":    replyMessage.ChatId,
			"replyType": replyMessage.ReplyType,
		})
	)

	chatId, err := strconv.Atoi(replyMessage.ChatId)
	if err != nil {
		logInfo.WithError(err).Errorf("Error converting ChatId to int")
		return
	}

	switch replyMessage.ReplyType {
	case reply.TypeMessage:
		if replyMessage.Content != nil {
			msg = telegram.NewMessage(int64(chatId), replyMessage.Content.(string))
		}
		break
	case reply.TypePhoto:
		msg = telegram.NewPhotoUpload(int64(chatId), replyMessage.Content)
		break
	default:
		return
	}

	if msg != nil {
		_, sendErr := t.botAPI.Send(msg)
		if sendErr != nil {
			logInfo.WithError(sendErr).Errorf("Error sending the message to chat")
		}
	}
}
