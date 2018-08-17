package jsonapi

//go:generate easyjson -all -snake_case $GOPATH/src/github.com/gsmcwhirter/discord-bot-lib/discordapi/jsonapi/gateway_response.go

// GatewayResponse is the json object received from the discord api
// when requesting gateway connection information
//easyjson:json
type GatewayResponse struct {
	URL    string
	Shards int
}
