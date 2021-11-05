package etfapi

import "fmt"

type Code int

// See https://github.com/discordapp/erlpack/blob/master/cpp/discordapi.h

// These are the ETF type codes that this library knows about
const (
	Map       Code = 116
	Atom      Code = 100
	List      Code = 108
	Binary    Code = 109
	Int8      Code = 97
	Int32     Code = 98
	Float     Code = 70
	String    Code = 107
	EmptyList Code = 106
	SmallBig  Code = 110
	LargeBig  Code = 111
)

func (c Code) String() string {
	switch c {
	case Map:
		return "Map"
	case Atom:
		return "Atom"
	case List:
		return "List"
	case Binary:
		return "Binary"
	case Int32:
		return "Int32"
	case Int8:
		return "Int8"
	case Float:
		return "Float"
	case String:
		return "String"
	case EmptyList:
		return "EmptyList"
	case SmallBig:
		return "SmallBig"
	case LargeBig:
		return "LargeBig"
	default:
		return fmt.Sprintf("(unknown: %d)", int(c))
	}
}

// IsNumeric determines if an ETF Code is number-like
func (c Code) IsNumeric() bool {
	return c == Int8 || c == Int32 || c == Float || c == SmallBig || c == LargeBig
}

// IsCollection determines if an ETF Code is a collection of other elements
func (c Code) IsCollection() bool {
	return c == Map || c == List || c == EmptyList
}

// IsStringish determines if an ETF Code is string-like
func (c Code) IsStringish() bool {
	return c == Atom || c == String || c == Binary
}

// IsList determines if an ETF Code is a list
func (c Code) IsList() bool {
	return c == List || c == EmptyList
}
