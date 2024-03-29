package session

import (
	"github.com/gsmcwhirter/go-util/v10/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
)

// User represents the data about a discord user
type User struct {
	id            snowflake.Snowflake
	username      string
	discriminator string
	avatar        string
}

// UpdateFromElementMap updates the information about a user from the given data
//
// This will not remove information, only change and add information
func (u *User) UpdateFromElementMap(eMap map[string]etfapi.Element) error {
	var e2 etfapi.Element
	var ok bool
	var err error

	if e2, ok = eMap["username"]; ok {
		u.username, err = e2.ToString()
		if err != nil {
			return errors.Wrap(err, "could not get username")
		}
	}

	if e2, ok = eMap["discriminator"]; ok {
		u.discriminator, err = e2.ToString()
		if err != nil {
			return errors.Wrap(err, "could not get discriminator")
		}
	}

	if e2, ok = eMap["avatar"]; ok {
		u.avatar, err = e2.ToString()
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

	eMap, u.id, err = etfapi.MapAndIDFromElement(e)
	if err != nil {
		return u, err
	}

	err = u.UpdateFromElementMap(eMap)
	return u, errors.Wrap(err, "could not inflate user data")
}
