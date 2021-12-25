package logs

import (
	"fmt"
	builder "github.com/detecc/detecctor-v2/internal/command"
	"github.com/detecc/detecctor-v2/internal/model/command/logs"
	"github.com/detecc/detecctor-v2/pkg/payload"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CommandLogTestSuite struct {
	suite.Suite
}

func (suite *CommandLogTestSuite) SetupTest() {
}

func (suite *CommandLogTestSuite) TestNewCommandLog() {
	cmd := builder.NewCommandBuilder().Id("exampleId").FromChat("exampleChatId").WithName("/example").Build()
	payloads := []payload.Payload{{
		Id:             "1234",
		ServiceNodeKey: "example123",
		Data:           nil,
		Command:        "/example",
		Success:        true,
		Error:          "",
	}}

	expectedLog := logs.CommandLog{
		Command:        cmd,
		Errors:         nil,
		PluginPayloads: payloads,
	}

	cmdLog := NewCommandLog(cmd, WithErrors(nil), WithPayloads(payloads...))

	suite.Require().EqualValues(expectedLog.Command, cmdLog.Command)
	suite.Require().EqualValues(expectedLog.Errors, cmdLog.Errors)
	suite.Require().EqualValues(expectedLog.PluginPayloads, cmdLog.PluginPayloads)

	cmdLog = NewCommandLog(cmd, WithErrors(fmt.Errorf("sample123"), nil), WithPayloads(payloads...))

	suite.Require().EqualValues(expectedLog.Command, cmdLog.Command)
	suite.Require().EqualValues([]interface{}{"sample123"}, cmdLog.Errors)
	suite.Require().EqualValues(expectedLog.PluginPayloads, cmdLog.PluginPayloads)

	cmdLog = NewCommandLog(cmd, WithErrors(fmt.Errorf("sample123"), nil))

	suite.Require().EqualValues(expectedLog.Command, cmdLog.Command)
	suite.Require().EqualValues([]interface{}{"sample123"}, cmdLog.Errors)
	suite.Require().Nil(cmdLog.PluginPayloads)

}

func (suite *CommandLogTestSuite) TestNewCommandResponseLog() {
	expectedResponseLog := logs.CommandResponseLog{
		PayloadId:      "examplePayloadId",
		Errors:         nil,
		PluginResponse: nil,
	}

	responseLog := NewCommandResponseLog("examplePayloadId")
	suite.Require().EqualValues(expectedResponseLog.PayloadId, responseLog.PayloadId)
	suite.Require().EqualValues(expectedResponseLog.Errors, responseLog.Errors)
	suite.Require().EqualValues(expectedResponseLog.PluginResponse, responseLog.PluginResponse)

	responseLog = NewCommandResponseLog("examplePayloadId", WithResponse("sampleResponse123"))
	suite.Require().EqualValues(expectedResponseLog.PayloadId, responseLog.PayloadId)
	suite.Require().EqualValues(expectedResponseLog.Errors, responseLog.Errors)
	suite.Require().EqualValues("sampleResponse123", responseLog.PluginResponse)

	responseLog = NewCommandResponseLog("examplePayloadId", WithResponseError(fmt.Errorf("exampleErr")))
	suite.Require().EqualValues(expectedResponseLog.PayloadId, responseLog.PayloadId)
	suite.Require().EqualValues([]interface{}{"exampleErr"}, responseLog.Errors)
	suite.Require().EqualValues(expectedResponseLog.PluginResponse, responseLog.PluginResponse)

	responseLog = NewCommandResponseLog("examplePayloadId", WithResponse("sampleResponse123"), WithResponseError(fmt.Errorf("exampleErr")))
	suite.Require().EqualValues(expectedResponseLog.PayloadId, responseLog.PayloadId)
	suite.Require().EqualValues([]interface{}{"exampleErr"}, responseLog.Errors)
	suite.Require().EqualValues("sampleResponse123", responseLog.PluginResponse)
}

func TestCommandBuilder(t *testing.T) {
	suite.Run(t, new(CommandLogTestSuite))
}
