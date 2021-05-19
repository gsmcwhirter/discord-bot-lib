package cmdhandler

import (
	"reflect"
	"testing"
)

func TestNewMessageHandler(t *testing.T) {
	type args struct {
		f MessageHandlerFunc
	}
	tests := []struct {
		name string
		args args
		want MessageHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMessageHandler(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessageHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_msgHandlerFunc_HandleMessage(t *testing.T) {
	type fields struct {
		handler func(Message) (Response, error)
	}
	type args struct {
		msg Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lh := &msgHandlerFunc{
				handler: tt.fields.handler,
			}
			got, err := lh.HandleMessage(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgHandlerFunc.HandleMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgHandlerFunc.HandleMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
