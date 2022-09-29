package snowflake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnowflake(t *testing.T) {
	t.Parallel()
	var s Snowflake

	assert.Equal(t, "0", s.ToString())
	s, err := FromString("12345")
	if assert.Nil(t, err) {
		assert.Equal(t, Snowflake(12345), s)
	}
}
