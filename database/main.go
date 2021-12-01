package database

import (
	"github.com/detecc/detecctor-v2/database/mongo"
	"github.com/detecc/detecctor-v2/database/repositories"
	"github.com/detecc/detecctor-v2/model/configuration"
	log "github.com/sirupsen/logrus"
)

type Database string

const (
	Mongo = Database("Database")
)

var (
	databaseType = Mongo // Default database is mongo

	chatRepository         repositories.ChatRepository
	clientRepository       repositories.ClientRepository
	logRepository          repositories.LogRepository
	messageRepository      repositories.MessageRepository
	statistics             repositories.Statistics
	subscriptionRepository repositories.SubscriptionRepository
)

// initDatabase check if database exists, if not, create a database.
func initDatabase(credentials configuration.Database) {
	log.Info("Initializing database..")

	GetStatistics().CreateStatisticsIfNotExists(nil)
}

// Connect to the MongoDb instance specified in the configuration.
func Connect(credentials configuration.Database) {
	switch databaseType {
	case Mongo:
		mongo.Connect(credentials)
		break
	default:
		log.Fatalf("Database unsupported: %s", databaseType)
	}

	initDatabase(credentials)
}

func GetChatRepository() repositories.ChatRepository {
	if chatRepository == nil {
		switch databaseType {
		case Mongo:
			chatRepository = mongo.NewChatRepository()
			break
		default:
			log.Fatalf("Database unsupported: %s", databaseType)
		}
	}
	return chatRepository
}

func GetClientRepository() repositories.ClientRepository {
	if clientRepository == nil {
		switch databaseType {
		case Mongo:
			clientRepository = mongo.NewClientRepository()
			break
		default:
			log.Fatalf("Database unsupported: %s", databaseType)
		}
	}

	return clientRepository
}

func GetStatistics() repositories.Statistics {
	if statistics == nil {
		switch databaseType {
		case Mongo:
			statistics = mongo.NewStatistics()
			break
		default:
			log.Fatalf("Database unsupported: %s", databaseType)
		}
	}

	return statistics
}

func GetMessageRepository() repositories.MessageRepository {
	if messageRepository == nil {
		switch databaseType {
		case Mongo:
			messageRepository = mongo.NewMessageRepository()
			break
		default:
			log.Fatalf("Database unsupported: %s", databaseType)
		}
	}

	return messageRepository
}

func GetLogRepository() repositories.LogRepository {
	if logRepository == nil {
		switch databaseType {
		case Mongo:
			logRepository = mongo.NewLogRepository()
			break
		default:
			log.Fatalf("Database unsupported: %s", databaseType)
		}
	}

	return logRepository
}

func GetSubscriptionsRepository() repositories.SubscriptionRepository {
	if subscriptionRepository == nil {
		switch databaseType {
		case Mongo:
			subscriptionRepository = mongo.NewSubscriptionRepository()
			break
		default:
			log.Fatalf("Database unsupported: %s", databaseType)
		}
	}

	return subscriptionRepository
}
