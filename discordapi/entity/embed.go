package entity

// Embed is the data about a message embed received from the json api
type Embed struct {
	Title       string        `json:"title"`
	Type        string        `json:"type"`
	Description string        `json:"description"`
	URL         string        `json:"url"`
	Timestamp   string        `json:"timestamp"` // ISO8601
	Color       int           `json:"color"`
	Footer      EmbedFooter   `json:"footer"`
	Image       EmbedImage    `json:"image"`
	Thumbnail   EmbedImage    `json:"thumbnail"`
	Video       EmbedImage    `json:"video"`
	Provider    EmbedProvider `json:"provider"`
	Author      EmbedAuthor   `json:"author"`
	Fields      []EmbedField  `json:"fields"`
}

// EmbedFooter is the data about an embed footer recevied from the json api
type EmbedFooter struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url"`
	ProxyIconURL string `json:"proxy_icon_url"`
}

// EmbedField is the data about an embed field received from the json api
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// EmbedImage is the data about an embed thumbnail received from the json api
type EmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

// EmbedProvider is the data about an embed provider recevied from the json api
type EmbedProvider struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// EmbedAuthor is the data about an embed author recevied from the json api
type EmbedAuthor struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	IconURL      string `json:"icon_url"`
	ProxyIconURL string `json:"proxy_icon_url"`
}
