package entity

import (
	"fmt"

	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v23/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v23/snowflake"
)

// Reaction represents the data about a message in a discord channel
type Reaction struct {
	UserIDString    string `json:"user_id"`
	ChannelIDString string `json:"channel_id"`
	MessageIDString string `json:"message_id"`
	GuildIDString   string `json:"guild_id"`

	UserIDSnowflake    snowflake.Snowflake `json:"-"`
	ChannelIDSnowflake snowflake.Snowflake `json:"-"`
	MessageIDSnowflake snowflake.Snowflake `json:"-"`
	GuildIDSnowflake   snowflake.Snowflake `json:"-"`
	// member    GuildMember
	EmojiString string `json:"emoji"`
}

// UserID returns the ID of the reactor
func (r *Reaction) UserID() snowflake.Snowflake {
	return r.UserIDSnowflake
}

// ChannelID returns the ID of the channel the reaction was to
func (r *Reaction) ChannelID() snowflake.Snowflake {
	return r.ChannelIDSnowflake
}

// MessageID returns the ID of the message the reaction was to
func (r *Reaction) MessageID() snowflake.Snowflake {
	return r.MessageIDSnowflake
}

// GuildID returns the ID of the guild the reaction was to
func (r *Reaction) GuildID() snowflake.Snowflake {
	return r.GuildIDSnowflake
}

// Emoji returns the emoji of the reaction
// ChannelID returns the ID of the channel the reaction was to
func (r *Reaction) Emoji() string {
	return r.EmojiString
}

// ReactionFromElementMap generates a new Reaction object from the given data
func ReactionFromElementMap(eMap map[string]etfapi.Element) (Reaction, error) {
	var r Reaction
	var err error

	e2, ok := eMap["user_id"]
	if ok {
		r.UserIDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
		if err != nil {
			return r, errors.Wrap(err, "could not get user_id snowflake.Snowflake")
		}
		r.UserIDString = r.UserIDSnowflake.ToString()
	}

	e2, ok = eMap["channel_id"]
	if ok {
		r.ChannelIDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
		if err != nil {
			return r, errors.Wrap(err, "could not get channel_id snowflake.Snowflake")
		}
		r.ChannelIDString = r.ChannelIDSnowflake.ToString()
	}

	e2, ok = eMap["message_id"]
	if ok {
		r.MessageIDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
		if err != nil {
			return r, errors.Wrap(err, "could not get message_id snowflake.Snowflake")
		}
		r.MessageIDString = r.MessageIDSnowflake.ToString()
	}

	e2, ok = eMap["guild_id"]
	if ok && !e2.IsNil() {
		r.GuildIDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
		if err != nil {
			return r, errors.Wrap(err, "could not get guild_id snowflake.Snowflake")
		}
		r.GuildIDString = r.GuildIDSnowflake.ToString()
	}

	e2, ok = eMap["emoji"]
	if ok {
		eMap2, err := e2.ToMap()
		if err != nil {
			return r, errors.Wrap(err, "could not get emoji map")
		}

		e3, ok := eMap2["name"]
		if ok {
			r.EmojiString, err = e3.ToString()
			if err != nil {
				return r, errors.Wrap(err, "could not get emoji")
			}
		}
	}

	return r, nil
}

// ReactionFromElement generates a new Reaction object from the given Element
func ReactionFromElement(e etfapi.Element) (Reaction, error) {
	eMap, err := e.ToMap()
	if err != nil {
		return Reaction{}, errors.Wrap(err, fmt.Sprintf("could not inflate element to map: %v", e))
	}

	r, err := ReactionFromElementMap(eMap)

	return r, err
}
