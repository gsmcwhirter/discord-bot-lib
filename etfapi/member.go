package etfapi

import (
	"github.com/gsmcwhirter/go-util/v5/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v11/snowflake"
)

// GuildMember represents the information about a known guild membership
type GuildMember struct {
	id    snowflake.Snowflake
	user  User
	roles []snowflake.Snowflake
}

// UpdateFromElementMap updates the information from the given data
//
// This will not remove data; it will only add and change data
func (m *GuildMember) UpdateFromElementMap(eMap map[string]Element) error {
	var eMap2 map[string]Element
	var rEList []Element
	var roleID snowflake.Snowflake
	var userID snowflake.Snowflake

	var err error

	if e, ok := eMap["user"]; ok {
		eMap2, userID, err = MapAndIDFromElement(e)
		if err != nil {
			return errors.Wrap(err, "could not inflate guild member user to map")
		}

		if m.user.id == 0 {
			m.user.id = userID
		}

		if m.id == 0 {
			m.id = m.user.id
		}

		err = m.user.UpdateFromElementMap(eMap2)
		if err != nil {
			return errors.Wrap(err, "could not update user record")
		}
	}

	if rList, ok := eMap["roles"]; ok {
		rEList, err = rList.ToList()
		if err != nil {
			return errors.Wrap(err, "could not inflate guild member role ids")
		}

		m.roles = make([]snowflake.Snowflake, 0, len(rEList))
		for _, re := range rEList {
			roleID, err = SnowflakeFromElement(re)
			if err != nil {
				return errors.Wrap(err, "could not inflate snowflake for guild member role")
			}
			m.roles = append(m.roles, roleID)
		}
	}

	return nil
}

// GuildMemberFromElement generates a new GuildMember object from the given Element
func GuildMemberFromElement(e Element) (GuildMember, error) {
	var m GuildMember

	var rEList []Element
	var roleID snowflake.Snowflake

	eMap, err := e.ToMap()
	if err != nil {
		return m, errors.Wrap(err, "could not inflate guild member to element map")
	}

	err = m.UpdateFromElementMap(eMap)
	if err != nil {
		return m, err
	}

	m.user, err = UserFromElement(eMap["user"])
	if err != nil {
		return m, errors.Wrap(err, "could not inflate guild member user")
	}
	m.id = m.user.id

	rList, ok := eMap["roles"]
	if ok {
		rEList, err = rList.ToList()
		if err != nil {
			return m, errors.Wrap(err, "could not inflate guild member role ids")
		}

		m.roles = make([]snowflake.Snowflake, 0, len(rEList))
		for _, re := range rEList {
			roleID, err = SnowflakeFromElement(re)
			if err != nil {
				return m, errors.Wrap(err, "could not inflate snowflake for guild member role")
			}
			m.roles = append(m.roles, roleID)
		}
	}

	return m, nil
}
