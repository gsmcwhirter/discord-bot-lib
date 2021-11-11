package entity

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v22/snowflake"
)

// Emoji is the data about an emoji recevied from the json api
type Emoji struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Roles         []Role `json:"roles"`
	User          User   `json:"user"`
	RequireColons bool   `json:"require_colons"`
	Managed       bool   `json:"managed"`
	Animated      bool   `json:"animated"`
	Available     bool   `json:"available"`

	IDSnowflake snowflake.Snowflake
}

func (er *Emoji) Snowflakify() error {
	var err error

	if er.ID != "" {
		if er.IDSnowflake, err = snowflake.FromString(er.ID); err != nil {
			return errors.Wrap(err, "could not snowflakify ID")
		}
	}

	if err = er.User.Snowflakify(); err != nil {
		return errors.Wrap(err, "could not snowflakify User")
	}

	for i := range er.Roles {
		m := er.Roles[i]
		if err = m.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify Roles")
		}
		er.Roles[i] = m
	}

	return nil
}
