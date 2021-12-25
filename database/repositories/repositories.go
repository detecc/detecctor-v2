package repositories

import (
	"context"
	. "github.com/detecc/detecctor-v2/internal/command/logs"
	. "github.com/detecc/detecctor-v2/internal/model/chat"
	"github.com/detecc/detecctor-v2/internal/model/client"
	. "github.com/detecc/detecctor-v2/internal/model/command"
)

type (
	ChatRepository interface {
		GetChatWithId(ctx context.Context, chatId string) (*Chat, error)
		GetChats(ctx context.Context) ([]Chat, error)
		AuthorizeChat(ctx context.Context, chatId string) error
		IsChatAuthorized(ctx context.Context, chatId string) bool
		RevokeChatAuthorization(ctx context.Context, chatId string) error
		AddChatIfDoesntExist(ctx context.Context, chatId string, name string) error
		GetLanguage(ctx context.Context, chatId string) (string, error)
		SetLanguage(ctx context.Context, chatId string, lang string) error
	}

	ClientRepository interface {
		GetClient(ctx context.Context, clientId string) (*client.Client, error)
		GetClientWithServiceNodeKey(ctx context.Context, serviceNodeKey string) (*client.Client, error)
		GetClients(ctx context.Context) ([]client.Client, error)
		IsOnline(ctx context.Context, clientId string) bool
		IsClientAuthorized(ctx context.Context, clientId string) bool
		AuthorizeClient(ctx context.Context, clientId, serviceNodeKey string) error
		UpdateClientStatus(ctx context.Context, clientId string, status client.Status) error
		UpdateLastOnline(ctx context.Context, clientId string) error
		CreateClientIfNotExists(ctx context.Context, clientId, IP, SNKey string) (*client.Client, error)
	}

	LogRepository interface {
		AddCommandResponse(ctx context.Context, payloadId string, option ...ResponseOption) error
		AddCommandLog(ctx context.Context, command Command, option ...Option) (string, error)
		UpdateCommandLogWithId(ctx context.Context, messageId string, options ...Option) error
	}

	MessageRepository interface {
		GetMessageFromChat(ctx context.Context, chatId int) (*Message, error)
		GetMessagesFromChat(ctx context.Context, chatId string) ([]Message, error)
		GetMessageWithId(ctx context.Context, messageId string) (*Message, error)
		NewMessage(ctx context.Context, chatId string, messageId string, content string) (*Message, error)
	}

	Statistics interface {
		GetStatistics(ctx context.Context) (*client.Statistics, error)
		UpdateLastMessageId(ctx context.Context, lastMessageId string) error
		CreateStatisticsIfNotExists(ctx context.Context)
	}

	SubscriptionRepository interface {
		GetSubscribedChats(ctx context.Context, nodeId, command string) ([]Chat, error)
		SubscribeToAll(ctx context.Context, chatId string) error
		SubscribeTo(ctx context.Context, chatId string, clients []string, commands []string) error
		UnSubscribeFromAll(ctx context.Context, chatId string) error
		UnSubscribeFrom(ctx context.Context, chatId string, clients []string, commands []string) error
	}
)
