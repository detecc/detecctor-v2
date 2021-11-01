package logs

import (
	"github.com/detecc/detecctor-v2/model/command"
	"github.com/detecc/detecctor-v2/model/payload"
	"github.com/kamva/mgm/v3"
)

type (
	// CommandLog contains information about a Command that was issued to the server, processing of the command
	// and the payloads produced by the plugin as well as any errors that occurred.
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

	Option func(cmd *CommandLog)

	ResponseOption func(cmd *CommandResponseLog)
)

func (c *CommandLog) CollectionName() string {
	return "command_logs"
}

func (c *CommandResponseLog) CollectionName() string {
	return "command_logs"
}

func WithPayloads(payloads ...payload.Payload) Option {
	return func(cmd *CommandLog) {
		cmd.PluginPayloads = payloads
	}
}

func WithResponse(response interface{}) ResponseOption {
	return func(cmd *CommandResponseLog) {
		cmd.PluginResponse = response
	}
}

func WithErrors(errors ...error) Option {
	return func(cmd *CommandLog) {
		for _, err := range errors {
			if err != nil {
				cmd.Errors = append(cmd.Errors, err.Error())
			}
		}
	}
}

func WithResponseError(errors ...error) ResponseOption {
	return func(cmd *CommandResponseLog) {
		for _, err := range errors {
			if err != nil {
				cmd.Errors = append(cmd.Errors, err.Error())
			}
		}
	}
}

//NewCommandLog creates a new command log
func NewCommandLog(command command.Command, options ...Option) *CommandLog {
	cmdLog := &CommandLog{
		Command:        command,
		Errors:         nil,
		PluginPayloads: nil,
	}

	for _, option := range options {
		option(cmdLog)
	}

	return cmdLog
}

//NewCommandResponseLog creates a new response log for a command
func NewCommandResponseLog(payloadId string, options ...ResponseOption) *CommandResponseLog {
	responseLog := &CommandResponseLog{
		PayloadId:      payloadId,
		Errors:         nil,
		PluginResponse: nil,
	}

	for _, option := range options {
		option(responseLog)
	}

	return responseLog
}
