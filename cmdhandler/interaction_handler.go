package cmdhandler

// InteractionHandlerFunc is the api of a function that handles messages
type InteractionHandlerFunc func(Interaction) (Response, error)

// InteractionHandler is the api of a message handler
type InteractionHandler interface {
	HandleInteraction(Interaction) (Response, error)
}

type interactionHandlerFunc struct {
	handler func(Interaction) (Response, error)
}

// NewInteractionHandler wraps a InteractionHandlerFunc into a InteractionHandler
func NewInteractionHandler(f InteractionHandlerFunc) InteractionHandler {
	return &interactionHandlerFunc{handler: f}
}

func (lh *interactionHandlerFunc) HandleInteraction(msg Interaction) (Response, error) {
	return lh.handler(msg)
}
