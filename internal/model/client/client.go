package client

import (
	"github.com/detecc/detecctor-v2/internal/model/timestamp"
	"github.com/kamva/mgm/v3"
)

const (
	StatusOnline       = Status("online")
	StatusOffline      = Status("offline")
	StatusUnauthorized = Status("unauthorized")
)

type (
	Status string

	// Client object contains some basic information of a client.
	Client struct {
		mgm.DefaultModel `bson:",inline"`
		IP               string              `json:"IP" bson:"IP"`
		ClientId         string              `json:"clientId" bson:"clientId"`
		ServiceNodeKey   string              `json:"serviceNodeKey" bson:"serviceNodeKey"`
		LastOnline       *timestamp.DateTime `json:"lastOnline" bson:"lastOnline"`
		Status           Status              `json:"status" bson:"status"`
	}

	// Statistics for the telegram bot.
	// ActiveClients is a number of currently active/connected clients.
	// TotalClients is a number of all known connections and is used to calculate the number of offline clients.
	// LastMessageId is a number of the last known messageId.
	Statistics struct {
		mgm.DefaultModel `bson:",inline"`
		ActiveClients    int    `json:"activeClients" bson:"activeClients"`
		TotalClients     int    `json:"totalClients" bson:"totalClients"`
		LastMessageId    string `json:"lastMessageId" bson:"lastMessageId"`
	}
)

func (s Status) String() string {
	return string(s)
}
