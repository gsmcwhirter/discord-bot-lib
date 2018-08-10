package etfapi

import (
	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

// State TODOC
type State struct {
	user            User
	guilds          map[snowflake.Snowflake]Guild
	privateChannels map[snowflake.Snowflake]Channel
}

// NewState TODOC
func NewState() *State {
	return &State{
		guilds:          map[snowflake.Snowflake]Guild{},
		privateChannels: map[snowflake.Snowflake]Channel{},
	}
}

// UpdateFromReady TODOC
func (s *State) UpdateFromReady(p *Payload) (err error) {
	var ok bool
	var e Element
	var e2 Element
	var c Channel
	var g Guild
	var gMap map[string]Element
	var gid snowflake.Snowflake

	e, ok = p.Data["user"]
	if !ok {
		err = errors.Wrap(ErrMissingData, "missing user")
		return
	}
	s.user, err = UserFromElement(e)
	if err != nil {
		err = errors.Wrap(err, "could not inflate session user")
		return
	}

	e, ok = p.Data["private_channels"]
	if !ok {
		err = errors.Wrap(ErrMissingData, "missing private_channels")
		return
	}
	if !e.Code.IsList() {
		err = errors.Wrap(ErrBadData, "private_channels was not a list")
		return
	}
	for _, e2 = range e.Vals {
		c, err = ChannelFromElement(e2)
		if err != nil {
			err = errors.Wrap(err, "could not inflate session channel")
			return
		}
		s.privateChannels[c.id] = c
	}

	e, ok = p.Data["guilds"]
	if !ok {
		err = errors.Wrap(ErrMissingData, "missing guilds")
		return
	}
	if !e.Code.IsList() {
		err = errors.Wrap(ErrBadData, "guilds was not a list")
		return
	}
	for _, e2 = range e.Vals {
		gMap, gid, err = MapAndIDFromElement(e2)
		if err != nil {
			err = errors.Wrap(err, "could not inflate session guild to map")
			return
		}

		g, ok = s.guilds[gid]
		if !ok {
			g, err = GuildFromElementMap(gMap)
			if err != nil {
				err = errors.Wrap(err, "could not inflate session guild map to guild")
				return
			}
		} else {
			g.UpdateFromElementMap(gMap)
		}
		s.guilds[gid] = g
	}

	return
}

// UpsertGuildFromElement TODOC
func (s *State) UpsertGuildFromElement(e Element) (err error) {
	eMap, id, err := MapAndIDFromElement(e)
	if err != nil {
		err = errors.Wrap(err, "UpsertGuildFromElement could not inflate element to find guild")
		return
	}

	g, ok := s.guilds[id]
	if !ok {
		s.guilds[id], err = GuildFromElement(e)
		if err != nil {
			err = errors.Wrap(err, "UpsertGuildFromElement could not insert guild into the session")
			return
		}
		return
	}

	err = g.UpdateFromElementMap(eMap)
	if err != nil {
		err = errors.Wrap(err, "UpsertGuildFromElement could not update guild into the session")
		return
	}
	s.guilds[id] = g

	return
}

// UpsertGuildFromElementMap TODOC
func (s *State) UpsertGuildFromElementMap(eMap map[string]Element) (err error) {
	var id snowflake.Snowflake

	e, ok := eMap["id"]
	if !ok {
		err = errors.Wrap(ErrBadElementData, "UpsertGuildFromElementMap could not find guild id map element")
		return
	}

	id, err = SnowflakeFromElement(e)
	if err != nil {
		err = errors.Wrap(err, "UpsertGuildFromElementMap could not find guild id")
		return
	}

	g, ok := s.guilds[id]
	if !ok {
		g, err = GuildFromElementMap(eMap)
		if err != nil {
			err = errors.Wrap(err, "UpsertGuildFromElementMap could not insert guild into the session")
		} else {
			s.guilds[id] = g
		}
		return
	}

	err = g.UpdateFromElementMap(eMap)
	if err != nil {
		err = errors.Wrap(err, "UpsertGuildFromElementMap could not update guild into the session")
		return
	}
	s.guilds[id] = g

	return
}

// UpsertGuildMemberFromElementMap TODOC
func (s *State) UpsertGuildMemberFromElementMap(eMap map[string]Element) (err error) {
	var id snowflake.Snowflake

	e, ok := eMap["guild_id"]
	if !ok {
		err = errors.Wrap(ErrBadElementData, "UpsertGuildMemberFromElementMap could not find guild id map element")
		return
	}

	id, err = SnowflakeFromElement(e)
	if err != nil {
		err = errors.Wrap(err, "UpsertGuildMemberFromElementMap could not find guild id")
		return
	}

	g, ok := s.guilds[id]
	if !ok {
		err = errors.Wrap(ErrBadElementData, "UpsertGuildMemberFromElementMap could not find the guild to add a member to")
		return
	}

	if _, err = g.UpsertMemberFromElementMap(eMap); err != nil {
		err = errors.Wrap(err, "UpsertGuildMemberFromElementMap could not upsert guild member into the session")
		return
	}
	s.guilds[id] = g

	return
}

// UpsertGuildRoleFromElementMap TODOC
func (s *State) UpsertGuildRoleFromElementMap(eMap map[string]Element) (err error) {
	var id snowflake.Snowflake

	e, ok := eMap["guild_id"]
	if !ok {
		err = errors.Wrap(ErrBadElementData, "UpsertGuildRoleFromElementMap could not find guild id map element")
		return
	}

	id, err = SnowflakeFromElement(e)
	if err != nil {
		err = errors.Wrap(err, "UpsertGuildRoleFromElementMap could not find guild id")
		return
	}

	g, ok := s.guilds[id]
	if !ok {
		err = errors.Wrap(ErrBadElementData, "UpsertGuildRoleFromElementMap could not find the guild to add a role to")
		return
	}

	if _, err = g.UpsertRoleFromElementMap(eMap); err != nil {
		err = errors.Wrap(err, "UpsertGuildRoleFromElementMap could not upsert guild role into the session")
		return
	}
	s.guilds[id] = g

	return
}

// UpsertChannelFromElement TODOC
func (s *State) UpsertChannelFromElement(e Element) (err error) {
	eMap, id, err := MapAndIDFromElement(e)
	if err != nil {
		err = errors.Wrap(err, "could not inflate element to find channel")
		return
	}

	c, ok := s.privateChannels[id]
	if !ok {
		s.privateChannels[id], err = ChannelFromElement(e)
		if err != nil {
			err = errors.Wrap(err, "could not insert channel into the session")
			return
		}
		return
	}

	err = c.UpdateFromElementMap(eMap)
	if err != nil {
		err = errors.Wrap(err, "could not update channel into the session")
		return
	}
	s.privateChannels[id] = c

	return
}

// UpsertChannelFromElementMap TODOC
func (s *State) UpsertChannelFromElementMap(eMap map[string]Element) (err error) {
	var id snowflake.Snowflake

	e, ok := eMap["id"]
	if !ok {
		err = errors.Wrap(ErrBadElementData, "UpsertChannelFromElementMap could not find channel id map element")
		return
	}

	id, err = SnowflakeFromElement(e)
	if err != nil {
		err = errors.Wrap(err, "UpsertChannelFromElementMap could not find channel id")
		return
	}

	c, ok := s.privateChannels[id]
	if !ok {
		s.privateChannels[id], err = ChannelFromElementMap(eMap)
		if err != nil {
			err = errors.Wrap(err, "could not insert channel into the session")
			return
		}
		return
	}

	err = c.UpdateFromElementMap(eMap)
	if err != nil {
		err = errors.Wrap(err, "could not update channel into the session")
		return
	}
	s.privateChannels[id] = c

	return
}

// GuildOfChannel TODOC
func (s *State) GuildOfChannel(cid snowflake.Snowflake) (gid snowflake.Snowflake, ok bool) {
	for _, g := range s.guilds {
		if g.OwnsChannel(cid) {
			gid = g.id
			ok = true
			return
		}
	}

	return
}

// Guild TODOC
func (s *State) Guild(gid snowflake.Snowflake) (g Guild, err error) {
	var ok bool
	g, ok = s.guilds[gid]
	if !ok {
		err = ErrNotFound
		return
	}

	return
}

// ChannelName TODOC
func (s *State) ChannelName(cid snowflake.Snowflake) (string, bool) {
	for _, g := range s.guilds {
		c, ok := g.channels[cid]
		if !ok {
			continue
		}

		return c.name, true
	}
	return "", false
}

// EveryoneRoleID TODOC
func (s *State) EveryoneRoleID(gid snowflake.Snowflake) (rid snowflake.Snowflake, ok bool) {
	g, ok := s.guilds[gid]
	if !ok {
		return
	}

	for _, r := range g.roles {
		if r.name == "everyone" {
			rid = r.id
			return
		}
	}

	ok = false
	return
}
