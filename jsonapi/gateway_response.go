package jsonapi

//go:generate easyjson -all

// GatewayResponse is the json object received from the discord api
// when requesting gateway connection information
//easyjson:json
type GatewayResponse struct {
	URL    string `json:"url"`
	Shards int    `json:"shards"`
}
