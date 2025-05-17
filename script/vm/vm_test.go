package vm

import (
	"reflect"
	"testing"

	"github.com/bjatkin/nook/script/ast"
	"github.com/bjatkin/nook/script/token"
)

func TestVM_Eval(t *testing.T) {
	type fields struct {
		scope *scope
	}
	type args struct {
		expr ast.Expr
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		{
			name:   "add integers",
			fields: fields{},
			args: args{expr: ast.SExpr{
				Operator: ast.Identifier{Tok: token.Token{Pos: 1, Value: "+", Kind: token.Plus}, Name: "+"},
				Operands: []ast.Expr{
					ast.Int{Tok: token.Token{Pos: 3, Value: "5", Kind: token.Int}, Value: 5},
					ast.Int{Tok: token.Token{Pos: 5, Value: "3", Kind: token.Int}, Value: 3},
				},
			}},
			want:    Value{value: int64(8), kind: Int},
			wantErr: false,
		},
		{
			name:   "add floats",
			fields: fields{},
			args: args{expr: ast.SExpr{
				Operator: ast.Identifier{Tok: token.Token{Pos: 1, Value: "+", Kind: token.Plus}, Name: "+"},
				Operands: []ast.Expr{
					ast.Float{Tok: token.Token{Pos: 3, Value: "1.43", Kind: token.Float}, Value: 1.43},
					ast.Float{Tok: token.Token{Pos: 5, Value: "9.35", Kind: token.Float}, Value: 9.35},
				},
			}},
			want:    Value{value: 10.78, kind: Float},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := &VM{
				scope: tt.fields.scope,
			}
			got, err := vm.Eval(tt.args.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("VM.Eval() err %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VM.Eval() = %#v(%T), want %#v(%T)", got, got, tt.want, tt.want)
			}
		})
	}
}
