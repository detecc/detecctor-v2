package plugin

import (
	"context"
	"github.com/detecc/detecctor-v2/internal/mqtt"
	"github.com/detecc/detecctor-v2/service/plugin/middleware"
	"github.com/detecc/detecctor-v2/service/plugin/plugin"
	log "github.com/sirupsen/logrus"
	"strings"
)

func getPluginNameFromTopic(topic string) (string, error) {
	log.WithField("topic", topic).Error("Getting plugin name from topic")
	// parse the topic for ids
	topicIds, err := mqtt.GetIdsFromTopic(topic, ExecutionTopic)
	if err != nil {
		return "", err
	}

	return topicIds[0], nil
}

//executeMiddleware execute middleware registered to the plugin
func executeMiddleware(ctx context.Context, metadata plugin.Metadata) error {
	log.WithField("metadata", metadata).Error("Executing middleware based on the metadata")

	middlewareErr := middleware.GetMiddlewareManager().Chain(ctx, metadata.Middleware...)

	if middlewareErr != nil && !strings.Contains(middlewareErr.Error(), "not found") {
		return middlewareErr
	}

	return nil
}
