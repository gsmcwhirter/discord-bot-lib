package payloads

import (
	"github.com/gsmcwhirter/discord-bot-lib/discordapi"
	"github.com/gsmcwhirter/discord-bot-lib/etfapi"
	"github.com/pkg/errors"
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
	Properties     IdentifyPayloadProperties
	LargeThreshold int
	Shard          IdentifyPayloadShard
	Presence       IdentifyPayloadPresence
}

// Payload converts the specialized payload to a generic etfapi.Payload
func (ip *IdentifyPayload) Payload() (p etfapi.Payload, err error) {
	p.OpCode = discordapi.Identify
	p.Data = map[string]etfapi.Element{}

	// TOKEN
	p.Data["token"], err = etfapi.NewStringElement(ip.Token)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for token")
		return
	}

	// PROPERTIES
	propMap := map[string]etfapi.Element{}

	propMap["$os"], err = etfapi.NewStringElement(ip.Properties.OS)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Properties OS")
		return
	}

	propMap["$browser"], err = etfapi.NewStringElement(ip.Properties.Browser)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Properties Browser")
		return
	}

	propMap["$device"], err = etfapi.NewStringElement(ip.Properties.Device)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Properties Device")
		return
	}

	p.Data["properties"], err = etfapi.NewMapElement(propMap)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Properties map")
		return
	}

	// COMPRESS -- FALSE
	p.Data["compress"], err = etfapi.NewBoolElement(false)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Compress")
		return
	}

	// LARGE_THRESHOLD
	p.Data["large_threshold"], err = etfapi.NewInt32Element(ip.LargeThreshold)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for LargeThreshold")
		return
	}

	// SHARD
	shardData := make([]etfapi.Element, 2)

	shardData[0], err = etfapi.NewInt32Element(ip.Shard.ID)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Shard ID")
		return
	}

	shardData[1], err = etfapi.NewInt32Element(ip.Shard.MaxID + 1)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Shard Total")
		return
	}

	p.Data["shard"], err = etfapi.NewListElement(shardData)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Shard")
		return
	}

	// PRESENCE
	presMap := map[string]etfapi.Element{}

	presMap["status"], err = etfapi.NewStringElement(ip.Presence.Status)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Presence Status")
		return
	}

	if ip.Presence.Since > 0 {
		presMap["since"], err = etfapi.NewInt32Element(ip.Presence.Since)
		if err != nil {
			err = errors.Wrap(err, "could not create Element for Presence Since")
			return
		}
	}

	gameMap := map[string]etfapi.Element{}

	gameMap["name"], err = etfapi.NewStringElement(ip.Presence.Game.Name)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Presence Game Name")
		return
	}

	gameMap["type"], err = etfapi.NewInt32Element(ip.Presence.Game.Type)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Presence Game Type")
		return
	}

	presMap["game"], err = etfapi.NewMapElement(gameMap)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Presence Game map")
		return
	}

	presMap["afk"], err = etfapi.NewBoolElement(ip.Presence.AFK)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Presence AFK")
		return
	}

	p.Data["presence"], err = etfapi.NewMapElement(presMap)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for Presence map")
		return
	}

	return
}