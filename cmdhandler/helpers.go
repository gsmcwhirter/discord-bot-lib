package cmdhandler

import (
	"fmt"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

// UserMentionString TODOC
func UserMentionString(uid snowflake.Snowflake) string {
	return fmt.Sprintf("<@!%s>", uid.ToString())
}

// ChannelMentionString TODOC
func ChannelMentionString(cid snowflake.Snowflake) string {
	return fmt.Sprintf("<#%s>", cid.ToString())
}

// RoleMentionString TODOC
func RoleMentionString(rid snowflake.Snowflake) string {
	return fmt.Sprintf("<@&%s>", rid.ToString())
}
