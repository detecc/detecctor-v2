package mongo

import (
	"context"
	. "github.com/detecc/detecctor-v2/model/command"
	. "github.com/detecc/detecctor-v2/model/command/logs"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogRepository struct{}

func NewLogRepository() *LogRepository {
	return &LogRepository{}
}

func (l *LogRepository) AddCommandResponse(ctx context.Context, payloadId string, option ...ResponseOption) error {
	log.WithField("payloadId", payloadId).Debug("Adding a response for a command")

	commandResponse := NewCommandResponseLog(payloadId, option...)
	return addNewCommandResponse(commandResponse)
}

func (l *LogRepository) AddCommandLog(ctx context.Context, command Command, option ...Option) (string, error) {
	log.WithField("messageId", command.MessageId).Debug("Adding a log for command")

	commandLog := NewCommandLog(command, option...)

	err := addNewCommandLog(commandLog)
	if err != nil {
		return "", err
	}

	return commandLog.ID.String(), nil
}

func (l *LogRepository) UpdateCommandLogWithId(ctx context.Context, messageId string, options ...Option) error {
	log.WithField("messageId", messageId).Debug("Updating a command log")

	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		cmd := &CommandLog{}
		err := mgm.Coll(&CommandLog{}).FirstWithCtx(sc, bson.M{"command": bson.M{"messageId": messageId}}, cmd)
		if err != nil {
			return err
		}

		for _, opt := range options {
			opt(cmd)
		}

		err = mgm.Coll(&CommandLog{}).Update(cmd)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})
}

func addNewCommandLog(commandLog *CommandLog) error {
	return mgm.Coll(&CommandLog{}).Create(commandLog)
}

func addNewCommandResponse(commandResponse *CommandResponseLog) error {
	return mgm.Coll(&CommandResponseLog{}).Create(commandResponse)
}
