package entity

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v22/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v22/snowflake"
)

type InteractionType int

const (
	InteractionPing               InteractionType = 1
	InteractionApplicationCommand InteractionType = 2
	InteractionMessageComponent   InteractionType = 3
)

// InteractionTypeFromElement generates a InteractionType representation from the given
// interaction-type Element
func InteractionTypeFromElement(e etfapi.Element) (InteractionType, error) {
	temp, err := e.ToInt()
	t := InteractionType(temp)
	return t, errors.Wrap(err, "could not unmarshal InteractionType")
}

type Interaction struct {
	Type    InteractionType  `json:"type"`
	Data    *InteractionData `json:"data"`
	Member  *GuildMember     `json:"member"`
	User    *User            `json:"user"`
	Token   string           `json:"token"`
	Version int              `json:"version"`
	Message *Message         `json:"message"`

	IDString            string `json:"id"`
	ApplicationIDString string `json:"application_id"`
	GuildIDString       string `json:"guild_id"`
	ChannelIDString     string `json:"channel_id"`

	IDSnowflake            snowflake.Snowflake `json:"-"`
	ApplicationIDSnowflake snowflake.Snowflake `json:"-"`
	GuildIDSnowflake       snowflake.Snowflake `json:"-"`
	ChannelIDSnowflake     snowflake.Snowflake `json:"-"`
}

// InteractionFromElementMap generates a new Interaction object from the given data
func InteractionFromElementMap(eMap map[string]etfapi.Element) (Interaction, error) {
	var i Interaction
	var err error

	e2, ok := eMap["id"]
	if ok {
		i.IDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return i, errors.Wrap(err, "could not get id snowflake.Snowflake")
		}

		i.IDString = i.IDSnowflake.ToString()
	}

	e2, ok = eMap["application_id"]
	if ok {
		i.ApplicationIDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return i, errors.Wrap(err, "could not get application_id snowflake.Snowflake")
		}

		i.ApplicationIDString = i.ApplicationIDSnowflake.ToString()
	}

	e2, ok = eMap["channel_id"]
	if ok && !e2.IsNil() {
		i.ChannelIDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return i, errors.Wrap(err, "could not get channel_id snowflake.Snowflake")
		}
		i.ChannelIDString = i.ChannelIDSnowflake.ToString()
	}

	e2, ok = eMap["guild_id"]
	if ok && !e2.IsNil() {
		i.GuildIDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return i, errors.Wrap(err, "could not get guild_id snowflake.Snowflake")
		}
		i.GuildIDString = i.GuildIDSnowflake.ToString()
	}

	i.Type, err = InteractionTypeFromElement(eMap["type"])
	if err != nil {
		return i, errors.Wrap(err, "could not get interactionType")
	}

	e2 = eMap["token"]
	i.Token, err = e2.ToString()
	if err != nil {
		return i, errors.Wrap(err, "could not get token")
	}

	e2 = eMap["version"]
	i.Version, err = e2.ToInt()
	if err != nil {
		return i, errors.Wrap(err, "could not get version")
	}

	e2, ok = eMap["data"]
	if ok && !e2.IsNil() {
		v, err := InteractionDataFromElement(e2)
		if err != nil {
			return i, errors.Wrap(err, "could not inflate interaction data")
		}
		i.Data = &v
	}

	e2, ok = eMap["member"]
	if ok && !e2.IsNil() {
		v, err := GuildMemberFromElement(e2)
		if err != nil {
			return i, errors.Wrap(err, "could not inflate guild member")
		}
		i.Member = &v
	}

	e2, ok = eMap["user"]
	if ok && !e2.IsNil() {
		v, err := UserFromElement(e2)
		if err != nil {
			return i, errors.Wrap(err, "could not inflate interaction user")
		}
		i.User = &v
	}

	e2, ok = eMap["message"]
	if ok && !e2.IsNil() {
		v, err := MessageFromElement(e2)
		if err != nil {
			return i, errors.Wrap(err, "could not inflate interaction message")
		}
		i.Message = &v
	}

	return i, nil
}

type InteractionData struct {
	Name     string
	Type     ApplicationCommandType
	Resolved ResolvedData
	Options  []ApplicationCommandInteractionOption
	// CustomID string
	// ComponentType
	// Values

	IDString       string `json:"id"`
	TargetIDString string `json:"target_id"`

	IDSnowflake       snowflake.Snowflake `json:"-"`
	TargetIDSnowflake snowflake.Snowflake `json:"-"`
}

