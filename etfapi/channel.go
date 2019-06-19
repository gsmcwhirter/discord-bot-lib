package etfapi

import (
	"fmt"

	"github.com/gsmcwhirter/go-util/v4/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v9/snowflake"
)

// ChannelType represents the type of a discord channel
type ChannelType int

// These are the known discord channel types
const (
	GuildTextChannel     ChannelType = 0
	DMChannel            ChannelType = 1
	GuildVoiceChannel    ChannelType = 2
	GroupDMChannel       ChannelType = 3
	GuildCategoryChannel ChannelType = 4
)

// ChannelTypeFromElement extracts the channel type from a etf Element
func ChannelTypeFromElement(e Element) (ChannelType, error) {
	temp, err := e.ToInt()
	return ChannelType(temp), errors.Wrap(err, "could not unmarshal channelType")
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
func (c *Channel) UpdateFromElementMap(eMap map[string]Element) error {
	var ok bool
	var e2 Element
	var u User
	var err error

	c.channelType, err = ChannelTypeFromElement(eMap["type"])
	if err != nil {
		return errors.Wrap(err, "could not get channelType")
	}

	e2, ok = eMap["guild_id"]
	if ok && !e2.IsNil() {
		c.guildID, err = SnowflakeFromElement(e2)
		if err != nil {
			return errors.Wrap(err, "could not get guild_id snowflake.Snowflake")
		}
	}

	e2, ok = eMap["name"]
	if ok {
		c.name, err = e2.ToString()
		if err != nil {
			return errors.Wrap(err, "could not get name")
		}
	}

	e2, ok = eMap["topic"]
	if ok {
		c.topic, err = e2.ToString()
		if err != nil {
			return errors.Wrap(err, "could not get topic")
		}
	}

	e2, ok = eMap["last_message_id"]
	if ok && !e2.IsNil() {
		c.lastMessageID, err = SnowflakeFromElement(e2)
		if err != nil {
			return errors.Wrap(err, "could not get last_message_id snowflake.Snowflake")
		}
	}

	e2, ok = eMap["parent_id"]
	if ok && !e2.IsNil() {
		c.parentID, err = SnowflakeFromElement(e2)
		if err != nil {
			return errors.Wrap(err, "could not get parent_id snowflake.Snowflake")
		}
	}

	e2, ok = eMap["owner_id"]
	if ok && !e2.IsNil() {
		c.ownerID, err = SnowflakeFromElement(e2)
		if err != nil {
			return errors.Wrap(err, "could not get owner_id snowflake.Snowflake")
		}
	}

	e2, ok = eMap["application_id"]
	if ok && !e2.IsNil() {
		c.applicationID, err = SnowflakeFromElement(e2)
		if err != nil {
			return errors.Wrap(err, "could not get application_id snowflake.Snowflake")
		}
	}

	e2, ok = eMap["recipients"]
	if ok {
		c.recipients = make([]User, 0, len(e2.Vals))
		for _, e3 := range e2.Vals {
			u, err = UserFromElement(e3)
			if err != nil {
				return errors.Wrap(err, "could not inflate channel user")
			}
			c.recipients = append(c.recipients, u)
		}
	}

	return nil
}

// ChannelFromElement creates a new Channel object from the given etf Element.
// The element should be a Map-type Element
func ChannelFromElement(e Element) (Channel, error) {
	var c Channel
	var eMap map[string]Element
	var err error

	eMap, c.id, err = MapAndIDFromElement(e)
	if err != nil {
		return c, err
	}

	err = c.UpdateFromElementMap(eMap)
	return c, err
}

// ChannelFromElementMap creates a new Channel object from the given data map.
func ChannelFromElementMap(eMap map[string]Element) (Channel, error) {
	var c Channel
	var err error

	c.id, err = SnowflakeFromElement(eMap["id"])
	if err != nil {
		return c, errors.Wrap(err, "could not get channel id")
	}

	err = c.UpdateFromElementMap(eMap)
	return c, errors.Wrap(err, "could not create a channel")
}
