package discordapi

import (
	"sync"

	"github.com/gsmcwhirter/discord-bot-lib/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
	"github.com/pkg/errors"
)

// Session TODOC
type Session struct {
	lock      sync.RWMutex
	sessionID string
	state     etfapi.State
}

// NewSession TODOC
func NewSession() Session {
	return Session{
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

// UpsertChannelFromElement TODOC
func (s *Session) UpsertChannelFromElement(e etfapi.Element) (err error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	err = s.state.UpsertChannelFromElement(e)

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
