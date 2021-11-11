package entity

// Gateway is the json object received from the discord api
// when requesting gateway connection information
type Gateway struct {
	URL    string `json:"url"`
	Shards int    `json:"shards"`
}
