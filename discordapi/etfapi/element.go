package etfapi

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Element TODOC
type Element struct {
	Code ETFCode
	Val  []byte
	Vals []Element
}

// NewElement TODOC
func NewElement(code ETFCode, val interface{}) (e Element, err error) {
	e.Code = code

	if v, ok := val.([]Element); ok {
		if !code.IsCollection() {
			err = ErrBadElementData
			return
		}
		e.Vals = v

		return
	}

	if v, ok := val.([]byte); ok {
		if code.IsCollection() {
			err = ErrBadElementData
			return
		}
		e.Val = v

		return
	}

	return
}

// NewNilElement TODOC
func NewNilElement() (e Element, err error) {
	e, err = NewAtomElement("nil")
	err = errors.Wrap(err, "could not create Nil Element")
	return
}

// NewBoolElement TODOC
func NewBoolElement(val bool) (e Element, err error) {
	if val {
		e, err = NewAtomElement("true")
	} else {
		e, err = NewAtomElement("false")
	}

	err = errors.Wrap(err, "could not create Bool Element")
	return
}

// NewInt8Element TODOC
func NewInt8Element(val int) (e Element, err error) {
	var v []byte
	v, err = intToInt8Slice(val)
	if err != nil {
		err = errors.Wrap(err, "could not convert to int8 slice")
		return
	}

	e, err = NewElement(Int8, v)
	err = errors.Wrap(err, "could not create Int8 Element")
	return
}

// NewInt32Element TODOC
func NewInt32Element(val int) (e Element, err error) {
	var v []byte
	v, err = intToInt32Slice(val)
	if err != nil {
		err = errors.Wrap(err, "could not convert to int32 slice")
		return
	}

	e, err = NewElement(Int32, v)
	err = errors.Wrap(err, "could not create Int32 Element")
	return
}

// NewBinaryElement TODOC
func NewBinaryElement(val []byte) (e Element, err error) {
	e, err = NewElement(Binary, val)
	err = errors.Wrap(err, "could not create binary Element")
	return
}

// NewAtomElement TODOC
func NewAtomElement(val string) (e Element, err error) {
	e, err = NewElement(Atom, []byte(val))
	err = errors.Wrap(err, "could not create atom Element")
	return
}

// NewStringElement TODOC
func NewStringElement(val string) (e Element, err error) {
	e, err = NewElement(Binary, []byte(val))
	err = errors.Wrap(err, "could not create string Element")
	return
}

// NewMapElement TODOC
func NewMapElement(val map[string]Element) (e Element, err error) {
	e2, err := ElementMapToElementSlice(val)
	if err != nil {
		err = errors.Wrap(err, "could not create element slice")
		return
	}

	e, err = NewElement(Map, e2)
	err = errors.Wrap(err, "could not create map Element")
	return
}

// NewListElement TODOC
func NewListElement(val []Element) (e Element, err error) {
	e, err = NewElement(List, val)
	err = errors.Wrap(err, "could not create list Element")
	return
}

// String TODOC
func (e Element) String() string {
	switch e.Code {
	case Map:
		fallthrough
	case List:
		return fmt.Sprintf("Element{Code: %v, Vals: %v}", e.Code, e.Vals)
	default:
		return fmt.Sprintf("Element{Code: %v, Val: %v}", e.Code, e.Val)
	}
}

// WriteTo TODOC
func (e *Element) WriteTo(b io.Writer) (int64, error) {
	var tmp interface{}
	if e.Val != nil {
		tmp = e.Val
	} else if e.Vals != nil {
		tmp = e.Vals
	} else {
		tmp = nil
	}

	data, err := marshalInterface(e.Code, tmp)
	if err != nil {
		return 0, errors.Wrap(err, "couldn't marshal element")
	}

	n, err := b.Write(data)
	return int64(n), err
}

// ToString TODOC
func (e *Element) ToString() (string, error) {
	switch e.Code {
	case Atom:
		fallthrough
	case String:
		fallthrough
	case Binary:
		return string(e.Val), nil
	default:
		return "", errors.Wrap(ErrBadTarget, "cannot convert to string")
	}
}

// ToBytes TODOC
func (e *Element) ToBytes() ([]byte, error) {
	switch e.Code {
	case Atom:
		fallthrough
	case String:
		fallthrough
	case Binary:
		b := make([]byte, len(e.Val))
		copy(b, e.Val)
		return b, nil
	default:
		return nil, errors.Wrap(ErrBadTarget, "cannot convert to []byte")
	}
}

