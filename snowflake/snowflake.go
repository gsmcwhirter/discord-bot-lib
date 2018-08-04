package snowflake

import (
	"strconv"
)

// Snowflake TODOC
type Snowflake uint64

// ToString TODOC
func (s Snowflake) ToString() string {
	return strconv.FormatInt(int64(s), 10)
}

// FromString TODOC
func FromString(v string) (s Snowflake, err error) {
	i, err := strconv.ParseUint(v, 10, 64)
	s = Snowflake(i)
	return
}
