package etfapi_test

import (
	"reflect"
	"testing"

	"github.com/gsmcwhirter/discord-bot-lib/v15/discordapi"
	"github.com/gsmcwhirter/discord-bot-lib/v15/etfapi"
)

func TestPayload_Marshal(t *testing.T) {
	s := new(int)
	*s = 3

	type fields struct {
		OpCode    discordapi.OpCode
		SeqNum    *int
		EventName string
		Data      map[string]etfapi.Element
		DataList  []etfapi.Element
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				OpCode:    1,
				SeqNum:    s,
				EventName: "test",
				Data: map[string]etfapi.Element{
					"test": {
						Code: etfapi.Int8,
						Val:  []byte{128},
					},
				},
			},
			want: []byte{131, 116, 0, 0, 0, 3, 109, 0, 0, 0, 2, 111, 112, 97, 1, 109, 0, 0, 0, 1, 100, 116, 0, 0, 0, 1, 109, 0, 0, 0, 4, 116, 101, 115, 116, 97, 128, 109, 0, 0, 0, 1, 115, 98, 0, 0, 0, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &etfapi.Payload{
				OpCode:    tt.fields.OpCode,
				SeqNum:    tt.fields.SeqNum,
				EventName: tt.fields.EventName,
				Data:      tt.fields.Data,
				DataList:  tt.fields.DataList,
			}
			got, err := p.Marshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Payload.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Payload.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	s := new(int)
	*s = 3

	type args struct {
		raw []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *etfapi.Payload
		wantErr bool
	}{
		{
			name: "ok",
			args: args{[]byte{131, 116, 0, 0, 0, 3, 109, 0, 0, 0, 2, 111, 112, 97, 1, 109, 0, 0, 0, 1, 100, 116, 0, 0, 0, 1, 109, 0, 0, 0, 4, 116, 101, 115, 116, 97, 128, 109, 0, 0, 0, 1, 115, 98, 0, 0, 0, 3}},
			want: &etfapi.Payload{
				OpCode: 1,
				SeqNum: s,
				Data: map[string]etfapi.Element{
					"test": {
						Code: etfapi.Int8,
						Val:  []byte{128},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := etfapi.Unmarshal(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
