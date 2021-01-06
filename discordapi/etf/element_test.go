package etf_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gsmcwhirter/discord-bot-lib/v19/discordapi/etf"
)

func TestNewCollectionElement(t *testing.T) {
	type args struct {
		code etf.Code
		val  []etf.Element
	}
	tests := []struct {
		name    string
		args    args
		wantE   etf.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				code: etf.List,
				val: []etf.Element{
					{
						Code: etf.EmptyList,
					},
				},
			},
			wantE: etf.Element{
				Code: etf.List,
				Vals: []etf.Element{
					{
						Code: etf.EmptyList,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "not ok",
			args: args{
				code: etf.Int8,
				val: []etf.Element{
					{
						Code: etf.EmptyList,
					},
				},
			},
			wantE:   etf.Element{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := etf.NewCollectionElement(tt.args.code, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCollectionElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("NewCollectionElement() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNewBasicElement(t *testing.T) {
	type args struct {
		code etf.Code
		val  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantE   etf.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				code: etf.Int8,
				val:  []byte{1},
			},
			wantE: etf.Element{
				Code: etf.Int8,
				Val:  []byte{1},
			},
			wantErr: false,
		},
		{
			name: "not ok",
			args: args{
				code: etf.List,
				val:  []byte{2},
			},
			wantE:   etf.Element{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := etf.NewBasicElement(tt.args.code, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBasicElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("NewBasicElement() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNewNilElement(t *testing.T) {
	tests := []struct {
		name    string
		wantE   etf.Element
		wantErr bool
	}{
		{
			name: "ok",
			wantE: etf.Element{
				Code: etf.Atom,
				Val:  []byte("nil"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := etf.NewNilElement()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNilElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("NewNilElement() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNewBoolElement(t *testing.T) {
	type args struct {
		val bool
	}
	tests := []struct {
		name    string
		args    args
		wantE   etf.Element
		wantErr bool
	}{
		{
			name: "ok true",
			args: args{true},
			wantE: etf.Element{
				Code: etf.Atom,
				Val:  []byte("true"),
			},
		},
		{
			name: "ok false",
			args: args{false},
			wantE: etf.Element{
				Code: etf.Atom,
				Val:  []byte("false"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := etf.NewBoolElement(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBoolElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("NewBoolElement() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNewInt8Element(t *testing.T) {
	type args struct {
		val int
	}
	tests := []struct {
		name    string
		args    args
		wantE   etf.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{1},
			wantE: etf.Element{
				Code: etf.Int8,
				Val:  []byte{1},
			},
		},
		{
			name:    "not ok",
			args:    args{1024},
			wantE:   etf.Element{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := etf.NewInt8Element(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInt8Element() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("NewInt8Element() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNewInt32Element(t *testing.T) {
	type args struct {
		val int
	}
	tests := []struct {
		name    string
		args    args
		wantE   etf.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{1 << 31},
			wantE: etf.Element{
				Code: etf.Int32,
				Val:  []byte{128, 0, 0, 0},
			},
		},
		{
			name:    "not ok",
			args:    args{1 << 42},
			wantE:   etf.Element{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := etf.NewInt32Element(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInt32Element() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("NewInt32Element() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNewSmallBigElement(t *testing.T) {
	type args struct {
		val int64
	}
	tests := []struct {
		name    string
		args    args
		wantE   etf.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{1 << 42},
			wantE: etf.Element{
				Code: etf.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := etf.NewSmallBigElement(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSmallBigElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("NewSmallBigElement() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNewBinaryElement(t *testing.T) {
	type args struct {
		val []byte
	}
	tests := []struct {
		name    string
		args    args
		wantE   etf.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{[]byte("test")},
			wantE: etf.Element{
				Code: etf.Binary,
				Val:  []byte{116, 101, 115, 116},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := etf.NewBinaryElement(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBinaryElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("NewBinaryElement() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNewAtomElement(t *testing.T) {
	type args struct {
		val []byte
	}
	tests := []struct {
		name    string
		args    args
		wantE   etf.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{[]byte("test")},
			wantE: etf.Element{
				Code: etf.Atom,
				Val:  []byte{116, 101, 115, 116},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := etf.NewAtomElement(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAtomElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("NewAtomElement() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNewStringElement(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name    string
		args    args
		wantE   etf.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{"test"},
			wantE: etf.Element{
				Code: etf.Binary,
				Val:  []byte{116, 101, 115, 116},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := etf.NewStringElement(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStringElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("NewStringElement() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNewMapElement(t *testing.T) {
	type args struct {
		val map[string]etf.Element
	}
	tests := []struct {
		name    string
		args    args
		wantE   etf.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{map[string]etf.Element{
				"test": {Code: etf.Int8, Val: []byte{1}},
			}},
			wantE: etf.Element{
				Code: etf.Map,
				Vals: []etf.Element{
					{
						Code: etf.Binary,
						Val:  []byte{116, 101, 115, 116},
					},
					{
						Code: etf.Int8,
						Val:  []byte{1},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := etf.NewMapElement(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMapElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("NewMapElement() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestNewListElement(t *testing.T) {
	type args struct {
		val []etf.Element
	}
	tests := []struct {
		name    string
		args    args
		wantE   etf.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{[]etf.Element{
				{
					Code: etf.Int8,
					Val:  []byte{2},
				},
				{
					Code: etf.Int8,
					Val:  []byte{1},
				},
			}},
			wantE: etf.Element{
				Code: etf.List,
				Vals: []etf.Element{
					{
						Code: etf.Int8,
						Val:  []byte{2},
					},
					{
						Code: etf.Int8,
						Val:  []byte{1},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := etf.NewListElement(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewListElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("NewListElement() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestElement_Marshal(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "map",
			fields: fields{
				Code: etf.Map,
				Vals: []etf.Element{
					{
						Code: etf.Binary,
						Val:  []byte{116, 101, 115, 116},
					},
					{
						Code: etf.Int8,
						Val:  []byte{1},
					},
				},
			},
			want:    []byte{116, 0, 0, 0, 1, 109, 0, 0, 0, 4, 116, 101, 115, 116, 97, 1},
			wantErr: false,
		},
		{
			name: "empty list",
			fields: fields{
				Code: etf.EmptyList,
			},
			want:    []byte{106},
			wantErr: false,
		},
		{
			name: "list",
			fields: fields{
				Code: etf.List,
				Vals: []etf.Element{
					{
						Code: etf.Binary,
						Val:  []byte{116, 101, 115, 116},
					},
					{
						Code: etf.Int8,
						Val:  []byte{1},
					},
				},
			},
			want:    []byte{108, 0, 0, 0, 2, 109, 0, 0, 0, 4, 116, 101, 115, 116, 97, 1, 106},
			wantErr: false,
		},
		{
			name: "atom",
			fields: fields{
				Code: etf.Atom,
				Val:  []byte{116, 101, 115, 116},
			},
			want:    []byte{100, 0, 4, 116, 101, 115, 116},
			wantErr: false,
		},
		{
			name: "string",
			fields: fields{
				Code: etf.String,
				Val:  []byte{116, 101, 115, 116},
			},
			want:    []byte{107, 0, 4, 116, 101, 115, 116},
			wantErr: false,
		},
		{
			name: "binary",
			fields: fields{
				Code: etf.Binary,
				Val:  []byte{116, 101, 115, 116},
			},
			want:    []byte{109, 0, 0, 0, 4, 116, 101, 115, 116},
			wantErr: false,
		},
		{
			name: "int8",
			fields: fields{
				Code: etf.Int8,
				Val:  []byte{100},
			},
			want:    []byte{97, 100},
			wantErr: false,
		},
		{
			name: "int8 oob",
			fields: fields{
				Code: etf.Int8,
				Val:  []byte{100, 1},
			},
			wantErr: true,
		},
		{
			name: "int32",
			fields: fields{
				Code: etf.Int32,
				Val:  []byte{1, 0, 0, 0},
			},
			want:    []byte{98, 1, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "int32 short",
			fields: fields{
				Code: etf.Int32,
				Val:  []byte{1, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "int32 long",
			fields: fields{
				Code: etf.Int32,
				Val:  []byte{1, 0, 0, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "smallbig",
			fields: fields{
				Code: etf.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0, 0},
			},
			want:    []byte{110, 8, 0, 0, 0, 0, 0, 0, 4, 0, 0},
			wantErr: false,
		},
		{
			name: "smallbig neg",
			fields: fields{
				Code: etf.SmallBig,
				Val:  []byte{1, 0, 0, 0, 0, 0, 4, 0, 0},
			},
			want:    []byte{110, 8, 1, 0, 0, 0, 0, 0, 4, 0, 0},
			wantErr: false,
		},
		{
			name: "smallbig short",
			fields: fields{
				Code: etf.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0},
			},
			wantErr: true,
		},
		{
			name: "smallbig long",
			fields: fields{
				Code: etf.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "largebig",
			fields: fields{
				Code: etf.LargeBig,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			got, err := e.Marshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Element.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_MarshalTo(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "map",
			fields: fields{
				Code: etf.Map,
				Vals: []etf.Element{
					{
						Code: etf.Binary,
						Val:  []byte{116, 101, 115, 116},
					},
					{
						Code: etf.Int8,
						Val:  []byte{1},
					},
				},
			},
			want:    []byte{116, 0, 0, 0, 1, 109, 0, 0, 0, 4, 116, 101, 115, 116, 97, 1},
			wantErr: false,
		},
		{
			name: "empty list",
			fields: fields{
				Code: etf.EmptyList,
			},
			want:    []byte{106},
			wantErr: false,
		},
		{
			name: "list",
			fields: fields{
				Code: etf.List,
				Vals: []etf.Element{
					{
						Code: etf.Binary,
						Val:  []byte{116, 101, 115, 116},
					},
					{
						Code: etf.Int8,
						Val:  []byte{1},
					},
				},
			},
			want:    []byte{108, 0, 0, 0, 2, 109, 0, 0, 0, 4, 116, 101, 115, 116, 97, 1, 106},
			wantErr: false,
		},
		{
			name: "atom",
			fields: fields{
				Code: etf.Atom,
				Val:  []byte{116, 101, 115, 116},
			},
			want:    []byte{100, 0, 4, 116, 101, 115, 116},
			wantErr: false,
		},
		{
			name: "string",
			fields: fields{
				Code: etf.String,
				Val:  []byte{116, 101, 115, 116},
			},
			want:    []byte{107, 0, 4, 116, 101, 115, 116},
			wantErr: false,
		},
		{
			name: "binary",
			fields: fields{
				Code: etf.Binary,
				Val:  []byte{116, 101, 115, 116},
			},
			want:    []byte{109, 0, 0, 0, 4, 116, 101, 115, 116},
			wantErr: false,
		},
		{
			name: "int8",
			fields: fields{
				Code: etf.Int8,
				Val:  []byte{100},
			},
			want:    []byte{97, 100},
			wantErr: false,
		},
		{
			name: "int8 oob",
			fields: fields{
				Code: etf.Int8,
				Val:  []byte{100, 1},
			},
			wantErr: true,
		},
		{
			name: "int32",
			fields: fields{
				Code: etf.Int32,
				Val:  []byte{1, 0, 0, 0},
			},
			want:    []byte{98, 1, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "int32 short",
			fields: fields{
				Code: etf.Int32,
				Val:  []byte{1, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "int32 long",
			fields: fields{
				Code: etf.Int32,
				Val:  []byte{1, 0, 0, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "smallbig",
			fields: fields{
				Code: etf.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0, 0},
			},
			want:    []byte{110, 8, 0, 0, 0, 0, 0, 0, 4, 0, 0},
			wantErr: false,
		},
		{
			name: "smallbig neg",
			fields: fields{
				Code: etf.SmallBig,
				Val:  []byte{1, 0, 0, 0, 0, 0, 4, 0, 0},
			},
			want:    []byte{110, 8, 1, 0, 0, 0, 0, 0, 4, 0, 0},
			wantErr: false,
		},
		{
			name: "smallbig short",
			fields: fields{
				Code: etf.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0},
			},
			wantErr: true,
		},
		{
			name: "smallbig long",
			fields: fields{
				Code: etf.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "largebig",
			fields: fields{
				Code: etf.LargeBig,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			b := &bytes.Buffer{}
			if err := e.MarshalTo(b); (err != nil) != tt.wantErr {
				t.Errorf("Element.MarshalTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if gotB := b.Bytes(); !bytes.Equal(gotB, tt.want) {
				t.Errorf("Element.MarshalTo() = %v, want %v", gotB, tt.want)
			}
		})
	}
}

func TestElement_ToString(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "atom",
			fields: fields{
				Code: etf.Atom,
				Val:  []byte("test"),
			},
			want: "test",
		},
		{
			name: "binary",
			fields: fields{
				Code: etf.Binary,
				Val:  []byte("test"),
			},
			want: "test",
		},
		{
			name: "string",
			fields: fields{
				Code: etf.String,
				Val:  []byte("test"),
			},
			want: "test",
		},
		{
			name: "bad",
			fields: fields{
				Code: etf.Map,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			got, err := e.ToString()
			if (err != nil) != tt.wantErr {
				t.Errorf("Element.ToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Element.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_ToBytes(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "atom",
			fields: fields{
				Code: etf.Atom,
				Val:  []byte("test"),
			},
			want: []byte("test"),
		},
		{
			name: "binary",
			fields: fields{
				Code: etf.Binary,
				Val:  []byte("test"),
			},
			want: []byte("test"),
		},
		{
			name: "string",
			fields: fields{
				Code: etf.String,
				Val:  []byte("test"),
			},
			want: []byte("test"),
		},
		{
			name: "bad",
			fields: fields{
				Code: etf.Map,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			got, err := e.ToBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Element.ToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_ToInt(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "int8",
			fields: fields{
				Code: etf.Int8,
				Val:  []byte{123},
			},
			want: 123,
		},
		{
			name: "int32",
			fields: fields{
				Code: etf.Int32,
				Val:  []byte{0, 1, 0, 1},
			},
			want: 65537,
		},
		{
			name: "smallbig",
			fields: fields{
				Code: etf.SmallBig,
				Val:  []byte{0, 1, 0, 1, 0, 0, 0, 0, 0},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			got, err := e.ToInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("Element.ToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Element.ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_ToInt64(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		{
			name: "int8",
			fields: fields{
				Code: etf.Int8,
				Val:  []byte{123},
			},
			want: 123,
		},
		{
			name: "int32",
			fields: fields{
				Code: etf.Int32,
				Val:  []byte{0, 1, 0, 1},
			},
			want: 65537,
		},
		{
			name: "smallbig",
			fields: fields{
				Code: etf.SmallBig,
				Val:  []byte{0, 1, 0, 1, 0, 0, 0, 0, 0},
			},
			want: 65537,
		},
		{
			name: "smallbig neg",
			fields: fields{
				Code: etf.SmallBig,
				Val:  []byte{1, 1, 0, 1, 0, 0, 0, 0, 0},
			},
			want: -65537,
		},
		{
			name: "map",
			fields: fields{
				Code: etf.Map,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			got, err := e.ToInt64()
			if (err != nil) != tt.wantErr {
				t.Errorf("Element.ToInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Element.ToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_ToMap(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]etf.Element
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Code: etf.Map,
				Vals: []etf.Element{
					{
						Code: etf.Atom,
						Val:  []byte("test"),
					},
					{
						Code: etf.Int32,
						Val:  []byte{0, 1, 0, 1},
					},
				},
			},
			want: map[string]etf.Element{
				"test": {
					Code: etf.Int32,
					Val:  []byte{0, 1, 0, 1},
				},
			},
			wantErr: false,
		},
		{
			name: "not map",
			fields: fields{
				Code: etf.List,
				Vals: []etf.Element{
					{
						Code: etf.Atom,
						Val:  []byte("test"),
					},
					{
						Code: etf.Int32,
						Val:  []byte{0, 1, 0, 1},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "bad parity",
			fields: fields{
				Code: etf.Map,
				Vals: []etf.Element{
					{
						Code: etf.Atom,
						Val:  []byte("test"),
					},
					{
						Code: etf.Int32,
						Val:  []byte{0, 1, 0, 1},
					},
					{
						Code: etf.Atom,
						Val:  []byte("test2"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "bad key",
			fields: fields{
				Code: etf.Map,
				Vals: []etf.Element{
					{
						Code: etf.Int32,
						Val:  []byte{0, 1, 0, 1},
					},
					{
						Code: etf.Atom,
						Val:  []byte("test2"),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			got, err := e.ToMap()
			if (err != nil) != tt.wantErr {
				t.Errorf("Element.ToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_ToList(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name    string
		fields  fields
		want    []etf.Element
		wantErr bool
	}{
		{
			name: "some",
			fields: fields{
				Code: etf.List,
				Vals: []etf.Element{
					{
						Code: etf.Atom,
						Val:  []byte("test"),
					},
					{
						Code: etf.Int32,
						Val:  []byte{0, 1, 0, 1},
					},
				},
			},
			want: []etf.Element{
				{
					Code: etf.Atom,
					Val:  []byte("test"),
				},
				{
					Code: etf.Int32,
					Val:  []byte{0, 1, 0, 1},
				},
			},
		},
		{
			name: "empty",
			fields: fields{
				Code: etf.EmptyList,
			},
			want: nil,
		},
		{
			name: "map",
			fields: fields{
				Code: etf.Map,
				Vals: []etf.Element{
					{
						Code: etf.Atom,
						Val:  []byte("test"),
					},
					{
						Code: etf.Int32,
						Val:  []byte{0, 1, 0, 1},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			got, err := e.ToList()
			if (err != nil) != tt.wantErr {
				t.Errorf("Element.ToList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.ToList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_IsNumeric(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "int8",
			fields: fields{
				Code: etf.Int8,
			},
			want: true,
		},
		{
			name: "int32",
			fields: fields{
				Code: etf.Int32,
			},
			want: true,
		},
		{
			name: "float",
			fields: fields{
				Code: etf.Float,
			},
			want: true,
		},
		{
			name: "smallbit",
			fields: fields{
				Code: etf.SmallBig,
			},
			want: true,
		},
		{
			name: "largebig",
			fields: fields{
				Code: etf.LargeBig,
			},
			want: true,
		},
		{
			name: "atom",
			fields: fields{
				Code: etf.Atom,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			if got := e.IsNumeric(); got != tt.want {
				t.Errorf("Element.IsNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_IsCollection(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "list",
			fields: fields{
				Code: etf.List,
			},
			want: true,
		},
		{
			name: "empty list",
			fields: fields{
				Code: etf.EmptyList,
			},
			want: true,
		},
		{
			name: "map",
			fields: fields{
				Code: etf.Map,
			},
			want: true,
		},
		{
			name: "atom",
			fields: fields{
				Code: etf.Atom,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			if got := e.IsCollection(); got != tt.want {
				t.Errorf("Element.IsCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_IsStringish(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "atom",
			fields: fields{
				Code: etf.Atom,
			},
			want: true,
		},
		{
			name: "binary",
			fields: fields{
				Code: etf.Binary,
			},
			want: true,
		},
		{
			name: "string",
			fields: fields{
				Code: etf.String,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			if got := e.IsStringish(); got != tt.want {
				t.Errorf("Element.IsStringish() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_IsList(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "list",
			fields: fields{
				Code: etf.List,
			},
			want: true,
		},
		{
			name: "empty list",
			fields: fields{
				Code: etf.EmptyList,
			},
			want: true,
		},
		{
			name: "map",
			fields: fields{
				Code: etf.Map,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			if got := e.IsList(); got != tt.want {
				t.Errorf("Element.IsList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_IsNil(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "ok",
			fields: fields{
				Code: etf.Atom,
				Val:  []byte("nil"),
			},
			want: true,
		},
		{
			name: "not atom",
			fields: fields{
				Code: etf.Binary,
				Val:  []byte("nil"),
			},
			want: false,
		},
		{
			name: "not false",
			fields: fields{
				Code: etf.Atom,
				Val:  []byte("foo"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			if got := e.IsNil(); got != tt.want {
				t.Errorf("Element.IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_IsTrue(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "ok",
			fields: fields{
				Code: etf.Atom,
				Val:  []byte("true"),
			},
			want: true,
		},
		{
			name: "not atom",
			fields: fields{
				Code: etf.Binary,
				Val:  []byte("true"),
			},
			want: false,
		},
		{
			name: "not false",
			fields: fields{
				Code: etf.Atom,
				Val:  []byte("foo"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			if got := e.IsTrue(); got != tt.want {
				t.Errorf("Element.IsTrue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_IsFalse(t *testing.T) {
	type fields struct {
		Code etf.Code
		Val  []byte
		Vals []etf.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "ok",
			fields: fields{
				Code: etf.Atom,
				Val:  []byte("false"),
			},
			want: true,
		},
		{
			name: "not atom",
			fields: fields{
				Code: etf.Binary,
				Val:  []byte("false"),
			},
			want: false,
		},
		{
			name: "not false",
			fields: fields{
				Code: etf.Atom,
				Val:  []byte("foo"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &etf.Element{
				Code: tt.fields.Code,
				Val:  tt.fields.Val,
				Vals: tt.fields.Vals,
			}
			if got := e.IsFalse(); got != tt.want {
				t.Errorf("Element.IsFalse() = %v, want %v", got, tt.want)
			}
		})
	}
}
