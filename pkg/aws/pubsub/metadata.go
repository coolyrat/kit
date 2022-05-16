package kitsqs

import "github.com/aws/aws-sdk-go-v2/service/sqs/types"

func actionMessage(attrs map[string]types.MessageAttributeValue) (string, bool) {
	if attrs == nil {
		return "", false
	}

	action, ok := attrs["action"]
	if !ok {
		return "", false
	}
	return *action.StringValue, true
}
