package cmdhandler

import "github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/entity"

// InteractionHandlerFunc is the api of a function that handles messages
type InteractionHandlerFunc func(*Interaction) (Response, []Response, error)

// InteractionHandler is the api of a message handler
type InteractionHandler interface {
	HandleInteraction(*Interaction) (Response, []Response, error)
}

type interactionHandlerFunc struct {
	handler func(*Interaction) (Response, []Response, error)
}

// NewInteractionHandler wraps a InteractionHandlerFunc into a InteractionHandler
func NewInteractionHandler(f InteractionHandlerFunc) InteractionHandler {
	return &interactionHandlerFunc{handler: f}
}

func (lh *interactionHandlerFunc) HandleInteraction(ix *Interaction) (Response, []Response, error) {
	return lh.handler(ix)
}

// AutocompleteHandlerFunc is the api of a function that handles autocomplete requests
type AutocompleteHandlerFunc func(*Interaction) ([]entity.ApplicationCommandOptionChoice, error)

// AutocompleteHandler is the api of an autocomplete handler
type AutocompleteHandler interface {
	Autocomplete(*Interaction) ([]entity.ApplicationCommandOptionChoice, error)
}

type autocompleteHandlerFunc struct {
	handler func(*Interaction) ([]entity.ApplicationCommandOptionChoice, error)
}

// NewAutocompleteHandler wraps an AutocompleteHandlerFunc into a AutocompleteHandler
func NewAutocompleteHandler(f AutocompleteHandlerFunc) AutocompleteHandler {
	return &autocompleteHandlerFunc{handler: f}
}

func (ah *autocompleteHandlerFunc) Autocomplete(ix *Interaction) ([]entity.ApplicationCommandOptionChoice, error) {
	return ah.handler(ix)
}
