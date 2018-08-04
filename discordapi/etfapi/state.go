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
func NewState() State {
	return State{
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

	s.privateChannels = map[snowflake.Snowflake]Channel{}
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

	s.guilds = map[snowflake.Snowflake]Guild{}
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
		g, err = GuildFromElement(e2)
		if err != nil {
			err = errors.Wrap(err, "could not inflate session guild")
			return
		}
		s.guilds[g.id] = g
	}

	return
}

// UpsertGuildFromElement TODOC
func (s *State) UpsertGuildFromElement(e Element) (err error) {
	eMap, id, err := MapAndIDFromElement(e)
	if err != nil {
		err = errors.Wrap(err, "could not inflate element to find guild")
		return
	}

	g, ok := s.guilds[id]
	if !ok {
		s.guilds[id], err = GuildFromElement(e)
		if err != nil {
			err = errors.Wrap(err, "could not insert guild into the session")
			return
		}
		return
	}

	err = g.UpdateFromElementMap(eMap)
	if err != nil {
		err = errors.Wrap(err, "could not update guild into the session")
		return
	}
	s.guilds[id] = g

	return
}

// UpsertGuildFromElementMap TODOC
func (s *State) UpsertGuildFromElementMap(eMap map[string]Element) (err error) {
	e := eMap["id"]
	id, err := SnowflakeFromElement(e)
	if err != nil {
		err = errors.Wrap(err, "could not find guild id")
		return
	}

	g, ok := s.guilds[id]
	if !ok {
		s.guilds[id], err = GuildFromElement(e)
		if err != nil {
			err = errors.Wrap(err, "could not insert guild into the session")
			return
		}
		return
	}

	err = g.UpdateFromElementMap(eMap)
	if err != nil {
		err = errors.Wrap(err, "could not update guild into the session")
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
		err = errors.Wrap(err, "could not update guild into the session")
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
