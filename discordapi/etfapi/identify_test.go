package etfapi

import (
	"reflect"
	"testing"
)

func TestIdentifyPayload_Payload(t *testing.T) {
	t.Skip()
	t.Parallel()

	type fields struct {
		Token          string
		Intents        int
		Properties     IdentifyPayloadProperties
		LargeThreshold int
		Shard          IdentifyPayloadShard
		Presence       IdentifyPayloadPresence
	}
	tests := []struct {
		name    string
		fields  fields
		wantP   Payload
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				Token:   "token",
				Intents: 0,
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ip := IdentifyPayload{
				Token:          tt.fields.Token,
				Intents:        tt.fields.Intents,
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
