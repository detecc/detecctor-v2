package logs

import (
	"github.com/detecc/detecctor-v2/internal/model/command"
	"github.com/detecc/detecctor-v2/pkg/payload"
	"github.com/kamva/mgm/v3"
)

type (
	// CommandLog contains information about a Command that was issued to the server, processing of the command
	// and the payloads produced by the cmd as well as any errors that occurred.
	CommandLog struct {
		mgm.DefaultModel `bson:",inline"`
		Command          command.Command   `json:"command" bson:"command" validate:"required"`
		Errors           []interface{}     `json:"errors,omitempty" bson:"errors,omitempty"`
		PluginPayloads   []payload.Payload `json:"payloads" bson:"payloads"`
	}

	// CommandResponseLog contains the client's response to a specific command and payload.
	CommandResponseLog struct {
		mgm.DefaultModel `bson:",inline"`
		PayloadId        string        `json:"payloadId" bson:"payloadId" validate:"required"`
		Errors           []interface{} `json:"errors,omitempty" bson:"errors,omitempty"`
		PluginResponse   interface{}   `json:"pluginResponse" bson:"pluginResponse"`
	}
)

func (c *CommandLog) CollectionName() string {
	return "command_logs"
}

func (c *CommandResponseLog) CollectionName() string {
	return "command_logs"
}
