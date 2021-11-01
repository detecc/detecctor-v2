package command

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type CommandBuilderTestSuite struct {
	suite.Suite
	builder *Builder
}

func (suite *CommandBuilderTestSuite) SetupTest() {
	suite.builder = NewCommandBuilder()
}

func (suite *CommandBuilderTestSuite) TestEmptyBuild() {
	expectedCommand := Command{
		Name:   "",
		Args:   []string{},
		ChatId: "0",
	}
	suite.Require().Equal(expectedCommand, suite.builder.Build())
}

func (suite *CommandBuilderTestSuite) TestLegitBuild() {
	expectedCommand := Command{
		Name:   "/name123",
		Args:   []string{"a", "b"},
		ChatId: "chatId123",
	}
	cmd := suite.builder.WithName("name123").WithArgs([]string{"a", "b"}).FromChat("chatId123").Build()
	suite.Equal(expectedCommand, cmd)
}

func (suite *CommandBuilderTestSuite) TestBuildWithoutArgs() {
	expectedCommand := Command{
		Name:   "/name123",
		Args:   []string{},
		ChatId: "chatId123",
	}

	cmd := suite.builder.WithName("name123").FromChat("chatId123").Build()
	suite.Require().Equal(expectedCommand, cmd)
}

func TestCommandBuilder(t *testing.T) {
	suite.Run(t, new(CommandBuilderTestSuite))
}
