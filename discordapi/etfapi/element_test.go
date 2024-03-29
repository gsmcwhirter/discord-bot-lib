package etfapi_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/etfapi"
)

func TestNewCollectionElement(t *testing.T) {
	t.Parallel()

	type args struct {
		code etfapi.Code
		val  []etfapi.Element
	}
	tests := []struct {
		name    string
		args    args
		wantE   etfapi.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				code: etfapi.List,
				val: []etfapi.Element{
					{
						Code: etfapi.EmptyList,
					},
				},
			},
			wantE: etfapi.Element{
				Code: etfapi.List,
				Vals: []etfapi.Element{
					{
						Code: etfapi.EmptyList,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "not ok",
			args: args{
				code: etfapi.Int8,
				val: []etfapi.Element{
					{
						Code: etfapi.EmptyList,
					},
				},
			},
			wantE:   etfapi.Element{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotE, err := etfapi.NewCollectionElement(tt.args.code, tt.args.val)
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
	t.Parallel()

	type args struct {
		code etfapi.Code
		val  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantE   etfapi.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				code: etfapi.Int8,
				val:  []byte{1},
			},
			wantE: etfapi.Element{
				Code: etfapi.Int8,
				Val:  []byte{1},
			},
			wantErr: false,
		},
		{
			name: "not ok",
			args: args{
				code: etfapi.List,
				val:  []byte{2},
			},
			wantE:   etfapi.Element{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotE, err := etfapi.NewBasicElement(tt.args.code, tt.args.val)
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
	t.Parallel()

	tests := []struct {
		name    string
		wantE   etfapi.Element
		wantErr bool
	}{
		{
			name: "ok",
			wantE: etfapi.Element{
				Code: etfapi.Atom,
				Val:  []byte("nil"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotE, err := etfapi.NewNilElement()
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
	t.Parallel()

	type args struct {
		val bool
	}
	tests := []struct {
		name    string
		args    args
		wantE   etfapi.Element
		wantErr bool
	}{
		{
			name: "ok true",
			args: args{true},
			wantE: etfapi.Element{
				Code: etfapi.Atom,
				Val:  []byte("true"),
			},
		},
		{
			name: "ok false",
			args: args{false},
			wantE: etfapi.Element{
				Code: etfapi.Atom,
				Val:  []byte("false"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotE, err := etfapi.NewBoolElement(tt.args.val)
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
	t.Parallel()

	type args struct {
		val int
	}
	tests := []struct {
		name    string
		args    args
		wantE   etfapi.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{1},
			wantE: etfapi.Element{
				Code: etfapi.Int8,
				Val:  []byte{1},
			},
		},
		{
			name:    "not ok",
			args:    args{1024},
			wantE:   etfapi.Element{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotE, err := etfapi.NewInt8Element(tt.args.val)
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
	t.Parallel()

	type args struct {
		val int
	}
	tests := []struct {
		name    string
		args    args
		wantE   etfapi.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{1 << 31},
			wantE: etfapi.Element{
				Code: etfapi.Int32,
				Val:  []byte{128, 0, 0, 0},
			},
		},
		{
			name:    "not ok",
			args:    args{1 << 42},
			wantE:   etfapi.Element{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotE, err := etfapi.NewInt32Element(tt.args.val)
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
	t.Parallel()

	type args struct {
		val int64
	}
	tests := []struct {
		name    string
		args    args
		wantE   etfapi.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{1 << 42},
			wantE: etfapi.Element{
				Code: etfapi.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0, 0},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotE, err := etfapi.NewSmallBigElement(tt.args.val)
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
	t.Parallel()

	type args struct {
		val []byte
	}
	tests := []struct {
		name    string
		args    args
		wantE   etfapi.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{[]byte("test")},
			wantE: etfapi.Element{
				Code: etfapi.Binary,
				Val:  []byte{116, 101, 115, 116},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotE, err := etfapi.NewBinaryElement(tt.args.val)
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
	t.Parallel()

	type args struct {
		val []byte
	}
	tests := []struct {
		name    string
		args    args
		wantE   etfapi.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{[]byte("test")},
			wantE: etfapi.Element{
				Code: etfapi.Atom,
				Val:  []byte{116, 101, 115, 116},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotE, err := etfapi.NewAtomElement(tt.args.val)
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
	t.Parallel()

	type args struct {
		val string
	}
	tests := []struct {
		name    string
		args    args
		wantE   etfapi.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{"test"},
			wantE: etfapi.Element{
				Code: etfapi.Binary,
				Val:  []byte{116, 101, 115, 116},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotE, err := etfapi.NewStringElement(tt.args.val)
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
	t.Parallel()

	type args struct {
		val map[string]etfapi.Element
	}
	tests := []struct {
		name    string
		args    args
		wantE   etfapi.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{map[string]etfapi.Element{
				"test": {Code: etfapi.Int8, Val: []byte{1}},
			}},
			wantE: etfapi.Element{
				Code: etfapi.Map,
				Vals: []etfapi.Element{
					{
						Code: etfapi.Binary,
						Val:  []byte{116, 101, 115, 116},
					},
					{
						Code: etfapi.Int8,
						Val:  []byte{1},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotE, err := etfapi.NewMapElement(tt.args.val)
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
	t.Parallel()

	type args struct {
		val []etfapi.Element
	}
	tests := []struct {
		name    string
		args    args
		wantE   etfapi.Element
		wantErr bool
	}{
		{
			name: "ok",
			args: args{[]etfapi.Element{
				{
					Code: etfapi.Int8,
					Val:  []byte{2},
				},
				{
					Code: etfapi.Int8,
					Val:  []byte{1},
				},
			}},
			wantE: etfapi.Element{
				Code: etfapi.List,
				Vals: []etfapi.Element{
					{
						Code: etfapi.Int8,
						Val:  []byte{2},
					},
					{
						Code: etfapi.Int8,
						Val:  []byte{1},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotE, err := etfapi.NewListElement(tt.args.val)
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
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
				Code: etfapi.Map,
				Vals: []etfapi.Element{
					{
						Code: etfapi.Binary,
						Val:  []byte{116, 101, 115, 116},
					},
					{
						Code: etfapi.Int8,
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
				Code: etfapi.EmptyList,
			},
			want:    []byte{106},
			wantErr: false,
		},
		{
			name: "list",
			fields: fields{
				Code: etfapi.List,
				Vals: []etfapi.Element{
					{
						Code: etfapi.Binary,
						Val:  []byte{116, 101, 115, 116},
					},
					{
						Code: etfapi.Int8,
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
				Code: etfapi.Atom,
				Val:  []byte{116, 101, 115, 116},
			},
			want:    []byte{100, 0, 4, 116, 101, 115, 116},
			wantErr: false,
		},
		{
			name: "string",
			fields: fields{
				Code: etfapi.String,
				Val:  []byte{116, 101, 115, 116},
			},
			want:    []byte{107, 0, 4, 116, 101, 115, 116},
			wantErr: false,
		},
		{
			name: "binary",
			fields: fields{
				Code: etfapi.Binary,
				Val:  []byte{116, 101, 115, 116},
			},
			want:    []byte{109, 0, 0, 0, 4, 116, 101, 115, 116},
			wantErr: false,
		},
		{
			name: "int8",
			fields: fields{
				Code: etfapi.Int8,
				Val:  []byte{100},
			},
			want:    []byte{97, 100},
			wantErr: false,
		},
		{
			name: "int8 oob",
			fields: fields{
				Code: etfapi.Int8,
				Val:  []byte{100, 1},
			},
			wantErr: true,
		},
		{
			name: "int32",
			fields: fields{
				Code: etfapi.Int32,
				Val:  []byte{1, 0, 0, 0},
			},
			want:    []byte{98, 1, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "int32 short",
			fields: fields{
				Code: etfapi.Int32,
				Val:  []byte{1, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "int32 long",
			fields: fields{
				Code: etfapi.Int32,
				Val:  []byte{1, 0, 0, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "smallbig",
			fields: fields{
				Code: etfapi.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0, 0},
			},
			want:    []byte{110, 8, 0, 0, 0, 0, 0, 0, 4, 0, 0},
			wantErr: false,
		},
		{
			name: "smallbig neg",
			fields: fields{
				Code: etfapi.SmallBig,
				Val:  []byte{1, 0, 0, 0, 0, 0, 4, 0, 0},
			},
			want:    []byte{110, 8, 1, 0, 0, 0, 0, 0, 4, 0, 0},
			wantErr: false,
		},
		{
			name: "smallbig short",
			fields: fields{
				Code: etfapi.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0},
			},
			wantErr: true,
		},
		{
			name: "smallbig long",
			fields: fields{
				Code: etfapi.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "largebig",
			fields: fields{
				Code: etfapi.LargeBig,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
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
				Code: etfapi.Map,
				Vals: []etfapi.Element{
					{
						Code: etfapi.Binary,
						Val:  []byte{116, 101, 115, 116},
					},
					{
						Code: etfapi.Int8,
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
				Code: etfapi.EmptyList,
			},
			want:    []byte{106},
			wantErr: false,
		},
		{
			name: "list",
			fields: fields{
				Code: etfapi.List,
				Vals: []etfapi.Element{
					{
						Code: etfapi.Binary,
						Val:  []byte{116, 101, 115, 116},
					},
					{
						Code: etfapi.Int8,
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
				Code: etfapi.Atom,
				Val:  []byte{116, 101, 115, 116},
			},
			want:    []byte{100, 0, 4, 116, 101, 115, 116},
			wantErr: false,
		},
		{
			name: "string",
			fields: fields{
				Code: etfapi.String,
				Val:  []byte{116, 101, 115, 116},
			},
			want:    []byte{107, 0, 4, 116, 101, 115, 116},
			wantErr: false,
		},
		{
			name: "binary",
			fields: fields{
				Code: etfapi.Binary,
				Val:  []byte{116, 101, 115, 116},
			},
			want:    []byte{109, 0, 0, 0, 4, 116, 101, 115, 116},
			wantErr: false,
		},
		{
			name: "int8",
			fields: fields{
				Code: etfapi.Int8,
				Val:  []byte{100},
			},
			want:    []byte{97, 100},
			wantErr: false,
		},
		{
			name: "int8 oob",
			fields: fields{
				Code: etfapi.Int8,
				Val:  []byte{100, 1},
			},
			wantErr: true,
		},
		{
			name: "int32",
			fields: fields{
				Code: etfapi.Int32,
				Val:  []byte{1, 0, 0, 0},
			},
			want:    []byte{98, 1, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "int32 short",
			fields: fields{
				Code: etfapi.Int32,
				Val:  []byte{1, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "int32 long",
			fields: fields{
				Code: etfapi.Int32,
				Val:  []byte{1, 0, 0, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "smallbig",
			fields: fields{
				Code: etfapi.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0, 0},
			},
			want:    []byte{110, 8, 0, 0, 0, 0, 0, 0, 4, 0, 0},
			wantErr: false,
		},
		{
			name: "smallbig neg",
			fields: fields{
				Code: etfapi.SmallBig,
				Val:  []byte{1, 0, 0, 0, 0, 0, 4, 0, 0},
			},
			want:    []byte{110, 8, 1, 0, 0, 0, 0, 0, 4, 0, 0},
			wantErr: false,
		},
		{
			name: "smallbig short",
			fields: fields{
				Code: etfapi.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0},
			},
			wantErr: true,
		},
		{
			name: "smallbig long",
			fields: fields{
				Code: etfapi.SmallBig,
				Val:  []byte{0, 0, 0, 0, 0, 0, 4, 0, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "largebig",
			fields: fields{
				Code: etfapi.LargeBig,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
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
				Code: etfapi.Atom,
				Val:  []byte("test"),
			},
			want: "test",
		},
		{
			name: "binary",
			fields: fields{
				Code: etfapi.Binary,
				Val:  []byte("test"),
			},
			want: "test",
		},
		{
			name: "string",
			fields: fields{
				Code: etfapi.String,
				Val:  []byte("test"),
			},
			want: "test",
		},
		{
			name: "bad",
			fields: fields{
				Code: etfapi.Map,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
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
				Code: etfapi.Atom,
				Val:  []byte("test"),
			},
			want: []byte("test"),
		},
		{
			name: "binary",
			fields: fields{
				Code: etfapi.Binary,
				Val:  []byte("test"),
			},
			want: []byte("test"),
		},
		{
			name: "string",
			fields: fields{
				Code: etfapi.String,
				Val:  []byte("test"),
			},
			want: []byte("test"),
		},
		{
			name: "bad",
			fields: fields{
				Code: etfapi.Map,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
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
				Code: etfapi.Int8,
				Val:  []byte{123},
			},
			want: 123,
		},
		{
			name: "int32",
			fields: fields{
				Code: etfapi.Int32,
				Val:  []byte{0, 1, 0, 1},
			},
			want: 65537,
		},
		{
			name: "smallbig",
			fields: fields{
				Code: etfapi.SmallBig,
				Val:  []byte{0, 1, 0, 1, 0, 0, 0, 0, 0},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
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
				Code: etfapi.Int8,
				Val:  []byte{123},
			},
			want: 123,
		},
		{
			name: "int32",
			fields: fields{
				Code: etfapi.Int32,
				Val:  []byte{0, 1, 0, 1},
			},
			want: 65537,
		},
		{
			name: "smallbig",
			fields: fields{
				Code: etfapi.SmallBig,
				Val:  []byte{0, 1, 0, 1, 0, 0, 0, 0, 0},
			},
			want: 65537,
		},
		{
			name: "smallbig neg",
			fields: fields{
				Code: etfapi.SmallBig,
				Val:  []byte{1, 1, 0, 1, 0, 0, 0, 0, 0},
			},
			want: -65537,
		},
		{
			name: "map",
			fields: fields{
				Code: etfapi.Map,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]etfapi.Element
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Code: etfapi.Map,
				Vals: []etfapi.Element{
					{
						Code: etfapi.Atom,
						Val:  []byte("test"),
					},
					{
						Code: etfapi.Int32,
						Val:  []byte{0, 1, 0, 1},
					},
				},
			},
			want: map[string]etfapi.Element{
				"test": {
					Code: etfapi.Int32,
					Val:  []byte{0, 1, 0, 1},
				},
			},
			wantErr: false,
		},
		{
			name: "not map",
			fields: fields{
				Code: etfapi.List,
				Vals: []etfapi.Element{
					{
						Code: etfapi.Atom,
						Val:  []byte("test"),
					},
					{
						Code: etfapi.Int32,
						Val:  []byte{0, 1, 0, 1},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "bad parity",
			fields: fields{
				Code: etfapi.Map,
				Vals: []etfapi.Element{
					{
						Code: etfapi.Atom,
						Val:  []byte("test"),
					},
					{
						Code: etfapi.Int32,
						Val:  []byte{0, 1, 0, 1},
					},
					{
						Code: etfapi.Atom,
						Val:  []byte("test2"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "bad key",
			fields: fields{
				Code: etfapi.Map,
				Vals: []etfapi.Element{
					{
						Code: etfapi.Int32,
						Val:  []byte{0, 1, 0, 1},
					},
					{
						Code: etfapi.Atom,
						Val:  []byte("test2"),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
	}
	tests := []struct {
		name    string
		fields  fields
		want    []etfapi.Element
		wantErr bool
	}{
		{
			name: "some",
			fields: fields{
				Code: etfapi.List,
				Vals: []etfapi.Element{
					{
						Code: etfapi.Atom,
						Val:  []byte("test"),
					},
					{
						Code: etfapi.Int32,
						Val:  []byte{0, 1, 0, 1},
					},
				},
			},
			want: []etfapi.Element{
				{
					Code: etfapi.Atom,
					Val:  []byte("test"),
				},
				{
					Code: etfapi.Int32,
					Val:  []byte{0, 1, 0, 1},
				},
			},
		},
		{
			name: "empty",
			fields: fields{
				Code: etfapi.EmptyList,
			},
			want: nil,
		},
		{
			name: "map",
			fields: fields{
				Code: etfapi.Map,
				Vals: []etfapi.Element{
					{
						Code: etfapi.Atom,
						Val:  []byte("test"),
					},
					{
						Code: etfapi.Int32,
						Val:  []byte{0, 1, 0, 1},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "int8",
			fields: fields{
				Code: etfapi.Int8,
			},
			want: true,
		},
		{
			name: "int32",
			fields: fields{
				Code: etfapi.Int32,
			},
			want: true,
		},
		{
			name: "float",
			fields: fields{
				Code: etfapi.Float,
			},
			want: true,
		},
		{
			name: "smallbit",
			fields: fields{
				Code: etfapi.SmallBig,
			},
			want: true,
		},
		{
			name: "largebig",
			fields: fields{
				Code: etfapi.LargeBig,
			},
			want: true,
		},
		{
			name: "atom",
			fields: fields{
				Code: etfapi.Atom,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "list",
			fields: fields{
				Code: etfapi.List,
			},
			want: true,
		},
		{
			name: "empty list",
			fields: fields{
				Code: etfapi.EmptyList,
			},
			want: true,
		},
		{
			name: "map",
			fields: fields{
				Code: etfapi.Map,
			},
			want: true,
		},
		{
			name: "atom",
			fields: fields{
				Code: etfapi.Atom,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "atom",
			fields: fields{
				Code: etfapi.Atom,
			},
			want: true,
		},
		{
			name: "binary",
			fields: fields{
				Code: etfapi.Binary,
			},
			want: true,
		},
		{
			name: "string",
			fields: fields{
				Code: etfapi.String,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "list",
			fields: fields{
				Code: etfapi.List,
			},
			want: true,
		},
		{
			name: "empty list",
			fields: fields{
				Code: etfapi.EmptyList,
			},
			want: true,
		},
		{
			name: "map",
			fields: fields{
				Code: etfapi.Map,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "ok",
			fields: fields{
				Code: etfapi.Atom,
				Val:  []byte("nil"),
			},
			want: true,
		},
		{
			name: "not atom",
			fields: fields{
				Code: etfapi.Binary,
				Val:  []byte("nil"),
			},
			want: false,
		},
		{
			name: "not false",
			fields: fields{
				Code: etfapi.Atom,
				Val:  []byte("foo"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "ok",
			fields: fields{
				Code: etfapi.Atom,
				Val:  []byte("true"),
			},
			want: true,
		},
		{
			name: "not atom",
			fields: fields{
				Code: etfapi.Binary,
				Val:  []byte("true"),
			},
			want: false,
		},
		{
			name: "not false",
			fields: fields{
				Code: etfapi.Atom,
				Val:  []byte("foo"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
	t.Parallel()

	type fields struct {
		Code etfapi.Code
		Val  []byte
		Vals []etfapi.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "ok",
			fields: fields{
				Code: etfapi.Atom,
				Val:  []byte("false"),
			},
			want: true,
		},
		{
			name: "not atom",
			fields: fields{
				Code: etfapi.Binary,
				Val:  []byte("false"),
			},
			want: false,
		},
		{
			name: "not false",
			fields: fields{
				Code: etfapi.Atom,
				Val:  []byte("foo"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &etfapi.Element{
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
