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

// UserFromElement TODOC
func UserFromElement(e Element) (u User, err error) {
	var e2 Element

	var eMap map[string]Element
	eMap, u.id, err = MapAndIDFromElement(e)
	if err != nil {
		return
	}

	e2 = eMap["username"]
	u.username, err = e2.ToString()
	if err != nil {
		err = errors.Wrap(err, "could not get username")
		return
	}

	e2 = eMap["discriminator"]
	u.discriminator, err = e2.ToString()
	if err != nil {
		err = errors.Wrap(err, "could not get discriminator")
		return
	}

	e2 = eMap["avatar"]
	u.avatar, err = e2.ToString()
	if err != nil {
		err = errors.Wrap(err, "could not get avatar")
		return
	}

	return
}
