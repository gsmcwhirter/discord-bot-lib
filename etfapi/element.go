package etfapi

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Element is a container for arbitrary etf-formatted data
type Element struct {
	Code ETFCode
	Val  []byte
	Vals []Element
}

var trueB = []byte("true")
var falseB = []byte("false")
var nilB = []byte("nil")

// NewCollectionElement generates a new Element to hold data for collection types.
func NewCollectionElement(code ETFCode, val []Element) (e Element, err error) {
	if !code.IsCollection() {
		err = ErrBadElementData
		return
	}

	e.Code = code
	e.Vals = val

	return
}

// NewBasicElement generates a new Element to hold data for non-collection types.
func NewBasicElement(code ETFCode, val []byte) (e Element, err error) {
	if code.IsCollection() {
		err = ErrBadElementData
		return
	}

	e.Code = code
	e.Val = val

	return

}

// NewNilElement generates a new Element representing "nil"
func NewNilElement() (e Element, err error) {
	e, err = NewAtomElement(nilB)
	err = errors.Wrap(err, "could not create Nil Element")
	return
}

// NewBoolElement generates a new Element representing a boolean value
func NewBoolElement(val bool) (e Element, err error) {
	if val {
		e, err = NewAtomElement(trueB)
	} else {
		e, err = NewAtomElement(falseB)
	}

	err = errors.Wrap(err, "could not create Bool Element")
	return
}

// NewInt8Element generates a new Element representing an 8-bit unsigned integer value;
// Bounds checking will happen inside the function.
func NewInt8Element(val int) (e Element, err error) {
	var v []byte
	v, err = intToInt8Slice(val)
	if err != nil {
		err = errors.Wrap(err, "could not convert to int8 slice")
		return
	}

	e, err = NewBasicElement(Int8, v)
	err = errors.Wrap(err, "could not create Int8 Element")
	return
}

// NewInt32Element generates a new Element representing a 32-bit unsigned integer value;
// Bounds checking will happen inside the function
func NewInt32Element(val int) (e Element, err error) {
	var v []byte
	v, err = intToInt32Slice(val)
	if err != nil {
		err = errors.Wrap(err, "could not convert to int32 slice")
		return
	}

	e, err = NewBasicElement(Int32, v)
	err = errors.Wrap(err, "could not create Int32 Element")
	return
}

// NewSmallBigElement generates a new Element representing a 64-bit unsigned integer value
func NewSmallBigElement(val int64) (e Element, err error) {
	var v []byte
	v, err = int64ToInt64Slice(val)
	if err != nil {
		err = errors.Wrap(err, "could not convert to int64 slice")
		return
	}

	e, err = NewBasicElement(SmallBig, v)
	err = errors.Wrap(err, "could not create SmallBig Element")
	return
}

// NewBinaryElement generates a new Element representing Binary data
func NewBinaryElement(val []byte) (e Element, err error) {
	e, err = NewBasicElement(Binary, val)
	err = errors.Wrap(err, "could not create binary Element")
	return
}

// NewAtomElement generates a new Element representing an Atom value
func NewAtomElement(val []byte) (e Element, err error) {
	e, err = NewBasicElement(Atom, val)
	err = errors.Wrap(err, "could not create atom Element")
	return
}

// NewStringElement generates a new Element representing a String value
func NewStringElement(val string) (e Element, err error) {
	e, err = NewBasicElement(Binary, []byte(val))
	err = errors.Wrap(err, "could not create string Element")
	return
}

// NewMapElement generates a new Element representing a Map
//
// Keys are encoded as Binary type elements
func NewMapElement(val map[string]Element) (e Element, err error) {
	e2, err := ElementMapToElementSlice(val)
	if err != nil {
		err = errors.Wrap(err, "could not create element slice")
		return
	}

	e, err = NewCollectionElement(Map, e2)
	err = errors.Wrap(err, "could not create map Element")
	return
}

// NewListElement generates a new Element representing a List
//
// NOTE: empty lists are likely not handled well
func NewListElement(val []Element) (e Element, err error) {
	e, err = NewCollectionElement(List, val)
	err = errors.Wrap(err, "could not create list Element")
	return
}

func (e Element) String() string {
	switch e.Code {
	case Map, List:
		return fmt.Sprintf("Element{Code: %v, Vals: %v}", e.Code, e.Vals)
	default:
		return fmt.Sprintf("Element{Code: %v, Val: %v}", e.Code, e.Val)
	}
}

// Marshal formats the data in the given element in etf binary format
func (e *Element) Marshal() ([]byte, error) {
	b := &bytes.Buffer{}
	err := e.MarshalTo(b)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal Element")
	}
	return b.Bytes(), nil
}

