package jsonapi

//go:generate easyjson -all -snake_case $GOPATH/src/github.com/gsmcwhirter/discord-bot-lib/discordapi/jsonapi/message.go

// Message TODOC
//easyjson:json
type Message struct {
	Content string
	Tts     bool
}
