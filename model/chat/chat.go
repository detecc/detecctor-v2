package chat

import "github.com/kamva/mgm/v3"

type (
	// Chat is the Chat the Bot is listening to.
	Chat struct {
		mgm.DefaultModel `bson:",inline"`
		ChatId           string         `json:"chatId" bson:"chatId"`
		Name             string         `json:"name" bson:"name"`
		IsAuthorized     bool           `json:"isAuthorized" bson:"isAuthorized"`
		Language         string         `json:"lang" bson:"lang"`
		Subscriptions    []Subscription `json:"subscriptions" bson:"subscriptions"`
	}

	// Subscription is a filter used for subscribing to a client messages.
	// If the chat/user is subscribed to all nodes and all topics, there should only be one entry with both subNode and subCommand values equal to "*".
	// Else, there are separate entries with values, "*" meaning all.
	// Example entry: subNode: "*", subCommand:"/ping" -> meaning subscribe to the /ping command on all nodes.
	Subscription struct {
		Client  string `json:"client" bson:"client"`
		Command string `json:"command" bson:"command"`
	}

	// Message is a Message that gets logged in the database.
	Message struct {
		mgm.DefaultModel `bson:",inline"`
		ChatId           string `json:"chatId" bson:"chatId"`
		MessageId        string `json:"messageId" bson:"messageId"`
		Content          string `json:"content" bson:"content"`
	}
)
