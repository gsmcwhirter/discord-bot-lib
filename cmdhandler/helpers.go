package cmdhandler

import (
	"fmt"
	"regexp"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
	"github.com/pkg/errors"
)

var userMentionRe = regexp.MustCompile(`^<@[!]?([0-9]+)>$`)
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

// ForceUserNicknameMention TODOC
func ForceUserNicknameMention(v string) (string, error) {
	matches := userMentionRe.FindStringSubmatch(v)
	if len(matches) < 2 || matches[0] == "" || matches[1] == "" {
		return "", errors.New("not a user mention")
	}

	return fmt.Sprintf("<@!%s>", matches[1]), nil
}

// ForceUserAccountMention TODOC
func ForceUserAccountMention(v string) (string, error) {
	matches := userMentionRe.FindStringSubmatch(v)
	if len(matches) < 2 || matches[0] == "" || matches[1] == "" {
		return "", errors.New("not a user mention")
	}

	return fmt.Sprintf("<@%s>", matches[1]), nil
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
