package cmdhandler

import (
	"fmt"
	"regexp"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
	"github.com/pkg/errors"
)

// ErrNotUserMention is the error returned when a user mention string is required but not provided
var ErrNotUserMention = errors.New("not a user mention")

var userMentionRe = regexp.MustCompile(`^<@[!]?([0-9]+)>$`)
var channelMentionRe = regexp.MustCompile(`^<#[0-9]+>$`)
var roleMentionRe = regexp.MustCompile(`^<@&[0-9]+>|@everyone|@here$`)

// UserMentionString generates a string that discord interprets as a mention of a user
// by their server nickname
func UserMentionString(uid snowflake.Snowflake) string {
	return fmt.Sprintf("<@!%s>", uid.ToString())
}

// IsUserMention determines if a string is a mention of a user (either by nickname or account name)
func IsUserMention(v string) bool {
	return userMentionRe.MatchString(v)
}

// ForceUserNicknameMention converts a user mention into a nickname mention
// (if it is not already a nickname mention)
func ForceUserNicknameMention(v string) (string, error) {
	matches := userMentionRe.FindStringSubmatch(v)
	if len(matches) < 2 || matches[0] == "" || matches[1] == "" {
		return "", ErrNotUserMention
	}

	return fmt.Sprintf("<@!%s>", matches[1]), nil
}

// ForceUserAccountMention converts a user mention into an account mention
// (if it is not already an account mention)
func ForceUserAccountMention(v string) (string, error) {
	matches := userMentionRe.FindStringSubmatch(v)
	if len(matches) < 2 || matches[0] == "" || matches[1] == "" {
		return "", ErrNotUserMention
	}

	return fmt.Sprintf("<@%s>", matches[1]), nil
}

// ChannelMentionString generates a string that discord interprets as a mention of a channel
func ChannelMentionString(cid snowflake.Snowflake) string {
	return fmt.Sprintf("<#%s>", cid.ToString())
}

// IsChannelMention determines if a string is a mention of a channel
func IsChannelMention(v string) bool {
	return channelMentionRe.MatchString(v)
}

// RoleMentionString generates a string that discord interprets as a mention of a server role
func RoleMentionString(rid snowflake.Snowflake) string {
	return fmt.Sprintf("<@&%s>", rid.ToString())
}

// IsRoleMention determines if a string is a mention of a server role
func IsRoleMention(v string) bool {
	return roleMentionRe.MatchString(v)
}
