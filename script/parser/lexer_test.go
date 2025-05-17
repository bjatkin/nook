package parser

import (
	"reflect"
	"testing"

	"github.com/bjatkin/nook/script/token"
)

func Test_lexer_lex(t *testing.T) {
	type fields struct {
		source               []byte
		pos                  uint
		includeIgnoredTokens bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []token.Token
	}{
		{
			name: "skip white space",
			fields: fields{
				source:               []byte("(+ 1 2 3)"),
				pos:                  0,
				includeIgnoredTokens: false,
			},
			want: []token.Token{
				{Pos: 0, Value: "(", Kind: token.OpenParen},
				{Pos: 1, Value: "+", Kind: token.Plus},
				{Pos: 3, Value: "1", Kind: token.Int},
				{Pos: 5, Value: "2", Kind: token.Int},
				{Pos: 7, Value: "3", Kind: token.Int},
				{Pos: 8, Value: ")", Kind: token.CloseParen},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				source:               tt.fields.source,
				pos:                  tt.fields.pos,
				includeIgnoredTokens: tt.fields.includeIgnoredTokens,
			}
			if got := l.Lex(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lexer.lex() = %v, want %v", got, tt.want)
			}
		})
	}
}
