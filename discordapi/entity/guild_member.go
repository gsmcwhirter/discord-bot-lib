package entity

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
)

// GuildMember is the data about a guild member recevied from the json api
type GuildMember struct {
	User         *User    `json:"user"`
	Nick         string   `json:"nick"`
	Roles        []string `json:"roles"`
	JoinedAt     string   `json:"joined_at"`     // ISO8601
	PremiumSince string   `json:"premium_since"` // ISO8601
	Deaf         bool     `json:"deaf"`
	Mute         bool     `json:"mute"`

	RoleSnowflakes []snowflake.Snowflake `json:"-"`
}

// Snowflakify converts snowflake strings into real sowflakes
func (gmr *GuildMember) Snowflakify() error {
	sfs := make([]snowflake.Snowflake, 0, len(gmr.Roles))
	for _, r := range gmr.Roles {
		sf, err := snowflake.FromString(r)
		if err != nil {
			return errors.Wrap(err, "could not convert role strings to snowflakes")
		}
		sfs = append(sfs, sf)
	}

	gmr.RoleSnowflakes = sfs

	return nil
}

// HasRole determines if the guild member is currently known to have the requested role
func (gmr GuildMember) HasRole(rid snowflake.Snowflake) bool {
	for _, r := range gmr.RoleSnowflakes {
		if r == rid {
			return true
		}
	}

	return false
}

// GuildMemberFromElement instantiates a GuildMember from an etf element
func GuildMemberFromElement(e etfapi.Element) (GuildMember, error) {
	var m GuildMember

	eMap, err := e.ToMap()
	if err != nil {
		return m, errors.Wrap(err, "could not inflate GuildMember from non-map")
	}

	e2, ok := eMap["user"]
	if ok && !e2.IsNil() {
		v, err := UserFromElement(e2)
		if err != nil {
			return m, errors.Wrap(err, "could not inflate User")
		}
		m.User = &v
	}

	e2 = eMap["nick"]
	m.Nick, err = e2.ToString()
	if err != nil {
		return m, errors.Wrap(err, "could not inflate Nick")
	}

	e2 = eMap["joined_at"]
	m.JoinedAt, err = e2.ToString()
	if err != nil {
		return m, errors.Wrap(err, "could not inflate JoinedAt")
	}

	e2 = eMap["premium_since"]
	m.PremiumSince, err = e2.ToString()
	if err != nil {
		return m, errors.Wrap(err, "could not inflate PremiumSince")
	}

	e2, ok = eMap["deaf"]
	if ok {
		m.Deaf, err = e2.ToBool()
		if err != nil {
			return m, errors.Wrap(err, "could not inflate Deaf", "element_type", e2.Code.String())
		}
	}

	e2, ok = eMap["mute"]
	if ok {
		m.Mute, err = e2.ToBool()
		if err != nil {
			return m, errors.Wrap(err, "could not inflate Mute", "element_type", e2.Code.String())
		}
	}

	e2, ok = eMap["roles"]
	if ok && !e2.IsNil() {
		el, err := e2.ToList()
		if err != nil {
			return m, errors.Wrap(err, "could not inflate roles list")
		}

		m.RoleSnowflakes = make([]snowflake.Snowflake, 0, len(el))
		m.Roles = make([]string, 0, len(el))

		for _, e3 := range el {
			s, err := etfapi.SnowflakeFromUnknownElement(e3)
			if err != nil {
				return m, errors.Wrap(err, "could not inflate snowflake for role")
			}

			m.RoleSnowflakes = append(m.RoleSnowflakes, s)
			m.Roles = append(m.Roles, s.ToString())
		}
	}

	return m, nil
}
