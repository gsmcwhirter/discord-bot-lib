package entity

import (
	"reflect"
	"testing"
)

func TestApplicationCommandOptionChoice_MarshalJSON(t *testing.T) {
	type fields struct {
		Name        string
		Type        ApplicationCommandOptionType
		ValueString string
		ValueInt    int
		ValueNumber float64
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "basic string",
			fields: fields{
				Name:        "s",
				Type:        OptTypeString,
				ValueString: "sv",
			},
			want:    []byte(`{"name":"s","value":"sv"}`),
			wantErr: false,
		},
		{
			name: "basic int",
			fields: fields{
				Name:     "i",
				Type:     OptTypeInteger,
				ValueInt: 7,
			},
			want:    []byte(`{"name":"i","value":7}`),
			wantErr: false,
		},
		{
			name: "basic number",
			fields: fields{
				Name:        "n",
				Type:        OptTypeNumber,
				ValueNumber: 3.14,
			},
			want:    []byte(`{"name":"n","value":3.14}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ApplicationCommandOptionChoice{
				Name:        tt.fields.Name,
				Type:        tt.fields.Type,
				ValueString: tt.fields.ValueString,
				ValueInt:    tt.fields.ValueInt,
				ValueNumber: tt.fields.ValueNumber,
			}
			got, err := c.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplicationCommandOptionChoice.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplicationCommandOptionChoice.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
