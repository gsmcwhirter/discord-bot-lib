package cmdhandler

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gsmcwhirter/discord-bot-lib/v21/snowflake"
)

func TestUserMentionString(t *testing.T) {
	type args struct {
		uid snowflake.Snowflake
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic test",
			args: args{
				uid: snowflake.Snowflake(1234),
			},
			want: "<@!1234>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserMentionString(tt.args.uid); got != tt.want {
				t.Errorf("UserMentionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUserMention(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "yes",
			args: args{
				v: "<@!1234>",
			},
			want: true,
		},
		{
			name: "no",
			args: args{
				v: "<?!1234>",
			},
		},
		{
			name: "degenerate",
			args: args{
				v: "<@!>",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUserMention(tt.args.v); got != tt.want {
				t.Errorf("IsUserMention() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForceUserNicknameMention(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "nothing to do",
			args: args{
				v: "<@!1234>",
			},
			want:    "<@!1234>",
			wantErr: false,
		},
		{
			name: "do convert",
			args: args{
				v: "<@1234>",
			},
			want:    "<@!1234>",
			wantErr: false,
		},
		{
			name: "not a mention",
			args: args{
				v: "<?!1234>",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ForceUserNicknameMention(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ForceUserNicknameMention() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ForceUserNicknameMention() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForceUserAccountMention(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "nothing to do",
			args: args{
				v: "<@1234>",
			},
			want:    "<@1234>",
			wantErr: false,
		},
		{
			name: "do convert",
			args: args{
				v: "<@!1234>",
			},
			want:    "<@1234>",
			wantErr: false,
		},
		{
			name: "not a mention",
			args: args{
				v: "<?!1234>",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ForceUserAccountMention(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ForceUserAccountMention() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ForceUserAccountMention() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelMentionString(t *testing.T) {
	type args struct {
		cid snowflake.Snowflake
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				cid: snowflake.Snowflake(5678),
			},
			want: "<#5678>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ChannelMentionString(tt.args.cid); got != tt.want {
				t.Errorf("ChannelMentionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsChannelMention(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "is",
			args: args{
				v: "<#5678>",
			},
			want: true,
		},
		{
			name: "not",
			args: args{
				v: "<?5678>",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsChannelMention(tt.args.v); got != tt.want {
				t.Errorf("IsChannelMention() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleMentionString(t *testing.T) {
	type args struct {
		rid snowflake.Snowflake
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				rid: snowflake.Snowflake(9012),
			},
			want: "<@&9012>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoleMentionString(tt.args.rid); got != tt.want {
				t.Errorf("RoleMentionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRoleMention(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "is",
			args: args{
				v: "<@&9012>",
			},
			want: true,
		},
		{
			name: "not",
			args: args{
				v: "<?9012>",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRoleMention(tt.args.v); got != tt.want {
				t.Errorf("IsRoleMention() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_textSplit(t *testing.T) {
	type args struct {
		text   string
		target int
		delim  string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "by lines",
			args: args{
				text:   "line1\nline2\nline3",
				target: 6,
				delim:  "\n",
			},
			want: []string{"line1\n", "line2\n", "line3"},
		},
		{
			name: "two lines",
			args: args{
				text:   "line1\nline2\nline3",
				target: 13,
				delim:  "\n",
			},
			want: []string{"line1\nline2\n", "line3"},
		},
		{
			name: "by words",
			args: args{
				text:   "line 1\nline 2\nline 3",
				target: 4,
				delim:  "\n",
			},
			want: []string{"line", "1\n", "line", "2\n", "line", "3"},
		},
		{
			name: "long words",
			args: args{
				text:   "lineline 1\nlineline 2\nlineline 3",
				target: 4,
				delim:  "\n",
			},
			want: []string{"line", "line", "1\n", "line", "line", "2\n", "line", "line", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("# %s\n", tt.name)
			if got := textSplit(tt.args.text, tt.args.target, tt.args.delim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("textSplit() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
