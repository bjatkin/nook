package parser

import (
	"reflect"
	"testing"

	"github.com/bjatkin/nook/script/token"
)

func Test_matchFlag(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want *match
	}{
		{
			name: "short flag",
			args: args{bytes: []byte("-m  ")},
			want: &match{
				len:  2,
				kind: token.Flag,
			},
		},
		{
			name: "long flag",
			args: args{bytes: []byte("--version  ")},
			want: &match{
				len:  9,
				kind: token.Flag,
			},
		},
		{
			name: "invalid flag",
			args: args{bytes: []byte("-")},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchFlag(tt.args.bytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchLongPath(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want *match
	}{
		{
			name: "valid path",
			args: args{bytes: []byte("./nook")},
			want: &match{
				len:  6,
				kind: token.Path,
			},
		},
		{
			name: "backwards dir",
			args: args{bytes: []byte("../nested")},
			want: &match{
				len:  9,
				kind: token.Path,
			},
		},
		{
			name: "invalid directory",
			args: args{bytes: []byte("test/this/path")},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchLongPath(tt.args.bytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchAtom(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want *match
	}{
		{
			name: "valid atom",
			args: args{bytes: []byte("'validAtom")},
			want: &match{len: 10, kind: token.Atom},
		},
		{
			name: "short atom",
			args: args{bytes: []byte("'ok")},
			want: &match{len: 3, kind: token.Atom},
		},
		{
			name: "invalid atom",
			args: args{bytes: []byte("test")},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchAtom(tt.args.bytes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchAtom() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchString(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want *match
	}{
		{
			name: "valid string",
			args: args{bytes: []byte("\"hello world\"")},
			want: &match{len: 13, kind: token.String},
		},
		{
			name: "invalid string",
			args: args{bytes: []byte("test")},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchString(tt.args.bytes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchString() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_lexer_matchIdentifier(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want *match
	}{
		{
			name: "invalid identifier",
			args: args{bytes: []byte("10")},
			want: nil,
		},
		{
			name: "valid identifier",
			args: args{bytes: []byte("test10")},
			want: &match{len: 6, kind: token.Identifier},
		},
		{
			name: "trailing space",
			args: args{bytes: []byte("a ")},
			want: &match{len: 1, kind: token.Identifier},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchIdentifier(tt.args.bytes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lexer.matchIdentifier() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchHex(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantLen uint
		wantOk  bool
	}{
		{
			name:    "single hex number",
			args:    args{bytes: []byte("A")},
			wantLen: 1,
			wantOk:  true,
		},
		{
			name:    "long hex number",
			args:    args{bytes: []byte("01_CF")},
			wantLen: 5,
			wantOk:  true,
		},
		{
			name:    "not a hex number",
			args:    args{bytes: []byte("Hello")},
			wantLen: 0,
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchHex(tt.args.bytes)
			if got != tt.wantLen {
				t.Errorf("matchHex() got = %v, wantLen %v", got, tt.wantLen)
			}
			if got1 != tt.wantOk {
				t.Errorf("matchHex() got1 = %v, wantOk %v", got1, tt.wantOk)
			}
		})
	}
}

func Test_matchOctal(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantLen uint
		wantOk  bool
	}{
		{
			name:    "octal number",
			args:    args{bytes: []byte("23_02")},
			wantLen: 5,
			wantOk:  true,
		},
		{
			name:    "not octal number",
			args:    args{bytes: []byte("Hello")},
			wantLen: 0,
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchOctal(tt.args.bytes)
			if got != tt.wantLen {
				t.Errorf("matchOctal() got = %v, want %v", got, tt.wantLen)
			}
			if got1 != tt.wantOk {
				t.Errorf("matchOctal() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}

func Test_matchBinary(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantLen uint
		wantOk  bool
	}{
		{
			name:    "binary number",
			args:    args{bytes: []byte("01_01")},
			wantLen: 5,
			wantOk:  true,
		},
		{
			name:    "not octal number",
			args:    args{bytes: []byte("Hello")},
			wantLen: 0,
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchBinary(tt.args.bytes)
			if got != tt.wantLen {
				t.Errorf("matchBinary() got = %v, want %v", got, tt.wantLen)
			}
			if got1 != tt.wantOk {
				t.Errorf("matchBinary() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}

func Test_matchNumber(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want *match
	}{
		{
			name: "decimal number",
			args: args{bytes: []byte("012_456")},
			want: &match{len: 7, kind: token.Int},
		},
		{
			name: "hex number",
			args: args{bytes: []byte("0xAB_C0")},
			want: &match{len: 7, kind: token.Int},
		},
		{
			name: "octal number",
			args: args{bytes: []byte("0o17_46")},
			want: &match{len: 7, kind: token.Int},
		},
		{
			name: "binary number",
			args: args{bytes: []byte("0b0110_1111")},
			want: &match{len: 11, kind: token.Int},
		},
		{
			name: "keyword",
			args: args{bytes: []byte("else")},
			want: nil,
		},
		{
			name: "invalid number",
			args: args{bytes: []byte("10)")},
			want: &match{len: 2, kind: token.Int},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchInt(tt.args.bytes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchNumber() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchFloat(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want *match
	}{
		{
			name: "float number",
			args: args{bytes: []byte("43.54")},
			want: &match{len: 5, kind: token.Float},
		},
		{
			name: "int number",
			args: args{bytes: []byte("1234")},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchFloat(tt.args.bytes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchFloat() got = %v, want %v", got, tt.want)
			}
		})
	}
}
