package session

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v20/discordapi/etf"
	"github.com/gsmcwhirter/discord-bot-lib/v20/snowflake"
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
func (s *state) UpdateFromReady(data map[string]etf.Element) error {
	var ok bool
	var e etf.Element
	var e2 etf.Element
	var c Channel
	var g Guild
	var gMap map[string]etf.Element
	var gid snowflake.Snowflake
	var err error

	e, ok = data["user"]
	if !ok {
		return errors.Wrap(ErrMissingData, "missing user")
	}
	s.user, err = UserFromElement(e)
	if err != nil {
		return errors.Wrap(err, "could not inflate session user")
	}

	e, ok = data["private_channels"]
	if !ok {
		return errors.Wrap(ErrMissingData, "missing private_channels")
	}
	if !e.Code.IsList() {
		return errors.Wrap(ErrBadData, "private_channels was not a list")
	}
	for _, e2 = range e.Vals {
		c, err = ChannelFromElement(e2)
		if err != nil {
			return errors.Wrap(err, "could not inflate session channel")
		}
		s.privateChannels[c.id] = c
	}

	e, ok = data["guilds"]
	if !ok {
		return errors.Wrap(ErrMissingData, "missing guilds")
	}
	if !e.Code.IsList() {
		return errors.Wrap(ErrBadData, "guilds was not a list")
	}
	for _, e2 = range e.Vals {
		gMap, gid, err = etf.MapAndIDFromElement(e2)
		if err != nil {
			return errors.Wrap(err, "could not inflate session guild to map")
		}

		g, ok = s.guilds[gid]
		if !ok {
			g, err = GuildFromElementMap(gMap)
			if err != nil {
				return errors.Wrap(err, "could not inflate session guild map to guild")
			}
		} else {
			err = g.UpdateFromElementMap(gMap)
			if err != nil {
				return errors.Wrap(err, "could not update guild from guild map")
			}
		}
		s.guilds[gid] = g
	}

	return nil
}

// UpsertGuildFromElement updates data in the session state for a guild based on the given Element
func (s *state) UpsertGuildFromElement(e etf.Element) (snowflake.Snowflake, error) {
	eMap, id, err := etf.MapAndIDFromElement(e)
	if err != nil {
		return 0, errors.Wrap(err, "UpsertGuildFromElement could not inflate element to find guild")
	}

	g, ok := s.guilds[id]
	if !ok {
		s.guilds[id], err = GuildFromElement(e)
		if err != nil {
			return id, errors.Wrap(err, "UpsertGuildFromElement could not insert guild into the session")
		}
		return id, nil
	}

	err = g.UpdateFromElementMap(eMap)
	if err != nil {
		return id, errors.Wrap(err, "UpsertGuildFromElement could not update guild into the session")
	}
	s.guilds[id] = g

	return id, nil
}

// UpsertGuildFromElementMap updates data in the session state for a guild based on the given data
func (s *state) UpsertGuildFromElementMap(eMap map[string]etf.Element) (snowflake.Snowflake, error) {

	e, ok := eMap["id"]
	if !ok {
		return 0, errors.Wrap(ErrMissingData, "UpsertGuildFromElementMap could not find guild id map element")
	}

	id, err := etf.SnowflakeFromElement(e)
	if err != nil {
		return id, errors.Wrap(err, "UpsertGuildFromElementMap could not find guild id")
	}

	g, ok := s.guilds[id]
	if !ok {
		g, err = GuildFromElementMap(eMap)
		if err != nil {
			return id, errors.Wrap(err, "UpsertGuildFromElementMap could not insert guild into the session")
		}
		s.guilds[id] = g

		return id, nil
	}

	err = g.UpdateFromElementMap(eMap)
	if err != nil {
		return id, errors.Wrap(err, "UpsertGuildFromElementMap could not update guild into the session")
	}

	s.guilds[id] = g
	return id, nil
}

// UpsertGuildMemberFromElementMap updates data in the session state for a guild member based on the given data
func (s *state) UpsertGuildMemberFromElementMap(eMap map[string]etf.Element) (snowflake.Snowflake, error) {
	e, ok := eMap["guild_id"]
	if !ok {
		return 0, errors.Wrap(ErrMissingData, "UpsertGuildMemberFromElementMap could not find guild id map element")
	}

	id, err := etf.SnowflakeFromElement(e)
	if err != nil {
		return 0, errors.Wrap(err, "UpsertGuildMemberFromElementMap could not find guild id")
	}

	g, ok := s.guilds[id]
	if !ok {
		return id, errors.Wrap(ErrNotFound, "UpsertGuildMemberFromElementMap could not find the guild to add a member to")
	}

	e, ok = eMap["user"]
	if !ok {
		return id, errors.Wrap(ErrMissingData, "UpsertGuildMemberFromElementMap could not find user element")
	}

	eMap, err = e.ToMap()
	if err != nil {
		return id, errors.Wrap(err, "Up[sertGuildMemberFromElementMap could not convert user element to a map")
	}

	if _, err = g.UpsertMemberFromElementMap(eMap); err != nil {
		return id, errors.Wrap(err, "UpsertGuildMemberFromElementMap could not upsert guild member into the session")
	}

	s.guilds[id] = g
	return id, nil
}

