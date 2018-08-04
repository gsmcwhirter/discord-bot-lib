package etfapi

import "fmt"

// ETFCode TODOC
type ETFCode int

// See https://github.com/discordapp/erlpack/blob/master/cpp/constants.h

// TODOC
const (
	Map       ETFCode = 116
	Atom              = 100
	List              = 108
	Binary            = 109
	Int8              = 97
	Int32             = 98
	Float             = 70
	String            = 107
	EmptyList         = 106
	SmallBig          = 110
	LargeBig          = 111
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

// IsNumeric TODOC
func (c ETFCode) IsNumeric() bool {
	return c == Int8 || c == Int32 || c == Float || c == SmallBig || c == LargeBig
}

// IsCollection TODOC
func (c ETFCode) IsCollection() bool {
	return c == Map || c == List || c == EmptyList
}

// IsStringish TODOC
func (c ETFCode) IsStringish() bool {
	return c == Atom || c == String || c == Binary
}

// IsList TODOC
func (c ETFCode) IsList() bool {
	return c == List || c == EmptyList
}
