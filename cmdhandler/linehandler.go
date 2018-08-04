package cmdhandler

// LineHandler TODOC
type LineHandler interface {
	HandleLine(user, guild string, line string) (Response, error)
}

type lineHandlerFunc struct {
	handler func(user, guild string, line string) (Response, error)
}

// NewLineHandler TODOC
func NewLineHandler(f func(string, string, string) (Response, error)) LineHandler {
	return &lineHandlerFunc{handler: f}
}

func (lh *lineHandlerFunc) HandleLine(user, guild string, line string) (Response, error) {
	return lh.handler(user, guild, line)
}
