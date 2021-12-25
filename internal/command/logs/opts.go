package logs

import (
	"github.com/detecc/detecctor-v2/internal/model/command"
	. "github.com/detecc/detecctor-v2/internal/model/command/logs"
	"github.com/detecc/detecctor-v2/pkg/payload"
)

type (
	Option func(cmd *CommandLog)

	ResponseOption func(cmd *CommandResponseLog)
)

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
