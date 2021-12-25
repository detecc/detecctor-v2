package mongo

import (
	"context"
	"fmt"
	. "github.com/detecc/detecctor-v2/internal/model/client"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatisticsRepository struct{}

func NewStatistics() *StatisticsRepository {
	return &StatisticsRepository{}
}

func (s *StatisticsRepository) GetStatistics(ctx context.Context) (*Statistics, error) {
	log.Debug("Getting statistics")
	stats := &Statistics{}

	err := mgm.Coll(stats).First(bson.M{}, stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *StatisticsRepository) UpdateLastMessageId(ctx context.Context, lastMessageId string) error {
	log.Debug("Updating last message ID")

	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		stats, err := getStatisticsWithCtx(sc)
		if err != nil {
			return err
		}

		stats.LastMessageId = lastMessageId

		err = updateStatisticsWithCtx(sc, stats)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})
}

func (s *StatisticsRepository) CreateStatisticsIfNotExists(ctx context.Context) {
	log.Debug("Creating statistics if they don't exist already")

	statistics := &Statistics{
		ActiveClients: 0,
		TotalClients:  0,
		LastMessageId: "0",
	}

	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		_, err := getStatisticsWithCtx(sc)

		switch err {
		case nil:
			return fmt.Errorf("statistics already exist")
		case mongo.ErrNoDocuments:
			err := createStatistics(statistics)
			if err != nil {
				return err
			}
		default:
			return err
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		log.Warning(err)
	}
}

func updateStatisticsWithCtx(ctx context.Context, stats *Statistics) error {
	return mgm.Coll(&Statistics{}).UpdateWithCtx(ctx, stats)
}

func createStatistics(statistics *Statistics) error {
	return mgm.Coll(&Statistics{}).Create(statistics)
}

func clientOnlineWithCtx(ctx context.Context) error {
	log.Debug("Updating statistics for client; a client is online")

	stats, err := getStatisticsWithCtx(ctx)
	if err != nil {
		return err
	}

	// number of active nodes cannot exceed the number of total nodes
	if stats.TotalClients >= stats.ActiveClients {
		stats.ActiveClients = stats.ActiveClients + 1
	}

	return updateStatisticsWithCtx(ctx, stats)
}

func clientOfflineWithCtx(ctx context.Context) error {
	log.Debug("Updating statistics for client; a client went offline")

	stats, err := getStatisticsWithCtx(ctx)
	if err != nil {
		return err
	}

	// number of active nodes cannot exceed the number of total nodes
	if stats.ActiveClients > 0 && stats.TotalClients > stats.ActiveClients {
		stats.ActiveClients = stats.ActiveClients - 1
	}

	return updateStatisticsWithCtx(ctx, stats)
}

func removeClientWithCtx(ctx context.Context) error {
	log.Debug("Removing a client")

	stats, err := getStatisticsWithCtx(ctx)
	if err != nil {
		return err
	}

	if stats.TotalClients > 0 {
		stats.TotalClients = stats.TotalClients - 1
	}

	return updateStatisticsWithCtx(ctx, stats)
}

func addClientToStatisticsWithCtx(ctx context.Context) error {
	log.Debug("Adding a client")
	stats, err := getStatisticsWithCtx(ctx)
	if err != nil {
		return err
	}

	stats.TotalClients = stats.TotalClients + 1

	return updateStatisticsWithCtx(ctx, stats)
}

func getStatisticsWithCtx(ctx context.Context) (*Statistics, error) {
	stats := &Statistics{}

	err := mgm.Coll(&Statistics{}).FirstWithCtx(ctx, bson.M{}, stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
