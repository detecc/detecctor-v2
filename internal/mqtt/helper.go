package mqtt

import (
	"fmt"
	"github.com/agrison/go-commons-lang/stringUtils"
	"strings"
)

// GetIdsFromTopic parses the topic received from the MQTT client and returns the ids based on the original subscription topic.
// For example:
// actual topic = some/exampleId1/subscription/exampleId2/topic
// subscription topic = some/+/subscription/+/topic
// should return ["exampleId1", "exampleId2"]
// If the topic are not the same length or don't contain the same words, it will return an error
func GetIdsFromTopic(actualTopic string, subTopic Topic) ([]string, error) {
	var ids []string

	actualTopicValues := strings.Split(actualTopic, "/")
	subTopicValues := strings.Split(string(subTopic), "/")

	// check if it is the same length, which would indicate the same topic
	if len(actualTopicValues) != len(subTopicValues) {
		return nil, fmt.Errorf("not the same topic")
	}

	//check if the subscription topic has at least one + or #
	if !strings.ContainsAny(string(subTopic), "+#") {
		return nil, fmt.Errorf("not a valid subscription topic")
	}

	for i, value := range subTopicValues {
		if value != actualTopicValues[i] && value != "+" {
			return nil, fmt.Errorf("not the subscribed topic")
		} else if value == "+" {
			ids = append(ids, actualTopicValues[i])
		}
	}

	return ids, nil
}

//CreateTopicWithIds replaces all the + sign in a topic used for subscription with ids. Works only if the number of pluses is matches the number of ids.
func CreateTopicWithIds(topicTemplate Topic, ids ...string) (string, error) {
	var finalString = string(topicTemplate)

	// check if the number of pluses is the same
	if strings.Count(string(topicTemplate), "+") != len(ids) {
		return "", fmt.Errorf("invalid number of arguments")
	}

	// any empty string is invalid
	if stringUtils.IsAnyEmpty(ids...) {
		return "", fmt.Errorf("ids cannot be an empty string")
	}

	for _, id := range ids {
		finalString = strings.Replace(finalString, "+", id, 1)
	}

	return finalString, nil
}
