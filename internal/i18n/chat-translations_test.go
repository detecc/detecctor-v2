package i18n

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TranslationMapSuite struct {
	suite.Suite
}

func (suite *TranslationMapSuite) SetupTest() {
}

func (suite *TranslationMapSuite) TestFullFunctionality() {
	expected := TranslationMap{
		MessageId: "messageId123",
		Plural:    1,
		Data:      map[string]interface{}{"key1": "val1"},
	}
	translationMap := NewTranslationMap(
		"messageId123",
		WithPlural(1),
		AddData("key1", "val1"),
	)

	suite.EqualValues(expected, translationMap)
}

func (suite *TranslationMapSuite) TestMultipleDataEntries() {
	expected := TranslationMap{
		MessageId: "messageId123",
		Plural:    1,
		Data:      map[string]interface{}{"key1": "val1", "key2": "val2"},
	}
	translationMap := NewTranslationMap(
		"messageId123",
		WithPlural(1),
		AddData("key1", "val1"),
		AddData("key2", "val2"),
	)

	suite.EqualValues(expected, translationMap)
}

func (suite *TranslationMapSuite) TestDefaults() {
	expected := TranslationMap{
		MessageId: "messageId123",
		Plural:    nil,
		Data:      map[string]interface{}{},
	}
	translationMap := NewTranslationMap("messageId123")

	suite.EqualValues(expected, translationMap)
}

func TestNewTranslationMap(t *testing.T) {
	suite.Run(t, new(TranslationMapSuite))
}
