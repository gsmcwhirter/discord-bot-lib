package entity

import (
	"github.com/gsmcwhirter/go-util/v10/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
)

// ApplicationCommandPermissions represents the permission sets for a command
type ApplicationCommandPermissions struct {
	IDString            string                         `json:"id"`
	ApplicationIDString string                         `json:"application_id"`
	GuildIDString       string                         `json:"guild_id"`
	Permissions         []ApplicationCommandPermission `json:"permissions"`

	IDSnowflake            snowflake.Snowflake `json:"-"`
	ApplicationIDSnowflake snowflake.Snowflake `json:"-"`
	GuildIDSnowflake       snowflake.Snowflake `json:"-"`
}

// Snowflakify converts snowflake strings into real sowflakes
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

// CommandPermissionType is the type of a permission
type CommandPermissionType int

// These are the possible CommandPermissionType values
const (
	CommandPermissionRole CommandPermissionType = 1
	CommandPermissionUser CommandPermissionType = 2
)

// ApplicationCommandPermission represents an individual permission for a command
type ApplicationCommandPermission struct {
	IDString   string                `json:"id"`
	Type       CommandPermissionType `json:"type"`
	Permission bool                  `json:"permission"`

	IDSnowflake snowflake.Snowflake `json:"-"`
}

// Snowflakify converts snowflake strings into real sowflakes
func (p *ApplicationCommandPermission) Snowflakify() error {
	var err error

	if p.IDString != "" {
		if p.IDSnowflake, err = snowflake.FromString(p.IDString); err != nil {
			return errors.Wrap(err, "could not snowflakify ID")
		}
	}

	return nil
}
