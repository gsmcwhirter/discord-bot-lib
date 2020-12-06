package jsonapi

import "github.com/gsmcwhirter/discord-bot-lib/v15/snowflake"

//go:generate easyjson -all

// UserResponse is the data about a user recevied from the json api
//easyjson:json
type UserResponse struct {
	ID            snowflake.Snowflake  `json:"id"`
	Username      string               `json:"username"`
	Discriminator string               `json:"discriminator"`
	Avatar        string               `json:"avatar"`
	Bot           bool                 `json:"bot"`
	System        bool                 `json:"system"`
	MFAEnabled    bool                 `json:"mfa_enabled"`
	Locale        string               `json:"locale"`
	Verified      bool                 `json:"verified"`
	Email         string               `json:"email"`
	Flags         int                  `json:"flags"`
	PremiumType   int                  `json:"premium_type"`
	PublicFlags   int                  `json:"public_flags"`
	Member        *GuildMemberResponse `json:"member"`
}