// ToInt TODOC
func (e *Element) ToInt() (int, error) {
	switch e.Code {
	case Int8:
		return int8SliceToInt(e.Val)
	case Int32:
		return int32SliceToInt(e.Val)
	default:
		return 0, errors.Wrap(ErrBadTarget, "cannot convert to int")
	}
}

// ToInt64 TODOC
func (e *Element) ToInt64() (int64, error) {
	switch e.Code {
	case Int8:
		v, err := int8SliceToInt(e.Val)
		return int64(v), err
	case Int32:
		v, err := int32SliceToInt(e.Val)
		return int64(v), err
	case SmallBig:
		fallthrough
	case LargeBig:
		v, err := intNSliceToInt64(e.Val[1:])
		if e.Val[0] == 1 {
			v *= -1
		}

		return v, err
	default:
		return 0, errors.Wrap(ErrBadTarget, "cannot convert to int64")
	}
}

// ToMap TODOC
func (e *Element) ToMap() (map[string]Element, error) {
	switch e.Code {
	case Map:
		v := map[string]Element{}

		if len(e.Vals)%2 != 0 {
			return nil, ErrBadParity
		}

		for i := 0; i < len(e.Vals); i += 2 {
			k, err := e.Vals[i].ToString()
			if err != nil {
				return nil, ErrBadFieldType
			}
			v[k] = e.Vals[i+1]
		}

		return v, nil
	default:
		return nil, errors.Wrap(ErrBadTarget, "cannot convert to map")
	}
}

// ToList TODOC
func (e *Element) ToList() ([]Element, error) {
	if !e.IsList() {
		return nil, errors.Wrap(ErrBadTarget, "cannot convert to list")
	}

	return e.Vals, nil
}

// IsNumeric TODOC
func (e *Element) IsNumeric() bool {
	return e.Code.IsNumeric()
}

// IsCollection TODOC
func (e *Element) IsCollection() bool {
	return e.Code.IsCollection()
}

// IsStringish TODOC
func (e *Element) IsStringish() bool {
	return e.Code.IsStringish()
}

// IsList TODOC
func (e *Element) IsList() bool {
	return e.Code.IsList()
}

// IsNil TODOC
func (e *Element) IsNil() bool {
	return e.Code == Atom && string(e.Val) == "nil"
}

// IsTrue TODOC
func (e *Element) IsTrue() bool {
	return e.Code == Atom && string(e.Val) == "true"
}

// IsFalse TODOC
func (e *Element) IsFalse() bool {
	return e.Code == Atom && string(e.Val) == "false"
}

// PrettyString TODOC
func (e *Element) PrettyString(indent string, skipFirstIndent bool) string {
	b := bytes.Buffer{}

	if e.Code.IsStringish() {
		if skipFirstIndent {
			indent = ""
		}
		_, _ = b.WriteString(fmt.Sprintf("%s%s", indent, string(e.Val)))
		return b.String()
	}

	if skipFirstIndent {
		_, _ = b.WriteString("Element{\n")
	} else {
		_, _ = b.WriteString(fmt.Sprintf("%sElement{\n", indent))
	}

	_, _ = b.WriteString(fmt.Sprintf("%s  Type: %v\n", indent, e.Code))
	if e.Code == List {
		_, _ = b.WriteString(fmt.Sprintf("%s  Vals: [\n", indent))
		for _, v := range e.Vals {
			_, _ = b.WriteString(v.PrettyString(indent+"     ", false))
			_, _ = b.WriteString("\n")
		}
		_, _ = b.WriteString(fmt.Sprintf("%s  ]", indent))
	} else if e.Code == Map {
		_, _ = b.WriteString(fmt.Sprintf("%s  Vals: {\n", indent))
		for i := 0; i < len(e.Vals); i += 2 {
			_, _ = b.WriteString(e.Vals[i].PrettyString(indent+"     ", false))
			_, _ = b.WriteString(": ")
			_, _ = b.WriteString(e.Vals[i+1].PrettyString(indent+"     ", true))
			_, _ = b.WriteString("\n")
		}
		_, _ = b.WriteString(fmt.Sprintf("%s  }", indent))
	} else {
		_, _ = b.WriteString(fmt.Sprintf("%s  Val: %v", indent, e.Val))
	}

	return b.String()
}
