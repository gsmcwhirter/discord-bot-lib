package jsonapi

// GatewayResponse is the json object received from the discord api
// when requesting gateway connection information
type GatewayResponse struct {
	URL    string `json:"url"`
	Shards int    `json:"shards"`
}
