package etfapi

import (
	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
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
func (u *User) UpdateFromElementMap(eMap map[string]Element) (err error) {
	var e2 Element
	var ok bool

	if e2, ok = eMap["username"]; ok {
		u.username, err = e2.ToString()
		if err != nil {
			err = errors.Wrap(err, "could not get username")
			return
		}
	}

	if e2, ok = eMap["discriminator"]; ok {
		u.discriminator, err = e2.ToString()
		if err != nil {
			err = errors.Wrap(err, "could not get discriminator")
			return
		}
	}

	if e2, ok = eMap["avatar"]; ok {
		u.avatar, err = e2.ToString()
		if err != nil {
			err = errors.Wrap(err, "could not get avatar")
			return
		}
	}

	return
}

// UserFromElement generates a new User object from the given Element
func UserFromElement(e Element) (u User, err error) {
	var eMap map[string]Element
	eMap, u.id, err = MapAndIDFromElement(e)
	if err != nil {
		return
	}

	err = u.UpdateFromElementMap(eMap)
	err = errors.Wrap(err, "could not inflate user data")
	return
}
