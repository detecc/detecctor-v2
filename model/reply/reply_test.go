package reply

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ReplyBuilderTestSuite struct {
	suite.Suite
	builder *Builder
}

func (suite *ReplyBuilderTestSuite) SetupTest() {
	suite.builder = NewReplyBuilder()
}

func (suite *ReplyBuilderTestSuite) TestEmptyBuild() {
	expected := Reply{
		ChatId:    "0",
		ReplyType: -1,
		Content:   nil,
	}
	suite.Equal(expected, suite.builder.Build())
}

func (suite *ReplyBuilderTestSuite) TestLegitBuild() {
	expected := Reply{
		ChatId:    "chatId123",
		ReplyType: 0,
		Content:   "123",
	}
	reply := suite.builder.TypeMessage().ForChat("chatId123").WithContent("123").Build()
	suite.Equal(expected, reply)
}

func (suite *ReplyBuilderTestSuite) TestBuildWithoutOneAttribute() {
	expected := Reply{
		ChatId:    "0",
		ReplyType: 0,
		Content:   "123",
	}
	reply := suite.builder.TypeMessage().WithContent("123").Build()
	suite.Equal(expected, reply)
}

func TestReplyBuilder(t *testing.T) {
	suite.Run(t, new(ReplyBuilderTestSuite))
}
