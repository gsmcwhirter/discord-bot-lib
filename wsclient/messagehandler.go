package wsclient

// MessageHandlerFunc TODOC
type MessageHandlerFunc func(WSMessage, chan<- WSMessage)

// MessageHandler TODOC
type MessageHandler interface {
	HandleRequest(WSMessage, chan<- WSMessage)
}

type messageHandler struct {
	handler MessageHandlerFunc
}

// NewMessageHandler TODOC
func NewMessageHandler(h MessageHandlerFunc) MessageHandler {
	return messageHandler{handler: h}
}

func (mh messageHandler) HandleRequest(req WSMessage, resp chan<- WSMessage) {
	mh.handler(req, resp)
}
