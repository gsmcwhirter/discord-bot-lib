package entity

import (
	"github.com/gsmcwhirter/go-util/v10/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
)

// Role is the data about a role recevied from the json api
type Role struct {
	IDString          string `json:"id"`
	Name              string `json:"name"`
	Color             int    `json:"color"`
	Hoist             bool   `json:"hoist"`
	Position          int    `json:"position"`
	PermissionsString string `json:"permissions"`
	Managed           bool   `json:"managed"`
	Mentionable       bool   `json:"mentionable"`

	IDSnowflake          snowflake.Snowflake `json:"-"`
	PermissionsSnowflake snowflake.Snowflake `json:"-"`
}

// Snowflakify converts snowflake strings into real sowflakes
func (rr *Role) Snowflakify() error {
	var err error

	if rr.IDSnowflake, err = snowflake.FromString(rr.IDString); err != nil {
		return errors.Wrap(err, "could not snowflakify ID")
	}

	if rr.PermissionsSnowflake, err = snowflake.FromString(rr.PermissionsString); err != nil {
		return errors.Wrap(err, "could not snowflakify Permissions")
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
		r.IDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
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
	r.PermissionsSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
	if err != nil {
		return r, errors.Wrap(err, "could not inflate permissions")
	}
	r.PermissionsString = r.PermissionsSnowflake.ToString()

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
