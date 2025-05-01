package parser

import (
	"reflect"
	"testing"

	"github.com/bjatkin/nook/script/ast"
	"github.com/bjatkin/nook/script/token"
)

func TestParser_Parse(t *testing.T) {
	type fields struct {
		lexer lexer
	}
	tests := []struct {
		name   string
		fields fields
		want   ast.Expr
	}{
		{
			name:   "simple script",
			fields: fields{lexer: newLexer([]byte("(+ 5 10)"))},
			want: ast.SExpr{
				Operator: token.Token{Pos: 1, Value: "+", Kind: token.Plus},
				Operands: []ast.Expr{
					ast.Int{Tok: token.Token{Pos: 3, Value: "5", Kind: token.Int}, Value: 5},
					ast.Int{Tok: token.Token{Pos: 5, Value: "10", Kind: token.Int}, Value: 10},
				},
			},
		},
		{
			name:   "assignement",
			fields: fields{lexer: newLexer([]byte("(let a (- 8 0xFF))"))},
			want: ast.SExpr{
				Operator: token.Token{Pos: 1, Value: "let", Kind: token.Let},
				Operands: []ast.Expr{
					ast.Identifier{Value: token.Token{Pos: 5, Value: "a", Kind: token.Identifier}},
					ast.SExpr{
						Operator: token.Token{Pos: 8, Value: "-", Kind: token.Minus},
						Operands: []ast.Expr{
							ast.Int{Tok: token.Token{Pos: 10, Value: "8", Kind: token.Int}, Value: 8},
							ast.Int{Tok: token.Token{Pos: 12, Value: "0xFF", Kind: token.Int}, Value: 0xFF},
						},
					},
				},
			},
		},
		{
			name:   "add floats",
			fields: fields{lexer: newLexer([]byte("(+ 1.2 3.5 2.5)"))},
			want: ast.SExpr{
				Operator: token.Token{Pos: 1, Value: "+", Kind: token.Plus},
				Operands: []ast.Expr{
					ast.Float{Tok: token.Token{Pos: 3, Value: "1.2", Kind: token.Float}, Value: 1.2},
					ast.Float{Tok: token.Token{Pos: 7, Value: "3.5", Kind: token.Float}, Value: 3.5},
					ast.Float{Tok: token.Token{Pos: 11, Value: "2.5", Kind: token.Float}, Value: 2.5},
				},
			},
		},
		{
			name:   "run command",
			fields: fields{lexer: newLexer([]byte("($git 'status)"))},
			want: ast.SExpr{
				Operator: token.Token{Pos: 1, Value: "$git", Kind: token.Command},
				Operands: []ast.Expr{
					ast.Atom{Tok: token.Token{Pos: 6, Value: "'status", Kind: token.Atom}, Value: "'status"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				lexer: tt.fields.lexer,
			}
			if got := p.Parse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.Parse() = %v\nwant %v", got, tt.want)
			}
		})
	}
}
