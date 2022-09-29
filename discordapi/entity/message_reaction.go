package entity

import "github.com/gsmcwhirter/go-util/v10/errors"

// MessageReaction is the data about a reaction received from the json api
type MessageReaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}

// Snowflakify converts snowflake strings into real sowflakes
func (rr *MessageReaction) Snowflakify() error {
	if err := rr.Emoji.Snowflakify(); err != nil {
		return errors.Wrap(err, "could not snowflakify Emoji")
	}

	return nil
}
