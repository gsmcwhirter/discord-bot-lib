package etfapi

import (
	"bytes"
	"fmt"

	"github.com/gsmcwhirter/discord-bot-lib/discordapi/constants"
	"github.com/pkg/errors"
)

// Payload TODOC
type Payload struct {
	OpCode    constants.OpCode
	SeqNum    *int
	EventName string

	Data     map[string]Element
	DataList []Element
}

func (p Payload) String() string {
	if p.DataList != nil {
		return fmt.Sprintf("Payload{OpCode: %v, DataList: %v, SeqNum: %v, EventName: %v}", p.OpCode, p.DataList, p.SeqNum, p.EventName)
	}
	return fmt.Sprintf("Payload{OpCode: %v, Data: %+v, SeqNum: %v, EventName: %v}", p.OpCode, p.Data, p.SeqNum, p.EventName)
}

var opElement Element
var dElement Element
var sElement Element

func init() {
	var err error

	opElement, err = NewStringElement("op")
	if err != nil {
		panic(err)
	}

	dElement, err = NewStringElement("d")
	if err != nil {
		panic(err)
	}

	sElement, err = NewStringElement("s")
	if err != nil {
		panic(err)
	}
}

// Marshal code

// Marshal TODOC
func (p *Payload) Marshal() ([]byte, error) {
	var data []byte
	var err error

	b := bytes.Buffer{}
	b.WriteByte(131)

	len := 2
	if p.SeqNum != nil {
		len++
	}

	err = b.WriteByte(byte(Map))
	if err != nil {
		return nil, errors.Wrap(err, "unable to write outer map label")
	}
	err = writeLength32(&b, len)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write outer map length")
	}

	_, err = opElement.WriteTo(&b)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write 'op' key")
	}
	data, err = intToInt8Slice(int(p.OpCode))
	if err != nil {
		return nil, errors.Wrap(err, "unable to convert 'op' value to byte slice")
	}
	data, err = marshalInterface(Int8, data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal 'op' value")
	}
	_, err = b.Write(data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write 'op' value")
	}

	_, err = dElement.WriteTo(&b)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write 'd' key")
	}
	e2, err := ElementMapToElementSlice(p.Data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to slicify 'd' value")
	}
	data, err = marshalInterface(Map, e2)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal 'd' value")
	}
	_, err = b.Write(data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write 'd' value")
	}

	if p.SeqNum != nil {
		_, err = sElement.WriteTo(&b)
		if err != nil {
			return nil, errors.Wrap(err, "unable to write 's' key")
		}
		data, err = intToInt32Slice(*p.SeqNum)
		if err != nil {
			return nil, errors.Wrap(err, "unable to convert 's' value to byte slice")
		}
		data, err = marshalInterface(Int32, data)
		if err != nil {
			return nil, errors.Wrap(err, "unable to marshal 's' value")
		}
		_, err = b.Write(data)
		if err != nil {
			return nil, errors.Wrap(err, "unable to write 's' value")
		}
	}

	return b.Bytes(), nil
}

// UnMarshal code
func (p *Payload) unmarshal(key string, val Element) error {

	switch key {
	case "t":
		if val.Code != Atom {
			return errors.Wrap(ErrBadPayload, "'t' was not an Atom")
		}

		if !val.IsNil() {
			eName, err := val.ToString()
			if err != nil {
				return errors.Wrap(err, "bad payload")
			}

			p.EventName = eName
		}

	case "s":
		if !val.Code.IsNumeric() && !val.IsNil() {
			return errors.Wrap(ErrBadPayload, "'s' was not numeric")
		}

		if !val.IsNil() {
			eVal, err := val.ToInt()
			if err != nil {
				return errors.Wrap(err, "bad payload")
			}

			p.SeqNum = &eVal
		}

	case "op":
		if val.Code != Int8 {
			return errors.Wrap(ErrBadPayload, "'op' was not an Int8")
		}

		eVal, err := val.ToInt()
		if err != nil {
			return errors.Wrap(err, "bad payload")
		}
		p.OpCode = constants.OpCode(eVal)

	case "d":
		switch val.Code {
		case Map:
			var err error
			p.Data, err = val.ToMap()
			if err != nil {
				return errors.Wrap(err, "bad payload")
			}
		case Atom:
			if !val.IsNil() {
				return errors.Wrap(ErrBadPayload, "'d' was not a map or list")
			}
		case EmptyList:
			fallthrough
		case List:
			p.DataList = val.Vals
		default:
			return errors.Wrap(ErrBadPayload, "'d' was not map or list")
		}

	default:
		return errors.Wrap(ErrBadPayload, fmt.Sprintf("unknown key '%s'", key))
	}
	return nil
}

// Unmarshal TODOC
func Unmarshal(raw []byte) (*Payload, error) {
	if len(raw) < 2 {
		return nil, ErrBadPayload
	}
	v := int(raw[0])
	if v != 131 {
		return nil, ErrBadPayload
	}

	p := Payload{}

	_, eSlice, err := unmarshalSlice(raw[1:], 1)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal bytes")
	}

	if eSlice[0].Code != 116 { // not a map
		return nil, errors.Wrap(ErrBadPayload, "payload not a map")
	}

	if len(eSlice[0].Vals)%2 != 0 {
		return nil, errors.Wrap(ErrBadPayload, "map parity incorrect incomplete")
	}

	for i := 0; i < len(eSlice[0].Vals); i += 2 {
		err = p.unmarshal(string(eSlice[0].Vals[i].Val), eSlice[0].Vals[i+1])
		if err != nil {
			return nil, errors.Wrap(err, "could not unmarshal field")
		}
	}

	return &p, nil
}

// PrettyString TODOC
func (p *Payload) PrettyString(indent string, skipFirstIndent bool) string {
	b := bytes.Buffer{}
	if skipFirstIndent {
		_, _ = b.WriteString("Payload{\n")
	} else {
		_, _ = b.WriteString(fmt.Sprintf("%sPayload{\n", indent))
	}

	_, _ = b.WriteString(fmt.Sprintf("%s  OpCode: %v\n", indent, p.OpCode))
	if p.EventName != "" {
		_, _ = b.WriteString(fmt.Sprintf("%s  EventName: %v\n", indent, p.EventName))
	}
	if p.SeqNum != nil {
		_, _ = b.WriteString(fmt.Sprintf("%s  SeqNum: %v\n", indent, *p.SeqNum))
	}

	if p.DataList != nil {
		_, _ = b.WriteString(fmt.Sprintf("%s  DataList: [\n", indent))
		for _, v := range p.DataList {
			_, _ = b.WriteString(v.PrettyString(indent+"     ", false))
			_, _ = b.WriteString("\n")
		}
		_, _ = b.WriteString(fmt.Sprintf("%s  ]", indent))
	} else {
		_, _ = b.WriteString(fmt.Sprintf("%s  Data: {\n", indent))
		for k, v := range p.Data {
			_, _ = b.WriteString(fmt.Sprintf("%s    %v: ", indent, k))
			_, _ = b.WriteString(v.PrettyString(indent+"     ", true))
			_, _ = b.WriteString("\n")
		}
		_, _ = b.WriteString(fmt.Sprintf("%s  }", indent))
	}

	return b.String()
}
