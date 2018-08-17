package cmdhandler

// MessageHandlerFunc is the api of a function that handles messages
type MessageHandlerFunc func(Message) (Response, error)

// MessageHandler is the api of a message handler
type MessageHandler interface {
	HandleMessage(Message) (Response, error)
}

type msgHandlerFunc struct {
	handler func(Message) (Response, error)
}

// NewMessageHandler wraps a MessageHandlerFunc into a MessageHandler
func NewMessageHandler(f MessageHandlerFunc) MessageHandler {
	return &msgHandlerFunc{handler: f}
}

func (lh *msgHandlerFunc) HandleMessage(msg Message) (Response, error) {
	return lh.handler(msg)
}
