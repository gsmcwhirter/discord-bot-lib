package entity

import (
	"github.com/gsmcwhirter/go-util/v10/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
)

// Attachment is the data about an attachment recevied from the json api
type Attachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`

	IDSnowflake snowflake.Snowflake
}

// Snowflakify converts snowflake strings into real sowflakes
func (ar *Attachment) Snowflakify() error {
	var err error
	if ar.IDSnowflake, err = snowflake.FromString(ar.ID); err != nil {
		return errors.Wrap(err, "could not snowflakify ID")
	}

	return nil
}
