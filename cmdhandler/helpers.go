package cmdhandler

import (
	"fmt"
	"regexp"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

var userMentionRe = regexp.MustCompile(`^<@[!]?[0-9]+>$`)
var channelMentionRe = regexp.MustCompile(`^<#[0-9]+>$`)
var roleMentionRe = regexp.MustCompile(`^<@&[0-9]+>|@everyone|@here$`)

// UserMentionString TODOC
func UserMentionString(uid snowflake.Snowflake) string {
	return fmt.Sprintf("<@!%s>", uid.ToString())
}

// IsUserMention TODOC
func IsUserMention(v string) bool {
	return userMentionRe.MatchString(v)
}

// ChannelMentionString TODOC
func ChannelMentionString(cid snowflake.Snowflake) string {
	return fmt.Sprintf("<#%s>", cid.ToString())
}

// IsChannelMention TODOC
func IsChannelMention(v string) bool {
	return channelMentionRe.MatchString(v)
}

// RoleMentionString TODOC
func RoleMentionString(rid snowflake.Snowflake) string {
	return fmt.Sprintf("<@&%s>", rid.ToString())
}

// IsRoleMention TODOC
func IsRoleMention(v string) bool {
	return roleMentionRe.MatchString(v)
}
