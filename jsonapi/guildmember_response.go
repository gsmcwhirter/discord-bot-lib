package jsonapi

import (
	"github.com/gsmcwhirter/go-util/v7/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v16/snowflake"
)

//go:generate easyjson -all

// GuildMemberResponse is the data about a guild member recevied from the json api
//easyjson:json
type GuildMemberResponse struct {
	User         *UserResponse `json:"user"`
	Nick         string        `json:"nick"`
	Roles        []string      `json:"roles"`
	JoinedAt     string        `json:"joined_at"`     // ISO8601
	PremiumSince string        `json:"premium_since"` // ISO8601
	Deaf         bool          `json:"deaf"`
	Mute         bool          `json:"mute"`

	RoleSnowflakes []snowflake.Snowflake
}

func (gmr *GuildMemberResponse) Snowflakify() error {
	sfs := make([]snowflake.Snowflake, 0, len(gmr.Roles))
	for _, r := range gmr.Roles {
		sf, err := snowflake.FromString(r)
		if err != nil {
			return errors.Wrap(err, "could not convert role strings to snowflakes")
		}
		sfs = append(sfs, sf)
	}

	gmr.RoleSnowflakes = sfs

	return nil
}

func (gmr GuildMemberResponse) HasRole(rid snowflake.Snowflake) bool {
	for _, r := range gmr.RoleSnowflakes {
		if r == rid {
			return true
		}
	}

	return false
}
