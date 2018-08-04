package etfapi

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

// Guild TODOC
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

// ID TODOC
func (g *Guild) ID() snowflake.Snowflake {
	return g.id
}

// OwnsChannel TODOC
func (g *Guild) OwnsChannel(cid snowflake.Snowflake) bool {
	fmt.Printf("*** %v -> %+v\n", g.id, g.channels)
	_, ok := g.channels[cid]
	return ok
}

// IsAdmin TODOC
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

// UpdateFromElementMap TODOC
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

// GuildFromElementMap TODOC
func GuildFromElementMap(eMap map[string]Element) (g Guild, err error) {
	g.channels = map[snowflake.Snowflake]Channel{}
	g.members = map[snowflake.Snowflake]GuildMember{}
	g.roles = map[snowflake.Snowflake]Role{}

	g.id, err = SnowflakeFromElement(eMap["id"])
	err = errors.Wrap(err, "could not get guild id")

	err = g.UpdateFromElementMap(eMap)
	err = errors.Wrap(err, "could not create a guild")
	return
}

// GuildFromElement TODOC
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
