package cmdhandler

// InteractionHandlerFunc is the api of a function that handles messages
type InteractionHandlerFunc func(Interaction) (Response, []Response, error)

// InteractionHandler is the api of a message handler
type InteractionHandler interface {
	HandleInteraction(Interaction) (Response, []Response, error)
}

type interactionHandlerFunc struct {
	handler func(Interaction) (Response, []Response, error)
}

// NewInteractionHandler wraps a InteractionHandlerFunc into a InteractionHandler
func NewInteractionHandler(f InteractionHandlerFunc) InteractionHandler {
	return &interactionHandlerFunc{handler: f}
}

func (lh *interactionHandlerFunc) HandleInteraction(msg Interaction) (Response, []Response, error) {
	return lh.handler(msg)
}
