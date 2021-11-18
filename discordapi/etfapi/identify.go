package etfapi

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v23/discordapi"
)

// IdentifyPayloadProperties holds the data about the os, etc. of the bot when identifying
type IdentifyPayloadProperties struct {
	OS      string
	Browser string
	Device  string
}

// IdentifyPayloadShard holds the data about the shards being identified for
type IdentifyPayloadShard struct {
	ID    int
	MaxID int
}

// IdentifyPayloadGame holds the data about the "game" portion of the presence
type IdentifyPayloadGame struct {
	Name string
	Type int
}

// IdentifyPayloadPresence holds the data about the "presence" portion of the identify payload
type IdentifyPayloadPresence struct {
	Game   IdentifyPayloadGame
	Status string
	Since  int
	AFK    bool
}

// IdentifyPayload is the specialized payload for sending "Identify" events to the discord gateway websocket
type IdentifyPayload struct {
	Token          string
	Intents        int
	Properties     IdentifyPayloadProperties
	LargeThreshold int
	Shard          IdentifyPayloadShard
	Presence       IdentifyPayloadPresence
}

// Payload converts the specialized payload to a generic Payload
func (ip *IdentifyPayload) Payload() (Payload, error) {
	var err error

	p := Payload{
		OpCode: discordapi.Identify,
		Data:   map[string]Element{},
	}

	// TOKEN
	p.Data["token"], err = NewStringElement(ip.Token)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for token")
	}

	p.Data["intents"], err = NewInt32Element(ip.Intents)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for intents")
	}

	// PROPERTIES
	propMap := map[string]Element{}

	propMap["$os"], err = NewStringElement(ip.Properties.OS)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Properties OS")
	}

	propMap["$browser"], err = NewStringElement(ip.Properties.Browser)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Properties Browser")
	}

	propMap["$device"], err = NewStringElement(ip.Properties.Device)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Properties Device")
	}

	p.Data["properties"], err = NewMapElement(propMap)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Properties map")
	}

	// COMPRESS -- FALSE
	p.Data["compress"], err = NewBoolElement(false)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Compress")
	}

	// LARGE_THRESHOLD
	p.Data["large_threshold"], err = NewInt32Element(ip.LargeThreshold)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for LargeThreshold")
	}

	// SHARD
	shardData := make([]Element, 2)

	shardData[0], err = NewInt32Element(ip.Shard.ID)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Shard ID")
	}

	shardData[1], err = NewInt32Element(ip.Shard.MaxID + 1)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Shard Total")
	}

	p.Data["shard"], err = NewListElement(shardData)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Shard")
	}

	// PRESENCE
	presMap := map[string]Element{}

	presMap["status"], err = NewStringElement(ip.Presence.Status)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Presence Status")
	}

	if ip.Presence.Since > 0 {
		presMap["since"], err = NewInt32Element(ip.Presence.Since)
		if err != nil {
			return p, errors.Wrap(err, "could not create Element for Presence Since")
		}
	}

	gameMap := map[string]Element{}

	gameMap["name"], err = NewStringElement(ip.Presence.Game.Name)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Presence Game Name")
	}

	gameMap["type"], err = NewInt32Element(ip.Presence.Game.Type)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Presence Game Type")
	}

	presMap["game"], err = NewMapElement(gameMap)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Presence Game map")
	}

	presMap["afk"], err = NewBoolElement(ip.Presence.AFK)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Presence AFK")
	}

	p.Data["presence"], err = NewMapElement(presMap)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for Presence map")
	}

	return p, nil
}
