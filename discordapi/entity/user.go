package entity

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v21/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v21/snowflake"
)

// User is the data about a user recevied from the json api
type User struct {
	ID            string       `json:"id"`
	Username      string       `json:"username"`
	Discriminator string       `json:"discriminator"`
	Avatar        string       `json:"avatar"`
	Bot           bool         `json:"bot"`
	System        bool         `json:"system"`
	MFAEnabled    bool         `json:"mfa_enabled"`
	Locale        string       `json:"locale"`
	Verified      bool         `json:"verified"`
	Email         string       `json:"email"`
	Flags         int          `json:"flags"`
	PremiumType   int          `json:"premium_type"`
	PublicFlags   int          `json:"public_flags"`
	Member        *GuildMember `json:"member"`

	IDSnowflake snowflake.Snowflake
}

func (u *User) Snowflakify() error {
	var err error

	if u.ID != "" {
		if u.IDSnowflake, err = snowflake.FromString(u.ID); err != nil {
			return errors.Wrap(err, "could not snowflakify ID")
		}
	}

	if u.Member != nil {
		if err = u.Member.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify Member")
		}
	}

	return nil
}

// UpdateFromElementMap updates the information about a user from the given data
//
// This will not remove information, only change and add information
func (u *User) UpdateFromElementMap(eMap map[string]etfapi.Element) error {
	var e2 etfapi.Element
	var ok bool
	var err error

	if e2, ok = eMap["username"]; ok {
		u.Username, err = e2.ToString()
		if err != nil {
			return errors.Wrap(err, "could not get username")
		}
	}

	if e2, ok = eMap["discriminator"]; ok {
		u.Discriminator, err = e2.ToString()
		if err != nil {
			return errors.Wrap(err, "could not get discriminator")
		}
	}

	if e2, ok = eMap["avatar"]; ok {
		u.Avatar, err = e2.ToString()
		if err != nil {
			return errors.Wrap(err, "could not get avatar")
		}
	}

	return nil
}

// UserFromElement generates a new User object from the given Element
func UserFromElement(e etfapi.Element) (User, error) {
	var u User
	var eMap map[string]etfapi.Element
	var err error

	eMap, u.IDSnowflake, err = etfapi.MapAndIDFromElement(e)
	if err != nil {
		return u, err
	}

	u.ID = u.IDSnowflake.ToString()

	err = u.UpdateFromElementMap(eMap)
	return u, errors.Wrap(err, "could not inflate user data")
}
