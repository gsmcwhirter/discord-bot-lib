package entity

import (
	stdjson "encoding/json" //nolint:depguard // we need this for RawMessage

	"github.com/gsmcwhirter/go-util/v8/errors"
	"github.com/gsmcwhirter/go-util/v8/json"

	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
)

// InteractionType represents the type of an interaction
type InteractionType int

// These are the known InteractionType values
const (
	InteractionPing               InteractionType = 1
	InteractionApplicationCommand InteractionType = 2
	InteractionMessageComponent   InteractionType = 3
	InteractionAutocomplete       InteractionType = 4
)

// InteractionTypeFromElement generates a InteractionType representation from the given
// interaction-type Element
func InteractionTypeFromElement(e etfapi.Element) (InteractionType, error) {
	temp, err := e.ToInt()
	t := InteractionType(temp)
	return t, errors.Wrap(err, "could not unmarshal InteractionType", "raw", e.Val)
}

// Interaction represents a new slash-command usage instance
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
		i.IDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
		if err != nil {
			return i, errors.Wrap(err, "could not get id snowflake.Snowflake")
		}

		i.IDString = i.IDSnowflake.ToString()
	}

	e2, ok = eMap["application_id"]
	if ok {
		i.ApplicationIDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
		if err != nil {
			return i, errors.Wrap(err, "could not get application_id snowflake.Snowflake")
		}

		i.ApplicationIDString = i.ApplicationIDSnowflake.ToString()
	}

	e2, ok = eMap["channel_id"]
	if ok && !e2.IsNil() {
		i.ChannelIDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
		if err != nil {
			return i, errors.Wrap(err, "could not get channel_id snowflake.Snowflake")
		}
		i.ChannelIDString = i.ChannelIDSnowflake.ToString()
	}

	e2, ok = eMap["guild_id"]
	if ok && !e2.IsNil() {
		i.GuildIDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
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

// InteractionData represents the data provided in the interaction
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
		d.IDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
		if err != nil {
			return d, errors.Wrap(err, "could not get id snowflake.Snowflake")
		}

		d.IDString = d.IDSnowflake.ToString()
	}

	e2, ok = eMap["target_id"]
	if ok {
		d.TargetIDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
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

// ResolvedData is the resolved references from entities in the InteractionData
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

// ApplicationCommandInteractionOption represents ?
type ApplicationCommandInteractionOption struct {
	Name    string                                `json:"name"`
	Type    ApplicationCommandOptionType          `json:"type"`
	Value   stdjson.RawMessage                    `json:"value"`
	Options []ApplicationCommandInteractionOption `json:"options"`
	Focused bool                                  `json:"focused,omitempty"`

	ValueSubCommand      string              `json:"-"`
	ValueSubCommandGroup string              `json:"-"`
	ValueString          string              `json:"-"`
	ValueInt             int                 `json:"-"`
	ValueBool            bool                `json:"-"`
	ValueUser            snowflake.Snowflake `json:"-"`
	ValueChannel         snowflake.Snowflake `json:"-"`
	ValueRole            snowflake.Snowflake `json:"-"`
	ValueNumber          float64             `json:"-"`
	// TODO: ValueMentionable
}

// ResolveValue unmarshals the value into the appropriate field
func (i *ApplicationCommandInteractionOption) ResolveValue() error {
	var s string
	var err error

	switch i.Type {
	case OptTypeSubCommand:
		return json.Unmarshal([]byte(i.Value), &i.ValueSubCommand)
	case OptTypeSubCommandGroup:
		return json.Unmarshal([]byte(i.Value), &i.ValueSubCommandGroup)
	case OptTypeString:
		return json.Unmarshal([]byte(i.Value), &i.ValueString)
	case OptTypeInteger:
		return json.Unmarshal([]byte(i.Value), &i.ValueInt)
	case OptTypeBoolean:
		return json.Unmarshal([]byte(i.Value), &i.ValueBool)
	case OptTypeUser:
		if err := json.Unmarshal([]byte(i.Value), &s); err != nil {
			return errors.Wrap(err, "could not unmarshal User id to string")
		}
		i.ValueUser, err = snowflake.FromString(s)
		return errors.Wrap(err, "could not snowflakify User id string")
	case OptTypeRole:
		if err := json.Unmarshal([]byte(i.Value), &s); err != nil {
			return errors.Wrap(err, "could not unmarshal Role id to string")
		}
		i.ValueRole, err = snowflake.FromString(s)
		return errors.Wrap(err, "could not snowflakify Role id string")
	case OptTypeChannel:
		if err := json.Unmarshal([]byte(i.Value), &s); err != nil {
			return errors.Wrap(err, "could not unmarshal Channel id to string")
		}
		i.ValueChannel, err = snowflake.FromString(s)
		return errors.Wrap(err, "could not snowflakify Channel id string")
	case OptTypeMentionable:
		return ErrBadOptType // TODO: implement
	case OptTypeNumber:
		return json.Unmarshal([]byte(i.Value), &i.ValueNumber)
	default:
		return ErrBadOptType
	}
}

// PackValue marshals the value in the appropriate manner
func (i *ApplicationCommandInteractionOption) PackValue() error {
	var b []byte
	var err error

	switch i.Type {
	case OptTypeSubCommand:
		b, err = json.Marshal(i.ValueSubCommand)
	case OptTypeSubCommandGroup:
		b, err = json.Marshal(i.ValueSubCommandGroup)
	case OptTypeString:
		b, err = json.Marshal(i.ValueString)
	case OptTypeInteger:
		b, err = json.Marshal(i.ValueInt)
	case OptTypeBoolean:
		b, err = json.Marshal(i.ValueBool)
	case OptTypeUser:
		b, err = json.Marshal(i.ValueUser)
	case OptTypeRole:
		b, err = json.Marshal(i.ValueRole)
	case OptTypeChannel:
		b, err = json.Marshal(i.ValueChannel)
	case OptTypeMentionable:
		return ErrBadOptType // TODO: implement
	case OptTypeNumber:
		b, err = json.Marshal(i.ValueNumber)
	default:
		return ErrBadOptType
	}

	i.Value = stdjson.RawMessage(b)
	return err
}

// ApplicationCommandInteractionOptionFromElement instantiates an ApplicationCommandInteractionOption from an etf element
func ApplicationCommandInteractionOptionFromElement(e etfapi.Element) (ApplicationCommandInteractionOption, error) {
	var o ApplicationCommandInteractionOption

	eMap, err := e.ToMap()
	if err != nil {
		return o, errors.Wrap(err, "could not inflate ApplicationCommandInteractionOption from non-map")
	}

	e2, ok := eMap["options"]
	if ok && !e2.IsNil() {
		el, err := e2.ToList()
		if err != nil {
			return o, errors.Wrap(err, "could not inflate Options")
		}

		o.Options = make([]ApplicationCommandInteractionOption, 0, len(el))

		for _, e3 := range el {
			o2, err := ApplicationCommandInteractionOptionFromElement(e3)
			if err != nil {
				return o, errors.Wrap(err, "could not inflate sub-option")
			}

			o.Options = append(o.Options, o2)
		}
	}

	e2 = eMap["name"]
	o.Name, err = e2.ToString()
	if err != nil {
		return o, errors.Wrap(err, "could not inflate name")
	}

	e2 = eMap["type"]
	o.Type, err = ApplicationCommandOptionTypeFromElement(e2)
	if err != nil {
		return o, errors.Wrap(err, "could not inflate type")
	}

	e2, ok = eMap["focused"]
	if ok {
		o.Focused, err = e2.ToBool()
		if err != nil {
			return o, errors.Wrap(err, "could not inflate focused")
		}
	}

	e2, ok = eMap["value"]
	if ok && !e2.IsNil() {
		switch o.Type {
		case OptTypeSubCommand:
			o.ValueSubCommand, err = e2.ToString()
		case OptTypeSubCommandGroup:
			o.ValueSubCommandGroup, err = e2.ToString()
		case OptTypeString:
			o.ValueString, err = e2.ToString()
		case OptTypeInteger:
			o.ValueInt, err = e2.ToInt()
		case OptTypeBoolean:
			o.ValueBool, err = e2.ToBool()
		case OptTypeUser:
			o.ValueUser, err = etfapi.SnowflakeFromUnknownElement(e2)
		case OptTypeRole:
			o.ValueRole, err = etfapi.SnowflakeFromUnknownElement(e2)
		case OptTypeChannel:
			o.ValueChannel, err = etfapi.SnowflakeFromUnknownElement(e2)
		case OptTypeMentionable:
			err = ErrBadOptType
		case OptTypeNumber:
			o.ValueNumber, err = e2.ToFloat64()
		default:
			err = ErrBadOptType
		}

		if err != nil {
			return o, errors.Wrap(err, "could not inflate value", "raw", e2.Val, "type", o.Type)
		}
	}

	if err := o.PackValue(); err != nil {
		return o, errors.Wrap(err, "could not pack value")
	}

	return o, nil
}
