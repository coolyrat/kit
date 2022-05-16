package kitsqs

// type ActionReceiver struct {
// 	Client   *sqs.Client
// 	Wg       *sync.WaitGroup
// 	Pool     *ants.Pool
// 	handlers map[string]HandlersChain
//
// 	URL             string
// 	WaitTimeSeconds int32
// 	IdleWaitTime    time.Duration
// }
//
// func NewActionReceiver(conf *QueueConfig) *ActionReceiver {
// 	return &ActionReceiver{
// 		Client:   conf.Client,
// 		Wg:       conf.Wg,
// 		Pool:     conf.Pool,
// 		handlers: conf.Handlers,
// 		URL:      conf.URL,
// 		WaitTimeSeconds: conf.WaitTimeSeconds,
// 		IdleWaitTime:    conf.IdleWaitTime,
// 	}
// }
//
// func (r *ActionReceiver) ActionHandlers(action string, handlers ...HandlerFunc[any]) {
// 	r.handlers[action] = append(r.handlers[action], handlers...)
// }
//
// func (r *ActionReceiver) Receive(ctx context.Context) {
// 	fmt.Println("Subscribing to queue", r.URL)
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			goto exit
// 		default:
// 			resp, err := r.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
// 				QueueUrl:            aws.String(r.URL),
// 				MaxNumberOfMessages: 10,
// 				WaitTimeSeconds:     r.WaitTimeSeconds,
// 				AttributeNames: []types.QueueAttributeName{
// 					"All",
// 				},
// 				MessageAttributeNames: []string{
// 					"All",
// 				},
// 			})
// 			if err != nil {
// 				if !errors.Is(err, context.Canceled) {
// 					panic(err)
// 				} else {
// 					continue
// 				}
// 			}
//
// 			var isIdle bool
// 			if len(resp.Messages) == 0 {
// 				isIdle = true
// 			}
//
// 			for _, msg := range resp.Messages {
// 				// handle message
// 				// msg.
// 				// var m TopicMessage
// 				// if err := json.Unmarshal([]byte(*msg.Body), &m); err != nil {
// 				// 	panic(err)
// 				// }
// 				// m.receiptHandle = msg.ReceiptHandle
// 				// m.queueMsgID = msg.MessageId
//
// 				// if v, ok := m.MessageAttributes["action"]; ok {
// 				// 	// have an action field, pass it to the action handler
// 				// 	if chain, ok := r.actions[*v.StringValue]; ok {
// 				//
// 				// 		if err := h(ctx, &m); err != nil {
// 				// 			// TODO: handle error, should but to dead letter queue and save the reason?
// 				// 		}
// 				//
// 				// 		r.DeleteHandledMsg(m.receiptHandle)
// 				//
// 				// 	} else {
// 				// 		// TODO: no action handler, warn and put in dead letter queue
// 				// 		panic(fmt.Errorf("action %s not found", *v.StringValue))
// 				// 	}
// 				// } else {
// 				// 	// have no action, unknown message
// 				// }
//
// 				// spew.Dump(m)
// 				spew.Dump(msg)
// 			}
//
// 			if isIdle {
// 				time.Sleep(r.IdleWaitTime)
// 			}
// 		}
// 	}
// exit:
// 	fmt.Println(r.URL, " Subscriber stopped")
// 	r.Wg.Done()
//
// }
