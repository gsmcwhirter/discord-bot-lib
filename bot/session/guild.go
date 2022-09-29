package session

import (
	"strings"

	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
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
	name = strings.ToLower(name)

	for _, c := range g.channels {
		if strings.ToLower(c.name) == name {
			return c.id, true
		}
	}

	return 0, false
}

// RoleWithName finds the role id for the role with the provided name
func (g *Guild) RoleWithName(name string) (snowflake.Snowflake, bool) {
	name = strings.ToLower(name)

	for _, r := range g.roles {
		if strings.ToLower(r.name) == name {
			return r.id, true
		}
	}

	return 0, false
}

// HasRole determines if the user with the provided ID has the role with the provided id
func (g *Guild) HasRole(uid, rid snowflake.Snowflake) bool {
	gm, ok := g.members[uid]
	if !ok {
		return false
	}

	for _, rid2 := range gm.roles {
		if rid2 == rid {
			return true
		}
	}

	return false
}

// RoleIsAdministrator determines if a role has admin powers in the guild
func (g *Guild) RoleIsAdministrator(rid snowflake.Snowflake) bool {
	r, ok := g.roles[rid]
	if !ok {
		return false
	}

	return r.IsAdmin()
}

// AllAdministratorRoleIDs gets the role ids of all roles known to be administrators of the guild
func (g *Guild) AllAdministratorRoleIDs() []snowflake.Snowflake {
	rids := make([]snowflake.Snowflake, 0, len(g.roles))
	for _, r := range g.roles {
		if r.IsAdmin() {
			rids = append(rids, r.id)
		}
	}

	return rids
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
		if g.RoleIsAdministrator(rid) {
			return true
		}
	}

	return false
}

// UpdateFromElementMap updates information about the guild from the provided data
//
// This will not delete data; it will only add and change data
func (g *Guild) UpdateFromElementMap(eMap map[string]etfapi.Element) error {
	var ok bool
	var e2 etfapi.Element
	var m GuildMember
	var c Channel
	var r Role
	var err error

	e2, ok = eMap["owner_id"]
	if ok && !e2.IsNil() {
		g.ownerID, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return errors.Wrap(err, "could not get owner_id snowflake.Snowflake")
		}
	}

	e2, ok = eMap["application_id"]
	if ok && !e2.IsNil() {
		g.applicationID, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return errors.Wrap(err, "could not get application_id snowflake.Snowflake")
		}
	}

	e2, ok = eMap["name"]
	if ok {
		g.name, err = e2.ToString()
		if err != nil {
			return errors.Wrap(err, "could not get name")
		}
	}

	if !g.available {
		e2, ok = eMap["unavailable"]
		if ok {
			var uavStr string
			uavStr, err = e2.ToString()
			if err != nil {
				return errors.Wrap(err, "could not get unavailable status")
			}

			switch uavStr {
			case "true":
				g.available = false
			case "false":
				g.available = true
			default:
				return errors.Wrap(ErrBadData, "did not get true or false availability")
			}
		}
	}

	e2, ok = eMap["members"]
	if ok {
		for _, e3 := range e2.Vals {
			m, err = GuildMemberFromElement(e3)
			if err != nil {
				return errors.Wrap(err, "could not inflate guild member")
			}
			g.members[m.id] = m
		}
	}

	e2, ok = eMap["channels"]
	if ok {
		for _, e3 := range e2.Vals {
			c, err = ChannelFromElement(e3)
			if err != nil {
				return errors.Wrap(err, "could not inflate guild channel")
			}
			g.channels[c.id] = c
		}
	}

	e2, ok = eMap["roles"]
	if ok {
		for _, e3 := range e2.Vals {
			r, err = RoleFromElement(e3)
			if err != nil {
				return errors.Wrap(err, "could not inflate guild role")
			}
			g.roles[r.id] = r
		}
	}

	return nil
}

// UpsertMemberFromElementMap upserts a GuildMemeber in the guild from the given data
func (g *Guild) UpsertMemberFromElementMap(eMap map[string]etfapi.Element) (GuildMember, error) {
	var m GuildMember

	mid, err := etfapi.SnowflakeFromElement(eMap["id"])
	if err != nil {
		return m, errors.Wrap(err, "could not get member id")
	}

	m, ok := g.members[mid]
	if !ok {
		m.id = mid
	}
	err = m.UpdateFromElementMap(eMap)
	if err != nil {
		return m, err
	}
	g.members[mid] = m

	return m, nil
}

// UpsertRoleFromElementMap upserts a Role in the guild from the given data
func (g *Guild) UpsertRoleFromElementMap(eMap map[string]etfapi.Element) (Role, error) {
	var r Role

	rid, err := etfapi.SnowflakeFromElement(eMap["id"])
	if err != nil {
		return r, errors.Wrap(err, "could not get role id")
	}

	r, ok := g.roles[rid]
	if !ok {
		r.id = rid
	}
	err = r.UpdateFromElementMap(eMap)
	if err != nil {
		return r, err
	}

	g.roles[rid] = r
	return r, nil
}

// GuildFromElementMap creates a new Guild object from the given data
func GuildFromElementMap(eMap map[string]etfapi.Element) (Guild, error) {
	g := Guild{
		channels: map[snowflake.Snowflake]Channel{},
		members:  map[snowflake.Snowflake]GuildMember{},
		roles:    map[snowflake.Snowflake]Role{},
	}

	var err error

	g.id, err = etfapi.SnowflakeFromElement(eMap["id"])
	if err != nil {
		return g, errors.Wrap(err, "could not get guild id")
	}

	err = g.UpdateFromElementMap(eMap)
	return g, errors.Wrap(err, "could not create a guild")
}

// GuildFromElement creates a new Guild object from the given Element
func GuildFromElement(e etfapi.Element) (Guild, error) {
	eMap, _, err := etfapi.MapAndIDFromElement(e)
	if err != nil {
		return Guild{}, errors.Wrap(err, "could not create guild map")
	}

	return GuildFromElementMap(eMap)
}
