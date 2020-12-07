package jsonapi

import (
	"github.com/gsmcwhirter/go-util/v7/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v17/snowflake"
)

//go:generate easyjson -all

// UserResponse is the data about a user recevied from the json api
//easyjson:json
type UserResponse struct {
	ID            string               `json:"id"`
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

	IDSnowflake snowflake.Snowflake
}

func (ur *UserResponse) Snowflakify() error {
	var err error

	if ur.IDSnowflake, err = snowflake.FromString(ur.ID); err != nil {
		return errors.Wrap(err, "could not snowflakify ID")
	}

	if ur.Member != nil {
		if err = ur.Member.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify Member")
		}
	}

	return nil
}
