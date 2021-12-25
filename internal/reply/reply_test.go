package reply

import (
	"github.com/detecc/detecctor-v2/internal/model/reply"
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
	expected := reply.Reply{
		ChatId:    "0",
		ReplyType: -1,
		Content:   nil,
	}
	suite.Equal(expected, suite.builder.Build())
}

func (suite *ReplyBuilderTestSuite) TestLegitBuild() {
	expected := reply.Reply{
		ChatId:    "chatId123",
		ReplyType: 0,
		Content:   "123",
	}
	r := suite.builder.TypeMessage().ForChat("chatId123").WithContent("123").Build()
	suite.Equal(expected, r)
}

func (suite *ReplyBuilderTestSuite) TestBuildWithoutOneAttribute() {
	expected := reply.Reply{
		ChatId:    "0",
		ReplyType: 0,
		Content:   "123",
	}
	r := suite.builder.TypeMessage().WithContent("123").Build()
	suite.Equal(expected, r)
}

func TestReplyBuilder(t *testing.T) {
	suite.Run(t, new(ReplyBuilderTestSuite))
}