// UpsertGuildRoleFromElementMap updates data in the session state for a guild role based on the given data
func (s *state) UpsertGuildRoleFromElementMap(eMap map[string]etf.Element) (snowflake.Snowflake, error) {
	e, ok := eMap["guild_id"]
	if !ok {
		return 0, errors.Wrap(ErrMissingData, "UpsertGuildRoleFromElementMap could not find guild id map element")
	}

	id, err := etf.SnowflakeFromElement(e)
	if err != nil {
		return 0, errors.Wrap(err, "UpsertGuildRoleFromElementMap could not find guild id")
	}

	g, ok := s.guilds[id]
	if !ok {
		return id, errors.Wrap(ErrNotFound, "UpsertGuildRoleFromElementMap could not find the guild to add a role to")
	}

	e, ok = eMap["role"]
	if !ok {
		return id, errors.Wrap(ErrMissingData, "UpsertGuildRoleFromElementMap could not find the role element")
	}

	eMap, err = e.ToMap()
	if err != nil {
		return id, errors.Wrap(err, "UpsertGuildRoleFromElementMap could not convert role element into a map")
	}

	if _, err = g.UpsertRoleFromElementMap(eMap); err != nil {
		return id, errors.Wrap(err, "UpsertGuildRoleFromElementMap could not upsert guild role into the session")
	}

	s.guilds[id] = g
	return id, nil
}

// UpsertChannelFromElement updates data in the session state for a channel based on the given Element
func (s *state) UpsertChannelFromElement(e etf.Element) (snowflake.Snowflake, error) {
	eMap, id, err := etf.MapAndIDFromElement(e)
	if err != nil {
		return 0, errors.Wrap(err, "could not inflate element to find channel")
	}

	gidE, ok := eMap["guild_id"]
	if !ok || e.IsNil() { // private channel
		c, found := s.privateChannels[id]
		if !found {
			s.privateChannels[id], err = ChannelFromElement(e)
			if err != nil {
				return 0, errors.Wrap(err, "could not insert channel into the session")
			}
			return 0, nil
		}

		err = c.UpdateFromElementMap(eMap)
		if err != nil {
			return 0, errors.Wrap(err, "could not update channel into the session")
		}
		s.privateChannels[id] = c

		return 0, nil
	}

	gid, err := etf.SnowflakeFromElement(gidE)
	if err != nil {
		return 0, errors.Wrap(err, "could not get guild_id from element")
	}

	g, ok := s.guilds[gid]
	if !ok {
		return gid, errors.Wrap(ErrNotFound, "could not find the guild_id")
	}

	c, ok := g.channels[id]
	if !ok { // new channel
		g.channels[id], err = ChannelFromElement(e)
		if err != nil {
			return gid, errors.Wrap(err, "could not insert channel into the session")
		}

		s.guilds[gid] = g
		return gid, nil
	}

	if err = c.UpdateFromElementMap(eMap); err != nil {
		return gid, errors.Wrap(err, "could not update channel into the session")
	}
	g.channels[id] = c
	s.guilds[gid] = g

	return gid, nil
}

// UpsertChannelFromElementMap updates data in the session state for a channel based on the given data
func (s *state) UpsertChannelFromElementMap(eMap map[string]etf.Element) (snowflake.Snowflake, error) {
	var id snowflake.Snowflake

	e, ok := eMap["id"]
	if !ok {
		return 0, errors.Wrap(ErrMissingData, "UpsertChannelFromElementMap could not find channel id map element")
	}

	id, err := etf.SnowflakeFromElement(e)
	if err != nil {
		return 0, errors.Wrap(err, "UpsertChannelFromElementMap could not find channel id")
	}

	gidE, ok := eMap["guild_id"]
	if !ok || gidE.IsNil() { // private channel
		c, found := s.privateChannels[id]
		if !found {
			s.privateChannels[id], err = ChannelFromElementMap(eMap)
			if err != nil {
				return 0, errors.Wrap(err, "could not insert channel into the session")
			}
			return 0, nil
		}

		err = c.UpdateFromElementMap(eMap)
		if err != nil {
			return 0, errors.Wrap(err, "could not update channel into the session")
		}
		s.privateChannels[id] = c

		return 0, nil
	}

	gid, err := etf.SnowflakeFromElement(gidE)
	if err != nil {
		return 0, errors.Wrap(err, "could not get guild_id from element")
	}

	g, ok := s.guilds[gid]
	if !ok {
		return gid, errors.Wrap(ErrNotFound, "could not find the guild_id")
	}

	c, ok := g.channels[id]
	if !ok { // new channel
		g.channels[id], err = ChannelFromElementMap(eMap)
		if err != nil {
			return gid, errors.Wrap(err, "could not insert channel into the session")
		}

		s.guilds[gid] = g
		return gid, nil
	}

	if err = c.UpdateFromElementMap(eMap); err != nil {
		return gid, errors.Wrap(err, "could not update channel into the session")
	}
	g.channels[id] = c
	s.guilds[gid] = g

	return gid, nil
}

// GuildOfChannel returns the id of the guild that owns the channel with the provided id, if one is known
//
// The second return value will be false if no such guild was found
func (s *state) GuildOfChannel(cid snowflake.Snowflake) (snowflake.Snowflake, bool) {
	for _, g := range s.guilds {
		if g.OwnsChannel(cid) {
			return g.id, true
		}
	}

	return 0, false
}

// Guild finds a guild with the given ID in the current session state, if it exists
//
// The second return value will be false if no such guild was found
func (s *state) Guild(gid snowflake.Snowflake) (Guild, bool) {
	g, ok := s.guilds[gid]
	return g, ok
}

// GuildIDs returns the ids of all the guilds in the current session state
func (s *state) GuildIDs() []snowflake.Snowflake {
	gids := make([]snowflake.Snowflake, 0, len(s.guilds))
	for gid := range s.guilds {
		gids = append(gids, gid)
	}

	return gids
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
