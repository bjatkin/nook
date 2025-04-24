package ast

import "github.com/bjatkin/nook/script/token"

type Expr interface {
	expr()
}

type SExpr struct {
	Expr
	Operator token.Token
	Operands []Expr
}

type Int struct {
	Expr
	Tok   token.Token
	Value int64
}

type Float struct {
	Expr
	Tok   token.Token
	Value float64
}

type String struct {
	Expr
	Tok   token.Token
	Value string
}

type Atom struct {
	Expr
	Tok   token.Token
	Value string
}

type Bool struct {
	Expr
	Tok   token.Token
	Value bool
}

type Identifier struct {
	Expr
	Value token.Token
}
