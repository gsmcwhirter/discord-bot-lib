package entity

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v22/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v22/snowflake"
)

type ChannelType int

const (
	ChannelGuildText          ChannelType = 0
	ChannelDM                 ChannelType = 1
	ChannelGuildVoice         ChannelType = 2
	ChannelGroupDM            ChannelType = 3
	ChannelGuildCategory      ChannelType = 4
	ChannelGuildNews          ChannelType = 5
	ChannelGuildStore         ChannelType = 6
	ChannelGuildNewsThread    ChannelType = 10
	ChannelGuildPublicThread  ChannelType = 11
	ChannelGuildPrivateThread ChannelType = 12
	ChannelGuildStageVoice    ChannelType = 13
)

// ChannelTypeFromElement generates a ChannelType representation from the given
// channel-type Element
func ChannelTypeFromElement(e etfapi.Element) (ChannelType, error) {
	temp, err := e.ToInt()
	t := ChannelType(temp)
	return t, errors.Wrap(err, "could not unmarshal ChannelType")
}

type Channel struct {
	IDString            string `json:"id"`
	GuildIDString       string `json:"guild_id"`
	OwnerIDString       string `json:"owner_id"`
	ApplicationIDString string `json:"application_id"`
	LastMessageIDString string `json:"last_message_id"`
	ParentIDString      string `json:"parent_id"`
	Permissions         string `json:"permissions"`

	Type                 ChannelType           `json:"type"`
	Position             int                   `json:"position"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	Name                 string                `json:"name"`
	Topic                string                `json:"topic"`
	NSFW                 bool                  `json:"nsfw"`
	Bitrate              int                   `json:"bitrate"`
	UserLimit            int                   `json:"user_limit"`
	Recipients           []User                `json:"recipients"`
	Icon                 string                `json:"icon"`
	LastPinTimestamp     string                `json:"last_pin_timestamp"`
	RTCRegion            string                `json:"rtc_region"`
	VideoQualityMode     int                   `json:"video_quality_mode"`
	MessageCount         int                   `json:"message_count"`
	MemberCount          int                   `json:"member_count"`

	IDSnowflake            snowflake.Snowflake `json:"-"`
	GuildIDSnowflake       snowflake.Snowflake `json:"-"`
	LastMessageIDSnowflake snowflake.Snowflake `json:"-"`
	OwnerIDSnowflake       snowflake.Snowflake `json:"-"`
	ApplicationIDSnowflake snowflake.Snowflake `json:"-"`
	ParentIDSnowflake      snowflake.Snowflake `json:"-"`
	PermissionsSnowflake   snowflake.Snowflake `json:"-"`
}

func ChannelFromElement(e etfapi.Element) (Channel, error) {
	var c Channel

	eMap, err := e.ToMap()
	if err != nil {
		return c, errors.Wrap(err, "could not inflate Channel from map")
	}

	e2, ok := eMap["id"]
	if ok {
		c.IDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return c, errors.Wrap(err, "could not get id snowflake.Snowflake")
		}

		c.IDString = c.IDSnowflake.ToString()
	}

	e2, ok = eMap["guild_id"]
	if ok {
		c.GuildIDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return c, errors.Wrap(err, "could not get guild_id snowflake.Snowflake")
		}

		c.GuildIDString = c.GuildIDSnowflake.ToString()
	}

	e2, ok = eMap["last_message_id"]
	if ok {
		c.LastMessageIDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return c, errors.Wrap(err, "could not get last_message_id snowflake.Snowflake")
		}

		c.LastMessageIDString = c.LastMessageIDSnowflake.ToString()
	}

	e2, ok = eMap["owner_id"]
	if ok {
		c.OwnerIDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return c, errors.Wrap(err, "could not get owner_id snowflake.Snowflake")
		}

		c.OwnerIDString = c.OwnerIDSnowflake.ToString()
	}

	e2, ok = eMap["application_id"]
	if ok {
		c.ApplicationIDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return c, errors.Wrap(err, "could not get application_id snowflake.Snowflake")
		}

		c.ApplicationIDString = c.ApplicationIDSnowflake.ToString()
	}

	e2, ok = eMap["parent_id"]
	if ok {
		c.ParentIDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return c, errors.Wrap(err, "could not get parent_id snowflake.Snowflake")
		}

		c.ParentIDString = c.ParentIDSnowflake.ToString()
	}

	e2, ok = eMap["permissions"]
	if ok {
		c.PermissionsSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return c, errors.Wrap(err, "could not get guild_id snowflake.Snowflake")
		}

		c.Permissions = c.PermissionsSnowflake.ToString()
	}

	e2 = eMap["type"]
	c.Type, err = ChannelTypeFromElement(e2)
	if err != nil {
		return c, errors.Wrap(err, "could not inflate channel type")
	}

	e2, ok = eMap["position"]
	if ok {
		c.Position, err = e2.ToInt()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate position")
		}
	}

	e2, ok = eMap["name"]
	if ok {
		c.Name, err = e2.ToString()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate name")
		}
	}

	e2, ok = eMap["topic"]
	if ok {
		c.Topic, err = e2.ToString()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate topic")
		}
	}

	e2, ok = eMap["nsfw"]
	if ok {
		c.NSFW, err = e2.ToBool()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate nsfw")
		}
	}

	e2, ok = eMap["bitrate"]
	if ok {
		c.Bitrate, err = e2.ToInt()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate bitrate")
		}
	}

	e2, ok = eMap["user_limit"]
	if ok {
		c.UserLimit, err = e2.ToInt()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate user_limit")
		}
	}

	e2, ok = eMap["icon"]
	if ok {
		c.Icon, err = e2.ToString()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate icon")
		}
	}

	e2, ok = eMap["last_pin_timestamp"]
	if ok {
		c.LastPinTimestamp, err = e2.ToString()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate last_pin_timestamp")
		}
	}

	e2, ok = eMap["rtc_region"]
	if ok {
		c.RTCRegion, err = e2.ToString()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate rtc_region")
		}
	}

	e2, ok = eMap["video_quality_mode"]
	if ok {
		c.VideoQualityMode, err = e2.ToInt()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate video_quality_mode")
		}
	}

	e2, ok = eMap["message_count"]
	if ok {
		c.MessageCount, err = e2.ToInt()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate message_count")
		}
	}

	e2, ok = eMap["member_count"]
	if ok {
		c.MemberCount, err = e2.ToInt()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate member_count")
		}
	}

	e2, ok = eMap["recipients"]
	if ok && !e2.IsNil() {
		el, err := e2.ToList()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate recipients list")
		}

		c.Recipients = make([]User, 0, len(el))
		for _, e3 := range el {
			u, err := UserFromElement(e3)
			if err != nil {
				return c, errors.Wrap(err, "could not inflate recipient")
			}

			c.Recipients = append(c.Recipients, u)
		}
	}

	e2, ok = eMap["permission_overwrites"]
	if ok && !e2.IsNil() {
		el, err := e2.ToList()
		if err != nil {
			return c, errors.Wrap(err, "could not inflate permission_overwrites list")
		}

		c.PermissionOverwrites = make([]PermissionOverwrite, 0, len(el))
		for _, e3 := range el {
			ow, err := PermissionOverwriteFromElement(e3)
			if err != nil {
				return c, errors.Wrap(err, "could not inflate permission overwrite")
			}

			c.PermissionOverwrites = append(c.PermissionOverwrites, ow)
		}
	}

	return c, nil
}

type PermissionOverwrite struct {
	IDString    string `json:"id"`
	Type        int    `json:"type"`
	AllowString string `json:"allow"`
	DenyString  string `json:"deny"`

	IDSnowflake    snowflake.Snowflake `json:"-"`
	AllowSnowflake snowflake.Snowflake `json:"-"`
	DenySnowflake  snowflake.Snowflake `json:"-"`
}

func PermissionOverwriteFromElement(e etfapi.Element) (PermissionOverwrite, error) {
	var ow PermissionOverwrite

	eMap, err := e.ToMap()
	if err != nil {
		return ow, errors.Wrap(err, "could not inflate PermissionOverwrite from map")
	}

	e2, ok := eMap["id"]
	if ok {
		ow.IDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return ow, errors.Wrap(err, "could not get id snowflake.Snowflake")
		}

		ow.IDString = ow.IDSnowflake.ToString()
	}

	e2, ok = eMap["allow"]
	if ok {
		ow.AllowSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return ow, errors.Wrap(err, "could not get allow snowflake.Snowflake")
		}

		ow.AllowString = ow.AllowSnowflake.ToString()
	}

	e2, ok = eMap["deny"]
	if ok {
		ow.DenySnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return ow, errors.Wrap(err, "could not get deny snowflake.Snowflake")
		}

		ow.DenyString = ow.DenySnowflake.ToString()
	}

	e2 = eMap["type"]
	ow.Type, err = e2.ToInt()
	if err != nil {
		return ow, errors.Wrap(err, "could not inflate type")
	}

	return ow, nil
}
