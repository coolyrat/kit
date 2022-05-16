package kitsqs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/panjf2000/ants/v2"
)

type QueueSubscriber interface {
	Subscribe(ctx context.Context)
}

type QueueConfig struct {
	URL             string
	WaitTimeSeconds int32
	IdleWaitTime    time.Duration
}

type Queue struct {
	Client   *sqs.Client
	Wg       *sync.WaitGroup
	Pool     *ants.Pool
	handlers map[string]HandlerInfo

	URL             string
	WaitTimeSeconds int32
	IdleWaitTime    time.Duration
}

type HandlerInfo struct {
	action      string
	evtTyp      reflect.Type
	handlerFunc reflect.Value
}

func (q *Queue) ActionHandlers(action string, handlers ...any) {
	for _, h := range handlers {
		handlerTyp := reflect.TypeOf(h)
		handlerVal := reflect.ValueOf(h)
		if handlerTyp.Kind() != reflect.Func {
			panic(fmt.Sprintf("%s is not a function", handlerTyp.String()))
		}

		evtTyp := handlerTyp.In(1)
		q.handlers[action] = HandlerInfo{
			action:      action,
			evtTyp:      evtTyp,
			handlerFunc: handlerVal,
		}
	}
}

func (q *Queue) Receive(ctx context.Context) {
	fmt.Println("Subscribing to queue", q.URL)
	for {
		select {
		case <-ctx.Done():
			goto exit
		default:
			resp, err := q.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(q.URL),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     q.WaitTimeSeconds,
				AttributeNames: []types.QueueAttributeName{
					"All",
				},
				MessageAttributeNames: []string{
					"All",
				},
			})
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					panic(err)
				} else {
					continue
				}
			}

			var isIdle bool
			if len(resp.Messages) == 0 {
				isIdle = true
			}

			for _, msg := range resp.Messages {
				action, ok := actionMessage(msg.MessageAttributes)
				if !ok {
					// handle non-action message
				}

				handlerInfo, ok := q.handlers[action]
				if !ok {
					// handle unknown action
				}

				evt := reflect.New(handlerInfo.evtTyp).Interface()
				if err := json.Unmarshal([]byte(*msg.Body), evt); err != nil {

				}

				handlerInfo.handlerFunc.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(evt)})

				// handle message
				// msg.
				// var m TopicMessage
				// if err := json.Unmarshal([]byte(*msg.Body), &m); err != nil {
				// 	panic(err)
				// }
				// m.receiptHandle = msg.ReceiptHandle
				// m.queueMsgID = msg.MessageId

				// if v, ok := m.MessageAttributes["action"]; ok {
				// 	// have an action field, pass it to the action handler
				// 	if chain, ok := q.actions[*v.StringValue]; ok {
				//
				// 		if err := h(ctx, &m); err != nil {
				// 			// TODO: handle error, should but to dead letter queue and save the reason?
				// 		}
				//
				// 		q.DeleteHandledMsg(m.receiptHandle)
				//
				// 	} else {
				// 		// TODO: no action handler, warn and put in dead letter queue
				// 		panic(fmt.Errorf("action %s not found", *v.StringValue))
				// 	}
				// } else {
				// 	// have no action, unknown message
				// }

				// spew.Dump(m)
				spew.Dump(msg)
			}

			if isIdle {
				time.Sleep(q.IdleWaitTime)
			}
		}
	}
exit:
	fmt.Println(q.URL, " Subscriber stopped")
	q.Wg.Done()
}

// func (q *Queue) DeleteHandledMsg(receiptHandle *string) error {
// 	q.Client.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
// 		QueueUrl:      aws.String(q.URL),
// 		ReceiptHandle: receiptHandle,
// 	})
// }

type HandlerFunc[T any] func(ctx context.Context, msg T)

type HandlersChain []HandlerFunc[any]

// Last returns the last handler in the chain. ie. the last handler is the main one.
func (c HandlersChain) Last() HandlerFunc[any] {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

type ActionHandler struct {
	msg *types.Message
}
