package etfapi

import (
	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

const (
	administrator = 0x00000008
)

// Role TODOC
type Role struct {
	id          snowflake.Snowflake
	name        string
	permissions int64
}

// IsAdmin TODOC
func (r *Role) IsAdmin() bool {
	return r.permissions&administrator == administrator
}

// UpdateFromElementMap TODOC
func (r *Role) UpdateFromElementMap(eMap map[string]Element) (err error) {
	var ok bool
	var e2 Element

	e2, ok = eMap["name"]
	if ok {
		r.name, err = e2.ToString()
		if err != nil {
			err = errors.Wrap(err, "could not get name string")
			return
		}
	}

	e2, ok = eMap["permissions"]
	if ok {
		r.permissions, err = e2.ToInt64()
		if err != nil {
			err = errors.Wrap(err, "could not get permissions")
			return
		}
	}

	return
}

// RoleFromElement TODOC
func RoleFromElement(e Element) (r Role, err error) {
	var eMap map[string]Element

	eMap, r.id, err = MapAndIDFromElement(e)
	if err != nil {
		return
	}

	err = r.UpdateFromElementMap(eMap)
	err = errors.Wrap(err, "could not create a role")

	return
}
