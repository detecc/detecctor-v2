package mqtt

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type MqttTestSuite struct {
	suite.Suite
}

func (suite *MqttTestSuite) SetupTest() {
}

func (suite *MqttTestSuite) TestGetIdsFromTopic() {
	expectedIds := []string{"examplePlugin"}
	ids, err := GetIdsFromTopic("cmd/examplePlugin/execute", "cmd/+/execute")
	suite.Require().NoError(err)
	suite.Require().Equal(expectedIds, ids)

	ids, err = GetIdsFromTopic("cmd/execute", "cmd/+/execute")
	suite.Require().Error(err)

	ids, err = GetIdsFromTopic("ploogin/examplePlugin/execute", "cmd/+/execute")
	suite.Require().Error(err)

	ids, err = GetIdsFromTopic("ploogin/examplePlugin/execute", "cmd/execute")
	suite.Require().Error(err)

	ids, err = GetIdsFromTopic("cmd/examplePlugin/execute", "cmd/examplePlugin/execute")
	suite.Require().Error(err)

	ids, err = GetIdsFromTopic("cmd/examplePlugin/execute/example2/abc", "cmd/+/execute/+/abc")
	suite.Require().NoError(err)
	suite.Require().Equal([]string{"examplePlugin", "example2"}, ids)
}

func (suite *MqttTestSuite) TestCreateTopicWithIds() {
	ids, err := CreateTopicWithIds("cmd/+/execute", "exampleId")
	suite.Require().NoError(err)
	suite.Require().Equal("cmd/exampleId/execute", ids)

	ids, err = CreateTopicWithIds("cmd/+/execute/+/", "exampleId1", "exampleId2")
	suite.Require().NoError(err)
	suite.Require().Equal("cmd/exampleId1/execute/exampleId1/", ids)

	ids, err = CreateTopicWithIds("cmd/+/execute/+/", "exampleId")
	suite.Require().Error(err)

	ids, err = CreateTopicWithIds("cmd/+/execute/+/", "exampleId", "")
	suite.Require().Error(err)
}

func TestGetIdsFromTopic(t *testing.T) {
	suite.Run(t, new(MqttTestSuite))
}
