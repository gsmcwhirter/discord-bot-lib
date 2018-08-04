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

// GuildMemberFromElement TODOC
func GuildMemberFromElement(e Element) (m GuildMember, err error) {
	var rEList []Element
	var roleID snowflake.Snowflake

	eMap, err := e.ToMap()
	if err != nil {
		err = errors.Wrap(err, "could not inflate guild member to element map")
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
