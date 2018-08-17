package etfapi

import (
	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

// GuildMember TODOC
type GuildMember struct {
	id    snowflake.Snowflake
	user  User
	roles []snowflake.Snowflake
}

// UpdateFromElementMap TODOC
func (m *GuildMember) UpdateFromElementMap(eMap map[string]Element) (err error) {
	var eMap2 map[string]Element
	var rEList []Element
	var roleID snowflake.Snowflake
	var userID snowflake.Snowflake

	if e, ok := eMap["user"]; ok {
		eMap2, userID, err = MapAndIDFromElement(e)
		if err != nil {
			err = errors.Wrap(err, "could not inflate guild member user to map")
			return
		}

		if m.user.id == 0 {
			m.user.id = userID
		}

		if m.id == 0 {
			m.id = m.user.id
		}

		err = m.user.UpdateFromElementMap(eMap2)
		if err != nil {
			err = errors.Wrap(err, "could not update user record")
			return
		}
	}

	if rList, ok := eMap["roles"]; ok {
		rEList, err = rList.ToList()
		if err != nil {
			err = errors.Wrap(err, "could not inflate guild member role ids")
			return
		}

		m.roles = make([]snowflake.Snowflake, 0, len(rEList))
		for _, re := range rEList {
			roleID, err = SnowflakeFromElement(re)
			if err != nil {
				err = errors.Wrap(err, "could not inflate snowflake for guild member role")
				return
			}
			m.roles = append(m.roles, roleID)
		}
	}

	return
}

// GuildMemberFromElement TODOC
func GuildMemberFromElement(e Element) (m GuildMember, err error) {
	var rEList []Element
	var roleID snowflake.Snowflake

	eMap, err := e.ToMap()
	if err != nil {
		err = errors.Wrap(err, "could not inflate guild member to element map")
		return
	}

	err = m.UpdateFromElementMap(eMap)
	if err != nil {
		return
	}

	m.user, err = UserFromElement(eMap["user"])
	if err != nil {
		err = errors.Wrap(err, "could not inflate guild member user")
		return
	}
	m.id = m.user.id

	rList, ok := eMap["roles"]
	if ok {
		rEList, err = rList.ToList()
		if err != nil {
			err = errors.Wrap(err, "could not inflate guild member role ids")
			return
		}

		m.roles = make([]snowflake.Snowflake, 0, len(rEList))
		for _, re := range rEList {
			roleID, err = SnowflakeFromElement(re)
			if err != nil {
				err = errors.Wrap(err, "could not inflate snowflake for guild member role")
				return
			}
			m.roles = append(m.roles, roleID)
		}
	}

	return
}
