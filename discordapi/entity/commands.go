package entity

import (
	"bytes"
	stdjson "encoding/json" //nolint:depguard // we need this for RawMessage

	"github.com/gsmcwhirter/go-util/v8/errors"
	"github.com/gsmcwhirter/go-util/v8/json"

	"github.com/gsmcwhirter/discord-bot-lib/v23/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v23/snowflake"
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
		// fmt.Printf("snowflaking option i=%d total=%d\n", i, len(c.Options))
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
	ChannelTypes []ChannelType                    `json:"channel_types,omitempty"`
	Autocomplete bool                             `json:"autocomplete,omitempty"`
}

func (o *ApplicationCommandOption) Snowflakify() error {
	var err error

	for _, opt := range o.Options {
		// fmt.Printf("snowflaking sub-option i=%d total=%d\n", i, len(o.Options))
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
