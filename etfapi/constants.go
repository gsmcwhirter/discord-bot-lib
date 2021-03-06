package etfapi

import "fmt"

// ETFCode is a type alias for representing ETF type codes
type ETFCode int

// See https://github.com/discordapp/erlpack/blob/master/cpp/discordapi.h

// These are the ETF type codes that this library knows about
const (
	Map       ETFCode = 116
	Atom      ETFCode = 100
	List      ETFCode = 108
	Binary    ETFCode = 109
	Int8      ETFCode = 97
	Int32     ETFCode = 98
	Float     ETFCode = 70
	String    ETFCode = 107
	EmptyList ETFCode = 106
	SmallBig  ETFCode = 110
	LargeBig  ETFCode = 111
)

func (c ETFCode) String() string {
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

// IsNumeric determines if an ETFCode is number-like
func (c ETFCode) IsNumeric() bool {
	return c == Int8 || c == Int32 || c == Float || c == SmallBig || c == LargeBig
}

// IsCollection determines if an ETFCode is a collection of other elements
func (c ETFCode) IsCollection() bool {
	return c == Map || c == List || c == EmptyList
}

// IsStringish determines if an ETFCode is string-like
func (c ETFCode) IsStringish() bool {
	return c == Atom || c == String || c == Binary
}

// IsList determines if an ETFCode is a list
func (c ETFCode) IsList() bool {
	return c == List || c == EmptyList
}
