package snowflake

import (
	"strconv"
)

// Snowflake represents a discord-like snowflake id
type Snowflake uint64

// ToString converts the snowflake into a string representation
func (s Snowflake) ToString() string {
	return strconv.FormatInt(int64(s), 10)
}

// FromString converts a string representation of a snowflake into a snowflake
func FromString(v string) (s Snowflake, err error) {
	i, err := strconv.ParseUint(v, 10, 64)
	s = Snowflake(i)
	return
}