// InteractionDataFromElement generates a new Interaction object from the given data
func InteractionDataFromElement(e etfapi.Element) (InteractionData, error) {
	var d InteractionData

	eMap, err := e.ToMap()
	if err != nil {
		return d, errors.Wrap(err, "could not inflate InteractionData from non-map")
	}

	e2, ok := eMap["id"]
	if ok {
		d.IDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return d, errors.Wrap(err, "could not get id snowflake.Snowflake")
		}

		d.IDString = d.IDSnowflake.ToString()
	}

	e2, ok = eMap["target_id"]
	if ok {
		d.TargetIDSnowflake, err = etfapi.SnowflakeFromElement(e2)
		if err != nil {
			return d, errors.Wrap(err, "could not get target_id snowflake.Snowflake")
		}

		d.TargetIDString = d.TargetIDSnowflake.ToString()
	}

	e2 = eMap["name"]
	d.Name, err = e2.ToString()
	if err != nil {
		return d, errors.Wrap(err, "could not get name")
	}

	d.Type, err = ApplicationCommandTypeFromElement(eMap["type"])
	if err != nil {
		return d, errors.Wrap(err, "could not get type")
	}

	e2, ok = eMap["resolved"]
	if ok && !e2.IsNil() {
		d.Resolved, err = ResolvedDataFromElement(e2)
		if err != nil {
			return d, errors.Wrap(err, "could not get resolved data")
		}
	}

	e2, ok = eMap["options"]
	if ok && !e2.IsNil() {
		el, err := e2.ToList()
		if err != nil {
			return d, errors.Wrap(err, "options was not a list")
		}

		d.Options = make([]ApplicationCommandInteractionOption, 0, len(el))
		for _, e3 := range el {
			o, err := ApplicationCommandInteractionOptionFromElement(e3)
			if err != nil {
				return d, errors.Wrap(err, "could not inflate option")
			}
			d.Options = append(d.Options, o)
		}
	}

	return d, nil
}

type ResolvedData struct {
	Users    map[snowflake.Snowflake]User
	Members  map[snowflake.Snowflake]GuildMember
	Roles    map[snowflake.Snowflake]Role
	Channels map[snowflake.Snowflake]Channel
}

// ResolvedDataFromElement generates a new Interaction object from the given data
func ResolvedDataFromElement(e etfapi.Element) (ResolvedData, error) {
	var d ResolvedData
	var m map[snowflake.Snowflake]etfapi.Element

	eMap, err := e.ToMap()
	if err != nil {
		return d, errors.Wrap(err, "could not inflate ResolvedData from non-map")
	}

	e2, ok := eMap["users"]
	if ok && !e2.IsNil() {
		m, err = e2.ToSnowflakeMap()
		if err != nil {
			return d, errors.Wrap(err, "could not inflate users map")
		}

		d.Users = make(map[snowflake.Snowflake]User, len(m))
		for k, v := range m {
			if !v.IsNil() {
				d.Users[k], err = UserFromElement(v)
				if err != nil {
					return d, errors.Wrap(err, "could not inflate user")
				}
			}
		}
	}

	e2, ok = eMap["members"]
	if ok && !e2.IsNil() {
		m, err = e2.ToSnowflakeMap()
		if err != nil {
			return d, errors.Wrap(err, "could not inflate members map")
		}

		d.Members = make(map[snowflake.Snowflake]GuildMember, len(m))
		for k, v := range m {
			if !v.IsNil() {
				d.Members[k], err = GuildMemberFromElement(v)
				if err != nil {
					return d, errors.Wrap(err, "could not inflate member")
				}
			}
		}
	}

	e2, ok = eMap["roles"]
	if ok && !e2.IsNil() {
		m, err = e2.ToSnowflakeMap()
		if err != nil {
			return d, errors.Wrap(err, "could not inflate roles map")
		}

		d.Roles = make(map[snowflake.Snowflake]Role, len(m))
		for k, v := range m {
			if !v.IsNil() {
				d.Roles[k], err = RoleFromElement(v)
				if err != nil {
					return d, errors.Wrap(err, "could not inflate role")
				}
			}
		}
	}

	e2, ok = eMap["channels"]
	if ok && !e2.IsNil() {
		m, err = e2.ToSnowflakeMap()
		if err != nil {
			return d, errors.Wrap(err, "could not inflate channels map")
		}

		d.Channels = make(map[snowflake.Snowflake]Channel, len(m))
		for k, v := range m {
			if !v.IsNil() {
				d.Channels[k], err = ChannelFromElement(v)
				if err != nil {
					return d, errors.Wrap(err, "could not inflate channel")
				}
			}
		}
	}

	return d, nil
}
