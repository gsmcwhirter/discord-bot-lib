package entity

import (
	"bytes"
	stdjson "encoding/json" //nolint:depguard // we need this for RawMessage

	"github.com/gsmcwhirter/go-util/v8/errors"
	"github.com/gsmcwhirter/go-util/v8/json"

	"github.com/gsmcwhirter/discord-bot-lib/v22/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v22/snowflake"
)

type ApplicationCommandOptionType int

const (
	OptTypeSubCommand      ApplicationCommandOptionType = 1
	OptTypeSubCommandGroup ApplicationCommandOptionType = 2
	OptTypeString          ApplicationCommandOptionType = 3
	OptTypeInteger         ApplicationCommandOptionType = 4
	OptTypeBoolean         ApplicationCommandOptionType = 5
	OptTypeUser            ApplicationCommandOptionType = 6
	OptTypeChannel         ApplicationCommandOptionType = 7
	OptTypeRole            ApplicationCommandOptionType = 8
	OptTypeMentionable     ApplicationCommandOptionType = 9
	OptTypeNumber          ApplicationCommandOptionType = 10
)

// ApplicationCommandOptionTypeFromElement generates a ApplicationCommandOptionType representation from the given
// application-command-option-type Element
func ApplicationCommandOptionTypeFromElement(e etfapi.Element) (ApplicationCommandOptionType, error) {
	temp, err := e.ToInt()
	t := ApplicationCommandOptionType(temp)
	return t, errors.Wrap(err, "could not unmarshal ApplicationCommandOptionType")
}

type ApplicationCommandType int

const (
	CmdTypeChatInput ApplicationCommandType = 1
	CmdTypeUser      ApplicationCommandType = 2
	CmdTypeMessage   ApplicationCommandType = 3
)

// ApplicationCommandTypeFromElement generates a ApplicationCommandType representation from the given
// application-command-type Element
func ApplicationCommandTypeFromElement(e etfapi.Element) (ApplicationCommandType, error) {
	temp, err := e.ToInt()
	t := ApplicationCommandType(temp)
	return t, errors.Wrap(err, "could not unmarshal ApplicationCommandType")
}

var ErrBadOptType = errors.New("bad option type value")

