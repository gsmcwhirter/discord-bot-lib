package etfapi

import (
	"sync"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
	"github.com/pkg/errors"
)

// Session represents a discord bot's session with an api gateway
//
// The primary purpose of this wrapper is to wrap and lock access to a
// State field for safe access from multiple goroutines
type Session struct {
	lock      *sync.RWMutex
	sessionID string
	state     *state
}

// NewSession creates a new session object in an unlocked state and with empty session id
func NewSession() *Session {
	return &Session{
		lock:  &sync.RWMutex{},
		state: newState(),
	}
}

// ID returns the session id of the current session (or an empty string if an id has not been set)
func (s *Session) ID() string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s == nil {
		return ""
	}
	return s.sessionID
}

// Guild finds a guild with the given ID in the current session state, if it exists
//
// The second return value will be false if no such guild was found
func (s *Session) Guild(gid snowflake.Snowflake) (Guild, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.state.Guild(gid)
}

// GuildOfChannel returns the id of the guild that owns the channel with the provided id, if one is known
//
// The second return value will be false if no such guild was found
func (s *Session) GuildOfChannel(cid snowflake.Snowflake) (snowflake.Snowflake, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.state.GuildOfChannel(cid)
}

// ChannelName returns the name of the channel with the provided id, if one is known
//
// The second return value will be alse if not such channel was found
func (s *Session) ChannelName(cid snowflake.Snowflake) (string, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.state.ChannelName(cid)
}

// IsGuildAdmin returns true if the user with the given uid has Admin powers in the guild with
// the given gid. If the guild is not found, this will return false
func (s *Session) IsGuildAdmin(gid snowflake.Snowflake, uid snowflake.Snowflake) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	g, ok := s.state.Guild(gid)
	if !ok {
		return false
	}

	return g.IsAdmin(uid)
}

// UpsertGuildFromElement updates data in the session state for a guild based on the given Element
func (s *Session) UpsertGuildFromElement(e Element) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertGuildFromElement(e)
	return
}

// UpsertGuildFromElementMap updates data in the session state for a guild based on the given data
func (s *Session) UpsertGuildFromElementMap(eMap map[string]Element) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertGuildFromElementMap(eMap)
	return
}

// UpsertGuildMemberFromElementMap updates data in the session state for a guild member based on the given data
func (s *Session) UpsertGuildMemberFromElementMap(eMap map[string]Element) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertGuildMemberFromElementMap(eMap)
	return
}

// UpsertGuildRoleFromElementMap updates data in the session state for a guild role based on the given data
func (s *Session) UpsertGuildRoleFromElementMap(eMap map[string]Element) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertGuildRoleFromElementMap(eMap)
	return
}

// UpsertChannelFromElement updates data in the session state for a channel based on the given Element
func (s *Session) UpsertChannelFromElement(e Element) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertChannelFromElement(e)

	return
}

// UpsertChannelFromElementMap updates data in the session state for a channel based on the given data
func (s *Session) UpsertChannelFromElementMap(eMap map[string]Element) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertChannelFromElementMap(eMap)

	return
}

// UpdateFromReady updates data in the session state from a session ready message, and updates the session id
func (s *Session) UpdateFromReady(p *Payload) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	e, ok := p.Data["session_id"]
	if !ok {
		err = errors.Wrap(ErrMissingData, "missing session_id")
		return
	}

	s.sessionID, err = e.ToString()
	if err != nil {
		err = errors.Wrap(err, "could not inflate session_id")
		return
	}

	err = s.state.UpdateFromReady(p)

	return
}
