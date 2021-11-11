package entity

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v22/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v22/snowflake"
)

// Role is the data about a role recevied from the json api
type Role struct {
	IDString    string `json:"id"`
	Name        string `json:"name"`
	Color       int    `json:"color"`
	Hoist       bool   `json:"hoist"`
	Position    int    `json:"position"`
	Permissions int64  `json:"permissions"`
	Managed     bool   `json:"managed"`
	Mentionable bool   `json:"mentionable"`

	IDSnowflake snowflake.Snowflake `json:"-"`
}

func (rr *Role) Snowflakify() error {
	var err error

	if rr.IDSnowflake, err = snowflake.FromString(rr.IDString); err != nil {
		return errors.Wrap(err, "could not snowflakify ID")
	}

	return nil
}

// RoleFromElement generates a new Role object from the given data
func RoleFromElement(e etfapi.Element) (Role, error) {
	var r Role

	eMap, err := e.ToMap()
	if err != nil {
		return r, errors.Wrap(err, "could not inflate Role from non-map")
	}

	e2, ok := eMap["id"]
	if ok {
		r.IDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return r, errors.Wrap(err, "could not get id snowflake.Snowflake")
		}

		r.IDString = r.IDSnowflake.ToString()
	}

	e2 = eMap["name"]
	r.Name, err = e2.ToString()
	if err != nil {
		return r, errors.Wrap(err, "could not inflate name")
	}

	e2 = eMap["color"]
	r.Color, err = e2.ToInt()
	if err != nil {
		return r, errors.Wrap(err, "could not inflate color")
	}

	e2 = eMap["hoist"]
	r.Hoist, err = e2.ToBool()
	if err != nil {
		return r, errors.Wrap(err, "could not inflate hoist")
	}

	e2 = eMap["position"]
	r.Position, err = e2.ToInt()
	if err != nil {
		return r, errors.Wrap(err, "could not inflate position")
	}

	e2 = eMap["permissions"]
	r.Permissions, err = e2.ToInt64()
	if err != nil {
		return r, errors.Wrap(err, "could not inflate permissions")
	}

	e2 = eMap["managed"]
	r.Managed, err = e2.ToBool()
	if err != nil {
		return r, errors.Wrap(err, "could not inflate managed")
	}

	e2 = eMap["mentionable"]
	r.Mentionable, err = e2.ToBool()
	if err != nil {
		return r, errors.Wrap(err, "could not inflate mentionable")
	}

	return r, nil
}
