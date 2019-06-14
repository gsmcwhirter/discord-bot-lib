package cmdhandler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gsmcwhirter/discord-bot-lib/v8/snowflake"
)

const longContent = `"""
CR (#raid-signup)
aa (#raid-signup)
description=light your torches, sharpen your pitchforks, and kill those chickens! cloudrest will be held saturday 3/23 at 4pm! please, do not sign up if you do not plan on attending!

in effect as of now will be a weekly vote for your saturday raids! pick your favorite down below to cast your vote.
☁️ -cloudrest
<:emoji_1:549001147294416931>  -sanctum ophidia
🦄-aetherian archive
🤖-asylum sanctorium
🐲-hel ra citadel
🦂-maw of lorkhaj
<:kryscry:500530805244690444>-halls of fabrication (#)
description=light your torches, sharpen your pitchforks, and kill those chickens! cloudrest will be held saturday 3/23 at 4pm! please, do not sign up if you do not plan on attending!

in effect as of now will be a weekly vote for your saturday raids! pick your favorite down below to cast your vote.
☁️ -cloudrest
👹  -sanctum ophidia
🦄-aetherian archive
🤖-asylum sanctorium
🐲-hel ra citadel
🦂-maw of lorkhaj
☠️ -halls of fabrication (#)
description=light your torches, sharpen your pitchforks, and kill those chickens! cloudrest will be held saturday 3/23 at 4pm! please, do not sign up if you do not plan on attending!

in effect as of now will be a weekly vote for your saturday raids! pick your favorite down below to cast your vote.
☁️  -cloudrest
👹 -sanctum ophidia
🦄 -aetherian archive
🤖 -asylum sanctorium
🐲 -hel ra citadel
🦂 -maw of lorkhaj
☠️ -halls of fabrication (#)
description=light your torches, sharpen your pitchforks, and kill those chickens! cloudrest will be held saturday 3/23 at 4pm! please, do not sign up if you do not plan on attending!

in effect as of now will be a weekly vote for your saturday raids! pick your favorite down below to cast your vote.
☁️-cloudrest    <<:emoji_1:549001147294416931>549001147294416931> -sanctum ophidia   🦄 -aetherian archive
🤖-asylum sanctorium
🐲-hel ra citadel
🦂-maw of lorkhaj
<<:kryscry:500530805244690444>500530805244690444> -halls of fabrication (#)
"""`

func TestSimpleResponseSplit(t *testing.T) {
	r := SimpleResponse{
		To:        "test",
		Content:   longContent,
		ToChannel: 1,
	}

	parts := r.Split()

	assert.Equal(t, 3, len(parts))

	for _, p := range parts {
		assert.Equal(t, "test", p.(*SimpleResponse).To)
		assert.Equal(t, snowflake.Snowflake(1), p.(*SimpleResponse).ToChannel)
	}
}

func TestSimpleEmbedResponseSplit(t *testing.T) {
	r := SimpleEmbedResponse{
		To:          "test",
		Title:       "test title",
		Description: longContent,
		Color:       5,
		FooterText:  "test footer",
		ToChannel:   1,
	}

	parts := r.Split()

	assert.Equal(t, 3, len(parts))

	for i, p := range parts {
		assert.Equal(t, "test", p.(*SimpleEmbedResponse).To)
		if i == 0 {
			assert.Equal(t, "test title", p.(*SimpleEmbedResponse).Title)
		} else {
			assert.Equal(t, "", p.(*SimpleEmbedResponse).Title)
		}
		assert.Equal(t, 5, p.(*SimpleEmbedResponse).Color)
		assert.Equal(t, "test footer", p.(*SimpleEmbedResponse).FooterText)
		assert.Equal(t, snowflake.Snowflake(1), p.(*SimpleEmbedResponse).ToChannel)
	}
}

func TestEmbedResponseSplit(t *testing.T) {
	r := EmbedResponse{
		To:          "test",
		Title:       "test title",
		Description: "",
		Fields: []EmbedField{
			{
				Name: "test",
				Val:  longContent,
			},
		},
		Color:      5,
		FooterText: "test footer",
		ToChannel:  1,
	}

	parts := r.Split()

	assert.Equal(t, 3, len(parts))

	for i, p := range parts {
		assert.Equal(t, "test", p.(*EmbedResponse).To)
		if i == 0 {
			assert.Equal(t, "test title", p.(*EmbedResponse).Title)
		} else {
			assert.Equal(t, "", p.(*EmbedResponse).Title)
		}
		assert.Equal(t, 5, p.(*EmbedResponse).Color)
		assert.Equal(t, "test footer", p.(*EmbedResponse).FooterText)
		assert.Equal(t, snowflake.Snowflake(1), p.(*EmbedResponse).ToChannel)
		assert.Equal(t, 1, len(p.(*EmbedResponse).Fields))
		assert.True(t, len(p.(*EmbedResponse).Fields[0].Val) > 2)
	}
}
