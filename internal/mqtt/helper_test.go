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
	expectedIds := []string{"examplePloogin"}
	ids, err := GetIdsFromTopic("plugin/examplePloogin/execute", "plugin/+/execute")
	suite.Require().NoError(err)
	suite.Require().Equal(expectedIds, ids)

	ids, err = GetIdsFromTopic("plugin/execute", "plugin/+/execute")
	suite.Require().Error(err)

	ids, err = GetIdsFromTopic("ploogin/examplePloogin/execute", "plugin/+/execute")
	suite.Require().Error(err)

	ids, err = GetIdsFromTopic("ploogin/examplePloogin/execute", "plugin/execute")
	suite.Require().Error(err)

	ids, err = GetIdsFromTopic("plugin/examplePloogin/execute", "plugin/examplePloogin/execute")
	suite.Require().Error(err)

	ids, err = GetIdsFromTopic("plugin/examplePloogin/execute/example2/abc", "plugin/+/execute/+/abc")
	suite.Require().NoError(err)
	suite.Require().Equal([]string{"examplePloogin", "example2"}, ids)
}

func (suite *MqttTestSuite) TestCreateTopicWithIds() {
	ids, err := CreateTopicWithIds("plugin/+/execute", "exampleId")
	suite.Require().NoError(err)
	suite.Require().Equal("plugin/exampleId/execute", ids)

	ids, err = CreateTopicWithIds("plugin/+/execute/+/", "exampleId1", "exampleId2")
	suite.Require().NoError(err)
	suite.Require().Equal("plugin/exampleId1/execute/exampleId1/", ids)

	ids, err = CreateTopicWithIds("plugin/+/execute/+/", "exampleId")
	suite.Require().Error(err)

	ids, err = CreateTopicWithIds("plugin/+/execute/+/", "exampleId", "")
	suite.Require().Error(err)
}

func TestGetIdsFromTopic(t *testing.T) {
	suite.Run(t, new(MqttTestSuite))
}
