package etfapi

import (
	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

// Guild represents the known data about a discord guild
type Guild struct {
	id            snowflake.Snowflake
	ownerID       snowflake.Snowflake
	applicationID snowflake.Snowflake
	name          string
	available     bool
	members       map[snowflake.Snowflake]GuildMember
	channels      map[snowflake.Snowflake]Channel
	roles         map[snowflake.Snowflake]Role
}

// ID returns the guild ID
func (g *Guild) ID() snowflake.Snowflake {
	return g.id
}

// OwnsChannel determines if this guild owns a channel with the provided id
func (g *Guild) OwnsChannel(cid snowflake.Snowflake) bool {
	_, ok := g.channels[cid]
	return ok
}

// ChannelWithName finds the channel id for the channel with the provided name
func (g *Guild) ChannelWithName(name string) (snowflake.Snowflake, bool) {
	for _, c := range g.channels {
		if c.name == name {
			return c.id, true
		}
	}

	return 0, false
}

// IsAdmin determines if the user with the provided ID has administrator powers
// in the guild
func (g *Guild) IsAdmin(uid snowflake.Snowflake) bool {
	if g.ownerID != 0 && uid == g.ownerID {
		return true
	}

	gm, ok := g.members[uid]
	if !ok {
		return false
	}

	for _, rid := range gm.roles {
		r, ok := g.roles[rid]
		if !ok {
			continue
		}

		if r.IsAdmin() {
			return true
		}
	}

	return false
}

// UpdateFromElementMap updates information about the guild from the provided data
//
// This will not delete data; it will only add and change data
func (g *Guild) UpdateFromElementMap(eMap map[string]Element) (err error) {
	var ok bool
	var e2 Element
	var m GuildMember
	var c Channel
	var r Role

	e2, ok = eMap["owner_id"]
	if ok && !e2.IsNil() {
		g.ownerID, err = SnowflakeFromElement(e2)
		if err != nil {
			err = errors.Wrap(err, "could not get owner_id snowflake.Snowflake")
			return
		}
	}

	e2, ok = eMap["application_id"]
	if ok && !e2.IsNil() {
		g.applicationID, err = SnowflakeFromElement(e2)
		if err != nil {
			err = errors.Wrap(err, "could not get application_id snowflake.Snowflake")
			return
		}
	}

	e2, ok = eMap["name"]
	if ok {
		g.name, err = e2.ToString()
		if err != nil {
			err = errors.Wrap(err, "could not get name")
			return
		}
	}

	if !g.available {
		e2, ok = eMap["unavailable"]
		if ok {
			var uavStr string
			uavStr, err = e2.ToString()
			if err != nil {
				err = errors.Wrap(err, "could not get unavailable status")
				return
			}

			switch uavStr {
			case "true":
				g.available = false
			case "false":
				g.available = true
			default:
				err = errors.Wrap(ErrBadData, "did not get true or false availability")
				return
			}
		}
	}

	e2, ok = eMap["members"]
	if ok {
		for _, e3 := range e2.Vals {
			m, err = GuildMemberFromElement(e3)
			if err != nil {
				err = errors.Wrap(err, "could not inflate guild member")
				return
			}
			g.members[m.id] = m
		}
	}

	e2, ok = eMap["channels"]
	if ok {
		for _, e3 := range e2.Vals {
			c, err = ChannelFromElement(e3)
			if err != nil {
				err = errors.Wrap(err, "could not inflate guild channel")
				return
			}
			g.channels[c.id] = c
		}
	}

	e2, ok = eMap["roles"]
	if ok {
		for _, e3 := range e2.Vals {
			r, err = RoleFromElement(e3)
			if err != nil {
				err = errors.Wrap(err, "could not inflate guild role")
				return
			}
			g.roles[r.id] = r
		}
	}

	return
}

// UpsertMemberFromElementMap upserts a GuildMemeber in the guild from the given data
func (g *Guild) UpsertMemberFromElementMap(eMap map[string]Element) (m GuildMember, err error) {
	mid, err := SnowflakeFromElement(eMap["id"])
	if err != nil {
		err = errors.Wrap(err, "could not get member id")
		return
	}

	m, ok := g.members[mid]
	if !ok {
		m.id = mid
	}
	err = m.UpdateFromElementMap(eMap)
	if err != nil {
		return
	}
	g.members[mid] = m

	return
}

// UpsertRoleFromElementMap upserts a Role in the guild from the given data
func (g *Guild) UpsertRoleFromElementMap(eMap map[string]Element) (r Role, err error) {
	rid, err := SnowflakeFromElement(eMap["id"])
	if err != nil {
		err = errors.Wrap(err, "could not get role id")
		return
	}

	r, ok := g.roles[rid]
	if !ok {
		r.id = rid
	}
	err = r.UpdateFromElementMap(eMap)
	if err != nil {
		return
	}

	g.roles[rid] = r

	return
}

// GuildFromElementMap creates a new Guild object from the given data
func GuildFromElementMap(eMap map[string]Element) (g Guild, err error) {
	g.channels = map[snowflake.Snowflake]Channel{}
	g.members = map[snowflake.Snowflake]GuildMember{}
	g.roles = map[snowflake.Snowflake]Role{}

	g.id, err = SnowflakeFromElement(eMap["id"])
	if err != nil {
		err = errors.Wrap(err, "could not get guild id")
		return
	}

	err = g.UpdateFromElementMap(eMap)
	err = errors.Wrap(err, "could not create a guild")
	return
}

// GuildFromElement creates a new Guild object from the given Element
func GuildFromElement(e Element) (g Guild, err error) {
	var eMap map[string]Element

	eMap, _, err = MapAndIDFromElement(e)
	if err != nil {
		err = errors.Wrap(err, "could not create guild map")
		return
	}

	g, err = GuildFromElementMap(eMap)
	return
}