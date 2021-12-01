package mongo

import (
	"fmt"
	"github.com/detecc/detecctor-v2/model/configuration"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect to the MongoDb instance specified in the configuration.
func Connect(credentials configuration.Database) {
	mongoDbConnection := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		credentials.Username,
		credentials.Password,
		credentials.Host,
		credentials.Port,
		credentials.Database,
	)

	log.WithFields(log.Fields{
		"host": credentials.Host,
		"port": credentials.Port,
	}).Info("Connecting to MongoDB")

	defer log.Info("Connected to the database!")
	dbOptions := options.Client().ApplyURI(mongoDbConnection)

	err := mgm.SetDefaultConfig(nil, credentials.Database, dbOptions)
	if err != nil {
		log.Fatal(err)
	}
}
