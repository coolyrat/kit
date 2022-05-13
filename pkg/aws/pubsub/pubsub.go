package kitsqs

import "context"

type Publisher interface {
	// Publish sends a message to a topic.
	// TODO:
	// 	MsgAttr could not more than 10
	// 	builder pattern contains:
	// 	Message(m any), required
	// 	Topic(topic string), if there's a preset topic, then it's not required, this will override the default topic
	// 	To(app ...string), if not set, will fan out to all apps
	// 	AddAttr(key string, v *MessageAttributeValue)
	// 	Build() will build a PublishInput
	Publish(ctx context.Context, b SNSMessageBuilder) error
}

type QueueSender interface {
	// Send sends a message to a queue.
	// TODO:
	// 	builder pattern contains:
	// 	Message(m any), required
	// 	Queue(queue string), if there's a preset queue, then it's not required, this will override the default queue
	// 	DelaySec(s int64), if a message requires delay
	// 	AddAttr(key string, v *MessageAttributeValue)
	// 	Build() will build a SendMessageInput
	Send(ctx context.Context, b SQSMessageBuilder) error
}

type ActionQueueReceiver interface {
	AddActionHandlers(action string, handlers ...HandlerFunc)
}

type SNSMessageBuilder struct{}

type SQSMessageBuilder struct{}
