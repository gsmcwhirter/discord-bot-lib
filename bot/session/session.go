package session

import (
	"sync"

	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v22/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v22/snowflake"
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
	s.lock.RLock()         //nolint:staticcheck // this is technically possibly a nil dereference, but won't be in practice
	defer s.lock.RUnlock() //nolint:staticcheck // this is technically possibly a nil dereference, but won't be in practice

	if s == nil { //nolint:staticcheck // safety measure
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

// GuildIDs finds all the currently stored guild ids in the session state
func (s *Session) GuildIDs() []snowflake.Snowflake {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.state.GuildIDs()
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
func (s *Session) IsGuildAdmin(gid, uid snowflake.Snowflake) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	g, ok := s.state.Guild(gid)
	if !ok {
		return false
	}

	return g.IsAdmin(uid)
}

// UpsertGuildFromElement updates data in the session state for a guild based on the given Element
func (s *Session) UpsertGuildFromElement(e etfapi.Element) (snowflake.Snowflake, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.state.UpsertGuildFromElement(e)
}

// UpsertGuildFromElementMap updates data in the session state for a guild based on the given data
func (s *Session) UpsertGuildFromElementMap(eMap map[string]etfapi.Element) (snowflake.Snowflake, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.state.UpsertGuildFromElementMap(eMap)
}

// UpsertGuildMemberFromElementMap updates data in the session state for a guild member based on the given data
func (s *Session) UpsertGuildMemberFromElementMap(eMap map[string]etfapi.Element) (snowflake.Snowflake, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.state.UpsertGuildMemberFromElementMap(eMap)
}

// UpsertGuildRoleFromElementMap updates data in the session state for a guild role based on the given data
func (s *Session) UpsertGuildRoleFromElementMap(eMap map[string]etfapi.Element) (snowflake.Snowflake, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.state.UpsertGuildRoleFromElementMap(eMap)
}

// UpsertChannelFromElement updates data in the session state for a channel based on the given Element
func (s *Session) UpsertChannelFromElement(e etfapi.Element) (snowflake.Snowflake, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.state.UpsertChannelFromElement(e)
}

// UpsertChannelFromElementMap updates data in the session state for a channel based on the given data
func (s *Session) UpsertChannelFromElementMap(eMap map[string]etfapi.Element) (snowflake.Snowflake, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.state.UpsertChannelFromElementMap(eMap)
}

// UpdateFromReady updates data in the session state from a session ready message, and updates the session id
func (s *Session) UpdateFromReady(data map[string]etfapi.Element) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	var err error

	e, ok := data["session_id"]
	if !ok {
		return errors.Wrap(ErrMissingData, "missing session_id")
	}

	s.sessionID, err = e.ToString()
	if err != nil {
		return errors.Wrap(err, "could not inflate session_id")
	}

	return s.state.UpdateFromReady(data)
}
