package etfapi

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

// User TODOC
type User struct {
	id            snowflake.Snowflake
	username      string
	discriminator string
	avatar        string
}

// IDString TODOC
func (u *User) IDString() string {
	return fmt.Sprintf("<@!%s>", u.id.ToString())
}

// UpdateFromElementMap TODOC
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

// UserFromElement TODOC
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
