package payloads

import (
	"reflect"
	"testing"

	"github.com/gsmcwhirter/discord-bot-lib/etfapi"
)

func TestIdentifyPayload_Payload(t *testing.T) {
	t.Skip()

	type fields struct {
		Token          string
		Properties     IdentifyPayloadProperties
		LargeThreshold int
		Shard          IdentifyPayloadShard
		Presence       IdentifyPayloadPresence
	}
	tests := []struct {
		name    string
		fields  fields
		wantP   etfapi.Payload
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				Token: "token",
				Properties: IdentifyPayloadProperties{
					OS:      "golang",
					Browser: "eso-discord bot",
					Device:  "eso-discord bot",
				},
				LargeThreshold: 250,
				Shard: IdentifyPayloadShard{
					ID:    0,
					MaxID: 0,
				},
				Presence: IdentifyPayloadPresence{
					Game: IdentifyPayloadGame{
						Name: "List Manager 2018",
						Type: 0,
					},
					Status: "online",
					Since:  0,
					AFK:    false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := IdentifyPayload{
				Token:          tt.fields.Token,
				Properties:     tt.fields.Properties,
				LargeThreshold: tt.fields.LargeThreshold,
				Shard:          tt.fields.Shard,
				Presence:       tt.fields.Presence,
			}
			gotP, err := ip.Payload()
			if (err != nil) != tt.wantErr {
				t.Errorf("IdentifyPayload.Payload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotP, tt.wantP) {
				t.Errorf("IdentifyPayload.Payload() = %+v, want %v", gotP, tt.wantP)
			}
		})
	}
}
