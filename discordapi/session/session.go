package session

import (
	"sync"

	"github.com/gsmcwhirter/discord-bot-lib/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
	"github.com/pkg/errors"
)

// Session TODOC
type Session struct {
	lock      *sync.RWMutex
	sessionID string
	state     *etfapi.State
}

// NewSession TODOC
func NewSession() *Session {
	return &Session{
		lock:  &sync.RWMutex{},
		state: etfapi.NewState(),
	}
}

// ID TODOC
func (s *Session) ID() string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s == nil {
		return ""
	}
	return s.sessionID
}

// Guild TODOC
func (s *Session) Guild(gid snowflake.Snowflake) (etfapi.Guild, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.state.Guild(gid)
}

// GuildOfChannel TODOC
func (s *Session) GuildOfChannel(cid snowflake.Snowflake) (snowflake.Snowflake, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.state.GuildOfChannel(cid)
}

// UpsertGuildFromElement TODOC
func (s *Session) UpsertGuildFromElement(e etfapi.Element) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertGuildFromElement(e)
	return
}

// UpsertGuildFromElementMap TODOC
func (s *Session) UpsertGuildFromElementMap(eMap map[string]etfapi.Element) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertGuildFromElementMap(eMap)
	return
}

// UpsertGuildMemberFromElementMap TODOC
func (s *Session) UpsertGuildMemberFromElementMap(eMap map[string]etfapi.Element) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertGuildMemberFromElementMap(eMap)
	return
}

// UpsertGuildRoleFromElementMap TODOC
func (s *Session) UpsertGuildRoleFromElementMap(eMap map[string]etfapi.Element) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertGuildRoleFromElementMap(eMap)
	return
}

// UpsertChannelFromElement TODOC
func (s *Session) UpsertChannelFromElement(e etfapi.Element) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertChannelFromElement(e)

	return
}

// UpsertChannelFromElementMap TODOC
func (s *Session) UpsertChannelFromElementMap(eMap map[string]etfapi.Element) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertChannelFromElementMap(eMap)

	return
}

// UpdateFromReady TODOC
func (s *Session) UpdateFromReady(p *etfapi.Payload) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	e, ok := p.Data["session_id"]
	if !ok {
		err = errors.Wrap(etfapi.ErrMissingData, "missing session_id")
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

// ChannelName TODOC
func (s *Session) ChannelName(cid snowflake.Snowflake) (string, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.state.ChannelName(cid)
}

// IsGuildAdmin TODOC
func (s *Session) IsGuildAdmin(gid snowflake.Snowflake, uid snowflake.Snowflake) bool {
	g, err := s.Guild(gid)
	if err != nil {
		return false
	}

	return g.IsAdmin(uid)
}

// EveryoneRoleID TODOC
func (s *Session) EveryoneRoleID(gid snowflake.Snowflake) (snowflake.Snowflake, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.state.EveryoneRoleID(gid)
}
