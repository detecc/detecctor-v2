package mongo

import (
	"context"
	. "github.com/detecc/detecctor-v2/model/chat"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriptionRepository struct{}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{}
}

// GetSubscribedChats get all the chats that include subscription(s) where the nodeId == nodeId and command == command
// or either node == * or command == *.
func (s *SubscriptionRepository) GetSubscribedChats(ctx context.Context, nodeId, command string) ([]Chat, error) {
	return getChats(
		bson.M{"subscriptions.nodeId": bson.M{
			operator.In: bson.A{nodeId, "*"},
		},
			"subscriptions.command": bson.M{
				operator.In: bson.A{command, "*"},
			},
		},
	)
}

func (s *SubscriptionRepository) SubscribeToAll(ctx context.Context, chatId string) error {
	log.WithField("chatId", chatId).Debug("Subscribing to all clients and commands")

	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		chat, err := getChatWithCtx(sc, bson.M{"chatId": chatId})
		if err != nil {
			return err
		}

		// this overwrites any previous subscriptions.
		chat.Subscriptions = []Subscription{
			{
				Client:  "*",
				Command: "*",
			},
		}

		err = updateChatWithCtx(sc, chat)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})
}

func (s *SubscriptionRepository) SubscribeTo(ctx context.Context, chatId string, clients []string, commands []string) error {
	log.WithFields(log.Fields{
		"chatId":   chatId,
		"clients":  clients,
		"commands": commands,
	}).Debug("Subscribing to clients and/or commands")

	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		chat, err := getChatWithCtx(sc, bson.M{"chatId": chatId})
		if err != nil {
			return err
		}

		// check if there is a subscription to all clients and commands
		if len(chat.Subscriptions) == 1 {
			firstSubscription := chat.Subscriptions[0]
			if firstSubscription.Client == "*" && firstSubscription.Command == "*" {
				// replace the all subscription with provided subscriptions
				chat.Subscriptions = createSubscriptions(clients, commands)
				err = updateChatWithCtx(sc, chat)
				if err != nil {
					return err
				}

				return session.CommitTransaction(sc)
			}
		}

		subs := createSubscriptions(clients, commands)

		for _, sub := range subs {
			isDuplicateFound := false
			// check if there is an existing subscription for a node and command
			for _, subscription := range chat.Subscriptions {
				if sub.Client == subscription.Client && subscription.Command == sub.Command {
					isDuplicateFound = true
					break
				}
			}

			if !isDuplicateFound {
				chat.Subscriptions = append(chat.Subscriptions, sub)
			}
		}

		err = updateChatWithCtx(sc, chat)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})
}

func (s *SubscriptionRepository) UnSubscribeFromAll(ctx context.Context, chatId string) error {
	log.WithField("chatId", chatId).Debug("Unsubscribing from all")

	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		chat, err := getChatWithCtx(sc, bson.M{"chatId": chatId})
		if err != nil {
			return err
		}

		// this overwrites any previous subscriptions.
		chat.Subscriptions = []Subscription{}

		err = updateChatWithCtx(sc, chat)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})
}

func (s *SubscriptionRepository) UnSubscribeFrom(ctx context.Context, chatId string, clients []string, commands []string) error {
	log.WithFields(log.Fields{
		"chatId":   chatId,
		"clients":  clients,
		"commands": commands,
	}).Debug("Unsubscribing from clients and/or commands")

	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		chat, err := getChatWithCtx(sc, bson.M{"chatId": chatId})
		if err != nil {
			return err
		}

		// check if the first subscription is for all clients
		if len(chat.Subscriptions) == 1 {
			firstSubscription := chat.Subscriptions[0]
			if firstSubscription.Client == "*" && firstSubscription.Command == "*" {
				chat.Subscriptions = []Subscription{}

				err = updateChatWithCtx(sc, chat)
				if err != nil {
					return err
				}

				return session.CommitTransaction(sc)
			}
		}

		// remove any subscriptions listed
		for i, subscription := range chat.Subscriptions {
			for _, node := range clients {
				if node == "*" {
					if i+1 < len(chat.Subscriptions) {
						chat.Subscriptions = append(chat.Subscriptions[:i], chat.Subscriptions[i+1:]...)
						continue
					}

					chat.Subscriptions = append(chat.Subscriptions[:i])
					continue
				}

				for _, command := range commands {
					if command == "*" || (command == subscription.Command && node == subscription.Client) {

						if i+1 < len(chat.Subscriptions) {
							chat.Subscriptions = append(chat.Subscriptions[:i], chat.Subscriptions[i+1:]...)
							continue
						}

						chat.Subscriptions = append(chat.Subscriptions[:i])
						continue
					}
				}
			}
		}

		err = updateChatWithCtx(sc, chat)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})
}

func createSubscriptions(clients []string, commands []string) []Subscription {
	var subscriptions []Subscription

	for _, clientId := range clients {
		// check if the node exists
		_, err := getClient(bson.M{"serviceNodeKey": clientId})
		if err != nil && clientId != "*" {
			log.WithField("clientId", clientId).Debug("Error creating subscription; client doesn't exist")
			continue
		}

		for _, command := range commands {
			subscriptions = append(subscriptions, Subscription{Client: clientId, Command: command})
		}
	}

	return subscriptions
}
