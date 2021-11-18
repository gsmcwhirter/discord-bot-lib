package entity

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v23/snowflake"
)

type ApplicationCommandPermissions struct {
	IDString            string                         `json:"id"`
	ApplicationIDString string                         `json:"application_id"`
	GuildIDString       string                         `json:"guild_id"`
	Permissions         []ApplicationCommandPermission `json:"permissions"`

	IDSnowflake            snowflake.Snowflake `json:"-"`
	ApplicationIDSnowflake snowflake.Snowflake `json:"-"`
	GuildIDSnowflake       snowflake.Snowflake `json:"-"`
}

func (p *ApplicationCommandPermissions) Snowflakify() error {
	var err error

	if p.IDString != "" {
		if p.IDSnowflake, err = snowflake.FromString(p.IDString); err != nil {
			return errors.Wrap(err, "could not snowflakify ID")
		}
	}

	if p.GuildIDString != "" {
		if p.GuildIDSnowflake, err = snowflake.FromString(p.GuildIDString); err != nil {
			return errors.Wrap(err, "could not snowflakify GuildID")
		}
	}

	if p.ApplicationIDString != "" {
		if p.ApplicationIDSnowflake, err = snowflake.FromString(p.ApplicationIDString); err != nil {
			return errors.Wrap(err, "could not snowflakify ApplicationID")
		}
	}

	for _, perm := range p.Permissions {
		if err = perm.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify all Permissions")
		}
	}

	return nil
}

type CommandPermissionType int

const (
	CommandPermissionRole CommandPermissionType = 1
	CommandPermissionUser CommandPermissionType = 2
)

type ApplicationCommandPermission struct {
	IDString   string                `json:"id"`
	Type       CommandPermissionType `json:"type"`
	Permission bool                  `json:"permission"`

	IDSnowflake snowflake.Snowflake `json:"-"`
}

func (p *ApplicationCommandPermission) Snowflakify() error {
	var err error

	if p.IDString != "" {
		if p.IDSnowflake, err = snowflake.FromString(p.IDString); err != nil {
			return errors.Wrap(err, "could not snowflakify ID")
		}
	}

	return nil
}
