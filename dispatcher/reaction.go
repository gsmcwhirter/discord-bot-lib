package dispatcher

import (
	"fmt"

	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v19/discordapi/etf"
	"github.com/gsmcwhirter/discord-bot-lib/v19/snowflake"
)

// Message represents the data about a message in a discord channel
type Reaction struct {
	userID    snowflake.Snowflake
	channelID snowflake.Snowflake
	messageID snowflake.Snowflake
	guildID   snowflake.Snowflake
	// member    GuildMember
	emoji string
}

// UserID returns the ID of the reactor
func (r *Reaction) UserID() snowflake.Snowflake {
	return r.userID
}

// ChannelID returns the ID of the channel the reaction was to
func (r *Reaction) ChannelID() snowflake.Snowflake {
	return r.channelID
}

// MessageID returns the ID of the message the reaction was to
func (r *Reaction) MessageID() snowflake.Snowflake {
	return r.messageID
}

// GuildID returns the ID of the guild the reaction was to
func (r *Reaction) GuildID() snowflake.Snowflake {
	return r.guildID
}

// Emoji returns the emoji of the reaction
// ChannelID returns the ID of the channel the reaction was to
func (r *Reaction) Emoji() string {
	return r.emoji
}

// ReactionFromElementMap generates a new Reaction object from the given data
func ReactionFromElementMap(eMap map[string]etf.Element) (Reaction, error) {
	var r Reaction
	var err error

	e2, ok := eMap["user_id"]
	if ok {
		r.userID, err = etf.SnowflakeFromElement(e2)
		if err != nil {
			return r, errors.Wrap(err, "could not get user_id snowflake.Snowflake")
		}
	}

	e2, ok = eMap["channel_id"]
	if ok {
		r.channelID, err = etf.SnowflakeFromElement(e2)
		if err != nil {
			return r, errors.Wrap(err, "could not get channel_id snowflake.Snowflake")
		}
	}

	e2, ok = eMap["message_id"]
	if ok {
		r.messageID, err = etf.SnowflakeFromElement(e2)
		if err != nil {
			return r, errors.Wrap(err, "could not get message_id snowflake.Snowflake")
		}
	}

	e2, ok = eMap["guild_id"]
	if ok && !e2.IsNil() {
		r.guildID, err = etf.SnowflakeFromElement(e2)
		if err != nil {
			return r, errors.Wrap(err, "could not get guild_id snowflake.Snowflake")
		}
	}

	e2, ok = eMap["emoji"]
	if ok {
		eMap2, err := e2.ToMap()
		if err != nil {
			return r, errors.Wrap(err, "could not get emoji map")
		}

		e3, ok := eMap2["name"]
		if ok {
			r.emoji, err = e3.ToString()
			if err != nil {
				return r, errors.Wrap(err, "could not get emoji")
			}
		}
	}

	return r, nil
}

// ReactionFromElement generates a new Reaction object from the given Element
func ReactionFromElement(e etf.Element) (Reaction, error) {
	eMap, err := e.ToMap()
	if err != nil {
		return Reaction{}, errors.Wrap(err, fmt.Sprintf("could not inflate element to map: %v", e))
	}

	r, err := ReactionFromElementMap(eMap)

	return r, err
}
