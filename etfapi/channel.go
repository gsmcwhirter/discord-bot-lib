package etfapi

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

// ChannelType represents the type of a discord channel
type ChannelType int

// These are the known discord channel types
const (
	GuildTextChannel     ChannelType = 0
	DMChannel                        = 1
	GuildVoiceChannel                = 2
	GroupDMChannel                   = 3
	GuildCategoryChannel             = 4
)

// ChannelTypeFromElement extracts the channel type from a etf Element
func ChannelTypeFromElement(e Element) (t ChannelType, err error) {
	temp, err := e.ToInt()
	if err != nil {
		err = errors.Wrap(err, "could not unmarshal channelType")
	}
	t = ChannelType(temp)
	return
}

func (t ChannelType) String() string {
	switch t {
	case GuildTextChannel:
		return "GUILD_TEXT"
	case DMChannel:
		return "DM"
	case GuildVoiceChannel:
		return "GUILD_VOICE"
	case GroupDMChannel:
		return "GROUP_DM"
	case GuildCategoryChannel:
		return "GUILD_CATEGORY"
	default:
		return fmt.Sprintf("(unknown: %d)", int(t))
	}
}

// Channel represents known information about a discord channel
type Channel struct {
	id            snowflake.Snowflake
	guildID       snowflake.Snowflake
	ownerID       snowflake.Snowflake
	applicationID snowflake.Snowflake
	lastMessageID snowflake.Snowflake
	parentID      snowflake.Snowflake
	channelType   ChannelType
	name          string
	topic         string
	recipients    []User
}

// ID returns the channel's ID
func (c *Channel) ID() snowflake.Snowflake {
	return c.id
}

// UpdateFromElementMap updates information about the channel
// This will not remove known data, only replace it
func (c *Channel) UpdateFromElementMap(eMap map[string]Element) (err error) {
	var ok bool
	var e2 Element
	var u User

	c.channelType, err = ChannelTypeFromElement(eMap["type"])
	if err != nil {
		err = errors.Wrap(err, "could not get channelType")
		return
	}

	e2, ok = eMap["guild_id"]
	if ok && !e2.IsNil() {
		c.guildID, err = SnowflakeFromElement(e2)
		if err != nil {
			err = errors.Wrap(err, "could not get guild_id snowflake.Snowflake")
			return
		}
	}

	e2, ok = eMap["name"]
	if ok {
		c.name, err = e2.ToString()
		if err != nil {
			err = errors.Wrap(err, "could not get name")
			return
		}
	}

	e2, ok = eMap["topic"]
	if ok {
		c.topic, err = e2.ToString()
		if err != nil {
			err = errors.Wrap(err, "could not get topic")
			return
		}
	}

	e2, ok = eMap["last_message_id"]
	if ok && !e2.IsNil() {
		c.lastMessageID, err = SnowflakeFromElement(e2)
		if err != nil {
			err = errors.Wrap(err, "could not get last_message_id snowflake.Snowflake")
			return
		}
	}

	e2, ok = eMap["parent_id"]
	if ok && !e2.IsNil() {
		c.parentID, err = SnowflakeFromElement(e2)
		if err != nil {
			err = errors.Wrap(err, "could not get parent_id snowflake.Snowflake")
			return
		}
	}

	e2, ok = eMap["owner_id"]
	if ok && !e2.IsNil() {
		c.ownerID, err = SnowflakeFromElement(e2)
		if err != nil {
			err = errors.Wrap(err, "could not get owner_id snowflake.Snowflake")
			return
		}
	}

	e2, ok = eMap["application_id"]
	if ok && !e2.IsNil() {
		c.applicationID, err = SnowflakeFromElement(e2)
		if err != nil {
			err = errors.Wrap(err, "could not get application_id snowflake.Snowflake")
			return
		}
	}

	e2, ok = eMap["recipients"]
	if ok {
		c.recipients = make([]User, 0, len(e2.Vals))
		for _, e3 := range e2.Vals {
			u, err = UserFromElement(e3)
			if err != nil {
				err = errors.Wrap(err, "could not inflate channel user")
				return
			}
			c.recipients = append(c.recipients, u)
		}
	}

	return
}

// ChannelFromElement creates a new Channel object from the given etf Element.
// The element should be a Map-type Element
func ChannelFromElement(e Element) (c Channel, err error) {
	var eMap map[string]Element
	eMap, c.id, err = MapAndIDFromElement(e)
	if err != nil {
		return
	}

	err = c.UpdateFromElementMap(eMap)
	return
}

// ChannelFromElementMap creates a new Channel object from the given data map.
func ChannelFromElementMap(eMap map[string]Element) (c Channel, err error) {
	c.id, err = SnowflakeFromElement(eMap["id"])
	if err != nil {
		err = errors.Wrap(err, "could not get channel id")
		return
	}

	err = c.UpdateFromElementMap(eMap)
	err = errors.Wrap(err, "could not create a channel")
	return
}
