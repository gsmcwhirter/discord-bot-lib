package session

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v21/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v21/snowflake"
)

const (
	administrator = 0x00000008
)

// Role represents a discord guild role
type Role struct {
	id          snowflake.Snowflake
	name        string
	permissions int64
}

// IsAdmin determines if a role is a server admin
func (r *Role) IsAdmin() bool {
	return r.permissions&administrator == administrator
}

// UpdateFromElementMap updates the data in a role from the given information
//
// This will not remove data, only change and add data
func (r *Role) UpdateFromElementMap(eMap map[string]etfapi.Element) error {
	var ok bool
	var e2 etfapi.Element

	var err error

	e2, ok = eMap["name"]
	if ok {
		r.name, err = e2.ToString()
		if err != nil {
			return errors.Wrap(err, "could not get name string")
		}
	}

	e2, ok = eMap["permissions"]
	if ok {
		r.permissions, err = e2.ToInt64()
		if err != nil {
			return errors.Wrap(err, "could not get permissions")
		}
	}

	return nil
}

// RoleFromElement generates a new Role object from the given Element
func RoleFromElement(e etfapi.Element) (Role, error) {
	var eMap map[string]etfapi.Element
	var r Role
	var err error

	eMap, r.id, err = etfapi.MapAndIDFromElement(e)
	if err != nil {
		return r, err
	}

	err = r.UpdateFromElementMap(eMap)
	return r, errors.Wrap(err, "could not create a role")
}
