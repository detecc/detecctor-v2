package database

import (
	"context"
	"fmt"
	. "github.com/detecc/detecctor-v2/model/client"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func createClient(newClient *Client) error {
	return mgm.Coll(&Client{}).Create(newClient)
}

func createClientWithCtx(ctx context.Context, newClient *Client) error {
	return mgm.Coll(&Client{}).CreateWithCtx(ctx, newClient)
}

func updateClient(client *Client) error {
	return mgm.Coll(&Client{}).Update(client)
}

func updateClientWithCtx(ctx context.Context, client *Client) error {
	return mgm.Coll(&Client{}).UpdateWithCtx(ctx, client)
}

func getClient(filter interface{}) (*Client, error) {
	client := &Client{}
	err := mgm.Coll(&Client{}).First(filter, client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getClientWithCtx(ctx context.Context, filter interface{}) (*Client, error) {
	client := &Client{}
	err := mgm.Coll(&Client{}).FirstWithCtx(ctx, filter, client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetClient(clientId string) (*Client, error) {
	log.WithField("clientId", clientId).Debug("Getting a client with Id")
	return getClient(bson.M{"clientId": clientId})
}

func GetClientWithServiceNodeKey(serviceNodeKey string) (*Client, error) {
	return getClient(bson.M{"serviceNodeKey": serviceNodeKey})
}

// GetClients returns all clients
func GetClients() ([]Client, error) {
	var (
		serviceNode = &Client{}
		results     []Client
	)
	// find all clients
	cursor, err := mgm.Coll(serviceNode).Find(mgm.Ctx(), bson.M{})
	if err = cursor.All(mgm.Ctx(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func IsClientAuthorized(clientId string) bool {
	log.WithField("clientId", clientId).Debug("Checking if client is authorized")

	sn, err := GetClient(clientId)
	if err != nil {
		return false
	}

	return sn.Status != StatusUnauthorized
}

func AuthorizeClient(clientId, serviceNodeKey string) error {
	log.WithFields(log.Fields{
		"clientId":       clientId,
		"serviceNodeKey": serviceNodeKey,
	}).Debug("Authorizing client")

	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		client, err := getClientWithCtx(sc, bson.M{"clientId": clientId})
		if err != nil {
			return err
		}

		client.ServiceNodeKey = serviceNodeKey

		err = updateClientWithCtx(sc, client)
		if err != nil {
			return err
		}

		updateClientStatusWithCtx(sc, client.ClientId, StatusOffline)

		addNodeErr := addClientWithCtx(sc)
		if addNodeErr != nil {
			log.Warning(addNodeErr)
		}

		return session.CommitTransaction(sc)
	})
}

func updateClientStatusWithCtx(ctx context.Context, clientId, status string) error {
	log.WithFields(log.Fields{
		"clientId": clientId,
		"status":   status,
	}).Debug("Updating a client")

	client, err := getClientWithCtx(ctx, bson.M{"clientId": clientId})
	if err != nil {
		return err
	}

	switch status {
	case StatusUnauthorized:
		client.Status = status
		break
	case StatusOffline:
		client.Status = status
		clientOfflineWithCtx(ctx)
		break
	case StatusOnline:
		client.Status = status
		clientOnlineWithCtx(ctx)
		break
	default:
		return fmt.Errorf("invalid client status")
	}

	return updateClientWithCtx(ctx, client)
}

func UpdateClientStatus(clientId, status string) error {
	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := updateClientStatusWithCtx(sc, clientId, status)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})

}

func CreateClientIfNotExists(clientId, IP, SNKey string) *Client {
	log.WithFields(log.Fields{
		"clientId":       clientId,
		"IP":             IP,
		"serviceNodeKey": SNKey,
	}).Debug("Creating a new client")

	client := &Client{
		IP:             IP,
		ClientId:       clientId,
		ServiceNodeKey: SNKey,
		LastOnline:     nil,
		Status:         StatusUnauthorized,
	}

	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		_, err := getClientWithCtx(sc, bson.M{"clientId": clientId})

		switch err {
		case nil:
			return fmt.Errorf("client already exists")
		case mongo.ErrNoDocuments:
			err = createClientWithCtx(sc, client)
			if err != nil {
				return err
			}
			break
		default:
			return err
		}

		return session.CommitTransaction(sc)
	})
	if err != nil {
		return nil
	}

	return client
}
