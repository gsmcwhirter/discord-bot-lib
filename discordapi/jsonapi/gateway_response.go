package jsonapi

//go:generate easyjson -all -snake_case $GOPATH/src/github.com/gsmcwhirter/discord-bot-lib/discordapi/jsonapi/gateway_response.go

// GatewayResponse TODOC
//easyjson:json
type GatewayResponse struct {
	URL    string
	Shards int
}