type ApplicationCommand struct {
	ID                string                     `json:"id,omitempty"`
	Type              ApplicationCommandType     `json:"type"`
	ApplicationID     string                     `json:"application_id,omitempty"`
	GuildID           string                     `json:"guild_id,omitempty"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	Options           []ApplicationCommandOption `json:"options"`
	DefaultPermission bool                       `json:"default_permission"`
	Version           string                     `json:"version,omitempty"`

	IDSnowflake            snowflake.Snowflake `json:"-"`
	ApplicationIDSnowflake snowflake.Snowflake `json:"-"`
	GuildIDSnowflake       snowflake.Snowflake `json:"-"`
	VersionSnowflake       snowflake.Snowflake `json:"-"`
}

func (c *ApplicationCommand) Snowflakify() error {
	var err error

	if c.ID != "" {
		if c.IDSnowflake, err = snowflake.FromString(c.ID); err != nil {
			return errors.Wrap(err, "could not snowflakify ID")
		}
	}

	if c.GuildID != "" {
		if c.GuildIDSnowflake, err = snowflake.FromString(c.GuildID); err != nil {
			return errors.Wrap(err, "could not snowflakify GuildID")
		}
	}

	if c.ApplicationID != "" {
		if c.ApplicationIDSnowflake, err = snowflake.FromString(c.ApplicationID); err != nil {
			return errors.Wrap(err, "could not snowflakify ApplicationID")
		}
	}

	if c.Version != "" {
		if c.VersionSnowflake, err = snowflake.FromString(c.Version); err != nil {
			return errors.Wrap(err, "could not snowflakify Version")
		}
	}

	for _, opt := range c.Options {
		if err = opt.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify all Options")
		}
	}

	return nil
}

type ApplicationCommandOption struct {
	Type         ApplicationCommandOptionType     `json:"type"`
	Name         string                           `json:"name"`
	Description  string                           `json:"description"`
	Required     bool                             `json:"required"`
	Choices      []ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options      []ApplicationCommandOption       `json:"options,omitempty"`
	ChannelTypes int                              `json:"channel_types,omitempty"`
}

func (o *ApplicationCommandOption) Snowflakify() error {
	var err error

	for _, opt := range o.Options {
		if err = opt.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify all Options")
		}
	}

	return nil
}

type ApplicationCommandOptionChoice struct {
	Name  string             `json:"name"`
	Value stdjson.RawMessage `json:"value"`

	Type        ApplicationCommandOptionType `json:"-"`
	ValueString string                       `json:"-"`
	ValueInt    int                          `json:"-"`
	ValueNumber float64                      `json:"-"`
}

func (c *ApplicationCommandOptionChoice) MarshalJSON() ([]byte, error) {
	var b2 []byte
	var err error

	if err = c.FillValue(); err != nil {
		return nil, errors.Wrap(err, "could not FillValue")
	}

	b := &bytes.Buffer{}
	if _, err = b.WriteString(`{"name":`); err != nil {
		return nil, errors.Wrap(err, "could not write to buffer")
	}

	b2, err = json.Marshal(c.Name)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal name")
	}

	if _, err = b.Write(b2); err != nil {
		return nil, errors.Wrap(err, "could not write to buffer")
	}

	if _, err = b.WriteString(`,"value":`); err != nil {
		return nil, errors.Wrap(err, "could not write to buffer")
	}

	b2, err = json.Marshal(c.Value)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal value")
	}

	if _, err = b.Write(b2); err != nil {
		return nil, errors.Wrap(err, "could not write to buffer")
	}

	if _, err = b.WriteString(`}`); err != nil {
		return nil, errors.Wrap(err, "could not write to buffer")
	}

	return b.Bytes(), nil
}

func (c *ApplicationCommandOptionChoice) ResolveValue() error {
	switch c.Type {
	case OptTypeString:
		return json.Unmarshal([]byte(c.Value), &c.ValueString)
	case OptTypeInteger:
		return json.Unmarshal([]byte(c.Value), &c.ValueInt)
	case OptTypeNumber:
		return json.Unmarshal([]byte(c.Value), &c.ValueNumber)
	default:
		return ErrBadOptType
	}
}

func (c *ApplicationCommandOptionChoice) FillValue() error {
	var b []byte
	var err error

	switch c.Type {
	case OptTypeInteger:
		b, err = json.Marshal(c.ValueInt)
		if err != nil {
			return errors.Wrap(err, "could not marshal ValueInt")
		}
	case OptTypeNumber:
		b, err = json.Marshal(c.ValueNumber)
		if err != nil {
			return errors.Wrap(err, "could not marshal ValueNumber")
		}
	case OptTypeString:
		b, err = json.Marshal(c.ValueString)
		if err != nil {
			return errors.Wrap(err, "could not marshal ValueString")
		}
	default:
		return ErrBadOptType
	}

	return errors.Wrap(c.Value.UnmarshalJSON(b), "could not unmarshal to RawMessage")
}

type ApplicationCommandInteractionOption struct {
	Name    string                                `json:"name"`
	Type    ApplicationCommandOptionType          `json:"type"`
	Value   stdjson.RawMessage                    `json:"value"`
	Options []ApplicationCommandInteractionOption `json:"options"`
	Focused bool                                  `json:"focused,omitempty"`

	ValueSubCommand      string   `json:"-"`
	ValueSubCommandGroup string   `json:"-"`
	ValueString          string   `json:"-"`
	ValueInt             int      `json:"-"`
	ValueBool            bool     `json:"-"`
	ValueUser            *User    `json:"-"`
	ValueChannel         *Channel `json:"-"`
	ValueRole            *Role    `json:"-"`
	ValueNumber          float64  `json:"-"`
	// TODO: ValueMentionable
}

func (i *ApplicationCommandInteractionOption) ResolveValue() error {
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
		i.ValueUser = new(User)
		return json.Unmarshal([]byte(i.Value), i.ValueUser)
	case OptTypeRole:
		i.ValueRole = new(Role)
		return json.Unmarshal([]byte(i.Value), i.ValueRole)
	case OptTypeChannel:
		i.ValueChannel = new(Channel)
		return json.Unmarshal([]byte(i.Value), &i.ValueChannel)
	case OptTypeMentionable:
		return ErrBadOptType // TODO
	case OptTypeNumber:
		return json.Unmarshal([]byte(i.Value), &i.ValueNumber)
	default:
		return ErrBadOptType
	}
}

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
		return ErrBadOptType // TODO
	case OptTypeNumber:
		b, err = json.Marshal(i.ValueNumber)
	default:
		return ErrBadOptType
	}

	i.Value = stdjson.RawMessage(b)
	return err
}

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

	e2 = eMap["value"]
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
		if !e2.IsNil() {
			v, err2 := UserFromElement(e2)
			err = err2
			o.ValueUser = &v
		}
	case OptTypeRole:
		if !e2.IsNil() {
			v, err2 := RoleFromElement(e2)
			err = err2
			o.ValueRole = &v
		}
	case OptTypeChannel:
		if !e2.IsNil() {
			v, err2 := ChannelFromElement(e2)
			err = err2
			o.ValueChannel = &v
		}
	case OptTypeMentionable:
		err = ErrBadOptType
	case OptTypeNumber:
		o.ValueNumber, err = e2.ToFloat64()
	default:
		err = ErrBadOptType
	}

	if err != nil {
		return o, errors.Wrap(err, "could not inflate value")
	}

	if err := o.PackValue(); err != nil {
		return o, errors.Wrap(err, "could not pack value")
	}

	return o, nil
}
