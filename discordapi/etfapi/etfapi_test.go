package etfapi

import (
	"reflect"
	"testing"
)

func TestPayload_Marshal(t *testing.T) {
	t.Parallel()

	// var one = 1
	tests := []struct {
		name    string
		p       Payload
		want    []byte
		wantErr bool
	}{
		// {
		// 	name: "base case",
		// 	p: Payload{
		// 		OpCode: 10,
		// 		Data: map[string]Element{
		// 			"_trace": {
		// 				Code: 108,
		// 				Val:  nil,
		// 				Vals: []Element{
		// 					{
		// 						Code: 109,
		// 						Val:  []byte{103, 97, 116, 101, 119, 97, 121, 45, 112, 114, 100, 45, 109, 97, 105, 110, 45, 118, 109, 116, 107},
		// 						Vals: nil,
		// 					},
		// 				},
		// 			},
		// 			"heartbeat_interval": {
		// 				Code: 98,
		// 				Val:  []byte{0, 0, 161, 34},
		// 				Vals: nil,
		// 			},
		// 		},
		// 		SeqNum:    &one,
		// 		EventName: "",
		// 	},
		// 	want:    []byte{131, 116, 0, 0, 0, 3, 109, 0, 0, 0, 2, 111, 112, 97, 10, 109, 0, 0, 0, 1, 100, 116, 0, 0, 0, 2, 109, 0, 0, 0, 6, 95, 116, 114, 97, 99, 101, 108, 0, 0, 0, 1, 109, 0, 0, 0, 21, 103, 97, 116, 101, 119, 97, 121, 45, 112, 114, 100, 45, 109, 97, 105, 110, 45, 118, 109, 116, 107, 106, 109, 0, 0, 0, 18, 104, 101, 97, 114, 116, 98, 101, 97, 116, 95, 105, 110, 116, 101, 114, 118, 97, 108, 98, 0, 0, 161, 34, 109, 0, 0, 0, 1, 115, 98, 0, 0, 0, 1},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.p.Marshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Payload.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	type args struct {
		raw []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Payload
		wantErr bool
	}{
		{
			name: "base case",
			args: args{[]byte{
				131,             // start code
				116, 0, 0, 0, 4, // map 1 length 4

				100, 0, 1, 100, // map 1[0] key
				116, 0, 0, 0, 2, // map1[0] val (map 2 length 2)

				100, 0, 6, 95, 116, 114, 97, 99, 101, // map 2[0] key
				108, 0, 0, 0, 1, // map 2[0] val (list length 1)
				109, 0, 0, 0, 21, 103, 97, 116, 101, 119, 97, 121, 45, 112, 114, 100, 45, 109, 97, 105, 110, 45, 118, 109, 116, 107, 106, // list entry binary

				100, 0, 18, 104, 101, 97, 114, 116, 98, 101, 97, 116, 95, 105, 110, 116, 101, 114, 118, 97, 108, // map 2[1] key
				98, 0, 0, 161, 34, // map 2[1] val

				100, 0, 2, 111, 112, // map 1[1] key
				97, 10, // map 1[1] val

				100, 0, 1, 115, // map 1[2] key
				100, 0, 3, 110, 105, 108, // map 1[2] val

				100, 0, 1, 116, // map 1[3] key
				100, 0, 3, 110, 105, 108, // map 1[3] val
			}},
			want: &Payload{
				OpCode: 10,
				Data: map[string]Element{
					"_trace": {
						Code: 108,
						Val:  nil,
						Vals: []Element{
							{
								Code: 109,
								Val:  []byte{103, 97, 116, 101, 119, 97, 121, 45, 112, 114, 100, 45, 109, 97, 105, 110, 45, 118, 109, 116, 107},
								Vals: nil,
							},
						},
					},
					"heartbeat_interval": {
						Code: 98,
						Val:  []byte{0, 0, 161, 34},
						Vals: nil,
					},
				},
				SeqNum: nil,
				EName:  "",
			},
			wantErr: false,
		},
		// {
		// 	name: "identify",
		// 	args: args{[]byte{
		// 		131,

		// 		116, 0, 0, 0, 2,

		// 		100, 0, 2, 111, 112,
		// 		97, 2,

		// 		100, 0, 1, 100,
		// 		116, 0, 0, 0, 5,

		// 		100, 0, 5, 115, 104, 97, 114, 100,
		// 		108, 0, 0, 0, 2,
		// 		98, 0, 0, 0, 0,
		// 		98, 0, 0, 0, 1,
		// 		106,

		// 		100, 0, 8, 112, 114, 101, 115, 101, 110, 99, 101,
		// 		116, 0, 0, 0, 4,

		// 		100, 0, 6, 115, 116, 97, 116, 117, 115,
		// 		100, 0, 6, 111, 110, 108, 105, 110, 101,

		// 		100, 0, 5, 115, 105, 110, 99, 101,
		// 		98, 0, 0, 0, 0,

		// 		100, 0, 4, 103, 97, 109, 101,
		// 		116, 0, 0, 0, 2,

		// 		100, 0, 4, 110, 97, 109, 101,
		// 		100, 0, 17, 76, 105, 115, 116, 32, 77, 97, 110, 97, 103, 101, 114, 32, 50, 48, 49, 56,

		// 		100, 0, 4, 116, 121, 112, 101,
		// 		98, 0, 0, 0, 0,

		// 		100, 0, 3, 97, 102, 107,
		// 		100, 0, 5, 102, 97, 108, 115, 101,

		// 		100, 0, 5, 116, 111, 107, 101, 110,
		// 		100, 0, 59, 78, 68, 77, 51, 78, 122, 81, 50, 78, 122, 107, 50, 77, 106, 65, 51, 78, 84, 77, 52, 77, 84, 99, 50, 46, 68, 99, 85, 67, 110, 65, 46, 69, 86, 51, 70, 104, 106, 119, 53, 87, 51, 100, 99, 107, 50, 86, 83, 107, 89, 118, 112, 70, 89, 50, 49, 103, 87, 103,

		// 		100, 0, 10, 112, 114, 111, 112, 101, 114, 116, 105, 101, 115,
		// 		116, 0, 0, 0, 3,

		// 		100, 0, 3, 36, 111, 115,
		// 		100, 0, 6, 103, 111, 108, 97, 110, 103,

		// 		100, 0, 8, 36, 98, 114, 111, 119, 115, 101, 114,
		// 		100, 0, 15, 101, 115, 111, 45, 100, 105, 115, 99, 111, 114, 100, 32, 98, 111, 116,

		// 		100, 0, 7, 36, 100, 101, 118, 105, 99, 101,
		// 		100, 0, 15, 101, 115, 111, 45, 100, 105, 115, 99, 111, 114, 100, 32, 98, 111, 116,

		// 		100, 0, 15, 108, 97, 114, 103, 101, 95, 116, 104, 114, 101, 115, 104, 111, 108, 100,
		// 		98, 0, 0, 0, 250,
		// 	}},
		// 	want:    nil,
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Unmarshal(tt.args.raw)
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
