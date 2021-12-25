package plugin

import (
	"context"
	"github.com/detecc/detecctor-v2/internal/model/reply"
	"github.com/detecc/detecctor-v2/pkg/payload"
)

// constants for Metadata.Type
const (
	ServerOnly   = "serverOnly"
	ServerClient = "serverClient"
)

type (
	Handler interface {
		// Response is called when the client(s) have responded and should
		// return a Reply object to send to the Notification service and an error if there was an error producing the reply.
		Response(ctx context.Context, payload payload.Payload) (*reply.Reply, error)

		// Execute method is called when the command is issued by the user.
		// The method must return a Payload array with data to be sent to the client(s) or an error if anything went wrong.
		Execute(ctx context.Context, args ...string) ([]payload.Payload, error)

		// GetMetadata returns the metadata about the cmd.
		GetMetadata() Metadata
	}

	// Metadata is used to determine the role of a cmd registered in the PluginManager.
	Metadata struct {

		// The Type of the cmd will determine the behaviour of the server and execution of the cmd(s).
		Type string

		// The Middleware list is used to determine if the Plugin has any Middleware to execute before calling the Execute method.
		// Will be skipped if the cmd itself is registered as middleware.
		Middleware []string
	}
)