// MarshalTo formats the data in the given element in etf binary format
// and writes it to the provided writer
func (e *Element) MarshalTo(b io.Writer) error {
	var err error

	_, err = b.Write([]byte{byte(e.Code)})
	if err != nil {
		return errors.Wrap(err, "could not marshal element code")
	}

	switch e.Code {
	case Map:
		err = marshalMapTo(b, e.Vals)
	case EmptyList:
		err = nil
	case List:
		err = marshalListTo(b, e.Vals)
	case Atom, String:
		err = marshalStringTo(b, e.Val)
	case Binary:
		err = marshalBinaryTo(b, e.Val)
	case Int8:
		err = marshalInt8To(b, e.Val)
	case Int32:
		err = marshalInt32To(b, e.Val)
	case SmallBig:
		err = marshalInt64To(b, e.Val)
	default:
		err = errors.Wrap(ErrBadMarshalData, "unsupported etf element code")
	}

	return errors.Wrap(err, "could not marshal element data")
}

// ToString converts a string-like element to a real string, if possible
func (e *Element) ToString() (string, error) {
	switch e.Code {
	case Atom, String, Binary:
		return string(e.Val), nil
	default:
		return "", errors.Wrap(ErrBadTarget, "cannot convert to string")
	}
}

// ToBytes converts a string-like Element to a []byte, if possible
func (e *Element) ToBytes() ([]byte, error) {
	switch e.Code {
	case Atom, String, Binary:
		b := make([]byte, len(e.Val))
		copy(b, e.Val)
		return b, nil
	default:
		return nil, errors.Wrap(ErrBadTarget, "cannot convert to []byte")
	}
}

// ToInt converts an int-like Element to an int, if possible
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

// ToInt64 converts an int-like Element to an int64, if possible
func (e *Element) ToInt64() (int64, error) {
	switch e.Code {
	case Int8:
		v, err := int8SliceToInt(e.Val)
		return int64(v), err
	case Int32:
		v, err := int32SliceToInt(e.Val)
		return int64(v), err
	case SmallBig, LargeBig:
		v, err := intNSliceToInt64(e.Val[1:])
		if e.Val[0] == 1 {
			v *= -1
		}

		return v, err
	default:
		return 0, errors.Wrap(ErrBadTarget, "cannot convert to int64")
	}
}

// ToMap converts a map Element to a map
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

// ToList converts a list Element to a list
func (e *Element) ToList() ([]Element, error) {
	if !e.IsList() {
		return nil, errors.Wrap(ErrBadTarget, "cannot convert to list")
	}

	return e.Vals, nil
}

// IsNumeric determines if an element is number-like
func (e *Element) IsNumeric() bool {
	return e.Code.IsNumeric()
}

// IsCollection determines if an element is a collection (map or list)
func (e *Element) IsCollection() bool {
	return e.Code.IsCollection()
}

// IsStringish determines if an element is string-like
func (e *Element) IsStringish() bool {
	return e.Code.IsStringish()
}

// IsList determines if an element is a list (with or without members)
func (e *Element) IsList() bool {
	return e.Code.IsList()
}

// IsNil determines if an element represents a "nil" value
func (e *Element) IsNil() bool {
	return e.Code == Atom && bytes.Equal(e.Val, nilB)
}

// IsTrue determines if an element represents a "true" value
func (e *Element) IsTrue() bool {
	return e.Code == Atom && bytes.Equal(e.Val, trueB)
}

// IsFalse determines if an element represents a "false" value
func (e *Element) IsFalse() bool {
	return e.Code == Atom && bytes.Equal(e.Val, falseB)
}

// // PrettyString generates a pretty, human-readable representation of an Element
// func (e *Element) PrettyString(indent string, skipFirstIndent bool) string {
// 	b := bytes.Buffer{}

// 	if e.Code.IsStringish() {
// 		if skipFirstIndent {
// 			indent = ""
// 		}
// 		_, _ = b.WriteString(fmt.Sprintf("%s%s", indent, string(e.Val)))
// 		return b.String()
// 	}

// 	if skipFirstIndent {
// 		_, _ = b.WriteString("Element{\n")
// 	} else {
// 		_, _ = b.WriteString(fmt.Sprintf("%sElement{\n", indent))
// 	}

// 	_, _ = b.WriteString(fmt.Sprintf("%s  Type: %v\n", indent, e.Code))
// 	if e.Code == List {
// 		_, _ = b.WriteString(fmt.Sprintf("%s  Vals: [\n", indent))
// 		for _, v := range e.Vals {
// 			_, _ = b.WriteString(v.PrettyString(indent+"     ", false))
// 			_, _ = b.WriteString("\n")
// 		}
// 		_, _ = b.WriteString(fmt.Sprintf("%s  ]", indent))
// 	} else if e.Code == Map {
// 		_, _ = b.WriteString(fmt.Sprintf("%s  Vals: {\n", indent))
// 		for i := 0; i < len(e.Vals); i += 2 {
// 			_, _ = b.WriteString(e.Vals[i].PrettyString(indent+"     ", false))
// 			_, _ = b.WriteString(": ")
// 			_, _ = b.WriteString(e.Vals[i+1].PrettyString(indent+"     ", true))
// 			_, _ = b.WriteString("\n")
// 		}
// 		_, _ = b.WriteString(fmt.Sprintf("%s  }", indent))
// 	} else {
// 		_, _ = b.WriteString(fmt.Sprintf("%s  Val: %v", indent, e.Val))
// 	}

// 	return b.String()
// }
