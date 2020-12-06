package jsonapi

import "github.com/gsmcwhirter/discord-bot-lib/v13/snowflake"

//go:generate easyjson -all

// GuildMemberResponse is the data about a guild member recevied from the json api
//easyjson:json
type GuildMemberResponse struct {
	User         *UserResponse         `json:"user"`
	Nick         string                `json:"nick"`
	Roles        []snowflake.Snowflake `json:"roles"`
	JoinedAt     string                `json:"joined_at"`     // ISO8601
	PremiumSince string                `json:"premium_since"` //ISO8601
	Deaf         bool                  `json:"deaf"`
	Mute         bool                  `json:"mute"`
}

func (gmr GuildMemberResponse) HasRole(rid snowflake.Snowflake) bool {
	for _, r := range gmr.Roles {
		if r == rid {
			return true
		}
	}

	return false
}
