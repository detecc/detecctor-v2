package payload

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PayloadTestSuite struct {
	suite.Suite
}

func (suite *PayloadTestSuite) SetupTest() {
}

func (suite *PayloadTestSuite) TestEmpty() {
	expected := Payload{
		Id:             "",
		ServiceNodeKey: "",
		Data:           nil,
		Command:        "",
		Success:        false,
		Error:          "",
	}
	suite.Equal(expected, NewPayload())
}

func (suite *PayloadTestSuite) TestOk() {
	expected := Payload{
		Id:             "",
		ServiceNodeKey: "snKey123",
		Data:           "testData",
		Command:        "/cmd",
		Success:        true,
		Error:          "",
	}
	actual := NewPayload(ForClient("snKey123"), WithData("testData"), ForCommand("/cmd"), Successful())
	suite.Equal(expected, actual)
}

func (suite *PayloadTestSuite) TestWithPayloadError() {
	expected := Payload{
		Id:             "",
		ServiceNodeKey: "snKey123",
		Data:           "testData",
		Command:        "/cmd",
		Success:        false,
		Error:          "example error",
	}
	actual := NewPayload(ForClient("snKey123"), WithData("testData"), ForCommand("/cmd"), WithError(fmt.Errorf("example error")))
	suite.Equal(expected, actual)
}

func (suite *PayloadTestSuite) TestWithoutClient() {
	expected := Payload{
		Id:             "",
		ServiceNodeKey: "",
		Data:           "testData",
		Command:        "/cmd",
		Success:        false,
		Error:          "example error",
	}
	actual := NewPayload(WithData("testData"), ForCommand("/cmd"), WithError(fmt.Errorf("example error")))
	suite.Equal(expected, actual)
}

func TestNewPayload(t *testing.T) {
	suite.Run(t, new(PayloadTestSuite))
}

func TestPayload_SetError(t *testing.T) {
	expected := Payload{
		Id:             "",
		ServiceNodeKey: "snKey123",
		Data:           "testData",
		Command:        "/cmd",
		Success:        false,
		Error:          "example error 2",
	}
	payload := NewPayload(ForClient("snKey123"), WithData("testData"), ForCommand("/cmd"), Successful())

	payload.SetError(fmt.Errorf("example error 2"))
	assert.Equal(t, expected, payload)
}
