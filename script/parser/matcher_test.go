package parser

import (
	"reflect"
	"testing"

	"github.com/bjatkin/nook/script/token"
)

// TODO: move the other matcher tests here from the lexer_test.go file

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

func Test_matchPath(t *testing.T) {
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
			name: "current dir",
			args: args{bytes: []byte(".")},
			want: &match{
				len:  1,
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
			if got := matchPath(tt.args.bytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
