package etfapi

import (
	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

// state represents the state of a current bot session
//
// This object is not concurrency safe; it should be accessed through a Session
// object which handles appropriate locking
type state struct {
	user            User
	guilds          map[snowflake.Snowflake]Guild
	privateChannels map[snowflake.Snowflake]Channel
}

// newState constructs a new, empty state
func newState() *state {
	return &state{
		guilds:          map[snowflake.Snowflake]Guild{},
		privateChannels: map[snowflake.Snowflake]Channel{},
	}
}

// UpdateFromReady updates the session state from the given "ready" payload
func (s *state) UpdateFromReady(p *Payload) (err error) {
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
			err = g.UpdateFromElementMap(gMap)
			if err != nil {
				err = errors.Wrap(err, "could not update guild from guild map")
				return
			}
		}
		s.guilds[gid] = g
	}

	return
}

// UpsertGuildFromElement updates data in the session state for a guild based on the given Element
func (s *state) UpsertGuildFromElement(e Element) (err error) {
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

// UpsertGuildFromElementMap updates data in the session state for a guild based on the given data
func (s *state) UpsertGuildFromElementMap(eMap map[string]Element) (err error) {
	var id snowflake.Snowflake

	e, ok := eMap["id"]
	if !ok {
		err = errors.Wrap(ErrMissingData, "UpsertGuildFromElementMap could not find guild id map element")
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

// UpsertGuildMemberFromElementMap updates data in the session state for a guild member based on the given data
func (s *state) UpsertGuildMemberFromElementMap(eMap map[string]Element) (err error) {
	var id snowflake.Snowflake

	e, ok := eMap["guild_id"]
	if !ok {
		err = errors.Wrap(ErrMissingData, "UpsertGuildMemberFromElementMap could not find guild id map element")
		return
	}

	id, err = SnowflakeFromElement(e)
	if err != nil {
		err = errors.Wrap(err, "UpsertGuildMemberFromElementMap could not find guild id")
		return
	}

	g, ok := s.guilds[id]
	if !ok {
		err = errors.Wrap(ErrNotFound, "UpsertGuildMemberFromElementMap could not find the guild to add a member to")
		return
	}

	e, ok = eMap["user"]
	if !ok {
		err = errors.Wrap(ErrMissingData, "UpsertGuildMemberFromElementMap could not find user element")
		return
	}

	eMap, err = e.ToMap()
	if err != nil {
		err = errors.Wrap(err, "Up[sertGuildMemberFromElementMap could not convert user element to a map")
		return
	}

	if _, err = g.UpsertMemberFromElementMap(eMap); err != nil {
		err = errors.Wrap(err, "UpsertGuildMemberFromElementMap could not upsert guild member into the session")
		return
	}
	s.guilds[id] = g

	return
}

// UpsertGuildRoleFromElementMap updates data in the session state for a guild role based on the given data
func (s *state) UpsertGuildRoleFromElementMap(eMap map[string]Element) (err error) {
	var id snowflake.Snowflake

	e, ok := eMap["guild_id"]
	if !ok {
		err = errors.Wrap(ErrMissingData, "UpsertGuildRoleFromElementMap could not find guild id map element")
		return
	}

	id, err = SnowflakeFromElement(e)
	if err != nil {
		err = errors.Wrap(err, "UpsertGuildRoleFromElementMap could not find guild id")
		return
	}

	g, ok := s.guilds[id]
	if !ok {
		err = errors.Wrap(ErrNotFound, "UpsertGuildRoleFromElementMap could not find the guild to add a role to")
		return
	}

	e, ok = eMap["role"]
	if !ok {
		err = errors.Wrap(ErrMissingData, "UpsertGuildRoleFromElementMap could not find the role element")
		return
	}

	eMap, err = e.ToMap()
	if err != nil {
		err = errors.Wrap(err, "UpsertGuildRoleFromElementMap could not convert role element into a map")
		return
	}

	if _, err = g.UpsertRoleFromElementMap(eMap); err != nil {
		err = errors.Wrap(err, "UpsertGuildRoleFromElementMap could not upsert guild role into the session")
		return
	}
	s.guilds[id] = g

	return
}

// UpsertChannelFromElement updates data in the session state for a channel based on the given Element
func (s *state) UpsertChannelFromElement(e Element) (err error) {
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

// UpsertChannelFromElementMap updates data in the session state for a channel based on the given data
func (s *state) UpsertChannelFromElementMap(eMap map[string]Element) (err error) {
	var id snowflake.Snowflake

	e, ok := eMap["id"]
	if !ok {
		err = errors.Wrap(ErrMissingData, "UpsertChannelFromElementMap could not find channel id map element")
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

// GuildOfChannel returns the id of the guild that owns the channel with the provided id, if one is known
//
// The second return value will be false if no such guild was found
func (s *state) GuildOfChannel(cid snowflake.Snowflake) (gid snowflake.Snowflake, ok bool) {
	for _, g := range s.guilds {
		if g.OwnsChannel(cid) {
			gid = g.id
			ok = true
			return
		}
	}

	return
}

// Guild finds a guild with the given ID in the current session state, if it exists
//
// The second return value will be false if no such guild was found
func (s *state) Guild(gid snowflake.Snowflake) (g Guild, ok bool) {
	g, ok = s.guilds[gid]

	return
}

// ChannelName returns the name of the channel with the provided id, if one is known
//
// The second return value will be alse if not such channel was found
func (s *state) ChannelName(cid snowflake.Snowflake) (string, bool) {
	for _, g := range s.guilds {
		c, ok := g.channels[cid]
		if !ok {
			continue
		}

		return c.name, true
	}
	return "", false
}
