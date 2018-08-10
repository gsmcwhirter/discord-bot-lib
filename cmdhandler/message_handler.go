package cmdhandler

// MessageHandler TODOC
type MessageHandler interface {
	HandleLine(Message) (Response, error)
}

type msgHandlerFunc struct {
	handler func(Message) (Response, error)
}

// NewMessageHandler TODOC
func NewMessageHandler(f func(Message) (Response, error)) MessageHandler {
	return &msgHandlerFunc{handler: f}
}

func (lh *msgHandlerFunc) HandleLine(msg Message) (Response, error) {
	return lh.handler(msg)
}
